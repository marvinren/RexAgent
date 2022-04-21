package file

import (
	"os"
	"rexagent/pkg/task"
	"testing"
)

func TestDownloadFileJob(t *testing.T) {

	job := &task.DownloadFileJob{
		Id:          "1",
		Name:        "download",
		DownloadUrl: "https://archiva-maven-storage-prod.oss-cn-beijing.aliyuncs.com/repository/central/%23/log4/log4j-2.17.0/log4-log4j-2.17.0.jar",
		TargetPath:  "/tmp/log4-log4j-2.17.0.jar",
	}
	e := FileDownloadExecutor{}
	fileName, err := e.Execute(job, nil)
	t.Log(string(fileName))
	if err != nil {
		t.Errorf("downdload err, %s", err)
	}

	_, err2 := os.Stat(job.TargetPath)
	if err2 != nil || os.IsNotExist(err2) {

		t.Errorf("file not found")
	}

}
