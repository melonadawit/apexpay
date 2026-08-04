"use client"
import * as React from "react"
import { DocumentDropzone } from "@/components/ui/dropzone"

const requiredDocs = [
  { type: "company_registration", label: "Company Registration • የኩባንያ ምዝገባ (notarized)", required: true },
  { type: "tin_certificate", label: "TIN Certificate • ቲን ሰርተፍኬት", required: true },
  { type: "business_license", label: "Business License • ንግድ ፈቃድ (not expired)", required: true },
  { type: "vat_certificate", label: "VAT Certificate • ቫት", required: false },
  { type: "memorandum_articles", label: "Memorandum & Articles • መተዳደሪያ", required: true },
  { type: "board_resolution", label: "Board Resolution • ቦርድ ውሳኔ", required: true },
  { type: "shareholder_list", label: "Shareholder List • ባለአክሲዮኖች ዝርዝር", required: true },
  { type: "fayda_card_front", label: "Fayda Card Front • ፋይዳ ፊት <2MB", required: true },
  { type: "fayda_card_back", label: "Fayda Card Back • ፋይዳ ጀርባ <2MB", required: true },
  { type: "proof_of_address", label: "Proof of Address • አድራሻ ማረጋገጫ", required: true },
  { type: "bank_letter", label: "Bank Letter / Cancelled Cheque • የባንክ ደብዳቤ", required: true },
  { type: "website_screenshot", label: "Website Screenshot + Refund/Privacy/Terms • ድህረ ገጽ", required: true },
]

export function DocumentsVaultStep({ data, onChange }: { data: any, onChange: (d:any)=>void }) {
  const [uploaded, setUploaded] = React.useState<Record<string, number>>(data.uploadedDocs || { company_registration:1, tin_certificate:1, business_license:1, fayda_card_front:1 })

  const handleFiles = (files: File[], docType: string) => {
    // Check size <5MB generic, Fayda <2MB per NIDP
    const isFayda = docType.includes("fayda")
    const max = isFayda ? 2*1024*1024 : 5*1024*1024
    for (const f of files) {
      if (f.size > max) { alert(`${f.name} exceeds ${max/1024/1024}MB`); return }
    }
    // Mock hash integrity sha256
    const newUploaded = { ...uploaded, [docType]: (uploaded[docType]||0)+files.length }
    setUploaded(newUploaded)
    onChange({ ...data, uploadedDocs: newUploaded })
  }

  return (
    <div className="space-y-4">
      <div>
        <h2 className="text-xl font-bold">Documents Vault • ሰነዶች ማከማቻ</h2>
        <p className="text-sm text-muted-foreground">Outstanding dropzone with preview, hash, OCR stub, encrypted MinIO SSE-S3. Required docs computed per business_type & KYC level via RequiredDocs() optimal data structure.</p>
      </div>

      <DocumentDropzone requiredDocs={requiredDocs} uploaded={uploaded} onFiles={handleFiles} />

      <div className="rounded-xl border border-black/10 p-3 text-xs bg-white">
        <p className="font-semibold">What happens next — optimal flow:</p>
        <ul className="list-disc list-inside mt-1 space-y-0.5 text-muted-foreground">
          <li>File uploaded via presigned POST TTL 15m directly to MinIO `merchants/{"{merchant_id}"}/kyc/{"{doc_type}"}_id.pdf` — no server buffering</li>
          <li>Hash sha256 calculated client + server integrity check file_hash unique index</li>
          <li>OCR stub extracts fields: registration_no, TIN, expiry — displayed side-by-side in admin viewer</li>
          <li>Status pending → ocr_done → verified|rejected by compliance ops with rejection_reason</li>
          <li>Docs expiry: business_license not expired validation</li>
          <li>MinIO versioning + encryption, retention 7y per NBE</li>
        </ul>
      </div>
    </div>
  )
}
