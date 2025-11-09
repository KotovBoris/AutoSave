'use client';

import { useEffect, useState } from 'react';
import { accountsAPI, goalsAPI, loansAPI, analysisAPI } from '@/lib/api';
import { Account, Goal, Loan } from '@/types';
import { formatCurrency, calculateProgress, calculateMonthsToGoal } from '@/lib/utils';
import Link from 'next/link';
import { FiArrowRight, FiRefreshCw } from 'react-icons/fi';

export default function DashboardPage() {
  const [loading, setLoading] = useState(true);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [goals, setGoals] = useState<Goal[]>([]);
  const [loans, setLoans] = useState<Loan[]>([]);
  const [savingsCapacity, setSavingsCapacity] = useState<number>(0);

  useEffect(() => {
    loadData();
    
    // Listen for emergency withdraw success event
    const handleEmergencyWithdrawSuccess = () => {
      loadData(); // Reload data when emergency withdraw is completed
    };
    
    window.addEventListener('emergencyWithdrawSuccess', handleEmergencyWithdrawSuccess);
    
    return () => {
      window.removeEventListener('emergencyWithdrawSuccess', handleEmergencyWithdrawSuccess);
    };
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [accountsData, goalsData, loansData, capacity] = await Promise.all([
        accountsAPI.getAccounts(),
        goalsAPI.getGoals(),
        loansAPI.getLoans(),
        analysisAPI.getSavingsCapacity(),
      ]);

      setAccounts(accountsData);
      setGoals(goalsData);
      setLoans(loansData);
      setSavingsCapacity(capacity);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const totalBalance = accounts.reduce((sum, acc) => sum + acc.balance, 0);
  const totalInGoals = goals.reduce((sum, goal) => sum + goal.currentAmount, 0);
  const totalDebt = loans.reduce((sum, loan) => sum + loan.debt, 0);
  const activeGoal = goals.length > 0 ? goals[0] : null;

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Дашборд</h1>
        <button
          onClick={loadData}
          className="btn btn-secondary flex items-center gap-2"
        >
          <FiRefreshCw className="w-4 h-4" />
          Обновить
        </button>
      </div>

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card">
          <h3 className="text-sm font-medium text-gray-500 mb-2">Всего на счетах</h3>
          <p className="text-3xl font-bold text-gray-900">{formatCurrency(totalBalance)}</p>
        </div>

        <div className="card">
          <h3 className="text-sm font-medium text-gray-500 mb-2">В целях</h3>
          <p className="text-3xl font-bold text-primary-600">{formatCurrency(totalInGoals)}</p>
        </div>

        <div className="card">
          <h3 className="text-sm font-medium text-gray-500 mb-2">В долгах</h3>
          <p className="text-3xl font-bold text-red-600">{formatCurrency(totalDebt)}</p>
        </div>
      </div>

      {/* Active Goal */}
      {activeGoal ? (
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-900">Текущая цель: {activeGoal.name}</h2>
            <Link
              href="/goals"
              className="text-primary-600 hover:text-primary-700 font-medium flex items-center gap-1"
            >
              Все цели <FiArrowRight className="w-4 h-4" />
            </Link>
          </div>

          <div className="space-y-4">
            <div>
              <div className="flex justify-between text-sm text-gray-600 mb-2">
                <span>Прогресс</span>
                <span>
                  {formatCurrency(activeGoal.currentAmount)} из {formatCurrency(activeGoal.targetAmount)}
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-4">
                <div
                  className="bg-primary-600 h-4 rounded-full transition-all"
                  style={{ width: `${calculateProgress(activeGoal.currentAmount, activeGoal.targetAmount)}%` }}
                />
              </div>
            </div>

            <div className="grid grid-cols-3 gap-4 text-sm">
              <div>
                <span className="text-gray-500">Откладываем в месяц:</span>
                <p className="font-semibold text-gray-900">{formatCurrency(activeGoal.monthlyAmount)}</p>
              </div>
              <div>
                <span className="text-gray-500">Следующий депозит:</span>
                <p className="font-semibold text-gray-900">{new Date(activeGoal.nextDeposit).toLocaleDateString('ru-RU')}</p>
              </div>
              <div>
                <span className="text-gray-500">Осталось месяцев:</span>
                <p className="font-semibold text-gray-900">
                  {calculateMonthsToGoal(activeGoal.targetAmount, activeGoal.currentAmount, activeGoal.monthlyAmount)}
                </p>
              </div>
            </div>
          </div>
        </div>
      ) : (
        <div className="card text-center py-12">
          <div className="text-6xl mb-4">🎯</div>
          <h3 className="text-xl font-semibold text-gray-900 mb-2">У вас пока нет целей</h3>
          <p className="text-gray-600 mb-6">Создайте первую цель накопления</p>
          <Link href="/goals" className="btn btn-primary inline-flex items-center gap-2">
            Создать цель <FiArrowRight className="w-4 h-4" />
          </Link>
        </div>
      )}

      {/* Next Actions */}
      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Следующие действия</h2>
        <div className="space-y-3">
          <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
            <div>
              <p className="font-medium text-gray-900">Вы можете откладывать</p>
              <p className="text-sm text-gray-600">Рекомендуемая сумма на основе анализа</p>
            </div>
            <p className="text-2xl font-bold text-primary-600">{formatCurrency(savingsCapacity)}/мес</p>
          </div>

          {loans.length > 0 && (
            <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
              <div>
                <p className="font-medium text-gray-900">Следующий платеж по кредиту</p>
                <p className="text-sm text-gray-600">{loans[0].name}</p>
              </div>
              <div className="text-right">
                <p className="font-semibold text-gray-900">{formatCurrency(loans[0].monthlyPayment)}</p>
                <p className="text-sm text-gray-600">{new Date(loans[0].nextPayment).toLocaleDateString('ru-RU')}</p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Quick Links */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Link href="/accounts" className="card hover:shadow-lg transition-shadow">
          <h3 className="font-semibold text-gray-900 mb-2">Все счета</h3>
          <p className="text-sm text-gray-600">Просмотр всех счетов и транзакций</p>
        </Link>

        <Link href="/goals" className="card hover:shadow-lg transition-shadow">
          <h3 className="font-semibold text-gray-900 mb-2">Мои цели</h3>
          <p className="text-sm text-gray-600">Управление целями накопления</p>
        </Link>

        <Link href="/loans" className="card hover:shadow-lg transition-shadow">
          <h3 className="font-semibold text-gray-900 mb-2">Кредиты</h3>
          <p className="text-sm text-gray-600">Управление кредитами и автоплатежами</p>
        </Link>
      </div>
    </div>
  );
}

