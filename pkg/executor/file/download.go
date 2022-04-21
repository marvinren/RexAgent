package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"rexagent/pkg/task"
	"rexagent/pkg/variable"
)

type FileDownloadExecutor struct {
	id   string
	name string
}

func (fd *FileDownloadExecutor) Execute(t *task.DownloadFileJob, ctx *task.Context) (outBuf []byte, err error) {

	// 对job进行校验
	err = fd.Validate(t)
	if err != nil {
		return nil, err
	}

	//必要的时候初始化上下文Context
	if ctx == nil {
		ctx = &task.Context{}
	}
	//初始化上下文的参数
	if t.Parameters != nil {
		for _, p := range t.Parameters {
			ctx.SetVar(p.Name, p.Value)
		}
	}

	// 获取下载urls地址
	url := variable.ExpandVars(t.DownloadUrl, func(name string) string {
		v, _ := ctx.GetVar(name)
		return v
	})

	// 下载
	filePath, err2 := fd.Download(url, t.TargetPath)

	return []byte(filePath), err2
}

func (fd *FileDownloadExecutor) Validate(t *task.DownloadFileJob) error {
	if t == nil {
		return errors.New("the job is null")
	}

	if t.DownloadUrl == "" {
		return errors.New("the download url is null")
	}

	return nil

}

func (fd *FileDownloadExecutor) Download(url string, target string) (string, error) {
	var filePath string
	//fileName := path.Base(url)
	res, err := http.Get(url)

	if err != nil {
		fmt.Println("A error occurred!")
		return "", err
	}
	defer res.Body.Close()

	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	filePath = target
	file, err := os.Create(filePath)
	if err != nil {
		return filePath, err
	}

	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	_, err2 := io.Copy(writer, reader)
	return filePath, err2

}
