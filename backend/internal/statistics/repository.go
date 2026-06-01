package statistics

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticRepository interface {
	Summary(ctx context.Context) (*Summary, error)
	StatisticsByNeighborhood(ctx context.Context) ([]NeighborhoodAlertCount, error)
}

type statisticRepository struct {
	pool *pgxpool.Pool
}

func NewStatisticsRepository(pool *pgxpool.Pool) StatisticRepository {
	return &statisticRepository{pool: pool}
}

func (r *statisticRepository) Summary(ctx context.Context) (*Summary, error) {
	var total, reviewed, education, health, social int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*), COALESCE(SUM(CASE WHEN reviewed THEN 1 ELSE 0 END), 0) FROM children").
		Scan(&total, &reviewed)
	if err != nil {
		return nil, fmt.Errorf("failed to get summary: %w", err)
	}

	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM alert_education").Scan(&education)
	if err != nil {
		return nil, fmt.Errorf("failed to get education alerts: %w", err)
	}

	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM alert_health").Scan(&health)
	if err != nil {
		return nil, fmt.Errorf("failed to get health alerts: %w", err)
	}

	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM alert_social_assistance").Scan(&social)
	if err != nil {
		return nil, fmt.Errorf("failed to get social assistance alerts: %w", err)
	}

	alertsByArea := make(map[string]int)
	alertsByArea["education"] = education
	alertsByArea["health"] = health
	alertsByArea["social_assistance"] = social

	return &Summary{
		TotalChildren: total,
		Reviewed:      reviewed,
		PendingReview: total - reviewed,
		AlertsByArea:  alertsByArea,
	}, nil
}

func (r *statisticRepository) StatisticsByNeighborhood(ctx context.Context) ([]NeighborhoodAlertCount, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			c.neighborhood,
			COUNT(DISTINCT ah.id) AS health,
			COUNT(DISTINCT ae.id) AS education,
			COUNT(DISTINCT asa.id) AS social_assistance
		FROM children c
		LEFT JOIN health h ON c.id = h.child_id
		LEFT JOIN alert_health ah ON h.id = ah.health_id
		LEFT JOIN education e ON c.id = e.child_id
		LEFT JOIN alert_education ae ON e.id = ae.education_id
		LEFT JOIN social_assistance s ON c.id = s.child_id
		LEFT JOIN alert_social_assistance asa ON s.id = asa.social_assistance_id
		GROUP BY c.neighborhood
		ORDER BY c.neighborhood`)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics by neighborhood: %w", err)
	}
	defer rows.Close()

	var items []NeighborhoodAlertCount
	for rows.Next() {
		var item NeighborhoodAlertCount
		err := rows.Scan(
			&item.Neighborhood,
			&item.Health,
			&item.Education,
			&item.SocialAssistance,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan neighborhood alert count: %w", err)
		}
		items = append(items, item)
	}

	if items == nil {
		items = []NeighborhoodAlertCount{}
	}

	return items, nil
}
