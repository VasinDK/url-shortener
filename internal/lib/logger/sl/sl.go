package sl

import (
	"golang.org/x/exp/slog"
)

// slog свое дополнение

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
