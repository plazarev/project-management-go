package data

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Tree store allows to interact with data in tree structure
type TreeStore struct {
	*ItemsStore
}

func NewTreeStore(itemsProcessor *ItemsProcessorProvider) *TreeStore {
	return &TreeStore{
		ItemsStore: NewItemsStore(itemsProcessor),
	}
}

// Returns an index of item with the given id
func (s *TreeStore) GetItemIndex(ctx *DBContext, id int) (int, error) {
	item := Item{}
	err := ctx.DB.Take(&item, id).Error
	return item.Index, err
}

// Returns an index of item with the given id
func (s *TreeStore) GetByIndex(ctx *DBContext, project, parent, index int, dest IItem) error {
	item := Item{}
	err := ctx.DB.
		Take(&item, "project_id = ? AND parent_id = ? AND `index` = ?", project, parent, index).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	dest.PutItem(item)
	return err
}

// Returns a maximum index in the given branch
func (s *TreeStore) MaxBranchIndex(ctx *DBContext, project, parent int) (int, error) {
	item := Item{}
	err := ctx.DB.
		Where("project_id = ? AND parent_id = ?", project, parent).
		Order("`index` DESC").
		Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return -1, nil
	}

	return item.Index, err
}

// Shifts the index of items in the given branch
func (s *TreeStore) ShiftIndex(ctx *DBContext, project, parent, from, offset int) error {
	if from < 0 {
		from = 0
	}
	if from+offset < 0 {
		return fmt.Errorf("cannot use offset causing negative indices")
	}

	err := ctx.DB.Model(&Item{}).
		Where("project_id = ? AND parent_id = ? AND `index` >= ?", project, parent, from).
		Update("index", gorm.Expr("`index` + ?", offset)).Error
	return err
}

// Returns all children of the parent node
func (s *TreeStore) GetAllChildren(ctx *DBContext, project, parent int, dest IItemsList) error {
	// parent ID may not always identify the branch,
	// because the parent can be equals 0 in many projects,
	// to make sure the branch is correct need to use project ID

	items, err := s.findProjectItems(ctx, project)
	if err != nil {
		return err
	}
	children := findAllChildren(items, parent)

	dest.PutItems(children)

	return err
}

// Returns array of IDs of all children of the parent node
func (s *TreeStore) GetAllChildrenIDs(ctx *DBContext, project, parent int) ([]int, error) {
	// parent ID may not always identify the branch,
	// because the parent can be equals 0 in many projects,
	// to make sure the branch is correct need to use project ID

	items, err := s.findProjectItems(ctx, project)
	if err != nil {
		return nil, err
	}

	children := findAllChildren(items, parent)

	return toIds(children), err
}

// Returns only the branch children (only the level of nesting)
func (s *TreeStore) GetBranchChildren(ctx *DBContext, project, parent int, dest IItemsList) error {
	// parent ID may not always identify the branch,
	// because the parent can be equals 0 in many projects,
	// to make sure the branch is correct need to use project ID

	children, err := s.findBranchChildren(ctx, project, parent)
	if err != nil {
		return err
	}

	dest.PutItems(children)

	return err
}

// Returns array of IDs of only the branch children (only the level of nesting)
func (s *TreeStore) GetBranchChildrenIDs(ctx *DBContext, project, parent int) ([]int, error) {
	// parent ID may not always identify the branch,
	// because the parent can be equals 0 in many projects,
	// to make sure the branch is correct need to use project ID

	children, err := s.findBranchChildren(ctx, project, parent)
	if err != nil {
		return nil, err
	}

	return toIds(children), err
}

// Deletes an item and its all associated entities
func (s *TreeStore) DeleteCascade(ctx *DBContext, id int) ([]int, error) {
	item, err := s.getById(ctx, id)
	if err != nil {
		return nil, err
	}

	// find item children
	children, err := s.GetAllChildrenIDs(ctx, item.ProjectID, id)
	if err != nil {
		return nil, err
	}

	ids := append([]int{id}, children...)

	err = s.Delete(ctx, ids...)
	if err != nil {
		return nil, err
	}

	// should update branch order
	items, err := s.findBranchChildren(ctx, item.ProjectID, item.ParentID)
	for i := range items {
		item := &items[i]
		item.Index = i

		err := ctx.DB.Save(item).Error
		if err != nil {
			return nil, err
		}
	}

	return children, err
}

func (s *TreeStore) findProjectItems(ctx *DBContext, project int) ([]Item, error) {
	items := make([]Item, 0)
	err := ctx.DB.
		Order("parent_id, `index`").
		Find(&items, "project_id = ?", project).
		Error

	return items, err
}

func (s *TreeStore) findBranchChildren(ctx *DBContext, project, parent int) ([]Item, error) {
	children := make([]Item, 0)
	err := ctx.DB.
		Order("`index`").
		Find(&children, "parent_id = ? AND project_id = ?", parent, project).
		Error
	return children, err
}

func findAllChildren(items []Item, id int) []Item {
	buff := make([]Item, 0)
	for _, item := range items {
		if item.ParentID == id {
			buff = append(buff, item)
			buff = append(buff, findAllChildren(items, item.ID)...)
		}
	}
	return buff
}

func toIds(items []Item) []int {
	ids := make([]int, len(items))
	for i := range items {
		ids[i] = items[i].ID
	}
	return ids
}
