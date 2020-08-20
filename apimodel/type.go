package apimodel

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	//DayTimeFormat ...
	DayTimeFormat = "2006-01-02"
	//SecondTimeFormat ...
	SecondTimeFormat = "2006-01-02 15-04-05"
	//MillisecondTimeFormat ...
	MillisecondTimeFormat = "2006-01-02 15-04-05.999"
	//MicrosecondTimeFormat ...
	MicrosecondTimeFormat = "2006-01-02 15-04-05.999999"
	//NanosecondTimeFormat ...
	NanosecondTimeFormat = "2006-01-02 15-04-05.999999999"
)

//DayTime ...
type DayTime time.Time

//UnmashalJSON 用于json Unmashal
func (day *DayTime) UnmashalJSON(data []byte) (err error) {
	t, err := time.Parse(DayTimeFormat, string(data))
	if err != nil {
		return fmt.Errorf("时间格式错误")
	}
	*day = DayTime(t)
	return nil
}

//MashalJSON 用于json Mashal
func (day DayTime) MashalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DayTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(day).AppendFormat(b, DayTimeFormat)
	b = append(b, '"')
	return b, nil
}

//Value 实现 database driver 的 Valuer 接口
func (day DayTime) Value() (driver.Value, error) {
	return time.Time(day), nil
}

//Scan sqlx的Scan接口，用于反序列化到结构体
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

//SecondTime ...
type SecondTime time.Time

//UnmashalJSON 用于json Unmashal
func (second *SecondTime) UnmashalJSON(data []byte) (err error) {
	t, err := time.Parse(SecondTimeFormat, string(data))
	if err != nil {
		return fmt.Errorf("时间格式错误")
	}
	*second = SecondTime(t)
	return nil
}

//MashalJSON 用于json Mashal
func (second SecondTime) MashalJSON() ([]byte, error) {
	b := make([]byte, 0, len(SecondTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(second).AppendFormat(b, SecondTimeFormat)
	b = append(b, '"')
	return b, nil
}

//Value 实现 database driver 的 Valuer 接口
func (second SecondTime) Value() (driver.Value, error) {
	return time.Time(second), nil
}

//Scan sqlx的Scan接口，用于反序列化到结构体
func (second *SecondTime) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch src.(type) {
	case time.Time:
		*second = SecondTime(src.(time.Time))
		return nil
	case []byte:
		t, err := time.Parse(SecondTimeFormat, string(src.([]byte)))
		*second = SecondTime(t)
		return err
	case string:
		t, err := time.Parse(SecondTimeFormat, src.(string))
		*second = SecondTime(t)
		return err
	default:
		return fmt.Errorf("时间格式错误，%v", src)
	}
}
