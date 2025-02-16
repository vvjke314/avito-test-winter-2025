package tests

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vvjke314/avito-test-winter-2025/internal/ds"
	"github.com/vvjke314/avito-test-winter-2025/internal/repository"
)

func TestCreateEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := &repository.Repo{DB: db}

	t.Run("success", func(t *testing.T) {
		employee := &ds.Employee{
			Id:           uuid.New(),
			Username:     "newuser",
			PasswordHash: "newhashedpassword",
			Coins:        50,
		}

		mock.ExpectQuery("INSERT INTO employees \\(id, username, password_hash, coins\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
			WithArgs(employee.Id, employee.Username, employee.PasswordHash, employee.Coins).
			WillReturnRows(sqlmock.NewRows([]string{}))

		err := r.CreateEmployee(employee)
		require.NoError(t, err)
	})

	t.Run("query error", func(t *testing.T) {
		employee := &ds.Employee{
			Id:           uuid.New(),
			Username:     "erroruser",
			PasswordHash: "errorhashedpassword",
			Coins:        0,
		}

		mock.ExpectQuery("INSERT INTO employees \\(id, username, password_hash, coins\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
			WithArgs(employee.Id, employee.Username, employee.PasswordHash, employee.Coins).
			WillReturnError(errors.New("insert failed"))

		err := r.CreateEmployee(employee)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "[QueryRow.Err] error occured")
	})
}

func TestFindEmployeeByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := &repository.Repo{DB: db}

	t.Run("success", func(t *testing.T) {
		expectedEmployee := &ds.Employee{
			Id:           uuid.New(),
			Username:     "testuser",
			PasswordHash: "hashedpassword",
			Coins:        100,
		}

		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "coins"}).
			AddRow(expectedEmployee.Id, expectedEmployee.Username, expectedEmployee.PasswordHash, expectedEmployee.Coins)

		mock.ExpectQuery("SELECT id, username, password_hash, coins FROM employees WHERE username = \\$1").
			WithArgs("testuser").
			WillReturnRows(rows)

		employee, err := r.FindEmployeeByUsername("testuser")
		require.NoError(t, err)
		assert.Equal(t, expectedEmployee, employee)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password_hash, coins FROM employees WHERE username = \\$1").
			WithArgs("unknownuser").
			WillReturnError(sql.ErrNoRows)

		employee, err := r.FindEmployeeByUsername("unknownuser")
		require.NoError(t, err)
		assert.Nil(t, employee)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password_hash, coins FROM employees WHERE username = \\$1").
			WithArgs("erroruser").
			WillReturnError(errors.New("query failed"))

		employee, err := r.FindEmployeeByUsername("erroruser")
		require.Error(t, err)
		assert.Nil(t, employee)
		assert.Contains(t, err.Error(), "[QueryRow.Scan] error occured")
	})
}

func TestTransferCoins(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository.Repo{DB: db}

	t.Run("insufficient_funds", func(t *testing.T) {
		sender := ds.Employee{
			Id:           uuid.New(),
			Username:     "sender",
			PasswordHash: "somehash",
			Coins:        10, // У отправителя недостаточно монет
		}
		amount := 50

		// Ожидание запроса для получения отправителя
		mock.ExpectQuery(`SELECT id, username, password_hash, coins FROM employees WHERE username = \$1`).
			WithArgs(sender.Username).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "username", "password_hash", "coins"}).
					AddRow(sender.Id, sender.Username, sender.PasswordHash, sender.Coins),
			)

		// Вызов метода
		err := repo.TransferCoins(sender.Username, "recipient", amount)

		// Проверка ошибки
		require.Error(t, err)
		require.Contains(t, err.Error(), "insufficient funds")

		// Проверка всех ожиданий
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetEmployeeInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository.Repo{DB: db}
	t.Run("success", func(t *testing.T) {
		// Настроить тестовые данные
		employee := ds.Employee{Id: uuid.New(), Username: "buyer", Coins: 100}
		item := ds.Merch{Id: uuid.New(), Name: "item1", Price: 50}

		// Ожидания запросов
		mock.ExpectQuery(`SELECT id, username, password_hash, coins FROM employees WHERE username = \$1`).
			WithArgs(employee.Username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins"}).
				AddRow(employee.Id, employee.Username, "hash", employee.Coins))

		mock.ExpectQuery(`SELECT merch.name, SUM\(purchases.amount\) AS quantity FROM purchases`).
			WithArgs(employee.Id).
			WillReturnRows(sqlmock.NewRows([]string{"name", "quantity"}).
				AddRow(item.Name, 2))

		mock.ExpectQuery(`SELECT employees.username, transfers.amount FROM transfers`).
			WithArgs(employee.Id).
			WillReturnRows(sqlmock.NewRows([]string{"username", "amount"}).
				AddRow("sender", 50))

		mock.ExpectQuery(`SELECT employees.username, transfers.amount FROM transfers`).
			WithArgs(employee.Id).
			WillReturnRows(sqlmock.NewRows([]string{"username", "amount"}).
				AddRow("receiver", 50))

		// Вызов метода
		info, err := repo.GetClientInfo(employee.Username)

		// Проверка результата
		require.NoError(t, err)
		require.NotNil(t, info)
		require.Equal(t, employee.Coins, info.Coins)
		require.Len(t, info.Inventory, 1)
		require.Len(t, info.CoinHistory.Received, 1)
		require.Len(t, info.CoinHistory.Sent, 1)
	})
}

func TestPurchase(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository.Repo{DB: db}
	t.Run("success", func(t *testing.T) {
		// Настроить тестовые данные
		employee := ds.Employee{Id: uuid.New(), Username: "buyer", Coins: 100}
		item := ds.Merch{Id: uuid.New(), Name: "item1", Price: 50}

		// Ожидание запросов
		mock.ExpectQuery(`SELECT id, username, password_hash, coins FROM employees WHERE username = \$1`).
			WithArgs(employee.Username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins"}).
				AddRow(employee.Id, employee.Username, "hash", employee.Coins))

		mock.ExpectQuery(`SELECT id, name, price FROM items WHERE name = \$1`).
			WithArgs(item.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow(item.Id, item.Name, item.Price))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE employees SET coins = coins - \$1 WHERE id = \$2`).
			WithArgs(item.Price, employee.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`INSERT INTO purchases \(id, emp_id, item_id, amount\) VALUES \(\$1, \$2, \$3, \$4\)`).
			WithArgs(sqlmock.AnyArg(), employee.Id, item.Id, item.Price).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Вызов метода
		err := repo.BuyItem(employee.Username, item.Name)

		// Проверка результата
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}
