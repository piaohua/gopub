package controllers

import (
	"encoding/json"
	"fmt"
	"gopub/app/entity"
	"gopub/app/libs"
	"gopub/app/service"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
	"utils"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type MainController struct {
	BaseController
}

// 首页
func (this *MainController) Index() {
	this.Data["pageTitle"] = "我的概况"
	//活跃用户
	popProjects := make([]map[string]interface{}, 0, 4)
	this.Data["popProjects"] = popProjects
	//基础信息
	username := this.auth.GetUser().UserName
	cashNum, _ := service.AgencyService.GetMyCashTotal(username)
	this.Data["cashNum"] = cashNum / 100 //转换为元
	this.Data["buildNum"], _ = service.AgencyService.GetMyAgencyTotal(username, bson.M{})
	extractNum, _ := service.AgencyService.GetMyExtractTotal(username)
	this.Data["extractNum"] = extractNum / 100 //转换为元
	this.Data["payNum"], _ = service.LoggerService.GetAgencyPayTotal(username)
	//最新动态
	feeds, _ := service.ActionService.GetList(username, 1, 7)
	this.Data["feeds"] = feeds
	//系统信息
	//this.Data["hostname"], _ = os.Hostname()
	this.Data["gover"] = runtime.Version()
	this.Data["os"] = runtime.GOOS
	this.Data["goroutineNum"] = runtime.NumGoroutine()
	this.Data["cpuNum"] = runtime.NumCPU()
	this.Data["arch"] = runtime.GOARCH
	up, day, hour, min, sec := this.getUptime()
	this.Data["uptime"] = fmt.Sprintf("%s，已运行 %d天 %d小时 %d分钟 %d秒", beego.Date(up, "Y-m-d H:i:s"), day, hour, min, sec)
	this.display()
}

func (this *MainController) getUptime() (up time.Time, day, hour, min, sec int) {
	ts, _ := beego.AppConfig.Int64("up_time")
	up = time.Unix(ts, 0)
	uptime := int(time.Now().Sub(up) / time.Second)
	if uptime >= 86400 {
		day = uptime / 86400
		uptime %= 86400
	}
	if uptime >= 3600 {
		hour = uptime / 3600
		uptime %= 3600
	}
	if uptime >= 60 {
		min = uptime / 60
		uptime %= 60
	}
	sec = uptime
	return
}

// 绑定统计
func (this *MainController) GetPubStat() {
	username := this.auth.GetUser().UserName
	rangeType := this.GetString("range")
	result := service.LoggerService.GetPubStat(username, rangeType)

	ticks := make([]interface{}, 0)
	chart := make([]interface{}, 0)
	json := make(map[string]interface{}, 0)
	switch rangeType {
	case "this_month":
		year, month, _ := time.Now().Date()
		maxDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).AddDate(0, 0, -1).Day()

		for i := 1; i <= maxDay; i++ {
			var row [3]interface{}
			row[0] = i
			row[1] = fmt.Sprintf("%02d", i)
			row[2] = fmt.Sprintf("%d-%02d-%02d", year, month, i)
			ticks = append(ticks, row)
			if v, ok := result[i]; ok {
				chart = append(chart, []int{i, v})
			} else {
				chart = append(chart, []int{i, 0})
			}
		}
	case "last_month":
		year, month, _ := time.Now().AddDate(0, -1, 0).Date()
		maxDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).AddDate(0, 0, -1).Day()

		for i := 1; i <= maxDay; i++ {
			var row [3]interface{}
			row[0] = i
			row[1] = fmt.Sprintf("%02d", i)
			row[2] = fmt.Sprintf("%d-%02d-%02d", year, month, i)
			ticks = append(ticks, row)
			if v, ok := result[i]; ok {
				chart = append(chart, []int{i, v})
			} else {
				chart = append(chart, []int{i, 0})
			}
		}
	case "this_year":
		year := time.Now().Year()
		for i := 1; i <= 12; i++ {
			var row [3]interface{}
			row[0] = i
			row[1] = fmt.Sprintf("%d月", i)
			row[2] = fmt.Sprintf("%d年%d月", year, i)
			ticks = append(ticks, row)
			if v, ok := result[i]; ok {
				chart = append(chart, []int{i, v})
			} else {
				chart = append(chart, []int{i, 0})
			}
		}
	case "last_year":
		year := time.Now().Year() - 1
		for i := 1; i <= 12; i++ {
			var row [3]interface{}
			row[0] = i
			row[1] = fmt.Sprintf("%d月", i)
			row[2] = fmt.Sprintf("%d年%d月", year, i)
			ticks = append(ticks, row)
			if v, ok := result[i]; ok {
				chart = append(chart, []int{i, v})
			} else {
				chart = append(chart, []int{i, 0})
			}
		}
	}

	json["ticks"] = ticks
	json["chart"] = chart
	this.Data["json"] = json
	this.ServeJSON()
}

// 个人信息
func (this *MainController) Profile() {
	beego.ReadFromRequest(&this.Controller)
	user := this.auth.GetUser()

	if this.isPost() {
		flash := beego.NewFlash()
		email := this.GetString("email")
		sex, _ := this.GetInt("sex")
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")

		user.Email = email
		user.Sex = sex
		fileds := bson.M{"email": user.Email, "sex": user.Sex}
		service.UserService.UpdateUser(user, fileds)
		if password1 != "" {
			if len(password1) < 6 {
				flash.Error("密码长度必须大于6位")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else if password2 != password1 {
				flash.Error("两次输入的密码不一致")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else {
				service.UserService.ModifyPassword(this.userId, password1)
			}
		}
		service.ActionService.UpdateProfile(this.auth.GetUser().UserName, this.userId)
		flash.Success("修改成功！")
		flash.Store(&this.Controller)
		this.redirect(beego.URLFor(".Profile"))
	}

	this.Data["pageTitle"] = "个人信息"
	this.Data["user"] = user
	this.display()
}

// 登录
func (this *MainController) Login() {
	if this.userId != "" {
		this.redirect("/")
	}
	beego.ReadFromRequest(&this.Controller)
	if this.isPost() {
		flash := beego.NewFlash()
		username := this.GetString("username")
		password := this.GetString("password")
		remember := this.GetString("remember")
		if username != "" && password != "" {
			token, err := this.auth.Login(username, password)
			if err != nil {
				flash.Error(err.Error())
				flash.Store(&this.Controller)
				this.redirect("/login")
			} else {
				if remember == "yes" {
					this.Ctx.SetCookie("auth", token, 7*86400)
				} else {
					this.Ctx.SetCookie("auth", token)
				}
				service.ActionService.Login(username, this.auth.GetUserId(), this.getClientIp())
				this.redirect(beego.URLFor(".Index"))
			}
		}
	}
	this.TplName = "main/login.html"
}

// 退出登录
func (this *MainController) Logout() {
	service.ActionService.Logout(this.auth.GetUser().UserName, this.auth.GetUserId(), this.getClientIp())
	this.auth.Logout()
	this.Ctx.SetCookie("auth", "")
	this.redirect(beego.URLFor(".Login"))
}

// 注册
func (this *MainController) Regist() {
	beego.ReadFromRequest(&this.Controller)
	if this.isPost() {
		flash := beego.NewFlash()
		username := this.GetString("username")
		password := this.GetString("password")
		password1 := this.GetString("password1")
		agent := this.GetString("agent")
		weixin := this.GetString("weixin")
		phone := this.GetString("phone")
		qq := this.GetString("qq")
		address := this.GetString("address")
		valid := validation.Validation{}
		valid.Required(username, "username").Message("账号不能为空")
		valid.Required(password, "password").Message("密码不能为空")
		valid.Required(password1, "password1").Message("请输入确认密码")
		valid.Required(agent, "agent").Message("游戏id不能为空")
		valid.Required(weixin, "weixin").Message("微信号不能为空")
		//valid.Required(qq, "qq").Message("QQ不能为空")
		//valid.Required(address, "address").Message("地址不能为空")
		valid.Mobile(phone, "phone").Message("手机号码不正确")
		valid.MinSize(password, 6, "password").Message("密码长度不能小于6位")
		//valid.Numeric(password, "password").Message("密码不能是纯数字")
		valid.Numeric(agent, "agent").Message("代理ID不正确")
		valid.Match(password, regexp.MustCompile(`^`+regexp.QuoteMeta(password1)+`$`), "password1").Message("两次输入的密码不一致")
		if !valid.HasErrors() {
			err := this.auth.Regist(username, password, agent, phone, weixin, qq, address)
			if err != nil {
				flash.Error(err.Error())
				flash.Store(&this.Controller)
				this.redirect("/regist")
			} else {
				service.ActionService.Regist(username, agent, this.getClientIp())
				this.redirect(beego.URLFor(".Login"))
			}
		} else {
			for _, err := range valid.Errors {
				flash.Error(err.Message)
				flash.Store(&this.Controller)
				break
			}
			this.Data["username"] = username
			this.Data["agent"] = agent
			this.Data["weixin"] = weixin
			this.Data["phone"] = phone
			this.Data["qq"] = qq
			this.Data["address"] = address
		}
	}
	this.TplName = "main/regist.html"
}

// 服务器列表
func (this *MainController) Servers() {
	userid := this.GetString("userid")
	platform := this.GetString("platform")
	//fmt.Println("userid ", userid, "platform ", platform)
	//this.getClientIp()
	//servers := make([]entity.Server, 0)
	server := new(entity.Server)
	content, err := ioutil.ReadFile("conf/servers.json")
	if err != nil {
		beego.Error(err.Error())
	}
	//err = json.Unmarshal(content, &servers)
	srv := new(entity.Srv)
	err = json.Unmarshal(content, srv)
	if err != nil {
		beego.Error(err.Error())
	}
	server.Name = srv.Name
	server.Host = srv.Host
	server.Port = srv.Port
	server.Status = srv.Status
	server.Info = srv.Info
	server.UpdateUrl = srv.UpdateUrl
	server.UpdateDes = srv.UpdateDes
	server.Version = srv.Version
	//
	switch platform {
	case "ios":
		server.Version = srv.VersionIos
		server.UpdateUrl = srv.UpdateUrlIos
	case "android":
	default:
		str := this.getUserAgent()
		if (len(utils.Split(str, "CFNetwork")) < 2 &&
			len(utils.Split(str, "Darwin")) < 2 &&
			len(utils.Split(str, "iPhone")) < 2 &&
			len(utils.Split(str, "iPad")) < 2) ||
			len(utils.Split(str, "Android")) > 1 ||
			str == "" { //android
		} else { //ios
			server.Version = srv.VersionIos
			server.UpdateUrl = srv.UpdateUrlIos
		}
	}
	if userid != "" && srv.Userids != "" {
		str := utils.Split(srv.Userids, ",")
		for _, v := range str {
			if v == userid {
				server.Version = srv.VersionBeta //内测版本
				break
			}
		}
	}
	//this.jsonResult(servers)
	this.jsonResult(server)
}

// 文件列表
func (this *MainController) Files() {
	lm := make([]entity.ListFiles, 0)
	//遍历目录，读出文件名、大小
	filepath.Walk(service.GetFilePath(), func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		//fmt.Println("path -> ", path, fi.Name())
		str := utils.Split(path, "assets/files/")
		if len(str) != 2 {
			return nil
		}
		var m entity.ListFiles
		m.Name = str[1]
		m.Size = utils.String(fi.Size() / 1024)
		m.Time = fi.ModTime()
		lm = append(lm, m)
		return nil
	})
	this.Data["list"] = lm
	this.Data["pageTitle"] = "文件列表"
	this.TplName = "main/files.html"
}

// 删除文件
func (this *MainController) DeleteFile() {
	fileName := this.GetString("fileName")
	path := service.GetFilePath() + "/" + fileName //文件目录
	err := os.Remove(path)                         //删除文件
	if err != nil {
		beego.Error(err.Error())
	}
	this.redirect(beego.URLFor("MainController.Files"))
}

// 上传文件
func (this *MainController) UploadFile() {
	if this.isPost() {
		file, h, err := this.GetFile("fileName") //获取上传的文件
		//fmt.Println("err : ", err)
		if err != nil {
			beego.Error(err.Error())
		} else if file != nil {
			file.Close() //关闭上传的文件
			//file.Size()
			dir := service.GetFilePath()
			path := dir + "/" + h.Filename //文件目录
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				beego.Error(err.Error())
			}
			if f != nil {
				io.Copy(f, file)
				f.Close()
			}
			str := utils.Split(h.Filename, "zip")
			if len(str) > 1 && str[1] == "" {
				beego.Trace("unzip dir ", dir, " name ", h.Filename)
				stdout, stderr, err := libs.ExecCmdDir(dir, "unzip", "-ou", h.Filename)
				beego.Trace("unzip ", stdout, stderr, err)
				//cmd := "tar -xf " + h.Filename
				//stdout, stderr, err = libs.ExecCmdDir(dir, "/bin/bash", "-c", cmd)
				//fmt.Println(stdout, stderr, err)
			}
		}
	}
	this.redirect(beego.URLFor("MainController.Files"))
}

// 下载文件
func (this *MainController) DownFile() {
	fileName := this.GetString("fileName")
	path := service.GetFilePath() + "/" + fileName //文件目录
	if fileName == "version.manifest" {
		platform := this.GetString("platform")
		switch platform {
		case "ios":
			fileName2 := "version.manifest.ios"
			path = service.GetFilePath() + "/" + fileName2 //文件目录
		case "android":
			fileName2 := "version.manifest.android"
			path = service.GetFilePath() + "/" + fileName2 //文件目录
		default:
			str := this.getUserAgent()
			//hnmj_xt-mobile/10 CFNetwork/811.4.18 Darwin/16.5.0
			if (len(utils.Split(str, "CFNetwork")) < 2 &&
				len(utils.Split(str, "Darwin")) < 2 &&
				len(utils.Split(str, "iPhone")) < 2 &&
				len(utils.Split(str, "iPad")) < 2) ||
				len(utils.Split(str, "Android")) > 1 ||
				str == "" { //android
				fileName2 := "version.manifest.android"
				path = service.GetFilePath() + "/" + fileName2 //文件目录
			}
		}
		this.Ctx.Output.Download(path, fileName)
	} else {
		this.Ctx.Output.Download(path)
	}
}
