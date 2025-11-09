package repository

import (
    "context"
    "time"
    
    "github.com/KotovBoris/AutoSave/backend/internal/models"
    "github.com/jmoiron/sqlx"
)

type Repositories struct {
    User        UserRepository
    Bank        BankRepository
    Account     AccountRepository
    Transaction TransactionRepository
    Goal        GoalRepository
    Deposit     DepositRepository
    Loan        LoanRepository
    Operation   OperationRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
    return &Repositories{
        User:        NewUserRepository(db),
        Bank:        NewBankRepository(db),
        Account:     NewAccountRepository(db),
        Transaction: NewTransactionRepository(db),
        Goal:        NewGoalRepository(db),
        Deposit:     NewDepositRepository(db),
        Loan:        NewLoanRepository(db),
        Operation:   NewOperationRepository(db),
    }
}

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id int) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    UpdateFinancialProfile(ctx context.Context, userID int, avgSalary, avgExpenses, savingsCapacity float64, salaryDates []int) error
    UpdateAutopilot(ctx context.Context, userID int, enabled bool) error
}

type BankRepository interface {
    GetAll(ctx context.Context) ([]models.Bank, error)
    GetByID(ctx context.Context, id string) (*models.Bank, error)
    CreateConnection(ctx context.Context, conn *models.BankConnection) error
    GetUserConnections(ctx context.Context, userID int) ([]models.BankConnection, error)
    GetConnection(ctx context.Context, userID int, bankID string) (*models.BankConnection, error)
    GetConnectionByID(ctx context.Context, id int) (*models.BankConnection, error)
    UpdateConnection(ctx context.Context, conn *models.BankConnection) error
    DeleteConnection(ctx context.Context, userID int, bankID string) error
}

type AccountRepository interface {
    Create(ctx context.Context, account *models.Account) error
    CreateBatch(ctx context.Context, accounts []models.Account) error
    GetByID(ctx context.Context, id int) (*models.Account, error)
    GetByExternalID(ctx context.Context, userBankID int, externalID string) (*models.Account, error)
    GetUserAccounts(ctx context.Context, userID int) ([]models.Account, error)
    GetBankAccounts(ctx context.Context, userID int, bankID string) ([]models.Account, error)
    Update(ctx context.Context, account *models.Account) error
    UpdateBalance(ctx context.Context, id int, balance float64) error
    Delete(ctx context.Context, id int) error
}

type TransactionRepository interface {
    Create(ctx context.Context, tx *models.Transaction) error
    CreateBatch(ctx context.Context, transactions []models.Transaction) error
    GetByID(ctx context.Context, id int) (*models.Transaction, error)
    GetByExternalID(ctx context.Context, accountID int, externalID string) (*models.Transaction, error)
    GetAccountTransactions(ctx context.Context, filter models.TransactionFilter) ([]models.Transaction, error)
    GetUserTransactions(ctx context.Context, userID int, fromDate, toDate time.Time) ([]models.Transaction, error)
    GetSalaryTransactions(ctx context.Context, userID int) ([]models.Transaction, error)
    MarkAsSalary(ctx context.Context, transactionIDs []int) error
    CountAccountTransactions(ctx context.Context, accountID int) (int, error)
}

type GoalRepository interface {
    Create(ctx context.Context, goal *models.Goal) error
    GetByID(ctx context.Context, id int) (*models.Goal, error)
    GetUserGoals(ctx context.Context, userID int) ([]models.Goal, error)
    GetActiveGoal(ctx context.Context, userID int) (*models.Goal, error)
    Update(ctx context.Context, goal *models.Goal) error
    UpdatePosition(ctx context.Context, goalID int, position int) error
    UpdateStatus(ctx context.Context, goalID int, status string) error
    UpdateCurrentAmount(ctx context.Context, goalID int, amount float64) error
    Delete(ctx context.Context, id int) error
    GetMaxPosition(ctx context.Context, userID int) (int, error)
}

type DepositRepository interface {
    Create(ctx context.Context, deposit *models.Deposit) error
    GetByID(ctx context.Context, id int) (*models.Deposit, error)
    GetByAgreementID(ctx context.Context, agreementID string) (*models.Deposit, error)
    GetGoalDeposits(ctx context.Context, goalID int) ([]models.Deposit, error)
    GetUserDeposits(ctx context.Context, userID int) ([]models.Deposit, error)
    GetActiveDeposits(ctx context.Context, userID int) ([]models.Deposit, error)
    Update(ctx context.Context, deposit *models.Deposit) error
    UpdateStatus(ctx context.Context, id int, status string) error
    Close(ctx context.Context, id int, closedAt time.Time) error
}

type LoanRepository interface {
    Create(ctx context.Context, loan *models.Loan) error
    GetByID(ctx context.Context, id int) (*models.Loan, error)
    GetUserLoans(ctx context.Context, userID int) ([]models.Loan, error)
    GetActiveLoans(ctx context.Context, userID int) ([]models.Loan, error)
    Update(ctx context.Context, loan *models.Loan) error
    UpdateDebt(ctx context.Context, id int, currentDebt float64) error
    UpdateStatus(ctx context.Context, id int, status string) error
    Delete(ctx context.Context, id int) error
    CreatePayment(ctx context.Context, payment *models.LoanPayment) error
    GetPayments(ctx context.Context, loanID int) ([]models.LoanPayment, error)
    GetScheduledPayments(ctx context.Context, date time.Time) ([]models.LoanPayment, error)
    UpdatePayment(ctx context.Context, payment *models.LoanPayment) error
}

type OperationRepository interface {
    Create(ctx context.Context, operation *models.Operation) error
    GetByID(ctx context.Context, id int) (*models.Operation, error)
    GetUserOperations(ctx context.Context, userID int, limit int) ([]models.Operation, error)
    GetByType(ctx context.Context, userID int, operationType string) ([]models.Operation, error)
}

