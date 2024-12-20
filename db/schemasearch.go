package db

import (
	"blockexchange/types"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/lib/pq"
)

type SchemaSearchRepository struct {
	DB *sql.DB
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func (r *SchemaSearchRepository) buildWhereQuery(query *strings.Builder, search *types.SchemaSearchRequest, with_order bool) []any {
	query.WriteString(" from schema s where true=true")
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
		// remove non-alpha/non-numbers
		keywords := nonAlphanumericRegex.ReplaceAllString(*search.Keywords, "")
		// remove double spaces
		keywords = strings.ReplaceAll(keywords, "  ", " ")
		params = append(params, strings.Join(strings.Split(strings.TrimSpace(keywords), " "), "&"))
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

	if search.ModName != nil {
		query.WriteString(fmt.Sprintf(" and uid in (select schema_uid from schemamod where mod_name = $%d)", i))
		params = append(params, *search.ModName)
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

	if search.TagName != nil {
		query.WriteString(fmt.Sprintf(" and uid in (select schema_uid from schematag where tag_uid = (select uid from tag where name = $%d))", i))
		params = append(params, *search.TagName)
		i++
	}

	if search.CollectionUID != nil {
		query.WriteString(fmt.Sprintf(" and collection_uid = $%d", i))
		params = append(params, *search.CollectionUID)
		i++
	}

	if search.CollectionName != nil {
		query.WriteString(fmt.Sprintf(" and collection_uid in (select c.uid from collection c where c.name = $%d)", i))
		params = append(params, *search.CollectionName)
		i++
	}

	if search.FromMtime != nil {
		query.WriteString(fmt.Sprintf(" and mtime > $%d", i))
		params = append(params, *search.FromMtime)
		i++
	}

	if search.UntilMtime != nil {
		query.WriteString(fmt.Sprintf(" and mtime < $%d", i))
		params = append(params, *search.UntilMtime)
		i++
	}

	if search.WithCollection != nil {
		if *search.WithCollection {
			// only schematics with collection
			query.WriteString(" and collection_uid is not null")
		} else {
			// without collection
			query.WriteString(" and collection_uid is null")
		}
	}

	if with_order {
		if search.OrderColumn != nil && search.OrderDirection != nil && types.OrderColumns[*search.OrderColumn] && types.OrderDirections[*search.OrderDirection] {
			query.WriteString(fmt.Sprintf(" order by $%d %s", i, *search.OrderDirection))
			params = append(params, *search.OrderColumn)
			i++
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
	query.WriteString("select count(*) as count")
	params := r.buildWhereQuery(&query, search, false)
	var c int64
	rows, err := r.DB.Query(query.String(), params...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	return c, rows.Scan(&c)
}

func (r *SchemaSearchRepository) Search(search *types.SchemaSearchRequest) ([]*types.SchemaSearchResponse, error) {
	query := strings.Builder{}
	fields := []string{
		"uid",
		"created",
		"mtime",
		"user_uid",
		"collection_uid",
		"name",
		"description",
		"short_description",
		"cdb_collection",
		"complete",
		"size_x", "size_y", "size_z",
		"total_size",
		"total_parts",
		"downloads",
		"views",
		"license",
		"stars",
		"(select u.name from public.user u where u.uid = user_uid)",
		"(select c.name from collection c where c.uid = collection_uid)",
		"array(select sm.mod_name from schemamod sm where sm.schema_uid = uid)::text[]",
		"array(select t.name from schematag st join tag t on t.uid = st.tag_uid  where st.schema_uid = s.uid)",
	}
	query.WriteString(fmt.Sprintf("select %s", strings.Join(fields, ",")))
	params := r.buildWhereQuery(&query, search, true)
	list := []*types.SchemaSearchResponse{}

	rows, err := r.DB.Query(query.String(), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := &types.SchemaSearchResponse{
			Schema: &types.Schema{},
			Tags:   []string{},
			Mods:   []string{},
		}
		err = rows.Scan(
			&e.Schema.UID,
			&e.Schema.Created,
			&e.Schema.Mtime,
			&e.Schema.UserUID,
			&e.Schema.CollectionUID,
			&e.Schema.Name,
			&e.Schema.Description,
			&e.Schema.ShortDescription,
			&e.Schema.CDBCollection,
			&e.Schema.Complete,
			&e.Schema.SizeX, &e.Schema.SizeY, &e.Schema.SizeZ,
			&e.Schema.TotalSize,
			&e.Schema.TotalParts,
			&e.Schema.Downloads,
			&e.Schema.Views,
			&e.Schema.License,
			&e.Schema.Stars,
			&e.Username,
			&e.CollectionName,
			pq.Array(&e.Mods),
			pq.Array(&e.Tags),
		)

		if err != nil {
			return nil, err
		}

		list = append(list, e)
	}

	return list, nil
}
