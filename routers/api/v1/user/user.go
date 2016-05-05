// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/Unknwon/com"

	api "github.com/gigforks/go-gogs-client"

	"github.com/gigforks/gogs/models"
	"github.com/gigforks/gogs/modules/middleware"
)

// https://github.com/gigforks/go-gogs-client/wiki/Users#search-users
func Search(ctx *middleware.Context) {
	opt := models.SearchOption{
		Keyword: ctx.Query("q"),
		Limit:   com.StrTo(ctx.Query("limit")).MustInt(),
	}
	if opt.Limit == 0 {
		opt.Limit = 10
	}

	us, err := models.SearchUserByName(opt)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	results := make([]*api.User, len(us))
	for i := range us {
		results[i] = &api.User{
			ID:        us[i].Id,
			UserName:  us[i].Name,
			AvatarUrl: us[i].AvatarLink(),
			FullName:  us[i].FullName,
		}
		if ctx.IsSigned {
			results[i].Email = us[i].Email
		}
	}

	ctx.JSON(200, map[string]interface{}{
		"ok":   true,
		"data": results,
	})
}

// https://github.com/gigforks/go-gogs-client/wiki/Users#get-a-single-user
func GetInfo(ctx *middleware.Context) {
	u, err := models.GetUserByName(ctx.Params(":username"))
	if err != nil {
		if models.IsErrUserNotExist(err) {
			ctx.Error(404)
		} else {
			ctx.APIError(500, "GetUserByName", err)
		}
		return
	}

	// Hide user e-mail when API caller isn't signed in.
	if !ctx.IsSigned {
		u.Email = ""
	}
	ctx.JSON(200, &api.User{u.Id, u.Name, u.FullName, u.Email, u.AvatarLink()})
}
