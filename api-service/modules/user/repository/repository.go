package repository

import (
	"gorm.io/gorm"

	"stock-data-processing/api-service/entities"
)

type Repository interface {
	GetList(take, limit int) (users *[]entities.User, totalData int64, err error)
	InsertOne(user entities.User) error
	GetOne(id int) (entities.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return repository{
		db: db,
	}
}

func (r repository) GetList(take, limit int) (users *[]entities.User, totalData int64, err error) {
	res := r.db.Limit(limit).Take(take).Find(&users)
	if res.Error != nil {
		return nil, 0, err
	}

	r.db.Model(&entities.User{}).Group("id").Count(&totalData)
	return users, totalData, nil
}

func (r repository) InsertOne(user entities.User) error {
	crt := r.db.Create(&user)
	if crt.Error != nil {
		return crt.Error
	}

	return nil
}

func (r repository) GetOne(id int) (entities.User, error) {
	user := entities.User{}
	err := r.db.Where("id", "=", id).First(&user)
	if err != nil {
		return entities.User{}, nil
	}

	return user, nil
}
