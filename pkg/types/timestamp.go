package types

import (
	"database/sql/driver"
	"strconv"
	"time"
)

var timeformat = "2006-01-02 15:04:05"

// Timestamp be used to MySql timestamp converting.
type Timestamp int64

// UnmarshalJSON ...
func (jt *Timestamp) UnmarshalJSON(data []byte) (err error) {
	t, err := time.Parse(`"` + timeformat + `"`, string(data))
	if err != nil {
		return err
	}
	*jt = Timestamp(t.Unix())
	return nil
}

// MarshalJSON ...
func (jt *Timestamp) MarshalJSON() ([]byte, error) {
	timeStr := time.Unix(int64(*jt), 0).Format(timeformat)
	return []byte(`"` + timeStr + `"`), nil
}

// Scan scan time.
func (jt *Timestamp) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*jt = Timestamp(sc.Unix())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*jt = Timestamp(i)
	}
	return
}

// Value get time value.
func (jt Timestamp) Value() (driver.Value, error) {
	return time.Unix(int64(jt), 0), nil
}

// Time get time.
func (jt Timestamp) Time() time.Time {
	return time.Unix(int64(jt), 0)
}

// NowTimestamp ...
func NowTimestamp() Timestamp {
	return Timestamp(time.Now().Unix())
}
