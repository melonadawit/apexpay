import * as React from "react"
import { cn } from "@/lib/utils"
export const Input = React.forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>(({className,...props},ref)=>(
  <input className={cn("w-full rounded-xl border border-black/10 bg-white h-12 px-3 focus:ring-2 focus:ring-primary focus:outline-none",className)} ref={ref} {...props} />
))
Input.displayName="Input"
