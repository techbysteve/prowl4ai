# `prowl4ai` Roadmap

This roadmap tracks `prowl4ai` as a Go-native crawler inspired by `crawl4ai` (Python), with priorities based on the current Python codebase (`crawl4ai` v0.8.0).

## Current State (as of February 22, 2026)

Implemented in `prowl4ai`:

- [x] Single URL crawl command: `prowl4ai crawl <url>`
- [x] Playwright-backed navigation and rendered HTML capture
- [x] Wait strategy support in run config (`wait_until`, optional `wait_for`)
- [x] Readability-based extraction of cleaned HTML + basic article metadata
- [x] Markdown conversion from cleaned HTML
- [x] Structured JSON result with status code, headers, and redirected URL
- [x] Browser and run config defaults

Not implemented yet:

- [ ] Automated tests (unit/integration/e2e)
- [ ] Link extraction and classification
- [ ] Multi-URL crawling and dispatcher controls
- [ ] Cache modes and resumable deep crawl flows
- [ ] Hooks/plugin system
- [ ] Service/API mode

## Parity Targets From Python `crawl4ai`

Python repo capabilities we want to progressively match in Go:

- Multi-URL orchestration (`arun_many`) with bounded concurrency
- Rate limiting + retry/backoff controls
- Deep crawl strategies (BFS/DFS/Best-first) with filters/scorers
- Broader extraction options (links/media/structured selectors)
- Cache-aware crawling and resume/crash-recovery behavior
- Stronger operational tooling (monitoring, service APIs)

`prowl4ai` does not need a one-to-one port of every Python API. Priority is practical parity for reliability, crawl scale, and output quality.

## Phase 1 - MVP Hardening

Goal: make current single-page crawl deterministic and testable.

- [ ] Add unit tests for `internal/config` defaults and validation
- [ ] Add extraction pipeline tests (`clean -> markdown`) with golden files
- [ ] Add adapter/service tests for success, timeout, and invalid URL paths
- [ ] Normalize/validate URL input before navigation
- [ ] Introduce stage-specific typed errors (`validation`, `navigation`, `extraction`)
- [ ] Add fixture pages for deterministic local test runs

## Phase 2 - CLI and Output Maturity

Goal: expose existing internal controls and improve machine ergonomics.

- [ ] Expose `--wait-until`, `--wait-for`, and `--wait-for-timeout` on CLI
- [ ] Add output toggles (`--include-html`, `--include-cleaned-html`, `--metadata-only`)
- [ ] Add `--quiet` and `--verbose` behavior aligned with config
- [ ] Add stable exit codes by failure category
- [ ] Add `--version` and improved subcommand help text

## Phase 3 - Crawl4AI-Core Parity (Practical)

Goal: close the biggest feature gaps with Python `crawl4ai` for common workflows.

- [ ] Implement link extraction pipeline and include in `CrawlResult`
- [ ] Normalize links against redirected/final URL
- [ ] Classify internal vs external links, deduplicate, stable-sort
- [ ] Add selector-based content targeting (basic CSS selector extraction)
- [ ] Add optional file/raw-input crawl sources (`file://` / raw HTML mode)

## Phase 4 - Multi-URL Execution

Goal: move from single-crawl utility to controlled crawl runner.

- [ ] Add `crawl-many` command accepting URL list/file/STDIN
- [ ] Add worker-pool concurrency controls
- [ ] Add retry policy with status/error classification
- [ ] Add per-domain rate limiter
- [ ] Add run summary metrics (`success`, `failed`, `duration`, `retries`)
- [ ] Add JSONL output sink for large runs

## Phase 5 - Deep Crawl and Recovery

Goal: support long-running site traversal with resumability.

- [ ] Introduce deep-crawl strategies (BFS first, then DFS/best-first)
- [ ] Add crawl limits (max pages, depth, domain boundaries)
- [ ] Add filter/scorer interfaces for URL prioritization
- [ ] Persist crawl state for resume after interruption
- [ ] Add state-change callbacks/events for observability

## Phase 6 - Runtime and Integration Layer

Goal: make `prowl4ai` embeddable in systems, not only CLI-driven.

- [ ] Add optional HTTP service mode (`/crawl`, `/crawl-many`)
- [ ] Add hook points (pre-nav, post-fetch, post-extract)
- [ ] Add structured logging + metrics (timings, counters, error taxonomy)
- [ ] Add auth/session primitives (cookies, storage state, headers)
- [ ] Add proxy configuration parity with runtime controls

## Phase 7 - Release Readiness

Goal: establish repeatable quality and contributor velocity.

- [ ] CI for test/lint/build on Linux/macOS
- [ ] Versioning and changelog process (semver)
- [ ] Cross-platform Playwright setup checks
- [ ] Performance and memory benchmarks
- [ ] Security checks for dependency and input handling
- [ ] Contributor guide and architecture docs

## Backlog (After Phase 7)

- [ ] LLM-oriented chunking helpers for RAG workflows
- [ ] Pluggable structured extraction modules
- [ ] Snapshot bundles (HTML + Markdown + metadata artifacts)
- [ ] First-party Docker image and deployment docs

## Success Criteria for `v1.0.0`

`prowl4ai` reaches practical v1 when:

- single and multi-URL crawls are reliable under CI-tested scenarios,
- output schema is stable and documented,
- deep crawl resume works for interrupted long runs,
- and CLI/service modes are both production-usable.
