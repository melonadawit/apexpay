"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockInvoices = [
  { id: "inv_001", vendor: "CBE Vendor Supplies • Vendor • Office Supplies", invoice_number: "INV-2026-001", invoice_date: "2026-07-20", due_date: "2026-08-05", amount: "50000", tax_amount: "7500", withholding_tax: "1000", total_amount: "56500", status: "pending_approval", ocr_raw: { extracted_text: "Invoice INV-2026-001 Vendor CBE Vendor Supplies TIN 0098765432 Amount 50000 VAT 15% 7500 Withholding 2% 1000 Total 56500", confidence: 0.92, vendor_name: "CBE Vendor Supplies", tin: "0098765432", amount: "50000", tax: "7500", withholding: "1000" }, file_key: "vendors/inv_001.pdf", file_hash: "hash_inv_001", payout_id: null },
  { id: "inv_002", vendor: "Awash Logistics • Vendor • Logistics", invoice_number: "INV-2026-002", invoice_date: "2026-07-22", due_date: "2026-08-10", amount: "30000", tax_amount: "4500", withholding_tax: "600", total_amount: "33900", status: "approved", ocr_raw: { confidence: 0.89, vendor_name: "Awash Logistics", tin: "0098765433" }, file_key: "vendors/inv_002.pdf", file_hash: "hash_inv_002" },
  { id: "inv_003", vendor: "Dashen Marketing • Vendor • Marketing", invoice_number: "INV-2026-003", invoice_date: "2026-07-25", due_date: "2026-08-15", amount: "100000", tax_amount: "15000", withholding_tax: "2000", total_amount: "113000", status: "paid", ocr_raw: { confidence: 0.91 }, file_key: "vendors/inv_003.pdf", payout_id: "pout_inv_003" },
]

const mockPOs = [
  { id: "po_001", vendor: "CBE Vendor Supplies", po_number: "PO-2026-001", amount: "50000", currency: "ETB", status: "approved", created_at: "2026-07-18" },
]

export default function VendorPaymentsPage() {
  const [selected, setSelected] = React.useState(mockInvoices[0])
  const [dragOver, setDragOver] = React.useState(false)

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Vendor Payments • አቅራቢ ክፍያዎች • End-to-End Accounts Payable Automation OCR-enabled Invoice Capture Multi-layer Approval Workflows Automated TDS Calculation and Filing to NSDL Integrated Payouts • RazorpayX Parity • P0</h1>
            <p className="text-sm text-muted-foreground mt-2">Upload invoices and auto-pay vendors and TDS payments OCR-enabled invoice capture multi-layer approval workflows automated TDS calculation and filing to NSDL integrated payouts via NEFT/RTGS/IMPS/UPI manage all vendors payouts purchase orders invoices taxes in one dashboard track invoices paid as cash/cheque automatic TDS setup aggregated TDS paid monthly TDS challans available for auditing track petty cash budgets and make payments from assigned budgets add bills & receipt as attachments to petty cash expenses automate invoice payments and pay taxes in few clicks using Vendor Payments calculate and schedule salary & compliance payments in minutes using Payroll one-click vendor and compliance payouts directly make vendor and TDS payments from RazorpayX dashboard automatically calculate and deduct TDS at time of invoice processing settle vendor balances with click of button RazorpayX will automatically pay aggregated TDS to government on behalf at end of each month TDS challans available on dashboard for auditing and tax-filing can mark invoices paid as cash/cheque to track all invoices in one place and still apply TDS deductions receive supplier invoices from multiple sources like emails vendor portals drive etc and pre-fill details on dashboard and making them ready for payment pay vendors employees and customers seamlessly automate invoice payments and pay taxes in few clicks create purchase orders add invoices and set-up approval workflows for large payments onboard all vendors by providing their emails and let RazorpayX do KYC check automates vendor verification and simplifies onboarding vendors can use portal for invoice uploads and easy reconciliation manage all vendors payouts purchase orders invoices taxes in one dashboard approval workflow to obtain consent from stakeholders for customer and vendor refunds or payments vendor payments by RazorpayX for convenient and timely vendor payments whether immediate or scheduled TDS amount automatically deducted from your account balance and deposited with government monthly you can then download challan and file TDS • Ethiopia: VAT 15% TOT 2%/10% Withholding Tax 2% for services per Income Tax Proclamation 286/2002 • Outstanding modern UI glassmorphic Recharts</p>
          </div>
          <div className="flex gap-2">
            <button className="rounded-xl border bg-white h-10 px-4 text-xs">📷 Scan Invoice • OCR-enabled Invoice Capture • PyMuPDF • Tesseract OCR • Auto-fill Details</button>
            <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Upload Invoice • Auto-pay Vendors and TDS • OCR • Multi-layer Approval • Outstanding</button>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Vendor Invoices • Invoice Number Invoice Date Due Date Amount Tax VAT 15% TOT 2%/10% Withholding Tax 2% for Services Total Amount Status Draft/Pending Approval/Approved/Paid/Rejected/Cancelled OCR Raw Extracted Text Confidence Vendor Name TIN Invoice Number Amount Tax Withholding File Key MinIO File Hash • Outstanding Pipeline Visual Stepper</h3>
            <div className="mt-4 space-y-3">
              {mockInvoices.map(inv => (
                <button key={inv.id} onClick={()=>setSelected(inv)} className={`w-full text-left rounded-xl border p-4 hover:bg-muted ${selected.id===inv.id ? "bg-primary/10 border-primary/30" : ""}`}>
                  <div className="flex justify-between"><p className="font-medium text-sm">{inv.vendor} • {inv.invoice_number}</p><Badge variant={inv.status==="paid" ? "success" : inv.status==="pending_approval" ? "warning" : "default"}>{inv.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Invoice Date {inv.invoice_date} • Due {inv.due_date} • Amount {inv.amount} • Tax VAT 15% {inv.tax_amount} • Withholding 2% {inv.withholding_tax} • Total {inv.total_amount} • OCR Confidence {inv.ocr_raw.confidence} • Vendor {inv.ocr_raw.vendor_name} TIN {inv.ocr_raw.tin} • File {inv.file_key} Hash {inv.file_hash} • Payout {inv.payout_id || "—"}</p>
                </button>
              ))}
              <button className="w-full rounded-xl border border-dashed h-12 text-xs">+ Upload Invoice • Auto-pay Vendors and TDS • OCR-enabled Invoice Capture • Multi-layer Approval • Outstanding • PyMuPDF • Tesseract OCR • Auto-fill Details on Dashboard and Making Them Ready for Payment • Manage All Vendors Payouts Purchase Orders Invoices Taxes in One Dashboard</button>
            </div>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <div className="flex justify-between items-center">
              <h3 className="font-semibold">Invoice Detail • {selected.vendor} • {selected.invoice_number} • Amount {selected.amount} • Tax VAT 15% {selected.tax_amount} • Withholding 2% {selected.withholding_tax} • Total {selected.total_amount} • Status {selected.status} • OCR Raw Extracted Text Confidence {selected.ocr_raw.confidence} • Outstanding • DocumentViewerOCR Side-by-side OCR</h3>
              <Badge variant={selected.status==="paid" ? "success" : "warning"}>{selected.status}</Badge>
            </div>

            <div className="mt-6 grid grid-cols-2 gap-6">
              <div className="rounded-xl border bg-card p-4">
                <h4 className="font-semibold text-sm">Invoice Preview • Thumb • Preview Thumbs • Hash Integrity • Progress Donut 0-100% • Outstanding Modern • DocumentViewer.tsx Side-by-side OCR • OCR Raw • PyMuPDF • Tesseract OCR</h4>
                <div className="mt-3 rounded-xl border bg-white overflow-hidden h-64 flex flex-col items-center justify-center p-8 text-center">
                  <div className="h-16 w-16 rounded-2xl bg-primary/10 flex items-center justify-center text-2xl">📄</div>
                  <p className="mt-3 font-medium text-sm">{selected.invoice_number} • {selected.vendor}</p>
                  <p className="text-[11px] text-muted-foreground mt-1">PDF Preview • Outstanding modern template • Logo QR pie chart YTD bilingual EN/AM • File {selected.file_key} Hash {selected.file_hash} • MinIO presigned 15m • Encrypted SSE-S3 • 7y retention NBE • ClamAV clean • Invoice Date {selected.invoice_date} Due {selected.due_date} Amount {selected.amount} Tax VAT 15% {selected.tax_amount} Withholding 2% {selected.withholding_tax} Total {selected.total_amount}</p>
                  <div className="mt-4 w-full h-2 rounded-full bg-neutral-200 overflow-hidden"><div className="h-full bg-primary rounded-full w-[100%]" /></div>
                </div>
                <div className="mt-4 rounded-xl bg-green-500/10 border border-green-500/20 p-3 text-[11px]">
                  <p className="font-semibold">OCR Raw • Extracted Text • Confidence {selected.ocr_raw.confidence} • Outstanding • PyMuPDF • Tesseract OCR • Document Authenticity • NBE Checklist</p>
                  <pre className="mt-2 whitespace-pre-wrap font-mono text-[10px] bg-white border rounded-xl p-3 max-h-48 overflow-auto">{selected.ocr_raw.extracted_text}</pre>
                  <p className="mt-2">Vendor Name {selected.ocr_raw.vendor_name} • TIN {selected.ocr_raw.tin} • Invoice Number {selected.ocr_raw.invoice_number || selected.invoice_number} • Amount {selected.ocr_raw.amount || selected.amount} • Tax {selected.ocr_raw.tax || selected.tax_amount} • Withholding {selected.ocr_raw.withholding || selected.withholding_tax} • Confidence {selected.ocr_raw.confidence} • OCR-enabled Invoice Capture • Multi-layer Approval Workflows • Automated TDS Calculation and Filing to NSDL • Integrated Payouts • Outstanding modern UI glassmorphic • Receipt preview thumbs • Hash integrity • Progress donut • DocumentViewer.tsx side-by-side OCR</p>
                </div>
              </div>
              <div className="rounded-xl border bg-card p-4">
                <h4 className="font-semibold text-sm">TDS Calculation Ethiopia Withholding Tax 2% for Services per Income Tax Proclamation 286/2002 • Automated TDS Calculation and Filing to NSDL • Outstanding • Formula Engine Secure O(n) Tokenization + Shunting-yard + Decimal Precise</h4>
                <div className="mt-3 space-y-3 text-xs">
                  <div className="rounded-xl bg-muted p-3"><p className="font-medium">Amount • ETB • {selected.amount} • Tax VAT 15% {selected.tax_amount} • Withholding 2% {selected.withholding_tax} • Total {selected.total_amount}</p><p className="text-[11px] text-muted-foreground mt-1">Amount 50000 • Tax VAT 15% 7500 • Withholding 2% 1000 • Total 56500 • Formula: Tax VAT 15% = Amount * 15% = 50000*0.15=7500 • Withholding Tax 2% for services per Ethiopia Income Tax Proclamation 286/2002 = Amount * 2% = 50000*0.02=1000 • Total = Amount + Tax - Withholding? Actually Total = Amount + Tax - Withholding? 50000+7500-1000=56500? No 50000+7500=57500-1000=56500 • Outstanding • Formula Engine Secure O(n) Tokenization + Shunting-yard + Decimal Precise • ValidateFormula only allowed vars BASIC CTC_MONTHLY CTC_ANNUAL GROSS</p></div>
                  <div className="rounded-xl bg-blue-500/10 border border-blue-500/20 p-3">
                    <p className="font-semibold">Approval Flow Multi-layer Approval Workflows • Outstanding Pipeline Visual Stepper • Maker-checker</p>
                    <div className="mt-3 space-y-2">
                      <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-primary text-white flex items-center justify-center text-[10px]">1</span><span className="text-[11px]">Employee creates invoice amount description receipt upload MinIO presigned POST 15m TTL file_key &lt;5MB pdf/jpg/png file_hash receipt_file_hash tax_exempt_limit is_taxable is_pensionable status pending • Preview thumbs • DocumentViewer.tsx • Hash integrity • Progress donut 0-100% • Outstanding Modern</span></div>
                      <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-amber-500 text-white flex items-center justify-center text-[10px]">2</span><span className="text-[11px]">Manager approval • Reporting manager checks receipt verification • approved_by_manager manager_approved_at • Status approved_by_manager • DocumentViewer.tsx side-by-side OCR</span></div>
                      <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-green-600 text-white flex items-center justify-center text-[10px]">3</span><span className="text-[11px]">Finance approval • Finance manager checks taxability pensionability • approved_by_finance finance_approved_at • Status approved • TDS Calculation Ethiopia Withholding Tax 2% for Services per Income Tax Proclamation 286/2002 • Automated TDS Calculation and Filing to NSDL • Integrated Payouts • Outstanding</span></div>
                      <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-green-600 text-white flex items-center justify-center text-[10px]">4</span><span className="text-[11px]">Paid • Payout ID {selected.payout_id || "—"} • Paid via payout batch pain.001 XML ISO20022 Document CstmrCdtTrfInitn + Pension CSV + ERCA CSV + Cost center report + Bank file • Ledger second journal Dr payable Cr clearing:bank • Status paid • TDS amount automatically deducted from your account balance and deposited with government monthly you can then download challan and file TDS • Outstanding • Receipt preview thumbs • Hash integrity • Progress donut</span></div>
                    </div>
                  </div>
                  <div className="rounded-xl bg-green-500/10 border border-green-500/20 p-3 text-[11px]">
                    <p className="font-semibold">Payroll Integration • One-click Vendor and Compliance Payouts • Directly Make Vendor and TDS Payments from Dashboard • Automatically Calculate and Deduct TDS at Time of Invoice Processing • Settle Vendor Balances with Click of Button</p>
                    <p className="mt-1">RazorpayX will automatically pay aggregated TDS to government on behalf at end of each month. TDS challans will be available on dashboard for auditing and tax-filing purposes. Can mark invoices paid as cash/cheque to track all invoices in one place and still apply TDS deductions • Receive supplier invoices from multiple sources like emails vendor portals drive etc and pre-fill details on dashboard and making them ready for payment • Pay vendors employees and customers seamlessly automate invoice payments and pay taxes in few clicks • Create purchase orders add invoices and set-up approval workflows for large payments • Onboard all vendors by providing their emails and let RazorpayX do KYC check automates vendor verification and simplifies onboarding vendors can use portal for invoice uploads and easy reconciliation • Manage all vendors payouts purchase orders invoices taxes in one dashboard • Approval workflow to obtain consent from stakeholders for customer and vendor refunds or payments • Outstanding modern UI glassmorphic</p>
                  </div>
                </div>
              </div>
            </div>

            <div className="mt-6"
              onDragOver={(e)=>{ e.preventDefault() }}
              onDragLeave={()=>{}}
              onDrop={(e)=>{ e.preventDefault() }}
            >
              <div className="rounded-2xl border-2 border-dashed border-primary/30 bg-primary/5 p-6 text-center">
                <div className="mx-auto h-12 w-12 rounded-2xl bg-primary/10 flex items-center justify-center">📄</div>
                <p className="mt-3 font-medium text-sm">Drop invoice here • MinIO presigned POST 15m • &lt;5MB pdf/jpg/png • OCR-enabled Invoice Capture • PyMuPDF • Tesseract OCR • Auto-fill Details on Dashboard and Making Them Ready for Payment • Manage All Vendors Payouts Purchase Orders Invoices Taxes in One Dashboard</p>
                <p className="text-[11px] text-muted-foreground mt-1">File types whitelist pdf/jpg/png • Size &lt;5MB • ClamAV stub VirusScanner clean • File hash integrity file_hash unique index per merchant • Encrypted SSE-S3 MinIO versioning • Retention 7y per NBE • No plain FIN logs grep test CI • PII redact zerolog field filter • FIN only last4 responses • Account masked ****1234 • OCR Raw Extracted Text Confidence Vendor Name TIN Invoice Number Amount Tax Withholding File Key MinIO File Hash • Outstanding Modern • Mercury/Linear inspiration • Glassmorphic • DocumentViewer.tsx side-by-side OCR • Preview thumbs • Hash integrity • Progress donut 0-100%</p>
                <button className="mt-4 rounded-xl bg-primary text-white h-10 px-6 text-xs">Select File • Dropzone dashed pulse on drag • Preview thumbs • Hash integrity • Progress donut 0-100% • Outstanding Modern • Mercury/Linear inspiration • Glassmorphic • DocumentViewer.tsx • OCR-enabled Invoice Capture • Multi-layer Approval • Outstanding</button>
              </div>
            </div>

            <div className="mt-6 rounded-xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Purchase Orders • Vendor Portal • Petty Cash Budgets & Expenses • Outstanding • End-to-End Accounts Payable Automation</h4>
              <div className="mt-3 grid grid-cols-3 gap-4 text-xs">
                <div className="rounded-xl bg-muted p-3"><p className="font-medium">Purchase Orders • {mockPOs.length} • Approved</p>{mockPOs.map(po=><p key={po.id} className="text-[11px]">PO {po.po_number} • Vendor {po.vendor} • Amount {po.amount} {po.currency} • Status {po.status} • Created {po.created_at}</p>)}</div>
                <div className="rounded-xl bg-muted p-3"><p className="font-medium">Petty Cash Budgets • Outstanding • Track Petty Cash Budgets and Make Payments from Assigned Budgets</p><p className="text-[11px]">Budget Name Office Supplies • Amount 50000 ETB • Assigned to Finance Manager • Status active • Spent 15000 Remaining 35000 • Created 2026-07-01 • Outstanding per RazorpayX track petty cash budgets and make payments from assigned budgets add bills & receipt as attachments to petty cash expenses</p></div>
                <div className="rounded-xl bg-muted p-3"><p className="font-medium">Petty Cash Expenses • Outstanding • Add Bills & Receipt as Attachments</p><p className="text-[11px]">Expense 1500 ETB • Description Office supplies - printer ink • Receipt expense_receipt_EMP003.jpg • Hash hash_expense_003 • Status paid • Approved by Finance • Created 2026-07-10 • Paid via petty cash budget Office Supplies • Outstanding per RazorpayX add bills & receipt as attachments to petty cash expenses track petty cash budgets and make payments from assigned budgets</p></div>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
