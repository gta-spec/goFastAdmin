package UserService

import (
	"gota/app/common/dao/User"
	"gota/app/common/model"
)

func GetById(id uint) *model.User {
	return UserDao.GetById(id)
}
