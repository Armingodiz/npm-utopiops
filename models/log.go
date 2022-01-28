package models

import (
	"errors"
)

type Log struct {
	App         string
	Environment string
	LogGroup    string
	Exept       string
	Find        string
	From        int
}

func (l Log) IsValid() error {
	if l.App == "" {
		return errors.New("app is needed")
	}
	if l.Environment == "" {
		return errors.New("environment is needed")
	}
	return nil
}
