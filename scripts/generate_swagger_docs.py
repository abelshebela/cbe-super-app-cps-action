#!/usr/bin/env python3
"""
Generate Swagger/OpenAPI documentation from Go source code.

This script parses Go source files to extract API endpoint information
and generates or updates the Swagger/OpenAPI specification.
"""
import os
import re
import json
import argparse
from pathlib import Path
from typing import Dict, List, Optional, Any, Tuple

# Configuration
PROJECT_ROOT = Path(__file__).parent.parent
SWAGGER_FILE = PROJECT_ROOT / "docs" / "swagger" / "swagger.json"
HANDLERS_DIR = PROJECT_ROOT / "internal" / "handlers"
MODELS_DIR = PROJECT_ROOT / "internal" / "constants" / "dto"

# Regex patterns to extract API information
ROUTE_PATTERN = re.compile(r'router\.(?:Get|Post|Put|Delete|Patch|Options|Head|Connect|Trace)\([^,]+,\s*"([^"]+)"'
                          r'|r\.(?:Get|Post|Put|Delete|Patch|Options|Head|Connect|Trace)\([^,]+,\s*"([^"]+)"'
                          r'|chi\.Router\(\)\.(?:Get|Post|Put|Delete|Patch|Options|Head|Connect|Trace)\([^,]+,\s*"([^"]+)"'
                          r'|r\.With\([^)]+\)\.(?:Get|Post|Put|Delete|Patch|Options|Head|Connect|Trace)\([^,]+,\s*"([^"]+)"')

HANDLER_PATTERN = re.compile(r'func\s+([a-zA-Z0-9_]+)\s*\(')
COMMENT_PATTERN = re.compile(r'//\s*@(\w+)\s*(.*)')

# Swagger template
SWAGGER_TEMPLATE = {
    "openapi": "3.0.0",
    "info": {
        "title": "CPS Action API",
        "description": "API documentation for the CPS Action module of CBE Super App",
        "version": "1.0.0",
        "contact": {
            "name": "API Support",
            "email": "contact@eaglelionsystems.com"
        }
    },
    "servers": [
        {
            "url": "http://localhost:8080/api/v1/cbesuperapp/cps_action",
            "description": "Local Development Server"
        }
    ],
    "paths": {},
    "components": {
        "securitySchemes": {
            "BearerAuth": {
                "type": "http",
                "scheme": "bearer",
                "bearerFormat": "JWT"
            }
        },
        "schemas": {},
        "responses": {}
    },
    "security": [{"BearerAuth": []}]
}

class SwaggerGenerator:
    def __init__(self):
        self.swagger = SWAGGER_TEMPLATE.copy()
        self.current_path = ""
        self.current_method = ""
        self.current_component = ""
        
    def load_existing_swagger(self) -> bool:
        """Load existing Swagger file if it exists."""
        if SWAGGER_FILE.exists():
            try:
                with open(SWAGGER_FILE, 'r', encoding='utf-8') as f:
                    self.swagger = json.load(f)
                return True
            except json.JSONDecodeError:
                print(f"Warning: Could not parse {SWAGGER_FILE}. Creating a new one.")
        return False
    
    def save_swagger(self):
        """Save the Swagger specification to a file."""
        # Create directory if it doesn't exist
        SWAGGER_FILE.parent.mkdir(parents=True, exist_ok=True)
        
        # Write the Swagger file with pretty-printing
        with open(SWAGGER_FILE, 'w', encoding='utf-8') as f:
            json.dump(self.swagger, f, indent=2, ensure_ascii=False)
        
        print(f"Swagger documentation generated at {SWAGGER_FILE}")
    
    def process_go_files(self):
        """Process all Go files to extract API information."""
        # Process handler files
        for root, _, files in os.walk(HANDLERS_DIR):
            for file in files:
                if file.endswith('.go') and not file.endswith('_test.go'):
                    self.process_go_file(Path(root) / file)
        
        # Process model files
        for root, _, files in os.walk(MODELS_DIR):
            for file in files:
                if file.endswith('.go') and not file.endswith('_test.go'):
                    self.process_model_file(Path(root) / file)
    
    def process_go_file(self, file_path: Path):
        """Process a Go file to extract API endpoint information."""
        print(f"Processing {file_path}...")
        
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Extract route information
        for match in ROUTE_PATTERN.finditer(content):
            path = match.group(1) or match.group(3) or match.group(5) or match.group(7)
            method = match.group(0).split('.')[-1].split('(')[0].lower()
            
            # Initialize path if it doesn't exist
            if path not in self.swagger["paths"]:
                self.swagger["paths"][path] = {}
            
            # Add method to path
            self.swagger["paths"][path][method] = {
                "tags": [file_path.stem],
                "summary": f"{method.upper()} {path}",
                "responses": {
                    "200": {
                        "description": "Success",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/StandardResponse"
                                }
                            }
                        }
                    },
                    "default": {
                        "$ref": "#/components/responses/InternalServerError"
                    }
                },
                "security": [{"BearerAuth": []}]
            }
    
    def process_model_file(self, file_path: Path):
        """Process a Go file to extract model information."""
        print(f"Processing model {file_path}...")
        
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # This is a simplified example - in a real implementation, you would parse the Go structs
        # and convert them to JSON Schema
        
        # For now, we'll just add a placeholder for the model
        model_name = file_path.stem
        if model_name not in self.swagger["components"]["schemas"]:
            self.swagger["components"]["schemas"][model_name] = {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "string",
                        "description": f"Unique identifier for the {model_name}"
                    }
                }
            }

def main():
    parser = argparse.ArgumentParser(description='Generate Swagger/OpenAPI documentation from Go source code.')
    parser.add_argument('--output', '-o', default=str(SWAGGER_FILE),
                      help=f'Output file path (default: {SWAGGER_FILE})')
    args = parser.parse_args()
    
    # Initialize the generator
    generator = SwaggerGenerator()
    
    # Load existing Swagger file if it exists
    generator.load_existing_swagger()
    
    # Process Go files to extract API information
    generator.process_go_files()
    
    # Save the updated Swagger documentation
    generator.save_swagger()
    
    print("Swagger documentation generation complete!")

if __name__ == "__main__":
    main()
