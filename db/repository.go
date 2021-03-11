package db

type Processor func()

type Repository struct {
	table string
}

func (r Repository) GetBy(wherepart string, target interface{}) {

}

func NewRepository(table string) Repository {
	return Repository{
		table: table,
	}
}
