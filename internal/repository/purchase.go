package repository

import (
	"database/sql"
	"errors"
	"fmt"

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
	_, err = tx.Exec(`INSERT INTO purchases (id, emp_id, item_id, amount) VALUES ($1, $2, $3, $4)`,
		uuid.New(), employee.Id, item.Id, item.Price)
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
