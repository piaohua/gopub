package service

import (
	"errors"
	"fmt"
	"gopub/app/entity"
	"time"
	"utils"

	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2/bson"
)

type agencyService struct{}

// 获取代理商列表
func (this *agencyService) GetAgencyList(page, pageSize int, m bson.M) ([]entity.User, error) {
	var list []entity.User
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "update_time", false)
	m["agent"] = bson.M{"$ne": ""}
	err := Users.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, err
}

// 获取代理商总数
func (this *agencyService) GetAgencyTotal(m bson.M) (int64, error) {
	return int64(Count(Users, m)), nil
}

// 获取代理商类型
func (this *agencyService) GetAgencyType() ([]int, error) {
	var types []int
	ListByQ(Users, bson.M{"$group": bson.M{"level": "$level"}}, &types)
	return types, nil
}

// 获取邀请列表
func (this *agencyService) GetInviterList(page, pageSize int, m bson.M) ([]entity.User, error) {
	var list []entity.User
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "update_time", false)
	m["agent"] = bson.M{"$ne": ""}
	err := Users.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, err
}

// 获取邀请总数
func (this *agencyService) GetInviterTotal(m bson.M) (int64, error) {
	m["agent"] = bson.M{"$ne": ""}
	return int64(Count(Users, m)), nil
}

// 获取绑定我的玩家列表
func (this *agencyService) GetMyAgencyList(username string, page, pageSize int, m bson.M) ([]entity.PlayerUser, error) {
	var list []entity.PlayerUser
	if pageSize == -1 {
		pageSize = 100000
	}
	if username == "" {
		return list, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return list, err
	}
	if agency.Agent == "" {
		return list, nil
	}
	m["agent"] = agency.Agent
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "create_time", false)
	PlayerUsers.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取绑定我的玩家总数
func (this *agencyService) GetMyAgencyTotal(username string, m bson.M) (int64, error) {
	var count int64
	if username == "" {
		return count, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil || agency == nil {
		return count, err
	}
	return int64(agency.Builds), nil
}

// 获取全部提现记录
func (this *agencyService) GetCashList(page, pageSize int, m bson.M) ([]entity.ApplyCash, error) {
	var list []entity.ApplyCash
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)

	ApplyCashs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取全部提现记录总数
func (this *agencyService) GetCashListTotal(m bson.M) (int64, error) {
	var count int64
	count = int64(Count(ApplyCashs, m))
	return count, nil
}

// 获取我的提现
func (this *agencyService) GetMyCashList(username string, page, pageSize int, m bson.M) ([]entity.ApplyCash, error) {
	var list []entity.ApplyCash
	if pageSize == -1 {
		pageSize = 100000
	}
	if username == "" {
		return list, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return list, err
	}
	if agency.Agent == "" {
		return list, nil
	}

	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	m["agent"] = agency.Agent

	ApplyCashs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取我的提现记录总数
func (this *agencyService) GetMyCashListTotal(username string, m bson.M) (int64, error) {
	var count int64
	if username == "" {
		return count, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return count, err
	}
	if agency.Agent == "" {
		return count, nil
	}
	m["agent"] = agency.Agent
	count = int64(Count(ApplyCashs, m))
	return count, nil
}

// 获取我的总的已提现金额
func (this *agencyService) GetMyExtractTotal(username string) (float32, error) {
	var count float32
	if username == "" {
		return count, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return count, err
	}
	if agency.Agent == "" {
		return count, nil
	}
	return agency.Extract, nil
}

//代理获取可以提现金额
func (this *agencyService) GetMyCashTotal(username string) (float32, error) {
	if username == "" {
		return 0, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return 0, err
	}
	if agency.Agent == "" {
		return 0, errors.New("代理不存在")
	}
	return agency.Cash, nil
}

// 申请提现
func (this *agencyService) ApplyCashAdd(username, name, bankAddr string, bankCard int, cash float64) error {
	agency, err := UserService.GetUserByName(username)
	if agency.Agent == "" || err != nil {
		return errors.New("代理ID不存在")
	}
	if cash > float64(agency.Cash) {
		return errors.New("金额不足")
	}
	applyCash := new(entity.ApplyCash)
	applyCash.Id = bson.NewObjectId().Hex()
	applyCash.Agent = agency.Agent
	applyCash.Cash = float32(cash)
	applyCash.Status = 1 //表示等待处理
	applyCash.RealName = name
	applyCash.BankCard = bankCard
	applyCash.BankAddr = bankAddr
	applyCash.Ctime = bson.Now()
	if Insert(ApplyCashs, applyCash) {
		err1 := this.updateCash(agency.Id, (-1 * applyCash.Cash))
		if err1 != nil {
			Delete(ApplyCashs, bson.M{"_id": applyCash.Id})
			return errors.New("提取失败")
		}
		return nil
	}
	return errors.New("提取失败")
}

// 提现处理
func (this *agencyService) ExtractCash(username, orderid string) error {
	agency, err := UserService.GetUserByName(username)
	if err != nil || agency == nil {
		return err
	}
	//if agency.Agent == "" {
	//	return errors.New("代理ID不存在")
	//}
	m := bson.M{"_id": orderid}
	n := bson.M{"status": 0, "user_name": username, "utime": bson.Now()}
	if Update(ApplyCashs, m, bson.M{"$set": n}) {
		return nil
	}
	return errors.New("提现失败")
}

// 更新提现率
func (this *agencyService) UpdateAgencyRate(agency *entity.User) error {
	m := bson.M{"user_name": agency.UserName, "agent": agency.Agent}
	n := bson.M{"rate": agency.Rate, "update_time": bson.Now()}
	if Update(Users, m, bson.M{"$set": n}) {
		return nil
	}
	return errors.New("更新失败")
}

// 更新提现金额
func (this *agencyService) updateCash(id string, cash float32) error {
	m := bson.M{"_id": id}
	n := bson.M{"update_time": bson.Now()}
	c := bson.M{"cash": cash}
	if Update(Users, m, bson.M{"$set": n, "$inc": c}) {
		return nil
	}
	return errors.New("更新失败")
}

// 更新提现金额
func (this *agencyService) updateAgencyCash(agency *entity.User, cash float32) error {
	m := bson.M{"_id": agency.Id}
	n := bson.M{"cash_time": agency.CashTime, "update_time": bson.Now()}
	c := bson.M{"cash": cash}
	if Update(Users, m, bson.M{"$set": n, "$inc": c}) {
		return nil
	}
	return errors.New("更新失败")
}

// 定时统计
func (this *agencyService) stat() {
	list, err := this.GetAgencyList(1, -1, bson.M{})
	if err != nil {
		beego.Trace("stat err: ", err)
		fmt.Println("stat err: ", err)
	}
	for _, v := range list {
		this.statBuilds(&v) //绑定更新
		this.statCash(&v)   //提现更新
	}
}

// 获取我的总的可提现金额
func (this *agencyService) statCash(agency *entity.User) {
	if agency.Agent == "" {
		return
	}
	endTime := utils.LocalTime()
	startTime := agency.CashTime
	//统计属于代理的
	money, err2 := this.statAgentCash(startTime, endTime, agency.Agent)
	if err2 != nil {
		beego.Trace("statCash err2: ", agency.Id, err2)
		fmt.Println("statCash err2: ", agency.Id, err2)
		return
	}
	money = money * 0.6 //返现60%
	money = float32(utils.Float64(fmt.Sprintf("%.2f", money)))

	if money == 0 {
		return
	}
	//存在数据中
	agency.CashTime = endTime //截至统计时间
	err1 := this.updateAgencyCash(agency, money)
	if err1 != nil {
		beego.Trace("statCash err1: ", agency.Id, err1)
		fmt.Println("statCash err1: ", agency.Id, err1)
		return
	}
	beego.Trace("statCash ok: ", agency.Id, money)
	fmt.Println("statCash ok: ", agency.Id, money)
}

// 获取绑定我的玩家总数
func (this *agencyService) statBuilds(agency *entity.User) {
	if agency.Agent == "" {
		return
	}
	endTime := utils.LocalTime()
	startTime := agency.BuildsTime
	m := bson.M{"agent": agency.Agent}
	m["agent_time"] = bson.M{"$gte": startTime, "$lt": endTime}
	count, _ := PlayerService.GetTotal(0, m)
	if count == 0 {
		return
	}
	// 更新绑定人数
	q := bson.M{"_id": agency.Id}
	n := bson.M{"builds_time": endTime}
	c := bson.M{"builds": uint32(count)}
	if !Update(Users, q, bson.M{"$set": n, "$inc": c}) {
		beego.Trace("statBuilds err: ", agency.Id, count)
		fmt.Println("statBuilds err: ", agency.Id, count)
		return
	}
	beego.Trace("statBuilds ok: ", agency.Id, count)
	fmt.Println("statBuilds ok: ", agency.Id, count)
}

//统计属于代理的
func (this *agencyService) statAgentCash(startTime, endTime time.Time, agent string) (float32, error) {
	m := bson.M{
		"$match": bson.M{
			"agent":  agent,
			"result": entity.TradeSuccess,
			"ctime":  bson.M{"$gte": startTime, "$lt": endTime},
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": "$agent",
			"money": bson.M{
				"$sum": "$money",
			},
		},
	}
	operations := []bson.M{m, n}
	result := bson.M{}
	pipe := TradeRecords.Pipe(operations)
	err := pipe.One(&result)
	if err != nil {
		if err.Error() == "not found" {
			return 0, nil
		}
		return 0, err
	}
	if v, ok := result["money"]; ok {
		return float32(v.(int)), nil
	}
	return 0, nil
}
