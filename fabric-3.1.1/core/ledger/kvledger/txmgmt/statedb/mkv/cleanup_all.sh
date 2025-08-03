#!/bin/bash

# Script dọn dẹp hoàn chỉnh cho MKV Key Management System
# Xóa tất cả file tạm thời, key files, và build artifacts

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

echo "=========================================="
echo "    MKV Complete Cleanup Script"
echo "=========================================="
echo

# Step 1: Stop any running processes
print_info "1. Stopping any running processes..."
pkill -f "key_manager" 2>/dev/null || true
pkill -f "demo" 2>/dev/null || true

# Step 2: Remove key files
print_info "2. Removing key files..."
rm -f k1.key k0.key encrypted_k1.key decrypted_k1.key

# Step 3: Remove temporary C files and executables
print_info "3. Removing temporary C files and executables..."
rm -f encrypt_k1.c encrypt_k1
rm -f decrypt_k1.c decrypt_k1
rm -f test_*.c test_*

# Step 4: Remove build artifacts
print_info "4. Removing build artifacts..."
rm -f *.o
rm -f libmkv.so

# Step 5: Remove log files
print_info "5. Removing log files..."
rm -f /tmp/state_mkv.log
rm -f /root/state_mkv.log 2>/dev/null || true
rm -f *.log

# Step 6: Remove Go test artifacts
print_info "6. Removing Go test artifacts..."
rm -f mkv.test
rm -rf /tmp/go-build*

# Step 7: Remove any other temporary files
print_info "7. Removing other temporary files..."
rm -f *.tmp
rm -f *.bak
rm -f *.swp
rm -f .*.swp

# Step 8: Clean Go cache (optional)
print_info "8. Cleaning Go cache..."
go clean -cache -testcache -modcache 2>/dev/null || true

# Step 9: Show what's left
print_info "9. Remaining files:"
echo
ls -la
echo

# Step 10: Verify cleanup
print_info "10. Verifying cleanup..."
if [ -f "k1.key" ] || [ -f "k0.key" ] || [ -f "encrypted_k1.key" ] || [ -f "decrypted_k1.key" ]; then
    print_error "Some key files still exist!"
    exit 1
fi

if [ -f "decrypt_k1" ] || [ -f "encrypt_k1" ] || [ -f "decrypt_k1.c" ] || [ -f "encrypt_k1.c" ]; then
    print_error "Some temporary C files still exist!"
    exit 1
fi

if [ -f "libmkv.so" ]; then
    print_error "libmkv.so still exists!"
    exit 1
fi

print_success "Cleanup completed successfully!"
echo
echo "Files removed:"
echo "✓ Key files (k1.key, k0.key, encrypted_k1.key, decrypted_k1.key)"
echo "✓ Temporary C files (encrypt_k1.c, decrypt_k1.c, etc.)"
echo "✓ Executables (encrypt_k1, decrypt_k1)"
echo "✓ Build artifacts (*.o, libmkv.so)"
echo "✓ Log files (state_mkv.log)"
echo "✓ Go test artifacts"
echo
echo "Next steps:"
echo "1. Run 'make' to rebuild the library"
echo "2. Run './key_manager.sh init' to initialize the system"
echo "3. Run './demo.sh' to test everything"
echo
print_success "System is ready for fresh start!" 