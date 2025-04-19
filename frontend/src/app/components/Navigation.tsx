'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

export function Navigation() {
  const pathname = usePathname();

  const isActive = (path: string) => {
    return pathname === path || pathname.startsWith(`${path}/`);
  };

  return (
    <nav className="bg-white dark:bg-black/30 border-b border-black/[.08] dark:border-white/[.145] py-4 px-6">
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        <div className="flex items-center">
          <Link href="/" className="text-lg font-semibold">
            StatBot
          </Link>
        </div>
        
        <div className="flex space-x-6">
          <Link 
            href="/timeline" 
            className={`text-sm ${
              isActive('/timeline') 
                ? 'text-blue-600 dark:text-blue-400 font-medium' 
                : 'text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white'
            }`}
          >
            User Timeline
          </Link>
        </div>
      </div>
    </nav>
  );
} 