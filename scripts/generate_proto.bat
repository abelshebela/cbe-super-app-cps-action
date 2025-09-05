@echo off
echo Generating gRPC code from protobuf definitions...

REM Create the output directory if it doesn't exist
if not exist "internal\grpc\pb" mkdir "internal\grpc\pb"

REM Generate Go code from proto files
protoc --go_out=. --go_opt=paths=source_relative ^
       --go-grpc_out=. --go-grpc_opt=paths=source_relative ^
       proto/bank/bank.proto

echo gRPC code generation completed!
