package conf

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type XConfig struct {
	XMLName xml.Name `xml:"config"`
	Server  XServer  `xml:"server"`
	Logger  XLog      `xml:"log"`
	Vars      []XVars     `xml:"vars"`

}

type XServer struct {
	Listen string `xml:"listen"`
}

type XLog struct {
	Path  string `xml:"path"`
	Level int32 `xml:"level"`
}

type XVars struct {
	Name string `xml:"id,attr"`
	Vars []XVar `xml:"var"`
}

type XVar struct {
	Name     string   `xml:"id,attr"`
	Value    string   `xml:"value"`
	Readonly bool     `xml:"readonly,attr"`
	Patterns []string `xml:"pattern"`
	Expand   bool     `xml:"expand,attr"`
}



type XValidator struct {
	Name string `xml:"name,attr"`
	//class  string
	Pattern string `xml:",chardata"`
}

func XConfigFromData(data []byte, entities map[string]string) (*XConfig, error) {
	ret := XConfig{}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.Entity = entities
	err := decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func XConfigFromReader(reader io.Reader, entities map[string]string) (*XConfig, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return XConfigFromData(data, entities)
}

func XConfigFromFile(confPath string, entities map[string]string) (*XConfig, error) {
	entities["__file__"] = confPath
	entities["__dir__"] = path.Dir(confPath)
	reader, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}
	return XConfigFromReader(reader, entities)
}

func (conf *XConfig) ToConfig() *Config {
	ret := Config{}
	conf.IntoConfig(&ret)
	return &ret
}

func (conf *XConfig) IntoConfig(ret *Config) {
	if ret.Server.Listen == "" {
		ret.Server = Server{
			Listen: conf.Server.Listen,
		}
	}

	if ret.Log.Path == "" {
		ret.Log = Logger{
			Path: conf.Logger.Path,
			Level: conf.Logger.Level,
		}

	}

	if ret.Vars == nil {
		ret.Vars = make(map[string]*Vars)
	}
	for _, vars := range conf.Vars {
		vname := vars.Name
		if ret.Vars[vname] == nil {
			ret.Vars[vname] = &Vars{
				Vars: make(map[string]*Var),
			}
		}
		for _, v := range vars.Vars {
			ret.Vars[vname].Vars[v.Name] = &Var{
				Value:    v.Value,
				Patterns: v.Patterns,
				Readonly: v.Readonly,
				Expand:   v.Expand,
			}
		}
	}

}

type LoadConfigError struct {
	Path string
	Err  error
}

func (self LoadConfigError) Error() string {
	return fmt.Sprintf("load %s failed: %s", self.Path, self.Err.Error())
}

func LoadXmlConfig(files, dirs []string, params map[string]string) (config Config, err error) {
	for _, confPath := range files {
		confPath, _ = filepath.Abs(confPath)
		params["__file__"] = confPath
		params["__dir__"] = path.Dir(confPath)
		xconf, err := XConfigFromFile(confPath, params)
		if err != nil {
			return config, LoadConfigError{Path: confPath, Err: err}
		}
		xconf.IntoConfig(&config)
	}
	for _, confDirPath := range dirs {
		filesInfo, err := ioutil.ReadDir(confDirPath)
		if err != nil {
			return config, LoadConfigError{Path: confDirPath, Err: err}
		}
		for _, fileInfo := range filesInfo {
			filename := fileInfo.Name()
			if strings.HasSuffix(fileInfo.Name(), ".conf") || strings.HasSuffix(fileInfo.Name(), ".xml") {
				if fileInfo.IsDir() {
					continue
				}
				confPath := filepath.Join(confDirPath, filename)
				confPath, _ = filepath.Abs(confPath)

				xconf, err := XConfigFromFile(confPath, params)
				if err != nil {
					return config, LoadConfigError{Path: confPath, Err: err}
				}
				xconf.IntoConfig(&config)
			}
		}
	}
	return
}
