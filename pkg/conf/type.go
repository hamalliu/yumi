package conf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TimeDuration struct {
	time.Duration
}

func (d *TimeDuration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(strings.TrimSpace(string(text)))
	return err
}

type SpaceSize struct {
	Size int64
}

//parse "b", "kb", "mb", "gb", "tb" to int64
func (s *SpaceSize) UnmarshalText(text []byte) error {
	str := strings.TrimSpace(string(text))
	var err error
	if ok, _ := regexp.MatchString(`^\d+tb`, str); ok {
		s.Size, err = strconv.ParseInt(strings.TrimSuffix(str, "tb"), 10, 64)
		s.Size = s.Size * 2 << 40
		return err
	}
	if ok, _ := regexp.MatchString(`^\d+gb`, str); ok {
		s.Size, err = strconv.ParseInt(strings.TrimSuffix(str, "gb"), 10, 64)
		s.Size = s.Size * 2 << 30
		return err
	}
	if ok, _ := regexp.MatchString(`^\d+mb`, str); ok {
		s.Size, err = strconv.ParseInt(strings.TrimSuffix(str, "mb"), 10, 64)
		s.Size = s.Size * 2 << 20
		return err
	}
	if ok, _ := regexp.MatchString(`^\d+kb`, str); ok {
		s.Size, err = strconv.ParseInt(strings.TrimSuffix(str, "kb"), 10, 64)
		s.Size = s.Size * 2 << 10
		return err
	}
	if ok, _ := regexp.MatchString(`^\d+b`, str); ok {
		s.Size, err = strconv.ParseInt(strings.TrimSuffix(str, "b"), 10, 64)
		return err
	}

	return fmt.Errorf("空间单位未识别，%s", str)
}
