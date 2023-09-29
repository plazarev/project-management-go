package items

import (
	"project-manager-go/data"
)

type Item struct {
	data.Item
}

func (t *Item) PutItem(item data.Item) {
	*t = Item{
		Item: item,
	}
}

func (t Item) FillItem(item *data.Item) {
	*item = t.Item
}
