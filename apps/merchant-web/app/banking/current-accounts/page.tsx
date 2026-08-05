"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockCurrentAccounts = [
  { id: "ca_001", account_number: "ETB-CBE-1000123456789", account_name: "Apex Trading PLC • Primary Settlement • አፔክስ", bank_code: "CBE", partner_bank_name: "Commercial Bank of Ethiopia • የኢትዮጵያ ንግድ ባንክ", account_type: "current", currency: "ETB", status: "active", balance: "1250430", available_balance: "1250430", overdraft_limit: "0", is_primary: true, is_lite: false, is_virtual: false, cheque_book_issued: true, debit_card_issued: true, debit_card_type: "both", created_at: "2026-01-15" },
  { id: "ca_002", account_number: "ETB-AWASH-2000987654321", account_name: "Apex Trading PLC • Vendor Payments • አቅራቢ ክፍያ", bank_code: "AWASH", partner_bank_name: "Awash Bank • አዋሽ ባንክ", account_type: "current", currency: "ETB", status: "active", balance: "500000", available_balance: "500000", is_primary: false, is_lite: false, is_virtual: false, cheque_book_issued: true, debit_card_issued: false, created_at: "2026-02-01" },
  { id: "ca_003", account_number: "ETB-DASHEN-VA-3000112233", account_name: "Apex Trading PLC • Collections • Smart Collect • ስብስብ", bank_code: "DASHEN", partner_bank_name: "Dashen Bank • ዳሽን ባንክ", account_type: "virtual", currency: "ETB", status: "active", balance: "0", available_balance: "0", is_primary: false, is_lite: false, is_virtual: true, cheque_book_issued: false, debit_card_issued: false, created_at: "2026-03-01" },
  { id: "ca_004", account_number: "ETB-CBE-ESCROW-4000998877", account_name: "Apex Trading PLC • Escrow • Marketplace • እስክሮ", bank_code: "CBE", partner_bank_name: "Commercial Bank of Ethiopia • የኢትዮጵያ ንግድ ባንክ", account_type: "escrow", currency: "ETB", status: "active", balance: "100000", available_balance: "0", is_primary: false, is_lite: false, is_virtual: false, cheque_book_issued: false, debit_card_issued: false, created_at: "2026-04-01" },
]

const mockChequeBooks = [
  { id: "chq_001", current_account_id: "ca_001", cheque_book_number: "CHQ-CBE-001-2026", start_cheque_number: 1, end_cheque_number: 25, total_cheques: 25, used_cheques: 5, status: "active", issued_at: "2026-01-20" },
]

const mockDebitCards = [
  { id: "dc_001", current_account_id: "ca_001", card_number_masked: "****1234", card_type: "both", card_network: "visa", status: "active", daily_limit: "50000", monthly_limit: "500000", cardholder_name: "Abebe Kebede", is_contactless: true, created_at: "2026-01-20" },
]

const mockEscrow = [
  { id: "esc_001", agreement_id: "agr_001", account_number: "ETB-CBE-ESCROW-4000998877", account_name: "Order #12345 • Buyer Merkato Seller Habesha • Escrow • 1000 ETB • Platform Fee 10% 100 ETB Seller 90% 900 ETB Withholding 2% 20 ETB", amount: "1000", currency: "ETB", status: "held", held_at: "2026-07-25", buyer: "Merkato Trading PLC", seller: "Habesha Crafts PLC", order_id: "ORD-12345", order_amount: "1000", platform_fee: "100", seller_amount: "900", withholding_tax: "20", expires_at: "2026-08-01" },
]

const mockCorporateCards = [
  { id: "cc_001", card_number_masked: "****5678", card_type: "virtual", card_network: "visa", cardholder_name: "Meron Hailu • Engineering", status: "active", credit_limit: "2000000", available_credit: "1850000", daily_limit: "50000", monthly_limit: "500000", category_restrictions: ["SaaS", "Cloud", "Marketing"], cashback_percent: "1.00", forex_markup_percent: "2.50", interest_free_days: 45, is_addon: false, created_at: "2026-02-15" },
]

export default function CurrentAccountsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Current Accounts • የሂሳብ አስተዳደር • Real Partner Bank CBE/Awash/Dashen + Cheque Book + Debit Card + Lite Interim + Multiple Accounts Balance Snapshot • RazorpayX Parity • Ethiopia Business Banking Core P0</h1>
            <p className="text-sm text-muted-foreground mt-2">Online opening &lt;24h paperless, no minimum balance, no monthly maintenance fee, free tier genuinely free, issued by partner banks CBE/Awash/Dashen (Ethiopia equivalent of RBL/ICICI/Axis/YES in RazorpayX India), comes with cheque book, debit card, real-time online dashboard, unlimited cash deposits/withdrawals, unlimited transfers no transfer limits, instant loans without collaterals only after 3 months transactions history, dedicated relationship manager priority support, multiple accounts view all balances listed, snapshot payouts insights over last 30 days, go to insights dashboard detailed analytics, track activity day/week/month overall analysis, get in-depth reporting into cash flow trends, real-time cash flow insights intuitive dashboard immediate visibility generate financial reports in minutes, can use as default settlement account for payment gateway settlements • Heart of business finances per RazorpayX App Store description</p>
          </div>
          <div className="flex gap-2">
            <button className="rounded-xl border bg-white h-10 px-4 text-xs">Refresh Balance • Instant • SWR • Real-time • Latest balance instantly using refresh button • RazorpayX mobile app feature</button>
            <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Open Current Account • Online &lt;24h Paperless • No Min Balance • CBE/Awash/Dashen • Partner Bank</button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {mockCurrentAccounts.map(ca => (
            <Card key={ca.id} className="p-5 hover:shadow-medium transition-shadow">
              <div className="flex justify-between items-start"><p className="font-medium text-sm">{ca.account_name}</p><Badge variant={ca.status==="active" ? "success" : "warning"}>{ca.status}</Badge></div>
              <p className="text-[11px] text-muted-foreground mt-1">Account Number {ca.account_number} • {ca.partner_bank_name} • Bank Code {ca.bank_code} • Type {ca.account_type} • Currency {ca.currency} • Primary {ca.is_primary ? "Yes • Primary Settlement • Heart of business finances" : "No"} • Lite {ca.is_lite ? "Yes • RazorpayX Lite interim until current account active" : "No"} • Virtual {ca.is_virtual ? "Yes • Smart Collect Virtual Account • Automatically reconcile incoming NEFT RTGS IMPS UPI payments using Virtual Accounts & UPI-IDs" : "No"}</p>
              <p className="text-2xl font-bold mt-3">ETB {parseInt(ca.balance).toLocaleString()} • Balance • Available {ca.available_balance}</p>
              <p className="text-[11px] text-muted-foreground mt-1">Overdraft Limit {ca.overdraft_limit} • Cheque Book Issued {ca.cheque_book_issued ? "Yes • Cheque Book Number CHQ-CBE-001-2026 Start 1 End 25 Total 25 Used 5 Status active" : "No"} • Debit Card Issued {ca.debit_card_issued ? "Yes • Debit Card Type " + ca.debit_card_type + " • Virtual + Physical • Both • Card Network Visa • Daily Limit 50000 Monthly 500000 Cardholder Abebe Kebede Contactless true" : "No"} • Created {ca.created_at} • Multiple accounts on RazorpayX view all balances listed</p>
              <div className="mt-3 flex gap-2">
                <button className="rounded-xl border h-8 px-3 text-[11px]">View Statement • Cheque Book • Debit Card • Unlimited Deposits/Withdrawals • No Transfer Limits</button>
                <button className="rounded-xl border h-8 px-3 text-[11px]">Transfer • IMPS NEFT RTGS UPI Card • Single REST endpoint • Payouts API • Bulk Payouts 50k with one OTP</button>
              </div>
            </Card>
          ))}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Cheque Books • Issuance Tracking per Current Account • Start Cheque Number End Cheque Number Total Cheques Used Cheques Status Ordered/Issued/Active/Used Up/Blocked/Cancelled • Outstanding</h3>
            <div className="mt-4 space-y-3">
              {mockChequeBooks.map(chq => (
                <div key={chq.id} className="rounded-xl border p-3 text-xs">
                  <p className="font-medium">Cheque Book {chq.cheque_book_number} • Account {chq.current_account_id} • Start {chq.start_cheque_number} End {chq.end_cheque_number} Total {chq.total_cheques} Used {chq.used_cheques} • Status {chq.status} • Issued {chq.issued_at}</p>
                  <div className="mt-2 h-2 rounded-full bg-neutral-100 overflow-hidden"><div className="h-full bg-primary rounded-full" style={{ width: `${(chq.used_cheques/chq.total_cheques)*100}%` }} /></div>
                  <p className="text-[11px] mt-1">Cheque Book per current account: start_cheque_number end_cheque_number total_cheques used_cheques status ordered/issued/active/used_up/blocked/cancelled issued_at issued_by • Outstanding per RazorpayX current account comes with cheque book • Unlimited deposits/withdrawals • No transfer limits</p>
                </div>
              ))}
              <button className="w-full rounded-xl border border-dashed h-10 text-xs">+ Order Cheque Book • 25 cheques • Start 26 End 50 • Status ordered • Issued by Finance Manager • Cheque book number CHQ-CBE-002-2026</button>
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="font-semibold">Debit Cards • Virtual + Physical Issuance Tracking • Card Number Masked ****1234 Card Number Hash SHA256 Hash Last4 Card Type Virtual/Physical/Both Card Network Visa/Mastercard/Verve/EthSwitch Status Ordered/Active/Blocked/Expired/Cancelled Daily Limit Monthly Limit Cardholder Name Expiry Month Year CVV Hash Is Contactless • Outstanding</h3>
            <div className="mt-4 space-y-3">
              {mockDebitCards.map(card => (
                <div key={card.id} className="rounded-xl border p-3 text-xs">
                  <p className="font-medium">Debit Card {card.card_number_masked} • Account {card.current_account_id} • Type {card.card_type} • Network {card.card_network} • Status {card.status} • Cardholder {card.cardholder_name} • Contactless {card.is_contactless ? "Yes" : "No"} • Daily Limit {card.daily_limit} Monthly {card.monthly_limit} • Created {card.created_at}</p>
                  <p className="text-[11px] text-muted-foreground mt-1">Debit Cards Virtual + Physical Issuance Tracking: card_number_masked ****1234 card_number_hash sha256 hash last4 card_type virtual/physical/both card_network visa/mastercard/verve/ethswitch status ordered/active/blocked/expired/cancelled daily_limit monthly_limit cardholder_name expiry_month year cvv_hash is_contactless • Outstanding per RazorpayX current account comes with debit card • Real-time online dashboard • Unlimited deposits/withdrawals</p>
                </div>
              ))}
              <button className="w-full rounded-xl border border-dashed h-10 text-xs">+ Issue Debit Card • Virtual + Physical • Both • Card Network Visa • Daily Limit 50000 Monthly 500000 Cardholder Abebe Kebede Contactless true • Status ordered • Card Number Masked ****5678 Card Number Hash SHA256 • Outstanding</button>
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="font-semibold">Escrow Accounts Automated Marketplace P2P Hold & Release Funds Under Defined Conditions • Platform Fee 10% 100 ETB Seller 90% 900 ETB Withholding 2% 20 ETB Hold in Escrow Until Delivery Confirmed • Outstanding</h3>
            <div className="mt-4 space-y-3">
              {mockEscrow.map(esc => (
                <div key={esc.id} className="rounded-xl border p-3 text-xs">
                  <p className="font-medium">{esc.account_name}</p>
                  <p className="text-[11px] text-muted-foreground mt-1">Account Number {esc.account_number} • Agreement {esc.agreement_id} • Amount {esc.amount} {esc.currency} • Status {esc.status} • Held at {esc.held_at} • Buyer {esc.buyer} • Seller {esc.seller} • Order {esc.order_id} • Order Amount {esc.order_amount} • Platform Fee {esc.platform_fee} (10%) • Seller Amount {esc.seller_amount} (90%) • Withholding Tax {esc.withholding_tax} (2% per Ethiopia Income Tax Proclamation) • Expires {esc.expires_at} • Ledger Book ID escrow book per agreement book_type escrow • Auto Release After Days 7 • Conditions [{`{type: delivery_confirmed, days: 7}, {type: inspection_period, days: 3}`}] • Outstanding per RazorpayX automated escrow accounts added 2024 for marketplaces P2P platforms hold and release funds between parties under defined conditions reduces legal overhead</p>
                  <div className="mt-2 flex gap-2"><button className="rounded-xl bg-primary text-white h-7 px-3 text-[10px]">Release • Seller 90% 900 ETB Minus Fee 100 ETB Minus Withholding 20 ETB = 780 ETB? Actually seller_amount 900 • Platform Fee 100 • Withholding 20 • Net seller 880? • Ledger M4 Dr escrow Cr clearing</button><button className="rounded-xl border h-7 px-3 text-[10px]">Return • Buyer • Dispute • Expired • Return to buyer</button></div>
                </div>
              ))}
              <button className="w-full rounded-xl border border-dashed h-10 text-xs">+ Create Escrow Account • Marketplace Seller Settlement Split • Order Total 1000 ETB • Platform Fee 10% 100 ETB • Seller 90% 900 ETB • Withholding Tax 2% 20 ETB • Hold in Escrow Until Delivery Confirmed • Auto Release After Days 7 • Conditions delivery_confirmed inspection_period • Ledger Book Per Agreement Book Type Escrow</button>
            </div>
          </Card>
        </div>

        <Card className="p-6">
          <h3 className="font-semibold">Corporate Cards • Collateral-free Credit Cards Virtual + Physical Dynamic Limit up to 2Cr ETB Equivalent Up to 45-50 Day Interest-free Period 2.5% Forex Markup Flat 1% Cashback 30% Off SaaS Custom Spending Controls Real-time Expense Tracking Multi-level Approvals Virtual & Physical Add-on Cards Unlimited • Outstanding per RazorpayX</h3>
          <div className="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4">
            {mockCorporateCards.map(card => (
              <div key={card.id} className="rounded-xl border p-4 text-xs">
                <p className="font-medium">Corporate Card {card.card_number_masked} • Type {card.card_type} • Network {card.card_network} • Cardholder {card.cardholder_name} • Status {card.status} • Credit Limit {card.credit_limit} ETB • Available {card.available_credit} • Daily {card.daily_limit} Monthly {card.monthly_limit} • Category Restrictions {card.category_restrictions.join(", ")} • Cashback {card.cashback_percent}% • Forex Markup {card.forex_markup_percent}% • Interest Free Days {card.interest_free_days} • Addon {card.is_addon ? "Yes Unlimited Add-on Cards" : "No"} • Created {card.created_at}</p>
                <p className="text-[11px] text-muted-foreground mt-2">Corporate Cards Collateral-free Credit Cards Virtual + Physical Dynamic Limit up to 2Cr ETB Equivalent Up to 45-50 Day Interest-free Period 2.5% Forex Markup Flat 1% Cashback 30% Off SaaS Custom Spending Controls Real-time Expense Tracking Multi-level Approvals Virtual & Physical Add-on Cards Unlimited Minimum Documentation No Joining Fee Ideal for SaaS Subscriptions Cloud Services Marketing Expenses • Outstanding per RazorpayX • Up to ₹2Cr unsecured credit up to 45–50 day interest-free period 2.5% forex markup flat 1% cashback 30% off on 500+ SaaS/marketing tools custom spending controls real-time expense tracking seamless integration virtual and physical cards multi-level approvals</p>
                <div className="mt-2 flex gap-2"><button className="rounded-xl bg-primary text-white h-7 px-3 text-[10px]">View Transactions • Real-time Expense Tracking • Monitor all transactions instantly • Outstanding</button><button className="rounded-xl border h-7 px-3 text-[10px]">Block • Spending Controls • Category Restrictions • Daily Monthly Limit • Multi-level Approvals</button></div>
              </div>
            ))}
          </div>
          <button className="mt-4 w-full rounded-xl border border-dashed h-12 text-xs">+ Issue Corporate Card • Virtual + Physical • Both • Card Network Visa • Cardholder Meron Hailu Engineering • Credit Limit 2000000 ETB • Available Credit 2000000 • Daily Limit 50000 Monthly 500000 • Category Restrictions SaaS Cloud Marketing • Cashback 1% Forex Markup 2.5% Interest Free Days 45 • Addon No • Parent Card ID null • Outstanding per RazorpayX • Up to 2Cr unsecured credit up to 45–50 day interest-free period 2.5% forex markup flat 1% cashback 30% off on 500+ SaaS/marketing tools • Custom spending controls set specific limits and controls on card usage to prevent overspending • Real-time expense tracking monitor all transactions instantly • Seamless integration connect effortlessly with accounting software to automate expense reporting • Virtual and physical cards • Multi-level approvals • Ideal for SaaS subscriptions cloud services marketing expenses</button>
        </Card>
      </div>
    </div>
  )
}
