'use client';

import { useState } from 'react';
import Modal from '../ui/Modal';
import { operationsAPI } from '@/lib/api';
import { EmergencyWithdraw } from '@/types';
import { formatCurrency } from '@/lib/utils';

interface EmergencyWithdrawModalProps {
  onClose: () => void;
}

export default function EmergencyWithdrawModal({ onClose }: EmergencyWithdrawModalProps) {
  const [amount, setAmount] = useState('');
  const [loading, setLoading] = useState(false);
  const [calculation, setCalculation] = useState<EmergencyWithdraw | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleCalculate = async () => {
    if (!amount || parseFloat(amount) <= 0) {
      setError('Введите корректную сумму');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const result = await operationsAPI.emergencyWithdraw(parseFloat(amount));
      setCalculation(result);
    } catch (err) {
      setError('Не удалось рассчитать снятие. Попробуйте снова.');
    } finally {
      setLoading(false);
    }
  };

  const handleConfirm = async () => {
    if (!calculation) return;
    
    setLoading(true);
    setError(null);
    
    try {
      await operationsAPI.confirmEmergencyWithdraw(calculation.amount);
      // Success - dispatch event to notify other components
      window.dispatchEvent(new CustomEvent('emergencyWithdrawSuccess', { 
        detail: { amount: calculation.amount } 
      }));
      
      // Close modal first
      onClose();
      
      // Show success notification after a short delay
      setTimeout(() => {
        alert(`✅ ${formatCurrency(calculation.amount)} возвращены на счет VBank`);
      }, 100);
    } catch (err: any) {
      setError(err.message || 'Не удалось выполнить снятие. Банк недоступен.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal isOpen={true} onClose={onClose} title="Экстренное снятие" size="md">
      <div className="space-y-6">
        {!calculation ? (
          <>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Сколько нужно вернуть?
              </label>
              <input
                type="number"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="Введите сумму"
                className="input"
                min="0"
                step="1000"
              />
            </div>

            {error && (
              <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
                {error}
              </div>
            )}

            <div className="flex gap-4 justify-end">
              <button onClick={onClose} className="btn btn-secondary">
                Отмена
              </button>
              <button
                onClick={handleCalculate}
                disabled={loading || !amount}
                className="btn btn-primary"
              >
                {loading ? 'Рассчитываем...' : 'Рассчитать'}
              </button>
            </div>
          </>
        ) : (
          <>
            <div>
              <p className="text-lg font-semibold text-gray-900 mb-4">
                Нужно: {formatCurrency(calculation.amount)}
              </p>

              <div className="space-y-4">
                <div>
                  <h3 className="font-medium text-gray-900 mb-2">
                    Будут закрыты вклады:
                  </h3>
                  <ul className="space-y-2">
                    {calculation.depositsToClose.map((depositId) => (
                      <li key={depositId} className="flex items-center gap-2">
                        <span className="text-green-600">✓</span>
                        <span className="text-gray-700">
                          Вклад #{depositId}: {formatCurrency(calculation.amount / calculation.depositsToClose.length)}
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>

                <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
                  <p className="text-sm text-yellow-800">
                    Потеря процентов: {formatCurrency(calculation.lostInterest)}
                  </p>
                </div>

                {calculation.affectedGoals.length > 0 && (
                  <div className="p-4 bg-orange-50 border border-orange-200 rounded-lg">
                    <p className="text-sm text-orange-800">
                      ⚠️ Цель будет отложена
                    </p>
                  </div>
                )}
              </div>
            </div>

            {error && (
              <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
                {error}
              </div>
            )}

            <div className="flex gap-4 justify-end">
              <button
                onClick={() => {
                  setCalculation(null);
                  setAmount('');
                  setError(null);
                }}
                className="btn btn-secondary"
                disabled={loading}
              >
                Назад
              </button>
              <button
                onClick={handleConfirm}
                disabled={loading}
                className="btn btn-danger"
              >
                {loading ? 'Выполняется...' : 'Подтвердить снятие'}
              </button>
            </div>
          </>
        )}
      </div>
    </Modal>
  );
}

