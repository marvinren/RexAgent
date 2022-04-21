package tests

import (
	"rexagent/pkg/conf"
	"testing"
)

func TestConfigStrReader(t *testing.T) {

	config_str := `<?xml version="1.0" encoding="UTF-8"?>
<config>
    <server>
        <listen>:2465</listen>
    </server>
    <log>
        <path>/data/logs/rexagent/rexagent.log</path>
        <level>4</level>
    </log>
</config>`

	xconf, err := conf.XConfigFromData([]byte(config_str), map[string]string{})
	if err != nil {
		t.Errorf("parse error: %s", err.Error())
		return
	}

	config := xconf.ToConfig()
	t.Log(config)
	if config.Server.Listen != ":2465" {
		t.Errorf("parse error: %s", err)
	}

}

func TestConfigVarReader(t *testing.T) {
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
            <value>${world}</value>
        </var>
    </vars>

</config>`
	xconf, err := conf.XConfigFromData([]byte(config_str), map[string]string{})
	if err != nil {
		t.Errorf("parse error: %s", err.Error())
		return
	}

	config := xconf.ToConfig()
	t.Log(config)
	if config.Vars == nil {
		t.Errorf("not parse Vars")
		return
	}
	if _, ok := config.Vars["vars"]; !ok {
		t.Errorf("not parse Vars[vars]")
		return
	}
	if v, ok2 := config.Vars["vars"].Vars["foo"]; !ok2 || v.Value != "bar" {
		t.Errorf("parse foo varible error")
	}
}
