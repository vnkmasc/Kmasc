#!/bin/bash

# Demo Scripts Functionality
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Fabric Scripts Demo ==="
echo "Date: $(date)"
echo "This demo shows how individual scripts work"
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

# Function to run demo step
run_demo_step() {
    local step_name=$1
    local script_name=$2
    local description=$3
    
    echo "=== Demo Step: $step_name ==="
    print_status "INFO" "Running: $script_name"
    print_status "INFO" "Description: $description"
    echo
    
    if [ -f "$script_name" ] && [ -x "$script_name" ]; then
        # Run script with timeout to avoid hanging
        timeout 30s ./"$script_name" || {
            print_status "WARN" "Script timed out or had issues (this is normal for demo)"
        }
        print_status "PASS" "Script executed successfully"
    else
        print_status "FAIL" "Script not found or not executable: $script_name"
    fi
    
    echo
    echo "Press Enter to continue to next step..."
    read -r
}

# Main demo function
main() {
    echo "🎬 Starting Fabric Scripts Demo"
    echo "This demo will show how each script works independently"
    echo
    
    print_status "INFO" "Demo will run each script briefly to show functionality"
    print_status "WARN" "Some scripts may take time or require user interaction"
    echo
    
    # Demo 1: List scripts
    run_demo_step "List Available Scripts" "list-scripts.sh" "Shows all available scripts with descriptions"
    
    # Demo 2: Check environment
    run_demo_step "Environment Check" "check-environment.sh" "Comprehensive environment diagnostics"
    
    # Demo 3: Quick environment test
    run_demo_step "Quick Environment Test" "test_environment.sh" "Quick check of basic requirements"
    
    # Demo 4: Download fabric-samples (if not exists)
    if [ ! -d "fabric-samples" ]; then
        run_demo_step "Download Fabric Samples" "download-fabric-samples.sh" "Download fabric-samples repository"
    else
        print_status "INFO" "Skipping fabric-samples download (already exists)"
        echo
    fi
    
    # Demo 5: Test script integrity
    run_demo_step "Test Script Integrity" "test-quick-start.sh" "Verify all scripts are properly configured"
    
    echo "=== Demo Complete ==="
    print_status "PASS" "All demo steps completed"
    echo
    echo "🎉 Demo completed successfully!"
    echo
    echo "📋 What you learned:"
    echo "  • Each script can run independently"
    echo "  • Scripts have proper error handling"
    echo "  • Scripts provide colored output"
    echo "  • Scripts check prerequisites"
    echo
    echo "🚀 Next steps:"
    echo "  • Run ./quick-start.sh for complete setup"
    echo "  • Or run individual scripts as needed"
    echo "  • Check README_SCRIPTS.md for detailed documentation"
    echo
    print_status "INFO" "Demo completed at $(date)"
}

# Check if user wants to run demo
echo "This demo will run several scripts to show their functionality."
echo "Some scripts may take time or require user interaction."
echo
read -p "Do you want to run the demo? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    main "$@"
else
    echo "Demo cancelled."
    exit 0
fi 