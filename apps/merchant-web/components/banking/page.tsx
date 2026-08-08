"use client"
import * as React from "react"
import { api } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useRequireAuth } from "@/lib/api/require-auth"

// Shared banking page scaffold: enforces a dashboard session, loads data from a given
// loader, and renders a title + a table-ish list. Keeps banking pages consistent and
// real-data-backed instead of duplicated mock scaffolds.

export type Column<T> = {
  key: keyof T
  label: string
  render?: (row: T) => React.ReactNode
}

export function BankingPage<T extends { id: string }>({
  title,
  subtitle,
  columns,
  loader,
}: {
  title: string
  subtitle?: string
  columns: Column<T>[]
  loader: () => Promise<T[]>
}) {
  const { checking } = useRequireAuth()
  const { data, loading } = useData(loader, [])

  if (checking) {
    return (
      <Centered>
        <p className="text-sm text-muted-foreground">Checking session…</p>
      </Centered>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold">{title}</h1>
          {subtitle && <p className="text-sm text-muted-foreground mt-2">{subtitle}</p>}
        </div>

        <div className="rounded-2xl border bg-card overflow-hidden">
          <div className="grid gap-2 bg-muted p-3 text-[11px] font-semibold"
            style={{ gridTemplateColumns: `repeat(${columns.length}, minmax(0,1fr))` }}>
            {columns.map((c) => (
              <span key={String(c.key)}>{c.label}</span>
            ))}
          </div>
          {loading && <p className="p-4 text-sm text-muted-foreground">Loading…</p>}
          {!loading && (data ?? []).length === 0 && (
            <p className="p-4 text-sm text-muted-foreground">No records yet.</p>
          )}
          {(data ?? []).map((row) => (
            <div key={row.id} className="grid gap-2 p-3 border-t text-xs hover:bg-muted/50"
              style={{ gridTemplateColumns: `repeat(${columns.length}, minmax(0,1fr))` }}>
              {columns.map((c) => (
                <span key={String(c.key)}>
                  {c.render ? c.render(row) : String(row[c.key] ?? "—")}
                </span>
              ))}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-neutral-50">
      {children}
    </div>
  )
}

export { api }
