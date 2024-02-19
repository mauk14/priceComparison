package usecase

import (
	"golang.org/x/net/context"
	"priceComp/services/SearchManager/internal/domain"
	"priceComp/services/SearchManager/internal/repository"
)

type SearchManager interface {
	ListProducts(context.Context, string, string, string, []string, int, int, int, int) ([]*domain.Products, error)
}

type searchManager struct {
	searchRep repository.SearchRepository
}

func New(searchRep repository.SearchRepository) SearchManager {
	return &searchManager{searchRep: searchRep}
}

func (s *searchManager) ListProducts(ctx context.Context, searchKeyword, category, sortBy string, brand []string, limit, offset, priceFrom, priceTo int) ([]*domain.Products, error) {
	products, err := s.searchRep.GetAll(ctx, searchKeyword, category, sortBy, brand, limit, offset, priceFrom, priceTo)
	if err != nil {
		return nil, err
	}
	return products, nil
}
