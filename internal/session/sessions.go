package session

import (
	"strings"
	"sync"
	"time"
)

type Sessions struct {
	rwmux          sync.RWMutex
	userList       map[string]Session
	expireDuration time.Duration
}

type Session struct {
	User     string
	RealName string
	Token    string
	power    map[string]bool

	LastLoginTime time.Time
	LoginTime     time.Time
	LastHeart     time.Time
}

var sssns *Sessions

func initSessions(ed time.Duration) {
	sssns.userList = make(map[string]Session)
	sssns.expireDuration = ed
	go sssns.checkExpire()
}

func LoadSeesion(user string) error {
	sssns.rwmux.Lock()
	defer sssns.rwmux.Unlock()
	//登录成功加入队列

	return nil
}

func UpdateLastHeart(user string) {
	sssns.rwmux.Lock()
	defer sssns.rwmux.Unlock()
	//验签成功后，更新心跳
	if _tmp, ok := sssns.userList[user]; ok {
		_tmp.LastHeart = time.Now()
	}
	return
}

func GetUser(user string) (Session, bool) {
	sssns.rwmux.RLock()
	defer sssns.rwmux.RUnlock()
	_empty := Session{}
	if u, ok := sssns.userList[user]; ok {
		_now := time.Now().Add(-1 * sssns.expireDuration)
		if u.LastHeart.Before(_now) {
			return _empty, false
		} else {
			return u, true
		}
	}
	return _empty, false
}

func HavePower(user, code string) bool {
	codes := strings.Split(code, ",")
	for i := range codes {
		if sssns.userList[user].power[codes[i]] {
			return true
		}
	}

	return false
}

func Remove(user string) {
	sssns.rwmux.Lock()
	defer sssns.rwmux.Unlock()
	//退出登录或登录超时
	if _, ok := sssns.userList[user]; ok {
		delete(sssns.userList, user)
	}
}

func (m *Sessions) checkExpire() {
	for {
		_now := time.Now().Add(-1 * m.expireDuration)
		for _, v := range m.userList {
			if v.LastHeart.Before(_now) {
				Remove(v.User)
			}
		}
		time.Sleep(time.Minute * time.Duration(1))
	}
}
