"use client"
import * as React from "react"
import { Loader2, Inbox, AlertTriangle, RefreshCw } from "lucide-react"
import { cn } from "@/lib/utils"

// Shared loading / empty / error state components for consistent, friendly UI.
// These replace ad-hoc "Loading…" / "No records yet." text so every list page
// presents the same polished feedback and surfaces API errors clearly.

export function LoadingState({ label = "Loading…", className }: { label?: string; className?: string }) {
  return (
    <div className={cn("flex items-center justify-center gap-2 p-6 text-sm text-muted-foreground", className)}>
      <Loader2 className="h-4 w-4 animate-spin" />
      <span>{label}</span>
    </div>
  )
}

export function EmptyState({
  title = "No records yet",
  description,
  action,
  className,
}: {
  title?: string
  description?: string
  action?: React.ReactNode
  className?: string
}) {
  return (
    <div className={cn("flex flex-col items-center justify-center gap-2 px-6 py-10 text-center", className)}>
      <Inbox className="h-8 w-8 text-muted-foreground/50" />
      <p className="text-sm font-medium">{title}</p>
      {description && <p className="max-w-sm text-xs text-muted-foreground">{description}</p>}
      {action}
    </div>
  )
}

export function ErrorState({
  message = "Something went wrong while loading.",
  onRetry,
  className,
}: {
  message?: string
  onRetry?: () => void
  className?: string
}) {
  return (
    <div className={cn("flex flex-col items-center justify-center gap-2 px-6 py-10 text-center", className)}>
      <AlertTriangle className="h-8 w-8 text-destructive/70" />
      <p className="text-sm font-medium text-destructive">Unable to load data</p>
      <p className="max-w-sm text-xs text-muted-foreground">{message}</p>
      {onRetry && (
        <button
          onClick={onRetry}
          className="mt-1 inline-flex items-center gap-1.5 rounded-lg border px-3 py-1.5 text-xs font-medium hover:bg-muted"
        >
          <RefreshCw className="h-3 w-3" /> Retry
        </button>
      )}
    </div>
  )
}

// TableShell wraps a table's rows so loading/empty/error states stay consistent.
export function TableStates({
  loading,
  error,
  empty,
  onRetry,
}: {
  loading: boolean
  error: string
  empty: boolean
  onRetry?: () => void
}) {
  if (loading) return <LoadingState />
  if (error) return <ErrorState message={error} onRetry={onRetry} />
  if (empty) return <EmptyState />
  return null
}
