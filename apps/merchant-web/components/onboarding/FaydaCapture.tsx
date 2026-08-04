"use client"
import * as React from "react"
import { motion, AnimatePresence } from "framer-motion"
import { Camera, Check, AlertTriangle, RotateCcw } from "lucide-react"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

type CaptureType = "front" | "back" | "selfie"

export function FaydaCapture({
  type,
  onCapture,
  onRetake,
  capturedUrl,
}: {
  type: CaptureType
  onCapture: (file: File) => void
  onRetake?: () => void
  capturedUrl?: string
}) {
  const videoRef = React.useRef<HTMLVideoElement>(null)
  const canvasRef = React.useRef<HTMLCanvasElement>(null)
  const [stream, setStream] = React.useState<MediaStream | null>(null)
  const [hasGlare, setHasGlare] = React.useState(false)
  const [brightness, setBrightness] = React.useState(0)

  const typeConfig = {
    front: { label: "Fayda Front • የፊት ገጽታ", hint: "Align card within corners • ካርዱን በማዕዘኖቹ ውስጥ ያስቀምጡ", icon: "🪪" },
    back: { label: "Fayda Back • የኋላ ገጽታ", hint: "Capture back with QR visible", icon: "🪪" },
    selfie: { label: "Selfie Liveness • የፊት ፎቶ", hint: "Look straight, remove hat • በቀጥታ ይመልከቱ", icon: "🤳" },
  }[type]

  React.useEffect(() => {
    async function startCamera() {
      try {
        const s = await navigator.mediaDevices.getUserMedia({
          video: { facingMode: type === "selfie" ? "user" : "environment", width: 1280, height: 720 },
        })
        setStream(s)
        if (videoRef.current) videoRef.current.srcObject = s
      } catch (e) {
        console.error("camera error", e)
      }
    }
    if (!capturedUrl) startCamera()
    return () => stream?.getTracks().forEach((t) => t.stop())
  }, [capturedUrl])

  // Glare detection loop - optimal sampling every 400ms
  React.useEffect(() => {
    if (!videoRef.current || !canvasRef.current || capturedUrl) return
    const interval = setInterval(() => {
      const video = videoRef.current!, canvas = canvasRef.current!, ctx = canvas.getContext("2d")!
      canvas.width = 100
      canvas.height = 100
      ctx.drawImage(video, 0, 0, 100, 100)
      const imageData = ctx.getImageData(0, 0, 100, 100)
      let total = 0, bright = 0
      for (let i = 0; i < imageData.data.length; i += 16) {
        const r = imageData.data[i], g = imageData.data[i + 1], b = imageData.data[i + 2]
        const bness = (r + g + b) / 3
        total += bness
        if (bness > 200) bright++
      }
      const avg = total / (imageData.data.length / 16)
      setBrightness(Math.round(avg))
      setHasGlare(bright / (imageData.data.length / 16) > 0.15)
    }, 400)
    return () => clearInterval(interval)
  }, [capturedUrl, stream])

  const capture = () => {
    const video = videoRef.current!, canvas = document.createElement("canvas")
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight
    const ctx = canvas.getContext("2d")!
    if (type === "selfie") { ctx.scale(-1, 1); ctx.drawImage(video, -canvas.width, 0) } else ctx.drawImage(video, 0, 0)
    canvas.toBlob((blob) => {
      if (blob) {
        const file = new File([blob], `fayda_${type}_${Date.now()}.jpg`, { type: "image/jpeg" })
        if (file.size > 2 * 1024 * 1024) alert("File must be <2MB per NIDP spec")
        else onCapture(file)
      }
    }, "image/jpeg", 0.85)
  }

  if (capturedUrl) {
    return (
      <div className="relative rounded-2xl overflow-hidden border border-black/10">
        <img src={capturedUrl} alt={`Fayda ${type}`} className="w-full h-64 object-cover" />
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
        <div className="absolute bottom-3 left-3 right-3 flex justify-between">
          <span className="text-white text-sm font-medium flex items-center gap-2"><Check size={16} className="bg-green-500 rounded-full p-0.5" /> Captured</span>
          <Button size="sm" variant="glass" onClick={onRetake}><RotateCcw size={14} className="mr-1" /> Retake</Button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center gap-3">
        <div className="h-10 w-10 rounded-xl bg-primary/10 flex items-center justify-center text-lg">{typeConfig.icon}</div>
        <div><p className="font-semibold text-sm">{typeConfig.label}</p><p className="text-xs text-muted-foreground">{typeConfig.hint}</p></div>
      </div>

      <div className="relative rounded-2xl overflow-hidden bg-black aspect-[4/3] border-2 border-black/10">
        <video ref={videoRef} autoPlay playsInline muted className={cn("w-full h-full object-cover", type === "selfie" && "-scale-x-100")} />

        {/* Outstanding corner guides animated pulse */}
        <div className="absolute inset-6 pointer-events-none">
          <div className="absolute top-0 left-0 w-8 h-8 border-l-4 border-t-4 border-white rounded-tl-xl animate-pulse" />
          <div className="absolute top-0 right-0 w-8 h-8 border-r-4 border-t-4 border-white rounded-tr-xl animate-pulse" />
          <div className="absolute bottom-0 left-0 w-8 h-8 border-l-4 border-b-4 border-white rounded-bl-xl animate-pulse" />
          <div className="absolute bottom-0 right-0 w-8 h-8 border-r-4 border-b-4 border-white rounded-br-xl animate-pulse" />
        </div>

        <canvas ref={canvasRef} className="hidden" />

        {/* Glare / brightness outstanding helper */}
        <AnimatePresence>
          {hasGlare && (
            <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0 }} className="absolute top-3 left-1/2 -translate-x-1/2 bg-amber-500 text-white text-xs px-3 py-1.5 rounded-full flex items-center gap-1 shadow-medium">
              <AlertTriangle size={12} /> Move to shade • ጥላ ውስጥ ይሂዱ (brightness {brightness})
            </motion.div>
          )}
        </AnimatePresence>

        <div className="absolute bottom-0 left-0 right-0 p-4 bg-gradient-to-t from-black/80 to-transparent flex justify-center">
          <Button onClick={capture} size="lg" className="rounded-full h-16 w-16 p-0 shadow-medium">
            <Camera />
          </Button>
        </div>
      </div>

      <p className="text-[11px] text-muted-foreground">✓ Encrypted at rest • ✓ Hash integrity • ✓ FIN never logged • File &lt;2MB per id.gov.et</p>
    </div>
  )
}
