"use client"
import * as React from "react"

const banks = [
  { code: "CBE", name: "Commercial Bank of Ethiopia", logo: "🏦" },
  { code: "AWASH", name: "Awash Bank", logo: "🏦" },
  { code: "DASHEN", name: "Dashen Bank", logo: "🏦" },
  { code: "ABYSSINIA", name: "Bank of Abyssinia", logo: "🏦" },
  { code: "BERHAN", name: "Berhan Bank", logo: "🏦" },
  { code: "WEGAGEN", name: "Wegagen Bank", logo: "🏦" },
]

export function BankAccountStep({ data, onChange }: { data: any, onChange: (a: any) => void }) {
  const [bankCode, setBankCode] = React.useState(data.bank_code || "CBE")
  const [acctName, setAcctName] = React.useState(data.account_name || data.legal_name || "")
  const [acctNo, setAcctNo] = React.useState(data.account_number || "")
  const [isDefault, setIsDefault] = React.useState(true)

  React.useEffect(() => onChange({ ...data, bank_code: bankCode, account_name: acctName, account_number: acctNo, is_settlement_default: isDefault }), [bankCode, acctName, acctNo, isDefault])

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold">Settlement Bank Account • የክፍያ ሂሳብ</h2>
        <p className="text-sm text-muted-foreground">Per PayAtlas: Bank account must match legal name, bank letter/cancelled cheque required, account hash stored not plain, masked ****1234 display.</p>
      </div>

      <div className="space-y-2">
        <label className="text-sm font-medium">Bank • ባንክ * — via GET /v1/banks</label>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
          {banks.map(b => (
            <button key={b.code} onClick={() => setBankCode(b.code)} className={`rounded-xl border p-3 flex items-center gap-3 text-left ${bankCode === b.code ? "border-primary bg-primary/5" : "border-border bg-card hover:bg-muted"}`}>
              <span className="text-xl">{b.logo}</span><span className="text-sm font-medium">{b.name} • {b.code}</span>{bankCode === b.code && <span className="ml-auto text-primary">✓</span>}
            </button>
          ))}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-1"><label className="text-sm font-medium">Account Name must == Legal Name • ሂሳብ ስም *</label><input value={acctName} onChange={e => setAcctName(e.target.value)} placeholder="Apex Trading PLC" className="w-full rounded-xl border h-12 px-3" /></div>
        <div className="space-y-1"><label className="text-sm font-medium">Account Number • ሂሳብ ቁጥር *</label><input value={acctNo} onChange={e => setAcctNo(e.target.value)} placeholder="1000123456789" className="w-full rounded-xl border h-12 px-3" /><p className="text-xs text-muted-foreground">Stored as hash sha256 + masked ****{acctNo.slice(-4)}</p></div>
        <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={isDefault} onChange={e => setIsDefault(e.target.checked)} /> Settlement default • ነባሪ ሂሳብ</label>
      </div>

      <div className="rounded-xl bg-blue-500/10 border border-blue-500/20 p-3 text-xs">
        <p className="font-semibold">Validation:</p>
        <ul className="list-disc list-inside mt-1">
          <li>Account name fuzzy match Levenshtein &lt;3 vs legal_name or require override note</li>
          <li>Bank verification method: bank_letter / micro_deposit / manual (ops)</li>
          <li>At least one settlement default required to submit per NBE</li>
          <li>Account hash stored, full number encrypted in MinIO if needed, masked display ****1234 outstanding UI</li>
        </ul>
      </div>
    </div>
  )
}
