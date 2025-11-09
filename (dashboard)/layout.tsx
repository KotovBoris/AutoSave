import Topbar from '@/components/layout/Topbar';
import Sidebar from '@/components/layout/Sidebar';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-gray-50">
      <Topbar />
      <Sidebar />
      <main className="ml-60 mt-16 p-8">
        {children}
      </main>
    </div>
  );
}

