package db

import (
	"blockexchange/types"
	"database/sql"
	"fmt"
	"strings"
)

type SchemaSearchRepository struct {
	DB *sql.DB
}

func (repo SchemaSearchRepository) buildWhereQuery(query *strings.Builder, search *types.SchemaSearchRequest) []any {
	params := []any{}
	bind_index := 1

	// complete flagW
	if search.Complete != nil {
		query.WriteString(fmt.Sprintf(" and s.complete = $%d", bind_index))
		bind_index++
		params = append(params, *search.Complete)
	}

	if search.Keywords != nil {
		query.WriteString(fmt.Sprintf(" and s.search_tokens @@ to_tsquery($%d)", bind_index))
		bind_index++
		params = append(params, *search.Keywords)
	}

	if search.SchemaName != nil {
		query.WriteString(fmt.Sprintf(" and s.name = $%d", bind_index))
		bind_index++
		params = append(params, *search.SchemaName)
	}

	if search.UserName != nil {
		query.WriteString(fmt.Sprintf(" and s.user_id = (select id from public.user where name = $%d)", bind_index))
		bind_index++
		params = append(params, *search.UserName)
	}

	if search.SchemaID != nil {
		query.WriteString(fmt.Sprintf(" and s.id = $%d", bind_index))
		bind_index++
		params = append(params, *search.SchemaID)
	}

	if search.UserID != nil {
		query.WriteString(fmt.Sprintf(" and s.user_id = $%d", bind_index))
		bind_index++
		params = append(params, *search.UserID)
	}

	if search.TagID != nil {
		query.WriteString(fmt.Sprintf(" and s.id in (select schema_id from schematag where tag_id = $%d)", bind_index))
		bind_index++
		params = append(params, *search.TagID)
	}

	return params
}

func (repo SchemaSearchRepository) buildOrderQuery(query *strings.Builder, search *types.SchemaSearchRequest) {
	if search.OrderColumn != nil && search.OrderDirection != nil {
		query.WriteString(fmt.Sprintf(" order by %s %s", *search.OrderColumn, *search.OrderDirection))
	} else {
		query.WriteString(" order by s.created desc")
	}
}

func (repo SchemaSearchRepository) Count(search *types.SchemaSearchRequest) (int, error) {
	query := strings.Builder{}
	query.WriteString("where true=true")

	// build query
	params := repo.buildWhereQuery(&query, search)
	return Count(repo.DB, &types.SchemaSearchResult{}, query.String(), params...)
}

func (repo SchemaSearchRepository) Search(search *types.SchemaSearchRequest, limit, offset int) ([]*types.SchemaSearchResult, error) {

	query := strings.Builder{}
	query.WriteString("where true=true")

	// build query
	params := repo.buildWhereQuery(&query, search)
	repo.buildOrderQuery(&query, search)

	// add limit and offset
	query.WriteString(fmt.Sprintf(" limit $%d offset $%d", len(params)+1, len(params)+2))
	params = append(params, limit, offset)

	return SelectMulti(repo.DB, func() *types.SchemaSearchResult { return &types.SchemaSearchResult{} }, query.String(), params...)
}
