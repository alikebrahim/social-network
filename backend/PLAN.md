# Project Restructuring and Logger Implementation Plan

This document outlines a detailed plan for:
1. Restructuring the project to move code from `/internal` to `/pkg`
2. Implementing a logging system for the application

## Implementation Phases

Each phase will be documented in a separate markdown file with a detailed checklist:

1. **PHASE-1.md**: Assessment Phase - evaluating current structure and dependencies
2. **PHASE-2.md**: Pre-Migration Checks - ensuring rollback capability
3. **PHASE-3.md**: Migration Execution - creating new structure and moving files
4. **PHASE-4.md**: Post-Migration Testing - verifying functionality
5. **PHASE-5.md**: Logger Design - defining the logger interface and features
6. **PHASE-6.md**: Logger Implementation - building the logger package
7. **PHASE-7.md**: Logger Integration - adding logging to the application components

## Phase Overview

### Assessment Phase
- Identifying all import dependencies
- Evaluating database path references
- Checking for hardcoded paths
- Creating dependency graph

### Pre-Migration Checks
- Creating branch backup
- Testing current application
- Establishing verification criteria

### Migration Execution
- Creating new directory structure
- Updating database connection code
- Migrating files in dependency order
- Updating import statements
- Updating migration scripts

### Post-Migration Testing
- Testing application functionality
- Verifying database operations
- Documenting rollback procedure

### Logger Design
- Defining logger interface
- Establishing configurable options
- Designing middleware approach

### Logger Implementation
- Creating basic logger structure
- Adding output options
- Implementing request logging
- Creating context utilities

### Logger Integration
- Adding logger to application initialization
- Integrating with API layer
- Adding logging to storage operations
- Implementing WebSocket logging

## New Directory Structure Goal

```
/pkg
  /api            (from internal/api)
  /db             (from pkg/db, merging with internal database code)
    /migrations   (remain as-is but with updated paths)
    /sqlite       (remain as-is but with updated paths)
  /domain         (from internal/domain)
  /errors         (already here)
  /logger         (to be added)
  /storage        (from internal/storage)
  /websocket      (from internal/websocket)
```

## Implementation Schedule

- Phase 1: Assessment (0.5-1 day)
- Phase 2: Pre-Migration Checks (0.5 day)
- Phase 3: Migration Execution (1-2 days)
- Phase 4: Post-Migration Testing (0.5 day)
- Phase 5: Logger Design (0.5 day)
- Phase 6: Logger Implementation (1-2 days)
- Phase 7: Logger Integration (1-2 days)

## Process

For each phase:
1. Create a detailed PHASE-X.md document with specific tasks
2. Execute tasks according to the checklist
3. Document completion and any issues encountered
4. Verify phase completion before proceeding to the next phase