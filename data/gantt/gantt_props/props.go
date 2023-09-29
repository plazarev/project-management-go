package gantt_props

type GanttItemProps struct {
	GanttTaskType string `gorm:"gantt_task_type" json:"gantt_task_type"`
}

type GanttProjectProps struct{}
