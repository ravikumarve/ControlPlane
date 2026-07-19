# ControlPlane AI — E2E Testing Pipeline Strategy

**Version**: 1.0 | **Date**: 2026-07-19 | **Status**: Draft

## 1. Testing Philosophy

Security products must be tested more rigorously than standard software. A misconfigured rule or latent injection bypass in an API gateway results in a data leak; the same defect in an agent security proxy allows autonomous code execution through an MCP tool call. ControlPlane AI's testing strategy is built on four foundational principles:

- **Test the three paths**: Every test suite must exercise ALLOW, BLOCK, and HITL (human-in-the-loop) decision paths. A proxy that correctly blocks known-dangerous tools but silently allows an unhandled edge case is not secure.
- **Audit integrity is non-negotiable**: The HMAC chain must always be verifiable. Any test run that produces audit entries must conclude with a chain verification step. If the chain is broken, the test fails — regardless of whether the functional assertion passed.
- **Integration tests are the source of truth**: Unit tests verify that individual components behave correctly in isolation, but only integration tests prove that the proxy, policy engine, audit logger, and MCP protocol handling work together as a system. When a unit test and an integration test disagree, the integration test wins.
- **Determinism over coverage**: A flaky E2E test is worse than no test. All integration tests must be fully deterministic — no network-dependent timeouts, no race-prone goroutine scheduling assumptions, no dependency on external services.

## 2. Test Pyramid for Security Proxy

```
                    ┌───────────────────────┐
                    │   E2E Integration     │
                    │  (few, high-value)    │
                    │                       │
                    │ full proxy round-trip │
                    └──────────┬────────────┘
                               │
                    ┌──────────▼────────────┐
                    │  Component Integration │
                    │                       │
                    │ policy + audit        │
                    │ proxy + schema        │
                    │ proxy + rate limiter  │
                    └──────────┬────────────┘
                               │
                    ┌──────────▼────────────┐
                    │     Unit Tests        │
                    │                       │
                    │ policy engine         │
                    │ audit logger          │
                    │ config validation     │
                    │ rate limiter          │
                    │ glob matcher          │
                    │ HMAC chain            │
                    │ JSON-RPC parser       │
                    │ schema pinner         │
                    │ injection detector    │
                    └───────────────────────┘
```

### 2.1 Distribution Guidelines

| Layer | Count Target | Run Frequency | Max Duration |
|-------|-------------|---------------|--------------|
| Unit | 200+ | Every commit | 10s |
| Component Integration | 30-50 | Every PR | 30s |
| E2E Integration | 10-15 | Every PR + pre-release | 60s |

## 3. E2E Test Topology

```
┌──────────────┐   JSON-RPC    ┌──────────────┐   JSON-RPC    ┌──────────────┐
│  Test Client  │ ────────────▶│  Proxy Under  │ ────────────▶│  Test MCP    │
│  (go test)    │ ◀────────────│  Test         │ ◀────────────│  Server      │
│               │              │  (real impl)  │              │  (mock)      │
└──────────────┘              └──────────────┘              └──────────────┘
       │                             │                              │
       │                             │                              │
       ▼                             ▼                              ▼
┌──────────────┐              ┌──────────────┐
│  Test Config  │              │  Audit Log   │
│  (temp YAML)  │              │  (inspected  │
│               │              │   post-run)  │
└──────────────┘              └──────────────┘
```

### 3.1 Topology Rules

- **Test client** sends real JSON-RPC 2.0 messages over a `net.Conn` (TCP mode) or `io.Pipe` (stdio mode). It is a `go test` function that drives the scenario.
- **Proxy under test** is the real `mcp-guard` binary compiled with `go build` or invoked via `internal/proxy` package directly. No mocks, no stubs for the proxy itself.
- **Test MCP server** is a lightweight Go `net/http` or raw TCP server that implements the MCP protocol just enough to respond to `initialize`, `tools/list`, and `tools/call`. It may return configurable responses to exercise edge cases (e.g., slow responses, oversized payloads, malformed JSON).
- **Test config** is a temporary YAML file written to a temp directory before each test run. The proxy is pointed at this config file via `--config`.
- **Audit log** is written to a temp file whose path is set in the test config. After each test, the audit log is read, parsed, and assertions are run against its entries.

### 3.2 Startup Sequence per Test

1. Generate temp directory with test config YAML and HMAC key
2. Start Test MCP Server on a random available port
3. Start Proxy Under Test, pointed at Test MCP Server
4. Wait for proxy readiness (poll `GET /health` or connect to listener)
5. Run test scenario (send JSON-RPC messages, collect responses)
6. Shutdown proxy (SIGTERM)
7. Shutdown Test MCP Server
8. Load and verify audit log from temp path
9. Assert on responses and audit entries

## 4. Test Scenarios

All scenarios assume a standard test policy file with the following rules loaded:

```yaml
policies:
  - name: allow-reads
    match:
      identity: "*"
      tools: ["read_*", "search_*", "get_*", "list_*"]
    action: allow
    rate_limit: 50/min

  - name: payment-hitl
    match:
      identity: "payment-bot"
      tools: ["execute_payout"]
    action: hitl
    constraints:
      max_amount: 1000

  - name: block-deletes
    match:
      tools: ["delete_*", "drop_*", "exec"]
    action: block
    alert: true
```

| # | Scenario | Steps | Expected |
|---|----------|-------|----------|
| 1 | **Handshake pass-through** | Send `initialize` request with protocol version and client info | Response from Test MCP Server, no policy evaluation logged |
| 2 | **tools/list discovery** | Send `tools/list` request | Tools definition returned from Test MCP Server; schema pin computed if enabled |
| 3 | **ALLOW path** | Send `tools/call` with tool name `read_file` (matches `allow-reads` policy) | Request forwarded to Test MCP Server; response returned to client; audit entry with `decision: allow` |
| 4 | **BLOCK path** | Send `tools/call` with tool name `delete_dataset` (matches `block-deletes` policy) | Error response `-32000` from proxy (not from server); audit entry with `decision: block` and `reason: "blocked by policy: block-deletes"` |
| 5 | **HITL path** | Send `tools/call` with tool name `execute_payout` and identity `payment-bot` (matches `payment-hitl` policy) | Pending approval response from proxy; webhook payload sent to configured endpoint; audit entry with `decision: pending` |
| 6 | **Default deny** | Send `tools/call` with tool name `unknown_tool_xyz` (no matching policy) | Blocked with error; audit entry with `decision: block` and `reason: "no matching policy"` |
| 7 | **Rate limit exceeded** | Send 51 `tools/call` requests for `read_file` in rapid succession (exceeds `50/min` limit) | First 50 allowed; 51st blocked with rate-limit error; audit entries show 50 `allow` + 1 `block` |
| 8 | **Injection detected** | Send `tools/call` with params containing shell injection pattern (e.g., `"; rm -rf /"`) | Blocked with error containing `injection_detected` reason; audit entry with `decision: block` and reason referencing injection scan |
| 9 | **Schema drift detection** | First connect: `tools/list` returns tool set A (pinned). Second connect: `tools/list` returns tool set B (different hash) with mode `block` | First connection succeeds. Second connection blocked; error indicates schema drift |
| 10 | **Audit integrity** | Run 5 mixed decisions (2 allow, 2 block, 1 hitl), then verify audit log | 5 entries in audit log; `mcp-guard logs --verify` or equivalent HMAC chain computation succeeds. `prev_hmac` of entry N equals `hmac` of entry N-1 |
| 11 | **Config reload** | Start proxy with policy A; send request (blocked). Update YAML with policy B that allows the tool; signal SIGHUP | First request blocked under policy A. After reload, same request allowed under policy B |
| 12 | **HITL approval flow** | Send HITL-triggering request; approve via callback endpoint | Request transitions from `pending` to `approved`; actual MCP call is executed; response returned to client |
| 13 | **HITL denial flow** | Send HITL-triggering request; deny via callback endpoint (or timeout) | Request transitions from `pending` to `denied`; MCP call is never executed; error returned to client |
| 14 | **Concurrent identity isolation** | Send requests from two identities simultaneously; identity A matches `payment-hitl`, identity B matches `allow-reads` | Identity A gets HITL-pending response. Identity B gets normal ALLOW response. No cross-contamination |

## 5. CI Pipeline Integration (GitHub Actions)

### 5.1 On Push / PR to `main`

```yaml
name: ci

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - run: go vet ./...
      - run: go install honnef.co/go/tools/cmd/staticcheck@latest
      - run: staticcheck ./...

  unit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - run: go test -race -count=1 -coverprofile=coverage.out ./...
      - run: go tool cover -func=coverage.out
      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | tr -d '%')
          if [ "$COVERAGE" -lt 80 ]; then
            echo "Coverage $COVERAGE% is below 80% threshold"
            exit 1
          fi

  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, darwin]
        goarch: [amd64, arm64]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - run: GOOS=${{ matrix.goos }} GOARCH=${{ matrix.goarch }} go build -o mcp-guard ./cmd/mcp-guard
      - name: Check binary size
        run: |
          SIZE=$(stat -c%s mcp-guard)
          if [ "$SIZE" -gt 15728640 ]; then
            echo "Binary size $SIZE bytes exceeds 15MB limit"
            exit 1
          fi

  e2e:
    runs-on: ubuntu-latest
    needs: [lint, unit, build]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - name: Build test binary
        run: go test -c -tags=e2e ./test/e2e/ -o e2e.test
      - name: Run E2E integration tests
        run: ./e2e.test -test.v -test.parallel=2
        env:
          MCP_GUARD_HMAC_KEY: "test-hmac-key-for-ci"
      - name: Verify audit integrity
        run: go run ./cmd/mcp-guard logs --verify --path /tmp/mcp-guard-e2e-audit.jsonl

  benchmark:
    runs-on: ubuntu-latest
    needs: [e2e]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - name: Run benchmarks
        run: go test -bench=. -benchtime=1x -count=1 -run=^$ ./internal/... 2>&1 | tee bench-output.txt
      - name: Compare with baseline
        run: |
          # Compare against stored baseline; fail if regression > 10%
          go install golang.org/x/perf/cmd/benchstat@latest
          benchstat baseline.txt bench-output.txt | tee benchstat-result.txt
          if grep -q "+[1-9][0-9]%\|+100%" benchstat-result.txt; then
            echo "Performance regression detected"
            exit 1
          fi
```

### 5.2 On Tag Push (`v*`)

All of the above jobs run first. On success, a release job executes:

```yaml
release:
  runs-on: ubuntu-latest
  needs: [lint, unit, build, e2e, benchmark]
  if: startsWith(github.ref, 'refs/tags/v')
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: "1.22"
    - name: Cross-platform release build
      run: |
        mkdir -p dist
        GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/mcp-guard-linux-amd64 ./cmd/mcp-guard
        GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/mcp-guard-linux-arm64 ./cmd/mcp-guard
        GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/mcp-guard-darwin-amd64 ./cmd/mcp-guard
        GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/mcp-guard-darwin-arm64 ./cmd/mcp-guard
    - name: Generate checksums
      run: |
        cd dist
        sha256sum mcp-guard-* > checksums.txt
    - name: Create GitHub Release
      uses: softprops/action-gh-release@v2
      with:
        files: dist/*
        generate_release_notes: true
```

## 6. Performance Benchmarks

Benchmarks are defined as Go benchmark functions in `internal/proxy/bench_test.go` and `internal/policy/bench_test.go`. They run on every PR and are compared against a committed baseline (`test/bench/baseline.txt`).

| Benchmark | Target | Measurement Method |
|-----------|--------|--------------------|
| **Proxy latency overhead (ALLOW path)** | < 100us p99 | Full round-trip through proxy vs. direct-to-server, measured with `go test -bench` and `-benchtime=1000x` |
| **Block path latency** | < 50us | Time from request receipt to error response; no upstream call made |
| **Injection scan latency** | < 500us for typical payloads | 1KB request body with embedded injection patterns; measured over 1000 iterations |
| **Rate limit check** | < 10us | Token bucket consume operation for an already-warm limiter |
| **Audit log write** | < 50us per entry | Sequential write + HMAC computation to temp file; 10K entries |
| **Concurrent connections** | 1000 simultaneous | `net.Conn` fan-out test; all 1000 connections maintain open JSON-RPC sessions |

### 6.1 Benchmark Environment

All benchmarks must be run on the same CI runner type (`ubuntu-latest`) with CPU governor set to `performance`. Results are cached and compared against the stored baseline. A regression of more than 10% in any benchmark triggers a CI failure.

## 7. Testing Tools & Frameworks

| Tool | Purpose | Usage |
|------|---------|-------|
| `go test` | Primary test runner | All unit, component, and E2E tests |
| `testing/quick` | Property-based testing | Policy engine edge cases (fuzz argument combinations) |
| `testify/assert` | Assertion helpers | Structured test assertions (if adopted by project) |
| `-race` flag | Race detector | Enabled for all unit and component tests |
| `-count=1` | Cache bypass | Ensures every test run executes fresh |
| `-tags=e2e` | Build tag isolation | E2E tests are gated behind `//go:build e2e` to keep `go test ./...` fast |
| `testcontainers-go` | Containerized dependencies | Optional — only if database-backed audit store is introduced post-MVP |
| Custom MCP test server | Protocol-level testing | Lightweight Go server in `test/mcp/server.go` — implements minimal MCP protocol for test scenarios |

### 7.1 Test Package Layout

```
test/
├── e2e/
│   ├── e2e_test.go          # E2E test suite (build tag: e2e)
│   ├── suite.go             # Shared topology setup/teardown
│   └── fixtures/            # Test YAML configs per scenario
├── mcp/
│   └── server.go            # Test MCP server implementation
├── bench/
│   └── baseline.txt         # Committed benchmark numbers
└── fixtures/
    ├── policies/
    │   ├── allow-reads.yaml
    │   ├── block-deletes.yaml
    │   └── hitl-payment.yaml
    └── audit/
        ├── valid-chain.jsonl
        └── tampered-chain.jsonl
```

## 8. Quality Gates

Every PR and release must pass the following quality gates. Failure of any gate blocks merge or release.

| Gate | Threshold | Enforcement |
|------|-----------|-------------|
| **Code coverage** | >= 80% | CI `go test -cover` with `-coverprofile`; checked in CI job |
| **Flaky tests** | Zero tolerance | Any E2E test that fails non-deterministically must be fixed or removed; flaky detection via `go test -count=3 -failfast` on E2E suite |
| **Three decision paths** | ALLOW, BLOCK, HITL all tested | Every PR must include or update tests for all three paths; checked via codeowners review + CI scenario enumeration |
| **Performance regression** | < 10% degradation vs baseline | `benchstat` comparison between PR branch and baseline; failure if any benchmark regresses > 10% |
| **Audit integrity** | HMAC chain must verify | Mandatory `mcp-guard logs --verify` step in E2E CI job; pre-release gate runs on the full audit log from the E2E suite |
| **Binary size** | < 15MB | `stat` check on `go build` output; enforced in build CI job |
| **Race condition** | Zero races | `-race` flag on all unit and component tests; race-free policy enforced in CI |
| **Lint** | Zero `go vet` and `staticcheck` warnings | Enforced in lint CI job |

### 8.1 Pre-Release Checklist

Before cutting a release (`v*` tag), the following manual and automated checks must pass:

1. All CI jobs pass on the target commit (lint, unit, build, E2E, benchmark)
2. Audit integrity check passes (full chain verification against E2E output)
3. No open `P0` or `P1` issues targeting this release
4. Binary size on linux/amd64 < 15MB
5. Benchmark comparison shows no regression > 10% against the previous release tag
6. Changelog updated with all changes since last release
7. At least one E2E test run with the release candidate binary (not just `go test`)

## Appendix A: Test Environment Requirements

- **CI runner**: `ubuntu-latest` (GitHub Actions)
- **Go version**: 1.22+
- **Temp directory**: Writable `/tmp` on runner — no persisted state between test runs
- **Network**: Loopback only; no external network calls in any test
- **Environment variables**: `MCP_GUARD_HMAC_KEY` must be set for all E2E tests; a fixed test key is used for deterministic HMAC output

## Appendix B: Adding a New Test Scenario

1. Create a new test config YAML in `test/e2e/fixtures/` (or reuse an existing one)
2. Add a test function in `test/e2e/e2e_test.go` following the existing pattern
3. Run locally: `go test -tags=e2e -v -count=1 ./test/e2e/`
4. Verify the audit log produced by the new test has an intact HMAC chain
5. Add the scenario to the table in Section 4 of this document
6. Ensure the test is deterministic: no sleeps, no network calls, no external state
