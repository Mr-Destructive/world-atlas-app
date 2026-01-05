#!/bin/bash

# Deploy Go backend to Render
# Usage: ./deploy-backend.sh

set -e

echo "📦 Building Go backend..."
cd server

# Build binary
go build -o server_bin

echo "✅ Backend built successfully!"
echo ""
echo "Manual deployment steps:"
echo "1. Commit and push to GitHub:"
echo "   git add ."
echo "   git commit -m 'Deploy backend'"
echo "   git push"
echo ""
echo "2. Go to Render dashboard:"
echo "   https://dashboard.render.com"
echo ""
echo "3. Your service will auto-deploy from the pushed code"
echo ""
echo "4. Get the backend URL from Render:"
echo "   https://wa-1-backend.onrender.com (or your custom name)"
