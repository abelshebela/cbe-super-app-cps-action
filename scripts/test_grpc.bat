@echo off
echo Testing gRPC Bank Service with grpcurl
echo =====================================

echo.
echo 1. List available services:
grpcurl -plaintext localhost:9090 list

echo.
echo 2. List methods for BankService:
grpcurl -plaintext localhost:9090 list BankService

echo.
echo 3. Describe GetBank method:
grpcurl -plaintext localhost:9090 describe BankService.GetBank

echo.
echo 4. Test GetAllBanks (example call):
grpcurl -plaintext -d "{\"page\": 1, \"per_page\": 10}" localhost:9090 BankService.GetAllBanks

echo.
echo 5. Test GetBank by ID (example call):
grpcurl -plaintext -d "{\"id\": \"1\"}" localhost:9090 BankService.GetBank

echo.
echo gRPC testing completed!
