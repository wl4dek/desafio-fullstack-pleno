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
	FindByID(ctx context.Context, id string) (*ChildById, error)
	MarkReviewed(ctx context.Context, id string, reviewedBy string) error
	Summary(ctx context.Context) (*Summary, error)
	Count(ctx context.Context) (int, error)
	ListAlertsByChildID(ctx context.Context, id string) ([]Alerts, error)
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

func (r *childRepository) ListAlertsByChildID(ctx context.Context, id string) ([]Alerts, error) {
	var alerts []Alerts

	rows, err := r.pool.Query(ctx, `
		SELECT 'health' AS category, code, message FROM alert_health ah
		JOIN health h ON h.id = ah.health_id WHERE h.child_id = $1
		UNION ALL
		SELECT 'education' AS category, code, message FROM alert_education ae
		JOIN education e ON e.id = ae.education_id WHERE e.child_id = $1
		UNION ALL
		SELECT 'social_assistance' AS category, code, message FROM alert_social_assistance asa
		JOIN social_assistance s ON s.id = asa.social_assistance_id WHERE s.child_id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find alerts by child ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a Alerts
		err := rows.Scan(&a.Category, &a.Code, &a.Message)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}
		alerts = append(alerts, a)
	}

	if alerts == nil {
		alerts = []Alerts{}
	}

	return alerts, nil
}

func (r *childRepository) List(ctx context.Context, filters Filters) ([]Child, error) {
	qb := NewQueryBuilder()

	if filters.Name != "" {
		qb.AddCondition(fmt.Sprintf("c.name ILIKE $%d", len(qb.args)+1), "%"+filters.Name+"%")
	}
	if filters.Neighborhood != "" {
		qb.AddCondition(fmt.Sprintf("c.neighborhood = $%d", len(qb.args)+1), filters.Neighborhood)
	}
	if filters.Alert != "" {
		n := len(qb.args) + 1
		qb.AddCondition(fmt.Sprintf(`(
			EXISTS(SELECT 1 FROM health h2 JOIN alert_health ah ON h2.id = ah.health_id WHERE h2.child_id = c.id AND ah.code = $%d)
			OR EXISTS(SELECT 1 FROM education e2 JOIN alert_education ae ON e2.id = ae.education_id WHERE e2.child_id = c.id AND ae.code = $%d)
			OR EXISTS(SELECT 1 FROM social_assistance s2 JOIN alert_social_assistance asa ON s2.id = asa.social_assistance_id WHERE s2.child_id = c.id AND asa.code = $%d)
		)`, n, n, n), filters.Alert)
	}
	if filters.Reviewed != nil {
		qb.AddCondition(fmt.Sprintf("c.reviewed = $%d", len(qb.args)+1), *filters.Reviewed)
	}
	if filters.HasAlert != nil && *filters.HasAlert {
		qb.AddConditionOnly(`(
			EXISTS(SELECT 1 FROM health h2 JOIN alert_health ah ON h2.id = ah.health_id WHERE h2.child_id = c.id)
			OR EXISTS(SELECT 1 FROM education e2 JOIN alert_education ae ON e2.id = ae.education_id WHERE e2.child_id = c.id)
			OR EXISTS(SELECT 1 FROM social_assistance s2 JOIN alert_social_assistance asa ON s2.id = asa.social_assistance_id WHERE s2.child_id = c.id)
		)`)
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
			&c.AlertCategories, &c.Reviewed, &c.ReviewedBy, &c.ReviewedAt,
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
		qb.AddCondition(fmt.Sprintf("c.neighborhood = $%d", len(qb.args)+1), filters.Neighborhood)
	}
	if filters.Alert != "" {
		n := len(qb.args) + 1
		qb.AddCondition(fmt.Sprintf(`(
			EXISTS(SELECT 1 FROM health h2 JOIN alert_health ah ON h2.id = ah.health_id WHERE h2.child_id = c.id AND ah.code = $%d)
			OR EXISTS(SELECT 1 FROM education e2 JOIN alert_education ae ON e2.id = ae.education_id WHERE e2.child_id = c.id AND ae.code = $%d)
			OR EXISTS(SELECT 1 FROM social_assistance s2 JOIN alert_social_assistance asa ON s2.id = asa.social_assistance_id WHERE s2.child_id = c.id AND asa.code = $%d)
		)`, n, n, n), filters.Alert)
	}
	if filters.Reviewed != nil {
		qb.AddCondition(fmt.Sprintf("c.reviewed = $%d", len(qb.args)+1), *filters.Reviewed)
	}
	if filters.HasAlert != nil && *filters.HasAlert {
		qb.AddConditionOnly(`(
			EXISTS(SELECT 1 FROM health h2 JOIN alert_health ah ON h2.id = ah.health_id WHERE h2.child_id = c.id)
			OR EXISTS(SELECT 1 FROM education e2 JOIN alert_education ae ON e2.id = ae.education_id WHERE e2.child_id = c.id)
			OR EXISTS(SELECT 1 FROM social_assistance s2 JOIN alert_social_assistance asa ON s2.id = asa.social_assistance_id WHERE s2.child_id = c.id)
		)`)
	}

	query, args := qb.BuildCount()
	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count filtered children: %w", err)
	}
	return count, nil
}

func (r *childRepository) FindByID(ctx context.Context, id string) (*ChildById, error) {
	var c ChildById
	qb := NewQueryBuilder()
	query, args := qb.BuildById(id)
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&c.ID, &c.Name, &c.Age, &c.Neighborhood, &c.AlertCategories,
		&c.Reviewed, &c.ReviewedBy, &c.ReviewedAt,
		&c.Notes, &c.CreatedAt, &c.Health.VaccinationsUpToDate, &c.Health.LastConsultation,
		&c.Education.SchoolName, &c.Education.FrequenciaPercent,
		&c.SocialAssistance.CadUnico, &c.SocialAssistance.ActiveBenefit,
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

func (r *childRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM children").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count children: %w", err)
	}
	return count, nil
}
