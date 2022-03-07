package db

import (
	"blockexchange/types"
	"fmt"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type SchemaSearchRepository interface {
	Search(search *types.SchemaSearchRequest, limit, offset int) ([]*types.SchemaSearchResult, error)
}

func NewSchemaSearchRepository(db *sqlx.DB) SchemaSearchRepository {
	return &DBSchemaSearchRepository{
		DB:             db,
		UserRepo:       DBUserRepository{DB: db},
		SchemaModRepo:  DBSchemaModRepository{DB: db},
		SchemaTagRepo:  DBSchemaTagRepository{DB: db},
		SchemaStarRepo: DBSchemaStarRepository{DB: db},
	}
}

type DBSchemaSearchRepository struct {
	DB             *sqlx.DB
	UserRepo       UserRepository
	SchemaModRepo  SchemaModRepository
	SchemaTagRepo  SchemaTagRepository
	SchemaStarRepo SchemaStarRepository
}

var schemaSearchHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "bx_schema_search_hist",
	Help:    "Histogram for the schema render time",
	Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
})

func (repo DBSchemaSearchRepository) Search(search *types.SchemaSearchRequest, limit, offset int) ([]*types.SchemaSearchResult, error) {
	timer := prometheus.NewTimer(schemaSearchHistogram)
	defer timer.ObserveDuration()

	query := strings.Builder{}
	query.WriteString("select * from schema where true=true")
	params := []interface{}{}
	bind_index := 1

	// complete flag
	if search.Complete != nil {
		query.WriteString(fmt.Sprintf(" and complete = $%d", bind_index))
		bind_index++
		params = append(params, *search.Complete)
	}

	if search.Keywords != nil {
		query.WriteString(fmt.Sprintf(" and search_tokens @@ to_tsquery($%d)", bind_index))
		bind_index++
		params = append(params, *search.Keywords)
	}

	if search.SchemaName != nil {
		query.WriteString(fmt.Sprintf(" and name = $%d", bind_index))
		bind_index++
		params = append(params, *search.SchemaName)
	}

	if search.UserName != nil {
		query.WriteString(fmt.Sprintf(" and user_id = (select id from public.user where name = $%d)", bind_index))
		bind_index++
		params = append(params, *search.UserName)
	}

	if search.SchemaID != nil {
		query.WriteString(fmt.Sprintf(" and id = $%d", bind_index))
		bind_index++
		params = append(params, *search.SchemaID)
	}

	if search.UserID != nil {
		query.WriteString(fmt.Sprintf(" and user_id = $%d", bind_index))
		bind_index++
		params = append(params, *search.UserID)
	}

	if search.UserID != nil {
		query.WriteString(fmt.Sprintf(" and id in (select schema_id from schematag where tag_id = $%d)", bind_index))
		bind_index++
		params = append(params, *search.TagID)
	}

	if search.OrderColumn != nil && search.OrderDirection != nil {
		query.WriteString(fmt.Sprintf(" order by %s %s", *search.OrderColumn, *search.OrderDirection))
	} else {
		query.WriteString(" order by created desc")
	}

	// add limit and offset
	query.WriteString(fmt.Sprintf(" limit $%d offset $%d", bind_index, bind_index+1))
	bind_index += 2
	params = append(params, limit, offset)

	logrus.WithFields(logrus.Fields{
		"query":  query.String(),
		"params": params,
	}).Trace("DBSchemaSearchRepository::Search")

	list := []*types.Schema{}
	err := repo.DB.Select(&list, query.String(), params...)
	if err != nil {
		return nil, err
	}

	return repo.enhance(list)
}

func (repo DBSchemaSearchRepository) enhance_schema(schema *types.Schema) (*types.SchemaSearchResult, error) {
	user, err := repo.UserRepo.GetUserById(schema.UserID)
	if err != nil {
		return nil, err
	}

	stars, err := repo.SchemaStarRepo.CountBySchemaID(schema.ID)
	if err != nil {
		return nil, err
	}

	mods, err := repo.SchemaModRepo.GetSchemaModsBySchemaID(schema.ID)
	if err != nil {
		return nil, err
	}

	mod_list := make([]string, len(mods))
	for i, mod := range mods {
		mod_list[i] = mod.ModName
	}

	tag_list, err := repo.SchemaTagRepo.GetBySchemaID(schema.ID)
	if err != nil {
		return nil, err
	}

	result := &types.SchemaSearchResult{
		Schema: *schema,
		Stars:  stars,
		User: &types.User{
			Name:    user.Name,
			ID:      user.ID,
			Created: user.Created,
			Type:    user.Type,
		},
		Mods: mod_list,
		Tags: tag_list,
	}

	return result, nil
}

func (repo DBSchemaSearchRepository) enhance(list []*types.Schema) ([]*types.SchemaSearchResult, error) {
	result := make([]*types.SchemaSearchResult, len(list))
	var wg sync.WaitGroup
	err_chan := make(chan error)

	for i, schema := range list {
		wg.Add(1)
		go func(schema *types.Schema, index int, wg *sync.WaitGroup, err_chan chan error) {
			enhanced_schema, err := repo.enhance_schema(schema)
			if err != nil {
				err_chan <- err
			} else {
				result[index] = enhanced_schema
			}
			wg.Done()
		}(schema, i, &wg, err_chan)
	}

	wg.Wait()

	select {
	case err := <-err_chan:
		return nil, err
	default:
		return result, nil
	}
}
