package reqModel

import "scnu_acm_rank/biz/model"

type CreateTeamReq struct {
	Name       string `form:"name,required" json:"name"`
	Key        string `form:"key,required" json:"key"`
	NcTeamName string `form:"nc_team_name,required" json:"nc_team_name"`
}

func (req *CreateTeamReq) GetModel() *model.Team {
	model := &model.Team{}
	model.Name = req.Name
	return model
}
