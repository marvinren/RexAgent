package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"rexagent/pkg/executor/command"
	"rexagent/pkg/executor/kube"
	"rexagent/pkg/task"
)

func ResponseSucc(message string, w http.ResponseWriter, data interface{}) {
	var res = map[string]interface{}{"result": 0, "message": message, "data": data}
	response, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func ResponseError(message string, w http.ResponseWriter) {
	var res = map[string]interface{}{"result": 1, "message": message}
	response, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func CommandJobHandler(w http.ResponseWriter, r *http.Request){
	var t task.CommandJob

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &t)
	if err != nil {
		ResponseError(err.Error(), w)
		return
	}

	executor := command.CommandExecutor{}
	ret, err := executor.Execute(&t, nil)
	if err != nil {
		ResponseError(err.Error(), w)
	}else{
		ResponseSucc("ok", w, string(ret))
	}

}

func KubeJobHandler(w http.ResponseWriter, r *http.Request){
	var t task.KubeJob
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &t)
	if err != nil {
		ResponseError(err.Error(), w)
		return
	}

	executor := kube.KubeExecutor{}
	var ret []byte
	ret, err = executor.Execute(&t)
	if err != nil {
		ResponseError(err.Error(), w)
	}else{
		ResponseSucc("ok", w, string(ret))
	}

}
