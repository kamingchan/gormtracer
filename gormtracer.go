package gormtracer

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

var (
	_ gorm.Plugin = NewGormTracer()
)

const (
	HookBefore = "opentracing:startSpan"
	HookAfter  = "opentracing:finishSpan"

	OperationName = "MySQL"

	InstanceSpanKey = "opentracing-span"

	KeyRowsAffected = "db.rows_affected"
)

type GormTracer struct {
}

func (g *GormTracer) Name() string {
	return "GormTracer"
}

func (g *GormTracer) Initialize(db *gorm.DB) error {
	db.Callback().Create().Before("*").Register(HookBefore, g.startSpan)
	db.Callback().Query().Before("*").Register(HookBefore, g.startSpan)
	db.Callback().Delete().Before("*").Register(HookBefore, g.startSpan)
	db.Callback().Update().Before("*").Register(HookBefore, g.startSpan)
	db.Callback().Row().Before("*").Register(HookBefore, g.startSpan)
	db.Callback().Raw().Before("*").Register(HookBefore, g.startSpan)

	db.Callback().Create().After("*").Register(HookAfter, g.finishSpan)
	db.Callback().Query().After("*").Register(HookAfter, g.finishSpan)
	db.Callback().Delete().After("*").Register(HookAfter, g.finishSpan)
	db.Callback().Update().After("*").Register(HookAfter, g.finishSpan)
	db.Callback().Row().After("*").Register(HookAfter, g.finishSpan)
	db.Callback().Raw().After("*").Register(HookAfter, g.finishSpan)

	return nil
}

func (g *GormTracer) startSpan(db *gorm.DB) {
	sp, _ := opentracing.StartSpanFromContext(db.Statement.Context, OperationName)
	db.InstanceSet(InstanceSpanKey, sp)
}

func (g *GormTracer) finishSpan(db *gorm.DB) {
	_sp, ok := db.InstanceGet(InstanceSpanKey)
	if !ok {
		return
	}
	sp, ok := _sp.(opentracing.Span)
	if !ok {
		return
	}

	ext.DBType.Set(sp, OperationName)
	// from: gorm.io/gorm/callbacks.go
	ext.DBStatement.Set(sp, db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...))
	sp.LogFields(log.Int64(KeyRowsAffected, db.RowsAffected))
	err := db.Error
	if err != nil {
		ext.Error.Set(sp, true)
		sp.LogFields(log.Error(err))
	}

	sp.Finish()
}

func NewGormTracer() *GormTracer {
	return &GormTracer{}
}
