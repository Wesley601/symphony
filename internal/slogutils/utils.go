package slogutils

import (
	"fmt"
	"log/slog"
)

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func InstanceName(i any) slog.Attr {
	return slog.String("instance", fmt.Sprintf("%T", i))
}
