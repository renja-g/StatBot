import Image from "next/image";
import { ApiResponse } from "@/app/components/ApiResponse";

export default function Home() {
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
        <div className="w-full max-w-xl p-6 bg-white dark:bg-black/30 rounded-lg shadow-sm border border-black/[.08] dark:border-white/[.145]">
          <h2 className="text-xl font-semibold mb-4">API Response</h2>
          <ApiResponse />
        </div>
      </main>
    </div>
  );
}
