package filter

import (
	"errors"
	"forum/internal/models"
	"forum/internal/repository/filter"
)

type FilterService struct {
	repo filter.Filter
}

type Filter interface {
	FilterByCategories(categories []int) ([]models.Post, error)
	FilterByLikes(userId int) ([]models.Post, error)
	GetCategoryByName(strings []string) ([]int, error)
}

func NewFilterService(repository filter.Filter) *FilterService {
	return &FilterService{
		repo: repository,
	}
}

func (f *FilterService) FilterByCategories(categories []int) ([]models.Post, error) {
	return f.repo.GetPostsByCategories(categories)
}

func (f *FilterService) FilterByLikes(userID int) ([]models.Post, error) {
	return f.repo.GetUsersByLikedPosts(userID)
}

func (f *FilterService) GetCategoryByName(strings []string) ([]int, error) {
	if len(strings) == 0 {
		return []int{4}, nil
	}
	res := []int{}
	category := map[string]int{
		"Detective": 1,
		"Horror":    2,
		"Comedy":    3,
		"Other":     4,
	}
	for _, str := range strings {
		if num, ok := category[str]; !ok {
			return []int{}, errors.New("invalid Post")
		} else {
			res = append(res, num)
		}
	}
	return res, nil
}
