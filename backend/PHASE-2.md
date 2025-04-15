# PHASE 2: Pre-Migration Checks

## Overview
Before proceeding with the migration, we need to ensure we have the necessary safeguards in place to allow for rollback if needed and to verify the application still works after migration.

## Checklist

### 1. Git Branch Backup
- [x] Ensure all current changes are committed
- [x] Create a new branch for the restructuring work (`restructure-internal-to-pkg`)
- [x] Confirm branch creation was successful

### 2. Application Functionality Verification
- [x] Attempt to run the application to verify current functionality
  - ISSUE: Import cycle detected between `internal/api` and `internal/api/handlers`
  - The application cannot compile or run due to this circular dependency
  - This is a key issue that our restructuring plan will address
- [ ] Test basic API endpoints (blocked by import cycle)
- [ ] Check database connections (will verify during restructuring)

### 3. Verification Criteria Definition
- [x] Define specific tests to run after migration
  - The application should compile without import cycle errors
  - The application should start successfully
  - The test_api.sh script should run without errors
  - API endpoints should return expected responses
- [x] Document expected behavior for key features
  - User registration should create a new user in the database
  - User login should generate a valid session token
  - Chat functionality should work with WebSockets
- [x] Create a verification script if possible
  - Will use existing test_api.sh script to verify API functionality

### 4. Rollback Plan
- [x] Document steps to revert to the original structure if needed
  - Use git reset --hard to discard changes
  - Switch back to the backend-ak-2 branch
  - Restore database backup if needed
- [x] Ensure understanding of git reset/revert commands
  - `git reset --hard` to discard all changes
  - `git checkout backend-ak-2` to return to the original branch
- [x] Identify any non-git-tracked files that might need manual backup
  - Only the SQLite database file (already backed up)

### 5. Database Backup
- [x] Create a backup of the database file (`main.db.bak`)
- [x] Verify the backup is valid
- [x] Document restoration procedure

## Commands to Execute

### Git Branch Operations
```bash
# Check current git status
git status

# Create a new branch
git checkout -b restructure-internal-to-pkg

# Verify branch creation
git branch
```

### Database Backup
```bash
# Create a backup of the database file
cp /home/alikebrahim/dev/social-network/backend/pkg/db/sqlite/main.db /home/alikebrahim/dev/social-network/backend/pkg/db/sqlite/main.db.bak

# Verify the backup
ls -la /home/alikebrahim/dev/social-network/backend/pkg/db/sqlite/
```

### Application Verification
```bash
# Run the application
cd /home/alikebrahim/dev/social-network/backend
go run ./cmd/server

# In a separate terminal, test the API
curl -X POST http://localhost:3000/auth/register -d '{"email":"test@example.com", "password":"password123", "first_name":"Test", "last_name":"User", "date_of_birth":"1990-01-01"}'
```

## Rollback Procedure
If the migration fails, execute the following steps:

1. Discard all changes and return to the main branch:
```bash
git reset --hard
git checkout main
```

2. Restore database backup if needed:
```bash
cp /home/alikebrahim/dev/social-network/backend/pkg/db/sqlite/main.db.bak /home/alikebrahim/dev/social-network/backend/pkg/db/sqlite/main.db
```

## Success Criteria
Pre-migration checks are successful when:
1. A new git branch has been created ✅
2. The application's current state is documented:
   - Import cycle preventing compilation ✅
   - This will be fixed during restructuring
3. Verification criteria for post-migration are defined ✅
4. A database backup exists ✅
5. The rollback procedure is documented and understood ✅

**Note:** Since the application currently has an import cycle that prevents it from compiling or running, we can't verify its functionality. Fixing this is a key goal of our restructuring effort.