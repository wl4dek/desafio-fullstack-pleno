package children

import (
	"context"
	"math"
)

type ChildService struct {
	repo ChildRepository
}

func NewChildService(repo ChildRepository) *ChildService {
	return &ChildService{repo: repo}
}

func (s *ChildService) List(ctx context.Context, filters Filters) (*PaginatedResponse, error) {
	filters.Normalize()

	children, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountFiltered(ctx, filters)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(filters.PerPage)))

	return &PaginatedResponse{
		Data: children,
		Pagination: Pagination{
			Page:       filters.Page,
			PerPage:    filters.PerPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *ChildService) GetByID(ctx context.Context, id string) (*Child, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ChildService) GetAreasByChildID(ctx context.Context, id string) (*Areas, error) {
	return s.repo.FindAreasByChildID(ctx, id)
}

func (s *ChildService) MarkReviewed(ctx context.Context, id string, reviewedBy string) error {
	return s.repo.MarkReviewed(ctx, id, reviewedBy)
}

func (s *ChildService) Summary(ctx context.Context) (*Summary, error) {
	return s.repo.Summary(ctx)
}

func (s *ChildService) ListNeighborhood(ctx context.Context) ([]string, error) {
	return s.repo.ListNeighborhood(ctx)
}
