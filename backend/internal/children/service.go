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

func (s *ChildService) GetByID(ctx context.Context, id string) (*ChildByIdResponse, error) {
	child, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if child == nil {
		return nil, nil
	}

	alerts, err := s.repo.ListAlertsByChildID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, a := range alerts {
		msg := ""
		if a.Message.Valid {
			msg = a.Message.String
		}
		switch a.Category {
		case "health":
			child.Health.Alerts = append(child.Health.Alerts, msg)
		case "education":
			child.Education.Alerts = append(child.Education.Alerts, msg)
		case "social_assistance":
			child.SocialAssistance.Alerts = append(child.SocialAssistance.Alerts, msg)
		}
	}

	if child.Health.Alerts == nil {
		child.Health.Alerts = []string{}
	}
	if child.Education.Alerts == nil {
		child.Education.Alerts = []string{}
	}
	if child.SocialAssistance.Alerts == nil {
		child.SocialAssistance.Alerts = []string{}
	}

	notes := ""
	if child.Notes.Valid {
		notes = child.Notes.String
	}

	alertCategories := child.AlertCategories
	if alertCategories == nil {
		alertCategories = []string{}
	}

	return &ChildByIdResponse{
		ChildResponse: ChildResponse{
			ID:               child.ID,
			Name:             child.Name,
			Age:              child.Age,
			Neighborhood:     child.Neighborhood,
			AlertCategories:  alertCategories,
			Reviewed:         child.Reviewed,
			ReviewedBy:       child.ReviewedBy,
			ReviewedAt:       child.ReviewedAt,
			Notes:            notes,
			CreatedAt:        child.CreatedAt,
		},
		Health:           child.Health,
		Education:        child.Education,
		SocialAssistance: child.SocialAssistance,
	}, nil
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
