package publisher

import (
	remote "github.com/mkozhukh/go-remote"
)

type MoveParams struct {
	ID        int  `json:"id"`
	ParentID  int  `json:"parent"`
	Index     int  `json:"index"`
	TargetID  int  `json:"targetId,omitempty"`
	ProjectID int  `json:"project,omitempty"`
	Reverse   bool `json:"reverse,omitempty"`
}

type IPublisherAPI interface {
	AddItem(ctx PublisherContext, id int) error
	UpdateItem(ctx PublisherContext, id int) error
	DeleteItem(ctx PublisherContext, id int, children []int) error

	AddProject(ctx PublisherContext, id int) error
	UpdateProject(ctx PublisherContext, id int) error
	DeleteProject(ctx PublisherContext, id int, children []int) error

	GetRemoteAPI() *remote.Server
}

type IMovePublisher interface {
	ChangeProject(ctx PublisherContext, id, projectId int) error
}

type ITreePublisher interface {
	Move(ctx PublisherContext, id int, params MoveParams) error
}

type BasePublisher struct {
	api        *remote.Server
	widgetType WidgetType
}

func (b *BasePublisher) GetRemoteAPI() *remote.Server {
	return b.api
}
