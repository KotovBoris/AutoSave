'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { analysisAPI } from '@/lib/api';
import { SalaryTransaction } from '@/types';
import { formatDate, formatCurrency } from '@/lib/utils';

export default function SalariesPage() {
  const router = useRouter();
  const [transactions, setTransactions] = useState<SalaryTransaction[]>([]);
  const [selected, setSelected] = useState<number[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadSalaries();
  }, []);

  const loadSalaries = async () => {
    try {
      const data = await analysisAPI.getSalaries();
      setTransactions(data);
      // Auto-select transactions marked as salary
      setSelected(data.filter((t) => t.isSalary).map((t) => t.id));
    } catch (error) {
      console.error('Error loading salaries:', error);
    } finally {
      setLoading(false);
    }
  };

  const toggleSelection = (id: number) => {
    if (selected.includes(id)) {
      setSelected(selected.filter((i) => i !== id));
    } else {
      setSelected([...selected, id]);
    }
  };

  const handleContinue = () => {
    router.push('/onboarding/result');
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-primary-50 to-primary-100 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
      </div>
    );
  }

  if (transactions.length === 0) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-primary-50 to-primary-100 flex items-center justify-center p-8">
        <div className="max-w-md w-full bg-white rounded-lg shadow-lg p-8 text-center">
          <div className="text-6xl mb-4">📋</div>
          <h2 className="text-2xl font-bold text-gray-900 mb-4">
            Мы не нашли регулярных зарплат
          </h2>
          <p className="text-gray-600 mb-6">
            Вы можете настроить это позже в разделе "Настройки"
          </p>
          <button onClick={handleContinue} className="btn btn-primary w-full">
            Пропустить
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-primary-100 p-8">
      <div className="max-w-2xl mx-auto bg-white rounded-lg shadow-lg p-8">
        <h2 className="text-2xl font-bold text-gray-900 mb-4">
          Подтвердите ваши зарплаты
        </h2>
        <p className="text-gray-600 mb-6">
          Мы нашли возможные зарплаты по транзакциям. Отметьте галочками правильные.
        </p>

        <div className="space-y-3 mb-6">
          {transactions.map((transaction) => (
            <label
              key={transaction.id}
              className="flex items-start gap-3 p-4 border rounded-lg cursor-pointer hover:bg-gray-50"
            >
              <input
                type="checkbox"
                checked={selected.includes(transaction.id)}
                onChange={() => toggleSelection(transaction.id)}
                className="mt-1 w-5 h-5 text-primary-600 rounded focus:ring-primary-500"
              />
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <span className="font-medium text-gray-900">
                    {formatDate(transaction.date)}
                  </span>
                  <span className="font-semibold text-gray-900">
                    {formatCurrency(transaction.amount)}
                  </span>
                </div>
                <p className="text-sm text-gray-600 mt-1">
                  {transaction.bankId.toUpperCase()} • От: {transaction.sender}
                </p>
              </div>
            </label>
          ))}
        </div>

        <div className="flex gap-4">
          <button
            onClick={() => router.back()}
            className="btn btn-secondary flex-1"
          >
            Назад
          </button>
          <button
            onClick={handleContinue}
            disabled={selected.length === 0}
            className="btn btn-primary flex-1 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Продолжить
          </button>
        </div>
      </div>
    </div>
  );
}

