package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func NewLogger(debug bool) (logger *zap.Logger, err error) {
	var l *zap.Logger

	if debug {
		l, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("initialize development logger error: %v", err)
		}
	}

	l, err = zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("initialize production logger error: %v", err)
	}

	return l, nil
}
