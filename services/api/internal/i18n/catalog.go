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

	"ok":      "እሺ",
	"success": "ተሳክቷል",
	"error":   "ስህተት",
}
