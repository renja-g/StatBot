import { TimelineChart } from "../components/TimelineChart";

export default function TimelinePage() {
  return (
    <div className="flex justify-center items-center min-h-screen p-8 bg-gray-100 dark:bg-gray-900">
      <div className="w-full max-w-5xl p-6 bg-white dark:bg-black/30 rounded-lg shadow-sm border border-black/[.08] dark:border-white/[.145]">
        <h2 className="text-xl font-semibold mb-4">User Status Timeline</h2>
        <div className="h-[300px] w-full">
          <TimelineChart />
        </div>
      </div>
    </div>
  );
} 