package db

import (
	"blockexchange/types"
	"database/sql"
	"fmt"
	"strings"

	"github.com/minetest-go/dbutil"
)

func NewSchemaSearchRepository(DB *sql.DB) *SchemaSearchRepository {
	return &SchemaSearchRepository{
		DB:  DB,
		dbu: dbutil.New(DB, dbutil.DialectPostgres, func() *types.SchemaSearchResult { return &types.SchemaSearchResult{} }),
	}
}

type SchemaSearchRepository struct {
	DB  *sql.DB
	dbu *dbutil.DBUtil[*types.SchemaSearchResult]
}

func (r *SchemaSearchRepository) buildWhereQuery(query *strings.Builder, search *types.SchemaSearchRequest) []any {
	params := []any{}

	// complete flagW
	if search.Complete != nil {
		query.WriteString(" and s.complete = %s")
		params = append(params, *search.Complete)
	}

	if search.Keywords != nil {
		query.WriteString(" and s.search_tokens @@ to_tsquery(%s)")
		params = append(params, *search.Keywords)
	}

	if search.SchemaName != nil {
		query.WriteString(" and s.name = %s")
		params = append(params, *search.SchemaName)
	}

	if search.UserName != nil {
		query.WriteString(" and s.user_id = (select id from public.user where name = %s)")
		params = append(params, *search.UserName)
	}

	if search.SchemaID != nil {
		query.WriteString(" and s.id = %s")
		params = append(params, *search.SchemaID)
	}

	if search.SchemaIDList != nil && len(search.SchemaIDList) > 0 {
		p2 := make([]string, len(search.SchemaIDList))
		for i, id := range search.SchemaIDList {
			p2[i] = fmt.Sprintf("%d", id)
		}
		params = append(params, fmt.Sprintf("{%s}", strings.Join(p2, ",")))
		query.WriteString(" and s.id = any(%s::int[])")
	}

	if search.UserID != nil {
		query.WriteString(" and s.user_id = %s")
		params = append(params, *search.UserID)
	}

	if search.TagID != nil {
		query.WriteString(" and s.id in (select schema_id from schematag where tag_id = %s)")
		params = append(params, *search.TagID)
	}

	return params
}

func (r *SchemaSearchRepository) buildOrderQuery(query *strings.Builder, search *types.SchemaSearchRequest) {
	if search.OrderColumn != nil && search.OrderDirection != nil && types.OrderColumns[*search.OrderColumn] && types.OrderDirections[*search.OrderDirection] {
		query.WriteString(fmt.Sprintf(" order by %s %s", *search.OrderColumn, *search.OrderColumn))
	} else {
		query.WriteString(" order by s.mtime desc")
	}
}

func (r *SchemaSearchRepository) Count(search *types.SchemaSearchRequest) (int, error) {
	query := strings.Builder{}
	query.WriteString("where true=true")

	// build query
	params := r.buildWhereQuery(&query, search)
	return r.dbu.Count(query.String(), params...)
}

func (r *SchemaSearchRepository) Search(search *types.SchemaSearchRequest, limit, offset int) ([]*types.SchemaSearchResult, error) {

	query := strings.Builder{}
	query.WriteString("where true=true")

	// build query
	params := r.buildWhereQuery(&query, search)
	r.buildOrderQuery(&query, search)

	// add limit and offset
	query.WriteString(" limit %s offset %s")
	params = append(params, limit, offset)

	return r.dbu.SelectMulti(query.String(), params...)
}
