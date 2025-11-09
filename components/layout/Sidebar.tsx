'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
  FiLayout,
  FiCreditCard,
  FiTarget,
  FiDollarSign,
  FiClock,
  FiSettings,
} from 'react-icons/fi';
import { cn } from '@/lib/utils';

const menuItems = [
  { href: '/dashboard', icon: FiLayout, label: 'Дашборд' },
  { href: '/accounts', icon: FiCreditCard, label: 'Все счета' },
  { href: '/goals', icon: FiTarget, label: 'Мои цели' },
  { href: '/loans', icon: FiDollarSign, label: 'Кредиты' },
  { href: '/history', icon: FiClock, label: 'История' },
  { href: '/settings', icon: FiSettings, label: 'Настройки' },
];

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="w-60 bg-gray-50 border-r border-gray-200 fixed left-0 top-16 bottom-0 overflow-y-auto">
      <nav className="p-4 space-y-1">
        {menuItems.map((item) => {
          const Icon = item.icon;
          const isActive = pathname === item.href;

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                'flex items-center gap-3 px-4 py-3 rounded-lg transition-colors',
                isActive
                  ? 'bg-primary-600 text-white shadow-md'
                  : 'text-gray-700 hover:bg-gray-100'
              )}
            >
              <Icon className="w-5 h-5" />
              <span className="font-medium">{item.label}</span>
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}

