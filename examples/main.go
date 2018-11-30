package main

import (
	"encoding/json"
	"fmt"

	nestedset "github.com/sulthonzh/nested-set"
	"github.com/sulthonzh/nested-set/examples/db"
	"github.com/sulthonzh/nested-set/examples/utils"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		utils.ExitOnFailure(err)
	}

	n := nestedset.NewNodeStorageCustomTableName(db, "scraper_paths")
	// root := nestedset.Node{
	// 	ID:    uuid.Must(uuid.NewV4()),
	// 	Type:  1,
	// 	Name:  "listings",
	// 	Left:  1,
	// 	Right: 2,
	// }
	parent, err := n.GetNodeByValue(1, "listings")
	// n.Plant(&root)
	// node.GetAllNodes(1)
	if err != nil {
		utils.ExitOnFailure(err)
	}

	// n.InsertChild(parent, &nestedset.Node{
	// 	ID:   uuid.Must(uuid.NewV4()),
	// 	Type: 1,
	// 	Name: "price",
	// })

	nodes, err := n.GetAllNodes(parent)
	if err != nil {
		utils.ExitOnFailure(err)
	}
	child := struct {
		Target interface{} `json:"target"`
		Field  string      `json:"field"`
		Info   string      `json:"info"`
	}{}

	treeNodes := nodes.GenerateTree(child)
	b, _ := treeNodes.Marshal()
	fmt.Println(string(b))

	err = json.Unmarshal(b, &nodes)
	if err != nil {
		utils.ExitOnFailure(err)
	}

}
