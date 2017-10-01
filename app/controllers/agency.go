package controllers

import (
	"encoding/json"
	"fmt"
	"gopub/app/entity"
	"gopub/app/libs"
	"gopub/app/service"
	"regexp"
	"strconv"
	"utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type AgencyController struct {
	BaseController
}

// 代理商列表
func (this *AgencyController) AgencyList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "create_time", "create_time")
	m["status"] = status
	count, _ := service.AgencyService.GetAgencyTotal(m)
	list, _ := service.AgencyService.GetAgencyList(page, this.pageSize, m)

	le := len(list)
	for i := 0; i < le; i++ {
		list[i].Cash = list[i].Cash / 100       //转换为元
		list[i].Extract = list[i].Extract / 100 //转换为元
	}

	this.Data["pageTitle"] = "代理列表"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.AgencyList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 赠送/扣除钻石操作
func (this *AgencyController) AgencyGive() {
	id := this.GetString("id")
	if id == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	}
	user, err := service.UserService.GetUser(id, false)
	this.checkError(err)
	if this.isPost() {
		diamond, err := this.GetInt("diamond")
		if err != nil {
			// handle error
			this.checkError(err)
		}

		if user.Agent == "" {
			this.checkError(fmt.Errorf("代理商不存在"))
		}

		reqMsg := &entity.ReqMsg{
			Userid: user.Agent,
			Rtype:  entity.LogType9,
			Itemid: entity.ITYPE1,
			Amount: int32(diamond),
		}
		data, err1 := json.Marshal(reqMsg)
		this.checkError(err1)
		_, err2 := service.Gm("ReqMsg", string(data))
		this.checkError(err2)

		service.ActionService.UpdateDiamond(this.auth.GetUser().UserName,
			utils.String(entity.LogType9), user.Agent, utils.String(diamond))
		this.redirect(beego.URLFor("AgencyController.AgencyList"))
	}

	this.Data["user"] = user
	this.Data["pageTitle"] = "钻石操作"
	this.display()
}

// 代理编辑
func (this *AgencyController) AgencyEdit() {
	id := this.GetString("id")
	if id == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	}
	user, err := service.UserService.GetUser(id, false)
	this.checkError(err)
	if this.isPost() {
		rate, err := this.GetInt("rate")
		if err != nil {
			this.checkError(err)
		}

		if user.Agent == "" {
			this.checkError(fmt.Errorf("代理商不存在"))
		}

		if rate >= 0 && rate <= 100 {
			user.Rate = uint32(rate)
		} else {
			this.checkError(fmt.Errorf("提取率错误"))
		}

		service.AgencyService.UpdateAgencyRate(user)

		service.ActionService.UpdateAgency(this.auth.GetUser().UserName,
			user.Agent, utils.String(rate))
		this.redirect(beego.URLFor("AgencyController.AgencyList"))
	}

	this.Data["user"] = user
	this.Data["pageTitle"] = "编辑操作"
	this.display()
}

// 我的充值日志
func (this *AgencyController) MyAgencyPay() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	id := this.GetString("id")
	if page < 1 {
		page = 1
	}

	p, err := service.PlayerService.GetPlayer(id)
	agent := this.auth.GetUser().Agent
	if id == "" || err != nil || p.Agent != agent || agent == "" || p.Agent == "" {
		this.checkError(fmt.Errorf("非法操作"))
	} else {
		m := service.FindByDate(startDate, endDate, "ctime", "ctime")
		m["userid"] = id
		m["result"] = entity.TradeSuccess
		count, _ := service.LoggerService.GetPayTotal(m)
		list, _ := service.LoggerService.GetPayList(page, this.pageSize, m)

		le := len(list)
		for i := 0; i < le; i++ {
			list[i].Money = list[i].Money / 100 //转换为元
		}

		this.Data["pageTitle"] = "充值记录"
		this.Data["count"] = count
		this.Data["list"] = list
		this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MyAgencyPay", "status", status, "id", id, "start_date", startDate, "end_date", endDate), true).ToString()
		this.Data["startDate"] = startDate
		this.Data["endDate"] = endDate
	}
	this.display()
}

// 我的代理,绑定我的用户
func (this *AgencyController) MyAgencyList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "agent_time", "agent_time")
	username := this.auth.GetUser().UserName
	list, _ := service.AgencyService.GetMyAgencyList(username, page, this.pageSize, m)
	count, _ := service.AgencyService.GetMyAgencyTotal(username, m)

	this.Data["pageTitle"] = "我的代理"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MyAgencyList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

/*
// 代理商赠送/扣除钻石操作,同时扣除代理商金额
func (this *AgencyController) MyAgencyEdit() {
	userid := this.GetString("id")
	if userid == "" {
		this.checkError(fmt.Errorf("用户ID不能为空"))
	}
	if this.isPost() {
		diamond, err := this.GetInt("diamond")
		if err != nil {
			// handle error
			this.checkError(err)
		}

		agent := this.auth.GetUser().Agent
		if agent == "" || diamond <= 0 {
			this.checkError(fmt.Errorf("无法赠送"))
		} else {
			reqMsg := &entity.ReqGiveDiamondMsg{
				Userid: userid,
				Agent:  agent,
				Rtype:  entity.LogType9,
				Itemid: entity.ITYPE1,
				Amount: int32(diamond),
			}
			data, err1 := json.Marshal(reqMsg)
			this.checkError(err1)
			_, err2 := service.Gm("ReqGiveDiamondMsg", string(data))
			this.checkError(err2)

			service.ActionService.UpdateDiamond(this.auth.GetUser().UserName,
				utils.String(entity.LogType9), userid, utils.String(diamond))
		}
		this.redirect(beego.URLFor("AgencyController.AgencyList"))
	}

	p, err := service.PlayerService.GetPlayer(userid)
	this.checkError(err)

	this.Data["id"] = userid
	this.Data["nickname"] = p.Nickname
	this.Data["pageTitle"] = "赠送钻石"
	this.display()
}
*/

// 提现记录
func (this *AgencyController) CashList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	m["status"] = status
	list, _ := service.AgencyService.GetCashList(page, this.pageSize, m)
	count, _ := service.AgencyService.GetCashListTotal(m)

	le := len(list)
	for i := 0; i < le; i++ {
		list[i].Cash = list[i].Cash / 100 //转换为元
	}

	this.Data["pageTitle"] = "提现记录"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.CashList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 我的提现
func (this *AgencyController) MyCashList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	m["status"] = status
	username := this.auth.GetUser().UserName
	list, _ := service.AgencyService.GetMyCashList(username, page, this.pageSize, m)
	count, _ := service.AgencyService.GetMyCashListTotal(username, m)

	le := len(list)
	for i := 0; i < le; i++ {
		list[i].Cash = list[i].Cash / 100 //转换为元
	}

	this.Data["pageTitle"] = "我的提现"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MyCashList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 申请提现,让玩家自己输入银行卡号，开户行，姓名
func (this *AgencyController) MyCashAdd() {
	payway := map[int]string{
		entity.PAYWAY1: "微信",
		entity.PAYWAY2: "支付宝",
		entity.PAYWAY3: "银行账号",
	}
	if this.isPost() {
		name := this.GetString("name")
		bankCard, _ := this.GetInt("bankCard")
		bankAddr := this.GetString("bankAddr")
		money, _ := this.GetFloat("money")
		valid := validation.Validation{}
		valid.Required(name, "name").Message("姓名不能为空")
		valid.Required(bankCard, "bankCard").Message("收款方式不能为空")
		valid.Required(bankAddr, "bankAddr").Message("收款账号不能为空")
		valid.Required(money, "money").Message("提现金额错误")
		if payway[bankCard] == "" {
			this.checkError(fmt.Errorf("收款方式错误"))
		} else if !valid.HasErrors() {
			username := this.auth.GetUser().UserName
			money = money * 100 //元转换为分
			err := service.AgencyService.ApplyCashAdd(username, name, bankAddr, bankCard, money)
			if err != nil {
				this.showMsg(err.Error(), MSG_ERR)
			} else {
				service.ActionService.AddApplyCash(username, utils.String(money))
				this.redirect(beego.URLFor("AgencyController.MyCashList"))
			}
		} else {
			for _, err := range valid.Errors {
				this.showMsg(err.Message, MSG_ERR)
				break
			}
			this.Data["name"] = name
			this.Data["bankCard"] = bankCard
			this.Data["bankAddr"] = bankAddr
			this.Data["money"] = money
		}
	}

	this.Data["payway"] = payway
	this.Data["pageTitle"] = "申请提现"
	this.display()
}

// 提现处理
func (this *AgencyController) AgencyExtract() {
	orderid := this.GetString("id")
	if orderid == "" {
		this.checkError(fmt.Errorf("订单不存在"))
	}

	username := this.auth.GetUser().UserName
	err := service.AgencyService.ExtractCash(username, orderid)
	if err != nil {
		this.checkError(err)
	} else {
		service.ActionService.ExtractApplyCash(username, orderid)
	}

	this.redirect(beego.URLFor("AgencyController.CashList"))
}

// 我的充值日志
func (this *AgencyController) MyPayList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	agent := this.auth.GetUser().Agent
	if agent == "" {
		this.checkError(fmt.Errorf("非法操作"))
	} else {
		m := service.FindByDate(startDate, endDate, "ctime", "ctime")
		m["userid"] = agent
		m["result"] = entity.TradeSuccess
		count, _ := service.LoggerService.GetPayTotal(m)
		list, _ := service.LoggerService.GetPayList(page, this.pageSize, m)

		le := len(list)
		for i := 0; i < le; i++ {
			list[i].Money = list[i].Money / 100 //转换为元
		}

		this.Data["pageTitle"] = "我的充值"
		this.Data["count"] = count
		this.Data["list"] = list
		this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.MyPayList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
		this.Data["startDate"] = startDate
		this.Data["endDate"] = endDate
	}
	this.display()
}

// 代理商邀请的用户
func (this *AgencyController) InviterList() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}
	inviter := this.auth.GetUser().UserName
	if inviter == "" {
		this.checkError(fmt.Errorf("非法操作"))
	} else {
		m := service.FindByDate(startDate, endDate, "create_time", "create_time")
		m["inviter"] = inviter
		list, _ := service.AgencyService.GetInviterList(page, this.pageSize, m)
		count, _ := service.AgencyService.GetInviterTotal(m)
		this.Data["count"] = count
		this.Data["list"] = list
		this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.InviterList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	}

	this.Data["pageTitle"] = "我的邀请"
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 代理商自己去邀请用户
func (this *AgencyController) InviterUser() {
	inviter := this.auth.GetUser().UserName
	if inviter == "" {
		this.checkError(fmt.Errorf("非法操作"))
	} else if this.isPost() {
		valid := validation.Validation{}

		username := this.GetString("username")
		phone := this.GetString("phone")
		agent := this.GetString("agent")
		weixin := this.GetString("weixin")
		sex, _ := this.GetInt("sex")
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")

		valid.Required(username, "username").Message("请输入用户名")
		valid.Required(agent, "agent").Message("请输入游戏ID")
		valid.Required(password1, "password1").Message("请输入密码")
		valid.Required(password2, "password2").Message("请输入确认密码")
		valid.MinSize(password1, 6, "password1").Message("密码长度不能小于6个字符")
		valid.Match(password1, regexp.MustCompile(`^`+regexp.QuoteMeta(password2)+`$`), "password2").Message("两次输入的密码不一致")
		if valid.HasErrors() {
			for _, err := range valid.Errors {
				this.showMsg(err.Message, MSG_ERR)
			}
		}

		user, err := service.UserService.InviterAddUser(username, "", inviter, agent, phone, weixin, password1, sex)
		if err == nil {
			service.ActionService.AddInviterUser(this.auth.GetUser().UserName, username)
		}
		this.checkError(err)

		// 更新角色
		roleIds := make([]string, 0)
		for _, v := range this.GetStrings("role_ids") {
			roleIds = append(roleIds, v)
		}
		service.UserService.UpdateUserRoles(user.Id, roleIds)

		this.redirect(beego.URLFor("AgencyController.InviterList"))
	}

	roleList, _ := service.RoleService.GetAllRoles()
	//非超级管理员,过滤掉超级管理员
	//if this.auth.GetUser().Id != "1" {
	for k, v := range roleList {
		if v.Id == "1" {
			roleList = append(roleList[:k], roleList[k+1:]...)
			break
		}
	}
	//}
	this.Data["pageTitle"] = "邀请用户"
	this.Data["roleList"] = roleList
	this.display()
}
