# ApexPay Makefile - senior engineer best practices
# Usage: make test, make lint, make migrate-up, make k6-smoke

GO := go
PKG := ./services/api/internal/...

.PHONY: test ledger-test test-race lint gosec migrate-up migrate-down k6-smoke fmt

test:
	$(GO) test $(PKG) -v -count=1 -race -cover

ledger-test:
	$(GO) test ./services/api/internal/ledger -v -count=1 -run TestLedger -cover
	$(GO) test ./services/api/internal/ledger -v -count=1 -run TestPayroll -cover
	$(GO) test ./services/api/internal/ledger -run TestLedgerBalancedProperty_10k -count=1 -v

test-race:
	$(GO) test ./services/api/internal/... -race -count=1

fmt:
	$(GO) fmt $(PKG)
	gofmt -w .

lint:
	golangci-lint run ./services/api/...
	# no float money custom lint
	! grep -R "float64.*amount\|amount.*float64" --include="*.go" services/api/internal/payment services/api/internal/ledger services/api/internal/refund services/api/internal/payout services/api/internal/payroll || (echo "Float money found! Use decimal.Decimal per DATABASE" && exit 1)

gosec:
	gosec ./services/api/...

migrate-up:
	goose -dir db/migrations postgres "postgres://apexpay:apexpay_dev@localhost:5432/apexpay?sslmode=disable" up

migrate-down:
	goose -dir db/migrations postgres "postgres://apexpay:apexpay_dev@localhost:5432/apexpay?sslmode=disable" down

k6-smoke:
	k6 run scripts/k6/smoke.js

rag-test:
	cd services/rag && python -m pytest tests/ -v
	cd services/rag && python -m app.eval

# Generate mocks
gen-mocks:
	mockgen -source=services/api/internal/onboarding/service.go -destination=services/api/internal/onboarding/mock_repo.go
	mockgen -source=services/api/internal/fayda/service.go -destination=services/api/internal/fayda/mock_repo.go
