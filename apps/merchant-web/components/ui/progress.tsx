"use client"
import * as React from "react"
import * as ProgressPrimitive from "@radix-ui/react-progress"
import { cn } from "@/lib/utils"

const Progress = React.forwardRef<React.ElementRef<typeof ProgressPrimitive.Root>, React.ComponentPropsWithoutRef<typeof ProgressPrimitive.Root>>(({ className, value, ...props }, ref) => (
  <ProgressPrimitive.Root ref={ref} className={cn("relative h-3 w-full overflow-hidden rounded-full bg-neutral-100", className)} {...props}>
    <ProgressPrimitive.Indicator className="h-full w-full flex-1 bg-primary transition-all duration-500 ease-out" style={{ transform: `translateX(-${100 - (value || 0)}%)` }} />
  </ProgressPrimitive.Root>
))
Progress.displayName = ProgressPrimitive.Root.displayName

// Donut progress for onboarding checklist - outstanding
export function DonutProgress({ value, size = 64 }: { value: number; size?: number }) {
  const radius = (size - 8) / 2
  const circ = 2 * Math.PI * radius
  const offset = circ - (value / 100) * circ
  return (
    <div className="relative" style={{ width: size, height: size }}>
      <svg width={size} height={size} className="-rotate-90">
        <circle cx={size/2} cy={size/2} r={radius} fill="none" stroke="#E4E4E7" strokeWidth="6" />
        <circle cx={size/2} cy={size/2} r={radius} fill="none" stroke="#0B6E4F" strokeWidth="6" strokeDasharray={circ} strokeDashoffset={offset} strokeLinecap="round" className="transition-all duration-700 ease-out" />
      </svg>
      <div className="absolute inset-0 flex items-center justify-center text-sm font-bold">{value}%</div>
    </div>
  )
}

export { Progress }
