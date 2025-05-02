-- Add account_id column to transactions
ALTER TABLE transactions ADD COLUMN account_id INTEGER REFERENCES accounts(id);

-- Create accounts table
CREATE TABLE accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    balance REAL NOT NULL DEFAULT 0,
    user_id INTEGER NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);