export interface User {
  id: number;
  email: string;
  avgSalary?: number;
  avgExpenses?: number;
  savingsCapacity?: number;
}

export interface Bank {
  id: string;
  name: string;
  logo: string;
  connected: boolean;
  connectedAt?: string;
}

export interface Account {
  id: number;
  bankId: string;
  number: string;
  balance: number;
  transactions?: Transaction[];
}

export interface Transaction {
  id: number;
  date: string;
  description: string;
  amount: number;
  type: 'income' | 'expense';
  category?: string;
  sender?: string;
}

export interface Goal {
  id: number;
  name: string;
  targetAmount: number;
  currentAmount: number;
  monthlyAmount: number;
  nextDeposit: string;
  bankId: string;
  order: number;
  deposits?: Deposit[];
}

export interface Deposit {
  id: number;
  goalId: number;
  amount: number;
  date: string;
  status: 'pending' | 'completed' | 'failed';
}

export interface Loan {
  id: number;
  name: string;
  debt: number;
  rate: number;
  monthlyPayment: number;
  nextPayment: string;
  bankId: string;
  autoPayment: boolean;
  paymentHistory?: Payment[];
}

export interface Payment {
  id: number;
  loanId: number;
  amount: number;
  date: string;
  status: 'pending' | 'completed' | 'failed';
}

export interface Operation {
  id: number;
  date: string;
  type: 'deposit' | 'loan_payment' | 'emergency_withdraw';
  amount: number;
  goal?: string;
  loan?: string;
  status: 'pending' | 'completed' | 'failed';
}

export interface SalaryTransaction {
  id: number;
  date: string;
  amount: number;
  bankId: string;
  sender: string;
  isSalary: boolean;
}

export interface EmergencyWithdraw {
  amount: number;
  depositsToClose: number[];
  lostInterest: number;
  affectedGoals: number[];
}

