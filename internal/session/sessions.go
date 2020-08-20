package session

import (
	"strings"
	"sync"
	"time"
)

//Sessions ...
type Sessions struct {
	rwmux          sync.RWMutex
	userList       map[string]Session
	expireDuration time.Duration
}

//Session ...
type Session struct {
	UserID   string
	UserName string
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

//LoadSeesion ...
func LoadSeesion(user string) error {
	sssns.rwmux.Lock()
	defer sssns.rwmux.Unlock()
	//登录成功加入队列

	return nil
}

//UpdateLastHeart ...
func UpdateLastHeart(user string) {
	sssns.rwmux.Lock()
	defer sssns.rwmux.Unlock()
	//验签成功后，更新心跳
	if tmp, ok := sssns.userList[user]; ok {
		tmp.LastHeart = time.Now()
	}
	return
}

//GetUser ...
func GetUser(user string) (Session, bool) {
	sssns.rwmux.RLock()
	defer sssns.rwmux.RUnlock()
	empty := Session{}
	if u, ok := sssns.userList[user]; ok {
		now := time.Now().Add(-1 * sssns.expireDuration)
		if u.LastHeart.Before(now) {
			return empty, false
		}
		
		return u, true
	}
	return empty, false
}

//HavePower ...
func HavePower(user, code string) bool {
	codes := strings.Split(code, ",")
	for i := range codes {
		if sssns.userList[user].power[codes[i]] {
			return true
		}
	}

	return false
}

//Remove ...
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
		now := time.Now().Add(-1 * m.expireDuration)
		for _, v := range m.userList {
			if v.LastHeart.Before(now) {
				Remove(v.UserID)
			}
		}
		time.Sleep(time.Minute * time.Duration(1))
	}
}
