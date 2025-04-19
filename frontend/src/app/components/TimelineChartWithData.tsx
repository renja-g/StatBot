'use client';

import { useEffect, useState } from 'react';
import { TimelineChart } from './TimelineChart';

interface StatusPeriod {
  startTime: string;
  endTime: string;
  status: string;
}

interface TimelineChartWithDataProps {
  guildId: string;
  userId: string;
  date: string;
}

export function TimelineChartWithData({ guildId, userId, date }: TimelineChartWithDataProps) {
  const [statusData, setStatusData] = useState<StatusPeriod[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStatusChanges = async () => {
      try {
        setLoading(true);
        setError(null);

        const apiUrl = `/api/status-changes/${guildId}/${userId}?date=${date}`;
        const response = await fetch(apiUrl);

        if (!response.ok) {
          throw new Error(`API responded with status: ${response.status}`);
        }

        const data = await response.json();
        setStatusData(data);
      } catch (err) {
        console.error('Error fetching status changes:', err);
        setError(err instanceof Error ? err.message : 'An unknown error occurred');
      } finally {
        setLoading(false);
      }
    };

    fetchStatusChanges();
  }, [guildId, userId, date]);

  const formatTimelineData = () => {
    if (!statusData.length) {
      return [];
    }

    return statusData.map((period, index) => {
      const startTime = new Date(period.startTime);
      const endTime = new Date(period.endTime);
      
      return [
        0, // Index for the y-axis (we only have one row)
        startTime.toISOString().replace('T', ' ').substring(0, 19), // Format: YYYY-MM-DD HH:MM:SS
        endTime.toISOString().replace('T', ' ').substring(0, 19),
        mapApiStatusToChartStatus(period.status)
      ];
    });
  };

  // Map API status values to the ones expected by the chart
  const mapApiStatusToChartStatus = (status: string): string => {
    // Map from API status names to what the chart expects
    const statusMap: Record<string, string> = {
      'online': 'Online',
      'idle': 'Idle',
      'dnd': 'Do Not Disturb',
      'offline': 'Offline'
    };

    return statusMap[status.toLowerCase()] || status;
  };

  return (
    <div className="w-full h-full">
      {loading && (
        <div className="flex justify-center items-center h-full">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 dark:border-white"></div>
        </div>
      )}

      {error && (
        <div className="flex justify-center items-center h-full">
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30 text-red-700 dark:text-red-400 p-4 rounded max-w-md">
            <p className="text-sm">Error: {error}</p>
            <p className="text-xs mt-1">Make sure the API service is running.</p>
          </div>
        </div>
      )}

      {!loading && !error && statusData.length === 0 && (
        <div className="flex justify-center items-center h-full">
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-900/30 text-blue-700 dark:text-blue-400 p-4 rounded max-w-md text-center">
            <p>No status data available for this user on the selected date.</p>
            <p className="text-xs mt-2">Try selecting a different date or check that the user ID is correct.</p>
          </div>
        </div>
      )}

      {!loading && !error && statusData.length > 0 && (
        <TimelineChart 
          timelineData={formatTimelineData()} 
          minDate={`${date} 00:00:00`} 
          maxDate={`${date} 23:59:59`}
        />
      )}
    </div>
  );
} 