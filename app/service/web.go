package service

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"gopub/app/entity"
	"net/http"
	"net/url"
	"utils"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

var gmHost string
var gmPort string
var gmPath string
var gmKey string

func init() {
	gmHost = beego.AppConfig.String("gm.host")
	gmPort = beego.AppConfig.String("gm.port")
	gmPath = beego.AppConfig.String("gm.path")
	gmKey = beego.AppConfig.String("gm.key")
}

var cstDialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//Gm操作,角色ID,操作类型,操作物品,操作数量
func Gm(msgName, msg string) (string, error) {
	addr := gmHost + ":" + gmPort
	//fmt.Println("addr -> ", addr)
	u := url.URL{Scheme: "wss", Host: addr, Path: gmPath}
	TimeStr := GmTime()
	Token := GmToke(TimeStr)
	//c, _, err := websocket.DefaultDialer.Dial(u.String(),
	//	http.Header{"Token": {Token}})
	d := cstDialer
	d.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c, _, err := d.Dial(u.String(),
		http.Header{"Token": {Token}})
	if err != nil {
		return "", errors.New(fmt.Sprintf("dial err -> %v", err))
	}
	sign := GmSign(msg, TimeStr)
	msg2 := GmMsg(msg, msgName, sign, TimeStr)
	if c != nil {
		c.WriteMessage(websocket.TextMessage, []byte(msg2))
		defer c.Close()
		_, message, err := c.ReadMessage()
		if err != nil {
			return "", errors.New(fmt.Sprintf("read err -> %v", err))
		}
		resp := new(entity.RespErr)
		err = json.Unmarshal(message, resp)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Unmarshal err -> %v", err))
		}
		if resp.ErrCode != 0 {
			return "", errors.New(fmt.Sprintf("resp.ErrMsgi %s", resp.ErrMsg))
		}
		return resp.Result, nil
	}
	return "", errors.New(fmt.Sprintf("c empty err -> %v", err))
}

//字符串时间
func GmTime() string {
	Time := utils.Timestamp()
	TimeStr := utils.String(Time)
	return TimeStr
}

// Sign := utils.Md5(Key+Now)
// Token := Sign+Now+RandNum
func GmToke(TimeStr string) string {
	Sign := utils.Md5(gmKey + TimeStr)
	Token := Sign + TimeStr + utils.RandStr(6)
	return Token
}

// Sign := TimeStr + Key + Md5(msg)
func GmSign(msg, TimeStr string) string {
	return utils.Md5(TimeStr + gmKey + utils.Md5(msg))
}

// Timestr|Sign|msg_name|msg
func GmMsg(msg, msgName, sign, TimeStr string) string {
	return TimeStr + "|" + sign + "|" + msgName + "|" + msg
}
