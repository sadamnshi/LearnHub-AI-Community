package repositories

import (
	"gin_demo/databases"
	"gin_demo/models"
)

// 定义一个接口处理用户相关的数据库操作
type UserRepository interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	UpdatePassword(userID uint, newHashedPassword string) error
}

// 定义一个结构体实现处理数据库数据的逻辑
type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

// 实现gorm插入一条记录到数据库的方法，user为传入的数据信息结构体
func (r *userRepository) Create(user *models.User) error {
	return databases.DB.Create(user).Error
}

// 实现传入的username参数查询数据库中是否存在对应的记录，如果存在则返回该记录的结构体指针，否则返回错误
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var u models.User
	if err := databases.DB.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByID 根据用户 ID 查询完整用户信息
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var u models.User
	if err := databases.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) UpdatePassword(userID uint, newHashedPassword string) error {
	return databases.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", newHashedPassword).Error
}
