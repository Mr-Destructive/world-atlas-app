# Quick Deploy - 5 Minutes

## 1️⃣ Push to GitHub

```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/wa-1.git
git push -u origin main
```

## 2️⃣ Deploy Backend (Go) on Render

1. Go to https://render.com/dashboard
2. Click **New +** → **Web Service**
3. Select your GitHub repo
4. Fill in:
   - **Service Name**: `wa-1-backend`
   - **Root Directory**: `server`
   - **Runtime**: Go
   - **Build Command**: `go build -o server_bin`
   - **Start Command**: `./server_bin`
5. Click **Create Web Service**
6. Wait 2-3 minutes, then get your URL: `https://wa-1-backend.onrender.com` (or custom)

## 3️⃣ Deploy Frontend (Vue) on Vercel

**Option A: Using Dashboard**
1. Go to https://vercel.com/new
2. Import GitHub repo
3. Select `client` as root directory
4. Click **Deploy**

**Option B: Using CLI**
```bash
npm install -g vercel
cd client
vercel --prod
```

## 4️⃣ Update Frontend to Connect to Backend

In `client/src/` (wherever you connect WebSocket), update to:

```typescript
const backendUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
const ws = new WebSocket(`${backendUrl.replace('http', 'ws')}/ws`);
```

On Vercel dashboard, add Environment Variable:
- **Key**: `VITE_API_URL`
- **Value**: `https://wa-1-backend.onrender.com`

Then redeploy frontend.

## ✅ Done!

Your app is now live:
- Backend: `https://wa-1-backend.onrender.com`
- Frontend: `https://wa-1.vercel.app` (or your custom domain)

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| WebSocket connection fails | Verify backend URL in Vercel env vars |
| Data file not found | Ensure `server/data/places.json` exists |
| 502 Bad Gateway | Check Render service logs |
| Frontend blank | Check browser console for errors |

See `deploy/DEPLOYMENT_GUIDE.md` for detailed setup.
