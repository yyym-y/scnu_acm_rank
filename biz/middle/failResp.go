package middle

import "github.com/cloudwego/hertz/pkg/common/utils"

func FailResp(err error) utils.H {
	return utils.H{
		"status": 1,
		"msg":    "",
		"data":   err,
	}
}

func FailRespWithMsg(msg string) utils.H {
	return utils.H{
		"status": 1,
		"msg":    msg,
		"data":   nil,
	}
}

func SuccessResp(msg string, data interface{}) utils.H {
	return utils.H{
		"status": 0,
		"msg":    msg,
		"data":   data,
	}
}
