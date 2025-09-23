package models

import (
	"fmt"
	"time"
)

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(d).Format("2006-01-02"))
	return []byte(formatted), nil
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	parsed, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}
	*d = DateOnly(parsed)
	return nil
}

func (d DateOnly) ToTime() time.Time {
	return time.Time(d)
}
