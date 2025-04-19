'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function TimelineIndexPage() {
  const router = useRouter();
  const [guildId, setGuildId] = useState('');
  const [userId, setUserId] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!guildId || !userId) {
      setError('Please enter both Guild ID and User ID');
      return;
    }
    
    if (!/^\d+$/.test(guildId) || !/^\d+$/.test(userId)) {
      setError('Guild ID and User ID must be numeric');
      return;
    }
    
    router.push(`/timeline/${guildId}/${userId}`);
  };

  return (
    <div className="flex justify-center items-center min-h-screen p-8 bg-gray-100 dark:bg-gray-900">
      <div className="w-full max-w-md p-6 bg-white dark:bg-black/30 rounded-lg shadow-sm border border-black/[.08] dark:border-white/[.145]">
        <h2 className="text-xl font-semibold mb-6">View User Status Timeline</h2>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="guildId" className="block text-sm font-medium mb-1">
              Guild ID
            </label>
            <input
              id="guildId"
              type="text"
              value={guildId}
              onChange={(e) => setGuildId(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md dark:bg-gray-800 dark:text-white"
              placeholder="Enter Guild ID"
            />
          </div>
          
          <div>
            <label htmlFor="userId" className="block text-sm font-medium mb-1">
              User ID
            </label>
            <input
              id="userId"
              type="text"
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md dark:bg-gray-800 dark:text-white"
              placeholder="Enter User ID"
            />
          </div>
          
          {error && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30 text-red-700 dark:text-red-400 p-3 rounded text-sm">
              {error}
            </div>
          )}
          
          <button
            type="submit"
            className="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors"
          >
            View Timeline
          </button>
        </form>
        
        <div className="mt-6 text-sm text-gray-500 dark:text-gray-400">
          <p>Enter the Discord Guild ID and User ID to view their status timeline.</p>
        </div>
      </div>
    </div>
  );
} 