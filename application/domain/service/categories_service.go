package service

import (
	inport "news-api/application/port/in"
	outport "news-api/application/port/out"

	"github.com/google/uuid"
)

type CategoriesService struct {
	categoriesPort outport.Categories
}

func NewCategoriesService(categoriesPort outport.Categories) *CategoriesService {
	return &CategoriesService{categoriesPort: categoriesPort}
}

func (g *CategoriesService) GetAll() ([]*inport.Category, error) {
	categoriesList, err := g.categoriesPort.GetAll()
	if err != nil {
		return nil, err
	}
	return func() []*inport.Category {
		result := make([]*inport.Category, len(categoriesList))
		for i, v := range categoriesList {
			result[i] = MapCategory(v)
		}
		return result
	}(), nil
}

func (g *CategoriesService) Search(keyword string) ([]*inport.Category, error) {
	categoriesList, err := g.categoriesPort.Search(keyword)
	if err != nil {
		return nil, err
	}
	return func() []*inport.Category {
		result := make([]*inport.Category, len(categoriesList))
		for i, v := range categoriesList {
			result[i] = MapCategory(v)
		}
		return result
	}(), nil
}

func (g *CategoriesService) Insert(category *inport.CreateCategoryPayload) error {

	return g.categoriesPort.Insert(outport.Category{
		ID:   uuid.New(),
		Name: category.Name,
	})
}

func (g *CategoriesService) Update(category *inport.UpdateCategoryPayload) error {
	return g.categoriesPort.Update(outport.Category{
		ID:   category.ID,
		Name: category.Name,
	})
}

func (g *CategoriesService) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}
