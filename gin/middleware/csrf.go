package middleware

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"yumi/gin"
)

func matchHostSuffix(suffix string) func(*url.URL) bool {
	return func(uri *url.URL) bool {
		return strings.HasSuffix(strings.ToLower(uri.Host), suffix)
	}
}

func matchPattern(pattern *regexp.Regexp) func(*url.URL) bool {
	return func(uri *url.URL) bool {
		return pattern.MatchString(strings.ToLower(uri.String()))
	}
}

// CSRF returns the csrf middleware to prevent invalid cross site request.
// Only referer is checked currently.
func CSRF(allowHosts []string, allowPattern []string) gin.HandlerFunc {
	validations := []func(*url.URL) bool{}

	addHostSuffix := func(suffix string) {
		validations = append(validations, matchHostSuffix(suffix))
	}
	addPattern := func(pattern string) {
		validations = append(validations, matchPattern(regexp.MustCompile(pattern)))
	}

	for _, r := range allowHosts {
		addHostSuffix(r)
	}
	for _, p := range allowPattern {
		addPattern(p)
	}

	return func(c *gin.Context) {
		referer := c.Request.Header.Get("Referer")
		if referer == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		illegal := true
		if uri, err := url.Parse(referer); err == nil && uri.Host != "" {
			for _, validate := range validations {
				if validate(uri) {
					illegal = false
					break
				}
			}
		}
		if illegal {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
