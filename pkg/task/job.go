package task

type Job func()

type CommandJob struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Lang        string      `json:"lang"`
	Timeout     int64       `json:"timeout"`
	Retries     int32       `json:"retries"`
	Background  int32       `json:"background"`
	Async       int32       `json:"async"`
	CallbackApi string      `json:"callback_api"`
	User        string      `json:"user"`
	Parameters  []Parameter `json:"parameters"`
}

type DownloadFileJob struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DownloadUrl string `json:"download_url"`
	TargetPath     string `json:"target_path"`
	Parameters  []Parameter `json:"parameters"`
}

type KubeJob struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	YamlStr      string `json:"yaml_str"`
	Type         string `json:"type"`
	KubeConfPath string `json:"kube_conf_path"`
	Namespace    string `json:"namespace"`
}

type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Desc  string `json:"description"`
}
