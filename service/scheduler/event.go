package scheduler

import (
	"project-manager-go/common"
	"project-manager-go/data"
	"time"
)

type Event struct {
	ID         common.TID  `json:"id"`
	CalendarID common.TID  `json:"type"`
	AllDay     bool        `json:"allDay"`
	Text       string      `json:"text"`
	Details    string      `json:"details"`
	StartDate  *time.Time  `json:"start_date"`
	EndDate    *time.Time  `json:"end_date"`

	Index int `json:"-"`
}

func (e Event) FillItem(item *data.Item) {
	item.ProjectID = int(e.CalendarID)
	item.AllDay = e.AllDay
	item.Text = e.Text
	item.Description = e.Details
	item.StartDate = e.StartDate
	item.EndDate = e.EndDate
}

func (e *Event) PutItem(item data.Item) {
	e.ID = common.TID(item.ID)
	e.CalendarID = common.TID(item.ProjectID)
	e.AllDay = item.AllDay
	e.Text = item.Text
	e.Details = item.Description
	e.StartDate = item.StartDate
	e.EndDate = item.EndDate
}
