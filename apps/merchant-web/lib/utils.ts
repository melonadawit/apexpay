import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

// cn - shadcn outstanding utility for conditional classes
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// format ETB with Ethiopian locale
export function formatETB(amount: number | string | { toString(): string }) {
  const num = typeof amount === "string" ? parseFloat(amount) : typeof amount === "number" ? amount : parseFloat(amount.toString())
  return new Intl.NumberFormat("en-ET", { style: "currency", currency: "ETB", minimumFractionDigits: 2 }).format(num)
}

// mask helpers - outstanding privacy
export function maskFINLast4(last4: string) {
  return `****-****-${last4}`
}

export function maskAccount(acc: string) {
  if (!acc) return "****"
  if (acc.length <= 4) return "****"
  return `****${acc.slice(-4)}`
}

// glare detection for Fayda capture - optimal canvas brightness check
export function detectGlare(imageData: ImageData, threshold = 200): { hasGlare: boolean; brightness: number } {
  const data = imageData.data
  let total = 0
  let brightCount = 0
  // sample every 4th pixel for performance O(n/4)
  for (let i = 0; i < data.length; i += 16) {
    const r = data[i], g = data[i+1], b = data[i+2]
    const brightness = (r + g + b) / 3
    total += brightness
    if (brightness > threshold) brightCount++
  }
  const avg = total / (data.length / 16)
  const glareRatio = brightCount / (data.length / 16)
  return { hasGlare: glareRatio > 0.15, brightness: avg }
}

// progress donut calculation - optimal
export function calcProgress(completed: number, total: number) {
  if (total === 0) return 0
  return Math.round((completed / total) * 100)
}

// debounce utility - optimal for search
export function debounce<T extends (...args: any[]) => any>(fn: T, delay: number) {
  let timeoutId: ReturnType<typeof setTimeout>
  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => fn(...args), delay)
  }
}
