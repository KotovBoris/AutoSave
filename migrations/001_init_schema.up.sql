-- 001_init_schema.up.sql
-- Initial database schema for AutoSave

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avg_salary DECIMAL(15,2),
    avg_expenses DECIMAL(15,2),
    savings_capacity DECIMAL(15,2),
    salary_dates INTEGER[],
    autopilot_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Banks table
CREATE TABLE IF NOT EXISTS banks (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    api_base_url VARCHAR(255) NOT NULL,
    deposit_rate DECIMAL(5,2) DEFAULT 8.0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User bank connections
CREATE TABLE IF NOT EXISTS user_banks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bank_id VARCHAR(50) NOT NULL REFERENCES banks(id),
    external_client_id VARCHAR(255) NOT NULL,
    bank_token TEXT,
    token_expires_at TIMESTAMP WITH TIME ZONE,
    account_consent_id VARCHAR(255),
    product_consent_id VARCHAR(255),
    payment_consent_id VARCHAR(255),
    connected BOOLEAN DEFAULT true,
    connected_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_sync_at TIMESTAMP WITH TIME ZONE,
    error TEXT,
    UNIQUE(user_id, bank_id)
);

-- Accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_bank_id INTEGER NOT NULL REFERENCES user_banks(id) ON DELETE CASCADE,
    bank_id VARCHAR(50) NOT NULL REFERENCES banks(id),
    external_id VARCHAR(255) NOT NULL,
    identification VARCHAR(255) NOT NULL,
    scheme_name VARCHAR(50),
    account_type VARCHAR(50),
    nickname VARCHAR(255),
    balance DECIMAL(15,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'RUB',
    servicer_name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_bank_id, external_id)
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    external_id VARCHAR(255) NOT NULL,
    booking_date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    value_date_time TIMESTAMP WITH TIME ZONE,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    description TEXT,
    credit_debit_indicator VARCHAR(10),
    counterparty_name VARCHAR(255),
    counterparty_account VARCHAR(255),
    category VARCHAR(50),
    is_salary BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(account_id, external_id)
);

-- Goals table
CREATE TABLE IF NOT EXISTS goals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    target_amount DECIMAL(15,2) NOT NULL,
    current_amount DECIMAL(15,2) DEFAULT 0,
    monthly_amount DECIMAL(15,2) NOT NULL,
    bank_id VARCHAR(50) NOT NULL REFERENCES banks(id),
    deposit_rate DECIMAL(5,2) NOT NULL,
    position INTEGER NOT NULL DEFAULT 1,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'waiting', 'completed', 'cancelled')),
    next_deposit_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Deposits table
CREATE TABLE IF NOT EXISTS deposits (
    id SERIAL PRIMARY KEY,
    goal_id INTEGER NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bank_id VARCHAR(50) NOT NULL REFERENCES banks(id),
    product_id VARCHAR(255),
    agreement_id VARCHAR(255),
    amount DECIMAL(15,2) NOT NULL,
    rate DECIMAL(5,2) NOT NULL,
    term_months INTEGER DEFAULT 12,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('pending', 'active', 'matured', 'closed', 'failed')),
    opened_at TIMESTAMP WITH TIME ZONE,
    matures_at TIMESTAMP WITH TIME ZONE,
    closed_at TIMESTAMP WITH TIME ZONE,
    accrued_interest DECIMAL(15,2) DEFAULT 0,
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Loans table
CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    original_debt DECIMAL(15,2) NOT NULL,
    current_debt DECIMAL(15,2) NOT NULL,
    rate DECIMAL(5,2) NOT NULL,
    monthly_payment DECIMAL(15,2) NOT NULL,
    autopay_enabled BOOLEAN DEFAULT false,
    autopay_bank_id VARCHAR(50) REFERENCES banks(id),
    autopay_day INTEGER CHECK (autopay_day >= 1 AND autopay_day <= 31),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'paid_off', 'cancelled')),
    next_payment_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    paid_off_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Loan payments table
CREATE TABLE IF NOT EXISTS loan_payments (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL,
    is_autopay BOOLEAN DEFAULT false,
    bank_payment_id VARCHAR(255),
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'processing', 'completed', 'failed')),
    scheduled_date DATE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Operations table (audit log)
CREATE TABLE IF NOT EXISTS operations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2),
    related_goal_id INTEGER REFERENCES goals(id),
    related_loan_id INTEGER REFERENCES loans(id),
    related_deposit_id INTEGER REFERENCES deposits(id),
    status VARCHAR(20) DEFAULT 'success' CHECK (status IN ('success', 'failed', 'pending')),
    error TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_banks_user_id ON user_banks(user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_user_bank_id ON accounts(user_bank_id);
CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_booking_date ON transactions(booking_date_time);
CREATE INDEX IF NOT EXISTS idx_transactions_is_salary ON transactions(is_salary);
CREATE INDEX IF NOT EXISTS idx_goals_user_id ON goals(user_id);
CREATE INDEX IF NOT EXISTS idx_goals_status ON goals(status);
CREATE INDEX IF NOT EXISTS idx_deposits_goal_id ON deposits(goal_id);
CREATE INDEX IF NOT EXISTS idx_deposits_user_id ON deposits(user_id);
CREATE INDEX IF NOT EXISTS idx_loans_user_id ON loans(user_id);
CREATE INDEX IF NOT EXISTS idx_loan_payments_loan_id ON loan_payments(loan_id);
CREATE INDEX IF NOT EXISTS idx_operations_user_id ON operations(user_id);
CREATE INDEX IF NOT EXISTS idx_operations_created_at ON operations(created_at);

-- Updated at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply updated_at triggers
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_goals_updated_at BEFORE UPDATE ON goals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_deposits_updated_at BEFORE UPDATE ON deposits
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_loans_updated_at BEFORE UPDATE ON loans
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

