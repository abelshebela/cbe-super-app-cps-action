#!/bin/bash

# Set path to .env file (you can pass it as an argument)
ENV_FILE=${1:-.env}

# Check if the .env file exists
if [ ! -f "$ENV_FILE" ]; then
  echo "❌ Error: $ENV_FILE not found!"
  exit 1
fi

# Read each line in the .env file
while IFS='=' read -r key value; do
  # Skip comments and empty lines
  if [[ -z "$key" || "$key" =~ ^# ]]; then
    continue
  fi

  # Remove surrounding quotes from the value
  value=$(echo "$value" | sed -e 's/^"//' -e 's/"$//' -e "s/^'//" -e "s/'$//")

  # Export the variable
  export "$key=$value"
done < "$ENV_FILE"

echo "✅ All environment variables from $ENV_FILE exported successfully."
