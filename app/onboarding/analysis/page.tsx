'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function AnalysisPage() {
  const router = useRouter();
  const [progress, setProgress] = useState(0);
  const [stage, setStage] = useState('Загружаем транзакции...');

  useEffect(() => {
    const stages = [
      { progress: 30, text: 'Загружаем транзакции...' },
      { progress: 60, text: 'Анализируем доходы...' },
      { progress: 90, text: 'Определяем зарплаты...' },
      { progress: 100, text: 'Рассчитываем возможности...' },
    ];

    let currentStage = 0;
    const interval = setInterval(() => {
      if (currentStage < stages.length) {
        setProgress(stages[currentStage].progress);
        setStage(stages[currentStage].text);
        currentStage++;
      } else {
        clearInterval(interval);
        setTimeout(() => {
          router.push('/onboarding/salaries');
        }, 500);
      }
    }, 1500);

    return () => clearInterval(interval);
  }, [router]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-primary-100 flex items-center justify-center p-8">
      <div className="max-w-md w-full bg-white rounded-lg shadow-lg p-8 text-center">
        <div className="mb-8">
          <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-primary-600 mx-auto mb-4" />
        </div>

        <h2 className="text-2xl font-bold text-gray-900 mb-4">
          Анализируем ваши транзакции
        </h2>

        <p className="text-gray-600 mb-6">{stage}</p>

        <div className="w-full bg-gray-200 rounded-full h-4">
          <div
            className="bg-primary-600 h-4 rounded-full transition-all duration-300"
            style={{ width: `${progress}%` }}
          />
        </div>

        <p className="text-sm text-gray-500 mt-4">{progress}%</p>
      </div>
    </div>
  );
}

