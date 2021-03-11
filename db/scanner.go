package db

type Scanner interface {
	Scan(dest ...interface{}) error
	Err() error
}
