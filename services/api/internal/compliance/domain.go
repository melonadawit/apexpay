package compliance

// Compliance console domain types.

type Status struct {
	MerchantID         string `json:"merchant_id"`
	OnboardingStatus   string `json:"onboarding_status"`
	KYCExpiryDate      string `json:"kyc_expiry_date,omitempty"`
	LicenseExpiry      string `json:"license_expiry,omitempty"`
	FaydaVerified      bool   `json:"fayda_verified"`
	RiskTier           string `json:"risk_tier"`
	NextERCADue        string `json:"next_erca_due,omitempty"`
	NextPensionDue     string `json:"next_pension_due,omitempty"`
	AnnualTaxFilingDue string `json:"annual_tax_filing_due,omitempty"`
	AMLDue             string `json:"aml_due,omitempty"`
	OverallStatus      string `json:"overall_status"`
	Notes              string `json:"notes,omitempty"`
}

type CheckLog struct {
	ID        string `json:"id"`
	CheckType string `json:"check_type"`
	Status    string `json:"status"`
	Detail    string `json:"detail"`
	CheckedAt string `json:"checked_at"`
}
