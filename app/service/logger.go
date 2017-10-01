package service

import (
	"errors"
	"fmt"
	"gopub/app/entity"
	"time"
	"utils"

	"gopkg.in/mgo.v2/bson"
)

type loggerService struct{}

// 获取注册日志列表
func (this *loggerService) GetRegistList(page, pageSize int, m bson.M) ([]entity.LogRegist, error) {
	var list []entity.LogRegist
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	RegistLogs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取注册日志总数
func (this *loggerService) GetRegistTotal(m bson.M) (int64, error) {
	return int64(Count(RegistLogs, m)), nil
}

// 获取登录日志列表
func (this *loggerService) GetLoginList(page, pageSize int, m bson.M) ([]entity.LogLogin, error) {
	var list []entity.LogLogin
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "login_time", false)
	LoginLogs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取登录日志总数
func (this *loggerService) GetLoginTotal(m bson.M) (int64, error) {
	return int64(Count(LoginLogs, m)), nil
}

// 获取充值列表
func (this *loggerService) GetPayList(page, pageSize int, m bson.M) ([]entity.TradeRecord, error) {
	var list []entity.TradeRecord
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	m["result"] = bson.M{"$ne": entity.Tradeing}
	TradeRecords.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取绑定代理充值数量
func (this *loggerService) GetAgencyPayTotal(username string) (int64, error) {
	var count int64
	if username == "" {
		return count, nil
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return count, err
	}
	if agency.Agent == "" {
		return count, errors.New("代理不存在")
	}
	return int64(Count(TradeRecords, bson.M{"agent": agency.Agent, "result": entity.TradeSuccess})), nil
}

// 获取充值总数
func (this *loggerService) GetPayTotal(m bson.M) (int64, error) {
	var count int64
	m["result"] = bson.M{"$ne": entity.Tradeing}
	count = int64(Count(TradeRecords, m))
	return count, nil
}

// 获取钻石日志列表
func (this *loggerService) GetDiamondList(page, pageSize int, m bson.M) ([]entity.LogDiamond, error) {
	var list []entity.LogDiamond
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	DiamondLogs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取注册日志总数
func (this *loggerService) GetDiamondTotal(m bson.M) (int64, error) {
	return int64(Count(DiamondLogs, m)), nil
}

// 获取金币日志列表
func (this *loggerService) GetCoinList(page, pageSize int, m bson.M) ([]entity.LogCoin, error) {
	var list []entity.LogCoin
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	CoinLogs.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取金币日志总数
func (this *loggerService) GetCoinTotal(m bson.M) (int64, error) {
	return int64(Count(CoinLogs, m)), nil
}

// 获取绑定日志列表
func (this *loggerService) GetBuildList(page, pageSize int, m bson.M) ([]entity.LogBuildAgency, error) {
	var list []entity.LogBuildAgency
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	LogBuildAgencys.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, nil
}

// 获取注册日志总数
func (this *loggerService) GetBuildTotal(m bson.M) (int64, error) {
	return int64(Count(LogBuildAgencys, m)), nil
}

// 绑定统计
func (this *loggerService) GetPubStat(username, rangeType string) map[int]int {
	result := make(map[int]int)

	if username == "" {
		return result
	}
	agency, err := UserService.GetUserByName(username)
	if err != nil {
		return result
	}
	if agency.Agent == "" {
		return result
	}

	m := bson.M{}
	var n bson.M

	switch rangeType {
	case "this_month":
		year, month, _ := time.Now().Date()
		startTimeStr := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
		endTimeStr := fmt.Sprintf("%d-%02d-31 23:59:59", year, month)
		startTime, _ := utils.Str2Unix(startTimeStr)
		endTime, _ := utils.Str2Unix(endTimeStr)
		m = bson.M{
			"$match": bson.M{
				"agent":     agency.Agent,
				"day_stamp": bson.M{"$gte": startTime, "$lte": endTime},
			},
		}
		n = bson.M{
			"$group": bson.M{
				"date": "$day",
				"count": bson.M{
					"$sum": 1,
				},
			},
		}
	case "last_month":
		year, month, _ := time.Now().AddDate(0, -1, 0).Date()
		startTimeStr := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
		endTimeStr := fmt.Sprintf("%d-%02d-31 23:59:59", year, month)
		startTime, _ := utils.Str2Unix(startTimeStr)
		endTime, _ := utils.Str2Unix(endTimeStr)
		m = bson.M{
			"$match": bson.M{
				"agent":     agency.Agent,
				"day_stamp": bson.M{"$gte": startTime, "$lte": endTime},
			},
		}
		n = bson.M{
			"$group": bson.M{
				"date": "$day",
				"count": bson.M{
					"$sum": 1,
				},
			},
		}
	case "this_year":
		year := time.Now().Year()
		startTimeStr := fmt.Sprintf("%d-01-01 00:00:00", year)
		endTimeStr := fmt.Sprintf("%d-12-31 23:59:59", year)
		startTime, _ := utils.Str2Unix(startTimeStr)
		endTime, _ := utils.Str2Unix(endTimeStr)
		m = bson.M{
			"$match": bson.M{
				"agent":     agency.Agent,
				"day_stamp": bson.M{"$gte": startTime, "$lte": endTime},
			},
		}
		n = bson.M{
			"$group": bson.M{
				"date": "$month",
				"count": bson.M{
					"$sum": 1,
				},
			},
		}
	case "last_year":
		year := time.Now().Year() - 1
		startTimeStr := fmt.Sprintf("%d-01-01 00:00:00", year)
		endTimeStr := fmt.Sprintf("%d-12-31 23:59:59", year)
		startTime, _ := utils.Str2Unix(startTimeStr)
		endTime, _ := utils.Str2Unix(endTimeStr)
		m = bson.M{
			"$match": bson.M{
				"agent":     agency.Agent,
				"day_stamp": bson.M{"$gte": startTime, "$lte": endTime},
			},
		}
		n = bson.M{
			"$group": bson.M{
				"date": "$month",
				"count": bson.M{
					"$sum": 1,
				},
			},
		}
	}

	operations := []bson.M{m, n}
	maps := []bson.M{}
	pipe := LogBuildAgencys.Pipe(operations)
	err1 := pipe.All(&maps)

	if err1 == nil && len(maps) > 0 {
		for _, v := range maps {
			date, _ := utils.Str2Int(v["date"].(string))
			count, _ := utils.Str2Int(v["count"].(string))
			result[date] = count
		}
	}
	return result
}
