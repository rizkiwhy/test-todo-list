package user

import (
	"rizkiwhy/test-todo-list/package/user/model"
	"rizkiwhy/test-todo-list/util/database"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gorm.io/gorm"
)

type Repository interface {
	IsExistsByEmail(email string) (bool, error)
	Create(user model.User) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

type RepositoryImpl struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{DB: db}
}

func (r *RepositoryImpl) countByFilter(filter database.MySQLFilter) (res int64, err error) {
	db := database.BuildMySQLFilter(r.DB, filter)
	err = db.Model(&model.User{}).Count(&res).Error
	if err != nil {
		log.Error().Err(err).Interface("filter", filter).Msg("[UserRepository][countByFilter] Failed to count user by filter")
	}

	return
}

func (r *RepositoryImpl) getByFilter(filter database.MySQLFilter) (user *model.User, err error) {
	err = database.BuildMySQLFilter(r.DB, filter).First(&user).Error
	if err != nil {
		log.Error().Err(err).Interface("filter", filter).Msg("[UserRepository][getByFilter] Failed to get user by filter")
	}

	return
}

func (r *RepositoryImpl) IsExistsByEmail(email string) (res bool, err error) {
	filter := database.MySQLFilter{Where: gin.H{"email": email}}
	totalUsers, err := r.countByFilter(filter)
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("[UserRepository][IsExistsByEmail] Failed to count user by email")
		return
	}

	return totalUsers > 0, nil
}

func (r *RepositoryImpl) Create(user model.User) (*model.User, error) {
	result := r.DB.Create(&user)
	if result.Error != nil {
		log.Error().Err(result.Error).Interface("user", user).Msg("[UserRepository][Create] Failed to create user")
		return nil, result.Error
	}

	return &user, nil
}

func (r *RepositoryImpl) GetByEmail(email string) (*model.User, error) {
	user, err := r.getByFilter(database.MySQLFilter{Where: gin.H{"email": email}})
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("[UserRepository][GetByEmail] Failed to get user by email")
	}

	return user, nil

}
