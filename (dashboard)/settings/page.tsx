'use client';

import { useEffect, useState } from 'react';
import { banksAPI, authAPI } from '@/lib/api';
import { Bank, User } from '@/types';

export default function SettingsPage() {
  const [user, setUser] = useState<User | null>(null);
  const [banks, setBanks] = useState<Bank[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [userData, banksData] = await Promise.all([
        authAPI.getCurrentUser(),
        banksAPI.getBanks(),
      ]);
      setUser(userData);
      setBanks(banksData);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDisconnectBank = async (bankId: string) => {
    if (!confirm('Отключить банк?')) return;
    try {
      await banksAPI.disconnectBank(bankId);
      setBanks(banks.map((b) => (b.id === bankId ? { ...b, connected: false } : b)));
    } catch (error) {
      console.error('Error disconnecting bank:', error);
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
      <h1 className="text-3xl font-bold text-gray-900">Настройки</h1>

      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Профиль</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Email</label>
            <input
              type="email"
              value={user?.email || ''}
              readOnly
              className="input bg-gray-50"
            />
          </div>
        </div>
      </div>

      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Подключенные банки</h2>
        <div className="space-y-4">
          {banks.map((bank) => (
            <div
              key={bank.id}
              className="flex items-center justify-between p-4 border rounded-lg"
            >
              <div className="flex items-center gap-4">
                <span className="text-2xl">{bank.logo}</span>
                <div>
                  <p className="font-medium text-gray-900">{bank.name}</p>
                  {bank.connected && bank.connectedAt && (
                    <p className="text-sm text-gray-500">
                      Подключен {new Date(bank.connectedAt).toLocaleDateString('ru-RU')}
                    </p>
                  )}
                </div>
              </div>
              {bank.connected ? (
                <button
                  onClick={() => handleDisconnectBank(bank.id)}
                  className="btn btn-danger"
                >
                  Отключить
                </button>
              ) : (
                <button
                  onClick={() => window.location.href = '/onboarding/banks'}
                  className="btn btn-primary"
                >
                  Подключить
                </button>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

