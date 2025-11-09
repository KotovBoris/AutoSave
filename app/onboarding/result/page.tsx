'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { analysisAPI } from '@/lib/api';
import { formatCurrency } from '@/lib/utils';

export default function ResultPage() {
  const router = useRouter();
  const [savingsCapacity, setSavingsCapacity] = useState<number>(0);

  useEffect(() => {
    loadCapacity();
  }, []);

  const loadCapacity = async () => {
    try {
      const capacity = await analysisAPI.getSavingsCapacity();
      setSavingsCapacity(capacity);
    } catch (error) {
      console.error('Error loading capacity:', error);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-primary-100 flex items-center justify-center p-8">
      <div className="max-w-2xl w-full bg-white rounded-lg shadow-lg p-8 text-center">
        <div className="text-6xl mb-6">✅</div>
        <h1 className="text-3xl font-bold text-gray-900 mb-8">Анализ завершен</h1>

        <div className="bg-primary-50 rounded-lg p-8 mb-6">
          <p className="text-gray-600 mb-4">Вы можете откладывать</p>
          <p className="text-5xl font-bold text-primary-600 mb-4">
            {formatCurrency(savingsCapacity)}
          </p>
          <p className="text-gray-600">в месяц</p>
        </div>

        <div className="text-left space-y-3 mb-8">
          <div className="flex justify-between">
            <span className="text-gray-600">Средняя зарплата:</span>
            <span className="font-semibold text-gray-900">{formatCurrency(85000)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600">Средние расходы:</span>
            <span className="font-semibold text-gray-900">{formatCurrency(60000)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600">Доступно для накоплений:</span>
            <span className="font-semibold text-primary-600">{formatCurrency(25000)}</span>
          </div>
        </div>

        <button
          onClick={() => router.push('/dashboard')}
          className="btn btn-primary w-full"
        >
          Перейти к дашборду
        </button>
      </div>
    </div>
  );
}

