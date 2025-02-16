-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    coins INT DEFAULT 0

);
CREATE INDEX IF NOT EXISTS idx_username ON employees(username);

CREATE TABLE IF NOT EXISTS merch (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_merch_name ON merch(name);

CREATE TABLE IF NOT EXISTS purchases (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL,
    merch_id UUID NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (merch_id) REFERENCES merch(id)
);
CREATE INDEX IF NOT EXISTS idx_purchase_employee ON purchases(employee_id);
CREATE INDEX IF NOT EXISTS idx_purchase_merch ON purchases(merch_id);

CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY,
    from_emp_id UUID NOT NULL,
    to_emp_id UUID NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_emp_id) REFERENCES employees(id),
    FOREIGN KEY (to_emp_id) REFERENCES employees(id)
);
CREATE INDEX IF NOT EXISTS idx_transfer_from_emp ON transfers(from_emp_id);
CREATE INDEX IF NOT EXISTS idx_transfer_to_emp ON transfers(to_emp_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transfers;
DROP TABLE purchases;
DROP TABLE merch;
DROP TABLE employees;


DROP INDEX idx_username;
DROP INDEX idx_merch_name;
DROP INDEX idx_purchase_employee;
DROP INDEX idx_purchase_merch;
DROP INDEX idx_transfer_from_emp;
DROP INDEX idx_transfer_to_emp;
-- +goose StatementEnd
