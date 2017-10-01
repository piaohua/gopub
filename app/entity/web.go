package entity

import "time"

const (
	ITYPE1 uint32 = 1 //钻石
	ITYPE2 uint32 = 2 //金币
)

//响应信息
type RespErr struct {
	ErrCode int    `json:errcode` //错误码
	ErrMsg  string `json:errmsg`  //错误信息
	Result  string `json:result`  //正常时返回信息
}

//货币变更请求
type ReqMsg struct {
	Userid string `json:userid` //角色ID
	Rtype  int    `json:rtype`  //类型
	Itemid uint32 `json:itemid` //物品,1钻石,2金币
	Amount int32  `json:amount` //数量
}

//在线状态请求
type ReqOnlineStatusMsg struct {
	Userid []string `json:userid` //角色ID
}

type RespOnlineStatusMsg struct {
	Userid map[string]int `json:userid` //角色ID,1在线
}

//发布公告
type ReqNoticeMsg struct {
	Id      string    `json:id`
	Rtype   int       `json:rtype`    //类型
	Acttype int       `json:act_type` //操作类型
	Top     int       `json:top`      //置顶
	Num     int       `json:num`      //广播次数
	Del     int       `json:del`      //是否移除
	Content string    `json:content`  //广播内容
	Etime   time.Time `json:etime`    //过期时间
	Ctime   time.Time `json:ctime`    //创建时间
}

//绑定请求(修改绑定)
type ReqBuildMsg struct {
	Userid string `json:userid` //角色ID
	Agent  string `json:agent`  //代理ID
}

//房间数据
type ReqRoomMsg struct {
	Userid string `json:userid` //角色ID
}

type RespRoomMsg struct {
	Userid     string `json:userid`      //角色ID
	DeskData   string `json:desk_data`   //代理ID
	DeskRecord string `json:desk_record` //代理ID
}

//发布商品
type ReqShopMsg struct {
	Id     string    `json:"id"`     //购买ID
	Status int       `json:"status"` //物品状态,1=热卖
	Propid int       `json:"propid"` //兑换的物品,1=钻石
	Payway int       `json:"payway"` //支付方式,1=RMB
	Number uint32    `json:"number"` //兑换的数量
	Give   uint32    `json:"give"`   //赠送的数量
	Price  uint32    `json:"price"`  //支付价格
	Name   string    `json:"name"`   //物品名字
	Info   string    `json:"info"`   //物品信息
	Del    int       `json:"del"`    //是否移除
	Etime  time.Time `json:"etime"`  //过期时间
	Ctime  time.Time `json:"ctime"`  //创建时间
}

//设置变量
//key      value
//regist   注册赠送
//build    绑定赠送
//agent    绑定赠给代理
type ReqEnvMsg struct {
	Key   string `json:key`   //key
	Value int32  `json:value` //value
}

type ReqGetEnvMsg struct {
	Key string `json:key` //key
}

type ReqDelEnvMsg struct {
	Key string `json:key` //key
}

//设置变量
type RespEnvMsg struct {
	List []ReqEnvMsg `json:list` //list
}
