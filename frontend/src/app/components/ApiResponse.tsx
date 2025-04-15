'use client';

import { useState, useEffect } from 'react';

export function ApiResponse() {
  const [data, setData] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const response = await fetch('/api/');
        
        if (!response.ok) {
          throw new Error(`API responded with status: ${response.status}`);
        }
        
        const text = await response.text();
        setData(text);
        setError(null);
      } catch (err) {
        console.error('Error fetching from API:', err);
        setError(err instanceof Error ? err.message : 'An unknown error occurred');
        setData(null);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  return (
    <div className="w-full">
      {loading && (
        <div className="flex justify-center items-center py-4">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900 dark:border-white"></div>
        </div>
      )}
      
      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30 text-red-700 dark:text-red-400 p-4 rounded">
          <p className="text-sm">Error: {error}</p>
          <p className="text-xs mt-1">Make sure the API service is running on port 8080.</p>
        </div>
      )}
      
      {!loading && !error && data && (
        <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-900/30 p-4 rounded">
          <p className="font-medium text-green-800 dark:text-green-400">Response:</p>
          <p className="mt-1 font-mono text-sm break-all">{data}</p>
        </div>
      )}
    </div>
  );
} 