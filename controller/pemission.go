package controller

type Pemission struct {
	UserIdCode map[string] /*userid*/ map[string]bool /*code*/
}

func (pem *Pemission) HavePower(user, code string) bool {
	if code == "" {
		return false
	}

	if !pem.UserIdCode[user][code] {
		return false
	}

	return true
}

func (pem *Pemission) SetUserCode(user string, codes []string) {
	for i := range codes {
		if pem.UserIdCode[user] == nil {
			pem.UserIdCode[user] = make(map[string]bool)
		}
		pem.UserIdCode[user][codes[i]] = true
	}
}
