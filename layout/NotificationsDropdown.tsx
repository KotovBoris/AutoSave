'use client';

import Link from 'next/link';
import { operationsAPI } from '@/lib/api';
import { useEffect, useState } from 'react';
import { Operation } from '@/types';
import { formatDateTime, formatCurrency } from '@/lib/utils';

interface NotificationsDropdownProps {
  onClose: () => void;
}

export default function NotificationsDropdown({ onClose }: NotificationsDropdownProps) {
  const [operations, setOperations] = useState<Operation[]>([]);

  useEffect(() => {
    operationsAPI.getOperations().then(setOperations).catch(console.error);
  }, []);

  return (
    <>
      <div
        className="fixed inset-0 z-10"
        onClick={onClose}
        aria-hidden="true"
      />
      <div className="absolute right-0 top-12 w-80 bg-white rounded-lg shadow-xl border border-gray-200 z-20 max-h-96 overflow-y-auto">
        <div className="p-4 border-b border-gray-200">
          <h3 className="font-semibold text-gray-900">Уведомления</h3>
        </div>

        <div className="divide-y divide-gray-200">
          {operations.slice(0, 5).map((operation) => (
            <div key={operation.id} className="p-4 hover:bg-gray-50">
              <div className="flex items-start gap-3">
                <div className="flex-shrink-0 w-2 h-2 bg-green-500 rounded-full mt-2" />
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-gray-900">
                    {operation.type === 'deposit' && `Открыт вклад "${operation.goal}"`}
                    {operation.type === 'loan_payment' && `Автоплатеж по кредиту "${operation.loan}"`}
                    {operation.type === 'emergency_withdraw' && 'Экстренное снятие'}
                  </p>
                  <p className="text-sm text-gray-500 mt-1">
                    {formatDateTime(operation.date)}
                  </p>
                  <p className="text-sm font-semibold text-gray-900 mt-1">
                    {formatCurrency(operation.amount)}
                  </p>
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="p-4 border-t border-gray-200">
          <Link
            href="/history"
            onClick={onClose}
            className="block text-center text-sm text-primary-600 hover:text-primary-700 font-medium"
          >
            Смотреть все
          </Link>
        </div>
      </div>
    </>
  );
}

