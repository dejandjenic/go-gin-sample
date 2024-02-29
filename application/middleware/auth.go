package middleware

import (
	"strings"

	"github.com/dejandjenic/go-gin-sample/application/configuration"
	"github.com/dejandjenic/go-gin-sample/authorization"

	"github.com/gin-gonic/gin"
)

type AuthClaims struct {
	authorization.Claims
	AllScope []string
}

func anyContains(s []string, i string) bool {

	for _, v := range s {
		if strings.Contains(i, v) {
			return true
		}
	}
	return false
}

func AuthHandler(cfg *configuration.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if anyContains(cfg.AuthExclusionUrl, c.Request.RequestURI) {
			return
		}

		r := c.Request.Header["Authorization"]
		if len(r) == 0 {
			c.AbortWithStatus(401)
			return
		}
		rs := strings.Split(r[0], " ")
		if len(rs) < 2 {
			c.AbortWithStatus(401)
			return
		}
		t := rs[1]

		res, err := authorization.Authorize(t, cfg.RealmConfigURL)

		if err != nil {
			c.AbortWithStatus(401)
		}

		c.Set("claims", AuthClaims{
			Claims:   res,
			AllScope: strings.Split(res.Scope, " "),
		})
	}
}
