package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//TimeDuration ...
type TimeDuration time.Duration

//Duration ...
func (d TimeDuration) Duration() time.Duration {
	return time.Duration(d)
}

//UnmarshalText ...
func (d *TimeDuration) UnmarshalText(text []byte) error {
	dur, err := time.ParseDuration(strings.TrimSpace(string(text)))
	*d = TimeDuration(dur)
	return err
}

//SpaceSize ...
type SpaceSize int64

//Size ...
func (s SpaceSize) Size() int64 {
	return int64(s)
}

//UnmarshalText ...
//parse "b", "kb", "mb", "gb", "tb" to int64
func (s *SpaceSize) UnmarshalText(text []byte) error {
	str := strings.TrimSpace(string(text))
	if ok, _ := regexp.MatchString(`^\d+tb`, str); ok {
		size, err := strconv.ParseInt(strings.TrimSuffix(str, "tb"), 10, 64)
		*s = SpaceSize(size * 2 << 40)
		return err
	}
	if ok, _ := regexp.MatchString(`^\d+gb`, str); ok {
		size, err := strconv.ParseInt(strings.TrimSuffix(str, "gb"), 10, 64)
		*s = SpaceSize(size * 2 << 30)
		return err
	}
	if ok, _ := regexp.MatchString(`^\d+mb`, str); ok {
		size, err := strconv.ParseInt(strings.TrimSuffix(str, "mb"), 10, 64)
		*s = SpaceSize(size * 2 << 20)
		return err
	}
	if ok, _ := regexp.MatchString(`^\d+kb`, str); ok {
		size, err := strconv.ParseInt(strings.TrimSuffix(str, "kb"), 10, 64)
		*s = SpaceSize(size * 2 << 10)
		return err
	}
	if ok, _ := regexp.MatchString(`^\d+b`, str); ok {
		size, err := strconv.ParseInt(strings.TrimSuffix(str, "b"), 10, 64)
		*s = SpaceSize(size)
		return err
	}

	return fmt.Errorf("空间单位未识别，%s", str)
}

func (s *SpaceSize) String() string {
	size := int64(*s)
	if size > 2<<40 {
		return fmt.Sprintf("%dtb", size/2<<40)
	}
	if size > 2<<30 {
		return fmt.Sprintf("%dgb", size/2<<30)
	}
	if size > 2<<20 {
		return fmt.Sprintf("%dmb", size/2<<20)
	}
	if size > 2<<10 {
		return fmt.Sprintf("%dkb", size/2<<10)
	}

	return fmt.Sprintf("%db", size)
}

//ArrayString ...
type ArrayString []string

//IndexOf ...
func (as ArrayString) IndexOf(elem string) int {
	for i, e := range as {
		if elem == e {
			return i
		}
	}

	return -1
}

type JsonBytes []byte

func (m *JsonBytes) MarshalJSON() ([]byte, error) {
	return *m, nil
}

func (m *JsonBytes) UnmarshalJSON(data []byte) (err error) {
	if data[0] == '"' {
		*m = data[1:len(data)-1]
	} else {
		*m = data
	}
	return nil
}
