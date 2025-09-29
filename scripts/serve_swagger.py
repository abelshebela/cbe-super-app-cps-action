#!/usr/bin/env python3
"""
Simple HTTP server to serve Swagger UI documentation.
"""
import http.server
import socketserver
import os
import webbrowser
from pathlib import Path

# Configuration
PORT = 8000
DOCS_DIR = Path(__file__).parent.parent / 'docs' / 'swagger'

class CORSRequestHandler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=str(DOCS_DIR), **kwargs)
    
    def end_headers(self):
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type, Authorization')
        super().end_headers()
    
    def do_OPTIONS(self):
        self.send_response(200, "ok")
        self.end_headers()

def main():
    # Change to the docs directory
    os.chdir(DOCS_DIR)
    
    # Start the server
    with socketserver.TCPServer(("", PORT), CORSRequestHandler) as httpd:
        print(f"Serving Swagger UI at http://localhost:{PORT}")
        print("Press Ctrl+C to stop the server")
        
        # Open the browser automatically
        webbrowser.open(f"http://localhost:{PORT}")
        
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\nShutting down the server...")
            httpd.server_close()

if __name__ == "__main__":
    main()
