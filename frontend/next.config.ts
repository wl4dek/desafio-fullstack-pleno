import type { NextConfig } from "next"

const API_URL = process.env.API_URL || "http://localhost:8080"

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      { source: "/auth/:path*", destination: `${API_URL}/auth/:path*` },
      { source: "/api/v1/:path*", destination: `${API_URL}/api/v1/:path*` },
    ]
  },
}

export default nextConfig
