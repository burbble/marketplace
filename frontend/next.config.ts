import type { NextConfig } from "next";

const backendUrl = process.env.BACKEND_URL || "http://localhost:38080";

const nextConfig: NextConfig = {
  output: "standalone",
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "store77.net",
      },
      {
        protocol: "https",
        hostname: "**.store77.net",
      },
    ],
  },
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${backendUrl}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;
