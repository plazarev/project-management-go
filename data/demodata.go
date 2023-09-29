package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type IDemodataProvider interface {
	Up(ctx *DBContext)
	Down(ctx *DBContext)
}

type AppDataProvider struct {
	folder string
}

type DemodataProvider struct{}

func NewDemodataProvider() *DemodataProvider {
	return &DemodataProvider{}
}

func NewAppDataProvider(folder string) *AppDataProvider {
	return &AppDataProvider{folder}
}

func (d AppDataProvider) Up(ctx *DBContext) {
	must(InitializeDemodata[User](ctx, path.Join(d.folder, "users.json")))
	must(InitializeDemodata[Project](ctx, path.Join(d.folder, "projects.json")))
	must(d.items(ctx))
}

func (d AppDataProvider) Down(ctx *DBContext) {
	must(ctx.DB.Delete(&Item{}, "1 = 1").Error)
	must(ctx.DB.Delete(&ItemUser{}, "1 = 1").Error)
	must(ctx.DB.Delete(&Project{}, "1 = 1").Error)
	must(ctx.DB.Delete(&User{}, "1 = 1").Error)
}

func (d AppDataProvider) items(ctx *DBContext) error {
	var items []Item
	err := ParseJSON(path.Join(d.folder, "items.json"), &items)
	if err != nil {
		return err
	}

	indexMap := make(map[string]int)
	for i := range items {
		item := &items[i]
		key := fmt.Sprintf("%d:%d", item.ProjectID, item.ParentID)
		item.Index = indexMap[key]
		indexMap[key]++
	}

	return ctx.DB.Create(&items).Error
}

func (d DemodataProvider) Restore(ctx *DBContext, providers ...IDemodataProvider) {
	for i := range providers {
		providers[i].Down(ctx)
		providers[i].Up(ctx)
	}
}

func (d DemodataProvider) Down(ctx *DBContext, providers ...IDemodataProvider) {
	for i := range providers {
		providers[i].Down(ctx)
	}
}

func (d DemodataProvider) Up(ctx *DBContext, providers ...IDemodataProvider) {
	for i := range providers {
		providers[i].Up(ctx)
	}
}

func InitializeDemodata[T any](ctx *DBContext, path string) error {
	data := make([]T, 0)
	err := ParseJSON(path, &data)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	err = ctx.DB.Create(&data).Error

	if Debug {
		if err == nil {
			fmt.Printf("--- Initialized %d rows of %T\n", len(data), data)
		}
	}

	return err
}

func ParseJSON(path string, dest any) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &dest)
}
