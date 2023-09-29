package todo

import "project-manager-go/data"

var tagValues = []string{
	"start",
	"#end",
}

type tags struct{}

func (s *tags) GetAll(dbCtx *data.DBContext) ([]string, error) {
	return tagValues, nil
}
