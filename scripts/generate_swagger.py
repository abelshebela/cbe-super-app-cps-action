import os
import json
import re
from typing import Dict, List, Optional, Any

# Configuration
PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
HANDLERS_DIR = os.path.join(PROJECT_ROOT, "internal", "handlers", "rest", "http")
OUTPUT_DIR = os.path.join(PROJECT_ROOT, "docs", "swagger")
SWAGGER_FILE = os.path.join(OUTPUT_DIR, "swagger.json")

# Initialize Swagger spec
swagger = {
    "openapi": "3.0.0",
    "info": {
        "title": "CPS Action API",
        "description": "API documentation for CPS Action service",
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
        "schemas": {
            "StandardResponse": {
                "type": "object",
                "properties": {
                    "ok": {"type": "boolean"},
                    "status": {"type": "integer"},
                    "timestamp": {"type": "string", "format": "date-time"},
                    "message": {"type": "string"},
                    "data": {"type": "object"},
                    "error": {"$ref": "#/components/schemas/ErrorDetail"}
                }
            },
            "ErrorDetail": {
                "type": "object",
                "properties": {
                    "code": {"type": "string"},
                    "message": {"type": "string"},
                    "status_code": {"type": "integer"},
                    "type": {"type": "string"},
                    "field_errors": {
                        "type": "array",
                        "items": {"$ref": "#/components/schemas/FieldError"}
                    },
                    "details": {"type": "object"}
                }
            },
            "FieldError": {
                "type": "object",
                "properties": {
                    "field": {"type": "string"},
                    "message": {"type": "string"},
                    "value": {"type": "string"},
                    "constraint": {"type": "string"}
                }
            }
        }
    },
    "security": [{"BearerAuth": []}]
}

def parse_go_file(file_path: str) -> Dict[str, Any]:
    """Parse a Go file to extract API endpoint information."""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # Extract route definitions
    routes = []
    for match in re.finditer(r'router\.(?:Get|Post|Put|Delete|Patch|Options|Head)\([^,]+,\s*"([^"]+)"', content):
        routes.append({
            "method": match.group(0).split('.')[1].split('(')[0].lower(),
            "path": match.group(1)
        })

    # Extract handler function documentation
    docs = {}
    for match in re.finditer(r'//\s*@(\w+)\s+(.*?)(?:\n\s*@|$)', content):
        key = match.group(1).lower()
        value = match.group(2).strip()
        if key not in docs:
            docs[key] = value
        elif isinstance(docs[key], list):
            docs[key].append(value)
        else:
            docs[key] = [docs[key], value]

    return {"routes": routes, "docs": docs}

def generate_path_item(method: str, path: str, docs: Dict) -> Dict:
    """Generate a path item object for the Swagger spec."""
    path_item = {
        method.lower(): {
            "summary": docs.get("summary", ""),
            "description": docs.get("description", ""),
            "tags": [docs.get("tag", "default")],
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
                "400": {
                    "description": "Bad Request",
                    "content": {
                        "application/json": {
                            "schema": {
                                "$ref": "#/components/schemas/ErrorDetail"
                            }
                        }
                    }
                },
                "401": {
                    "description": "Unauthorized"
                },
                "403": {
                    "description": "Forbidden"
                },
                "500": {
                    "description": "Internal Server Error"
                }
            }
        }
    }

    # Add parameters if any
    if "param" in docs:
        params = docs["param"] if isinstance(docs["param"], list) else [docs["param"]]
        path_item[method.lower()]["parameters"] = []
        for param in params:
            param_parts = param.split()
            if len(param_parts) >= 3:
                param_type = param_parts[0]
                param_name = param_parts[1]
                param_desc = " ".join(param_parts[2:]) if len(param_parts) > 2 else ""
                
                param_def = {
                    "name": param_name,
                    "in": param_type,
                    "description": param_desc,
                    "required": True,
                    "schema": {
                        "type": "string"
                    }
                }
                path_item[method.lower()]["parameters"].append(param_def)

    # Add request body if needed
    if method.lower() in ["post", "put", "patch"] and "body" in docs:
        path_item[method.lower()]["requestBody"] = {
            "description": "Request body",
            "required": True,
            "content": {
                "application/json": {
                    "schema": {
                        "type": "object"
                    }
                }
            }
        }

    return path_item

def main():
    # Create output directory if it doesn't exist
    os.makedirs(OUTPUT_DIR, exist_ok=True)

    # Process all handler files
    for root, _, files in os.walk(HANDLERS_DIR):
        for file in files:
            if file.endswith('.go') and not file.endswith('_test.go'):
                file_path = os.path.join(root, file)
                print(f"Processing {file_path}...")
                
                try:
                    file_info = parse_go_file(file_path)
                    
                    # Add paths to Swagger spec
                    for route in file_info.get("routes", []):
                        path = route["path"]
                        method = route["method"]
                        
                        # Normalize path
                        path = path.replace("//", "/")
                        
                        # Add to paths
                        if path not in swagger["paths"]:
                            swagger["paths"][path] = {}
                        
                        # Add method to path
                        swagger["paths"][path].update(
                            generate_path_item(method, path, file_info.get("docs", {}))
                        )
                except Exception as e:
                    print(f"Error processing {file_path}: {str(e)}")
    
    # Write Swagger spec to file
    with open(SWAGGER_FILE, 'w', encoding='utf-8') as f:
        json.dump(swagger, f, indent=2)
    
    print(f"Swagger documentation generated at {SWAGGER_FILE}")

if __name__ == "__main__":
    main()
