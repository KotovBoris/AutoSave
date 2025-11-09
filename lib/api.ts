import {
  User,
  Bank,
  Account,
  Goal,
  Loan,
  Operation,
  SalaryTransaction,
  EmergencyWithdraw,
} from '@/types';
import {
  MOCK_USER,
  MOCK_BANKS,
  MOCK_ACCOUNTS,
  MOCK_GOALS,
  MOCK_LOANS,
  MOCK_OPERATIONS,
  MOCK_SALARY_TRANSACTIONS,
} from './mockData';

// Simulate API delay
const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

// Auth API
export const authAPI = {
  async register(email: string, password: string): Promise<User> {
    await delay(1000);
    if (email === 'error@example.com') {
      throw new Error('Email уже используется');
    }
    return MOCK_USER;
  },

  async login(email: string, password: string): Promise<User> {
    await delay(1000);
    return MOCK_USER;
  },

  async getCurrentUser(): Promise<User> {
    await delay(500);
    return MOCK_USER;
  },
};

// Banks API
export const banksAPI = {
  async getBanks(): Promise<Bank[]> {
    await delay(500);
    return MOCK_BANKS;
  },

  async connectBank(bankId: string): Promise<Bank> {
    await delay(2000);
    const bank = MOCK_BANKS.find((b) => b.id === bankId);
    if (!bank) throw new Error('Банк не найден');
    bank.connected = true;
    bank.connectedAt = new Date().toISOString().split('T')[0];
    return bank;
  },

  async disconnectBank(bankId: string): Promise<void> {
    await delay(500);
    const bank = MOCK_BANKS.find((b) => b.id === bankId);
    if (bank) bank.connected = false;
  },

  async getConnectedBanks(): Promise<Bank[]> {
    await delay(500);
    return MOCK_BANKS.filter((b) => b.connected);
  },
};

// Accounts API
export const accountsAPI = {
  async getAccounts(): Promise<Account[]> {
    await delay(500);
    return MOCK_ACCOUNTS;
  },

  async getAccount(id: number): Promise<Account> {
    await delay(500);
    const account = MOCK_ACCOUNTS.find((a) => a.id === id);
    if (!account) throw new Error('Счет не найден');
    return account;
  },

  async getAccountTransactions(accountId: number): Promise<Account['transactions']> {
    await delay(500);
    const account = MOCK_ACCOUNTS.find((a) => a.id === accountId);
    return account?.transactions || [];
  },
};

// Goals API
let goalsData = [...MOCK_GOALS];
let nextGoalId = Math.max(...MOCK_GOALS.map((g) => g.id)) + 1;

export const goalsAPI = {
  async getGoals(): Promise<Goal[]> {
    await delay(500);
    return goalsData.sort((a, b) => a.order - b.order);
  },

  async createGoal(goal: Omit<Goal, 'id' | 'deposits'>): Promise<Goal> {
    await delay(1000);
    const newGoal: Goal = {
      ...goal,
      id: nextGoalId++,
      deposits: [],
    };
    goalsData.push(newGoal);
    return newGoal;
  },

  async updateGoal(id: number, updates: Partial<Goal>): Promise<Goal> {
    await delay(500);
    const goal = goalsData.find((g) => g.id === id);
    if (!goal) throw new Error('Цель не найдена');
    Object.assign(goal, updates);
    return goal;
  },

  async deleteGoal(id: number): Promise<void> {
    await delay(500);
    goalsData = goalsData.filter((g) => g.id !== id);
  },

  async reorderGoals(goalIds: number[]): Promise<void> {
    await delay(500);
    goalIds.forEach((id, index) => {
      const goal = goalsData.find((g) => g.id === id);
      if (goal) goal.order = index + 1;
    });
  },
};

// Loans API
let loansData = [...MOCK_LOANS];
let nextLoanId = Math.max(...MOCK_LOANS.map((l) => l.id)) + 1;

export const loansAPI = {
  async getLoans(): Promise<Loan[]> {
    await delay(500);
    return loansData;
  },

  async createLoan(loan: Omit<Loan, 'id' | 'paymentHistory'>): Promise<Loan> {
    await delay(1000);
    const newLoan: Loan = {
      ...loan,
      id: nextLoanId++,
      paymentHistory: [],
    };
    loansData.push(newLoan);
    return newLoan;
  },

  async updateLoan(id: number, updates: Partial<Loan>): Promise<Loan> {
    await delay(500);
    const loan = loansData.find((l) => l.id === id);
    if (!loan) throw new Error('Кредит не найден');
    Object.assign(loan, updates);
    return loan;
  },

  async deleteLoan(id: number): Promise<void> {
    await delay(500);
    loansData = loansData.filter((l) => l.id !== id);
  },
};

// Operations API
let operationsData = [...MOCK_OPERATIONS];
let nextOperationId = Math.max(...MOCK_OPERATIONS.map((o) => o.id)) + 1;

export const operationsAPI = {
  async getOperations(): Promise<Operation[]> {
    await delay(500);
    return operationsData.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
  },

  async emergencyWithdraw(amount: number): Promise<EmergencyWithdraw> {
    await delay(1000);
    // Simulate calculation - find goals that will be affected
    const carGoal = goalsData.find((g) => g.name === 'Машина');
    
    if (!carGoal) {
      throw new Error('Цель "Машина" не найдена');
    }
    
    // Validate that we have enough in the goal
    if (carGoal.currentAmount < amount) {
      throw new Error(`Недостаточно средств в цели. Доступно: ${carGoal.currentAmount} ₽`);
    }
    
    const affectedGoals = [carGoal.id];
    
    // Calculate how many deposits will be closed (simplified)
    const depositsToClose = Math.ceil(amount / (carGoal.currentAmount / (carGoal.deposits?.length || 1)));
    
    return {
      amount,
      depositsToClose: Array.from({ length: Math.min(depositsToClose, carGoal.deposits?.length || 1) }, (_, i) => i + 1),
      lostInterest: Math.round(amount * 0.015), // ~1.5% lost interest
      affectedGoals,
    };
  },

  async confirmEmergencyWithdraw(amount: number): Promise<void> {
    await delay(1000);
    
    // Find the "Машина" goal
    const carGoal = goalsData.find((g) => g.name === 'Машина');
    if (!carGoal) {
      throw new Error('Цель "Машина" не найдена');
    }

    // Validate that we have enough in the goal
    if (carGoal.currentAmount < amount) {
      throw new Error('Недостаточно средств в цели для снятия');
    }

    // 1. Уменьшить currentAmount цели "Машина" на введенную сумму
    carGoal.currentAmount = Math.max(0, carGoal.currentAmount - amount);
    
    // 2. Увеличить срок на 2 месяца (изменить nextDeposit)
    const currentDate = new Date(carGoal.nextDeposit);
    currentDate.setMonth(currentDate.getMonth() + 2);
    carGoal.nextDeposit = currentDate.toISOString().split('T')[0];

    // 3. Увеличить баланс первого счета (VBank счет #1) на эту сумму
    const firstAccount = MOCK_ACCOUNTS.find((a) => a.id === 1);
    if (firstAccount) {
      firstAccount.balance += amount;
      
      // Добавить транзакцию в историю счета
      if (!firstAccount.transactions) {
        firstAccount.transactions = [];
      }
      const transactionId = Math.max(...firstAccount.transactions.map((t) => t.id), 0) + 1;
      firstAccount.transactions.unshift({
        id: transactionId,
        date: new Date().toISOString().split('T')[0],
        description: `Экстренное снятие из цели "${carGoal.name}"`,
        amount: amount,
        type: 'income',
      });
    }

    // 4. Добавить операцию в историю
    const newOperation: Operation = {
      id: nextOperationId++,
      date: new Date().toISOString().split('T')[0],
      type: 'emergency_withdraw',
      amount: amount,
      goal: carGoal.name,
      status: 'completed',
    };
    operationsData.unshift(newOperation);
  },
};

// Analysis API
export const analysisAPI = {
  async getSalaries(): Promise<SalaryTransaction[]> {
    await delay(2000);
    return MOCK_SALARY_TRANSACTIONS;
  },

  async getSavingsCapacity(): Promise<number> {
    await delay(1000);
    return MOCK_USER.savingsCapacity || 15000;
  },
};

