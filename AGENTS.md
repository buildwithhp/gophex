# AGENTS.md - Development Guidelines

## 🚨 CRITICAL: Change Detection & Confirmation Protocol

### **BEFORE MAKING ANY CHANGES**
1. **Always check if files have been manually modified** since generation
2. **Never overwrite manual changes** without explicit user confirmation
3. **Break large changes into small, confirmable steps**
4. **Explain the purpose and impact** of each proposed change

### **Change Detection Workflow**
```bash
# Check if this is a generated project
if [ -f ".gophex-generated" ] || [ -f "go.mod" ]; then
    echo "⚠️  Generated project detected. Checking for manual changes..."
    
    # Look for signs of manual modifications:
    # - Git history with non-generation commits
    # - Files with timestamps newer than .gophex-generated
    # - Custom code patterns not in templates
fi
```

### **Manual Change Detection Indicators**
- Custom imports not in original templates
- Modified function signatures
- Additional methods or structs
- Custom business logic
- Modified database schemas
- Custom middleware or handlers
- Environment variables not in original config

### **Required Confirmation Process**

#### **Step 1: Detect Changes**
```
🔍 CHANGE DETECTION RESULTS:
- Modified files: [list files]
- Custom code detected in: [specific locations]
- Last generation: [timestamp]
- Manual changes since: [timestamp]
```

#### **Step 2: Explain Purpose**
```
📋 PROPOSED CHANGES:
Purpose: [Why this change is needed]
Impact: [What will be affected]
Risk Level: [Low/Medium/High]
Reversible: [Yes/No]
```

#### **Step 3: Get Explicit Confirmation**
```
❓ CONFIRMATION REQUIRED:
Do you want to proceed with these changes?
- [Y] Yes, proceed with all changes
- [S] Step-by-step confirmation
- [R] Review changes first
- [N] No, cancel operation
```

## 🔄 Step-by-Step Execution Protocol

### **For Large Changes (>5 files or >100 lines)**

#### **Phase 1: Planning**
1. **Break down into logical steps** (max 3-5 files per step)
2. **Identify dependencies** between changes
3. **Create execution plan** with rollback points
4. **Present plan to user** for approval

#### **Phase 2: Step-by-Step Execution**
```
📝 EXECUTION PLAN:
Step 1/5: Update database configuration
  - Files: internal/database/config.go, internal/database/factory.go
  - Purpose: Add new database type support
  - Risk: Low (backward compatible)
  
Continue with Step 1? [Y/n/s(skip)/q(quit)]
```

#### **Phase 3: Validation After Each Step**
```
✅ STEP 1 COMPLETED:
- Files modified: 2
- Build status: ✅ Success
- Tests status: ✅ All passing
- Lint status: ⚠️  2 warnings (non-breaking)

Continue to Step 2? [Y/n/r(rollback)/q(quit)]
```

### **Rollback Strategy**
- **Git-based rollback**: Create commits for each step
- **Backup strategy**: Keep copies of modified files
- **Validation points**: Test after each step
- **Safe exit**: Allow cancellation at any point

## 🛡️ Protection Mechanisms

### **File Protection Rules**
1. **Never modify** files with custom business logic without confirmation
2. **Always preserve** user-added imports and dependencies
3. **Maintain** existing function signatures unless explicitly requested
4. **Backup** original files before making changes
5. **Validate** changes don't break existing functionality

### **Code Pattern Recognition**
```go
// PROTECTED: Custom business logic
func (s *UserService) CustomBusinessMethod() error {
    // User-added code - DO NOT MODIFY
}

// SAFE TO MODIFY: Generated boilerplate
func (s *UserService) GetUser(id int) (*User, error) {
    // Template-generated code - can be updated
}
```

### **Database Schema Protection**
- **Never drop** existing tables/collections without confirmation
- **Always use** additive migrations when possible
- **Preserve** custom indexes and constraints
- **Backup** schema before major changes

## 📋 Change Categories & Protocols

### **Category 1: Safe Changes (Auto-approve)**
- Code formatting and linting fixes
- Adding new optional configuration
- Adding new endpoints (non-breaking)
- Documentation updates
- Adding new dependencies (non-conflicting)

### **Category 2: Medium Risk (Require Confirmation)**
- Modifying existing function signatures
- Changing database connection logic
- Updating middleware behavior
- Modifying error handling patterns
- Changing configuration structure

### **Category 3: High Risk (Step-by-step Required)**
- Database schema changes
- Breaking API changes
- Major refactoring
- Dependency version updates
- Security-related modifications

## 🔧 Implementation Guidelines

### **Before Any Code Change**
```bash
# 1. Run the change detection script (for generated projects)
if [ -f "scripts/detect-changes.sh" ]; then
    ./scripts/detect-changes.sh
else
    # Manual detection for non-generated projects
    git log --oneline --since="1 week ago" | head -10
    grep -r "// CUSTOM:\|// TODO:\|// FIXME:" . --exclude-dir=.git || true
fi

# 2. Check build status before changes
go build ./... && echo "✅ Build OK" || echo "❌ Build Failed"

# 3. Check test status
go test ./... && echo "✅ Tests OK" || echo "❌ Tests Failed"
```

### **Change Execution Template**
```
🔄 EXECUTING CHANGE:
File: [filename]
Purpose: [specific purpose]
Changes:
  - Line X: [what's being changed]
  - Line Y: [what's being added]
  
Proceed? [Y/n/d(diff)/s(skip)]
```

### **Validation After Changes**
```bash
# Always run after changes
go build ./...
go test ./...
go fmt ./...
go vet ./...

# Check for breaking changes
git diff --name-only HEAD~1 | xargs -I {} echo "Modified: {}"
```

## 🎯 Best Practices

### **Communication**
- **Always explain WHY** before explaining WHAT
- **Show impact** of changes on existing code
- **Provide alternatives** when possible
- **Respect user decisions** including "no"

### **Safety First**
- **Test incrementally** after each change
- **Maintain working state** at all times
- **Document changes** for future reference
- **Provide rollback instructions**

### **User Experience**
- **Clear progress indicators** for multi-step operations
- **Meaningful error messages** with suggested fixes
- **Consistent confirmation patterns**
- **Respect user workflow** and preferences

## 🚀 Gophex-Specific Workflows

### **Working with Generated Projects**
```bash
# Check if project was generated by Gophex
[ -f ".gophex-generated" ] && echo "Gophex project detected" || echo "Not a Gophex project"

# Run change detection (in generated projects)
./scripts/detect-changes.sh

# Database operations
./scripts/migrate.sh status    # Check migration status
./scripts/migrate.sh up        # Apply migrations (SQL databases)
./scripts/migrate.sh init      # Initialize collections (MongoDB)

# Development workflow
go run cmd/api/main.go         # Start API server
```

### **Template Development (Gophex Core)**
```bash
# Test template generation
go run main.go generate

# Build Gophex CLI
go build -o gophex .

# Test with different database configurations
# (Use interactive prompts to test various combinations)
```

### **Change Management Protocol**
```bash
# 1. Before making changes
./scripts/detect-changes.sh

# 2. Create backup branch (if Git repo)
git checkout -b backup-$(date +%Y%m%d-%H%M%S)
git checkout main

# 3. Make changes in small steps
# 4. Test after each step
go build ./... && go test ./...

# 5. Commit each step
git add . && git commit -m "Step X: [description]"
```

## Build/Test/Lint Commands
- `go build ./...` - Build all packages
- `go test ./...` - Run all tests
- `go test -v ./path/to/package` - Run single package tests
- `go test -run TestFunctionName` - Run specific test
- `go fmt ./...` - Format code
- `go vet ./...` - Static analysis
- `golangci-lint run` - Comprehensive linting (if available)

## Project Structure
This is a Go project following standard Go project layout patterns. See `projectstructures.md` for detailed architecture examples including API, web app, and microservice structures.

## Code Style Guidelines
- Follow standard Go conventions (gofmt, go vet)
- Use meaningful package names (lowercase, no underscores)
- Interfaces should be small and focused
- Error handling: always check errors, wrap with context using `fmt.Errorf`
- Naming: use camelCase for unexported, PascalCase for exported
- Imports: group standard library, third-party, then local packages
- Use `context.Context` for cancellation and timeouts
- Prefer composition over inheritance
- Write tests alongside code (package_test.go)
- Use dependency injection for testability