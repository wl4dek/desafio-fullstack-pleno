package statistics

import (
	"context"
)

type mockStatisticRepository struct {
	summaryFunc                func(ctx context.Context) (*Summary, error)
	statisticsByNeighborhoodFunc func(ctx context.Context) ([]NeighborhoodAlertCount, error)
}

func (m *mockStatisticRepository) Summary(ctx context.Context) (*Summary, error) {
	return m.summaryFunc(ctx)
}

func (m *mockStatisticRepository) StatisticsByNeighborhood(ctx context.Context) ([]NeighborhoodAlertCount, error) {
	return m.statisticsByNeighborhoodFunc(ctx)
}
