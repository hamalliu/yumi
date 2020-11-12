package middleware

import (
	"net/http"
	"reflect"
	"sort"
	"strings"
	"sync"

	"yumi/pkg/gin"
	"yumi/pkg/gin/valuer"
	"yumi/pkg/log"

	"github.com/casbin/casbin"
)

// Logic is the logical operation (AND/OR) used in permission checks
// in case multiple permissions or roles are specified.
type Logic int

const (
	// AND ...
	AND Logic = iota
	// OR ...
	OR
)

// Option is used to change some default behaviors.
type Option interface {
	apply(*options)
}

type options struct {
	logic Logic
}

type logicOption Logic

func (lo logicOption) apply(opts *options) {
	opts.logic = Logic(lo)
}

// WithLogic sets the logical operator used in permission or role checks.
func WithLogic(logic Logic) Option {
	return logicOption(logic)
}

var enforcer *casbin.Enforcer

// InitCasbin ...
func InitCasbin(modelFile string, policyAdapter interface{}) {
	once := sync.Once{}
	once.Do(func() {
		enforcer = casbin.NewEnforcer(modelFile, policyAdapter)
	})
}

// RequiresPermissions tries to find the current subject by calling SubjectFn
// and determine if the subject has the required permissions according to predefined Casbin policies.
// permissions are formatted strings. For example, "file:read" represents the permission to read a file.
// opts is some optional configurations such as the logical operator (default is AND) in case multiple permissions are specified.
func RequiresPermissions(permissions []string, opts ...Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(permissions) == 0 {
			c.Next()
			return
		}

		sub := c.Get(valuer.KeyUser).String()
		if sub == "" {
			log.Error("sub is null")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Here we provide default options.
		actualOptions := options{
			logic: AND,
		}
		// Apply actual options.
		for _, opt := range opts {
			opt.apply(&actualOptions)
		}

		// Enforce Casbin policies.
		if actualOptions.logic == AND {
			// Must pass all tests.
			for _, permission := range permissions {
				obj, act := parsePermissionStrings(permission)
				if obj == "" || act == "" {
					// Can not handle any illegal permission strings.
					log.Error("illegal permission string: ", permission)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				if ok := enforcer.Enforce(sub, obj, act); !ok {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
			}
			c.Next()
		} else {
			// Need to pass at least one test.
			for _, permission := range permissions {
				obj, act := parsePermissionStrings(permission)
				if obj == "" || act == "" {
					log.Error("illegal permission string: ", permission)
					c.AbortWithStatus(http.StatusInternalServerError)
					continue
				}

				if ok := enforcer.Enforce(sub, obj, act); ok {
					c.Next()
					return
				}
			}
			c.AbortWithStatus(401)
		}
	}
}

func parsePermissionStrings(str string) (string, string) {
	if !strings.Contains(str, ":") {
		return "", ""
	}
	vals := strings.Split(str, ":")
	return vals[0], vals[1]
}

// RequiresRoles tries to find the current subject by calling SubjectFn
// and determine if the subject has the required roles according to predefined Casbin policies.
// opts is some optional configurations such as the logical operator (default is AND) in case multiple roles are specified.
func RequiresRoles(requiredRoles []string, opts ...Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(requiredRoles) == 0 {
			c.Next()
			return
		}

		sub := c.Get(valuer.KeyUser).String()
		if sub == "" {
			log.Error("sub is null")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Here we provide default options.
		actualOptions := options{
			logic: AND,
		}
		// Apply actual options.
		for _, opt := range opts {
			opt.apply(&actualOptions)
		}

		actualRoles, err := enforcer.GetRolesForUser(sub)
		if err != nil {
			log.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Enforce Casbin policies.
		sort.Strings(requiredRoles)
		sort.Strings(actualRoles)
		if actualOptions.logic == AND {
			// Must have all required roles.
			if !reflect.DeepEqual(requiredRoles, actualRoles) {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.Next()
			}
		} else {
			// Need to have at least one of required roles.
			for _, requiredRole := range requiredRoles {
				if i := sort.SearchStrings(actualRoles, requiredRole); i >= 0 &&
					i < len(actualRoles) &&
					actualRoles[i] == requiredRole {
					c.Next()
					return
				}
			}
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
