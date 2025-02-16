package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vvjke314/avito-test-winter-2025/internal/ds"
)

// GetClientInfo получение информации о пользователе (монеты, инвентарь, история)
func (r *Repo) GetClientInfo(username string) (*ds.InfoResponse, error) {
	// Получаем информацию о пользователе
	employee, err := r.FindEmployeeByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("[FindEmployeeByUsername] error occured: %w", err)
	}
	if employee == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Получаем баланс монет
	coins := employee.Coins

	// Получаем инвентарь пользователя
	inventory, err := r.getInventory(employee.Id)
	if err != nil {
		return nil, fmt.Errorf("[getInventory] error occured: %w", err)
	}

	// Получаем историю переводов
	coinHistory, err := r.getTransferHistory(employee.Id)
	if err != nil {
		return nil, fmt.Errorf("[getTransferHistory] error occured: %w", err)
	}

	// Формируем итоговый ответ
	infoResponse := &ds.InfoResponse{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}

	return infoResponse, nil
}

// getInventory получает инвентарь пользователя (купленные предметы)
func (r *Repo) getInventory(employeeId uuid.UUID) ([]ds.ItemAmount, error) {
	// Получаем все предметы, которые были куплены пользователем
	query := `SELECT merch.name, SUM(purchases.amount) AS quantity 
			  FROM purchases 
			  JOIN merch ON purchases.merch_id = merch.id 
			  WHERE purchases.emp_id = $1 
			  GROUP BY merch.name`

	rows, err := r.DB.Query(query, employeeId)
	if err != nil {
		return nil, fmt.Errorf("[Query] error occured: %w", err)
	}
	defer rows.Close()

	var inventory []ds.ItemAmount
	for rows.Next() {
		var item ds.ItemAmount
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, fmt.Errorf("[rows.Scan] error occured: %w", err)
		}
		inventory = append(inventory, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[rows.Err] error occured: %w", err)
	}

	return inventory, nil
}

// getTransferHistory получает историю переводов монет
func (r *Repo) getTransferHistory(employeeId uuid.UUID) (ds.TransferHistory, error) {
	// Получаем все переводы монет, которые были получены пользователем
	receivedQuery := `SELECT employees.username, transfers.amount 
					  FROM transfers 
					  JOIN employees ON transfers.from_emp_id = employees.id 
					  WHERE transfers.to_emp_id = $1`
	rows, err := r.DB.Query(receivedQuery, employeeId)
	if err != nil {
		return ds.TransferHistory{}, fmt.Errorf("[Query] error occured: %w", err)
	}
	defer rows.Close()

	var received []ds.ReceiveRecord
	for rows.Next() {
		var record ds.ReceiveRecord
		if err := rows.Scan(&record.FromUser, &record.Amount); err != nil {
			return ds.TransferHistory{}, fmt.Errorf("[rows.Scan] error occured: %w", err)
		}
		received = append(received, record)
	}

	// Получаем все переводы монет, которые были отправлены пользователем
	sentQuery := `SELECT employees.username, transfers.amount 
				  FROM transfers 
				  JOIN employees ON transfers.to_emp_id = employees.id 
				  WHERE transfers.from_emp_id = $1`
	rows, err = r.DB.Query(sentQuery, employeeId)
	if err != nil {
		return ds.TransferHistory{}, fmt.Errorf("[Query] error occured: %w", err)
	}
	defer rows.Close()

	var sent []ds.SentRecord
	for rows.Next() {
		var record ds.SentRecord
		if err := rows.Scan(&record.ToUser, &record.Amount); err != nil {
			return ds.TransferHistory{}, fmt.Errorf("[rows.Scan] error occured: %w", err)
		}
		sent = append(sent, record)
	}

	// Возвращаем историю переводов
	return ds.TransferHistory{
		Received: received,
		Sent:     sent,
	}, nil
}
