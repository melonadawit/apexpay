import type { Metadata } from "next"
import "./globals.css"
import { Inter } from "next/font/google"
import { ThemeProvider } from "@/components/providers/theme-provider"
import { LanguageProvider } from "@/components/providers/language-provider"
import { Shell } from "@/components/layout/shell"

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" })

export const metadata: Metadata = {
  title: "ApexPay Merchant — NBE Onboarding + Fayda",
  description: "AI-native payment gateway for Ethiopia — outstanding merchant onboarding with Fayda ID verification, payouts, payroll, subscriptions.",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${inter.variable} font-sans antialiased`}>
        <ThemeProvider>
          <LanguageProvider>
            <Shell>
              {children}
            </Shell>
          </LanguageProvider>
        </ThemeProvider>
      </body>
    </html>
  )
}
