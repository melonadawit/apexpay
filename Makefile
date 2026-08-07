# ApexPay Makefile - senior engineer best practices
# Usage: make test, make lint, make migrate-up, make k6-smoke

# Go commands run only inside Docker; no Go toolchain/module cache is required locally.
DOCKER_BUILD := docker build -f deploy/docker/Dockerfile.api .

.PHONY: test ledger-test test-race lint gosec migrate-up migrate-down k6-smoke fmt

test:
	$(DOCKER_BUILD) --target test

ledger-test:
	$(DOCKER_BUILD) --target test

test-race:
	@echo "Add a dedicated Docker race-test target before enabling this command."
	@exit 1

fmt:
	@echo "Run gofmt only inside a dedicated Docker tooling target; do not install Go locally."
	@exit 1

lint:
	golangci-lint run ./services/api/...
	# no float money custom lint
	! grep -R "float64.*amount\|amount.*float64" --include="*.go" services/api/internal/payment services/api/internal/ledger services/api/internal/refund services/api/internal/payout services/api/internal/payroll || (echo "Float money found! Use decimal.Decimal per DATABASE" && exit 1)

gosec:
	gosec ./services/api/...

migrate-up:
	docker compose -f deploy/docker/docker-compose.yml run --rm migrate

migrate-down:
	@echo "Down migrations are intentionally disabled; use a disposable development database instead."
	@exit 1

k6-smoke:
	k6 run scripts/k6/smoke.js

rag-test:
	cd services/rag && python -m pytest tests/ -v
	cd services/rag && python -m app.eval

# Generate mocks
gen-mocks:
	mockgen -source=services/api/internal/onboarding/service.go -destination=services/api/internal/onboarding/mock_repo.go
	mockgen -source=services/api/internal/fayda/service.go -destination=services/api/internal/fayda/mock_repo.go
