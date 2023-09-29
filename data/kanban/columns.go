package kanban

import (
	"errors"
	"project-manager-go/data"

	"gorm.io/gorm"
)

type KanbanColumn struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Index int    `json:"index"`
}

type columns struct{}

func (s *columns) GetOne(ctx *data.DBContext, id int) (KanbanColumn, error) {
	column := KanbanColumn{}
	err := ctx.DB.Take(&column, id).Error
	return column, err
}

func (s *columns) GetAll(ctx *data.DBContext) ([]KanbanColumn, error) {
	columns := make([]KanbanColumn, 0)
	err := ctx.DB.Order("`index`").Find(&columns).Error
	return columns, err
}

func (s *columns) Add(ctx *data.DBContext, column KanbanColumn) (int, error) {
	err := ctx.DB.Create(&column).Error
	return column.ID, err
}

func (s *columns) Update(ctx *data.DBContext, id int, upd KanbanColumn) error {
	obj, err := s.GetOne(ctx, id)
	if err != nil {
		return err
	}

	obj.Label = upd.Label
	obj.Index = upd.Index

	err = ctx.DB.Save(&obj).Error

	return err
}

func (s *columns) Delete(ctx *data.DBContext, id int) ([]int, error) {
	// collect column items
	cards := make([]data.Item, 0)
	err := ctx.DB.Find(&cards, "kanban_column_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int, len(cards))
	for i := range cards {
		ids[i] = cards[i].ID
	}

	// delete column items
	err = ctx.DB.Delete(&data.Item{}, "kanban_column_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// delete column
	err = ctx.DB.Delete(&KanbanColumn{}, id).Error
	if err != nil {
		return nil, err
	}

	return ids, err
}

func (s *columns) GetMaxIndex(ctx *data.DBContext) (int, error) {
	column := KanbanColumn{}

	err := ctx.DB.Order("`index` DESC").Take(&column).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	return column.Index + 1, err
}

func (s *columns) GetByIndex(ctx *data.DBContext, index int) (KanbanColumn, error) {
	column := KanbanColumn{}
	err := ctx.DB.Order("`index` DESC").Take(&column, "`index` = ?", index).Error
	return column, err
}
