package repository

import (
	"fmt"
	"strings"
)

type SetQueryGenerator struct {
	SetValues []string
	Args      []interface{}
	ArgId     int
}

func NewSetQueryGenerator() *SetQueryGenerator {
	return &SetQueryGenerator{
		SetValues: make([]string, 0),
		Args:      make([]interface{}, 0),
		ArgId:     1,
	}
}

func (qg *SetQueryGenerator) Add(field string, value interface{}) error {
	qg.SetValues = append(qg.SetValues, fmt.Sprintf("%s=$%d", field, qg.ArgId))
	qg.Args = append(qg.Args, value)
	qg.ArgId = qg.ArgId + 1
	return nil
}

func (qg *SetQueryGenerator) GetSetQuery() (string, error) {
	setQuery := strings.Join(qg.SetValues, ", ")
	return setQuery, nil
}
