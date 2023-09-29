package kanban

import (
	"errors"
	"project-manager-go/data"

	"gorm.io/gorm"
)

type cards struct {
	*data.TreeStore
}

func (s *cards) GetAll(ctx *data.DBContext, dest data.IItemsList) error {
	cards := make([]data.Item, 0)
	err := ctx.DB.
		Preload("AssignedUsers").
		Order("kanban_index").
		Find(&cards).
		Error
	if err != nil {
		return err
	}

	dest.PutItems(cards)

	return nil
}

func (s *cards) GetMaxIndex(ctx *data.DBContext, row, column int) (int, error) {
	item := data.Item{}
	err := ctx.DB.
		Where("project_id = ? AND kanban_column_id = ?", row, column).
		Order("kanban_index DESC").
		Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	return item.Kanban_CardIndex + 1, err
}

func (s *cards) ShiftIndex(ctx *data.DBContext, from, offset int) error {
	err := ctx.DB.
		Model(data.Item{}).
		Where("kanban_index >= ?", from).
		Update("kanban_index", gorm.Expr("kanban_index + ?", offset)).
		Error
	return err
}
