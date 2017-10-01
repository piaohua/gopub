package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopub/app/entity"
	"gopub/app/libs"
	"gopub/app/service"
	"strconv"
	"utils"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type PlayerController struct {
	BaseController
}

// 玩家列表
func (this *PlayerController) List() {
	status, _ := this.GetInt("status")
	page, _ := strconv.Atoi(this.GetString("page"))
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	userid := this.GetString("userid")
	typeId, _ := this.GetInt("type_id")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "create_time", "create_time")
	if userid != "" {
		m["_id"] = userid
	}
	count, _ := service.PlayerService.GetTotal(typeId, m)
	list, _ := service.PlayerService.GetList(typeId, page, this.pageSize, m)

	ids := make([]string, 0)
	for _, v := range list {
		ids = append(ids, v.Userid)
	}

	msg := &entity.ReqOnlineStatusMsg{Userid: ids}
	data, err := json.Marshal(msg)
	fmt.Printf("data %s, err %v\n", string(data), err)
	this.checkError(err)
	flash := beego.NewFlash()
	if len(data) > 0 {
		result, err := service.Gm("ReqOnlineStatusMsg", string(data))
		fmt.Printf("result %s, err %v\n", result, err)
		if err != nil {
			flash.Error(fmt.Sprintf("%v", err))
			flash.Store(&this.Controller)
		} else {
			resp := new(entity.RespOnlineStatusMsg)
			err3 := json.Unmarshal([]byte(result), resp)
			if err3 != nil {
				flash.Error(fmt.Sprintf("%v", err3))
				flash.Store(&this.Controller)
			}
			//this.checkError(err3)
			for k, v := range list {
				list[k].State = resp.Userid[v.Userid]
			}
		}
	}

	typeList := map[int]string{
		0: "微信用户",
		1: "全部玩家",
		2: "其它用户",
	}

	this.Data["pageTitle"] = "玩家列表"
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["typeList"] = typeList
	this.Data["typeId"] = typeId
	this.Data["userid"] = userid
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("PlayerController.List", "status", status, "type_id", typeId, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 钻石操作
func (this *PlayerController) Edit() {
	userid := this.GetString("userid")
	if this.isPost() {
		diamond, err := this.GetInt("diamond")
		if userid == "" {
			this.checkError(fmt.Errorf("用户ID不能为空"))
		}
		if err != nil {
			// handle error
			this.checkError(err)
		}

		reqMsg := &entity.ReqMsg{
			Userid: userid,
			Rtype:  entity.LogType9,
			Itemid: entity.ITYPE1,
			Amount: int32(diamond),
		}
		data, err1 := json.Marshal(reqMsg)
		this.checkError(err1)
		_, err2 := service.Gm("ReqMsg", string(data))
		this.checkError(err2)

		service.ActionService.UpdateDiamond(this.auth.GetUser().UserName,
			utils.String(entity.LogType9), userid, utils.String(diamond))
		this.redirect(beego.URLFor("PlayerController.List"))
	} else {
		p, err := service.PlayerService.GetPlayer(userid)
		this.checkError(err)

		this.Data["player"] = p
		this.Data["pageTitle"] = "钻石操作"
		this.display()
	}
}

// 打印房间数据
func (this *PlayerController) Desk() {
	userid := this.GetString("userid")

	reqMsg := &entity.ReqRoomMsg{
		Userid: userid,
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqRoomMsg", string(data))
	this.checkError(err2)

	this.redirect(beego.URLFor("PlayerController.List"))
}

// 绑定代理操作
func (this *PlayerController) Build() {
	userid := this.GetString("userid")
	if this.isPost() {
		agent := this.GetString("agent")
		if userid == "" {
			this.checkError(fmt.Errorf("用户ID不能为空"))
		} else {
			reqMsg := &entity.ReqBuildMsg{
				Userid: userid,
				Agent:  agent,
			}
			data, err1 := json.Marshal(reqMsg)
			this.checkError(err1)
			_, err2 := service.Gm("ReqBuildMsg", string(data))
			this.checkError(err2)

			service.ActionService.UpdateBuild(this.auth.GetUser().UserName, userid, agent)
		}
		this.redirect(beego.URLFor("PlayerController.List"))
	} else {
		p, err := service.PlayerService.GetPlayer(userid)
		this.checkError(err)

		this.Data["player"] = p
		this.Data["pageTitle"] = "绑定操作"
		this.display()
	}
}

// 公告广播
func (this *PlayerController) NoticeList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	// 过期处理
	if status == 0 { //未过期
		m["etime"] = bson.M{"$gte": bson.Now()}
	} else { //已过期
		m["etime"] = bson.M{"$lt": bson.Now()}
	}
	m["del"] = status
	list, _ := service.PlayerService.GetNoticeList(page, this.pageSize, m)
	count, _ := service.PlayerService.GetNoticeListTotal(m)

	this.Data["pageTitle"] = "公告列表"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.NoticeList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 添加公告
func (this *PlayerController) NoticeAdd() {
	if this.isPost() {
		notice := new(entity.Notice)
		rtype, _ := this.GetInt("rtype")
		acttype, _ := this.GetInt("acttype")
		top, _ := this.GetInt("top")
		num, _ := this.GetInt("num")
		content := this.GetString("content")
		etime := this.GetString("end_date")
		e := fmt.Sprintf("%s 23:59:59", etime)
		endTime := utils.Str2Time(e) //过期时间
		//fmt.Println("e : ", e, endTime)
		if endTime.IsZero() {
			this.checkError(errors.New("参数错误"))
		}
		notice.Rtype = rtype
		notice.Acttype = acttype
		notice.Top = top
		notice.Num = num
		notice.Content = content
		notice.Etime = endTime
		err := this.validNotice(notice)
		this.checkError(err)

		err = service.PlayerService.AddNotice(notice)
		this.checkError(err)
		service.ActionService.AddNotice(this.auth.GetUser().UserName, notice.Id)
		this.redirect(beego.URLFor("PlayerController.NoticeList"))
	}

	types1 := map[int]string{
		0: "显示消息",
		1: "支付消息",
		2: "活动消息",
	}

	types2 := map[int]string{
		1: "活动公告",
		2: "广播消息",
	}

	tops := map[int]string{
		0: "否",
		1: "是",
	}

	this.Data["pageTitle"] = "添加公告"
	this.Data["types1"] = types1
	this.Data["types2"] = types2
	this.Data["tops"] = tops
	this.display()
}

func (this *PlayerController) validNotice(notice *entity.Notice) error {
	valid := validation.Validation{}
	valid.Required(notice.Rtype, "rtype").Message("消息类型不能为空")
	valid.Required(notice.Acttype, "acttype").Message("操作类型不能为空")
	valid.Required(notice.Num, "num").Message("公告次数不能为空")
	valid.Required(notice.Content, "content").Message("公告内容不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}

	return nil
}

// 公告广播
func (this *PlayerController) Notice() {
	id := this.GetString("id")

	notice, err := service.PlayerService.GetNotice(id)
	this.checkError(err)

	reqMsg := &entity.ReqNoticeMsg{
		Id:      notice.Id,
		Rtype:   notice.Rtype,
		Acttype: notice.Acttype,
		Top:     notice.Top,
		Num:     notice.Num,
		Del:     notice.Del, //是否移除
		Content: notice.Content,
		Etime:   notice.Etime,
		Ctime:   notice.Ctime,
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqNoticeMsg", string(data))
	this.checkError(err2)

	service.ActionService.Notice(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.NoticeList"))
}

// 移除公告广播
func (this *PlayerController) NoticeDel() {
	id := this.GetString("id")

	notice, err := service.PlayerService.GetNotice(id)
	this.checkError(err)

	reqMsg := &entity.ReqNoticeMsg{
		Id:      notice.Id,
		Rtype:   notice.Rtype,
		Acttype: notice.Acttype,
		Top:     notice.Top,
		Num:     notice.Num,
		Del:     1, //是否移除
		Content: notice.Content,
		Etime:   notice.Etime,
		Ctime:   notice.Ctime,
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqNoticeMsg", string(data))
	this.checkError(err2)

	if err2 == nil {
		err3 := service.PlayerService.DelNotice(id)
		this.checkError(err3)
	}

	service.ActionService.DelNotice(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.NoticeList"))
}

// 商城列表
func (this *PlayerController) ShopList() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	if page < 1 {
		page = 1
	}

	m := service.FindByDate(startDate, endDate, "ctime", "ctime")
	// 过期处理
	if status == 0 { //未过期
		m["etime"] = bson.M{"$gte": bson.Now()}
	} else { //已过期
		m["etime"] = bson.M{"$lt": bson.Now()}
	}
	m["del"] = status
	list, _ := service.PlayerService.GetShopList(page, this.pageSize, m)
	count, _ := service.PlayerService.GetShopListTotal(m)

	this.Data["pageTitle"] = "商城列表"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgencyController.ShopList", "status", status, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 添加公告
func (this *PlayerController) ShopAdd() {
	if this.isPost() {
		status, _ := this.GetInt("status")
		propid, _ := this.GetInt("propid")
		payway, _ := this.GetInt("payway")
		number, _ := this.GetInt("number")
		give, _ := this.GetInt("give")
		price, _ := this.GetInt("price")
		name := this.GetString("name")
		info := this.GetString("info")
		etime := this.GetString("end_date")
		e := fmt.Sprintf("%s 23:59:59", etime)
		endTime := utils.Str2Time(e) //过期时间
		//fmt.Println("e : ", e, endTime)
		if endTime.IsZero() {
			this.checkError(errors.New("参数错误"))
		}
		if give < 0 {
			give = 0
		}
		shop := new(entity.Shop)
		shop.Status = status
		shop.Propid = propid
		shop.Payway = payway
		shop.Number = uint32(number)
		shop.Give = uint32(give)
		shop.Price = uint32(price)
		shop.Name = name
		shop.Info = info
		shop.Etime = endTime
		err := this.validShop(shop)
		this.checkError(err)

		err = service.PlayerService.AddShop(shop)
		this.checkError(err)
		service.ActionService.AddShop(this.auth.GetUser().UserName, shop.Id)
		this.redirect(beego.URLFor("PlayerController.ShopList"))
	}

	types1 := map[int]string{
		1: "热卖",
		2: "普通",
	}

	types2 := map[int]string{
		1: "钻石",
	}

	types3 := map[int]string{
		1: "RMB",
	}

	this.Data["pageTitle"] = "添加商品"
	this.Data["types1"] = types1
	this.Data["types2"] = types2
	this.Data["types3"] = types3
	this.display()
}

func (this *PlayerController) validShop(shop *entity.Shop) error {
	valid := validation.Validation{}
	valid.Required(shop.Name, "name").Message("物品名称不能为空")
	valid.Required(shop.Info, "info").Message("物品描述不能为空")
	valid.Required(shop.Number, "number").Message("购买数量不能为空")
	//valid.Range(shop.Number, 1, 5000, "number").Message("购买数量不对")
	valid.Required(shop.Price, "price").Message("购买价格不能为空")
	//valid.Range(shop.Price, 1, 5000, "price").Message("购买价格不对")
	valid.Required(shop.Propid, "propid").Message("购买的物品不能为空")
	//valid.Range(shop.Propid, 1, 10, "propid").Message("购买的物品不对")
	valid.Required(shop.Payway, "payway").Message("支付方式不能为空")
	//valid.Range(shop.Payway, 1, 10, "payway").Message("支付方式不对")
	valid.Required(shop.Status, "status").Message("物品状态不能为空")
	//valid.Range(shop.Status, 1, 100, "status").Message("物品状态不对")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}

	return nil
}

// 公告广播
func (this *PlayerController) Shop() {
	id := this.GetString("id")

	shop, err := service.PlayerService.GetShop(id)
	this.checkError(err)

	reqMsg := &entity.ReqShopMsg{
		Id:     shop.Id,     //购买ID
		Status: shop.Status, //物品状态,1=热卖
		Propid: shop.Propid, //兑换的物品,1=钻石
		Payway: shop.Payway, //支付方式,1=RMB
		Number: shop.Number, //兑换的数量
		Give:   shop.Give,   //赠送的数量
		Price:  shop.Price,  //支付价格
		Name:   shop.Name,   //物品名字
		Info:   shop.Info,   //物品信息
		Del:    shop.Del,    //是否移除
		Etime:  shop.Etime,  //过期时间
		Ctime:  shop.Ctime,  //创建时间
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqShopMsg", string(data))
	this.checkError(err2)

	service.ActionService.Shop(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.ShopList"))
}

// 移除公告广播
func (this *PlayerController) ShopDel() {
	id := this.GetString("id")

	shop, err := service.PlayerService.GetShop(id)
	this.checkError(err)

	reqMsg := &entity.ReqShopMsg{
		Id:     shop.Id,     //购买ID
		Status: shop.Status, //物品状态,1=热卖
		Propid: shop.Propid, //兑换的物品,1=钻石
		Payway: shop.Payway, //支付方式,1=RMB
		Number: shop.Number, //兑换的数量
		Give:   shop.Give,   //赠送的数量
		Price:  shop.Price,  //支付价格
		Name:   shop.Name,   //物品名字
		Info:   shop.Info,   //物品信息
		Del:    1,           //是否移除
		Etime:  shop.Etime,  //过期时间
		Ctime:  shop.Ctime,  //创建时间
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	_, err2 := service.Gm("ReqShopMsg", string(data))
	this.checkError(err2)

	if err2 == nil {
		err3 := service.PlayerService.DelShop(id)
		this.checkError(err3)
	}

	service.ActionService.DelShop(this.auth.GetUserName(), id)

	this.redirect(beego.URLFor("PlayerController.ShopList"))
}

// 设置变量
func (this *PlayerController) EnvAdd() {

	types1 := entity.EnvTypeValue

	if this.isPost() {
		key, _ := this.GetInt("key")
		value, _ := this.GetInt("value")

		types2 := entity.EnvTypeKey
		if v, ok := types2[key]; ok {
			env := new(entity.ReqEnvMsg)
			env.Key = v
			env.Value = int32(value)
			data, err1 := json.Marshal(env)
			this.checkError(err1)
			_, err2 := service.Gm("ReqEnvMsg", string(data))
			this.checkError(err2)

			service.ActionService.EnvAdd(this.auth.GetUserName(), env.Key)
		} else {
			this.checkError(errors.New("参数错误"))
		}
		this.redirect(beego.URLFor("PlayerController.EnvList"))
	}

	this.Data["pageTitle"] = "添加变量"
	this.Data["types1"] = types1
	this.display()
}

// 删除变量
func (this *PlayerController) EnvDel() {
	key := this.GetString("key")

	env := new(entity.ReqDelEnvMsg)
	env.Key = key
	data, err1 := json.Marshal(env)
	this.checkError(err1)
	_, err2 := service.Gm("ReqDelEnvMsg", string(data))
	this.checkError(err2)

	service.ActionService.EnvDel(this.auth.GetUserName(), key)
	this.redirect(beego.URLFor("PlayerController.EnvList"))
}

// 变量列表
func (this *PlayerController) EnvList() {
	reqMsg := &entity.ReqGetEnvMsg{
		Key: "all",
	}
	data, err1 := json.Marshal(reqMsg)
	this.checkError(err1)
	r, err2 := service.Gm("ReqGetEnvMsg", string(data))
	if err2 != nil {
		//this.checkError(fmt.Errorf("请设置变量"))
	} else {
		resp := new(entity.RespEnvMsg)
		err3 := json.Unmarshal([]byte(r), resp)

		if err3 != nil {
			this.checkError(err3)
		}

		list := make(map[string]entity.ReqEnvMsg)
		for _, v := range resp.List {
			for k2, v2 := range entity.EnvTypeKey {
				if v2 == v.Key {
					list[entity.EnvTypeValue[k2]] = v
					break
				}
			}
		}

		this.Data["list"] = list
	}

	this.Data["pageTitle"] = "环境变量"
	this.display()
}
