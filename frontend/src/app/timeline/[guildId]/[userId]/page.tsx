'use client';

import { useState } from 'react';
import { TimelineChartWithData } from '@/app/components/TimelineChartWithData';
import { useParams } from 'next/navigation';

export default function UserTimelinePage() {
  const params = useParams();
  const guildId = params.guildId as string;
  const userId = params.userId as string;
  const [selectedDate, setSelectedDate] = useState<string>(
    new Date().toISOString().split('T')[0] // Today's date in YYYY-MM-DD format
  );

  const handleDateChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedDate(event.target.value);
  };

  return (
    <div className="flex justify-center items-center min-h-screen p-8 bg-gray-100 dark:bg-gray-900">
      <div className="w-full max-w-5xl p-6 bg-white dark:bg-black/30 rounded-lg shadow-sm border border-black/[.08] dark:border-white/[.145]">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">User Status Timeline</h2>
          <div>
            <input
              type="date"
              value={selectedDate}
              onChange={handleDateChange}
              className="px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md dark:bg-gray-800 dark:text-white"
            />
          </div>
        </div>
        <div className="h-[300px] w-full">
          <TimelineChartWithData guildId={guildId} userId={userId} date={selectedDate} />
        </div>
      </div>
    </div>
  );
} 