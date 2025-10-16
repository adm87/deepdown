package deepdown

import (
	"context"
	"log/slog"
)

type CtxKey string

const (
	CtxApplicationRoot CtxKey = "application_root"
	CtxAssetsRoot      CtxKey = "assets_root"
	CtxEmbeddedRoot    CtxKey = "embedded_root"
)

type Context interface {
	Ctx() context.Context
	Logger() *slog.Logger
	Set(key CtxKey, value any)
	Get(key CtxKey) any
}

type contextImpl struct {
	ctx    context.Context
	logger *slog.Logger
	values map[CtxKey]any
}

func NewContext() *contextImpl {
	return &contextImpl{
		ctx:    context.Background(),
		logger: slog.Default(),
		values: make(map[CtxKey]any),
	}
}

func (c *contextImpl) Ctx() context.Context {
	return c.ctx
}

func (c *contextImpl) Set(key CtxKey, value any) {
	c.values[key] = value
}

func (c *contextImpl) Get(key CtxKey) any {
	return c.values[key]
}

func (c *contextImpl) Logger() *slog.Logger {
	return c.logger
}
