# Cross-Platform Compatibility Guide

## ✅ **Full Cross-Platform Support**

Gophex is designed to work seamlessly across **Windows**, **macOS**, and **Linux** systems.

## 🖥️ **Platform Support Matrix**

| Feature | Windows | macOS | Linux | Notes |
|---------|---------|-------|-------|-------|
| **Core Generation** | ✅ | ✅ | ✅ | Full Go compatibility |
| **Interactive CLI** | ✅ | ✅ | ✅ | Survey library works on all platforms |
| **Database Support** | ✅ | ✅ | ✅ | All database types supported |
| **Migration Scripts** | ✅ | ✅ | ✅ | Platform-specific scripts generated |
| **Change Detection** | ✅ | ✅ | ✅ | Platform-specific implementations |
| **Directory Opening** | ✅ | ✅ | ✅ | Uses platform-specific commands |
| **Health Checks** | ✅ | ✅ | ✅ | HTTP client-based (cross-platform) |
| **Dependency Management** | ✅ | ✅ | ✅ | Go modules work everywhere |

## 🔧 **Platform-Specific Implementations**

### **Script Generation**

Gophex automatically generates the appropriate scripts for each platform:

#### **Unix/Linux/macOS**
- `scripts/migrate.sh` - Bash migration script
- `scripts/detect-changes.sh` - Bash change detection script

#### **Windows**
- `scripts/migrate.bat` - Batch migration script
- `scripts/detect-changes.bat` - Batch change detection script

### **Directory Opening**

```go
switch runtime.GOOS {
case "darwin":
    cmd = exec.Command("open", projectPath)      // macOS
case "linux":
    cmd = exec.Command("xdg-open", projectPath)  // Linux
case "windows":
    cmd = exec.Command("explorer", projectPath)  // Windows
}
```

### **Script Execution**

```go
if runtime.GOOS == "windows" {
    // Windows: cmd /c script.bat args
    return exec.Command("cmd", "/c", scriptPath, args...)
} else {
    // Unix: bash script.sh args
    return exec.Command("bash", scriptPath, args...)
}
```

### **Health Check Testing**

Uses HTTP client instead of curl for cross-platform compatibility:
```go
client := &http.Client{Timeout: 5 * time.Second}
resp, err := client.Get("http://localhost:8080/api/v1/health")
```

## 📋 **Platform-Specific Requirements**

### **Windows**
- **Go**: Go 1.19+ installed and in PATH
- **Git**: Git for Windows (optional, for change detection)
- **Database Tools**: 
  - golang-migrate (auto-installed by Gophex)
  - MongoDB shell (for MongoDB projects)

### **macOS**
- **Go**: Go 1.19+ installed
- **Git**: Xcode Command Line Tools or Git
- **Database Tools**:
  - golang-migrate (auto-installed by Gophex)
  - MongoDB shell: `brew install mongosh`

### **Linux**
- **Go**: Go 1.19+ installed
- **Git**: Git package
- **Database Tools**:
  - golang-migrate (auto-installed by Gophex)
  - MongoDB shell: Package manager or official installer

## 🚀 **Installation Instructions**

### **Windows**

```powershell
# Install Go
# Download from: https://golang.org/dl/

# Install Gophex
go install github.com/buildwithhp/gophex@latest

# Verify installation
gophex --help
```

### **macOS**

```bash
# Install Go (using Homebrew)
brew install go

# Install Gophex
go install github.com/buildwithhp/gophex@latest

# Verify installation
gophex --help
```

### **Linux**

```bash
# Install Go (Ubuntu/Debian)
sudo apt update
sudo apt install golang-go

# Install Gophex
go install github.com/buildwithhp/gophex@latest

# Verify installation
gophex --help
```

## 🔍 **Platform-Specific Features**

### **Windows Batch Scripts**

Generated Windows batch files include:
- Color-coded output (text-based)
- Error handling with proper exit codes
- Environment variable support
- Cross-platform command equivalents

### **Unix Shell Scripts**

Generated Unix shell scripts include:
- ANSI color support
- Advanced error handling
- Comprehensive file operations
- Git integration

### **MongoDB Shell Detection**

Automatically detects available MongoDB shell:
- `mongosh` (MongoDB 5.0+)
- `mongo` (Legacy shell)
- Provides platform-specific installation instructions

## 🧪 **Testing Cross-Platform Compatibility**

### **Build for All Platforms**

```bash
# Build for Windows
GOOS=windows GOARCH=amd64 go build -o gophex.exe .

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o gophex-darwin .

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o gophex-linux .
```

### **Test Generated Projects**

1. **Generate API project**
2. **Check generated scripts**:
   - Windows: `scripts\migrate.bat`, `scripts\detect-changes.bat`
   - Unix: `scripts/migrate.sh`, `scripts/detect-changes.sh`
3. **Test script execution**
4. **Verify database connectivity**

## 🛠️ **Troubleshooting**

### **Windows Issues**

**Script Execution Errors:**
```cmd
# If batch files don't execute
# Check file associations
assoc .bat

# Run with explicit command
cmd /c scripts\migrate.bat status
```

**PATH Issues:**
```cmd
# Check Go installation
go version

# Check GOPATH/GOBIN
echo %GOPATH%
echo %GOBIN%
```

### **macOS/Linux Issues**

**Permission Errors:**
```bash
# Make scripts executable
chmod +x scripts/*.sh

# Check script permissions
ls -la scripts/
```

**Shell Issues:**
```bash
# Use explicit bash
bash scripts/migrate.sh status

# Check bash availability
which bash
```

## 📊 **Performance Considerations**

### **Windows**
- Batch file execution is slower than shell scripts
- File operations may have different performance characteristics
- Antivirus software may affect build times

### **macOS/Linux**
- Native shell script execution
- Better file system performance
- Faster build and test cycles

## 🔄 **Future Enhancements**

### **Planned Improvements**
- PowerShell scripts for Windows (alternative to batch)
- Enhanced Windows terminal color support
- Cross-platform GUI tools integration
- Docker-based development environment

### **Community Contributions**
- Platform-specific optimizations welcome
- Testing on different OS versions
- Package manager integrations

## 📞 **Platform-Specific Support**

### **Windows Support**
- Windows 10/11 fully supported
- Windows Server 2019+ supported
- PowerShell and Command Prompt compatible

### **macOS Support**
- macOS 10.15+ supported
- Intel and Apple Silicon compatible
- Homebrew integration available

### **Linux Support**
- Ubuntu 18.04+ supported
- CentOS/RHEL 7+ supported
- Debian 10+ supported
- Arch Linux supported

---

**Gophex works everywhere Go works!** 🌍