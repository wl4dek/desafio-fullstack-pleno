package summary

import (
	"context"

	"backend/internal/children"
)

type SummaryService struct {
	childRepo children.ChildRepository
}

func NewSummaryService(childRepo children.ChildRepository) *SummaryService {
	return &SummaryService{childRepo: childRepo}
}

func (s *SummaryService) GetSummary(ctx context.Context) (*children.Summary, error) {
	return s.childRepo.Summary(ctx)
}
