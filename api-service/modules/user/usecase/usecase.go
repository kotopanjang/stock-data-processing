package usecase

import (
	"stock-data-processing/api-service/entities"
	"stock-data-processing/api-service/modules/user/repository"
	"stock-data-processing/api-service/utils"
)

type Usecase interface {
	GetList(take, limit int) (*utils.Pagination, error)
	Insert(user entities.User) error
}

type usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) usecase {
	return usecase{
		repo: repo,
	}
}

func (u usecase) GetList(take, limit int) (*utils.Pagination, error) {
	res, count, err := u.repo.GetList(take, limit)
	if err != nil {
		return nil, err
	}

	return &utils.Pagination{
		Data:  res,
		Count: int(count),
		Error: nil,
	}, nil
}

func (u usecase) Insert(user entities.User) error {
	err := u.repo.InsertOne(user)
	if err != nil {
		return err
	}

	return nil
}
