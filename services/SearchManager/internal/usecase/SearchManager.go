package usecase

import (
	"golang.org/x/net/context"
	"priceComp/services/SearchManager/internal/domain"
	"priceComp/services/SearchManager/internal/repository"
)

type SearchManager interface {
	ListProducts(context.Context, string, string, string, string, int, int) ([]*domain.Products, error)
}

type searchManager struct {
	searchRep repository.SearchRepository
}

func New(searchRep repository.SearchRepository) SearchManager {
	return &searchManager{searchRep: searchRep}
}

func (s *searchManager) ListProducts(ctx context.Context, searchKeyword, category, brand, sortBy string, limit, offset int) ([]*domain.Products, error) {
	products, err := s.searchRep.GetAll(ctx, searchKeyword, category, brand, sortBy, limit, offset)
	if err != nil {
		return nil, err
	}
	return products, nil
}
