/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    typedRoutes: false,
  },
  // Allow all images for demo
  images: {
    remotePatterns: [{ hostname: '**' }],
  },
}

module.exports = nextConfig
