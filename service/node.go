package service

import (
	"project-manager-go/data"
)

type Node struct {
	ID        int
	ParentID  int
	ProjectID int
	Index     int
}

func (n *Node) PutItem(item data.Item) {
	n.ID = item.ID
	n.ParentID = item.ParentID
	n.ProjectID = item.ProjectID
	n.Index = item.Index
}

func (n Node) FillItem(item *data.Item) {
	item.ParentID = n.ParentID
	item.ProjectID = n.ProjectID
	item.Index = n.Index
}
