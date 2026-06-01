package children

import (
	"context"
)

type mockChildRepository struct {
	listFunc              func(ctx context.Context, filters Filters) ([]Child, error)
	countFilteredFunc     func(ctx context.Context, filters Filters) (int, error)
	findByIDFunc          func(ctx context.Context, id string) (*ChildById, error)
	markReviewedFunc      func(ctx context.Context, id string, reviewedBy string) error
	countFunc             func(ctx context.Context) (int, error)
	listAlertsByChildIDFunc func(ctx context.Context, id string) ([]Alerts, error)
	listNeighborhoodFunc  func(ctx context.Context) ([]string, error)
}

func (m *mockChildRepository) List(ctx context.Context, filters Filters) ([]Child, error) {
	return m.listFunc(ctx, filters)
}

func (m *mockChildRepository) CountFiltered(ctx context.Context, filters Filters) (int, error) {
	return m.countFilteredFunc(ctx, filters)
}

func (m *mockChildRepository) FindByID(ctx context.Context, id string) (*ChildById, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockChildRepository) MarkReviewed(ctx context.Context, id string, reviewedBy string) error {
	return m.markReviewedFunc(ctx, id, reviewedBy)
}

func (m *mockChildRepository) Count(ctx context.Context) (int, error) {
	return m.countFunc(ctx)
}

func (m *mockChildRepository) ListAlertsByChildID(ctx context.Context, id string) ([]Alerts, error) {
	return m.listAlertsByChildIDFunc(ctx, id)
}

func (m *mockChildRepository) ListNeighborhood(ctx context.Context) ([]string, error) {
	return m.listNeighborhoodFunc(ctx)
}
