# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands
- Build/Run: `go run ./cmd/server`
- Install dependencies: `go mod download`
- API Tests: `./test_api.sh`

## Project Memory
- The project has been migrated from `/internal` directory structure to `/pkg`
- A new package `pkg/api/utils` was created to resolve import cycles
- The database is located at `./pkg/db/sqlite/main.db`

## Code Style Guidelines
- Format: Use standard Go formatting (`gofmt`)
- Imports: Group standard library imports first, third-party imports second, alphabetically sorted
- Types: Define in separate files with `types_` prefix, use CamelCase for struct fields with JSON tags
- Naming: Use CamelCase for functions (e.g., `HandleRegister`), prefix error vars with `Err`
- Error Handling: Check errors and return them up the call stack, log errors before returning
- File Organization: Use `handlers_` prefix for handler files, `storage_` prefix for storage files
- JSON: Use snake_case for JSON field names

## Project Structure
- Go 1.23+ with SQLite database
- REST API with WebSocket support
- Main dependencies: go-sqlite3, gorilla/websocket, golang.org/x/crypto
- Code organized in `/pkg` with domain-driven design