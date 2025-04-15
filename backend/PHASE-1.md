# PHASE 1: Project Structure Assessment

## Overview
This phase evaluates the current project structure, identifies dependencies, and prepares for restructuring from `/internal` to `/pkg`.

## Checklist

### 1. Directory Structure Analysis
- [ ] List all files and directories in `/internal`
- [ ] List all files and directories in `/pkg`
- [ ] Identify potential conflicts between the two structures
- [ ] Document the ownership and purpose of each directory

### 2. Import Dependencies Analysis
- [ ] Scan all `.go` files for import statements with `socialNetwork/internal/*`
- [ ] Create a dependency graph showing which packages depend on others
- [ ] Identify circular dependencies that need to be resolved
- [ ] Document all external package dependencies

### 3. Database Path References
- [ ] Identify all files that reference the SQLite database path
- [ ] Document the current path and how it's referenced (relative/absolute)
- [ ] Check for any environment variables or configuration that might affect database paths
- [ ] Analyze how database migrations reference the database file

### 4. Hardcoded Path Check
- [ ] Search for hardcoded paths in all `.go` files
- [ ] Identify any file operations that use hardcoded paths
- [ ] Check for path references in tests or scripts
- [ ] Document all places where paths would need to be updated

### 5. Namespace Analysis
- [ ] Review package naming conventions in both `/internal` and `/pkg`
- [ ] Ensure no package name conflicts will occur after migration
- [ ] Plan for consistent package naming in the new structure
- [ ] Document any package names that should be changed for clarity

### 6. Build Process Assessment
- [ ] Review how the application is built and deployed
- [ ] Check CI/CD scripts for path dependencies
- [ ] Verify Go module configuration impacts
- [ ] Identify any build scripts that might be affected

### 7. Test Coverage Evaluation
- [ ] Identify existing tests and their dependencies
- [ ] Verify how tests reference the code under test
- [ ] Check for test fixtures or data that might need updating
- [ ] Document test coverage and any gaps to address

### 8. Migration Risk Assessment
- [ ] Identify high-risk components that could break during migration
- [ ] Document components with many internal dependencies
- [ ] List potential failure points and mitigation strategies
- [ ] Define criteria for a successful migration

## Completion Criteria
- Complete documentation of all current project structure
- Full dependency graph with clear migration order
- Comprehensive list of all path references that need updating
- Identified risks and mitigation strategies
- Clear understanding of impact on build and test processes

## Tools and Commands

### Directory Structure Analysis
```bash
find /home/alikebrahim/dev/social-network/backend/internal -type f -name "*.go" | sort
find /home/alikebrahim/dev/social-network/backend/pkg -type f -name "*.go" | sort
```

### Import Dependency Analysis
```bash
grep -r "\"socialNetwork/internal/" --include="*.go" /home/alikebrahim/dev/social-network/backend
```

### Database Path References
```bash
grep -r "db.*Open" --include="*.go" /home/alikebrahim/dev/social-network/backend
grep -r "sqlite" --include="*.go" /home/alikebrahim/dev/social-network/backend
```

### Hardcoded Path Check
```bash
grep -r "\.\./" --include="*.go" /home/alikebrahim/dev/social-network/backend
grep -r "filepath\." --include="*.go" /home/alikebrahim/dev/social-network/backend
```

## Output Documents
- Dependency graph diagram
- List of all import statements to be updated
- Path reference inventory
- Risk assessment matrix