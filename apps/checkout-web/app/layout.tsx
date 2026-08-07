import "./globals.css"
import type { Metadata } from "next"
export const metadata: Metadata = { title: "ApexPay Checkout", description: "Hosted payment checkout" }
export default function RootLayout({ children }: { children: React.ReactNode }) { return <html lang="en"><body>{children}</body></html> }
