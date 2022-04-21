package kube

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"os/user"
	"path/filepath"
	"rexagent/pkg/log"
	"rexagent/pkg/task"
	sigyaml "sigs.k8s.io/yaml"
)

type KubeExecutor struct {
	id   string
	name string

	confPath  string
	namespace string
}

func (e *KubeExecutor) Execute(job *task.KubeJob) (retbytes []byte, err error) {
	e.namespace = job.Namespace
	e.confPath = job.KubeConfPath
	if e.confPath == "" {
		user, err := user.Current()
		if nil == err {
			homeDir := user.HomeDir
			e.confPath = filepath.Join(homeDir, ".kube", "config")
		} else {
			return nil, errors.New("can't find the config file")
		}
	}

	if e.namespace == "" {
		e.namespace = "default"
	}

	var ret *unstructured.Unstructured

	ret, err = e.Apply(job.YamlStr)
	log.Info("apply kube yaml config.")
	retbytes, err = sigyaml.Marshal(ret)
	return
}

func (d *KubeExecutor) Apply(filestr string) (*unstructured.Unstructured, error) {
	var ret *unstructured.Unstructured

	// 获取配置文件
	config, err := clientcmd.BuildConfigFromFlags("", d.confPath)
	if err != nil {
		return nil, err
	}
	// 获取kube api的客户端
	dynameicclient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// 解析yaml
	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(filestr), 4096)
	for {
		var rawObj runtime.RawExtension
		err = dec.Decode(&rawObj)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decode is err %v", err)
		}

		obj, _, err := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("rawobj is err%v", err)
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return nil, fmt.Errorf("tounstructured is err %v", err)
		}

		unstructureObj := &unstructured.Unstructured{Object: unstructuredMap}
		gvr, err := d.GtGVR(unstructureObj.GroupVersionKind())
		if err != nil {
			return nil, err
		}
		unstructuredYaml, err := sigyaml.Marshal(unstructureObj)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource as yaml: %w", err)
		}
		_, getErr := dynameicclient.Resource(gvr).Namespace(d.namespace).Get(context.Background(), unstructureObj.GetName(), metav1.GetOptions{})
		if getErr != nil {
			_, createErr := dynameicclient.Resource(gvr).Namespace(d.namespace).Create(context.Background(), unstructureObj, metav1.CreateOptions{})
			if createErr != nil {
				return nil, createErr
			}
		}

		force := true

		if d.namespace == unstructureObj.GetNamespace() {

			ret, err = dynameicclient.Resource(gvr).
				Namespace(d.namespace).
				Patch(context.Background(),
					unstructureObj.GetName(),
					types.ApplyPatchType,
					unstructuredYaml, metav1.PatchOptions{
						FieldManager: unstructureObj.GetName(),
						Force:        &force,
					})

			if err != nil {
				return nil, fmt.Errorf("unable to patch resource: %w", err)
			}

		} else {

			ret, err = dynameicclient.Resource(gvr).
				Patch(context.Background(),
					unstructureObj.GetName(),
					types.ApplyPatchType,
					unstructuredYaml, metav1.PatchOptions{
						Force:        &force,
						FieldManager: unstructureObj.GetName(),
					})
			if err != nil {
				return nil, fmt.Errorf("ns is nil unable to patch resource: %w", err)
			}

		}

		//get, err := dynameicclient.Resource(gvr).Namespace(d.namespace).Get(context.TODO(), unstructureObj.GetName(), metav1.GetOptions{})
		//if err!=nil {
		//	return fmt.Errorf("server: %w", err)
		//}
		//fmt.Println(gvr)
		//
		//
		//status := get.Object["status"].(map[string]interface{})["conditions"]
		//fmt.Println(status)

	}

	return ret, nil
}

func (d *KubeExecutor) GtGVR(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {

	// 获取配置文件
	config, err := clientcmd.BuildConfigFromFlags("", d.confPath)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(gr)

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	return mapping.Resource, nil
}

