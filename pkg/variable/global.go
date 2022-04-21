package variable

import (
	"rexagent/pkg/conf"
	"strings"
	"sync"
)

var GlobalParams = map[string]string{}

var varsLock = sync.Mutex{}

func InitGlobalVars(conf *conf.Config) error {
	for group_name, group_v := range conf.Vars {
		for n, v := range group_v.Vars {
			if v.Expand {
				value := ExpandVars(v.Value, func(name string) string { return GlobalParams[name]})
				SetGlobalGroupVar(group_name, n, value)
			} else {
				SetGlobalGroupVar(group_name, n, v.Value)
			}
		}
	}

	return nil
}

func GetGlobalGroupVar(g string, k string) (string, bool) {
	varsLock.Lock()
	defer varsLock.Unlock()
	if g == "" {
		ret, ok := GlobalParams[k]
		return ret, ok
	} else {
		ret, ok := GlobalParams[g+"."+k]
		return ret, ok
	}
}

func SetGlobalGroupVar(g string, k string, v string) {
	varsLock.Lock()
	defer varsLock.Unlock()
	if g == "" {
		GlobalParams[k] = v
	} else {
		GlobalParams[g+"."+k] = v
	}
}

func GetGlobalVarValue(name string) (string, bool) {
	n := strings.SplitN(name, ".", 2)
	if len(n) == 1 {
		return GetGlobalGroupVar("", name)
	} else if len(n) == 2 {
		return GetGlobalGroupVar(n[0], n[1])
	} else {
		return "", false
	}
}
