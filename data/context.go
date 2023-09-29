package data

import (
	"bytes"
	"fmt"
	"runtime/debug"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

type DBContext struct {
	DB     *gorm.DB
	tx     *gorm.DB
	level  int
	closed int
}

func (db *DBContext) inc() *DBContext {
	db.level += 1
	return db
}

func NewTCtx(ctx *DBContext) *DBContext {
	if ctx != nil {
		return ctx.inc()
	}

	tx := db.Begin()
	return &DBContext{DB: tx, tx: tx}
}

func NewCtx(ctx *DBContext) *DBContext {
	if ctx != nil {
		return ctx.inc()
	}

	return &DBContext{DB: db, tx: nil}
}

func (db *DBContext) End(err error) error {
	if db.tx == nil || db.closed == 1 {
		return err
	}

	if err != nil {
		logError(err)
		db.tx.Rollback()
		db.closed = 1
		return err
	}

	if db.level == 0 {
		err = db.tx.Commit().Error
		if err != nil {
			color.Red(err.Error())
			fmt.Println(niceStack(debug.Stack()))
		} else {
			color.Yellow("done")
		}
	} else {
		db.level -= 1
	}

	return err
}

func logError(err error) {
	if err != nil && Debug {
		color.Red(err.Error())
		fmt.Println(niceStack(debug.Stack()))
	}
}

func niceStack(data []byte) string {
	out := make([]byte, 0, len(data))

	lines := bytes.Split(data, []byte("\n"))
	for _, x := range lines {
		if bytes.HasPrefix(x, []byte("\t")) {
			if bytes.Contains(x, []byte("/go/src/")) || bytes.Contains(x, []byte("/go/pkg/")) {
				continue
			}
			out = append(out, x...)
			out = append(out, []byte("\n")...)

		}
	}

	return string(out)
}
