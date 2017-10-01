package entity

import "time"

//1注册赠送,2开房消耗,3房间解散返还
//4充值购买
//5商城购买
//6绑定赠送
//9后台操作
//10钻石房间
//11钻石房间抽成
const (
	LogType1  = 1
	LogType2  = 2
	LogType3  = 3
	LogType4  = 4
	LogType5  = 5
	LogType6  = 6
	LogType9  = 9
	LogType10 = 10
	LogType11 = 11
)

var LogType = map[int]string{
	0:         "全部",
	LogType1:  "注册赠送",
	LogType2:  "开房消耗",
	LogType3:  "房间解散返还",
	LogType4:  "充值购买",
	LogType5:  "商城购买",
	LogType6:  "绑定赠送",
	LogType9:  "后台操作",
	LogType10: "钻石房间",
	LogType11: "钻石房间抽成",
}

//注册日志
type LogRegist struct {
	Id       string    `bson:"_id"`
	Userid   string    `bson:"userid"`    //账户ID
	Nickname string    `bson:"nickname"`  //账户名称
	Ip       string    `bson:"ip"`        //注册IP
	DayStamp time.Time `bson:"day_stamp"` //regist Time Today
	DayDate  int       `bson:"day_date"`  //regist day date
	Ctime    time.Time `bson:"ctime"`     //create Time
}

//登录日志
type LogLogin struct {
	Id         string    `bson:"_id"`
	Userid     string    `bson:"userid"`      //账户ID
	Event      int       `bson:"event"`       //事件：0=登录,1=正常退出,2＝系统关闭时被迫退出,3＝被动退出,4＝其它情况导致的退出
	Ip         string    `bson:"ip"`          //登录IP
	DayStamp   time.Time `bson:"day_stamp"`   //login Time Today
	LoginTime  time.Time `bson:"login_time"`  //login Time
	LogoutTime time.Time `bson:"logout_time"` //logout Time
}

//钻石日志
type LogDiamond struct {
	Id     string    `bson:"_id"`
	Userid string    `bson:"userid"` //账户ID
	Type   int       `bson:"type"`   //类型
	Num    int32     `bson:"num"`    //数量
	Rest   uint32    `bson:"rest"`   //剩余数量
	Ctime  time.Time `bson:"ctime"`  //create Time
}

//金币日志
type LogCoin struct {
	Id     string    `bson:"_id"`
	Userid string    `bson:"userid"` //账户ID
	Type   int       `bson:"type"`   //类型
	Num    int32     `bson:"num"`    //数量
	Rest   uint32    `bson:"rest"`   //剩余数量
	Ctime  time.Time `bson:"ctime"`  //create Time
}

//绑定日志
type LogBuildAgency struct {
	Id       string    `bson:"_id"`
	Userid   string    `bson:"userid"`    //账户ID
	Agent    string    `bson:"agent"`     //绑定ID
	DayStamp time.Time `bson:"day_stamp"` //regist Time Today
	Day      int       `bson:"day"`       //regist day
	Month    int       `bson:"month"`     //regist month
	Ctime    time.Time `bson:"ctime"`     //create Time
}

//在线日志
type LogOnline struct {
	Id       string    `bson:"_id"`
	Num      int       `bson:"num"`       //online count
	DayStamp time.Time `bson:"day_stamp"` //Time Today
	Ctime    time.Time `bson:"ctime"`     //create Time
}
