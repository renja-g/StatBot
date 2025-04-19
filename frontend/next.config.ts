import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: process.env.NODE_ENV === 'development' 
          ? 'http://127.0.0.1:8080/api/:path*'
          : 'http://api:8080/api/:path*',
      },
    ];
  },
  /* config options here */
};

console.log(process.env.NODE_ENV);
export default nextConfig;
