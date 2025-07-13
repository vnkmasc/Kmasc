#!/bin/bash

# Quick fix for broken repositories
echo "=== Fixing Broken Repositories ==="

# Remove the problematic PPA
echo "Removing problematic PPA..."
sudo add-apt-repository --remove ppa:ubuntu-vn/ppa -y 2>/dev/null || true

# Clean up any broken packages
echo "Cleaning up..."
sudo apt-get clean
sudo apt-get autoclean

# Update with fix-missing
echo "Updating package lists..."
sudo apt-get update --fix-missing

echo "✅ Repository fix completed!"
echo "Now you can run: ./setup-environment.sh" 