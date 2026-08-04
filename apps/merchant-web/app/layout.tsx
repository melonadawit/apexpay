import type { Metadata } from "next"
import "./globals.css"
import { Inter } from "next/font/google"

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" })

export const metadata: Metadata = {
  title: "ApexPay Merchant — NBE Onboarding + Fayda",
  description: "AI-native payment gateway for Ethiopia — outstanding merchant onboarding with Fayda ID front/back verification",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className={`${inter.variable} font-sans antialiased bg-neutral-50`}>{children}</body>
    </html>
  )
}
