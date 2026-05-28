package children

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrChildNotFound = errors.New("child not found")

type ChildRepository interface {
	List(ctx context.Context, filters Filters) ([]Child, error)
	CountFiltered(ctx context.Context, filters Filters) (int, error)
	FindByID(ctx context.Context, id string) (*Child, error)
	MarkReviewed(ctx context.Context, id string, reviewedBy string) error
	Summary(ctx context.Context) (*Summary, error)
	Count(ctx context.Context) (int, error)
	FindAreasByChildID(ctx context.Context, id string) (*Areas, error)
	ListNeighborhood(ctx context.Context) ([]string, error)
}

type childRepository struct {
	pool *pgxpool.Pool
}

func NewChildRepository(pool *pgxpool.Pool) ChildRepository {
	return &childRepository{pool: pool}
}

func (r *childRepository) ListNeighborhood(ctx context.Context) ([]string, error) {
	rows, err := r.pool.Query(ctx, "SELECT DISTINCT neighborhood FROM children ORDER BY neighborhood")
	if err != nil {
		return nil, fmt.Errorf("failed to list neighborhoods: %w", err)
	}
	defer rows.Close()

	var neighborhoods []string
	for rows.Next() {
		var n string
		err := rows.Scan(&n)
		if err != nil {
			return nil, fmt.Errorf("failed to scan neighborhood: %w", err)
		}
		neighborhoods = append(neighborhoods, n)
	}

	return neighborhoods, nil
}

func (r *childRepository) FindAreasByChildID(ctx context.Context, id string) (*Areas, error) {
	var areas Areas
	err := r.pool.QueryRow(ctx, ` 
		(SELECT school_name, alerts, frequency_percent FROM education WHERE child_id = $1)
	`, id).Scan(&areas.Education.SchoolName, &areas.Education.Alerts, &areas.Education.FrequenciaPercent)
	if err != nil {
		return nil, fmt.Errorf("failed to find areas (education) by child ID: %w", err)
	}

	err = r.pool.QueryRow(ctx, ` 
		(SELECT vaccinations_up_to_date, alerts, last_consultation FROM health WHERE child_id = $1)
	`, id).Scan(&areas.Health.VaccinationsUpToDate, &areas.Health.Alerts, &areas.Health.LastConsultation)
	if err != nil {
		return nil, fmt.Errorf("failed to find areas (health) by child ID: %w", err)
	}

	err = r.pool.QueryRow(ctx, ` 
		(SELECT cad_unico, active_benefit, alerts FROM social_assistance WHERE child_id = $1)
	`, id).Scan(&areas.SocialAssistance.CadUnico, &areas.SocialAssistance.ActiveBenefit, &areas.SocialAssistance.Alerts)
	if err != nil {
		return nil, fmt.Errorf("failed to find areas (social_assistance) by child ID: %w", err)
	}
	return &areas, nil
}

func (r *childRepository) List(ctx context.Context, filters Filters) ([]Child, error) {
	qb := NewQueryBuilder()

	if filters.Name != "" {
		qb.AddCondition("name ILIKE $1", "%"+filters.Name+"%")
	}
	if filters.Neighborhood != "" {
		qb.AddCondition("neighborhood = $1", filters.Neighborhood)
	}
	if filters.HasAlert != nil {
		qb.AddCondition(fmt.Sprintf("has_alert = $%d", len(qb.args)+1), *filters.HasAlert)
	}
	if filters.Reviewed != nil {
		qb.AddCondition(fmt.Sprintf("reviewed = $%d", len(qb.args)+1), *filters.Reviewed)
	}

	query, args := qb.BuildPaginatedList(filters.PerPage, filters.Offset())
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list children: %w", err)
	}
	defer rows.Close()

	var children []Child
	for rows.Next() {
		var c Child
		err := rows.Scan(
			&c.ID, &c.Name, &c.Age, &c.Neighborhood,
			&c.HasAlert, &c.Reviewed, &c.ReviewedBy, &c.ReviewedAt,
			&c.Notes, &c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan child: %w", err)
		}
		children = append(children, c)
	}

	if children == nil {
		children = []Child{}
	}

	return children, nil
}

func (r *childRepository) CountFiltered(ctx context.Context, filters Filters) (int, error) {
	qb := NewQueryBuilder()

	if filters.Neighborhood != "" {
		qb.AddCondition("neighborhood = $1", filters.Neighborhood)
	}
	if filters.HasAlert != nil {
		qb.AddCondition(fmt.Sprintf("has_alert = $%d", len(qb.args)+1), *filters.HasAlert)
	}
	if filters.Reviewed != nil {
		qb.AddCondition(fmt.Sprintf("reviewed = $%d", len(qb.args)+1), *filters.Reviewed)
	}

	query, args := qb.BuildCount()
	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count filtered children: %w", err)
	}
	return count, nil
}

func (r *childRepository) FindByID(ctx context.Context, id string) (*Child, error) {
	var c Child
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, age, neighborhood, has_alert, reviewed, reviewed_by, reviewed_at, COALESCE(notes, ''), created_at
		FROM children WHERE id like $1`, id).
		Scan(
			&c.ID, &c.Name, &c.Age, &c.Neighborhood,
			&c.HasAlert, &c.Reviewed, &c.ReviewedBy, &c.ReviewedAt,
			&c.Notes, &c.CreatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find child: %w", err)
	}

	return &c, nil
}

func (r *childRepository) MarkReviewed(ctx context.Context, id string, reviewedBy string) error {
	now := time.Now().UTC()
	result, err := r.pool.Exec(ctx, `
		UPDATE children SET reviewed = true, reviewed_by = $1, reviewed_at = $2
		WHERE id = $3 AND reviewed = false`, reviewedBy, now, id)
	if err != nil {
		return fmt.Errorf("failed to mark reviewed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrChildNotFound
	}

	return nil
}

func (r *childRepository) Summary(ctx context.Context) (*Summary, error) {
	var total, reviewed, education, health, social int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*), COALESCE(SUM(CASE WHEN reviewed THEN 1 ELSE 0 END), 0) FROM children").
		Scan(&total, &reviewed)
	if err != nil {
		return nil, fmt.Errorf("failed to get summary: %w", err)
	}

	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM education").Scan(&education)
	if err != nil {
		return nil, fmt.Errorf("failed to get education alerts: %w", err)
	}

	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM health").Scan(&health)
	if err != nil {
		return nil, fmt.Errorf("failed to get health alerts: %w", err)
	}

	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM social_assistance").Scan(&social)
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

func (r *childRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM children").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count children: %w", err)
	}
	return count, nil
}
