package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/vvjke314/avito-test-winter-2025/internal/ds"
)

type Repo struct {
	DB *sql.DB
}

// FindEmployeeByUsername поиск записи сотрудника в БД по его никнейму
func (r *Repo) FindEmployeeByUsername(username string) (*ds.Employee, error) {
	query := `SELECT id, username, password_hash, coins 
	          FROM employees 
	          WHERE username = $1`

	row := r.DB.QueryRow(query, username)

	employee := &ds.Employee{}

	err := row.Scan(&employee.Id, &employee.Username, &employee.PasswordHash, &employee.Coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[QueryRow.Scan] error occured : %w", err)
	}

	return employee, nil
}

// CreateEmployee создание записи о сотруднике в БД
func (r *Repo) CreateEmployee(employee *ds.Employee) error {
	query := `INSERT INTO employees (id, username, password_hash, coins) 
	VALUES ($1, $2, $3, $4)`
	err := r.DB.QueryRow(
		query,
		employee.Id,
		employee.Username,
		employee.PasswordHash,
		employee.Coins).Err()
	if err != nil {
		return fmt.Errorf("[QueryRow.Err] error occured : %w", err)
	}

	return nil
}
