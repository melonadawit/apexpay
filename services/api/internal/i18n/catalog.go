package i18n

// Locale is a supported user language.
type Locale string

const (
	LocaleEnglish Locale = "en"
	LocaleAmharic Locale = "am"
)

// DefaultLocale is English unless a user explicitly chooses Amharic.
const DefaultLocale = LocaleEnglish

// Normalize coerces a raw preference string into a valid Locale (English on unknown/empty).
func Normalize(s string) Locale {
	switch Locale(s) {
	case LocaleAmharic:
		return LocaleAmharic
	default:
		return LocaleEnglish
	}
}

// IsValid reports whether s is a supported locale.
func IsValid(s string) bool {
	return s == string(LocaleEnglish) || s == string(LocaleAmharic)
}

// Catalog holds localized message strings for both languages.
// Messages that were previously concatenated EN+AM with "•" are split into clean,
// single-language entries keyed by a stable identifier.
type Catalog struct {
	en map[string]string
	am map[string]string
}

// New returns the built-in catalog.
func New() *Catalog {
	return &Catalog{
		en: enMessages,
		am: amMessages,
	}
}

// Get returns the localized message for a key in the given locale. Falls back to the key
// itself if missing, so translation gaps degrade gracefully rather than crash.
func (c *Catalog) Get(locale Locale, key string) string {
	var m map[string]string
	switch locale {
	case LocaleAmharic:
		m = c.am
	default:
		m = c.en
	}
	if s, ok := m[key]; ok {
		return s
	}
	if en, ok := c.en[key]; ok {
		return en
	}
	return key
}

// En exposes the English messages (for tests / tooling).
func (c *Catalog) En() map[string]string { return c.en }

var enMessages = map[string]string{
	// Auth / preferences
	"language_preference_updated":  "Language preference updated",
	"language_preference_required": "language preference must be 'en' or 'am'",
	"language":                     "Language",
	"english":                      "English",
	"amharic":                      "Amharic",

	// Payroll claim approval (previously mixed EN+AM with "•")
	"claim_created":          "Expense claim submitted",
	"claim_manager_approved": "Manager approved. Next step: finance approval.",
	"claim_finance_approved": "Finance approved. The reimbursement will be paid via the next payroll run.",
	"claim_requires_finance": "Finance approval pending.",
	"claim_paid":             "Reimbursement paid.",

	// Notifications
	"current_account_opened": "Current account opened",
	"invoice_sent":           "Invoice sent",

	// Payroll — clean, professional API success copy (no internals/marketing).
	"calendar_locked":           "Payroll calendar locked.",
	"calendar_unlocked":         "Payroll calendar unlocked.",
	"leave_approved":            "Leave request approved.",
	"leave_rejected":            "Leave request rejected.",
	"loan_emi_schedule":         "Loan repayment schedule generated.",
	"tax_cert_generated":        "Annual income tax certificate generated.",
	"payroll_register_ready":    "Payroll register is ready.",
	"payroll_register_need_run": "Select a payroll run to view the payroll register.",
	"variance_report_ready":     "Payroll variance report is ready.",
	"email_payslip_subject":     "Your payslip for %s",
	"email_compliance_subject":  "Compliance report for %02d/%d",
	"email_magiclink_subject":   "Your employee portal access link",
	"magiclink_sent":            "Access link sent via %s to %s.",

	// Onboarding / Fayda / verification
	"fayda_otp_sent":      "One-time password sent to your Fayda-registered phone.",
	"kyc_submitted":       "KYC submitted for compliance review.",
	"verification_queued": "Verification queued.",
	"otp_sent":            "One-time password sent.",

	// Payroll runs / reports (clean copy)
	"payroll_calculated":     "Payroll calculated.",
	"payroll_disbursed":      "Payroll disbursed.",
	"payslip_generated":      "Payslip generated.",
	"payslips_zip_ready":     "Payslips download ready.",
	"bank_disbursal_pending": "Bank disbursal file not yet generated.",
	"cost_center_report":     "Cost center report is ready.",
	"resend_queued":          "Webhook resend queued.",

	// Apex Assistant framing
	"assistant_overview":     "Here's your business overview: ",
	"assistant_found":        "Here's what I found for you: ",
	"assistant_no_results":   "I couldn't find anything for that. Try asking about payments, invoices, inventory, cash position, or (as an employee) your payslip, leave balance, or expense claims.",
	"assistant_unavailable":  "temporarily unavailable",
	"no_recent_payments":     "No recent payments yet.",
	"you_have_payments":      "You have %d recent payments.",
	"no_overdue_invoices":    "No overdue invoices.",
	"invoices_overdue":       "%d invoices are overdue totalling %s ETB.",
	"inventory_summary":      "You have %d products in inventory, %d below reorder threshold.",
	"no_loans":               "No loans on record.",
	"loans_outstanding":      "You have %d loans with %s ETB outstanding.",
	"ytd_pay":                "Your YTD gross is %s, net is %s, tax %s ETB.",
	"no_leave_balance":       "No leave balance on record for this year.",
	"leave_types_count":      "You have %d leave types on record this year.",
	"annual_leave_remaining": "You have %s annual leave days remaining this year.",
	"no_expense_claims":      "You have no expense claims on record.",
	"expense_claims_count":   "You have %d expense claims (%d pending) totalling %s ETB.",
	"cash_position":          "Cash position is %s ETB.",
	"cash_position_forecast": "Cash position is %s ETB. Net 90-day cash flow is forecast at %s ETB.",
	"statement_total":        "%s: %s %s.",
	"summary_tpv":            "Today's TPV is %s ETB across %s transactions.",

	// Generic
	"ok":      "OK",
	"success": "Success",
	"error":   "Error",
}

var amMessages = map[string]string{
	"language_preference_updated":  "የቋንቋ ምርጫ ተዘምኗል",
	"language_preference_required": "የቋንቋ ምርጫ 'en' ወይም 'am' መሆን አለበት",
	"language":                     "ቋንቋ",
	"english":                      "እንግሊዝኛ",
	"amharic":                      "አማርኛ",

	"claim_created":          "የወጪ ጥያቄ ቀርቧል",
	"claim_manager_approved": "ሥራ አስኪያጁ አጽድቋል። ቀጣይ ደረጃ፡ የፋይናንስ ማጽደቅ።",
	"claim_finance_approved": "ፋይናንስ አጽድቋል። ክፍያው በሚቀጥለው የደመወዝ ዙር ይከፈላል።",
	"claim_requires_finance": "የፋይናንስ ማጽደቅ በመጠባበቅ ላይ።",
	"claim_paid":             "ክፍያ ተፈጽሟል።",

	"current_account_opened": "የገንዘብ ሒሳብ ተከፍቷል",
	"invoice_sent":           "ደረሰኝ ተልኳል",

	"calendar_locked":           "የደመወዝ ካሌንደር ተቆልፏል።",
	"calendar_unlocked":         "የደመወዝ ካሌንደር ተከፍቷል።",
	"leave_approved":            "የፈቃድ ጥያቄ ጸድቋል።",
	"leave_rejected":            "የፈቃድ ጥያቄ ውድቅ ተደርጓል።",
	"loan_emi_schedule":         "የብድር ክፍያ መርሃ ግብር ተዘጋጅቷል።",
	"tax_cert_generated":        "የዓመታዊ የገቢ ግብር ሰርቲፊኬት ተዘጋጅቷል።",
	"payroll_register_ready":    "የደመወዝ መዝገብ ዝግጁ ነው።",
	"payroll_register_need_run": "እባክዎ የደመወዝ መዝገቡን ለማየት የደመወዝ ዙር ይምረጡ።",
	"variance_report_ready":     "የደመወዝ ልዩነት ሪፖርት ዝግጁ ነው።",
	"email_payslip_subject":     "የደመወዝ ደረሰኝዎ ለ%s",
	"email_compliance_subject":  "የተገዢነት ሪፖርት ለ%02d/%d",
	"email_magiclink_subject":   "የሰራተኛ መግቢያ አገናኝዎ",
	"magiclink_sent":            "የመግቢያ አገናኝ በ%s ወደ %s ተልኳል።",

	"fayda_otp_sent":      "የአንድ ጊዜ የይለፍ ቃል በFayda የተመዘገበ ስልክዎ ተልኳል።",
	"kyc_submitted":       "KYC ለተገዢነት ግምገማ ቀርቧል።",
	"verification_queued": "ማረጋገጫ በመስመር ላይ ተቀምጧል።",
	"otp_sent":            "የአንድ ጊዜ የይለፍ ቃል ተልኳል።",

	"payroll_calculated":     "ደመወዝ ተሰልቷል።",
	"payroll_disbursed":      "ደመወዝ ተከፋፍሏል።",
	"payslip_generated":      "የደመወዝ ደረሰኝ ተዘጋጅቷል።",
	"payslips_zip_ready":     "የደመወዝ ደረሰኞች ማውረድ ዝግጁ ነው።",
	"bank_disbursal_pending": "የባንክ ክፍያ ፋይል ገና አልተዘጋጀም።",
	"cost_center_report":     "የወጪ ማዕከል ሪፖርት ዝግጁ ነው።",
	"resend_queued":          "የድር ማንቂያ ዳግም መላክ በመስመር ላይ ተቀምጧል።",

	"assistant_overview":     "የንግድዎ አጠቃላይ እይታ እነሆ፡ ",
	"assistant_found":        "ለእርስዎ ያገኘሁት እነሆ፡ ",
	"assistant_no_results":   "ለዚያ ምንም አላገኘሁም። ስለ ክፍያዎች፣ ደረሰኞች፣ ክምችት፣ የገንዘብ አቅም፣ ወይም (እንደ ሰራተኛ) ስለ ደመወዝ፣ የፈቃድ ቀሪ እና የወጪ ጥያቄዎች ይጠይቁ።",
	"assistant_unavailable":  "ለጊዜው አይገኝም",
	"no_recent_payments":     "እስካሁን ምንም ክፍያ የለም።",
	"you_have_payments":      "%d የቅርብ ጊዜ ክፍያዎች አሉዎት።",
	"no_overdue_invoices":    "የዘገዩ ደረሰኞች የሉም።",
	"invoices_overdue":       "%d የዘገዩ ደረሰኞች በድምሩ %s ETB ናቸው።",
	"inventory_summary":      "%d ምርቶች በክምችት ውስጥ አሉዎት፣ %d ደግሞ ከድጋሚ ማዘዣ ገደብ በታች ናቸው።",
	"no_loans":               "ምንም የብድር መዝገብ የለም።",
	"loans_outstanding":      "%d ብድሮች በ%s ETB ቀሪ አሉዎት።",
	"ytd_pay":                "የዓመትዎ ጠቅላላ ገቢ %s፣ የተጣራ %s፣ ግብር %s ETB ነው።",
	"no_leave_balance":       "ለዚህ ዓመት ምንም የፈቃድ ቀሪ የለም።",
	"leave_types_count":      "በዚህ ዓመት %d አይነት የፈቃድ መዝገቦች አሉዎት።",
	"annual_leave_remaining": "በዚህ ዓመት %s የዓመታዊ ፈቃድ ቀሪ ቀናት አሉዎት።",
	"no_expense_claims":      "ምንም የወጪ ጥያቄ መዝገብ የለዎትም።",
	"expense_claims_count":   "%d የወጪ ጥያቄዎች (%d በመጠባበቅ ላይ) በድምሩ %s ETB አሉዎት።",
	"cash_position":          "የገንዘብ አቅም %s ETB ነው።",
	"cash_position_forecast": "የገንዘብ አቅም %s ETB ነው። የ90 ቀን የገንዘብ ፍሰት ትንበያ %s ETB ነው።",
	"statement_total":        "%s፡ %s %s።",
	"summary_tpv":            "የዛሬው TPV %s ETB በ%s ግብይቶች ነው።",

	"ok":      "እሺ",
	"success": "ተሳክቷል",
	"error":   "ስህተት",
}
