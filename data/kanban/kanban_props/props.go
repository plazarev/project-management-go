package kanban_props

type KanbanItemProps struct {
	Kanban_ColumnID  int `gorm:"column:kanban_column_id" json:"kanban_column_id"`
	Kanban_CardIndex int `gorm:"column:kanban_index" json:"kanban_card_index"`
}

type KanbanProjectProps struct {
	Kanban_RowIndex int `gorm:"column:kanban_row_index;" json:"kanban_row_index"`
}
