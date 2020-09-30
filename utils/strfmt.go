package utils

import (
	"fmt"
	"regexp"
)

// CheckNumber ...
func CheckNumber(str string) error {
	re := regexp.MustCompile(`\d+`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("要求全为数字")
	}

	return nil
}

// CheckUser ...
func CheckUser(str string) error {
	re := regexp.MustCompile(`\w+`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("要求字母，数字或下划线")
	}

	return nil
}

// CheckPhone ...
func CheckPhone(str string) error {
	re := regexp.MustCompile(`^1[3-9]\d{9}`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("错误的电话号码")
	}

	return nil
}

// CheckIP4 ...
func CheckIP4(str string) error {
	re := regexp.MustCompile(`(^[2][0-5][0-5]|^[1][0-9][0-9]|^[1-9][0-9]|^[0-9])\.([2][0-5][0-5]|[1][0-9][0-9]|[1-9][0-9]|[0-9])\.([2][0-5][0-5]|[1][0-9][0-9]|[1-9][0-9]|[0-9])\.([2][0-5][0-5]$|[1][0-9][0-9]$|[1-9][0-9]$|[0-9]$)`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("错误的ip地址")
	}

	return nil
}

// CheckPassword ...
func CheckPassword(str string) error {
	if len(str) < 10 || len(str) > 30 {
		return fmt.Errorf("密码字数在10~30之间")
	}

	specailChars := "[~`!@#\\$%\\^&\\*\\(\\)_\\+-=\\{}\\[]\\|\\\\:;\"'<>\\,\\./\\?]"
	alphabet := "[A-Za-z]"
	number := "[0-9]"

	//包含特殊字符和字母
	specialCharAndAlphabet := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", specailChars, alphabet, alphabet, specailChars)
	re := regexp.MustCompile(specialCharAndAlphabet)
	if string(re.Find([]byte(str))) == str {
		return nil
	}

	//包含特殊字符和数字
	specialCharAndNumber := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", specailChars, number, number, specailChars)
	re = regexp.MustCompile(specialCharAndNumber)
	if string(re.Find([]byte(str))) == str {
		return nil
	}

	//包含数字和字母
	alphabetAndNumber := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", alphabet, number, number, alphabet)
	re = regexp.MustCompile(alphabetAndNumber)
	if string(re.Find([]byte(str))) == str {
		return nil
	}

	return fmt.Errorf("密码必须包含特殊字符，数字，字母中的两种以上")
}
