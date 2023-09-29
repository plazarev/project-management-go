package data

type IItem interface {
	FillItem(*Item)
	PutItem(Item)
}

type IItemWrapper[T any] interface {
	IItem
	*T
}

type IItemsList interface {
	PutItems([]Item)
}

type ItemsList[T any, PT IItemWrapper[T]] struct {
	arr []T
}

func NewItemsList[T any, PT IItemWrapper[T]]() *ItemsList[T, PT] {
	return &ItemsList[T, PT]{
		arr: make([]T, 0),
	}
}

func (l *ItemsList[T, I]) GetArray() []T {
	return l.arr
}

func (l *ItemsList[T, PT]) PutItems(items []Item) {
	l.arr = make([]T, len(items))
	for i := range items {
		// create a pointer to a non-nil instance
		item := PT(new(T))
		// fill with the properties
		item.PutItem(items[i])
		l.arr[i] = *item
	}
}
