package apimodel

import (
	"database/sql/driver"
	"strconv"
	"time"
)

// Timestamp be used to MySql timestamp converting.
type Timestamp int64

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
