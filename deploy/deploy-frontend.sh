#!/bin/bash

# Deploy Vue frontend to Vercel
# Usage: ./deploy-frontend.sh

set -e

echo "🚀 Deploying Vue frontend to Vercel..."

# Check if Vercel CLI is installed
if ! command -v vercel &> /dev/null; then
    echo "❌ Vercel CLI not found. Install it:"
    echo "   npm install -g vercel"
    exit 1
fi

cd client

# Build frontend
echo "📦 Building frontend..."
npm install
npm run build

# Deploy to Vercel
echo "🌐 Deploying to Vercel..."
vercel --prod

echo ""
echo "✅ Frontend deployed successfully!"
echo "Your site URL will be shown above"
