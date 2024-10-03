package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"scnu_acm_rank/biz/middle"
	"scnu_acm_rank/biz/model"

	"github.com/cloudwego/hertz/pkg/app"
)

func EditTeam(ctx context.Context, c *app.RequestContext) {
	userInter, flag := c.Get("user")
	if !flag {
		c.JSON(http.StatusOK, middle.FailResp(errors.New("没有登录")))
		return
	}
	_, ok := userInter.(*model.User)
	if !ok {
		c.JSON(http.StatusOK, middle.FailResp(errors.New("没有登录")))
		return
	}
	req := model.Team{}
	err := c.BindForm(&req)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, middle.FailResp(err))
		return
	}
	if len(req.Name) > 255 || len(req.Key) > 255 || len(req.NcTeamName) > 255 {
		fmt.Println("队伍名/密钥/牛客队伍名长度超过限制")
		c.JSON(http.StatusOK, middle.FailRespWithMsg("队伍名/密钥/牛客队伍名长度超过限制"))
		return
	}
	leader := model.User{}
	model.DB.Model(&leader).Where("user.stu_id = ?", req.Leader).First(&leader)
	if leader.Id == 0 {
		fmt.Println("队长学号信息不存在")
		c.JSON(http.StatusOK, middle.FailRespWithMsg("队长学号信息不存在"))
		return
	}
	mutex.Lock()
	model.DB.Model(&req).Where("id = ?", req.Id).Update("name", req.Name).Update("key", req.Key).Update("leader", req.Leader).Update("nc_team_name", req.NcTeamName)
	mutex.Unlock()
	c.JSON(http.StatusOK, middle.SuccessResp("修改成功", nil))
}
