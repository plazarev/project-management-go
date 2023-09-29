package common

import "encoding/json"

type TempID string

func (f *TempID) UnmarshalJSON(data []byte) error {
	var temp string
	if data[0] == QuotesByte {
		err := json.Unmarshal(data, &temp)
		if err != nil {
			return err
		}
	}
	*f = TempID(temp)

	return nil
}
