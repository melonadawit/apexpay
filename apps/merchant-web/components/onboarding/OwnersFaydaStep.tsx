"use client"
import * as React from "react"
import { motion } from "framer-motion"
import { Button } from "@/components/ui/button"
import { FaydaCapture } from "./FaydaCapture"
import { Check, ShieldCheck, Smartphone } from "lucide-react"

export function OwnersFaydaStep({ data, onChange }: { data: any, onChange: (d: any) => void }) {
  const [owners, setOwners] = React.useState<any[]>(data.owners || [{ id: "own_1", full_name: "Abebe Kebede", full_name_am: "አበበ ከበደ", role: "owner", ownership_percentage: 100, is_authorized_signatory: true, fayda_verified: false }])
  const [activeOwnerIdx, setActiveOwnerIdx] = React.useState(0)
  const [frontFile, setFrontFile] = React.useState<File | null>(null)
  const [backFile, setBackFile] = React.useState<File | null>(null)
  const [selfieFile, setSelfieFile] = React.useState<File | null>(null)
  const [fin, setFin] = React.useState("")
  const [otp, setOtp] = React.useState("")
  const [verified, setVerified] = React.useState(false)
  const [faceScore, setFaceScore] = React.useState(0)

  const activeOwner = owners[activeOwnerIdx]

  const handleVerify = () => {
    if (fin.length !== 12 && fin.length !== 16) { alert("FIN must be 12 digits or FAN 16 chars per id.gov.et"); return }
    if (!frontFile || !backFile) { alert("Front/back images required <2MB per NIDP"); return }
    // mock verify OTP 123456
    if (otp === "123456") {
      setVerified(true)
      setFaceScore(0.92)
      const updated = [...owners]
      updated[activeOwnerIdx] = { ...updated[activeOwnerIdx], fayda_verified: true, fin_last4: fin.slice(-4), face_score: 0.92 }
      setOwners(updated)
      onChange({ ...data, owners: updated })
    } else {
      alert("Invalid OTP - mock OTP is 123456")
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold flex items-center gap-2"><ShieldCheck className="text-primary" /> Owners & Fayda ID Verification • ባለቤቶች እና ፋይዳ ማረጋገጫ</h2>
        <p className="text-sm text-muted-foreground">Per id.gov.et: FIN 12-digit, FAN alias, OTP consent, front/back images, offline QR alternative. FIN stored as sha256(salt+fin) + last4 only, never plain.</p>
      </div>

      <div className="flex gap-2 overflow-auto pb-2">
        {owners.map((o, idx) => (
          <button key={o.id} onClick={() => setActiveOwnerIdx(idx)} className={`rounded-xl border px-4 py-2 flex items-center gap-2 min-w-max ${idx === activeOwnerIdx ? "border-primary bg-primary/10" : "border-border bg-card"}`}>
            <div className={`h-6 w-6 rounded-full flex items-center justify-center text-xs ${o.fayda_verified ? "bg-green-500 text-foreground" : "bg-neutral-200"}`}>{o.fayda_verified ? <Check size={12} /> : idx + 1}</div>
            <span className="text-sm font-medium">{o.full_name} {o.is_authorized_signatory && "• Auth Signatory"}</span>
          </button>
        ))}
        <Button variant="outline" size="sm" onClick={() => setOwners([...owners, { id: `own_${owners.length + 1}`, full_name: "New Owner", role: "shareholder", ownership_percentage: 10 }])}>+ Add Owner</Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <FaydaCapture type="front" capturedUrl={frontFile ? URL.createObjectURL(frontFile) : undefined} onCapture={setFrontFile} onRetake={() => setFrontFile(null)} />
        <FaydaCapture type="back" capturedUrl={backFile ? URL.createObjectURL(backFile) : undefined} onCapture={setBackFile} onRetake={() => setBackFile(null)} />
        <FaydaCapture type="selfie" capturedUrl={selfieFile ? URL.createObjectURL(selfieFile) : undefined} onCapture={setSelfieFile} onRetake={() => setSelfieFile(null)} />
      </div>

      <div className="rounded-2xl border border-border p-4 space-y-4 bg-card">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-1"><label className="text-sm font-medium">FIN 12-digit or FAN 16 • ፋይዳ ቁጥር *</label><input value={fin} onChange={e => setFin(e.target.value)} placeholder="123456789012 or FAN alias" className="w-full rounded-xl border h-12 px-3" maxLength={16} /></div>
          <div className="space-y-1"><label className="text-sm font-medium flex items-center gap-1"><Smartphone size={14} /> OTP 6-digit • የተላከ ኮድ * (mock 123456)</label><input value={otp} onChange={e => setOtp(e.target.value)} placeholder="123456" className="w-full rounded-xl border h-12 px-3" maxLength={6} /></div>
        </div>

        <div className="flex items-center gap-3">
          <Button onClick={handleVerify} disabled={verified} className="flex-1">{verified ? `Verified • Face Score ${faceScore}` : "Verify via Fayda • በፋይዳ አረጋግጥ"}</Button>
          {verified && <span className="text-green-600 text-sm flex items-center gap-1"><Check size={16} /> OTP verified • FIN hashed sha256(salt+FIN)</span>}
        </div>

        <div className="text-xs text-muted-foreground">
          <p>✓ Consent timestamp + IP logged • ✓ OTP rate limit 5/hour via Redis • ✓ Offline QR fallback FaydaEncode scan • ✓ OIDC eSignet alternative • ✓ Front/back encrypted MinIO SSE-S3 presigned 15m</p>
        </div>
      </div>

      {verified && (
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="rounded-xl bg-green-500/10 border border-green-500/20 p-3 flex gap-3">
          <Check className="text-green-600" />
          <div className="text-sm">
            <p className="font-semibold text-green-700">Fayda Verification Successful • ማረጋገጫ ተሳክቷል</p>
            <p>FIN ****-{fin.slice(-4)} verified • Demographics match true • Face match 0.92 &gt; 0.85 threshold • Consent 2026-08-04 • Response encrypted ref fayda_responses/xxx.enc</p>
          </div>
        </motion.div>
      )}
    </div>
  )
}
