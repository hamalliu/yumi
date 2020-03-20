package pretreatment

import (
	"fmt"
	"regexp"
)

func CheckNumber(str string) error {
	re := regexp.MustCompile(`\d+`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("要求全为数字")
	}

	return nil
}

func CheckUser(str string) error {
	re := regexp.MustCompile(`\w+`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("要求字母，数字或下划线")
	}

	return nil
}

func CheckPhone(str string) error {
	re := regexp.MustCompile(`^1[3-9]\d{9}`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("错误的电话号码")
	}

	return nil
}

func CheckIp(str string) error {
	re := regexp.MustCompile(`(^[2][0-5][0-5]|^[1][0-9][0-9]|^[1-9][0-9]|^[0-9])\.([2][0-5][0-5]|[1][0-9][0-9]|[1-9][0-9]|[0-9])\.([2][0-5][0-5]|[1][0-9][0-9]|[1-9][0-9]|[0-9])\.([2][0-5][0-5]$|[1][0-9][0-9]$|[1-9][0-9]$|[0-9]$)`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("错误的ip地址")
	}

	return nil
}
