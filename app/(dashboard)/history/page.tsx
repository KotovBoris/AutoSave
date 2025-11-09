'use client';

import { useEffect, useState } from 'react';
import { operationsAPI } from '@/lib/api';
import { Operation } from '@/types';
import { formatDateTime, formatCurrency } from '@/lib/utils';

export default function HistoryPage() {
  const [operations, setOperations] = useState<Operation[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadOperations();
  }, []);

  const loadOperations = async () => {
    setLoading(true);
    try {
      const operationsData = await operationsAPI.getOperations();
      setOperations(operationsData);
    } catch (error) {
      console.error('Error loading operations:', error);
    } finally {
      setLoading(false);
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
        <h1 className="text-3xl font-bold text-gray-900">История операций</h1>
        <button onClick={loadOperations} className="btn btn-secondary">
          Обновить
        </button>
      </div>

      {operations.length === 0 ? (
        <div className="card text-center py-12">
          <div className="text-6xl mb-4">📈</div>
          <h3 className="text-xl font-semibold text-gray-900 mb-2">
            Нет операций
          </h3>
          <p className="text-gray-600">
            История операций появится здесь
          </p>
        </div>
      ) : (
        <div className="card">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left py-3 px-4 font-semibold text-gray-900">Дата</th>
                  <th className="text-left py-3 px-4 font-semibold text-gray-900">Операция</th>
                  <th className="text-left py-3 px-4 font-semibold text-gray-900">Сумма</th>
                  <th className="text-left py-3 px-4 font-semibold text-gray-900">Статус</th>
                </tr>
              </thead>
              <tbody>
                {operations.map((operation) => (
                  <tr key={operation.id} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="py-3 px-4 text-gray-700">
                      {formatDateTime(operation.date)}
                    </td>
                    <td className="py-3 px-4">
                      <div>
                        <p className="font-medium text-gray-900">
                          {operation.type === 'deposit' && `Открыт вклад "${operation.goal}"`}
                          {operation.type === 'loan_payment' && `Автоплатеж "${operation.loan}"`}
                          {operation.type === 'emergency_withdraw' && 'Экстренное снятие'}
                        </p>
                      </div>
                    </td>
                    <td className="py-3 px-4 font-semibold text-gray-900">
                      {formatCurrency(operation.amount)}
                    </td>
                    <td className="py-3 px-4">
                      <span
                        className={`px-2 py-1 rounded text-xs font-medium ${
                          operation.status === 'completed'
                            ? 'bg-green-100 text-green-800'
                            : operation.status === 'pending'
                            ? 'bg-yellow-100 text-yellow-800'
                            : 'bg-red-100 text-red-800'
                        }`}
                      >
                        {operation.status === 'completed' && 'Выполнено'}
                        {operation.status === 'pending' && 'В обработке'}
                        {operation.status === 'failed' && 'Ошибка'}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}

