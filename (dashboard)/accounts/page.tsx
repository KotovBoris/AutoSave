'use client';

import { useEffect, useState } from 'react';
import { accountsAPI } from '@/lib/api';
import { Account, Transaction } from '@/types';
import { formatCurrency, formatDate } from '@/lib/utils';
import Modal from '@/components/ui/Modal';

export default function AccountsPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedAccount, setSelectedAccount] = useState<Account | null>(null);

  useEffect(() => {
    loadAccounts();
  }, []);

  const loadAccounts = async () => {
    setLoading(true);
    try {
      const accountsData = await accountsAPI.getAccounts();
      setAccounts(accountsData);
    } catch (error) {
      console.error('Error loading accounts:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleViewTransactions = async (accountId: number) => {
    try {
      const transactions = await accountsAPI.getAccountTransactions(accountId);
      const account = accounts.find((a) => a.id === accountId);
      if (account) {
        setSelectedAccount({ ...account, transactions });
      }
    } catch (error) {
      console.error('Error loading transactions:', error);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Все счета</h1>
        <button onClick={loadAccounts} className="btn btn-secondary">
          Обновить данные
        </button>
      </div>

      {accounts.length === 0 ? (
        <div className="card text-center py-12">
          <div className="text-6xl mb-4">💳</div>
          <h3 className="text-xl font-semibold text-gray-900 mb-2">
            Нет подключенных счетов
          </h3>
          <p className="text-gray-600 mb-6">
            Подключите банки чтобы увидеть счета
          </p>
        </div>
      ) : (
        <div className="space-y-4">
          {accounts.map((account) => (
            <div key={account.id} className="card">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-1">
                    {account.bankId.toUpperCase()} - Счет {account.number}
                  </h3>
                  <p className="text-2xl font-bold text-gray-900">
                    {formatCurrency(account.balance)}
                  </p>
                </div>
                <button
                  onClick={() => handleViewTransactions(account.id)}
                  className="btn btn-outline"
                >
                  Последние операции
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {selectedAccount && (
        <Modal
          isOpen={true}
          onClose={() => setSelectedAccount(null)}
          title={`Транзакции - ${selectedAccount.bankId.toUpperCase()} Счет ${selectedAccount.number}`}
          size="lg"
        >
          <div className="space-y-2 max-h-96 overflow-y-auto">
            {selectedAccount.transactions && selectedAccount.transactions.length > 0 ? (
              selectedAccount.transactions.map((transaction) => (
                <div
                  key={transaction.id}
                  className="flex items-center justify-between p-4 border rounded-lg"
                >
                  <div>
                    <p className="font-medium text-gray-900">{transaction.description}</p>
                    <p className="text-sm text-gray-500">{formatDate(transaction.date)}</p>
                    {transaction.sender && (
                      <p className="text-sm text-gray-500">От: {transaction.sender}</p>
                    )}
                  </div>
                  <p
                    className={`font-semibold ${
                      transaction.amount > 0 ? 'text-green-600' : 'text-red-600'
                    }`}
                  >
                    {transaction.amount > 0 ? '+' : ''}
                    {formatCurrency(transaction.amount)}
                  </p>
                </div>
              ))
            ) : (
              <p className="text-center text-gray-500 py-8">Нет транзакций</p>
            )}
          </div>
        </Modal>
      )}
    </div>
  );
}

