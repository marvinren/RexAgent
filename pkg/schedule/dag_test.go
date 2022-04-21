package schedule

import (
	"io/ioutil"
	"testing"
)

func TestDAGBuild(t *testing.T){

	str, err := ioutil.ReadFile("./flow_example.json")
	if err != nil {
		t.Errorf(err.Error())
	}
	graph, err := NewDAGGraph(string(str))

	bfsNew := BFSNew(graph.GetRootVertex())

	for i, _ := range bfsNew {
		v := bfsNew[len(bfsNew) - i -1]
		t.Log("================")
		for _, n := range v {
			t.Log(n.Name)
		}
	}


}
