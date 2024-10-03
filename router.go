// Code generated by hertz generator.

package main

import (
	"log"
	handler "scnu_acm_rank/biz/handler"
	root2 "scnu_acm_rank/biz/handler/root"
	super2 "scnu_acm_rank/biz/handler/super"
	user2 "scnu_acm_rank/biz/handler/user"
	"scnu_acm_rank/biz/middle"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// customizeRegister registers customize routers.
func customizedRegister(r *server.Hertz) {

	auth, err := middle.GetJWT()
	if err != nil {
		panic(err)
	}
	errInit := auth.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.GET("/ping", handler.Ping)

	competition := r.Group("/competition")
	competition.GET("/person", handler.CompetitionPerson)
	competition.GET("/group", handler.CompetitionGroup)
	competition.GET("/detail", handler.CompetitionDetail) // 1

	// your code ...
	r.POST("/login", auth.LoginHandler)    //1
	r.POST("/register", handler.Register)  //1
	r.POST("sendEmail", handler.SendEmail) // 1
	r.GET("/personCompetitions", handler.UserCompetitions)
	r.GET("/groupCompetitions", handler.GroupCompetitions)
	r.GET("/TeamDetail", handler.TeamDetail)
	user := r.Group("/user")
	user.Use(auth.MiddlewareFunc())
	user.POST("/edit", user2.EditUser)                   // 1
	user.GET("/detail", handler.UserDetail)              // 1
	user.POST("/createTeam", user2.CreateTeam)           // 1
	user.POST("/joinTeam", user2.JoinTeam)               // 1
	user.POST("/editTeam", user2.EditTeam)               // 1
	user.POST("/totalTeamDetail", user2.TotalTeamDetail) // 1

	root := r.Group("/root")
	root.POST("/createCompetition", root2.CreateCompetition) // 1
	root.POST("/updateConfig", root2.UpdateConfig)           // 1

	super := r.Group("/super")
	super.POST("/addRoot", super2.AddRoot)
}
