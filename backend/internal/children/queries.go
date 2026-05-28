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
		baseQuery:  "SELECT id, name, age, neighborhood, has_alert, reviewed, reviewed_by, reviewed_at, COALESCE(notes, ''), created_at FROM children",
		countQuery: "SELECT COUNT(*) FROM children",
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
	query += " ORDER BY created_at DESC"
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
