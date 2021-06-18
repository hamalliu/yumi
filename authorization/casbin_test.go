package authorization

import (
	"fmt"
	"testing"

	"github.com/casbin/casbin"
)

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func TestRBAC(t *testing.T) {
	e := casbin.NewEnforcer("./rbac_model.conf", "./rbac_policy.csv")

	check(e, "alice", "data1", "read")
}
