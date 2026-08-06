"use client"
import * as React from "react"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend } from "recharts"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockCards = [
  { id: "cc_001", card_number_masked: "****5678", card_type: "virtual", card_network: "visa", cardholder_name: "Meron Hailu • Engineering", cardholder_email: "meron@apextrading.et", status: "active", credit_limit: "2000000", available_credit: "1850000", daily_limit: "50000", monthly_limit: "500000", category_restrictions: ["SaaS", "Cloud", "Marketing"], cashback_percent: "1.00", forex_markup_percent: "2.50", interest_free_days: 45, is_addon: false, created_at: "2026-02-15", transactions: 45, total_spent: "150000", cashback_earned: "1500", forex_fees: "2500" },
  { id: "cc_002", card_number_masked: "****9012", card_type: "physical", card_network: "mastercard", cardholder_name: "Abebe Kebede • Sales", cardholder_email: "abebe@apextrading.et", status: "active", credit_limit: "1000000", available_credit: "800000", daily_limit: "30000", monthly_limit: "300000", category_restrictions: ["Marketing", "Travel"], cashback_percent: "1.00", forex_markup_percent: "2.50", interest_free_days: 45, is_addon: false, created_at: "2026-03-01", transactions: 30, total_spent: "200000", cashback_earned: "2000", forex_fees: "0" },
  { id: "cc_003", card_number_masked: "****3456", card_type: "virtual", card_network: "verve", cardholder_name: "Almaz Tadesse • Marketing", cardholder_email: "almaz@apextrading.et", status: "blocked", credit_limit: "500000", available_credit: "500000", daily_limit: "20000", monthly_limit: "100000", category_restrictions: ["SaaS"], cashback_percent: "1.00", forex_markup_percent: "2.50", interest_free_days: 45, is_addon: true, parent_card_id: "cc_001", created_at: "2026-04-01", transactions: 5, total_spent: "0", cashback_earned: "0", forex_fees: "0" },
]

const mockTransactions = [
  { id: "cctx_001", card_id: "cc_001", amount: "5000", currency: "ETB", merchant_name: "AWS • Amazon Web Services • Cloud", merchant_category: "Cloud", status: "approved", cashback_amount: "50", forex_fee: "0", created_at: "2026-07-25T10:00:00Z" },
  { id: "cctx_002", card_id: "cc_001", amount: "2000", currency: "ETB", merchant_name: "Google Cloud • Cloud", merchant_category: "Cloud", status: "approved", cashback_amount: "20", forex_fee: "0", created_at: "2026-07-26T11:00:00Z" },
  { id: "cctx_003", card_id: "cc_001", amount: "10000", currency: "ETB", merchant_name: "Facebook Ads • Marketing", merchant_category: "Marketing", status: "approved", cashback_amount: "100", forex_fee: "0", created_at: "2026-07-27T12:00:00Z" },
  { id: "cctx_004", card_id: "cc_001", amount: "100", currency: "USD", merchant_name: "OpenAI • SaaS • International", merchant_category: "SaaS", status: "approved", cashback_amount: "1", forex_fee: "2.50", created_at: "2026-07-28T13:00:00Z" },
  { id: "cctx_005", card_id: "cc_001", amount: "1000", currency: "ETB", merchant_name: "Unknown Merchant • Blocked Category • Declined", merchant_category: "Gambling", status: "declined", decline_reason: "Category restricted per spending controls • Gambling/Crypto/Adult blocked per NBE", cashback_amount: "0", forex_fee: "0", created_at: "2026-07-28T14:00:00Z" },
]

export default function CorporateCardsPage() {
  const [selected, setSelected] = React.useState(mockCards[0])
  const chartData = [
    { category: "SaaS", spent: 20000, limit: 50000 },
    { category: "Cloud", spent: 7000, limit: 30000 },
    { category: "Marketing", spent: 10000, limit: 20000 },
    { category: "Travel", spent: 5000, limit: 15000 },
  ]
  const pieData = [
    { name: "SaaS Subscriptions", value: 150000, color: "#0B6E4F" },
    { name: "Cloud Services", value: 70000, color: "#10B981" },
    { name: "Marketing", value: 100000, color: "#EAB308" },
    { name: "Other", value: 20000, color: "#E4E4E7" },
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Corporate Cards • የድርጅት ካርዶች • Collateral-free Credit Cards Virtual + Physical Dynamic Limit up to 2Cr ETB Equivalent Up to 45-50 Day Interest-free Period 2.5% Forex Markup Flat 1% Cashback 30% Off SaaS Custom Spending Controls Real-time Expense Tracking Multi-level Approvals • RazorpayX Parity • P0</h1>
            <p className="text-sm text-muted-foreground mt-2">Collateral-free corporate credit card with dynamic credit limit up to Rs.20 lakhs ideal for SaaS subscriptions cloud services marketing expenses as these can only be made with credit card guaranteed collateral-free free corporate credit card along with unlimited add-on cards with minimum documentation and no joining fee access suite of smart apps and integrations like Payout Links Vendor Payments and Payroll — Only Neo Bank offers this feature up to ₹2Cr unsecured credit up to 45–50 day interest-free period 2.5% forex markup flat 1% cashback and 30% off on 500+ SaaS/marketing tools custom spending controls set specific limits and controls on card usage to prevent overspending and ensure adherence to company policies real-time expense tracking monitor all transactions instantly for improved visibility and timely reconciliation seamless integration connect effortlessly with accounting software to automate expense reporting and reduce manual errors virtual and physical cards utilize both virtual cards for online purchases and physical cards for offline expenses adapting to various business needs multi-level approvals implement approval workflows to maintain oversight and authorize spending before transactions occur ideal for SaaS subscriptions cloud services marketing expenses as these can only be made with credit card • Ethiopia: Corporate card issuing via banks CBE Dashen offer corporate cards but could implement virtual card generation via partner bank API mock physical card tracking spending controls per card per merchant per category expense tracking ledger multi-level approvals for card creation real-time transaction webhook cashback calculation + Up to 2Cr unsecured credit up to 45–50 day interest-free period 2.5% forex markup flat 1% cashback 30% off on 500+ SaaS/marketing tools • Outstanding modern UI glassmorphic Recharts</p>
          </div>
          <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Issue Corporate Card • Virtual + Physical • Both • Card Network Visa • Daily Limit 50000 Monthly 500000 • Category Restrictions SaaS Cloud Marketing • Outstanding</button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Corporate Cards • Virtual + Physical • Dynamic Limit up to 2Cr ETB Equivalent • Credit Limit • Available Credit • Daily Monthly Limit • Category Restrictions • Cashback 1% • Forex Markup 2.5% • Interest Free Days 45 • Addon Unlimited • Outstanding Pipeline Visual Stepper</h3>
            <div className="mt-4 space-y-3">
              {mockCards.map(card => (
                <button key={card.id} onClick={()=>setSelected(card)} className={`w-full text-left rounded-xl border p-4 hover:bg-muted ${selected.id===card.id ? "bg-primary/10 border-primary/30" : ""}`}>
                  <div className="flex justify-between"><p className="font-medium text-sm">{card.cardholder_name} • {card.card_number_masked} • {card.card_type} • {card.card_network}</p><Badge variant={card.status==="active" ? "success" : card.status==="blocked" ? "danger" : "warning"}>{card.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Credit Limit {card.credit_limit} ETB • Available {card.available_credit} • Daily {card.daily_limit} Monthly {card.monthly_limit} • Category Restrictions {card.category_restrictions.join(", ")} • Cashback {card.cashback_percent}% • Forex Markup {card.forex_markup_percent}% • Interest Free Days {card.interest_free_days} • Addon {card.is_addon ? "Yes Unlimited Add-on Cards" : "No"} • Created {card.created_at} • Transactions {card.transactions} • Total Spent {card.total_spent} ETB • Cashback Earned {card.cashback_earned} • Forex Fees {card.forex_fees} • Outstanding per RazorpayX up to ₹2Cr unsecured credit up to 45–50 day interest-free period 2.5% forex markup flat 1% cashback 30% off on 500+ SaaS/marketing tools</p>
                  <div className="mt-2 h-2 rounded-full bg-neutral-100 overflow-hidden"><div className="h-full bg-primary rounded-full" style={{ width: `${(parseInt(card.total_spent)/parseInt(card.credit_limit))*100}%` }} /></div>
                </button>
              ))}
              <button className="w-full rounded-xl border border-dashed h-12 text-xs">+ Issue Corporate Card • Virtual + Physical • Both • Card Network Visa • Cardholder Meron Hailu Engineering • Credit Limit 2000000 ETB • Available Credit 2000000 • Daily Limit 50000 Monthly 500000 • Category Restrictions SaaS Cloud Marketing • Cashback 1% Forex Markup 2.5% Interest Free Days 45 • Addon No • Parent Card ID null • Outstanding per RazorpayX • Up to 2Cr unsecured credit up to 45–50 day interest-free period 2.5% forex markup flat 1% cashback 30% off on 500+ SaaS/marketing tools • Custom spending controls set specific limits and controls on card usage to prevent overspending • Real-time expense tracking monitor all transactions instantly • Seamless integration connect effortlessly with accounting software to automate expense reporting • Virtual and physical cards • Multi-level approvals • Ideal for SaaS subscriptions cloud services marketing expenses</button>
            </div>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <div className="flex justify-between items-center">
              <h3 className="font-semibold">Card Detail • {selected.cardholder_name} • {selected.card_number_masked} • {selected.card_type} • {selected.card_network} • Credit Limit {selected.credit_limit} • Available {selected.available_credit} • Daily {selected.daily_limit} Monthly {selected.monthly_limit} • Category Restrictions {selected.category_restrictions.join(", ")} • Cashback {selected.cashback_percent}% • Forex Markup {selected.forex_markup_percent}% • Interest Free Days {selected.interest_free_days} • Addon {selected.is_addon ? "Yes Unlimited Add-on Cards" : "No"} • Spending Controls • Real-time Expense Tracking • Multi-level Approvals • Outstanding • Recharts Bar Spending vs Limit + Pie Category Breakdown</h3>
              <Badge variant={selected.status==="active" ? "success" : "danger"}>{selected.status}</Badge>
            </div>

            <div className="mt-6 grid grid-cols-2 gap-6">
              <div>
                <p className="text-xs font-semibold mb-2">Bar Chart Spending vs Limit per Category • Outstanding Recharts • O(n) per card n=categories • Category Restrictions SaaS Cloud Marketing • Spending Controls Daily Limit Monthly Limit Allowed Categories Blocked Merchants</p>
                <div className="h-64">
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={chartData}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="category" tick={{ fontSize: 10 }} />
                      <YAxis tick={{ fontSize: 10 }} />
                      <Tooltip />
                      <Bar dataKey="spent" fill="#0B6E4F" name="Spent" />
                      <Bar dataKey="limit" fill="#E4E4E7" name="Limit" />
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              </div>
              <div>
                <p className="text-xs font-semibold mb-2">Pie Chart Category Breakdown • SaaS Subscriptions Cloud Services Marketing Expenses • Outstanding modern template • Collateral-free Corporate Credit Card • SaaS Subscriptions Cloud Services Marketing Expenses as These Can Only be Made with a Credit Card</p>
                <div className="h-64">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie data={pieData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={80} label>
                        {pieData.map((entry, index) => <Cell key={`cell-${index}`} fill={entry.color} />)}
                      </Pie>
                      <Tooltip />
                      <Legend />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
              </div>
            </div>

            <div className="mt-6 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-7 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Transaction ID • Card ID • Amount • Currency • Merchant Name • AWS Google Cloud Facebook Ads • Merchant Category • SaaS Cloud Marketing • Status Pending/Approved/Declined/Reversed • Decline Reason • Category Restricted • Cashback 1% Flat • Forex Fee 2.5% • Created At</span></div>
              {mockTransactions.map(tx => (
                <div key={tx.id} className="grid grid-cols-7 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                  <span className="font-mono text-[11px]">{tx.id} • Card {tx.card_id} • Amount {tx.amount} {tx.currency} • Merchant {tx.merchant_name} • Category {tx.merchant_category} • Status {tx.status} • Decline {tx.decline_reason || "—"} • Cashback {tx.cashback_amount} • Forex Fee {tx.forex_fee} • Created {tx.created_at}</span>
                  <span className="font-bold">ETB {tx.amount} {tx.currency}</span>
                  <span>{tx.merchant_name} • {tx.merchant_category}</span>
                  <span><Badge variant={tx.status==="approved" ? "success" : tx.status==="declined" ? "danger" : "warning"}>{tx.status} • Decline {tx.decline_reason || "—"}</Badge></span>
                  <span>Cashback {tx.cashback_amount} • 1% flat • Forex Fee {tx.forex_fee} • 2.5% forex markup if international • Created {tx.created_at}</span>
                  <span className="flex flex-col gap-1"><button className="text-primary text-[11px]">View • Real-time Expense Tracking • Monitor all transactions instantly for improved visibility and timely reconciliation • Seamless integration connect effortlessly with accounting software to automate expense reporting and reduce manual errors</button><button className="text-red-500 text-[11px]">Block • Spending Controls • Category Restrictions • Daily Monthly Limit • Multi-level Approvals • Outstanding</button></span>
                </div>
              ))}
            </div>

            <div className="mt-6 grid grid-cols-2 gap-6 text-xs">
              <div className="rounded-xl bg-muted p-4"><p className="font-semibold">Spending Controls • Custom Spending Controls Set Specific Limits and Controls on Card Usage to Prevent Overspending and Ensure Adherence to Company Policies • Outstanding • Map O(1) Lookup</p><p className="mt-2 text-[11px]">Spending Controls per Card per Merchant per Category Expense Tracking Ledger Multi-level Approvals for Card Creation Real-time Transaction Webhook Cashback Calculation + Up to 2Cr unsecured credit up to 45–50 day interest-free period 2.5% forex markup flat 1% cashback 30% off on 500+ SaaS/marketing tools • Custom spending controls set specific limits and controls on card usage to prevent overspending and ensure adherence to company policies • Real-time expense tracking monitor all transactions instantly for improved visibility and timely reconciliation • Seamless integration connect effortlessly with accounting software to automate expense reporting and reduce manual errors • Virtual and physical cards utilize both virtual cards for online purchases and physical cards for offline expenses adapting to various business needs • Multi-level approvals implement approval workflows to maintain oversight and authorize spending before transactions occur • Ideal for SaaS subscriptions cloud services marketing expenses as these can only be made with credit card • Outstanding modern UI glassmorphic</p></div>
              <div className="rounded-xl bg-blue-500/10 border border-blue-500/20 p-4"><p className="font-semibold">Real-time Expense Tracking • Monitor All Transactions Instantly for Improved Visibility and Timely Reconciliation • Seamless Integration Connect Effortlessly with Accounting Software to Automate Expense Reporting and Reduce Manual Errors • Virtual and Physical Cards • Multi-level Approvals • Outstanding • Beyond RazorpayX • Ethiopia Business Banking Core P0</p><p className="mt-2 text-[11px]">Corporate Cards Collateral-free Credit Cards Virtual + Physical Dynamic Limit up to 2Cr ETB Equivalent Up to 45-50 Day Interest-free Period 2.5% Forex Markup Flat 1% Cashback 30% Off SaaS Custom Spending Controls Real-time Expense Tracking Multi-level Approvals Virtual & Physical Add-on Cards Unlimited Minimum Documentation No Joining Fee Ideal for SaaS Subscriptions Cloud Services Marketing Expenses • Outstanding per RazorpayX • Up to ₹2Cr unsecured credit up to 45–50 day interest-free period 2.5% forex markup flat 1% cashback 30% off on 500+ SaaS/marketing tools • Custom spending controls set specific limits and controls on card usage to prevent overspending • Real-time expense tracking monitor all transactions instantly • Seamless integration connect effortlessly with accounting software to automate expense reporting • Virtual and physical cards • Multi-level approvals • Ideal for SaaS subscriptions cloud services marketing expenses • Outstanding modern UI glassmorphic • Recharts bar principal vs interest pie deductions • Beyond RazorpayX • Ethiopia Business Banking Core P0</p></div>
            </div>

            <div className="mt-6 rounded-xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Issue Corporate Card • Virtual + Physical • Both • Card Network Visa • Cardholder Meron Hailu Engineering • Credit Limit 2000000 ETB • Available Credit 2000000 • Daily Limit 50000 Monthly 500000 • Category Restrictions SaaS Cloud Marketing • Cashback 1% Forex Markup 2.5% Interest Free Days 45 • Addon No • Parent Card ID null • Outstanding per RazorpayX</h4>
              <div className="mt-3 grid grid-cols-4 gap-3 text-xs">
                <div><label className="text-muted-foreground">Cardholder Name • Meron Hailu • Engineering</label><input placeholder="Meron Hailu • Engineering" className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div><label className="text-muted-foreground">Cardholder Email • meron@apextrading.et</label><input placeholder="meron@apextrading.et" className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div><label className="text-muted-foreground">Card Type • Virtual + Physical • Both • Card Network Visa • Card Network Visa/Mastercard/Verve/EthSwitch</label><select className="mt-1 w-full rounded-xl border h-9 px-3"><option>virtual • Virtual Card • Online Purchases • Instant Issuance</option><option>physical • Physical Card • Offline Expenses • Delivery Required</option><option>both • Both Virtual + Physical • Outstanding per RazorpayX</option></select></div>
                <div><label className="text-muted-foreground">Card Network • Visa/Mastercard/Verve/EthSwitch • Visa</label><select className="mt-1 w-full rounded-xl border h-9 px-3"><option>visa • Visa • International • 2.5% Forex Markup</option><option>mastercard • Mastercard • International</option><option>verve • Verve • Local Nigeria</option><option>ethswitch • EthSwitch • Local Ethiopia • EthSwitch QR</option></select></div>
                <div><label className="text-muted-foreground">Credit Limit • ETB • Up to 2Cr ETB Equivalent • 2000000 • 20L-2Cr INR in RazorpayX India</label><input type="number" defaultValue={2000000} className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div><label className="text-muted-foreground">Daily Limit • ETB • 50000 • Monthly Limit • 500000</label><input type="number" defaultValue={50000} className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div><label className="text-muted-foreground">Category Restrictions • SaaS Cloud Marketing • Allowed Categories • Blocked Merchants</label><input placeholder="SaaS, Cloud, Marketing • Allowed Categories • Blocked Merchants []" className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div className="flex items-end gap-2"><button className="rounded-xl bg-primary text-white h-9 px-6">Issue Corporate Card • Virtual + Physical • Both • Card Network Visa • Daily Limit 50000 Monthly 500000 • Category Restrictions SaaS Cloud Marketing • Cashback 1% Forex Markup 2.5% Interest Free Days 45 • Addon No • Parent Card ID null • Outstanding per RazorpayX • Up to 2Cr unsecured credit up to 45–50 day interest-free period 2.5% forex markup flat 1% cashback 30% off on 500+ SaaS/marketing tools • Custom spending controls • Real-time expense tracking • Seamless integration • Virtual and physical cards • Multi-level approvals • Ideal for SaaS subscriptions cloud services marketing expenses</button></div>
              </div>
              <p className="mt-3 text-[11px] text-muted-foreground">Logic: MerchantID CurrentAccountID CardNumberMasked ****1234 CardNumberHash sha256 hash last4 CardType virtual/physical/both CardNetwork visa/mastercard/verve/ethswitch CardholderName CardholderEmail Status ordered/active/blocked/expired/cancelled/suspended CreditLimit 2000000 up to 2Cr ETB equivalent 20L-2Cr INR in RazorpayX India AvailableCredit CreditLimit DailyLimit MonthlyLimit CategoryRestrictions [] SaaS Cloud Marketing SpendingControls jsonb daily_limit monthly_limit allowed_categories blocked_merchants cashback_percent 1% flat forex_markup_percent 2.5% interest_free_days 45 up to 45-50 day interest-free period is_addon parent_card_id unlimited add-on cards created_by approved_by approved_at expiry_month year created_at updated_at • Outstanding modern UI glassmorphic • Receipt preview thumbs • Hash integrity • Progress donut • Outstanding Modern • Mercury/Linear inspiration • Glassmorphic • QR Code Generator for Payout Links EthSwitch Interoperable QR Standard Spec • Scan & Pay Camera Permission Outstanding Dialog Overlay Rounded 260 Guides Corner Brackets Pulse Green + Supports FaydaEncode Offline QR + EthSwitch QR + Vibration for Merchant Scanning Payout Link QR to Pay • RazorpayX Parity • P0 • Outstanding Modern UI Glassmorphic</p>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
