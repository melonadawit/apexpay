"use client"
import * as React from "react"
import { useDropzone } from "react-dropzone"
import { motion, AnimatePresence } from "framer-motion"
import { UploadCloud, FileText, Image as ImageIcon, X, CheckCircle } from "lucide-react"
import { cn } from "@/lib/utils"

type FileWithPreview = File & { preview?: string; id: string; status: "pending" | "uploading" | "done" | "error"; progress: number }

export function DocumentDropzone({
  requiredDocs,
  uploaded,
  onFiles,
  accept = { "application/pdf": [".pdf"], "image/jpeg": [".jpg",".jpeg"], "image/png": [".png"] },
  maxSize = 5 * 1024 * 1024, // 5MB, Fayda 2MB override
}: {
  requiredDocs: { type: string; label: string; required: boolean }[]
  uploaded: Record<string, number>
  onFiles: (files: File[], docType: string) => void
  accept?: any
  maxSize?: number
}) {
  const [activeType, setActiveType] = React.useState(requiredDocs[0]?.type || "company_registration")

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    accept,
    maxSize,
    onDrop: (files) => onFiles(files, activeType),
  })

  const progress = Math.round((Object.keys(uploaded).length / requiredDocs.length) * 100)

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-semibold">Documents Vault • ሰነዶች</h3>
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          {Object.keys(uploaded).length}/{requiredDocs.length} • {progress}%
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="space-y-2 max-h-[400px] overflow-auto pr-2">
          {requiredDocs.map((doc) => {
            const done = (uploaded[doc.type] || 0) > 0
            return (
              <button
                key={doc.type}
                onClick={() => setActiveType(doc.type)}
                className={cn(
                  "w-full text-left rounded-xl border p-3 flex items-center gap-3 transition-all",
                  activeType === doc.type ? "border-primary bg-primary/5 shadow-soft" : "border-border hover:bg-muted",
                  done && "border-green-500/20 bg-green-500/10/50"
                )}
              >
                <div className={cn("h-8 w-8 rounded-lg flex items-center justify-center", done ? "bg-green-500/20 text-green-600" : "bg-muted/80")}>
                  {done ? <CheckCircle size={16} /> : <FileText size={16} />}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium truncate">{doc.label}</p>
                  <p className="text-xs text-muted-foreground">{doc.required ? "Required" : "Optional"} {done && "• Uploaded"}</p>
                </div>
              </button>
            )
          })}
        </div>

        <div className="md:col-span-2">
          <motion.div
            {...(getRootProps() as any)}
            className={cn(
              "relative rounded-2xl border-2 border-dashed p-8 text-center cursor-pointer transition-all",
              isDragActive ? "border-primary bg-primary/5 scale-[0.98] pulseGlow" : "border-border hover:border-primary/50 hover:bg-primary/[0.02]",
            )}
            whileTap={{ scale: 0.98 }}
          >
            <input {...getInputProps()} />
            <div className="mx-auto w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center mb-3">
              <UploadCloud className="text-primary" />
            </div>
            <p className="font-medium">Drop {activeType.replaceAll("_"," ")} here • እዚህ ጣል ያድርጉ</p>
            <p className="text-xs text-muted-foreground mt-1">PDF, JPG, PNG up to {maxSize/1024/1024}MB • Fayda images &lt;2MB per NIDP</p>

            <AnimatePresence>
              {isDragActive && (
                <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }} className="absolute inset-0 rounded-2xl bg-primary/10 backdrop-blur-sm flex items-center justify-center">
                  <p className="font-semibold text-primary">Drop now • አሁን ጣል ያድርጉ</p>
                </motion.div>
              )}
            </AnimatePresence>
          </motion.div>

          <div className="mt-4 rounded-xl bg-muted p-3 text-xs text-muted-foreground">
            <p>✓ Files encrypted at rest AES-256 SSE-S3 • ✓ Hash integrity sha256 • ✓ Presigned URLs 15m TTL • ✓ FIN never logged</p>
          </div>
        </div>
      </div>
    </div>
  )
}
