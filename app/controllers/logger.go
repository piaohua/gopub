package controllers

import (
	"gopub/app/entity"
	"gopub/app/libs"
	"gopub/app/service"
	"strconv"

	"github.com/astaxie/beego"
)

type LoggerController struct {
	BaseController
}

// 注册日志
func (this *LoggerController) RegistList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	count, _ := service.LoggerService.GetRegistTotal(m)
	list, _ := service.LoggerService.GetRegistList(page, this.pageSize, m)

	this.Data["pageTitle"] = "注册日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.RegistList", "status", status, "userid", userid, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 登录日志
func (this *LoggerController) LoginList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "login_time", "login_time")
	if userid != "" {
		m["userid"] = userid
	}
	count, _ := service.LoggerService.GetLoginTotal(m)
	list, _ := service.LoggerService.GetLoginList(page, this.pageSize, m)

	this.Data["pageTitle"] = "登录日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.LoginList", "status", status, "userid", userid, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 充值日志
func (this *LoggerController) PayList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	typeId, _ := this.GetInt("type_id")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	m["result"] = typeId

	count, _ := service.LoggerService.GetPayTotal(m)
	list, _ := service.LoggerService.GetPayList(page, this.pageSize, m)

	le := len(list)
	for i := 0; i < le; i++ {
		list[i].Money = list[i].Money / 100 //转换为元
	}

	typeList := entity.TradeResult

	this.Data["pageTitle"] = "充值日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["typeList"] = typeList
	this.Data["typeId"] = typeId
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.PayList", "status", status, "userid", userid, "type_id", typeId, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 钻石日志
func (this *LoggerController) DiamondList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	typeId, _ := this.GetInt("type_id")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	if typeId != 0 {
		m["type"] = typeId
	}
	count, _ := service.LoggerService.GetDiamondTotal(m)
	list, _ := service.LoggerService.GetDiamondList(page, this.pageSize, m)

	typeList := entity.LogType

	this.Data["pageTitle"] = "钻石日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["typeList"] = typeList
	this.Data["typeId"] = typeId
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.DiamondList", "status", status, "userid", userid, "type_id", typeId, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 金币日志
func (this *LoggerController) CoinList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	count, _ := service.LoggerService.GetCoinTotal(m)
	list, _ := service.LoggerService.GetCoinList(page, this.pageSize, m)

	this.Data["pageTitle"] = "金币日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.CoinList", "status", status, "userid", userid, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 绑定日志
func (this *LoggerController) BuildList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	agent := this.GetString("agent")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	if userid != "" {
		m["userid"] = userid
	}
	if agent != "" {
		m["agent"] = agent
	}
	count, _ := service.LoggerService.GetBuildTotal(m)
	list, _ := service.LoggerService.GetBuildList(page, this.pageSize, m)

	this.Data["pageTitle"] = "绑定日志"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("LoggerController.BuildList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}
