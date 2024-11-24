package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ParameterType string

const (
	TypeInt     ParameterType = "int"
	TypeDouble  ParameterType = "double"
	TypeFloat   ParameterType = "float"
	TypeString  ParameterType = "string"
	TypeChar    ParameterType = "char"
	TypeBool    ParameterType = "bool"
	TypeIntArr    ParameterType = "int[]"
	TypeDoubleArr ParameterType = "double[]"
	TypeFloatArr  ParameterType = "float[]"
	TypeStringArr ParameterType = "string[]"
	TypeCharArr   ParameterType = "char[]"
	TypeBoolArr   ParameterType = "bool[]"
)

type Parameter struct {
	Name     string       `json:"name"`
	Type     ParameterType `json:"type" binding:"required,oneof=int double float string char bool int[] double[] float[] string[] char[] bool[]"`
	Position int         `json:"position"`
}

type Parameters []Parameter

func (p Parameters) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Parameters) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &p)
} 