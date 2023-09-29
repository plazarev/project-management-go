package publisher

type IPublisherEvent interface {
	GetFrom() int
	GetWidget() WidgetType
	IsSelf() bool
}

type EventBase struct {
	Type   string     `json:"type"`
	From   int        `json:"-"`
	Self   bool       `json:"-"`
	Widget WidgetType `json:"-"`
}

func (b EventBase) GetFrom() int {
	return b.From
}

func (b EventBase) GetWidget() WidgetType {
	return b.Widget
}

func (b EventBase) IsSelf() bool {
	return b.Self
}
