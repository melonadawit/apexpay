package id

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

// New generates prefixed ULID: pay_01H..., merch_, kyc_, fayda_, ref_, sub_ etc
func New(prefix string) string {
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return prefix + "_" + id.String()
}

func NewMerchant() string      { return New("mer") }
func NewKYCProfile() string    { return New("kyc") }
func NewOwner() string         { return New("own") }
func NewFaydaVerification() string { return New("fayda") }
func NewDocument() string      { return New("doc") }
func NewPayment() string       { return New("pay") }
func NewRefund() string        { return New("ref") }
func NewCustomer() string      { return New("cust") }
func NewSubPlan() string       { return New("splan") }
func NewSubscription() string  { return New("sub") }
func NewBeneficiary() string   { return New("ben") }
func NewPayoutBatch() string   { return New("pbat") }
func NewPayout() string        { return New("pout") }
func NewEmployee() string      { return New("emp") }
func NewPayrollRun() string    { return New("prun") }
func NewPayrollItem() string   { return New("pitem") }
func NewLedgerBook() string    { return New("lbk") }
func NewLedgerJournal() string { return New("ljrn") }
func NewOutbox() string        { return New("outbox") }
func NewSwarmSession() string  { return New("swarm") }
func NewRAGDoc() string        { return New("rdoc") }
