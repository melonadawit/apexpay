"use client"
import * as React from "react"
import { motion } from "framer-motion"

const merchantsMock = [
  { id:"mer_01H", legal:"Apex Trading PLC", tin:"0023456789", status:"submitted", risk:42, fayda_verified:true, bank_verified:true, docs: "4/6" },
  { id:"mer_02H", legal:"Buna Coffee Export", tin:"0012345678", status:"fayda_pending", risk:78, fayda_verified:false, bank_verified:true, docs: "2/6" },
  { id:"mer_03H", legal:"Habesha Tech", tin:"0034567890", status:"compliance_check", risk:15, fayda_verified:true, bank_verified:true, docs: "6/6" },
]

export default function AdminPage() {
  const [selected, setSelected] = React.useState(merchantsMock[0])

  return (
    <div className="min-h-screen bg-neutral-50 p-6">
      <header className="mb-6 flex justify-between items-center">
        <h1 className="text-2xl font-bold">Admin — Onboarding Queue • ማስተዳደር</h1>
        <div className="text-sm text-muted-foreground">NBE Exam Console • 30-60 days pilot analogy</div>
      </header>

      <div className="grid grid-cols-12 gap-6">
        {/* Kanban board outstanding */}
        <div className="col-span-4 space-y-3">
          <h3 className="font-semibold">Queue • ወረፋ — {merchantsMock.length}</h3>
          {["submitted","fayda_pending","compliance_check","pending_approval","active"].map(status=> (
            <div key={status} className="rounded-xl border bg-white p-3">
              <p className="text-xs font-bold uppercase text-muted-foreground">{status.replace("_"," ")} • {[merchantsMock.filter(m=>m.status===status).length]}</p>
              <div className="mt-2 space-y-2">
                {merchantsMock.filter(m=>m.status===status).map(m=> (
                  <motion.button key={m.id} onClick={()=> setSelected(m)} whileHover={{ scale:1.01 }} className={`w-full text-left rounded-xl border p-3 ${selected.id===m.id ? "border-primary bg-primary/5" : "border-black/10 bg-white"}`}>
                    <p className="text-sm font-semibold">{m.legal}</p>
                    <p className="text-xs text-muted-foreground">TIN {m.tin} • Risk {m.risk}</p>
                    <div className="mt-1 flex gap-1">
                      <span className={`text-[10px] px-2 py-0.5 rounded-full ${m.fayda_verified ? "bg-green-100 text-green-700" : "bg-amber-100 text-amber-700"}`}>Fayda {m.fayda_verified ? "✓" : "pending"}</span>
                      <span className="text-[10px] px-2 py-0.5 rounded-full bg-neutral-100">{m.docs} docs</span>
                    </div>
                  </motion.button>
                ))}
              </div>
            </div>
          ))}
        </div>

        {/* Exam view split pane */}
        <div className="col-span-8 grid grid-cols-2 gap-4">
          <div className="rounded-xl border bg-white p-4 space-y-4">
            <h4 className="font-semibold">KYC Profile • {selected.legal}</h4>
            <div className="space-y-2 text-sm">
              <p><span className="text-muted-foreground">Legal:</span> {selected.legal}</p>
              <p><span className="text-muted-foreground">TIN:</span> {selected.tin}</p>
              <p><span className="text-muted-foreground">Risk:</span> <span className={`px-2 py-0.5 rounded-full text-xs ${selected.risk<50 ? "bg-green-100" : "bg-amber-100"}`}>{selected.risk}/100 Medium</span></p>
            </div>
            <div className="rounded-xl bg-amber-50 border border-amber-200 p-3 text-xs">
              <p className="font-semibold">Fayda Verification Chain</p>
              <p>FIN ****-****-1234 • OTP verified • Face 0.92 • Demographics match • Consent 2026-08-04 • Response encrypted ref fayda_responses/xxx.enc</p>
              <div className="mt-2 grid grid-cols-3 gap-2">
                <div className="rounded-lg bg-black/5 h-24 flex items-center justify-center text-[10px]">Front blurred — click to view auth</div>
                <div className="rounded-lg bg-black/5 h-24 flex items-center justify-center text-[10px]">Back blurred</div>
                <div className="rounded-lg bg-black/5 h-24 flex items-center justify-center text-[10px]">Selfie</div>
              </div>
            </div>
            <div className="rounded-xl border p-3 text-xs">
              <p className="font-semibold">Bank Settlement</p>
              <p>CBE • ****1234 • Account name == legal fuzzy match distance 1 ✓ Verified</p>
            </div>
          </div>

          <div className="rounded-xl border bg-white p-4 space-y-3">
            <h4 className="font-semibold">Compliance Checks • 8</h4>
            {[
              ["tin_validation","passed",95],
              ["business_license_validation","passed",90],
              ["bank_account_validation","passed",88],
              ["fayda_verification","passed",92],
              ["restricted_industry","passed",100],
              ["website_policy_check","needs_review",70],
              ["aml_screening","passed",85],
              ["risk_scoring","passed",42],
            ].map(([type,status,score])=> (
              <div key={type as string} className="flex items-center justify-between text-sm border-b last:border-0 py-2">
                <span>{type as string}</span>
                <span className={`text-xs px-2 py-0.5 rounded-full ${status==="passed" ? "bg-green-100" : "bg-amber-100"}`}>{status as string}</span>
                <span className="text-xs text-muted-foreground">{score as number}</span>
              </div>
            ))}
            <div className="flex gap-2 pt-2">
              <button className="flex-1 rounded-xl border h-10 text-sm">Request Info • መረጃ ጠይቅ</button>
              <button className="flex-1 rounded-xl bg-green-600 text-white h-10 text-sm">Approve • አጽድቅ</button>
              <button className="flex-1 rounded-xl bg-red-600 text-white h-10 text-sm">Reject • ውድቅ</button>
            </div>
            <p className="text-[11px] text-muted-foreground">Dual approval required if risk>=70 or TPV&gt;1M ETB. maker-checker enforced. Outbox merchant.activated triggers email confetti.</p>
          </div>

          <div className="col-span-2 rounded-xl border bg-white p-4">
            <h4 className="font-semibold">Onboarding Reviews Timeline • የጊዜ መስመር</h4>
            <div className="mt-3 relative pl-6 border-l-2 border-neutral-200 space-y-3">
              {[
                "2026-08-04 09:00 system submit risk_score 42",
                "2026-08-04 09:05 ops request_info missing refund_policy_doc",
                "2026-08-04 09:20 system Fayda verified face 0.92 OTP",
                "2026-08-04 09:30 compliance approve first approval needs second (TPV 500k medium)",
              ].map((t,i)=> (
                <div key={i} className="relative text-xs">
                  <div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-primary border-2 border-white" />
                  <p>{t}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
