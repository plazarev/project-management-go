package publisher

type WidgetType int

const (
	WidgetTodo WidgetType = iota + 1
	WidgetScheduler
	WidgetKanban
	WidgetGantt
)

type PublisherContext struct {
	UserID     int
	DeviceID   int
	FromWidget WidgetType
}
