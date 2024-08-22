package db

import (
	"database/sql"
	"fmt"
	"time"
)

type DBLock struct {
	DB *sql.DB
}

func NewLock(DB *sql.DB) *DBLock {
	return &DBLock{DB: DB}
}

func (l *DBLock) TryLock(id int64) (bool, error) {
	rows, err := l.DB.Query("select pg_try_advisory_lock($1)", id)
	if err != nil {
		return false, fmt.Errorf("could not try-lock: %v", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return false, fmt.Errorf("no try-result returned")
	}

	result := false
	return result, rows.Scan(&result)
}

func (l *DBLock) UnLock(id int64) (bool, error) {
	rows, err := l.DB.Query("select pg_advisory_unlock($1)", id)
	if err != nil {
		return false, fmt.Errorf("could not un-lock: %v", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return false, fmt.Errorf("no unlock-result returned")
	}

	result := false
	return result, rows.Scan(&result)
}

type RunLockedFunc func() error

var ErrTryLockTimeout = fmt.Errorf("trylock timed out")

func (l *DBLock) RunLocked(id int64, timeout time.Duration, fn RunLockedFunc) error {
	start := time.Now()
	for {
		result, err := l.TryLock(id)
		if err != nil {
			return fmt.Errorf("trylock loop error: %v", err)
		}
		if result {
			break
		}
		time.Sleep(500 * time.Millisecond)
		if time.Since(start) > timeout {
			return ErrTryLockTimeout
		}
	}

	fnerr := fn()

	_, err := l.UnLock(id)
	if err != nil {
		return fmt.Errorf("unlock error: %v", err)
	}

	if fnerr != nil {
		return fmt.Errorf("user-function error: %v", err)
	}

	return nil
}
