package controllers

import (
	"encoding/json"
	"gopub/app/entity"
	"gopub/app/service"
	"utils"
)

type ChartsController struct {
	BaseController
}

// 在线统计
func (this *ChartsController) Online() {

	list, _ := service.ChartsService.GetOnlineList(1, this.pageSize)
	//var list []entity.LogOnline
	//for i := 0; i < 5; i++ {
	//	l := entity.LogOnline{
	//		Num:   i,
	//		Ctime: utils.LocalTime(),
	//	}
	//	list = append(list, l)
	//}

	data1 := make([]entity.ChartData, 0)
	for _, v := range list {
		d := entity.ChartData{
			Label: utils.Format("Y-m-d H:i:s", v.Ctime),
			Value: utils.String(v.Num),
		}
		data1 = append(data1, d)
	}
	data, _ := json.Marshal(data1)
	/*
			data := `[
		              {
		                "label": "Mon",
		                "value": "4123"
		              },
		              {
		                "label": "Tue",
		                "value": "4633"
		              },
		              {
		                "label": "Wed",
		                "value": "5507"
		              },
		              {
		                "label": "Thu",
		                "value": "4910"
		              },
		              {
		                "label": "Fri",
		                "value": "5529"
		              },
		              {
		                "label": "Sat",
		                "value": "5803"
		              },
		              {
		                "label": "Sun",
		                "value": "6202"
		              }
		            ]`
	*/

	this.Data["pageTitle"] = "在线统计"
	this.Data["data"] = string(data)
	this.display()
}
