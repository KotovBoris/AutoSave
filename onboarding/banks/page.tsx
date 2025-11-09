'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { banksAPI } from '@/lib/api';
import { Bank } from '@/types';

export default function BanksPage() {
  const router = useRouter();
  const [banks, setBanks] = useState<Bank[]>([]);
  const [loading, setLoading] = useState<string | null>(null);
  const [connected, setConnected] = useState<string[]>([]);

  useEffect(() => {
    loadBanks();
  }, []);

  const loadBanks = async () => {
    try {
      const banksData = await banksAPI.getBanks();
      setBanks(banksData);
      setConnected(banksData.filter((b) => b.connected).map((b) => b.id));
    } catch (error) {
      console.error('Error loading banks:', error);
    }
  };

  const handleConnect = async (bankId: string) => {
    setLoading(bankId);
    try {
      await banksAPI.connectBank(bankId);
      setConnected([...connected, bankId]);
    } catch (error) {
      console.error('Error connecting bank:', error);
      alert('Не удалось подключить банк. Попробуйте снова.');
    } finally {
      setLoading(null);
    }
  };

  const handleContinue = () => {
    router.push('/onboarding/analysis');
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-primary-100 p-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-primary-600 mb-2">AutoSave</h1>
          <button
            onClick={() => router.push('/dashboard')}
            className="text-gray-600 hover:text-gray-900"
          >
            Пропустить пока
          </button>
        </div>

        <div className="bg-white rounded-lg shadow-lg p-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">
            Подключите банки для начала
          </h2>
          <p className="text-gray-600 mb-8">
            Мы проанализируем транзакции и определим сколько вы можете откладывать
          </p>

          <div className="space-y-4 mb-8">
            {banks.map((bank) => (
              <div
                key={bank.id}
                className="border rounded-lg p-6 flex items-center justify-between"
              >
                <div className="flex items-center gap-4">
                  <span className="text-4xl">{bank.logo}</span>
                  <span className="font-semibold text-lg">{bank.name}</span>
                </div>

                {connected.includes(bank.id) ? (
                  <span className="text-green-600 font-semibold">✅ Подключен</span>
                ) : loading === bank.id ? (
                  <button disabled className="btn btn-secondary">
                    ⏳ Подключение...
                  </button>
                ) : (
                  <button
                    onClick={() => handleConnect(bank.id)}
                    className="btn btn-primary"
                  >
                    Подключить
                  </button>
                )}
              </div>
            ))}
          </div>

          <button
            onClick={handleContinue}
            disabled={connected.length === 0}
            className="btn btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Продолжить {connected.length > 0 && `с ${connected.length} банком${connected.length > 1 ? 'ами' : ''}`}
          </button>
        </div>
      </div>
    </div>
  );
}

