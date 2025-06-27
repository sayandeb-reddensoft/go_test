package utils

import (
	"encoding/json"
	"errors"
)

func BindByteResponseToStruct(data []byte, responseStruct interface{}) error {
	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, responseStruct); err != nil {
		return errors.New("failed to bind byte response to struct: " + err.Error())
	}

	return nil
}

func ConvertDataToBytes(data any) ([]byte, error) {
	convertedData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return convertedData, nil
}
