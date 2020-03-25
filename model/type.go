package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	DayTimeFormat         = "2006-01-02"
	SecondTimeFormat      = "2006-01-02 15-04-05"
	MillisecondTimeFormat = "2006-01-02 15-04-05.999"
	MicrosecondTimeFormat = "2006-01-02 15-04-05.999999"
	NanosecondTimeFormat  = "2006-01-02 15-04-05.999999999"
)

type DayTime time.Time

func (day *DayTime) UnmashalJSON(data []byte) (err error) {
	t, err := time.Parse(DayTimeFormat, string(data))
	if err != nil {
		return fmt.Errorf("时间格式错误")
	}
	*day = DayTime(t)
	return nil
}

func (day DayTime) MashalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DayTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(day).AppendFormat(b, DayTimeFormat)
	b = append(b, '"')
	return b, nil
}

func (day DayTime) Value() (driver.Value, error) {
	return time.Time(day), nil
}

func (day *DayTime) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch src.(type) {
	case time.Time:
		*day = DayTime(src.(time.Time))
		return nil
	case []byte:
		t, err := time.Parse(DayTimeFormat, string(src.([]byte)))
		*day = DayTime(t)
		return err
	case string:
		t, err := time.Parse(DayTimeFormat, src.(string))
		*day = DayTime(t)
		return err
	default:
		return fmt.Errorf("时间格式错误，%v", src)
	}
}

type SecondTime time.Time

func (day *SecondTime) UnmashalJSON(data []byte) (err error) {
	t, err := time.Parse(SecondTimeFormat, string(data))
	if err != nil {
		return fmt.Errorf("时间格式错误")
	}
	*day = SecondTime(t)
	return nil
}

func (day SecondTime) MashalJSON() ([]byte, error) {
	b := make([]byte, 0, len(SecondTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(day).AppendFormat(b, SecondTimeFormat)
	b = append(b, '"')
	return b, nil
}

func (day SecondTime) Value() (driver.Value, error) {
	return time.Time(day), nil
}

func (day *SecondTime) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch src.(type) {
	case time.Time:
		*day = SecondTime(src.(time.Time))
		return nil
	case []byte:
		t, err := time.Parse(SecondTimeFormat, string(src.([]byte)))
		*day = SecondTime(t)
		return err
	case string:
		t, err := time.Parse(SecondTimeFormat, src.(string))
		*day = SecondTime(t)
		return err
	default:
		return fmt.Errorf("时间格式错误，%v", src)
	}
}
