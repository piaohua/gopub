package main

import (
	"fmt"
	"gopub/app/controllers"
	_ "gopub/app/mail"
	"gopub/app/service"
	"time"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

const VERSION = "2.0.1"

func main() {
	//service.Init()
	service.InitMgo()

	beego.AppConfig.Set("version", VERSION)
	if beego.AppConfig.String("runmode") == "dev" {
		beego.SetLevel(beego.LevelDebug)
	} else {
		beego.SetLevel(beego.LevelInformational)
		beego.SetLogger("file", `{"filename":"`+beego.AppConfig.String("log_file")+`"}`)
		beego.BeeLogger.DelLogger("console")
	}

	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/profile", &controllers.MainController{}, "*:Profile")
	beego.Router("/regist", &controllers.MainController{}, "*:Regist")
	beego.Router("/servers", &controllers.MainController{}, "*:Servers")
	beego.Router("/files", &controllers.MainController{}, "*:Files")

	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.RoleController{})
	beego.AutoRouter(&controllers.MailTplController{})
	beego.AutoRouter(&controllers.MainController{})
	beego.AutoRouter(&controllers.PlayerController{})
	beego.AutoRouter(&controllers.LoggerController{})
	beego.AutoRouter(&controllers.AgencyController{})
	beego.AutoRouter(&controllers.ChartsController{})

	// 记录启动时间
	beego.AppConfig.Set("up_time", fmt.Sprintf("%d", time.Now().Unix()))

	beego.AddFuncMap("i18n", i18n.Tr)

	beego.SetStaticPath("/assets", "assets")
	beego.SetStaticPath("/contract", "views/main/contract.html")
	beego.SetStaticPath("/rules", "views/main/rules.html")
	beego.SetStaticPath("/download", "views/main/download.html")
	beego.Run()
}
