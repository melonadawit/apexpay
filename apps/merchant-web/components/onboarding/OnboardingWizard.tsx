"use client"
import * as React from "react"
import { motion, AnimatePresence } from "framer-motion"
import { Check, ChevronRight, ChevronLeft, Sparkles } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, GlassCard } from "@/components/ui/card"
import { Progress, DonutProgress } from "@/components/ui/progress"
import { cn } from "@/lib/utils"
import { BusinessInfoStep } from "./BusinessInfoStep"
import { OwnersFaydaStep } from "./OwnersFaydaStep"
import { BankAccountStep } from "./BankAccountStep"
import { DocumentsVaultStep } from "./DocumentsVaultStep"
import { CompliancePreviewStep } from "./CompliancePreviewStep"
import { ReviewSubmitStep } from "./ReviewSubmitStep"

const steps = [
  { id: "business", label: "Business", labelAm: "ንግድ", desc: "Legal & TIN" },
  { id: "owners", label: "Owners & Fayda", labelAm: "ባለቤቶች እና ፋይዳ", desc: "Front/Back + OTP" },
  { id: "bank", label: "Bank", labelAm: "ባንክ", desc: "Settlement" },
  { id: "docs", label: "Documents", labelAm: "ሰነዶች", desc: "Vault" },
  { id: "compliance", label: "Compliance", labelAm: "ተገዢነት", desc: "Risk & Checks" },
  { id: "review", label: "Review", labelAm: "ግምገማ", desc: "Submit" },
]

export function OnboardingWizard() {
  const [current, setCurrent] = React.useState(0)
  const [completed, setCompleted] = React.useState<boolean[]>(Array(6).fill(false))
  const [formData, setFormData] = React.useState<any>({})

  const progress = Math.round((completed.filter(Boolean).length / steps.length) * 100)

  const next = () => {
    const newCompleted = [...completed]
    newCompleted[current] = true
    setCompleted(newCompleted)
    if (current < steps.length - 1) setCurrent(current + 1)
  }
  const back = () => { if (current > 0) setCurrent(current - 1) }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50 p-4 md:p-8">
      <div className="max-w-6xl mx-auto">
        {/* Header outstanding like Mercury + Linear */}
        <div className="mb-8 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Merchant Onboarding • የነጋዴ ምዝገባ</h1>
            <p className="text-muted-foreground">NBE PSO Gateway Operator per ONPS/02/2020 + Fayda ID verification id.gov.et</p>
          </div>
          <div className="flex items-center gap-4">
            <DonutProgress value={progress} />
            <div className="text-sm">
              <p className="font-semibold">{progress}% Complete</p>
              <p className="text-muted-foreground">{completed.filter(Boolean).length}/{steps.length} steps</p>
            </div>
          </div>
        </div>

        {/* Stepper outstanding with animated line */}
        <div className="mb-8 overflow-x-auto">
          <div className="flex items-center gap-2 min-w-max">
            {steps.map((s, i) => {
              const isDone = completed[i]
              const isCurrent = i === current
              const isPast = i < current
              return (
                <React.Fragment key={s.id}>
                  <motion.button
                    onClick={() => isPast && setCurrent(i)}
                    className={cn(
                      "flex items-center gap-3 rounded-xl border px-4 py-3 text-left transition-all min-w-[160px]",
                      isCurrent ? "border-primary bg-primary text-white shadow-soft scale-[1.02]" : "",
                      isDone && !isCurrent ? "border-green-200 bg-green-50 text-green-700" : "",
                      !isCurrent && !isDone ? "border-black/10 bg-white hover:bg-neutral-50" : ""
                    )}
                    whileHover={{ scale: isPast ? 1.02 : 1 }}
                    whileTap={{ scale: 0.98 }}
                  >
                    <div className={cn("h-8 w-8 rounded-full flex items-center justify-center text-sm font-bold", isCurrent ? "bg-white text-primary" : isDone ? "bg-green-500 text-white" : "bg-neutral-100")}>
                      {isDone ? <Check size={16} /> : i + 1}
                    </div>
                    <div>
                      <p className="text-sm font-semibold leading-none">{s.label}</p>
                      <p className="text-[11px] opacity-80">{s.labelAm} • {s.desc}</p>
                    </div>
                  </motion.button>
                  {i < steps.length - 1 && (
                    <div className="h-0.5 w-8 bg-neutral-200 relative overflow-hidden">
                      <motion.div className="absolute inset-0 bg-primary" initial={{ scaleX: 0 }} animate={{ scaleX: isDone || isPast ? 1 : 0 }} style={{ originX: 0 }} transition={{ duration: 0.5 }} />
                    </div>
                  )}
                </React.Fragment>
              )
            })}
          </div>
        </div>

        {/* Main content outstanding card with motion */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2">
            <Card className="p-0 overflow-hidden">
              <div className="h-1 w-full bg-neutral-100">
                <motion.div className="h-full bg-primary" initial={{ width: 0 }} animate={{ width: `${((current + 1) / steps.length) * 100}%` }} transition={{ duration: 0.7, ease: [0.22, 1, 0.36, 1] }} />
              </div>
              <div className="p-6 md:p-8">
                <AnimatePresence mode="wait">
                  <motion.div key={current} initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }} exit={{ opacity: 0, x: -20 }} transition={{ duration: 0.3, ease: [0.22, 1, 0.36, 1] }}>
                    {current === 0 && <BusinessInfoStep data={formData} onChange={setFormData} />}
                    {current === 1 && <OwnersFaydaStep data={formData} onChange={setFormData} />}
                    {current === 2 && <BankAccountStep data={formData} onChange={setFormData} />}
                    {current === 3 && <DocumentsVaultStep data={formData} onChange={setFormData} />}
                    {current === 4 && <CompliancePreviewStep data={formData} />}
                    {current === 5 && <ReviewSubmitStep data={formData} />}
                  </motion.div>
                </AnimatePresence>

                <div className="mt-8 flex justify-between">
                  <Button variant="outline" onClick={back} disabled={current === 0}><ChevronLeft size={16} className="mr-1" /> Back • ተመለስ</Button>
                  {current < steps.length - 1 ? (
                    <Button onClick={next}>Next • ቀጣይ <ChevronRight size={16} className="ml-1" /></Button>
                  ) : (
                    <Button variant="gold" onClick={() => alert("Submit to /v1/onboarding/submit")}><Sparkles size={16} className="mr-2" /> Submit • አስገባ</Button>
                  )}
                </div>
              </div>
            </Card>
          </div>

          {/* Sidebar outstanding - timeline + checklist */}
          <div className="space-y-4">
            <GlassCard className="p-4">
              <h4 className="font-semibold mb-3">Progress • እድገት</h4>
              <Progress value={progress} className="mb-3" />
              <div className="space-y-2 text-sm">
                {steps.map((s, i) => (
                  <div key={s.id} className="flex items-center gap-2">
                    <div className={cn("h-2 w-2 rounded-full", i === current ? "bg-primary animate-pulse" : completed[i] ? "bg-green-500" : "bg-neutral-300")} />
                    <span className={cn(i === current && "font-semibold")}>{s.label} • {s.labelAm}</span>
                    {completed[i] && <Check size={12} className="text-green-500 ml-auto" />}
                  </div>
                ))}
              </div>
            </GlassCard>

            <Card className="p-4 bg-gradient-to-br from-primary to-primary-light text-white">
              <h4 className="font-semibold">NBE Checklist • መስፈርቶች</h4>
              <ul className="mt-2 space-y-1 text-xs opacity-90 list-disc list-inside">
                <li>Company Registration notarized</li>
                <li>TIN 10-digit + Business License not expired</li>
                <li>Fayda front/back &lt;2MB + selfie + OTP consent id.gov.et</li>
                <li>Bank letter account name == legal name</li>
                <li>Website refund/privacy/terms per PayAtlas</li>
                <li>Min 1 authorized signatory Fayda verified</li>
              </ul>
            </Card>

            <Card className="p-4">
              <h4 className="font-semibold mb-2">Timeline • የጊዜ መስመር</h4>
              <div className="relative pl-6 border-l-2 border-neutral-200 space-y-4">
                {["Business info saved • draft", "Fayda verification OTP sent • 123456 mock", "Bank added • CBE ****1234", "Docs 4/6 uploaded • 66%"].map((t, i) => (
                  <div key={i} className="relative">
                    <div className="absolute -left-[29px] top-0 h-3 w-3 rounded-full bg-primary border-2 border-white" />
                    <p className="text-xs">{t}</p>
                    <p className="text-[11px] text-muted-foreground">{i + 1} min ago</p>
                  </div>
                ))}
              </div>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}
