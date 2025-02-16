package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TransferCoins отправка монет от одного сотрудника другому
func (r *Repo) TransferCoins(fromUsername, toUsername string, amount int) error {
	// Проверяем отправителя
	fromEmployee, err := r.FindEmployeeByUsername(fromUsername)
	if err != nil {
		return fmt.Errorf("[FindEmployeeByUsername] error occured: %w", err)
	}
	if fromEmployee == nil {
		return fmt.Errorf("sender with username '%s' not found", fromUsername)
	}
	if fromEmployee.Coins < amount {
		return fmt.Errorf("insufficient funds: sender has %d coins, but tried to send %d", fromEmployee.Coins, amount)
	}

	// Проверяем получателя
	toEmployee, err := r.FindEmployeeByUsername(toUsername)
	if err != nil {
		return fmt.Errorf("[FindEmployeeByUsername] error occured: %w", err)
	}
	if toEmployee == nil {
		return fmt.Errorf("recipient with username '%s' not found", toUsername)
	}

	// Начинаем транзакцию после всех проверок
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("[DB.Begin] error occured: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Обновляем количество монет у отправителя
	_, err = tx.Exec(`UPDATE employees SET coins = coins - $1 WHERE id = $2`, amount, fromEmployee.Id)
	if err != nil {
		return fmt.Errorf("[Update Sender] error occured: %w", err)
	}

	// Обновляем количество монет у получателя
	_, err = tx.Exec(`UPDATE employees SET coins = coins + $1 WHERE id = $2`, amount, toEmployee.Id)
	if err != nil {
		return fmt.Errorf("[Update Recipient] error occured: %w", err)
	}

	// Создаем запись о переводе
	_, err = tx.Exec(
		`INSERT INTO transfers (id, from_emp_id, to_emp_id, amount, created_at) 
		 VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), fromEmployee.Id, toEmployee.Id, amount, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("[Insert Transfer] error occured: %w", err)
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("[Tx.Commit] error occured: %w", err)
	}

	return nil
}
