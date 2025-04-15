# PHASE 4: Post-Migration Testing

## Overview
After completing the migration, thorough testing is required to ensure the application functions correctly with the new structure. This phase focuses on verifying that all components work as expected.

## Checklist

### 1. Build Verification
- [ ] Compile the project to check for build errors
- [ ] Address any compilation issues
- [ ] Verify successful build

### 2. Database Connection Testing
- [ ] Start the application
- [ ] Verify the database connection is established
- [ ] Check for any database-related errors in logs

### 3. API Endpoint Testing
- [ ] Test user registration endpoint
- [ ] Test user login endpoint
- [ ] Test profile endpoints
- [ ] Test social connections endpoints
- [ ] Test other critical endpoints

### 4. WebSocket Testing
- [ ] Test WebSocket connection
- [ ] Verify real-time messaging works

### 5. Error Logging Verification
- [ ] Trigger error scenarios
- [ ] Verify errors are properly logged
- [ ] Check error responses in API

### 6. Performance Check
- [ ] Compare application startup time
- [ ] Check response times for key endpoints
- [ ] Verify memory usage is consistent

## Commands to Execute

### 1. Build Verification
```bash
# Clean any build artifacts
cd /home/alikebrahim/dev/social-network/backend
go clean

# Build the project
go build ./cmd/server

# Check for errors
echo $?
```

### 2. Run the Application
```bash
# Start the server
cd /home/alikebrahim/dev/social-network/backend
go run ./cmd/server
```

### 3. API Testing
```bash
# In a separate terminal, test API endpoints
# Replace localhost with appropriate host if needed

# Test registration
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com", "password":"password123", "first_name":"Test", "last_name":"User", "date_of_birth":"1990-01-01"}'

# Test login
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com", "password":"password123"}'

# Save the token from login response
TOKEN="<token from login response>"

# Test get profile
curl -X GET http://localhost:3000/profiles/1 \
  -H "Cookie: session_token=$TOKEN"
```

### 4. WebSocket Testing
```bash
# A simple WebSocket client test can be done with wscat
# Install wscat if not available
npm install -g wscat

# Connect to WebSocket endpoint (adjust userId as needed)
wscat -c "ws://localhost:3000/ws/chat/1" -H "Cookie: session_token=$TOKEN"

# Send a test message
{"content":"Test message"}
```

## Documentation Updates

### 1. Update README.md
- Update any path references in the README
- Document the new project structure
- Update any build or run instructions that might have changed

### 2. Update Import Examples
- Update any code examples that show imports
- Ensure documentation reflects the new organization

## Success Criteria
Post-migration testing is successful when:
1. The application builds without errors
2. The application starts and connects to the database
3. All API endpoints function correctly
4. WebSocket connections work as expected
5. Error handling functions properly
6. Application performance is comparable to pre-migration

## Troubleshooting Common Issues

### 1. Import Path Issues
If you encounter import errors:
- Check for any remaining references to `/internal/`
- Verify all files were copied to the correct locations
- Ensure all import statements were updated

### 2. Database Connection Issues
If database connection fails:
- Verify the path in `sqlite.go` is correct
- Check file permissions on the database file
- Ensure the database file exists at the expected location

### 3. Build Errors
For build errors:
- Check the Go module configuration
- Verify there are no syntax errors from search/replace operations
- Ensure all dependencies are available