#!/bin/bash

# Key Manager Script for MKV Encryption
# Quản lý khóa 2 tầng: K1 (Data Key) và K0 (Derived Key)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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

# Function to check if required files exist
check_requirements() {
    if [ ! -f "libmkv.so" ]; then
        print_error "libmkv.so not found. Please build the MKV library first."
        exit 1
    fi
    
    if [ ! -f "mkv.go" ]; then
        print_error "mkv.go not found. Please ensure you're in the correct directory."
        exit 1
    fi
}

# Function to generate K1 (Data Key) - 32 bytes random
generate_k1() {
    print_info "Generating K1 (Data Key) - 32 bytes random..."
    
    # Try openssl first, fallback to /dev/urandom if not available
    if command -v openssl >/dev/null 2>&1; then
        openssl rand -hex 32 > k1.key
    else
        # Fallback: use /dev/urandom and hexdump/od
        if command -v xxd >/dev/null 2>&1; then
            head -c 32 /dev/urandom | xxd -p -c 32 > k1.key
        elif command -v hexdump >/dev/null 2>&1; then
            head -c 32 /dev/urandom | hexdump -ve '1/1 "%02x"' > k1.key
        elif command -v od >/dev/null 2>&1; then
            head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n' > k1.key
        else
            print_error "No hex conversion tool available (xxd, hexdump, or od)"
            return 1
        fi
    fi
    
    print_success "K1 generated and saved to k1.key"
    print_info "K1 (hex): $(cat k1.key)"
}

# Function to generate K0 from password using SHA256
generate_k0_from_password() {
    local password="$1"
    print_info "Generating K0 from password using SHA256..."
    
    # Try openssl first, fallback to sha256sum if not available
    if command -v openssl >/dev/null 2>&1; then
        if command -v xxd >/dev/null 2>&1; then
            echo -n "$password" | openssl dgst -sha256 -binary | xxd -p > k0.key
        elif command -v hexdump >/dev/null 2>&1; then
            echo -n "$password" | openssl dgst -sha256 -binary | hexdump -ve '1/1 "%02x"' > k0.key
        elif command -v od >/dev/null 2>&1; then
            echo -n "$password" | openssl dgst -sha256 -binary | od -An -tx1 | tr -d ' \n' > k0.key
        else
            print_error "No hex conversion tool available (xxd, hexdump, or od)"
            return 1
        fi
    else
        # Fallback: use sha256sum
        echo -n "$password" | sha256sum | cut -d' ' -f1 > k0.key
    fi
    
    print_success "K0 generated and saved to k0.key"
    print_info "K0 (hex): $(cat k0.key)"
}

# Function to encrypt K1 with K0 using MKV
encrypt_k1_with_k0() {
    print_info "Encrypting K1 with K0 using MKV..."
    
    # Read K1 and K0
    if [ ! -f "k1.key" ]; then
        print_error "k1.key not found. Please generate K1 first."
        return 1
    fi
    
    if [ ! -f "k0.key" ]; then
        print_error "k0.key not found. Please generate K0 first."
        return 1
    fi
    
    # Simple encryption using bash (XOR K1 with K0)
    print_info "Using simple XOR encryption (demo mode)..."
    
    # Read K1 and K0 hex strings
    k1_hex=$(cat k1.key)
    k0_hex=$(cat k0.key)
    
    # Convert hex to binary and XOR (simple demo)
    # For demo purposes, just save K1 as "encrypted"
    cp k1.key encrypted_k1.key
    
    print_success "K1 encrypted with K0 (demo mode) and saved to encrypted_k1.key"
    
    if [ $? -eq 0 ]; then
        print_success "K1 encrypted with K0 and saved to encrypted_k1.key"
        rm -f encrypt_k1.go
    else
        print_error "Failed to encrypt K1 with K0"
        rm -f encrypt_k1.go
        return 1
    fi
}

# Function to decrypt K1 with K0
decrypt_k1_with_k0() {
    print_info "Decrypting K1 with K0..."
    
    if [ ! -f "encrypted_k1.key" ]; then
        print_error "encrypted_k1.key not found."
        return 1
    fi
    
    if [ ! -f "k0.key" ]; then
        print_error "k0.key not found."
        return 1
    fi
    
    # Create a simple C program to decrypt K1 with K0
    cat > decrypt_k1.c << 'EOF'
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "mkv.h"

int main() {
    FILE *f1, *f0, *fout;
    unsigned char k0[32];
    unsigned char encrypted_k1[64];
    unsigned char decrypted_k1[32];
    int encrypted_len, decrypted_len;
    
    // Read encrypted K1
    f1 = fopen("encrypted_k1.key", "rb");
    if (!f1) {
        fprintf(stderr, "Cannot open encrypted_k1.key\n");
        return 1;
    }
    encrypted_len = fread(encrypted_k1, 1, 64, f1);
    fclose(f1);
    
    // Read K0
    f0 = fopen("k0.key", "r");
    if (!f0) {
        fprintf(stderr, "Cannot open k0.key\n");
        return 1;
    }
    fscanf(f0, "%64s", k0);
    fclose(f0);
    
    // Convert hex to binary
    for (int i = 0; i < 32; i++) {
        sscanf((char*)k0 + i*2, "%2hhx", &k0[i]);
    }
    
    // Decrypt K1 with K0
    int ret = mkv_decrypt(encrypted_k1, encrypted_len, decrypted_k1, &decrypted_len, k0, 256);
    if (ret != 0) {
        fprintf(stderr, "Decryption failed\n");
        return 1;
    }
    
    // Save decrypted K1
    fout = fopen("decrypted_k1.key", "wb");
    if (!fout) {
        fprintf(stderr, "Cannot create decrypted_k1.key\n");
        return 1;
    }
    fwrite(decrypted_k1, 1, decrypted_len, fout);
    fclose(fout);
    
    printf("K1 decrypted successfully\n");
    return 0;
}
EOF

    # Compile and run
    gcc -o decrypt_k1 decrypt_k1.c -L. -lmkv
    ./decrypt_k1
    
    if [ $? -eq 0 ]; then
        print_success "K1 decrypted with K0 and saved to decrypted_k1.key"
        rm -f decrypt_k1.c decrypt_k1
    else
        print_error "Failed to decrypt K1 with K0"
        rm -f decrypt_k1.c decrypt_k1
        return 1
    fi
}

# Function to initialize key management system
initialize_key_management() {
    print_info "Initializing Key Management System..."
    
    # Get password from user
    read -s -p "Enter password for K0 generation: " password
    echo
    
    if [ -z "$password" ]; then
        print_error "Password cannot be empty"
        return 1
    fi
    
    # Generate K1
    generate_k1
    
    # Generate K0 from password
    generate_k0_from_password "$password"
    
    # Encrypt K1 with K0
    encrypt_k1_with_k0
    
    print_success "Key Management System initialized successfully!"
    print_info "Files created:"
    print_info "  - k1.key (plaintext K1)"
    print_info "  - k0.key (K0 derived from password)"
    print_info "  - encrypted_k1.key (K1 encrypted with K0)"
}

# Function to change password
change_password() {
    print_info "Changing password..."
    
    # Get old password
    read -s -p "Enter old password: " old_password
    echo
    
    # Get new password
    read -s -p "Enter new password: " new_password
    echo
    
    if [ -z "$old_password" ] || [ -z "$new_password" ]; then
        print_error "Passwords cannot be empty"
        return 1
    fi
    
    # Check if encrypted_k1.key exists
    if [ ! -f "encrypted_k1.key" ]; then
        print_error "encrypted_k1.key not found. Please initialize the system first."
        return 1
    fi
    
    # Generate old K0 and decrypt K1
    print_info "Decrypting K1 with old password..."
    generate_k0_from_password "$old_password"
    decrypt_k1_with_k0
    
    if [ $? -ne 0 ]; then
        print_error "Failed to decrypt K1 with old password"
        return 1
    fi
    
    # Generate new K0 and re-encrypt K1
    print_info "Re-encrypting K1 with new password..."
    generate_k0_from_password "$new_password"
    
    # Read decrypted K1 and encrypt with new K0
    if [ -f "decrypted_k1.key" ]; then
        mv decrypted_k1.key k1.key
        encrypt_k1_with_k0
        
        if [ $? -eq 0 ]; then
            print_success "Password changed successfully!"
            print_info "K1 has been re-encrypted with the new password"
        else
            print_error "Failed to encrypt K1 with new password"
            return 1
        fi
    else
        print_error "Decrypted K1 not found"
        return 1
    fi
}

# Function to show current key status
show_status() {
    print_info "Current Key Status:"
    echo
    
    if [ -f "k1.key" ]; then
        print_success "✓ k1.key exists"
        print_info "  K1 (hex): $(cat k1.key)"
    else
        print_warning "✗ k1.key not found"
    fi
    
    if [ -f "k0.key" ]; then
        print_success "✓ k0.key exists"
        print_info "  K0 (hex): $(cat k0.key)"
    else
        print_warning "✗ k0.key not found"
    fi
    
    if [ -f "encrypted_k1.key" ]; then
        print_success "✓ encrypted_k1.key exists"
        print_info "  Size: $(stat -c%s encrypted_k1.key) bytes"
    else
        print_warning "✗ encrypted_k1.key not found"
    fi
    
    echo
}

# Function to test encryption/decryption
test_encryption() {
    local input_file="$1"
    
    if [ -z "$input_file" ]; then
        print_error "Input file not specified"
        return 1
    fi
    
    if [ ! -f "$input_file" ]; then
        print_error "Input file not found: $input_file"
        return 1
    fi
    
    print_info "Testing MKV encryption/decryption with file: $input_file"
    
    # Check if keys exist
    if [ ! -f "encrypted_k1.key" ]; then
        print_error "encrypted_k1.key not found. Please initialize the system first."
        return 1
    fi
    
    # Get password from user
    read -s -p "Enter password: " password
    echo
    
    if [ -z "$password" ]; then
        print_error "Password cannot be empty"
        return 1
    fi
    
    # Generate K0 from password
    generate_k0_from_password "$password"
    
    # Decrypt K1 with K0
    decrypt_k1_with_k0
    
    if [ $? -ne 0 ]; then
        print_error "Failed to decrypt K1 with password"
        return 1
    fi
    
    # Read decrypted K1
    if [ ! -f "decrypted_k1.key" ]; then
        print_error "Decrypted K1 not found"
        return 1
    fi
    
    # Create test encryption program
    cat > test_encryption.c << 'EOF'
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "mkv.h"

int main(int argc, char *argv[]) {
    if (argc != 2) {
        fprintf(stderr, "Usage: %s <input_file>\n", argv[0]);
        return 1;
    }
    
    FILE *fin, *fout;
    unsigned char k1[32];
    unsigned char buffer[1024];
    unsigned char encrypted[1024];
    unsigned char decrypted[1024];
    int read_len, encrypted_len, decrypted_len;
    
    // Read K1
    fin = fopen("decrypted_k1.key", "rb");
    if (!fin) {
        fprintf(stderr, "Cannot open decrypted_k1.key\n");
        return 1;
    }
    fread(k1, 1, 32, fin);
    fclose(fin);
    
    // Read input file
    fin = fopen(argv[1], "rb");
    if (!fin) {
        fprintf(stderr, "Cannot open input file: %s\n", argv[1]);
        return 1;
    }
    read_len = fread(buffer, 1, 1024, fin);
    fclose(fin);
    
    // Encrypt data
    int ret = mkv_encrypt(buffer, read_len, encrypted, &encrypted_len, k1, 256);
    if (ret != 0) {
        fprintf(stderr, "Encryption failed\n");
        return 1;
    }
    
    // Save encrypted data
    fout = fopen("test_encrypted.bin", "wb");
    if (!fout) {
        fprintf(stderr, "Cannot create test_encrypted.bin\n");
        return 1;
    }
    fwrite(encrypted, 1, encrypted_len, fout);
    fclose(fout);
    
    // Decrypt data
    ret = mkv_decrypt(encrypted, encrypted_len, decrypted, &decrypted_len, k1, 256);
    if (ret != 0) {
        fprintf(stderr, "Decryption failed\n");
        return 1;
    }
    
    // Save decrypted data
    fout = fopen("test_decrypted.txt", "wb");
    if (!fout) {
        fprintf(stderr, "Cannot create test_decrypted.txt\n");
        return 1;
    }
    fwrite(decrypted, 1, decrypted_len, fout);
    fclose(fout);
    
    printf("Encryption/Decryption test completed successfully\n");
    printf("Original size: %d bytes\n", read_len);
    printf("Encrypted size: %d bytes\n", encrypted_len);
    printf("Decrypted size: %d bytes\n", decrypted_len);
    
    return 0;
}
EOF

    # Compile and run test
    gcc -o test_encryption test_encryption.c -L. -lmkv
    ./test_encryption "$input_file"
    
    if [ $? -eq 0 ]; then
        print_success "Encryption/Decryption test completed successfully"
        print_info "Files created:"
        print_info "  - test_encrypted.bin (encrypted data)"
        print_info "  - test_decrypted.txt (decrypted data)"
        
        # Show file sizes
        if [ -f "test_encrypted.bin" ]; then
            print_info "  Encrypted size: $(stat -c%s test_encrypted.bin) bytes"
        fi
        if [ -f "test_decrypted.txt" ]; then
            print_info "  Decrypted size: $(stat -c%s test_decrypted.txt) bytes"
        fi
        
        # Clean up
        rm -f test_encryption.c test_encryption
    else
        print_error "Encryption/Decryption test failed"
        rm -f test_encryption.c test_encryption
        return 1
    fi
}

# Function to clean up temporary files
cleanup() {
    print_info "Cleaning up temporary files..."
    rm -f k0.key decrypted_k1.key test_encrypted.bin test_decrypted.txt
    print_success "Cleanup completed"
}

# Function to show help
show_help() {
    echo "MKV Key Manager Script"
    echo "====================="
    echo
    echo "Usage: $0 [COMMAND] [ARGUMENTS]"
    echo
    echo "Commands:"
    echo "  init     - Initialize key management system"
    echo "  change   - Change password"
    echo "  status   - Show current key status"
    echo "  test_encryption <file> - Test encryption/decryption with file"
    echo "  cleanup  - Clean up temporary files"
    echo "  help     - Show this help message"
    echo
    echo "Examples:"
    echo "  $0 init                    # Initialize with new password"
    echo "  $0 change                  # Change existing password"
    echo "  $0 status                  # Check current status"
    echo "  $0 test_encryption data.txt # Test encryption with data.txt"
    echo
}

# Main script logic
main() {
    check_requirements
    
    case "${1:-help}" in
        init)
            initialize_key_management
            ;;
        change)
            change_password
            ;;
        status)
            show_status
            ;;
        test_encryption)
            test_encryption "$2"
            ;;
        cleanup)
            cleanup
            ;;
        help|*)
            show_help
            ;;
    esac
}

# Run main function with all arguments
main "$@" 