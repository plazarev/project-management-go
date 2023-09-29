package gantt

import (
	"project-manager-go/data"
)

type GanttLink struct {
	ID     int `json:"id"`
	Source int `json:"source"`
	Target int `json:"target"`
	Type   int `json:"type"`
}

type links struct {
	data.ItemsStore
}

func (s *links) GetAll(ctx *data.DBContext) ([]GanttLink, error) {
	links := make([]GanttLink, 0)
	err := ctx.DB.
		Find(&links).
		Error

	return links, err
}

func (d *links) GetOne(ctx *data.DBContext, id int) (GanttLink, error) {
	link := GanttLink{}
	err := ctx.DB.Take(&link, id).Error
	return link, err
}

func (d *links) Add(ctx *data.DBContext, link GanttLink) (int, error) {
	err := ctx.DB.Create(&link).Error
	return link.ID, err
}

func (d *links) Update(ctx *data.DBContext, id int, upd GanttLink) error {
	link, err := d.GetOne(ctx, id)
	if err != nil {
		return err
	}

	link.Source = upd.Source
	link.Target = upd.Target
	link.Type = upd.Type

	err = ctx.DB.Save(&link).Error

	return err
}

func (d *links) Delete(ctx *data.DBContext, id int) error {
	err := ctx.DB.Delete(&GanttLink{}, id).Error
	return err
}
