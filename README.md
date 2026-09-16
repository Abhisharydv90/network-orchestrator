# 🌐 Network-Aware Test Orchestrator

> Kill flaky tests. Measure latency. Retry intelligently.

## The Problem

E2E tests on remote platforms (BrowserStack, Sauce Labs) are flaky due to network latency. Fixed wait times cause tests to pass locally but fail in CI. Developers waste hours debugging ghost failures.

## The Solution

A Go CLI that measures real network latency before tests run and dynamically calculates the optimal wait time. It also evaluates test errors and only retries on transient network errors (socket hang up, ECONNRESET, timeouts), while failing immediately on real test failures.

## 🚀 Features

- 📡 **Latency Measurement** — Measures real round-trip time to your test target
- ⏱️ **Dynamic Wait Calculation** — Uses `Math.max(500, latency * 1.5)` for adaptive waits
- 🧠 **Intelligent Retry Logic** — Distinguishes network errors from real test failures
- 📄 **JSON Output** — Machine-readable for CI/CD pipelines
- ⚙️ **GitHub Action** — Automatically posts reports on every PR

## ⚙️ Quick Start

Measure latency:
orchestrator measure https://wordpress.org

Evaluate a retry decision:
orchestrator retry "socket hang up"

## 🧠 How It Works

1. **Go CLI** measures real latency via HTTP
2. **Dynamic Wait Logic** calculates `max(500ms, latency * 1.5)`
3. **Retry Engine** checks for known transient network errors
4. **GitHub Action** runs both on every PR and posts a report

## 🎯 Why This Matters

- **BrowserStack** — Their #1 customer complaint is test flakiness from network latency
- **Automattic** — Jetpack PRs show this is a documented, ongoing CI issue
- **Every DevOps team** — Flaky tests waste millions of developer hours

## 👤 Author

Built by [Abhisharydv90](https://github.com/Abhisharydv90). 