package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SchemaScreenshotRepository interface {
	GetBySchemaID(schema_id int64) ([]types.SchemaScreenshot, error)
	GetByID(id int64) (*types.SchemaScreenshot, error)
	Create(screenshot *types.SchemaScreenshot) error
	Update(screenshot *types.SchemaScreenshot) error
	//Delete(id int64) error
}

type DBSchemaScreenshotRepository struct {
	DB *sqlx.DB
}

func (r DBSchemaScreenshotRepository) GetByID(id int64) (*types.SchemaScreenshot, error) {
	result := types.SchemaScreenshot{}
	err := r.DB.Get(&result, "select * from schema_screenshot where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

func (r DBSchemaScreenshotRepository) GetBySchemaID(schema_id int64) ([]types.SchemaScreenshot, error) {
	list := []types.SchemaScreenshot{}
	err := r.DB.Select(&list, "select * from schema_screenshot where schema_id = $1", schema_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (r DBSchemaScreenshotRepository) Create(screenshot *types.SchemaScreenshot) error {
	query := `
		insert into
		schema_screenshot(
			schema_id, type, title, data
		)
		values(
			:schema_id, :type, :title, :data
		)
		returning id
	`
	stmt, err := r.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&screenshot.ID, screenshot)
}

func (r DBSchemaScreenshotRepository) Update(screenshot *types.SchemaScreenshot) error {
	query := `
		update schema_screenshot
		set
			type = :type,
			title = :title,
			data = :data
		where id = :id
	`
	_, err := r.DB.NamedExec(query, screenshot)
	return err
}
