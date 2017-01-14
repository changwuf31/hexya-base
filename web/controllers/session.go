// Copyright 2016 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package controllers

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/npiganeau/yep/yep/models"
	"github.com/npiganeau/yep/yep/models/security"
	"github.com/npiganeau/yep/yep/models/types"
	"github.com/npiganeau/yep/yep/server"
)

func SessionInfo(sess sessions.Session) gin.H {
	var userContext *types.Context
	if sess.Get("uid") != nil {
		models.ExecuteInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			user := env.Pool("ResUsers").Filter("ID", "=", sess.Get("uid"))
			userContext = user.Call("ContextGet").(*types.Context)
		})
	}
	return gin.H{
		"session_id":   sess.Get("ID"),
		"uid":          sess.Get("uid"),
		"user_context": userContext.ToMap(),
		"db":           "default",
		"username":     sess.Get("login"),
		"company_id":   1,
	}
}

func GetSessionInfo(c *gin.Context) {
	sess := sessions.Default(c)
	server.RPC(c, http.StatusOK, SessionInfo(sess))
}

func Modules(c *gin.Context) {
	mods := make([]string, len(server.Modules))
	for i, m := range server.Modules {
		mods[i] = m.Name
	}
	server.RPC(c, http.StatusOK, mods)
}