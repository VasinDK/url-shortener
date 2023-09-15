package checkerr

import (
	"log/slog"
	"os"
)

const (
	KeyError = "error"
)

type Checkerr struct {
	log *slog.Logger
}

type Checkerrer interface {
	ErrExit(string, error)
	Err(string, error)
}

func New(log *slog.Logger) *Checkerr {
	return &Checkerr{log: log}
}

func (c *Checkerr) ErrExit(s1 string, err error) {
	if err != nil {
		c.log.Error(s1, slog.String(KeyError, err.Error()))
		os.Exit(1)
	}
}

func (c *Checkerr) Err(s1 string, err error) {
	if err != nil {
		c.log.Error(s1, slog.String(KeyError, err.Error()))
	}
}
