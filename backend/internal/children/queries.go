package children

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	baseQuery     string
	countQuery    string
	findByIDQuery string
	where         []string
	args          []interface{}
	argIdx        int
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		baseQuery: `SELECT c.id, c.name, c.age, c.neighborhood,
  ARRAY_REMOVE(ARRAY[
    CASE WHEN EXISTS(SELECT 1 FROM health h2 JOIN alert_health ah ON h2.id = ah.health_id WHERE h2.child_id = c.id) THEN 'health' END,
    CASE WHEN EXISTS(SELECT 1 FROM education e2 JOIN alert_education ae ON e2.id = ae.education_id WHERE e2.child_id = c.id) THEN 'education' END,
    CASE WHEN EXISTS(SELECT 1 FROM social_assistance s2 JOIN alert_social_assistance asa ON s2.id = asa.social_assistance_id WHERE s2.child_id = c.id) THEN 'social_assistance' END
  ], NULL) AS alert_categories,
  c.reviewed, c.reviewed_by, c.reviewed_at, notes, c.created_at
FROM children c`,
		countQuery: "SELECT COUNT(*) FROM children c",
		findByIDQuery: `SELECT c.id, c.name, c.age, c.neighborhood,
  ARRAY_REMOVE(ARRAY[
    CASE WHEN EXISTS(SELECT 1 FROM health h2 JOIN alert_health ah ON h2.id = ah.health_id WHERE h2.child_id = c.id) THEN 'health' END,
    CASE WHEN EXISTS(SELECT 1 FROM education e2 JOIN alert_education ae ON e2.id = ae.education_id WHERE e2.child_id = c.id) THEN 'education' END,
    CASE WHEN EXISTS(SELECT 1 FROM social_assistance s2 JOIN alert_social_assistance asa ON s2.id = asa.social_assistance_id WHERE s2.child_id = c.id) THEN 'social_assistance' END
  ], NULL) AS alert_categories,
  c.reviewed, c.reviewed_by, c.reviewed_at, notes, c.created_at,
  COALESCE(h.vaccinations_up_to_date, false), h.last_consultation,
  e.school_name, COALESCE(e.frequency_percent, 0),
  COALESCE(s.cad_unico, false), COALESCE(s.active_benefit, false)
FROM children c
LEFT JOIN health h ON c.id = h.child_id
LEFT JOIN education e ON c.id = e.child_id
LEFT JOIN social_assistance s ON c.id = s.child_id
WHERE c.id = $1`,
	}
}

func (qb *QueryBuilder) AddCondition(condition string, arg interface{}) {
	qb.where = append(qb.where, condition)
	qb.args = append(qb.args, arg)
}

func (qb *QueryBuilder) AddConditionOnly(condition string) {
	qb.where = append(qb.where, condition)
}

func (qb *QueryBuilder) BuildList() (string, []interface{}) {
	query := qb.baseQuery
	if len(qb.where) > 0 {
		query += " WHERE " + strings.Join(qb.where, " AND ")
	}
	query += " ORDER BY c.created_at DESC"
	return query, qb.args
}

func (qb *QueryBuilder) BuildById(id string) (string, []interface{}) {
	query := qb.findByIDQuery
	args := []interface{}{id}
	return query, args
}

func (qb *QueryBuilder) BuildPaginatedList(perPage, offset int) (string, []interface{}) {
	query, args := qb.BuildList()
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", perPage, offset)
	return query, args
}

func (qb *QueryBuilder) BuildCount() (string, []interface{}) {
	query := qb.countQuery
	if len(qb.where) > 0 {
		query += " WHERE " + strings.Join(qb.where, " AND ")
	}
	return query, qb.args
}
