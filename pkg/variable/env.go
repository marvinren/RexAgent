package variable

import (
	"os"
	"strings"
	"sync"
	"time"
)

var globalEnvParameters map[string]string
var envReadLock = sync.Mutex{}
var updateTime = time.Now()


func GetEnvParamValue(name string) (string, bool){
	var now = time.Now()
	if now.Sub(updateTime) > 1 * time.Minute || globalEnvParameters==nil{

		envReadLock.Lock()
		defer envReadLock.Unlock()

		globalEnvParameters = make(map[string]string)
		environ := os.Environ()
		for _, p := range environ {
			n := strings.SplitN(p, "=", 2)

			globalEnvParameters[n[0]] = n[1]
		}
	}
	return globalEnvParameters[name], true
}
