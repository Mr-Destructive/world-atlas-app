#!/bin/bash
# Build Client
echo "Building Client..."
cd client
npm run build
cd ..

# Run Server
echo "Starting Server..."
cd server
go run .
