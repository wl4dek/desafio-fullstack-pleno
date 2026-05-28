import type { NextConfig } from "next"

const API_URL = process.env.API_URL || "http://localhost:8080"

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/auth/:path*",
        destination: `${API_URL}/auth/:path*`,
      },
      {
        source: "/summary",
        destination: `${API_URL}/summary`,
      },
    ]
  },
}

export default nextConfig
