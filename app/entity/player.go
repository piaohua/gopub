package entity

import "time"

//玩家数据
type PlayerUser struct {
	Userid    string    `bson:"_id"`         // 用户id
	Nickname  string    `bson:"nickname"`    // 用户昵称
	Sex       uint32    `bson:"sex"`         // 用户性别,男1 女2 非男非女3
	Phone     string    `bson:"phone"`       // 绑定的手机号码
	Auth      string    `bson:"auth"`        // 密码验证码
	Pwd       string    `bson:"pwd"`         // MD5密码
	RegIp     string    `bson:"regist_ip"`   // 注册账户时的IP地址
	LoginIp   string    `bson:"login_ip"`    // 登录时的IP地址
	Coin      uint32    `bson:"coin"`        // 金币
	Diamond   uint32    `bson:"diamond"`     // 钻石
	RoomCard  uint32    `bson:"room_card"`   // 房卡
	Status    uint32    `bson:"status"`      // 正常1  锁定2  黑名单3
	Address   string    `bson:"address"`     // 物理地址
	Latitude  string    `bson:"latitude"`    // 纬度
	Lontitude string    `bson:"lontitude"`   // 经度
	Photo     string    `bson:"photo"`       // 头像
	Wxuid     string    `bson:"wxuid"`       // 微信uid
	Win       uint32    `bson:"win"`         // 赢
	Lost      uint32    `bson:"lost"`        // 输
	Ping      uint32    `bson:"ping"`        // 平
	Piao      uint32    `bson:"piao"`        // 漂
	Robot     bool      `bson:"robot"`       // 是否是机器人
	Money     uint32    `bson:"money"`       // 充值总金额(分)
	TopDia    uint32    `bson:"top_diamond"` // 钻石总金额
	Agent     string    `bson:"agent"`       // 代理ID
	Xftoken   string    `bson:"xftoken"`     // xftoken
	Atime     time.Time `bson:"agent_time"`  // 绑定代理时间
	Ctime     time.Time `bson:"create_time"` // 注册时间
	State     int       // 在线状态
}

const (
	TradeSuccess = 0 //交易成功
	TradeFail    = 1 //交易失败
	Tradeing     = 2 //交易中(下单状态)
	TradeGoods   = 3 //发货失败
)

var TradeResult = map[int]string{
	TradeSuccess: "成功",
	TradeFail:    "交易失败",
	//Tradeing:     "交易中",
	TradeGoods: "发货失败",
}

// 交易记录
type TradeRecord struct {
	Id        string    `bson:"_id"`       //商户订单号(游戏内自定义订单号)
	Transid   string    `bson:"transid"`   //交易流水号(计费支付平台的交易流水号,微信订单号)
	Userid    string    `bson:"userid"`    //用户在商户应用的唯一标识(userid)
	Itemid    string    `bson:"itemid"`    //购买商品ID
	Amount    string    `bson:"amount"`    //购买商品数量
	Diamond   uint32    `bson:"diamond"`   //购买钻石数量
	Money     uint32    `bson:"money"`     //交易总金额(单位为分)
	Transtime string    `bson:"transtime"` //交易完成时间 yyyy-mm-dd hh24:mi:ss
	Result    int       `bson:"result"`    //交易结果(0–交易成功,1–交易失败,2-交易中,3-发货中)
	Waresid   uint32    `bson:"waresid"`   //商品编码(平台为应用内需计费商品分配的编码)
	Currency  string    `bson:"currency"`  //货币类型(RMB,CNY)
	Transtype int       `bson:"transtype"` //交易类型(0–支付交易)
	Feetype   int       `bson:"feetype"`   //计费方式(表示商品采用的计费方式)
	Paytype   uint32    `bson:"paytype"`   //支付方式(表示用户采用的支付方式,403-微信支付)
	Clientip  string    `bson:"clientip"`  //客户端ip
	Agent     string    `bson:"agent"`     //绑定的父级代理商游戏ID
	Ctime     time.Time `bson:"ctime"`     //本条记录生成unix时间戳
}

const (
	NOTICE_TYPE1 = 1 //活动公告
	NOTICE_TYPE2 = 2 //广播消息
)

const (
	NOTICE_ACT_TYPE0 = 0 //无操作消息
	NOTICE_ACT_TYPE1 = 1 //支付消息
	NOTICE_ACT_TYPE2 = 2 //活动消息
)

//公告
type Notice struct {
	Id      string    `bson:"_id"`
	Rtype   int       `bson:"rtype"`    //类型,1=公告消息,2=广播消息
	Acttype int       `bson:"act_type"` //操作类型,0=无操作,1=支付,2=活动
	Top     int       `bson:"top"`      //置顶
	Num     int       `bson:"num"`      //广播次数
	Del     int       `bson:"del"`      //是否移除
	Content string    `bson:"content"`  //广播内容
	Etime   time.Time `bson:"etime"`    //过期时间
	Ctime   time.Time `bson:"ctime"`    //创建时间
}

//商城
type Shop struct {
	Id     string    `bson:"_id"`    //购买ID
	Status int       `bson:"status"` //物品状态,1=热卖
	Propid int       `bson:"propid"` //兑换的物品,1=钻石
	Payway int       `bson:"payway"` //支付方式,1=RMB
	Number uint32    `bson:"number"` //兑换的数量
	Give   uint32    `bson:"give"`   //赠送的数量
	Price  uint32    `bson:"price"`  //支付价格(单位元)
	Name   string    `bson:"name"`   //物品名字
	Info   string    `bson:"info"`   //物品信息
	Del    int       `bson:"del"`    //是否移除
	Etime  time.Time `bson:"etime"`  //过期时间
	Ctime  time.Time `bson:"ctime"`  //创建时间
}

const (
	EnvType1 = 1 //注册赠送
	EnvType2 = 2 //绑定赠送
	EnvType3 = 3 //绑定赠送给代理
	EnvType4 = 4 //钻石房间抽成
	EnvType5 = 5 //钻石房间进入限制
)

var EnvTypeValue = map[int]string{
	EnvType1: "注册赠送",
	EnvType2: "绑定赠送",
	EnvType3: "绑定赠送给代理",
	EnvType4: "钻石房间抽成",
	EnvType5: "钻石房间进入限制",
}

var EnvTypeKey = map[int]string{
	EnvType1: "regist",
	EnvType2: "build",
	EnvType3: "agent",
	EnvType4: "diamond_deduct",
	EnvType5: "diamond_enter",
}
