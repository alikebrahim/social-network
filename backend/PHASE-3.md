# PHASE 3: Migration Execution

## Overview
This phase involves creating the new directory structure, moving files, and updating import paths. The migration will be performed in a specific order to minimize issues.

A key goal of this phase is to resolve the import cycle between `internal/api` and `internal/api/handlers` that was discovered during Phase 2. This circular dependency is preventing the application from compiling and running.

## Checklist

### 1. Create New Directory Structure
- [x] Create necessary directories in `/pkg` if they don't exist
- [x] Verify directory permissions match current structure

### 2. Migrate Domain Models (First)
- [x] Move each domain model package to `/pkg/domain/`
- [x] Update import paths in all domain model files

### 3. Migrate Storage Interface
- [x] Move storage interface (`store.go`) to `/pkg/storage/`
- [x] Update import paths in storage interface
- [x] Verify interface references to domain models

### 4. Migrate SQLite Implementation
- [x] Move SQLite implementation files to `/pkg/storage/sqlite/`
- [x] Update database path reference in `sqlite.go`
- [x] Update import paths in all SQLite implementation files

### 5. Migrate WebSocket Package
- [x] Move WebSocket files to `/pkg/websocket/`
- [x] Update import paths in WebSocket files
- [x] Verify WebSocket references to other packages

### 6. Migrate API Components
- [x] Move API files to `/pkg/api/`
- [x] Update import paths in all API files
- [x] Verify API references to other packages
- [x] Break the import cycle:
  - [x] Create a separate `pkg/api/utils` or `pkg/api/common` package for shared functions
  - [x] Move `WriteJson` and `makeHTTPHandleFunc` to this package
  - [x] Update imports to eliminate circular dependencies

### 7. Update Main Application
- [x] Update import paths in `cmd/server/main.go`
- [x] Verify main application references to migrated packages

### 8. Verify Migration
- [x] Check for any remaining references to `/internal/`
- [x] Ensure all import paths have been updated
- [x] Verify no duplicate files exist in both locations

## Commands to Execute

### 1. Create Directory Structure
```bash
# Create necessary directories
mkdir -p /home/alikebrahim/dev/social-network/backend/pkg/api/handlers
mkdir -p /home/alikebrahim/dev/social-network/backend/pkg/storage/sqlite
mkdir -p /home/alikebrahim/dev/social-network/backend/pkg/domain/auth
mkdir -p /home/alikebrahim/dev/social-network/backend/pkg/domain/chat
mkdir -p /home/alikebrahim/dev/social-network/backend/pkg/domain/following
mkdir -p /home/alikebrahim/dev/social-network/backend/pkg/domain/groups
mkdir -p /home/alikebrahim/dev/social-network/backend/pkg/domain/posts
mkdir -p /home/alikebrahim/dev/social-network/backend/pkg/domain/profile
mkdir -p /home/alikebrahim/dev/social-network/backend/pkg/websocket
```

### 2. Migration Commands for Each Package

#### Domain Models
```bash
# Move domain models
cp /home/alikebrahim/dev/social-network/backend/internal/domain/auth/auth.go /home/alikebrahim/dev/social-network/backend/pkg/domain/auth/
cp /home/alikebrahim/dev/social-network/backend/internal/domain/chat/chat.go /home/alikebrahim/dev/social-network/backend/pkg/domain/chat/
cp /home/alikebrahim/dev/social-network/backend/internal/domain/following/following.go /home/alikebrahim/dev/social-network/backend/pkg/domain/following/
cp /home/alikebrahim/dev/social-network/backend/internal/domain/groups/groups.go /home/alikebrahim/dev/social-network/backend/pkg/domain/groups/
cp /home/alikebrahim/dev/social-network/backend/internal/domain/posts/posts.go /home/alikebrahim/dev/social-network/backend/pkg/domain/posts/
cp /home/alikebrahim/dev/social-network/backend/internal/domain/profile/profile.go /home/alikebrahim/dev/social-network/backend/pkg/domain/profile/

# Update import paths in domain files as needed
```

#### Storage Interface
```bash
# Move storage interface
cp /home/alikebrahim/dev/social-network/backend/internal/storage/store.go /home/alikebrahim/dev/social-network/backend/pkg/storage/

# Update import paths in storage interface
sed -i 's|"socialNetwork/internal/domain/|"socialNetwork/pkg/domain/|g' /home/alikebrahim/dev/social-network/backend/pkg/storage/store.go
```

#### SQLite Implementation
```bash
# Move SQLite implementation
cp /home/alikebrahim/dev/social-network/backend/internal/storage/sqlite/*.go /home/alikebrahim/dev/social-network/backend/pkg/storage/sqlite/

# Update database path in sqlite.go
sed -i 's|"../../pkg/db/sqlite/main.db"|"../../../pkg/db/sqlite/main.db"|g' /home/alikebrahim/dev/social-network/backend/pkg/storage/sqlite/sqlite.go

# Update import paths in SQLite files
sed -i 's|"socialNetwork/internal/domain/|"socialNetwork/pkg/domain/|g' /home/alikebrahim/dev/social-network/backend/pkg/storage/sqlite/*.go
sed -i 's|"socialNetwork/internal/storage"|"socialNetwork/pkg/storage"|g' /home/alikebrahim/dev/social-network/backend/pkg/storage/sqlite/*.go
```

#### WebSocket Package
```bash
# Move WebSocket files
cp /home/alikebrahim/dev/social-network/backend/internal/websocket/websocket.go /home/alikebrahim/dev/social-network/backend/pkg/websocket/

# Update import paths in WebSocket files
sed -i 's|"socialNetwork/internal/domain/|"socialNetwork/pkg/domain/|g' /home/alikebrahim/dev/social-network/backend/pkg/websocket/*.go
```

#### API Components
```bash
# Create utils directory for breaking the import cycle
mkdir -p /home/alikebrahim/dev/social-network/backend/pkg/api/utils

# Create utils file with shared functions
cat > /home/alikebrahim/dev/social-network/backend/pkg/api/utils/http.go << 'EOF'
package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// ApiError represents an API error response
type ApiError struct {
	Error string
}

// WriteJson sends a JSON response
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("json.NewEncoder error:", err)
		return err
	}
	return nil
}

// MakeHTTPHandleFunc is a helper that allows handlers to return errors
func MakeHTTPHandleFunc(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
EOF

# Move API files (excluding the functions we moved to utils)
cp /home/alikebrahim/dev/social-network/backend/internal/api/*.go /home/alikebrahim/dev/social-network/backend/pkg/api/
cp /home/alikebrahim/dev/social-network/backend/internal/api/handlers/*.go /home/alikebrahim/dev/social-network/backend/pkg/api/handlers/

# Update server.go to remove functions now in utils
sed -i '/type apiFunc/,/^}$/d' /home/alikebrahim/dev/social-network/backend/pkg/api/server.go
sed -i '/type ApiError/,/^}$/d' /home/alikebrahim/dev/social-network/backend/pkg/api/server.go
sed -i '/func WriteJson/,/^}$/d' /home/alikebrahim/dev/social-network/backend/pkg/api/server.go
sed -i '/func makeHTTPHandleFunc/,/^}$/d' /home/alikebrahim/dev/social-network/backend/pkg/api/server.go

# Update router.go to use utils package
sed -i 's/makeHTTPHandleFunc/utils.MakeHTTPHandleFunc/g' /home/alikebrahim/dev/social-network/backend/pkg/api/router.go

# Update import paths in API files
sed -i 's|"socialNetwork/internal/api/handlers"|"socialNetwork/pkg/api/handlers"|g' /home/alikebrahim/dev/social-network/backend/pkg/api/*.go
sed -i 's|"socialNetwork/internal/storage"|"socialNetwork/pkg/storage"|g' /home/alikebrahim/dev/social-network/backend/pkg/api/*.go
sed -i 's|"socialNetwork/internal/websocket"|"socialNetwork/pkg/websocket"|g' /home/alikebrahim/dev/social-network/backend/pkg/api/*.go
sed -i '1,/^)/ s|^)/|)\n\t"socialNetwork/pkg/api/utils"|' /home/alikebrahim/dev/social-network/backend/pkg/api/router.go

# Update import paths in handlers
sed -i 's|"socialNetwork/internal/api"|"socialNetwork/pkg/api/utils"|g' /home/alikebrahim/dev/social-network/backend/pkg/api/handlers/*.go
sed -i 's|api.WriteJson|utils.WriteJson|g' /home/alikebrahim/dev/social-network/backend/pkg/api/handlers/*.go
sed -i 's|"socialNetwork/internal/domain/|"socialNetwork/pkg/domain/|g' /home/alikebrahim/dev/social-network/backend/pkg/api/handlers/*.go
sed -i 's|"socialNetwork/internal/storage"|"socialNetwork/pkg/storage"|g' /home/alikebrahim/dev/social-network/backend/pkg/api/handlers/*.go
sed -i 's|"socialNetwork/internal/websocket"|"socialNetwork/pkg/websocket"|g' /home/alikebrahim/dev/social-network/backend/pkg/api/handlers/*.go
```

#### Main Application
```bash
# Update import paths in main.go
sed -i 's|"socialNetwork/internal/api"|"socialNetwork/pkg/api"|g' /home/alikebrahim/dev/social-network/backend/cmd/server/main.go
sed -i 's|"socialNetwork/internal/storage/sqlite"|"socialNetwork/pkg/storage/sqlite"|g' /home/alikebrahim/dev/social-network/backend/cmd/server/main.go
```

### 3. Verification Commands
```bash
# Check for any remaining internal references
grep -r "socialNetwork/internal" --include="*.go" /home/alikebrahim/dev/social-network/backend

# Ensure all go files have been moved with updated imports
find /home/alikebrahim/dev/social-network/backend/pkg -name "*.go" | wc -l
find /home/alikebrahim/dev/social-network/backend/internal -name "*.go" | wc -l
```

## Success Criteria
The migration is successful when:
1. All files have been moved to their new locations ✓
2. All import paths have been updated ✓
3. The database path reference in `sqlite.go` has been updated ✓
4. No references to `/internal/` remain in the codebase ✓
5. The import cycle between API packages has been resolved ✓
6. The application builds successfully without import cycle errors ✓
7. Tests run successfully with the migrated code ✓

## Migration Outcomes
- ✓ Successfully migrated all code from internal/ to pkg/
- ✓ Fixed import cycle by creating pkg/api/utils package for shared HTTP functions
- ✓ Implemented missing profile.go file for SQLite storage implementation
- ✓ Fixed database path reference for SQLite connection
- ✓ Application compiles and runs successfully
- ✓ All tests pass successfully

## Rollback Procedure
If any issues are encountered during migration:
1. Discard all changes and return to the main branch
2. If needed, restore from the backup created in Phase 2