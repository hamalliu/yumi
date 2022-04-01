package regexputility

import (
	"fmt"
	"regexp"
)

// RegexpNumber ...
func RegexpNumber(str string) bool {
	re := regexp.MustCompile(`\d+`)
	return string(re.Find([]byte(str))) == str
}

// RegexpUser ...
func RegexpUser(str string) bool {
	if len(str) < 10 || len(str) > 30 {
		return false
	}

	re := regexp.MustCompile(`\w+`)
	return string(re.Find([]byte(str))) == str
}

// RegexpPhone ...
func RegexpPhone(str string) bool {
	re := regexp.MustCompile(`^1[3-9]\d{9}`)
	return string(re.Find([]byte(str))) == str
}

// RegexpPassword ...
func RegexpPassword(str string) bool {
	if len(str) < 10 || len(str) > 30 {
		return false
	}

	specailChars := "[~`!@#\\$%\\^&\\*\\(\\)_\\+-=\\{}\\[]\\|\\\\:;\"'<>\\,\\./\\?]"
	alphabet := "[A-Za-z]"
	number := "[0-9]"

	//包含特殊字符和字母
	specialCharAndAlphabet := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", specailChars, alphabet, alphabet, specailChars)
	re := regexp.MustCompile(specialCharAndAlphabet)
	if string(re.Find([]byte(str))) == str {
		return true
	}

	//包含特殊字符和数字
	specialCharAndNumber := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", specailChars, number, number, specailChars)
	re = regexp.MustCompile(specialCharAndNumber)
	if string(re.Find([]byte(str))) == str {
		return true
	}

	//包含数字和字母
	alphabetAndNumber := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", alphabet, number, number, alphabet)
	re = regexp.MustCompile(alphabetAndNumber)
	return string(re.Find([]byte(str))) == str
}


// RegexpIP4 ...
func RegexpIP4(str string) bool {
	re := regexp.MustCompile(`(^[2][5][0-5]|^[2][0-4][0-9]|^[1][0-9][0-9]|^[1-9][0-9]|^[0-9])\.([2][5][0-5]|[2][0-4][0-9]|[1][0-9][0-9]|[1-9][0-9]|[0-9])\.([2][5][0-5]|[2][0-4][0-9]|[1][0-9][0-9]|[1-9][0-9]|[0-9])\.([2][5][0-5]$|[2][0-4][0-9]$|[1][0-9][0-9]$|[1-9][0-9]$|[0-9]$)`)
	fmt.Println(string(re.Find([]byte(str))))
	return string(re.Find([]byte(str))) == str
}

// RegexpIP6 ...
func RegexpIP6(str string) bool {
	re := regexp.MustCompile(`(([\da-fA-F]{1,4}):([\da-fA-F]{1,4}):([\da-fA-F]{1,4}):([\da-fA-F]{1,4}):([\da-fA-F]{1,4}):([\da-fA-F]{1,4}):([\da-fA-F]{1,4}):([\da-fA-F]{1,4}))$`)
	return string(re.Find([]byte(str))) == str
}