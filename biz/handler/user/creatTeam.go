package user

import (
	"context"
	"errors"
	"net/http"
	"runtime"
	"scnu_acm_rank/biz/middle"
	"scnu_acm_rank/biz/model"
	"scnu_acm_rank/biz/reqModel"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

var mutex sync.Mutex

func init() {
	mutex = sync.Mutex{}
	runtime.KeepAlive(mutex)
}

func CreateTeam(ctx context.Context, c *app.RequestContext) {
	req := reqModel.CreateTeamReq{}
	err := c.BindForm(&req)
	if err != nil {
		c.JSON(http.StatusOK, middle.FailResp(err))
		return
	}
	v, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusOK, middle.FailResp(errors.New("没有登录")))
		return
	}
	user := v.(*model.User)
	team := req.GetModel()
	team.Leader = user.StuId
	team.Key = req.Key
	team.Name = req.Name
	team.NcTeamName = req.NcTeamName
	team.Status = 1
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		tx.Save(&team).Find(&team)
		tx.Model(&user).Where("stu_id = ?", user.StuId).Update("group_id", team.Id).Find(&user)
		return tx.Error
	})
	if err != nil {
		c.JSON(http.StatusOK, middle.FailResp(err))
		return
	}
	go middle.SendKeyEmail([]string{user.Email}, req.Key)
	c.JSON(http.StatusOK, middle.SuccessResp("创建成功", map[string]interface{}{
		"group_id": team.Id,
	}))
}
