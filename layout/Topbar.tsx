'use client';

import { useState } from 'react';
import Link from 'next/link';
import { FiBell, FiAlertTriangle } from 'react-icons/fi';
import { useRouter } from 'next/navigation';
import EmergencyWithdrawModal from '../modals/EmergencyWithdrawModal';
import NotificationsDropdown from './NotificationsDropdown';

export default function Topbar() {
  const [showEmergencyModal, setShowEmergencyModal] = useState(false);
  const [showNotifications, setShowNotifications] = useState(false);
  const router = useRouter();

  return (
    <>
      <header className="bg-white shadow-sm border-b border-gray-200 h-16 flex items-center justify-between px-6 fixed top-0 left-0 right-0 z-40">
        <Link href="/dashboard" className="text-2xl font-bold text-primary-600">
          AutoSave
        </Link>

        <div className="flex items-center gap-4">
          <button
            onClick={() => setShowEmergencyModal(true)}
            className="btn border-2 border-danger-600 text-danger-600 hover:bg-danger-600 hover:text-white flex items-center gap-2"
          >
            <FiAlertTriangle className="w-5 h-5" />
            Экстренное снятие
          </button>

          <div className="relative">
            <button
              onClick={() => setShowNotifications(!showNotifications)}
              className="relative p-2 hover:bg-gray-100 rounded-lg transition-colors"
            >
              <FiBell className="w-6 h-6 text-gray-600" />
              <span className="absolute top-0 right-0 w-5 h-5 bg-danger-600 text-white text-xs rounded-full flex items-center justify-center">
                3
              </span>
            </button>

            {showNotifications && (
              <NotificationsDropdown onClose={() => setShowNotifications(false)} />
            )}
          </div>
        </div>
      </header>

      {showEmergencyModal && (
        <EmergencyWithdrawModal 
          onClose={() => setShowEmergencyModal(false)}
        />
      )}
    </>
  );
}

