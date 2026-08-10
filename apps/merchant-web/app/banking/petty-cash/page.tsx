"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type PettyCashBudget, type PettyCashExpense } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function PettyCashPage() {
  const { checking } = useRequireAuth()
  const { data: budgets, refetch: refetchBudgets } = useData(() => api.banking.pettyCashBudgets(), [])
  const { data: expenses, refetch: refetchExpenses } = useData(() => api.banking.pettyCashExpenses(), [])

  const [budgetName, setBudgetName] = React.useState("")
  const [budgetAmount, setBudgetAmount] = React.useState("")
  const [budgetErr, setBudgetErr] = React.useState("")
  const [creatingBudget, setCreatingBudget] = React.useState(false)

  const [budgetID, setBudgetID] = React.useState("")
  const [expenseAmount, setExpenseAmount] = React.useState("")
  const [expenseDesc, setExpenseDesc] = React.useState("")
  const [expenseErr, setExpenseErr] = React.useState("")
  const [creatingExpense, setCreatingExpense] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const createBudget = async () => {
    setCreatingBudget(true); setBudgetErr("")
    try {
      await api.banking.createPettyCashBudget({ budget_name: budgetName, amount: budgetAmount })
      setBudgetName(""); setBudgetAmount(""); refetchBudgets()
    } catch (e) { setBudgetErr((e as Error).message) } finally { setCreatingBudget(false) }
  }

  const createExpense = async () => {
    setCreatingExpense(true); setExpenseErr("")
    try {
      await api.banking.createPettyCashExpense({ budget_id: budgetID, amount: expenseAmount, description: expenseDesc })
      setExpenseAmount(""); setExpenseDesc(""); refetchExpenses()
    } catch (e) { setExpenseErr((e as Error).message) } finally { setCreatingExpense(false) }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Petty Cash • አነስተኛ ገንዘብ</h1>
          <p className="text-sm text-muted-foreground mt-2">Track petty cash budgets and expenses with receipts.</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Budgets</h3>
            <div className="flex gap-2">
              <input value={budgetName} onChange={(e) => setBudgetName(e.target.value)} placeholder="Budget name" className="flex-1 rounded-xl border h-11 px-3 text-sm" />
              <input value={budgetAmount} onChange={(e) => setBudgetAmount(e.target.value)} placeholder="Amount" className="w-28 rounded-xl border h-11 px-3 text-sm" />
              <button onClick={createBudget} disabled={creatingBudget} className="rounded-xl bg-primary text-white px-4 text-sm disabled:opacity-50">Add</button>
            </div>
            {budgetErr && <p className="text-sm text-red-600">{budgetErr}</p>}
            <div className="space-y-2">
              {(budgets ?? []).map((b) => (
                <div key={b.id} className="rounded-xl border p-3">
                  <div className="flex justify-between items-center">
                    <p className="text-sm font-medium">{b.budget_name}</p>
                    <span className="text-[11px] px-2 py-0.5 rounded-full bg-green-500/15 text-green-700">{b.status}</span>
                  </div>
                  <p className="text-xs text-muted-foreground mt-1">
                    ETB {b.amount} • Spent ETB {b.spent_amount} • Remaining ETB {b.remaining_amount}
                  </p>
                </div>
              ))}
            </div>
          </div>

          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Expenses</h3>
            <div className="space-y-2">
              <select value={budgetID} onChange={(e) => setBudgetID(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
                <option value="">Select budget…</option>
                {(budgets ?? []).map((b) => <option key={b.id} value={b.id}>{b.budget_name}</option>)}
              </select>
              <div className="flex gap-2">
                <input value={expenseAmount} onChange={(e) => setExpenseAmount(e.target.value)} placeholder="Amount" className="w-28 rounded-xl border h-11 px-3 text-sm" />
                <input value={expenseDesc} onChange={(e) => setExpenseDesc(e.target.value)} placeholder="Description" className="flex-1 rounded-xl border h-11 px-3 text-sm" />
                <button onClick={createExpense} disabled={creatingExpense} className="rounded-xl bg-primary text-white px-4 text-sm disabled:opacity-50">Add</button>
              </div>
            </div>
            {expenseErr && <p className="text-sm text-red-600">{expenseErr}</p>}
            <div className="space-y-2">
              {(expenses ?? []).map((e) => (
                <div key={e.id} className="rounded-xl border p-3 flex justify-between">
                  <div>
                    <p className="text-sm font-medium">ETB {e.amount} • {e.description}</p>
                    <p className="text-xs text-muted-foreground">Budget {e.budget_id}</p>
                  </div>
                  <span className="text-[11px] px-2 py-0.5 rounded-full bg-amber-500/15 text-amber-700">{e.status}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
