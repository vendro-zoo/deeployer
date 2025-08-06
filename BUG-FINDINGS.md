# Bug Findings Report

After systematically scanning the codebase, the following potential bugs and issues were identified:

## 1. **Path Validation Issues** (Medium Severity) ✅ FIXED

**File:** `internal/config/config.go` lines 81-83
- **Issue:** Project path validation only checks if the string is empty, not if the path actually exists or is accessible
- **Impact:** Deployment could fail later when trying to change to a non-existent directory
- **Location:** `Project.Validate()` method

## 2. **Missing Project Path Display** (Low Severity) ✅ FIXED

**File:** `cmd/list.go` line 49
- **Issue:** The list command shows output directory but not the project path, making it hard to see where projects are located
- **Impact:** Users can't see full project configuration at a glance

## 3. **Rsync Path Concatenation Logic** (Medium Severity) ✅ FIXED

**File:** `internal/rsync/rsync.go` lines 58, 80-84
- **Issue:** The `ensureTrailingSlash` function always adds a trailing slash, but rsync behavior differs between files and directories
- **Impact:** Could cause unintended directory structure on remote when syncing files vs directories

## 4. **SSH Agent Connection Leak** (Medium Severity) ✅ FIXED

**File:** `internal/ssh/ssh.go` lines 126-131
- **Issue:** SSH agent connection (`net.Dial`) is never closed, causing potential resource leak
- **Impact:** Long-running deployments could exhaust file descriptors

## 5. **No SSH Timeout Configuration** (Medium Severity) ✅ FIXED

**File:** `internal/ssh/ssh.go` line 84-88
- **Issue:** SSH client config has no timeout settings, could hang indefinitely on unresponsive hosts
- **Impact:** Deployments could hang forever on network issues

## 6. **Inconsistent Error Handling** (Low Severity) ✅ FIXED

**File:** `cmd/validate.go` lines 24-26
- **Issue:** Validation failure prints error and then returns it, causing double error output
- **Impact:** Confusing error messages for users

## 7. **Missing Input Validation** (Medium Severity) ⚠️ USER RESPONSIBILITY

**File:** `internal/config/config.go` lines 100-117
- **Issue:** Remote paths, hosts, and users are not validated for format or dangerous characters
- **Impact:** Could allow injection attacks or malformed configurations
- **Note:** Left unfixed as user is responsible for configuration content

## 8. **Directory Traversal Risk** (High Severity) ✅ FIXED

**File:** `cmd/deploy.go` line 87
- **Issue:** `filepath.Join(project.Path, project.OutputDir)` could be vulnerable to directory traversal if OutputDir contains "../"
- **Impact:** Could allow accessing files outside project directory

## 9. **Race Condition in Dry-Run** (Low Severity) ⚠️ USER RESPONSIBILITY

**File:** `internal/rsync/rsync.go` lines 50-51
- **Issue:** Dry-run flag is added to args after user-provided options, could be overridden by user options
- **Impact:** Dry-run might not work as expected if user provides conflicting flags
- **Note:** Left unfixed as user is responsible for rsync option configuration

## Summary

- **Total Issues Found:** 9
- **Fixed:** 7
- **User Responsibility:** 2
- **Critical Issues Fixed:** Directory traversal vulnerability
- **Security Improvements:** SSH timeout, connection leak fixes, path validation