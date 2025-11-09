'use client';

import { useEffect, useState } from 'react';
import { loansAPI, accountsAPI } from '@/lib/api';
import { Loan, Account } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { FiPlus, FiEdit, FiX, FiChevronDown, FiChevronUp } from 'react-icons/fi';
import Modal from '@/components/ui/Modal';

export default function LoansPage() {
  const [loans, setLoans] = useState<Loan[]>([]);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'list' | 'add'>('list');
  const [showAddModal, setShowAddModal] = useState(false);
  const [editingLoan, setEditingLoan] = useState<Loan | null>(null);
  const [deletingLoan, setDeletingLoan] = useState<Loan | null>(null);
  const [expandedLoan, setExpandedLoan] = useState<number | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [loansData, accountsData] = await Promise.all([
        loansAPI.getLoans(),
        accountsAPI.getAccounts(),
      ]);
      setLoans(loansData);
      setAccounts(accountsData);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!deletingLoan) return;
    try {
      await loansAPI.deleteLoan(deletingLoan.id);
      setLoans(loans.filter((l) => l.id !== deletingLoan.id));
      setDeletingLoan(null);
    } catch (error) {
      console.error('Error deleting loan:', error);
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
        <h1 className="text-3xl font-bold text-gray-900">Кредиты</h1>
      </div>

      <div className="border-b border-gray-200">
        <nav className="flex gap-4">
          <button
            onClick={() => setActiveTab('list')}
            className={`px-4 py-2 font-medium border-b-2 transition-colors ${
              activeTab === 'list'
                ? 'border-primary-600 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Мои кредиты
          </button>
          <button
            onClick={() => setActiveTab('add')}
            className={`px-4 py-2 font-medium border-b-2 transition-colors ${
              activeTab === 'add'
                ? 'border-primary-600 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Добавить кредит
          </button>
        </nav>
      </div>

      {activeTab === 'list' && (
        <>
          {loans.length === 0 ? (
            <div className="card text-center py-12">
              <div className="text-6xl mb-4">💰</div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">
                Кредитов не найдено
              </h3>
              <p className="text-gray-600 mb-6">
                Добавьте кредит для автоматизации платежей
              </p>
              <button
                onClick={() => setActiveTab('add')}
                className="btn btn-primary"
              >
                Добавить кредит
              </button>
            </div>
          ) : (
            <div className="space-y-4">
              {loans.map((loan) => (
                <div key={loan.id} className="card">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-4 mb-4">
                        <h3 className="text-xl font-semibold text-gray-900">{loan.name}</h3>
                        {loan.autoPayment && (
                          <span className="px-2 py-1 bg-green-100 text-green-800 text-xs font-medium rounded">
                            Автоплатеж
                          </span>
                        )}
                      </div>

                      <div className="grid grid-cols-4 gap-4 text-sm">
                        <div>
                          <span className="text-gray-500">Остаток:</span>
                          <p className="font-semibold text-gray-900">{formatCurrency(loan.debt)}</p>
                        </div>
                        <div>
                          <span className="text-gray-500">Ставка:</span>
                          <p className="font-semibold text-gray-900">{loan.rate}% годовых</p>
                        </div>
                        <div>
                          <span className="text-gray-500">Платеж:</span>
                          <p className="font-semibold text-gray-900">{formatCurrency(loan.monthlyPayment)}/мес</p>
                        </div>
                        <div>
                          <span className="text-gray-500">Следующий платеж:</span>
                          <p className="font-semibold text-gray-900">
                            {new Date(loan.nextPayment).toLocaleDateString('ru-RU')}
                          </p>
                        </div>
                      </div>

                      {expandedLoan === loan.id && (
                        <div className="mt-4 pt-4 border-t border-gray-200 space-y-4">
                          <div>
                            <h4 className="font-medium text-gray-900 mb-2">Настройки автоплатежа</h4>
                            <div className="space-y-2 text-sm">
                              <div className="flex items-center gap-2">
                                <input
                                  type="checkbox"
                                  checked={loan.autoPayment}
                                  readOnly
                                  className="w-4 h-4 text-primary-600 rounded"
                                />
                                <span>Автоплатеж включен</span>
                              </div>
                              <p className="text-gray-600">
                                Дата: {new Date(loan.nextPayment).getDate()} число каждого месяца
                              </p>
                              <p className="text-gray-600">Сумма: {formatCurrency(loan.monthlyPayment)}</p>
                            </div>
                          </div>
                        </div>
                      )}
                    </div>

                    <div className="flex flex-col gap-2 ml-4">
                      <button
                        onClick={() =>
                          setExpandedLoan(expandedLoan === loan.id ? null : loan.id)
                        }
                        className="p-2 hover:bg-gray-100 rounded-lg"
                      >
                        {expandedLoan === loan.id ? (
                          <FiChevronUp className="w-5 h-5" />
                        ) : (
                          <FiChevronDown className="w-5 h-5" />
                        )}
                      </button>
                      <button
                        onClick={() => setEditingLoan(loan)}
                        className="p-2 hover:bg-gray-100 rounded-lg"
                      >
                        <FiEdit className="w-5 h-5" />
                      </button>
                      <button
                        onClick={() => setDeletingLoan(loan)}
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
        </>
      )}

      {activeTab === 'add' && (
        <LoanFormModal
          accounts={accounts}
          onClose={() => setActiveTab('list')}
          onSuccess={() => {
            setActiveTab('list');
            loadData();
          }}
        />
      )}

      {editingLoan && (
        <LoanFormModal
          accounts={accounts}
          loan={editingLoan}
          onClose={() => setEditingLoan(null)}
          onSuccess={() => {
            setEditingLoan(null);
            loadData();
          }}
        />
      )}

      {deletingLoan && (
        <Modal
          isOpen={true}
          onClose={() => setDeletingLoan(null)}
          title={`Удалить кредит "${deletingLoan.name}"?`}
        >
          <div className="space-y-4">
            <p className="text-gray-700">
              Автоплатежи будут отменены.
            </p>
            <div className="flex gap-4 justify-end">
              <button onClick={() => setDeletingLoan(null)} className="btn btn-secondary">
                Отмена
              </button>
              <button onClick={handleDelete} className="btn btn-danger">
                Удалить
              </button>
            </div>
          </div>
        </Modal>
      )}
    </div>
  );
}

interface LoanFormModalProps {
  accounts: Account[];
  loan?: Loan;
  onClose: () => void;
  onSuccess: () => void;
}

function LoanFormModal({ accounts, loan, onClose, onSuccess }: LoanFormModalProps) {
  const [name, setName] = useState(loan?.name || '');
  const [debt, setDebt] = useState(loan?.debt.toString() || '');
  const [rate, setRate] = useState(loan?.rate.toString() || '');
  const [monthlyPayment, setMonthlyPayment] = useState(loan?.monthlyPayment.toString() || '');
  const [nextPayment, setNextPayment] = useState(
    loan?.nextPayment ? new Date(loan.nextPayment).getDate().toString() : '5'
  );
  const [autoPayment, setAutoPayment] = useState(loan?.autoPayment || false);
  const [bankId, setBankId] = useState(loan?.bankId || accounts[0]?.bankId || '');
  const [loading, setLoading] = useState(false);
  const [calculation, setCalculation] = useState<{
    months: number;
    overpayment: number;
  } | null>(null);

  const handleCalculate = () => {
    if (!debt || !rate || !monthlyPayment) return;

    const debtAmount = parseFloat(debt);
    const rateValue = parseFloat(rate);
    const payment = parseFloat(monthlyPayment);

    // Simplified calculation
    const months = Math.ceil(debtAmount / payment);
    const totalPayment = months * payment;
    const overpayment = totalPayment - debtAmount;

    setCalculation({ months, overpayment });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const nextPaymentDate = new Date();
      nextPaymentDate.setDate(parseInt(nextPayment));
      if (nextPaymentDate < new Date()) {
        nextPaymentDate.setMonth(nextPaymentDate.getMonth() + 1);
      }

      if (loan) {
        await loansAPI.updateLoan(loan.id, {
          name,
          debt: parseFloat(debt),
          rate: parseFloat(rate),
          monthlyPayment: parseFloat(monthlyPayment),
          nextPayment: nextPaymentDate.toISOString().split('T')[0],
          bankId,
          autoPayment,
        });
      } else {
        await loansAPI.createLoan({
          name,
          debt: parseFloat(debt),
          rate: parseFloat(rate),
          monthlyPayment: parseFloat(monthlyPayment),
          nextPayment: nextPaymentDate.toISOString().split('T')[0],
          bankId,
          autoPayment,
        });
      }
      onSuccess();
    } catch (error) {
      console.error('Error saving loan:', error);
      alert('Не удалось сохранить кредит. Попробуйте снова.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      isOpen={true}
      onClose={onClose}
      title={loan ? 'Редактирование кредита' : 'Добавить кредит'}
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
            placeholder="Например: Ипотека Сбер"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Остаток долга
            </label>
            <input
              type="number"
              value={debt}
              onChange={(e) => setDebt(e.target.value)}
              required
              min="0"
              step="1000"
              className="input"
              placeholder="1500000"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Процентная ставка (% годовых)
            </label>
            <input
              type="number"
              value={rate}
              onChange={(e) => setRate(e.target.value)}
              required
              min="0"
              step="0.1"
              className="input"
              placeholder="12"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Минимальный платеж
            </label>
            <input
              type="number"
              value={monthlyPayment}
              onChange={(e) => setMonthlyPayment(e.target.value)}
              required
              min="0"
              step="1000"
              className="input"
              placeholder="25000"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Дата платежа (число месяца)
            </label>
            <input
              type="number"
              value={nextPayment}
              onChange={(e) => setNextPayment(e.target.value)}
              required
              min="1"
              max="31"
              className="input"
              placeholder="5"
            />
          </div>
        </div>

        <div>
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={autoPayment}
              onChange={(e) => setAutoPayment(e.target.checked)}
              className="w-4 h-4 text-primary-600 rounded"
            />
            <span className="text-sm font-medium text-gray-700">Включить автоплатеж</span>
          </label>
        </div>

        {autoPayment && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Банк для списания
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
        )}

        {!calculation && (
          <button
            type="button"
            onClick={handleCalculate}
            className="btn btn-secondary w-full"
          >
            Рассчитать график
          </button>
        )}

        {calculation && (
          <div className="p-4 bg-gray-50 rounded-lg space-y-2">
            <p className="text-sm text-gray-600">
              Срок погашения: {calculation.months} месяцев
            </p>
            <p className="text-sm text-gray-600">
              Переплата: {formatCurrency(calculation.overpayment)}
            </p>
          </div>
        )}

        <div className="flex gap-4 justify-end pt-4">
          <button type="button" onClick={onClose} className="btn btn-secondary">
            Отмена
          </button>
          <button type="submit" disabled={loading} className="btn btn-primary">
            {loading ? 'Сохранение...' : loan ? 'Сохранить' : 'Добавить кредит'}
          </button>
        </div>
      </form>
    </Modal>
  );
}

