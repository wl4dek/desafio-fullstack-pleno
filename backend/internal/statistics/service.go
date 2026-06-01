package statistics

import (
	"context"
)

type StatisticsService struct {
	statisticRepo StatisticRepository
}

func NewStatisticsService(statisticRepo StatisticRepository) *StatisticsService {
	return &StatisticsService{statisticRepo: statisticRepo}
}

func (s *StatisticsService) GetStatistics(ctx context.Context) (*StatisticsResponse, error) {
	items, err := s.statisticRepo.StatisticsByNeighborhood(ctx)
	if err != nil {
		return nil, err
	}

	return &StatisticsResponse{Statistics: items}, nil
}

func (s *StatisticsService) GetSummary(ctx context.Context) (*Summary, error) {
	items, err := s.statisticRepo.Summary(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}
