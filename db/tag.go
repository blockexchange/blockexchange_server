package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TagRepository struct {
	DB *sqlx.DB
}

func (repo TagRepository) Create(tag *types.Tag) error {
	logrus.Trace("db.CreateTag", tag)
	query := `
		insert into
		tag(name, description)
		values(:name, :description)
		returning id
	`
	stmt, err := repo.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&tag.ID, tag)
}

func (repo TagRepository) Delete(id int64) error {
	_, err := repo.DB.Exec("delete from tag where id = $1", id)
	return err
}

func (repo TagRepository) Update(tag *types.Tag) error {
	query := `
	update tag
	set
		name = :name,
		description = :description
	where id = :id
`
	_, err := repo.DB.NamedExec(query, tag)
	return err
}

func (repo TagRepository) GetByID(id int64) (*types.Tag, error) {
	tags := types.Tag{}
	err := repo.DB.Get(&tags, "select * from tag where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &tags, nil
	}
}

func (repo TagRepository) GetAll() ([]*types.Tag, error) {
	list := []*types.Tag{}
	err := repo.DB.Select(&list, "select * from tag")
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}
