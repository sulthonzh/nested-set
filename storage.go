package nestedset

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

// NewNodeStorage initializes the storage
func NewNodeStorage(db *gorm.DB) *NodeStorage {
	return &NodeStorage{db: db}
}

// NewNodeStorageCustomTableName initializes the storage
func NewNodeStorageCustomTableName(db *gorm.DB, tableName string) *NodeStorage {
	node := Node{}
	node.SetTableName(tableName)

	return &NodeStorage{
		db:   db,
		node: node,
	}
}

// NodeStorage stores all nodes
type NodeStorage struct {
	db   *gorm.DB
	node Node
}

// IsPopulated tries to find root node in the database.
// If t.Root is not nil it does not do anything and returns true.
// Else, if node with left_key = 1 does not exist, it returns false
// If it exists, it sets t.Root to point to it
func (n NodeStorage) IsPopulated(nodeType int32) (root *Node, err error) {
	r := n.node
	err = n.db.First(&r, "type = ? AND lft = 1", nodeType).Error
	if err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("Node root with type %d does not exist", nodeType)
		return
	} else if err != nil {
		return
	}
	root = &r
	return
}

// Plant inserts root node of the tree to the database
// and sets t.Root to point to the root node.
// It does nothing and returns nil error if t.Root is not nil.
// It returns error in case if something is wrong.
func (n NodeStorage) Plant(root *Node) (err error) {
	root.SetTableName(n.node.TableName())
	root.Left = 1
	root.Right = 2
	err = n.db.Create(&root).Error
	return
}

// GetAllNodes traverses all the tree and returns all its' nodes
func (n NodeStorage) GetAllNodes(parent *Node) (nodes NodeCollection, err error) {
	var results NodeCollection
	sql := fmt.Sprintf(`SELECT node.*,(COUNT(parent.name) - (sub_tree.depth + 1))  AS depth
						FROM %s AS node,
								%s AS parent,
								%s AS sub_parent,
								(
										SELECT node.name, (COUNT(parent.name) - 1) AS depth
										FROM %s AS node,
										%s AS parent
										WHERE node.lft BETWEEN parent.lft AND parent.rgt
												AND node.type = ? AND node.name = ?
										GROUP BY node.name
										ORDER BY node.lft
								)AS sub_tree
						WHERE node.lft BETWEEN parent.lft AND parent.rgt
								AND node.lft BETWEEN sub_parent.lft AND sub_parent.rgt
								AND sub_parent.name = sub_tree.name
						GROUP BY node.name
						ORDER BY node.lft;
`, n.node.TableName(), n.node.TableName(), n.node.TableName(), n.node.TableName(), n.node.TableName())
	rows, err := n.db.Raw(sql, parent.Type, parent.Name).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := n.node
		n.db.ScanRows(rows, &r)
		results = append(results, &r)
	}
	nodes = results
	return
}

// GetNodeByValue search node by value in the tree and returns
// a pointer to it. If there is more than one node with the same value
// it returns the first found (ordered by left_key).
// If something goes wrong it returns non-nil error
func (n NodeStorage) GetNodeByValue(nodeType int32, nodeName string) (node *Node, err error) {
	r := n.node
	err = n.db.Where("type = ? AND name = ?", nodeType, nodeName).Order("lft").First(&r).Error
	if err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("Node with type '%d' and name '%s' does not exist", nodeType, nodeName)
		return
	} else if err != nil {
		return
	}
	node = &r
	return
}

// InsertChild creates new node with value and
// inserts it as child of the parent node.
// If something goes wrong it returns non-nil error.
// Please not that values of left_key, right_key in the parent node
// would be outdated after this operation.
func (n NodeStorage) InsertChild(parent *Node, child *Node) (err error) {
	// Update other childs of partner
	err = n.db.Table(n.node.TableName()).Where("type = ? AND lft > ?", parent.Type, parent.Right).Updates(map[string]interface{}{
		"lft": gorm.Expr("lft + 2"),
		"rgt": gorm.Expr("rgt + 2"),
	}).Error
	if err != nil {
		return
	}
	err = n.db.Table(n.node.TableName()).Where("type = ? AND rgt >= ? AND lft < ?", child.Type, parent.Right, parent.Right).Updates(map[string]interface{}{
		"rgt": gorm.Expr("rgt + 2"),
	}).Error
	if err != nil {
		return
	}

	child.Left = parent.Right
	child.Right = parent.Right + 1
	err = n.db.Create(&child).Error
	return
}

// DeleteNode deletes node n and all it's children
// from the tree. It returns non-nil error
// in case if something goes wrong
func (n NodeStorage) DeleteNode(node *Node) (err error) {
	node.SetTableName(n.node.TableName())
	// Removing node and all it's children
	err = n.db.Delete(n.node, "type = ? AND lft >= ? AND rgt <= ?", node.Type, node.Left, node.Right).Error
	if err != nil {
		return
	}

	// Update parent branch
	leftLim := node.Right - node.Left + 1
	err = n.db.Table(n.node.TableName()).Where("type = ? AND lft < ? AND rgt > ?", node.Type, node.Left, node.Right).Updates(map[string]interface{}{
		"rgt": gorm.Expr("rgt - ?", leftLim),
	}).Error
	if err != nil {
		return
	}

	err = n.db.Table(n.node.TableName()).Where("type = ? AND lft > ?", node.Type, node.Right).Updates(map[string]interface{}{
		"lft": gorm.Expr("lft - ?", leftLim),
		"rgt": gorm.Expr("rgt - ?", leftLim),
	}).Error
	return
}

// GetParent returns parent of node. If there is no parent node
// it returns ErrNodeDoesNotExist error. It also returns this error
// for root node, so user of this function should check if input node
// is root node himself
func (n NodeStorage) GetParent(node *Node) (parent *Node, err error) {
	err = n.db.Table(n.node.TableName()).Order("lft desc").First(&parent, "type = ? AND lft < ? AND rgt > ?", node.Type, node.Left, node.Right).Error
	if err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("Node parent with type %d does not exist", node.Type)
		return
	} else if err != nil {
		return
	}
	return
}

// IsDescendantOf checks if node child is really descendant of node parent
func IsDescendantOf(child, parent *Node) bool {
	if child.Left > parent.Left && child.Right < parent.Right {
		return true
	}
	return false
}

// MoveNode moves node to new parent newParent.
// It refuses to move node to it's descendant or move root node
// and returns error in that case.
// It does nothing and returns nil error in case of trying
// to move node to itself.
func (n NodeStorage) MoveNode(newParent, node *Node) (err error) {
	newParent.SetTableName(n.node.TableName())
	node.SetTableName(n.node.TableName())

	width := node.Right - node.Left + 1
	distance := node.Right - newParent.Right

	// Doing checks if operation is possible
	if distance == 0 {
		return nil
	}

	if IsDescendantOf(newParent, node) {
		return fmt.Errorf("Could not move node to it's own descendant")
	}

	_, err = n.GetParent(node)
	if err != nil {
		return fmt.Errorf("Not possible to move orphan node (or root node)")
	}

	log.Printf("Removing our branch")
	err = n.db.Table(n.node.TableName()).Where("type = ? AND lft >= ? AND rgt <= ? ", node.Type, node.Left, node.Right).Updates(map[string]interface{}{
		"lft": gorm.Expr("-lft"),
		"rgt": gorm.Expr("-rgt"),
	}).Error
	if err != nil {
		return
	}

	log.Print("Decrease right key for nodes after removal and parents")
	err = n.db.Table(n.node.TableName()).Where("type = ? AND rgt > ? ", node.Type, node.Right).Updates(map[string]interface{}{
		"rgt": gorm.Expr("rgt - ?", width),
	}).Error
	if err != nil {
		return
	}

	newParentUpdated := newParent.Right - width
	if distance > 0 {
		newParentUpdated = newParent.Right
	}
	log.Printf("Increasing right key after the place of insertion and new parents")
	err = n.db.Table(n.node.TableName()).Where("type = ? AND rgt >= ? ", node.Type, newParentUpdated).Updates(map[string]interface{}{
		"rgt": gorm.Expr("rgt + ?", width),
	}).Error
	if err != nil {
		return
	}

	log.Printf("Increasing left key after the place of insertion")
	err = n.db.Table(n.node.TableName()).Where("type = ? AND lft > ? ", node.Type, newParentUpdated).Updates(map[string]interface{}{
		"lft": gorm.Expr("lft + ?", width),
	}).Error
	if err != nil {
		return
	}

	var d int32
	if distance > 0 {
		newParentRK := newParent.Right + width
		d = node.Right - newParentRK + 1
	} else {
		d = distance + 1
	}
	log.Print("Actually moving our branch")
	err = n.db.Table(n.node.TableName()).Where("type = ? AND lft <= 0 ", node.Type).Updates(map[string]interface{}{
		"lft": gorm.Expr("-lft - ?", d),
		"rgt": gorm.Expr("-rgt - ?", d),
	}).Error
	if err != nil {
		return
	}
	return
}

// RenameNode updates value of node. It returns non-nil error
// in case if something goes wrong
func (n NodeStorage) RenameNode(node *Node) (err error) {
	node.SetTableName(n.node.TableName())
	err = n.db.Save(&node).Error
	return
}
