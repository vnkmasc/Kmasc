#!/bin/bash

echo "=== Environment Test ==="

# Test Go
echo "Testing Go..."
if command -v go >/dev/null 2>&1; then
    echo "✅ Go: $(go version)"
else
    echo "❌ Go not found"
fi

# Test GCC
echo "Testing GCC..."
if command -v gcc >/dev/null 2>&1; then
    echo "✅ GCC: $(gcc --version | head -n1)"
else
    echo "❌ GCC not found"
fi

# Test OpenSSL
echo "Testing OpenSSL..."
if pkg-config --modversion openssl >/dev/null 2>&1; then
    echo "✅ OpenSSL: $(pkg-config --modversion openssl)"
else
    echo "❌ OpenSSL not found"
fi

# Test Docker
echo "Testing Docker..."
if command -v docker >/dev/null 2>&1; then
    echo "✅ Docker: $(docker --version)"
else
    echo "❌ Docker not found"
fi

# Test CGO
echo "Testing CGO..."
if [ "$CGO_ENABLED" = "1" ]; then
    echo "✅ CGO is enabled"
else
    echo "❌ CGO is not enabled"
fi

echo "=== Test Complete ==="
