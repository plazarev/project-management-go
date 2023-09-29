package common

import (
	"encoding/json"
)

const QuotesByte = 34

type FuzzyInt int
type TID int

func (f *FuzzyInt) UnmarshalJSON(data []byte) error {
	var err error
	var temp int
	if data[0] == QuotesByte {
		// empty string => 0
		if len(data) > 2 {
			err = json.Unmarshal(data[1:len(data)-1], &temp)
		}
	} else {
		err = json.Unmarshal(data, &temp)
	}
	*f = FuzzyInt(temp)
	return err
}

func (t *TID) UnmarshalJSON(data []byte) error {
	var err error
	var temp int
	if data[0] == QuotesByte {
		// empty string => 0
		if len(data) > 2 {
			if len(data) > 8 {
				s := string(data)[1:8]
				if s == "temp://" {
					return nil
				}
			}
			err = json.Unmarshal(data[1:len(data)-1], &temp)
		}
	} else {
		err = json.Unmarshal(data, &temp)
	}
	*t = TID(temp)
	return err
}
