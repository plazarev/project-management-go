package service

import (
	"fmt"
	"project-manager-go/data"
)

// TreeService provides data manipalations such as movement items in the hierarhy structrue, getting children of the root and etc.
type TreeService struct {
	tree *data.TreeStore
}

type MiddleHandler func(dbCtx *data.DBContext, childId int, fielids map[string]any) error

func NewTreeService(store *data.TreeStore) *TreeService {
	return &TreeService{
		tree: store,
	}
}

// Changes project id for item with the given. For all children of the item project id also changed.
// cb MiddleHanler allows to define custom properties for each child of the item
func (s *TreeService) ChangeNodeProject(dbCtx *data.DBContext, id, newProjectId int, cb MiddleHandler) (err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	node := Node{}
	err = s.tree.GetOne(dbCtx, id, &node)
	if err != nil {
		return err
	}

	// get the max index in the given project
	i, err := s.tree.MaxBranchIndex(dbCtx, newProjectId, 0)
	if err != nil {
		return err
	}

	oldProjectId := node.ProjectID
	oldParentId := node.ParentID
	oldIndex := node.Index

	node.ParentID = 0
	node.ProjectID = newProjectId
	node.Index = i + 1
	err = s.tree.Update(dbCtx, node.ID, &node)
	if err != nil {
		return err
	}

	list := data.NewItemsList[Node]()
	err = s.tree.GetAllChildren(dbCtx, oldProjectId, id, list)
	if err != nil {
		return err
	}

	for _, child := range list.GetArray() {
		fields := make(map[string]any)

		// update project for the child
		fields["project_id"] = newProjectId

		if cb != nil {
			// set custom properties for the children
			err = cb(dbCtx, child.ID, fields)
			if err != nil {
				return err
			}
		}

		err = s.tree.UpdateFields(dbCtx, child.ID, fields)
		if err != nil {
			return err
		}
	}

	// remove missing position
	err = s.tree.ShiftIndex(dbCtx, oldProjectId, oldParentId, oldIndex+1, -1)

	return err
}

// Moves item with the given id to another position defined by project, parent and index values
func (s *TreeService) Move(dbCtx *data.DBContext, id, newProject, newParent, newIndex int) (err error) {
	if newIndex < 0 {
		return fmt.Errorf("new index must be not negative")
	}

	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	maxBranchIndex, err := s.tree.MaxBranchIndex(dbCtx, newProject, newParent)
	if err != nil {
		return err
	}
	if maxBranchIndex == -1 {
		maxBranchIndex = 0
	}

	node := Node{}
	err = s.tree.GetOne(dbCtx, id, &node)
	if err != nil {
		return err
	}
	sameBranch := node.ProjectID == newProject && node.ParentID == newParent

	if !sameBranch {
		maxBranchIndex++
	}

	if newIndex > maxBranchIndex {
		newIndex = maxBranchIndex
	}

	// remove the node from its current parent's child list
	parentChildList := data.NewItemsList[Node]()
	err = s.tree.GetBranchChildren(dbCtx, node.ProjectID, node.ParentID, parentChildList)
	if err != nil {
		return err
	}
	parentChildArr := parentChildList.GetArray()
	for i := range parentChildArr {
		child := &parentChildArr[i]
		if child.ID == node.ID {
			parentChildArr = append(parentChildArr[:i], parentChildArr[i+1:]...)
			break
		}
	}

	// find the new branch
	newParentChildList := data.NewItemsList[Node]()
	var newParentChildArr []Node
	if sameBranch {
		newParentChildArr = parentChildArr
	} else {
		err = s.tree.GetBranchChildren(dbCtx, newProject, newParent, newParentChildList)
		if err != nil {
			return err
		}
		newParentChildArr = newParentChildList.GetArray()
	}

	node.ParentID = newParent
	node.ProjectID = newProject

	if newIndex == len(newParentChildArr) {
		// push to the end of the branch
		newParentChildArr = append(newParentChildArr, node)
	} else {
		// insert between other nodes
		left := newParentChildArr[:newIndex]
		leftCopy := make([]Node, len(left))
		copy(leftCopy, left)

		right := newParentChildArr[newIndex:]
		newChildren := append(leftCopy, node)
		newChildren = append(newChildren, right...)
		newParentChildArr = newChildren
	}

	// ajust the index values of the new branch
	for i := range newParentChildArr {
		child := &newParentChildArr[i]
		child.Index = i
		err = s.tree.Update(dbCtx, child.ID, child)
		if err != nil {
			return err
		}
	}

	if !sameBranch {
		// if the branch to which the task was moved not the same,
		// then adjust indices in the old branch
		for i := range parentChildArr {
			child := &parentChildArr[i]
			child.Index = i
			err = s.tree.Update(dbCtx, child.ID, child)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
