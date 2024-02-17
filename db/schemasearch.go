package db

import (
	"blockexchange/types"
	"context"
	"fmt"
	"strings"

	"github.com/vingarcia/ksql"
)

type SchemaSearchRepository struct {
	kdb ksql.Provider
}

func (r *SchemaSearchRepository) buildWhereQuery(query *strings.Builder, search *types.SchemaSearchRequest, with_order bool) []any {
	query.WriteString("from schema where true=true")
	params := []any{}
	i := 1

	// complete flag
	if search.Complete != nil {
		query.WriteString(fmt.Sprintf(" and complete = $%d", i))
		params = append(params, *search.Complete)
		i++
	}

	if search.Keywords != nil {
		query.WriteString(fmt.Sprintf(" and to_tsvector('english', description || ' ' || name) @@ to_tsquery($%d)", i))
		params = append(params, *search.Keywords)
		i++
	}

	if search.SchemaName != nil {
		query.WriteString(fmt.Sprintf(" and name = $%d", i))
		params = append(params, *search.SchemaName)
		i++
	}

	if search.UserName != nil {
		query.WriteString(fmt.Sprintf(" and user_uid = (select uid from public.user where name = $%d)", i))
		params = append(params, *search.UserName)
		i++
	}

	if search.SchemaUID != nil {
		query.WriteString(fmt.Sprintf(" and uid = $%d", i))
		params = append(params, *search.SchemaUID)
		i++
	}

	if search.UserUID != nil {
		query.WriteString(fmt.Sprintf(" and user_uid = $%d", i))
		params = append(params, *search.UserUID)
		i++
	}

	if search.TagUID != nil {
		query.WriteString(fmt.Sprintf(" and uid in (select schema_uid from schematag where tag_uid = $%d)", i))
		params = append(params, *search.TagUID)
		i++
	}

	if with_order {
		if search.OrderColumn != nil && search.OrderDirection != nil && types.OrderColumns[*search.OrderColumn] && types.OrderDirections[*search.OrderDirection] {
			query.WriteString(fmt.Sprintf(" order by $%d $%d", i, i+1))
			params = append(params, *search.OrderColumn, *search.OrderColumn)
			i += 2
		} else {
			query.WriteString(" order by mtime desc")
		}
	}

	if search.Limit != nil {
		query.WriteString(fmt.Sprintf(" limit $%d", i))
		params = append(params, *search.Limit)
		i++
	}

	if search.Offset != nil {
		query.WriteString(fmt.Sprintf(" offset $%d", i))
		params = append(params, *search.Offset)
		i++
	}

	return params
}

func (r *SchemaSearchRepository) Count(search *types.SchemaSearchRequest) (int64, error) {
	query := strings.Builder{}
	params := r.buildWhereQuery(&query, search, false)
	c := &types.Count{}
	return c.Count, r.kdb.QueryOne(context.Background(), c, fmt.Sprintf("select count(*) as count %s", query.String()), params...)
}

func (r *SchemaSearchRepository) Search(search *types.SchemaSearchRequest) ([]*types.Schema, error) {
	query := strings.Builder{}
	params := r.buildWhereQuery(&query, search, true)
	list := []*types.Schema{}
	return list, r.kdb.Query(context.Background(), &list, query.String(), params...)
}
