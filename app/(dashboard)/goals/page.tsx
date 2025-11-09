'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { goalsAPI, accountsAPI } from '@/lib/api';
import { Goal, Account } from '@/types';
import { formatCurrency, calculateProgress, calculateMonthsToGoal } from '@/lib/utils';
import { FiPlus, FiEdit, FiX, FiArrowUp, FiArrowDown } from 'react-icons/fi';
import Modal from '@/components/ui/Modal';

export default function GoalsPage() {
  const [goals, setGoals] = useState<Goal[]>([]);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingGoal, setEditingGoal] = useState<Goal | null>(null);
  const [deletingGoal, setDeletingGoal] = useState<Goal | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [goalsData, accountsData] = await Promise.all([
        goalsAPI.getGoals(),
        accountsAPI.getAccounts(),
      ]);
      setGoals(goalsData);
      setAccounts(accountsData);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleMove = async (goalId: number, direction: 'up' | 'down') => {
    const currentIndex = goals.findIndex((g) => g.id === goalId);
    if (currentIndex === -1) return;

    const newIndex = direction === 'up' ? currentIndex - 1 : currentIndex + 1;
    if (newIndex < 0 || newIndex >= goals.length) return;

    const newGoals = [...goals];
    [newGoals[currentIndex], newGoals[newIndex]] = [newGoals[newIndex], newGoals[currentIndex]];

    // Update orders
    newGoals.forEach((goal, index) => {
      goal.order = index + 1;
    });

    setGoals(newGoals);
    await goalsAPI.reorderGoals(newGoals.map((g) => g.id));
  };

  const handleDelete = async () => {
    if (!deletingGoal) return;
    try {
      await goalsAPI.deleteGoal(deletingGoal.id);
      setGoals(goals.filter((g) => g.id !== deletingGoal.id));
      setDeletingGoal(null);
    } catch (error) {
      console.error('Error deleting goal:', error);
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
        <h1 className="text-3xl font-bold text-gray-900">Мои цели</h1>
        <button
          onClick={() => setShowCreateModal(true)}
          className="btn btn-primary flex items-center gap-2"
        >
          <FiPlus className="w-5 h-5" />
          Создать цель
        </button>
      </div>

      {goals.length === 0 ? (
        <div className="card text-center py-12">
          <div className="text-6xl mb-4">🎯</div>
          <h3 className="text-xl font-semibold text-gray-900 mb-2">
            У вас пока нет целей
          </h3>
          <p className="text-gray-600 mb-6">
            Создайте первую цель накопления
          </p>
          <button
            onClick={() => setShowCreateModal(true)}
            className="btn btn-primary"
          >
            Создать цель
          </button>
        </div>
      ) : (
        <div className="space-y-4">
          {goals.map((goal, index) => (
            <div key={goal.id} className="card">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-4 mb-4">
                    <h3 className="text-xl font-semibold text-gray-900">{goal.name}</h3>
                    <span className="text-sm text-gray-500">#{goal.order}</span>
                  </div>

                  <div className="mb-4">
                    <div className="flex justify-between text-sm text-gray-600 mb-2">
                      <span>Прогресс</span>
                      <span>
                        {formatCurrency(goal.currentAmount)} из {formatCurrency(goal.targetAmount)}
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-3">
                      <div
                        className="bg-primary-600 h-3 rounded-full transition-all"
                        style={{ width: `${calculateProgress(goal.currentAmount, goal.targetAmount)}%` }}
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-4 text-sm">
                    <div>
                      <span className="text-gray-500">В месяц:</span>
                      <p className="font-semibold text-gray-900">{formatCurrency(goal.monthlyAmount)}</p>
                    </div>
                    <div>
                      <span className="text-gray-500">Следующий депозит:</span>
                      <p className="font-semibold text-gray-900">
                        {new Date(goal.nextDeposit).toLocaleDateString('ru-RU')}
                      </p>
                    </div>
                    <div>
                      <span className="text-gray-500">Осталось месяцев:</span>
                      <p className="font-semibold text-gray-900">
                        {calculateMonthsToGoal(goal.targetAmount, goal.currentAmount, goal.monthlyAmount)}
                      </p>
                    </div>
                  </div>
                </div>

                <div className="flex flex-col gap-2 ml-4">
                  <button
                    onClick={() => handleMove(goal.id, 'up')}
                    disabled={index === 0}
                    className="p-2 hover:bg-gray-100 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <FiArrowUp className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => handleMove(goal.id, 'down')}
                    disabled={index === goals.length - 1}
                    className="p-2 hover:bg-gray-100 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <FiArrowDown className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => setEditingGoal(goal)}
                    className="p-2 hover:bg-gray-100 rounded-lg"
                  >
                    <FiEdit className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => setDeletingGoal(goal)}
                    className="p-2 hover:bg-red-100 rounded-lg text-red-600"
                  >
                    <FiX className="w-5 h-5" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {showCreateModal && (
        <GoalFormModal
          accounts={accounts}
          existingGoalsCount={goals.length}
          onClose={() => setShowCreateModal(false)}
          onSuccess={() => {
            setShowCreateModal(false);
            loadData();
          }}
        />
      )}

      {editingGoal && (
        <GoalFormModal
          accounts={accounts}
          goal={editingGoal}
          onClose={() => setEditingGoal(null)}
          onSuccess={() => {
            setEditingGoal(null);
            loadData();
          }}
        />
      )}

      {deletingGoal && (
        <Modal
          isOpen={true}
          onClose={() => setDeletingGoal(null)}
          title={`Закрыть цель "${deletingGoal.name}"?`}
        >
          <div className="space-y-4">
            <p className="text-gray-700">
              Все вклады будут закрыты. Потеря процентов: {formatCurrency(1200)}
            </p>
            <div className="flex gap-4 justify-end">
              <button onClick={() => setDeletingGoal(null)} className="btn btn-secondary">
                Отмена
              </button>
              <button onClick={handleDelete} className="btn btn-danger">
                Да, закрыть
              </button>
            </div>
          </div>
        </Modal>
      )}
    </div>
  );
}

interface GoalFormModalProps {
  accounts: Account[];
  goal?: Goal;
  onClose: () => void;
  onSuccess: () => void;
  existingGoalsCount?: number;
}

function GoalFormModal({ accounts, goal, onClose, onSuccess, existingGoalsCount = 0 }: GoalFormModalProps) {
  const [name, setName] = useState(goal?.name || '');
  const [targetAmount, setTargetAmount] = useState(goal?.targetAmount.toString() || '');
  const [monthlyAmount, setMonthlyAmount] = useState(goal?.monthlyAmount.toString() || '');
  const [bankId, setBankId] = useState(goal?.bankId || accounts[0]?.bankId || '');
  const [loading, setLoading] = useState(false);
  const [calculation, setCalculation] = useState<{
    months: number;
    interest: number;
  } | null>(null);

  const handleCalculate = () => {
    if (!targetAmount || !monthlyAmount) return;

    const target = parseFloat(targetAmount);
    const monthly = parseFloat(monthlyAmount);
    const current = goal?.currentAmount || 0;

    const months = calculateMonthsToGoal(target, current, monthly);
    const interest = (target * 0.07 * months) / 12; // Simplified interest calculation

    setCalculation({ months, interest });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      if (goal) {
        await goalsAPI.updateGoal(goal.id, {
          name,
          targetAmount: parseFloat(targetAmount),
          monthlyAmount: parseFloat(monthlyAmount),
          bankId,
        });
      } else {
        await goalsAPI.createGoal({
          name,
          targetAmount: parseFloat(targetAmount),
          monthlyAmount: parseFloat(monthlyAmount),
          bankId,
          currentAmount: 0,
          nextDeposit: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
          order: existingGoalsCount + 1,
        });
      }
      onSuccess();
    } catch (error) {
      console.error('Error saving goal:', error);
      alert('Не удалось сохранить цель. Попробуйте снова.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      isOpen={true}
      onClose={onClose}
      title={goal ? 'Редактирование цели' : 'Создание цели'}
      size="md"
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Название
          </label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            className="input"
            placeholder="Например: Машина"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Сумма цели
          </label>
          <input
            type="number"
            value={targetAmount}
            onChange={(e) => setTargetAmount(e.target.value)}
            required
            min="0"
            step="1000"
            className="input"
            placeholder="500000"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Сумма откладывать в месяц
          </label>
          <input
            type="number"
            value={monthlyAmount}
            onChange={(e) => setMonthlyAmount(e.target.value)}
            required
            min="0"
            step="1000"
            className="input"
            placeholder="15000"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Банк для вклада
          </label>
          <select
            value={bankId}
            onChange={(e) => setBankId(e.target.value)}
            required
            className="input"
          >
            {accounts.map((account) => (
              <option key={account.id} value={account.bankId}>
                {account.bankId.toUpperCase()} - Счет {account.number}
              </option>
            ))}
          </select>
        </div>

        {!calculation && (
          <button
            type="button"
            onClick={handleCalculate}
            className="btn btn-secondary w-full"
          >
            Рассчитать план
          </button>
        )}

        {calculation && (
          <div className="p-4 bg-gray-50 rounded-lg space-y-2">
            <p className="text-sm text-gray-600">Срок достижения: {calculation.months} месяцев</p>
            <p className="text-sm text-gray-600">
              Доход на процентах: {formatCurrency(calculation.interest)}
            </p>
          </div>
        )}

        <div className="flex gap-4 justify-end pt-4">
          <button type="button" onClick={onClose} className="btn btn-secondary">
            Отмена
          </button>
          <button type="submit" disabled={loading} className="btn btn-primary">
            {loading ? 'Сохранение...' : goal ? 'Сохранить' : 'Создать цель'}
          </button>
        </div>
      </form>
    </Modal>
  );
}

