#!/bin/bash

# Test script for Kafka consumer with vault configuration

set -e

echo "🚀 Starting Kafka Consumer Test"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Start Kafka infrastructure
echo "📦 Starting Kafka infrastructure..."
docker-compose -f docker-compose.kafka.yml up -d zookeeper kafka mongodb topic-creator

# Wait for Kafka to be ready
echo "⏳ Waiting for Kafka to be ready..."
sleep 30

# Check if topics were created
echo "🔍 Checking if topics were created..."
docker exec kafka kafka-topics --bootstrap-server kafka:29092 --list

# Send test messages
echo "📤 Sending test feedback messages..."
docker-compose -f docker-compose.kafka.yml up feedback-test-producer

# Wait a moment for messages to be sent
sleep 5

# Test consumer
echo "📥 Testing consumer..."
docker-compose -f docker-compose.kafka.yml up feedback-test-consumer

# Test with Go application
echo "🧪 Testing with Go application..."

# Set environment variables for testing
export KAFKA_BROKERS="localhost:9092"
export KAFKA_FEEDBACK_TOPIC="feedback-events"
export KAFKA_CONSUMER_GROUP="test-consumer-group"
export VAULT_ADDR="http://localhost:8200"
export VAULT_TOKEN="test-token"
export VAULT_PATH="secret/kafka"

# Run the test
echo "🔧 Running Kafka consumer tests..."
go test ./tests/domain/feedback -v

echo "✅ Kafka consumer test completed!"

# Optional: Start Kafka UI for monitoring
echo "🌐 Starting Kafka UI for monitoring..."
docker-compose -f docker-compose.kafka.yml up -d kafka-ui
echo "📊 Kafka UI available at: http://localhost:8082"

echo "🎉 All tests completed successfully!" 