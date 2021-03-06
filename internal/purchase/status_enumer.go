// Code generated by "enumer -type=Status -json -sql internal/purchase/purchase.go"; DO NOT EDIT.

//
package purchase

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const _StatusName = "CreatedPendingSuccessFailure"

var _StatusIndex = [...]uint8{0, 7, 14, 21, 28}

func (i Status) String() string {
	if i >= Status(len(_StatusIndex)-1) {
		return fmt.Sprintf("Status(%d)", i)
	}
	return _StatusName[_StatusIndex[i]:_StatusIndex[i+1]]
}

var _StatusValues = []Status{0, 1, 2, 3}

var _StatusNameToValueMap = map[string]Status{
	_StatusName[0:7]:   0,
	_StatusName[7:14]:  1,
	_StatusName[14:21]: 2,
	_StatusName[21:28]: 3,
}

// StatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func StatusString(s string) (Status, error) {
	if val, ok := _StatusNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Status values", s)
}

// StatusValues returns all values of the enum
func StatusValues() []Status {
	return _StatusValues
}

// IsAStatus returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Status) IsAStatus() bool {
	for _, v := range _StatusValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Status
func (i Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Status
func (i *Status) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Status should be a string, got %s", data)
	}

	var err error
	*i, err = StatusString(s)
	return err
}

func (i Status) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *Status) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a byte slice")
		}

		str = string(bytes[:])
	}

	val, err := StatusString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
