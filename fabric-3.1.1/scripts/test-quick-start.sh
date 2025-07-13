#!/bin/bash

# Test Quick Start Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Testing Quick Start Script ==="
echo "Date: $(date)"
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "PASS")
            echo -e "${GREEN}✅ PASS${NC}: $message"
            ;;
        "FAIL")
            echo -e "${RED}❌ FAIL${NC}: $message"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ️  INFO${NC}: $message"
            ;;
        "WARN")
            echo -e "${YELLOW}⚠️  WARN${NC}: $message"
            ;;
    esac
}

# Function to check if script exists and is executable
check_script() {
    local script_name=$1
    local description=$2
    
    if [ -f "$script_name" ] && [ -x "$script_name" ]; then
        print_status "PASS" "$description ($script_name)"
        return 0
    else
        print_status "FAIL" "$description ($script_name) - missing or not executable"
        return 1
    fi
}

# Test all required scripts
test_scripts() {
    echo "Testing all required scripts..."
    echo
    
    local all_passed=true
    
    # Main scripts
    check_script "quick-start.sh" "Main quick start script" || all_passed=false
    check_script "setup-environment.sh" "Environment setup" || all_passed=false
    check_script "build-fabric.sh" "Fabric build" || all_passed=false
    check_script "start-network.sh" "Network startup" || all_passed=false
    
    # Utility scripts
    check_script "download-fabric-samples.sh" "Download fabric-samples" || all_passed=false
    check_script "build-encryption.sh" "Build encryption library" || all_passed=false
    check_script "test-encryption.sh" "Test encryption integration" || all_passed=false
    check_script "check-environment.sh" "Environment check" || all_passed=false
    check_script "fix-repositories.sh" "Fix repositories" || all_passed=false
    check_script "test_environment.sh" "Quick environment test" || all_passed=false
    check_script "list-scripts.sh" "List scripts" || all_passed=false
    
    echo
    if [ "$all_passed" = true ]; then
        print_status "PASS" "All scripts are present and executable"
    else
        print_status "FAIL" "Some scripts are missing or not executable"
        exit 1
    fi
}

# Test script syntax
test_syntax() {
    echo "Testing script syntax..."
    echo
    
    local all_passed=true
    
    for script in *.sh; do
        if [ -f "$script" ]; then
            if bash -n "$script" 2>/dev/null; then
                print_status "PASS" "Syntax check: $script"
            else
                print_status "FAIL" "Syntax error in: $script"
                all_passed=false
            fi
        fi
    done
    
    echo
    if [ "$all_passed" = true ]; then
        print_status "PASS" "All scripts have valid syntax"
    else
        print_status "FAIL" "Some scripts have syntax errors"
        exit 1
    fi
}

# Test function availability in quick-start.sh
test_functions() {
    echo "Testing function availability in quick-start.sh..."
    echo
    
    local all_passed=true
    
    # Check if required functions exist in quick-start.sh
    if grep -q "run_script" quick-start.sh; then
        print_status "PASS" "run_script function found"
    else
        print_status "FAIL" "run_script function not found"
        all_passed=false
    fi
    
    if grep -q "step1_fix_repositories" quick-start.sh; then
        print_status "PASS" "step1_fix_repositories function found"
    else
        print_status "FAIL" "step1_fix_repositories function not found"
        all_passed=false
    fi
    
    if grep -q "step2_setup_environment" quick-start.sh; then
        print_status "PASS" "step2_setup_environment function found"
    else
        print_status "FAIL" "step2_setup_environment function not found"
        all_passed=false
    fi
    
    if grep -q "step3_download_fabric_samples" quick-start.sh; then
        print_status "PASS" "step3_download_fabric_samples function found"
    else
        print_status "FAIL" "step3_download_fabric_samples function not found"
        all_passed=false
    fi
    
    if grep -q "step5_build_encryption" quick-start.sh; then
        print_status "PASS" "step5_build_encryption function found"
    else
        print_status "FAIL" "step5_build_encryption function not found"
        all_passed=false
    fi
    
    if grep -q "step6_test_encryption" quick-start.sh; then
        print_status "PASS" "step6_test_encryption function found"
    else
        print_status "FAIL" "step6_test_encryption function not found"
        all_passed=false
    fi
    
    if grep -q "step7_build_fabric" quick-start.sh; then
        print_status "PASS" "step7_build_fabric function found"
    else
        print_status "FAIL" "step7_build_fabric function not found"
        all_passed=false
    fi
    
    if grep -q "step8_start_network" quick-start.sh; then
        print_status "PASS" "step8_start_network function found"
    else
        print_status "FAIL" "step8_start_network function not found"
        all_passed=false
    fi
    
    echo
    if [ "$all_passed" = true ]; then
        print_status "PASS" "All required functions found in quick-start.sh"
    else
        print_status "FAIL" "Some required functions missing in quick-start.sh"
        exit 1
    fi
}

# Main test execution
main() {
    echo "Starting quick start script tests..."
    echo
    
    test_scripts
    test_syntax
    test_functions
    
    echo
    print_status "PASS" "All tests completed successfully!"
    print_status "INFO" "Quick start script is ready to use"
    echo
    echo "🎉 Ready to run: ./quick-start.sh"
}

# Run main function
main "$@" 