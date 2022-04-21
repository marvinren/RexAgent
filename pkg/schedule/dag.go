package schedule

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Flow struct {
	NodeList []*Node `json:"nodeList"`
	LineList []*Line `json:"lineList"`
}

type Node struct {
	NodeId string `json:"nodeId"`
	Left   string `json:"left"`
	Top    string `json:"top"`
	Class  string `json:"class"`
	Icon   string `json:"icon"`
	Name   string `json:"name"`
	Id     string `json:"id"`
}

type Line struct {
	SourceId string `json:"sourceId"`
	TargetId string `json:"targetId"`
	Label    string `json:"label"`
}

type DAG struct {
	Vertexs []*Vertex
}

type Vertex struct {
	Key      string
	Name     string
	Value    interface{}
	Parents  []*Vertex
	Children []*Vertex
}

func (dag *DAG) AddVertex(v *Vertex) {
	dag.Vertexs = append(dag.Vertexs, v)
}

func (dag *DAG) AddEdge(from, to *Vertex) {
	from.Children = append(from.Children, to)
	to.Parents = append(to.Parents, from)
}

func (dag *DAG) GetRootVertex() (root *Vertex) {
	for _, v := range dag.Vertexs {
		if v == nil {
			continue
		}
		if v.Parents == nil || len(v.Parents) == 0 {
			return v
		}
	}
	return nil
}

func NewDAGGraph(flowStr string) (dag *DAG, err error) {
	var flow Flow = Flow{}
	nodeMap := make(map[string]*Vertex)
	dag = &DAG{}

	// unmarshal json str
	err = json.Unmarshal([]byte(flowStr), &flow)
	if err != nil {
		return nil, err
	}
	// initial vert list
	dag.Vertexs = make([]*Vertex, 0)

	// parse the graph node
	for _, tmpNode := range flow.NodeList {
		tmpVertex := &Vertex{Key: tmpNode.NodeId, Name: tmpNode.Name, Value: nil, Parents: make([]*Vertex, 0), Children: make([]*Vertex, 0)}
		nodeMap[tmpNode.NodeId] = tmpVertex
		dag.AddVertex(tmpVertex)
	}

	//parse the graph line
	for line := range flow.LineList {
		tmpLine := flow.LineList[line]
		sourceVertex, ok := nodeMap[tmpLine.SourceId]
		if !ok {
			return nil, errors.New(fmt.Sprintf("the line source not found %s", tmpLine.SourceId))
		}
		targetVertex, ok2 := nodeMap[tmpLine.TargetId]
		if !ok2 {
			return nil, errors.New(fmt.Sprintf("the line target not found %s", tmpLine.TargetId))
		}

		dag.AddEdge(sourceVertex, targetVertex)
	}

	return dag, nil
}

func BFSNew(root *Vertex) [][]*Vertex {
	q := NewQueue()
	q.Add(root)
	visited := make(map[string]*Vertex)
	all := make([][]*Vertex, 0)
	for q.Length() > 0 {
		qSize := q.Length()
		tmp := make([]*Vertex, 0)
		for i := 0; i < qSize; i++ {
			//pop vertex
			currVert := q.Pop().(*Vertex)
			if _, ok := visited[currVert.Key]; ok {
				continue
			}
			visited[currVert.Key] = currVert
			tmp = append(tmp, currVert)
			for _, val := range currVert.Children {
				if _, ok := visited[val.Key]; !ok {
					q.Add(val) //add child
				}
			}
		}
		all = append([][]*Vertex{tmp}, all...)
	}
	return all
}
