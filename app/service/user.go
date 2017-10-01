package service

import (
	"errors"
	"gopub/app/entity"
	"gopub/app/libs"
	"utils"

	"gopkg.in/mgo.v2/bson"
)

type userService struct{}

// 生成ID
func (this *userService) GetUserIDGen() (string, error) {
	gen := new(entity.UserIDGen)
	gen.Id = "last_user_id"
	Get(GenIDs, gen.Id, gen)
	if gen.LastUserId == "" {
		gen.LastUserId = "2"
	}
	id := gen.LastUserId
	gen.LastUserId = utils.StringAdd(id)
	if Upsert(GenIDs, bson.M{"_id": gen.Id}, gen) {
		return id, nil
	}
	return id, errors.New("生成错误")
}

// 根据用户id获取一个用户信息
func (this *userService) GetUser(userId string, getRoleInfo bool) (*entity.User, error) {
	user := new(entity.User)
	Get(Users, userId, user)
	if user.Id != "" && getRoleInfo {
		user.RoleList, _ = this.GetUserRoleList(user.Id)
		return user, nil
	}
	if user.Id != "" && !getRoleInfo {
		return user, nil
	}
	return user, errors.New("获取失败")
}

// 根据用户名获取用户信息
func (this *userService) GetUserByName(userName string) (*entity.User, error) {
	user := new(entity.User)
	GetByQ(Users, bson.M{"user_name": userName}, user)
	if user.Id != "" {
		return user, nil
	}
	return user, errors.New("获取失败")
}

// 获取用户总数
func (this *userService) GetTotal() (int64, error) {
	return int64(Count(Users, nil)), nil
}

// 分页获取用户列表
func (this *userService) GetUserList(page, pageSize int, getRoleInfo bool) ([]entity.User, error) {
	var users []entity.User
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "_id", false)
	err := Users.
		Find(nil).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&users)
	for k, user := range users {
		users[k].RoleList, _ = this.GetUserRoleList(user.Id)
	}
	return users, err
}

// 根据角色id获取用户列表
func (this *userService) GetUserListByRoleId(roleId string) ([]entity.User, error) {
	var users []entity.User
	var userRole []entity.UserRole
	ListByQ(UserRoles, bson.M{"role_id": roleId}, &userRole)
	if len(userRole) == 0 {
		return users, errors.New("角色不存在")
	}
	for _, v := range userRole {
		var user entity.User
		GetByQ(Users, bson.M{"_id": v.UserId}, &user)
		if user.Id != "" {
			users = append(users, user)
		}
	}
	return users, nil
}

// 获取某个用户的角色列表
// 为什么不直接连表查询role表？因为不想“越权”查询
func (this *userService) GetUserRoleList(userId string) ([]entity.Role, error) {
	var (
		roleRef  []entity.UserRole
		roleList []entity.Role
	)
	ListByQ(UserRoles, bson.M{"user_id": userId}, &roleRef)
	roleList = make([]entity.Role, 0, len(roleRef))
	for _, v := range roleRef {
		if role, err := RoleService.GetRole(v.RoleId); err == nil {
			roleList = append(roleList, *role)
		}
	}
	return roleList, nil
}

// 添加用户
func (this *userService) AddUser(userName, email, inviter, agent, phone, weixin, password string, sex int) (*entity.User, error) {
	if exists, _ := this.GetUserByName(userName); exists.Id != "" {
		return nil, errors.New("用户名已存在")
	}

	user := new(entity.User)
	user.UserName = userName
	user.Sex = sex
	user.Email = email
	user.Inviter = inviter
	user.Agent = agent
	user.Phone = phone
	user.Weixin = weixin
	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = libs.Md5([]byte(password + user.Salt))
	user.CreateTime = bson.Now()
	user.UpdateTime = bson.Now()
	user.LastLogin = bson.Now()
	// user.LastLogin = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	var err error
	user.Id, err = this.GetUserIDGen()
	if err != nil {
		return user, errors.New("添加失败")
	}
	if Insert(Users, user) {
		return user, nil
	}
	return user, errors.New("添加失败")
}

// 更新用户信息
func (this *userService) UpdateUser(user *entity.User, fileds bson.M) error {
	if len(fileds) < 1 {
		return errors.New("更新字段不能为空")
	}
	if Update(Users, bson.M{"_id": user.Id}, bson.M{"$set": fileds}) {
		return nil
	}
	return errors.New("更新失败")
}

// 修改密码
func (this *userService) ModifyPassword(userId string, password string) error {
	user, err := this.GetUser(userId, false)
	if err != nil {
		return err
	}
	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = libs.Md5([]byte(password + user.Salt))
	if Update(Users, bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"salt": user.Salt, "password": user.Password, "update_time": bson.Now()}}) {
		return nil
	}
	return errors.New("更新失败")
}

// 删除用户
func (this *userService) DeleteUser(userId string) error {
	if userId == "1" {
		return errors.New("不允许删除用户ID为1的用户")
	}
	if Delete(Users, bson.M{"_id": userId}) {
		return nil
	}
	return errors.New("删除用户失败")
}

// 设置用户角色
func (this *userService) UpdateUserRoles(userId string, roleIds []string) error {
	if _, err := this.GetUser(userId, false); err != nil {
		return err
	}
	DeleteAll(UserRoles, bson.M{"user_id": userId})
	for _, v := range roleIds {
		Insert(UserRoles, &entity.UserRole{Id: userId + "." + v, UserId: userId, RoleId: v})
	}
	return nil
}

//代理

func (this *userService) GetUserByPhone(phone string) (*entity.User, error) {
	user := new(entity.User)
	GetByQ(Users, bson.M{"phone": phone}, user)
	if user.Phone == "" {
		return user, errors.New("代理商不存在")
	}
	return user, nil
}

func (this *userService) GetUserByAgent(agent string) (*entity.User, error) {
	user := new(entity.User)
	GetByQ(Users, bson.M{"agent": agent}, user)
	if user.Agent == "" {
		return user, errors.New("代理商不存在")
	}
	return user, nil
}

// 添加代理商用户
func (this *userService) AddAgencyUser(userName, password, agent, phone, weixin, qq, address string) (*entity.User, error) {
	//if exists, _ := this.GetUserByName(userName); exists.Id != "" {
	//	return nil, errors.New("用户名已存在")
	//}

	user := new(entity.User)
	user.UserName = userName
	// user.Salt = salt
	// user.Password = password
	user.CreateTime = bson.Now()
	user.UpdateTime = bson.Now()
	user.LastLogin = bson.Now()
	// user.LastLogin = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = libs.Md5([]byte(password + user.Salt))
	user.Phone = phone
	user.Agent = agent
	user.Weixin = weixin
	user.QQ = qq
	user.Address = address
	user.Status = -1 //需要审核
	var err error
	user.Id, err = this.GetUserIDGen()
	if err != nil {
		return user, errors.New("添加失败")
	}
	if Insert(Users, user) {
		return user, nil
	}
	return user, errors.New("添加失败")
}

// 添加邀请
func (this *userService) InviterAddUser(userName, email, inviter, agent, phone, weixin, password string, sex int) (*entity.User, error) {
	if exists, _ := this.GetUserByName(userName); exists.Id != "" {
		return nil, errors.New("用户名已存在")
	}

	user := new(entity.User)
	user.UserName = userName
	user.Sex = sex
	user.Email = email
	user.Inviter = inviter
	user.Agent = agent
	user.Phone = phone
	user.Weixin = weixin
	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = libs.Md5([]byte(password + user.Salt))
	user.CreateTime = bson.Now()
	user.UpdateTime = bson.Now()
	user.LastLogin = bson.Now()
	// user.LastLogin = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	user.Status = -1 //需要审核
	var err error
	user.Id, err = this.GetUserIDGen()
	if err != nil {
		return user, errors.New("添加失败")
	}
	if Insert(Users, user) {
		return user, nil
	}
	return user, errors.New("添加失败")
}
