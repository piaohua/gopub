package service

import (
	"errors"
	"gopub/app/entity"

	"gopkg.in/mgo.v2/bson"
)

type playerService struct{}

// 获取
func (this *playerService) GetPlayer(userid string) (*entity.PlayerUser, error) {
	player := new(entity.PlayerUser)
	Get(PlayerUsers, userid, player)
	if player.Userid == "" {
		return player, errors.New("用户不存在")
	}
	return player, nil
}

// 获取所有
func (this *playerService) GetAllPlayer() ([]entity.PlayerUser, error) {
	return this.GetList(0, 1, -1, bson.M{})
}

// 获取列表
func (this *playerService) GetList(typeId, page, pageSize int, m bson.M) ([]entity.PlayerUser, error) {
	var list []entity.PlayerUser
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "_id", false)
	switch typeId {
	case 0: //微信
		m["phone"] = ""
	case 1: //全部玩家
	case 2: //机器人
		m["phone"] = bson.M{"$ne": ""}
	}
	err := PlayerUsers.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, err
}

// 获取总数
func (this *playerService) GetTotal(typeId int, m bson.M) (int64, error) {
	switch typeId {
	case 0: //微信
		m["phone"] = ""
	case 1: //全部玩家
	case 2: //机器人
		m["phone"] = bson.M{"$ne": ""}
	}
	return int64(Count(PlayerUsers, m)), nil
}

// 获取类型
func (this *playerService) GetPlayerType() ([]int, error) {
	var types []int
	types = []int{1}
	return types, nil
}

// 添加公告
func (this *playerService) AddNotice(notice *entity.Notice) error {
	notice.Id = bson.NewObjectId().Hex()
	notice.Ctime = bson.Now()
	if !Insert(Notices, notice) {
		return errors.New("写入失败:" + notice.Id)
	}
	return nil
}

// 获取
func (this *playerService) GetNotice(id string) (*entity.Notice, error) {
	notice := new(entity.Notice)
	Get(Notices, id, notice)
	if notice.Id == "" {
		return notice, errors.New("公告不存在")
	}
	return notice, nil
}

// 获取
func (this *playerService) DelNotice(id string) error {
	if Update(Notices, bson.M{"_id": id}, bson.M{"$set": bson.M{"del": 1}}) {
		return nil
	}
	return errors.New("移除失败")
}

// 获取列表
func (this *playerService) GetNoticeList(page, pageSize int, m bson.M) ([]entity.Notice, error) {
	var list []entity.Notice
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	err := Notices.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, err
}

// 获取总数
func (this *playerService) GetNoticeListTotal(m bson.M) (int64, error) {
	return int64(Count(Notices, m)), nil
}

// 添加商品
func (this *playerService) AddShop(shop *entity.Shop) error {
	shop.Id = bson.NewObjectId().Hex()
	shop.Ctime = bson.Now()
	if !Insert(Shops, shop) {
		return errors.New("写入失败:" + shop.Id)
	}
	return nil
}

// 获取商品
func (this *playerService) GetShop(id string) (*entity.Shop, error) {
	shop := new(entity.Shop)
	Get(Shops, id, shop)
	if shop.Id == "" {
		return shop, errors.New("商品不存在")
	}
	return shop, nil
}

// 移除商品
func (this *playerService) DelShop(id string) error {
	if Update(Shops, bson.M{"_id": id}, bson.M{"$set": bson.M{"del": 1}}) {
		return nil
	}
	return errors.New("移除失败")
}

// 获取列表
func (this *playerService) GetShopList(page, pageSize int, m bson.M) ([]entity.Shop, error) {
	var list []entity.Shop
	if pageSize == -1 {
		pageSize = 100000
	}
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	err := Shops.
		Find(m).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	return list, err
}

// 获取总数
func (this *playerService) GetShopListTotal(m bson.M) (int64, error) {
	return int64(Count(Shops, m)), nil
}
