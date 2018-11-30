package nestedset

import (
	"github.com/satori/go.uuid"
)

type (
	// Node represents a node of nested set tree
	Node struct {
		_tableName string         `sql:"-"`
		ID         uuid.UUID      `sql:"type:varbinary(36);NOT NULL" gorm:"primary_key;column:id" json:"-"`  //
		Left       int32          `sql:"type:int;NOT NULL" gorm:"column:lft" json:"-"`                       //
		Right      int32          `sql:"type:int;NOT NULL" gorm:"column:rgt" json:"-"`                       //
		Name       string         `sql:"type:varchar(50);NOT NULL" gorm:"column:name" json:"name,omitempty"` //
		Type       int32          `sql:"type:int;NOT NULL" gorm:"column:type" json:"type,omitempty"`         //
		JSONData   string         `sql:"type:JSON;default:NULL" gorm:"column:data" json:"-"`                 //
		Data       interface{}    `sql:"-" json:"data,omitempty"`                                            //
		Depth      int32          `sql:"-" json:"depth,omitempty"`                                           //
		Children   NodeCollection `sql:"-" json:"children,omitempty"`                                        //
	}
)

// NeWNode returns pointer to newly initialized Node
func NeWNode(left, right int32, name string) *Node {
	n := &Node{
		ID:    uuid.Must(uuid.NewV4()),
		Left:  left,
		Right: right,
		Name:  name,
	}
	return n
}

// TableName ::
func (n *Node) TableName() string {
	if n._tableName != "" {
		return n._tableName
	}
	return "nodes"
}

// SetTableName ::
func (n *Node) SetTableName(name string) {
	n._tableName = name
}
