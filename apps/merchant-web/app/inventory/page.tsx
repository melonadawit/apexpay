"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type Product, type Order } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function InventoryPage() {
  const { t } = useLanguage()
  const { checking } = useRequireAuth()
  const { data: products, refetch: refetchProducts } = useData(() => api.inventory.products(), [])
  const { data: orders, refetch: refetchOrders } = useData(() => api.inventory.orders(), [])

  // create product
  const [pName, setPName] = React.useState("")
  const [pPrice, setPPrice] = React.useState("")
  const [pStock, setPStock] = React.useState("10")
  const [err, setErr] = React.useState("")
  const [saving, setSaving] = React.useState(false)

  // create order
  const [selected, setSelected] = React.useState<{ product: Product; qty: number }[]>([])
  const [customer, setCustomer] = React.useState("")

  if (checking) return <Centered>Checking session…</Centered>

  const addProduct = async () => {
    setSaving(true); setErr("")
    try { await api.inventory.createProduct({ name: pName, price: pPrice, stock_qty: pStock, cost_price: "0" }); setPName(""); setPPrice(""); refetchProducts() }
    catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  const addToCart = (p: Product) => {
    setSelected((s) => {
      const found = s.find((x) => x.product.id === p.id)
      if (found) return s.map((x) => (x.product.id === p.id ? { ...x, qty: x.qty + 1 } : x))
      return [...s, { product: p, qty: 1 }]
    })
  }

  const checkout = async () => {
    setSaving(true); setErr("")
    try {
      const items = selected.map((x) => ({ product_id: x.product.id, description: x.product.name, quantity: String(x.qty), unit_price: x.product.price }))
      await api.inventory.createOrder({ customer_name: customer, items })
      setSelected([]); setCustomer(""); refetchProducts(); refetchOrders()
    } catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  const total = selected.reduce((s, x) => s + Number(x.product.price) * x.qty, 0)

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">{t("Inventory & Sales","እቃዎች")}</h1>
        <p className="text-sm text-muted-foreground">Products, stock, and online checkout (software POS) through the gateway.</p>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Add Product</h3>
            <input value={pName} onChange={(e) => setPName(e.target.value)} placeholder="Product name" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={pPrice} onChange={(e) => setPPrice(e.target.value)} placeholder="Price ETB" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={pStock} onChange={(e) => setPStock(e.target.value)} placeholder="Initial stock" className="w-full rounded-xl border h-11 px-3 text-sm" />
            {err && <p className="text-sm text-red-600">{err}</p>}
            <button onClick={addProduct} disabled={saving} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">Add Product</button>
          </div>

          <div className="lg:col-span-2 rounded-2xl border bg-card overflow-hidden">
            <h3 className="font-semibold p-4">Products</h3>
            <div className="grid grid-cols-5 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>Name</span><span>Price</span><span>Stock</span><span>VAT</span><span>Add</span>
            </div>
            {(products ?? []).map((p) => (
              <div key={p.id} className="grid grid-cols-5 gap-2 p-3 border-t text-xs">
                <span className="font-medium">{p.name}</span>
                <span>ETB {p.price}</span>
                <span className={Number(p.stock_qty) <= Number(p.low_stock_threshold) ? "text-red-600 font-medium" : ""}>{p.stock_qty}</span>
                <span>{p.vat_category}</span>
                <span><button onClick={() => addToCart(p)} disabled={Number(p.stock_qty) <= 0} className="text-primary disabled:opacity-40">+ Add</button></span>
              </div>
            ))}
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Checkout Cart (POS)</h3>
            <input value={customer} onChange={(e) => setCustomer(e.target.value)} placeholder="Customer name" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <div className="space-y-1">
              {selected.map((x) => (
                <div key={x.product.id} className="flex justify-between text-sm border-t pt-1">
                  <span>{x.product.name} × {x.qty}</span><span>ETB {(Number(x.product.price) * x.qty).toLocaleString()}</span>
                </div>
              ))}
            </div>
            <p className="text-lg font-bold">Total: ETB {total.toLocaleString()}</p>
            <button onClick={checkout} disabled={saving || selected.length === 0} className="w-full rounded-xl bg-primary text-white h-12 text-sm font-semibold disabled:opacity-50">Checkout • ክፍያ</button>
          </div>

          <div className="rounded-2xl border bg-card overflow-hidden">
            <h3 className="font-semibold p-4">Recent Orders</h3>
            <div className="grid grid-cols-4 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>Order</span><span>Customer</span><span>Total</span><span>Status</span>
            </div>
            {(orders ?? []).slice(0, 10).map((o) => (
              <div key={o.id} className="grid grid-cols-4 gap-2 p-3 border-t text-xs">
                <span className="font-mono text-[10px]">{o.order_number}</span>
                <span>{o.customer_name || "—"}</span>
                <span>ETB {o.total_amount}</span>
                <span className={`text-[11px] px-2 py-0.5 rounded-full w-fit ${o.status === "paid" ? "bg-green-500/15 text-green-700" : "bg-amber-500/15 text-amber-700"}`}>{o.status}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
