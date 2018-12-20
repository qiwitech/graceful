// validate.go
package api_v1

import (
	"encoding/json"
)

func (req *TransferRequest) Validate() *json.UnmarshalTypeError {
	return nil
}

func (req *AccountSettingsRequest) Validate() *json.UnmarshalTypeError {
	return nil
}

func (req *UpdateSettingsRequest) Validate() *json.UnmarshalTypeError {
	return nil
}

func (req *PrevHashRequest) Validate() *json.UnmarshalTypeError {
	return nil
}

func (req *HistoryRequest) Validate() *json.UnmarshalTypeError {
	return nil
}

func (req *AccountsRequest) Validate() *json.UnmarshalTypeError {
	return nil
}
