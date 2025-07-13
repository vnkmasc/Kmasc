#!/bin/bash

echo "=========================================="
echo "SO SÁNH PERFORMANCE: CÓ AES vs KHÔNG AES"
echo "=========================================="

echo ""
echo "1. Chạy test CÓ AES:"
echo "----------------------------------------"
cd test_perf_with_aes
go run main.go

echo ""
echo "2. Chạy test KHÔNG AES:"
echo "----------------------------------------"
cd ../test_perf_without_aes
go run main.go

echo ""
echo "=========================================="
echo "KẾT LUẬN:"
echo "- Test có AES sẽ sử dụng encryption thật"
echo "- Test không AES sẽ chạy trực tiếp không mã hóa"
echo "==========================================" 