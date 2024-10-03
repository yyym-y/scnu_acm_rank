package user

import (
	"context"
	"errors"
	"net/http"
	"scnu_acm_rank/biz/middle"
	"scnu_acm_rank/biz/model"

	"github.com/cloudwego/hertz/pkg/app"
)

func TotalTeamDetail(ctx context.Context, c *app.RequestContext) {
	userInter, flag := c.Get("user")
	if !flag {
		c.JSON(http.StatusOK, middle.FailResp(errors.New("没有登录")))
		return
	}
	user, ok := userInter.(*model.User)
	if !ok {
		c.JSON(http.StatusOK, middle.FailResp(errors.New("没有登录")))
		return
	}
	// 通过 user 中的学生的学号来找到属于这个学生的队伍信息
	res := make([]map[string]interface{}, 0)
	sql := `SELECT team.name, team.key, team.leader, team.nc_team_name FROM team WHERE id = (
    			SELECT group_id FROM user WHERE stu_id = ?
			);`
	model.DB.Raw(sql, user.StuId).Find(&res)
	c.JSON(http.StatusOK, middle.SuccessResp("测试", res[0]))
}
