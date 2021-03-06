// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	api "github.com/gigforks/go-gogs-client"

	"github.com/gigforks/gogs/models"
	"github.com/gigforks/gogs/modules/context"
	"github.com/gigforks/gogs/routers/api/v1/convert"
	"github.com/gigforks/gogs/routers/api/v1/user"
)

func CreateTeam(ctx *context.APIContext, form api.CreateTeamOption) {
	team := &models.Team{
		OrgID:       ctx.Org.Organization.Id,
		Name:        form.Name,
		Description: form.Description,
		Authorize:   models.ParseAccessMode(form.Permission),
	}
	if err := models.NewTeam(team); err != nil {
		if models.IsErrTeamAlreadyExist(err) {
			ctx.Error(422, "", err)
		} else {
			ctx.Error(500, "NewTeam", err)
		}
		return
	}

	ctx.JSON(201, convert.ToTeam(team))
}

func AddTeamMember(ctx *context.APIContext) {
	u := user.GetUserByParams(ctx)
	if ctx.Written() {
		return
	}
	if err := ctx.Org.Team.AddMember(u.Id); err != nil {
		ctx.Error(500, "AddMember", err)
		return
	}

	ctx.Status(204)
}

func RemoveTeamMember(ctx *context.APIContext) {
	u := user.GetUserByParams(ctx)
	if ctx.Written() {
		return
	}

	if err := ctx.Org.Team.RemoveMember(u.Id); err != nil {
		ctx.Error(500, "RemoveMember", err)
		return
	}

	ctx.Status(204)
}
