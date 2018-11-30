package nestedset

import "encoding/json"

type (
	// NodeCollection ::
	NodeCollection []*Node
)

// Marshal ::
func (nodes NodeCollection) Marshal() ([]byte, error) {
	return json.Marshal(nodes)
}

// GenerateTree ::
func (nodes NodeCollection) GenerateTree(childStruct interface{}) NodeCollection {
	query := NodeCollection{}
	for _, n := range nodes {
		if n.Depth == 0 {
			query = append(query, n)
		}
	}
	return nodes.runQueryRecursive(query, childStruct)
}

func (nodes NodeCollection) runQueryRecursive(query NodeCollection, v interface{}) (result NodeCollection) {
	for _, q := range query {
		var vTemp interface{}
		if q.JSONData != "" {
			vTemp = v
			b := []byte(q.JSONData)
			json.Unmarshal(b, &vTemp)
		}
		result = append(result, &Node{
			ID:       q.ID,
			Name:     q.Name,
			Data:     vTemp,
			Children: nodes.getChildrens(q, v),
		})
	}
	return result
}

func (nodes NodeCollection) getChildrens(parent *Node, v interface{}) NodeCollection {
	query := NodeCollection{}
	for _, n := range nodes {
		if n.Depth == parent.Depth+1 && n.Left > parent.Left && n.Right < parent.Right {
			query = append(query, n)
		}
	}
	return nodes.runQueryRecursive(query, v)
}
