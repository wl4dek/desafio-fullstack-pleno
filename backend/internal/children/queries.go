package children

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	baseQuery  string
	countQuery string
	where      []string
	args       []interface{}
	argIdx     int
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		baseQuery:  "SELECT c.id, c.name, c.age, c.neighborhood, STRING_AGG(distinct a.category, ', ') AS alert_categories, c.reviewed, c.reviewed_by, c.reviewed_at, notes, c.created_at FROM children c LEFT JOIN alert a ON c.id = a.child_id",
		countQuery: "SELECT COUNT(*) FROM children c LEFT JOIN alert a ON c.id = a.child_id",
	}
}

func (qb *QueryBuilder) AddCondition(condition string, arg interface{}) {
	qb.where = append(qb.where, condition)
	qb.args = append(qb.args, arg)
}

func (qb *QueryBuilder) BuildList() (string, []interface{}) {
	query := qb.baseQuery
	if len(qb.where) > 0 {
		query += " WHERE " + strings.Join(qb.where, " AND ")
	}
	query += " GROUP BY c.id ORDER BY c.created_at DESC"
	return query, qb.args
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
