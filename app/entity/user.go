package entity

import "time"

const (
	USER_STATUS0 = 0  //正常
	USER_STATUS1 = -1 //禁用
)

// 账号自增id
type UserIDGen struct {
	Id         string `bson:"_id"`
	LastUserId string `bson:"last_user_id"`
}

// 账号
type User struct {
	Id         string    `bson:"_id"`         // AUTO_INCREMENT, PRIMARY KEY (`id`),
	UserName   string    `bson:"user_name"`   // 用户名, UNIQUE KEY `user_name` (`user_name`)
	Password   string    `bson:"password"`    // 密码
	Salt       string    `bson:"salt"`        // 密码盐
	Sex        int       `bson:"sex"`         // 性别
	Email      string    `bson:"email"`       // 邮箱
	LastLogin  time.Time `bson:"last_login"`  // 最后登录时间
	LastIp     string    `bson:"last_ip"`     // 最后登录IP
	Status     int       `bson:"status"`      // 状态，0正常 -1禁用
	CreateTime time.Time `bson:"create_time"` // 创建时间
	UpdateTime time.Time `bson:"update_time"` // 更新时间
	RoleList   []Role    `bson:"role_list"`   // 角色列表
	//代理
	Phone      string    `bson:"phone"`       //绑定的手机号码(备用:非手机号注册时或多个手机时)
	Agent      string    `bson:"agent"`       //代理ID==Userid
	Level      int       `bson:"level"`       //代理等级ID:1级,2级...
	Weixin     string    `bson:"weixin"`      //微信ID
	Alipay     string    `bson:"alipay"`      //支付宝ID
	QQ         string    `bson:"qq"`          //qq号码
	Address    string    `bson:"address"`     //详细地址
	Number     uint32    `bson:"number"`      //当前余额
	Expend     uint32    `bson:"expend"`      //总消耗
	Cash       float32   `bson:"cash"`        //当前可提取额(分)
	Extract    float32   `bson:"extract"`     //已经提取额(分)
	Rate       uint32    `bson:"rate"`        //提现率,可配置，百分值(比如:80表示80%_)
	Builds     uint32    `bson:"builds"`      //绑定我的人数
	BuildsTime time.Time `bson:"builds_time"` //统计指定时间前所有
	CashTime   time.Time `bson:"cash_time"`   //提取指定时间前所有
	//邀请方,或者创建者
	Inviter string `bson:"inviter"` //邀请人
}

// 账号属于分组(可属于多个分组)
type UserRole struct {
	Id     string `bson:"_id"`     // UNIQUE KEY `user_id` (`user_id`,`role_id`)
	UserId string `bson:"user_id"` // 用户id
	RoleId string `bson:"role_id"` // 角色id
}
