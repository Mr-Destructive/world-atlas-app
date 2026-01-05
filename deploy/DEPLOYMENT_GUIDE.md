# Deployment Guide

Deploy Vue frontend on Vercel and Go backend on Render.

## Prerequisites

- Vercel account (https://vercel.com)
- Render account (https://render.com)
- GitHub account with this repo
- Go 1.24+
- Node.js 18+

---

## 1. Deploy Go Backend on Render

### Step 1: Push to GitHub

```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/your-repo.git
git push -u origin main
```

### Step 2: Create Render Service

1. Go to https://render.com and sign in
2. Click "New+" → "Web Service"
3. Connect GitHub repo
4. Configure:
   - **Name**: wa-1-backend
   - **Root Directory**: `server`
   - **Environment**: Go
   - **Build Command**: `go build -o server_bin`
   - **Start Command**: `./server_bin`
5. Add Environment Variable (optional, if using different data path)
   - `DATA_PATH=./data`
6. Create Web Service

### Step 3: Update Frontend API URL

After Render deploys, you'll get a URL like `https://wa-1-backend.onrender.com`

Update `client/src/` WebSocket connection to use this URL (update proxy settings if needed).

---

## 2. Deploy Vue Frontend on Vercel

### Step 1: Install Vercel CLI

```bash
npm install -g vercel
```

### Step 2: Deploy

Run from project root:

```bash
cd client
vercel
```

Or connect GitHub repo at https://vercel.com/new

### Step 3: Configure

1. Select GitHub repo
2. Framework Preset: `Vite`
3. Build Command: `npm run build`
4. Output Directory: `dist`
5. Environment Variables (if needed):
   - `VITE_API_URL=https://wa-1-backend.onrender.com`

### Step 4: Update WebSocket Connection

In your Vue app, update the WebSocket URL to use the backend:

```typescript
// src/main.ts or where you connect
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
const ws = new WebSocket(`${apiUrl.replace('http', 'ws')}/ws`);
```

---

## 3. Update CORS (if needed)

In `server/main.go`, add CORS middleware:

```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "https://your-vercel-url.vercel.app")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        next.ServeHTTP(w, r)
    })
}
```

---

## 4. Update Data Path for Production

Render runs from `server/` directory. Update `main.go`:

```go
dataPath := os.Getenv("DATA_PATH")
if dataPath == "" {
    dataPath = filepath.Join(".", "data", "places.json")
}
dict, err := game.NewDictionary(dataPath)
```

---

## Deployment Scripts

See `deploy/` directory for automated scripts:

- `deploy-backend.sh` - Deploy Go backend
- `deploy-frontend.sh` - Deploy Vue frontend

---

## Troubleshooting

**WebSocket Connection Fails**
- Ensure backend URL is correct in frontend
- Check Render logs: https://dashboard.render.com
- Verify CORS settings

**Data File Not Found**
- Ensure `data/places.json` exists in `server/data/`
- Update `DATA_PATH` environment variable on Render

**Frontend Can't Reach Backend**
- Check that Render URL is accessible
- Update `VITE_API_URL` on Vercel
- Verify firewall rules
