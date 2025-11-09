import { User, Bank, Account, Goal, Loan, Operation, Transaction } from '@/types';

export const MOCK_USER: User = {
  id: 1,
  email: 'demo@example.com',
  avgSalary: 85000,
  avgExpenses: 60000,
  savingsCapacity: 15000,
};

export const MOCK_BANKS: Bank[] = [
  { id: 'vbank', name: 'VBank', logo: '🏦', connected: true, connectedAt: '2024-01-15' },
  { id: 'abank', name: 'ABank', logo: '💰', connected: true, connectedAt: '2024-01-15' },
  { id: 'sbank', name: 'SBank', logo: '🏛️', connected: false },
];

export const MOCK_ACCOUNTS: Account[] = [
  {
    id: 1,
    bankId: 'vbank',
    number: '1234',
    balance: 150000,
    transactions: [
      { id: 1, date: '2024-01-15', description: 'Зарплата', amount: 85000, type: 'income', sender: 'ООО "Компания А"' },
      { id: 2, date: '2024-01-14', description: 'Продукты', amount: -2500, type: 'expense', category: 'Еда' },
      { id: 3, date: '2024-01-13', description: 'Кафе', amount: -800, type: 'expense', category: 'Развлечения' },
    ],
  },
  {
    id: 2,
    bankId: 'vbank',
    number: '5678',
    balance: 50000,
    transactions: [],
  },
  {
    id: 3,
    bankId: 'abank',
    number: '9012',
    balance: 80000,
    transactions: [],
  },
];

export const MOCK_GOALS: Goal[] = [
  {
    id: 1,
    name: 'Машина',
    targetAmount: 500000,
    currentAmount: 50000,
    monthlyAmount: 15000,
    nextDeposit: '2024-02-01',
    bankId: 'vbank',
    order: 1,
    deposits: [
      { id: 1, goalId: 1, amount: 15000, date: '2024-01-15', status: 'completed' },
    ],
  },
  {
    id: 2,
    name: 'Квартира',
    targetAmount: 2000000,
    currentAmount: 0,
    monthlyAmount: 20000,
    nextDeposit: '2024-02-01',
    bankId: 'abank',
    order: 2,
    deposits: [],
  },
];

export const MOCK_LOANS: Loan[] = [
  {
    id: 1,
    name: 'Ипотека Сбер',
    debt: 1500000,
    rate: 12,
    monthlyPayment: 25000,
    nextPayment: '2024-02-05',
    bankId: 'vbank',
    autoPayment: true,
    paymentHistory: [
      { id: 1, loanId: 1, amount: 25000, date: '2024-01-05', status: 'completed' },
    ],
  },
];

export const MOCK_OPERATIONS: Operation[] = [
  {
    id: 1,
    date: '2024-01-15',
    type: 'deposit',
    amount: 15000,
    goal: 'Машина',
    status: 'completed',
  },
  {
    id: 2,
    date: '2024-01-10',
    type: 'loan_payment',
    amount: 12000,
    loan: 'Ипотека Сбер',
    status: 'completed',
  },
  {
    id: 3,
    date: '2024-01-05',
    type: 'deposit',
    amount: 20000,
    goal: 'Квартира',
    status: 'completed',
  },
];

export const MOCK_SALARY_TRANSACTIONS = [
  {
    id: 1,
    date: '2024-01-15',
    amount: 85000,
    bankId: 'vbank',
    sender: 'ООО "Компания А"',
    isSalary: true,
  },
  {
    id: 2,
    date: '2023-12-15',
    amount: 85000,
    bankId: 'vbank',
    sender: 'ООО "Компания А"',
    isSalary: true,
  },
  {
    id: 3,
    date: '2024-01-20',
    amount: 50000,
    bankId: 'abank',
    sender: 'ИП Иванов',
    isSalary: false,
  },
];

