package entity

import "time"

// 服务器名称 | 服务器ID | Host | Port | Status | 提示信息
// Status: 1=正常,2=火爆,3=维护中,4=新区,5=等待开启
type Server struct {
	Name      string `json:"name"`      //服务器名称
	Host      string `json:"host"`      //服务器地址
	Port      string `json:"port"`      //服务器端口
	Status    string `json:"status"`    //服务器状态
	Info      string `json:"info"`      //服务器提示信息
	UpdateUrl string `json:"updateUrl"` //客户端下载地址
	UpdateDes string `json:"updateDes"` //客户端更新信息
	Version   string `json:"version"`   //客户端版本
}

// 文件
type ListFiles struct {
	Name string    `json:"name"`
	Size string    `json:"size"`
	Time time.Time `json:"time"`
}

type Srv struct {
	Name         string `json:"name"`         //服务器名称
	Host         string `json:"host"`         //服务器地址
	Port         string `json:"port"`         //服务器端口
	Status       string `json:"status"`       //服务器状态
	Info         string `json:"info"`         //服务器提示信息
	UpdateUrl    string `json:"updateUrl"`    //客户端下载地址
	UpdateDes    string `json:"updateDes"`    //客户端更新信息
	Version      string `json:"version"`      //android客户端版本
	VersionBeta  string `json:"versionBeta"`  //客户端内测版本号
	Userids      string `json:"userids"`      //白名单玩家id
	VersionIos   string `json:"versionIos"`   //ios客户端版本
	UpdateUrlIos string `json:"updateUrlIos"` //ios客户端下载地址
}
