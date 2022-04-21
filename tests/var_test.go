package tests

import (
	"rexagent/pkg/conf"
	"rexagent/pkg/variable"
	"testing"
)

func TestGlobalBuild(t *testing.T) {
	config_str := `<?xml version="1.0" encoding="UTF-8"?>
<config>
    <server>
        <listen>:2465</listen>
    </server>
    <log>
        <path>/tmp/log/</path>
        <level>4</level>
    </log>


    <vars id="vars">
        <var id="foo">
            <value>bar</value>
        </var>
        <var id="hello" expand="true">
            <value>${vars.foo}</value>
        </var>
    </vars>

</config>`
	xconf, err := conf.XConfigFromData([]byte(config_str), map[string]string{})
	if err != nil {
		t.Errorf("parse error: %s", err.Error())
		return
	}

	config := xconf.ToConfig()
	err = variable.InitGlobalVars(config)
	if err != nil {
		t.Errorf("init global vars error: %s", err.Error())
		return
	}
	if value, ok := variable.GetGlobalGroupVar("vars", "foo"); !ok || value != "bar" {
		t.Errorf("var error")
		return
	}

	t.Log(variable.GetGlobalGroupVar("vars", "hello"))

	if groupVar, ok2 := variable.GetGlobalGroupVar("vars", "hello"); !ok2 || groupVar != "bar" {
		t.Errorf("expand parameters error, ret: %s", groupVar)
	}

}

func TestExpandVar(t *testing.T) {
	mapping := map[string]string{
		"appName": "my_app",
		"appIP":   "0.0.0.0",
		"appPort": "8080",
	}
	expand := variable.ExpandVars("${appIP}:${appPort}", func(name string)string {return mapping[name]})

	t.Log(expand)
	if expand != "0.0.0.0:8080" {
		t.Errorf("expand error, ret: %s", expand)
	}
}

func TestEnvParams(t *testing.T) {

	path, ok := variable.GetEnvParamValue("PATH")
	if !ok {
		t.Errorf("get env var error")
	}
	t.Log(path)
}