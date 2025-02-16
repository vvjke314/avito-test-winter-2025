package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vvjke314/avito-test-winter-2025/internal/ds"
)

// BuyItem покупка предмета
func (r *Repo) BuyItem(username string, itemName string) error {
	// Получаем информацию о пользователе
	employee, err := r.FindEmployeeByUsername(username)
	if err != nil {
		return fmt.Errorf("[FindEmployeeByUsername] error occured: %w", err)
	}
	if employee == nil {
		return fmt.Errorf("user not found")
	}

	// Получаем информацию о предмете
	var item ds.Merch
	query := `SELECT id, name, price FROM items WHERE name = $1`
	err = r.DB.QueryRow(query, itemName).Scan(&item.Id, &item.Name, &item.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("item not found")
		}
		return fmt.Errorf("[QueryRow.Scan] error occured: %w", err)
	}

	// Проверяем, есть ли у пользователя достаточно монет
	if employee.Coins < item.Price {
		return fmt.Errorf("insufficient funds")
	}

	// Начинаем транзакцию
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("[DB.Begin] error occured: %w", err)
	}

	// Обновляем баланс пользователя (уменьшаем монеты)
	_, err = tx.Exec(`UPDATE employees SET coins = coins - $1 WHERE id = $2`, item.Price, employee.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[Update Employee Balance] error occured: %w", err)
	}

	// Записываем покупку в таблицу purchases
	_, err = tx.Exec(`INSERT INTO purchases (id, employee_id, merch_id, amount) VALUES ($1, $2, $3, $4)`,
		uuid.New(), employee.Id, item.Id, 1)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[Insert Purchase Record] error occured: %w", err)
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("[DB.Commit] error occured: %w", err)
	}

	return nil
}

func (r *Repo) FindItemByName(name string) (*ds.Merch, error) {
	query := `SELECT id, name, price FROM merch WHERE name = $1`
	row := r.DB.QueryRow(query, name)

	var item ds.Merch
	if err := row.Scan(&item.Id, &item.Name, &item.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *Repo) DecreaseUserCoins(tx *sql.Tx, userId string, amount int) error {
	var currentCoins int
	query := `SELECT coins FROM employees WHERE id = $1`
	err := tx.QueryRow(query, userId).Scan(&currentCoins)
	if err != nil {
		return fmt.Errorf("[QueryRow.Scan] error occurred: %w", err)
	}

	if currentCoins < amount {
		return fmt.Errorf("insufficient funds: required %d, but have %d", amount, currentCoins)
	}

	updateQuery := `UPDATE employees SET coins = coins - $1 WHERE id = $2`
	_, err = tx.Exec(updateQuery, amount, userId)
	if err != nil {
		return fmt.Errorf("[Exec] error occurred while updating coins: %w", err)
	}

	return nil
}

func (r *Repo) RecordPurchase(tx *sql.Tx, userId string, itemId string, quantity int) error {
	insertQuery := `INSERT INTO purchases (id, employee_id, merch_id, amount, created_at) 
	                VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(insertQuery, uuid.New(), userId, itemId, quantity, time.Now())
	if err != nil {
		return fmt.Errorf("[Exec] error occurred while inserting purchase record: %w", err)
	}

	return nil
}

// IncreaseUserCoins добавляет монеты пользователю
func (r *Repo) IncreaseUserCoins(tx *sql.Tx, userId string, amount int) error {
	query := `UPDATE employees SET coins = coins + $1 WHERE id = $2`
	_, err := tx.Exec(query, amount, userId)
	if err != nil {
		return fmt.Errorf("[Exec] error occurred while adding coins: %w", err)
	}
	return nil
}

// RecordCoinTransfer записывает информацию о переводе монет в таблицу
func (r *Repo) RecordCoinTransfer(tx *sql.Tx, senderId string, receiverId string, amount int) error {
	query := `INSERT INTO transfers (id, from_emp_id, to_emp_id, amount, created_at) 
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(query, uuid.New(), senderId, receiverId, amount, time.Now())
	if err != nil {
		return fmt.Errorf("[Exec] error occurred while inserting coin transfer record: %w", err)
	}
	return nil
}
