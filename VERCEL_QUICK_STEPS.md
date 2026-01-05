# Vercel Deployment - Quick Steps

## ✅ What's Already Done

- [x] Vue app updated to use `VITE_API_URL` environment variable
- [x] WebSocket connection updated in `client/src/App.vue`
- [x] Backend deployed on Render: `https://meetgor-api.onrender.com`
- [x] Vercel config created: `client/vercel.json`

## 📋 Deployment Steps

### Step 1: Commit Updated Code
```bash
cd /home/meet/code/playground/go/wa-1
git add client/src/App.vue
git commit -m "Update WebSocket to use VITE_API_URL environment variable"
git push origin main
```

### Step 2: Deploy on Vercel

**Option A: Dashboard (Easiest)**

1. Go to https://vercel.com/new
2. Click **Import Project**
3. Select your GitHub repo: `world-atlas-app`
4. Configure:
   - **Root Directory**: `./client`
   - **Framework**: Vite
   - **Build Command**: `npm run build`
   - **Output Directory**: `dist`
5. Click **Deploy**

**Option B: CLI**

```bash
npm install -g vercel
cd client
vercel --prod
```

### Step 3: Add Environment Variable

In Vercel Dashboard:

1. Go to your project
2. **Settings** → **Environment Variables**
3. Add:
   - **Name**: `VITE_API_URL`
   - **Value**: `https://meetgor-api.onrender.com`
   - **Environments**: Production, Preview, Development
4. Click **Save**

### Step 4: Redeploy

1. Go to **Deployments**
2. Click the latest deployment
3. Click **Redeploy** button

Wait 2-3 minutes for build to complete.

### Step 5: Test

1. Open your Vercel URL (e.g., `https://world-atlas-app.vercel.app`)
2. Start a game
3. Open DevTools (F12) → Network tab
4. Look for `wss://meetgor-api.onrender.com/ws` connection
5. Should see messages flowing

---

## 🔗 Your URLs

| Service | URL |
|---------|-----|
| Backend | `https://meetgor-api.onrender.com` |
| Frontend | `https://world-atlas-app.vercel.app` |

---

## ❌ Troubleshooting

**WebSocket connection fails**
- Check browser console (F12) for errors
- Verify `VITE_API_URL` is set in Vercel Settings
- Make sure Render backend is running

**Build fails on Vercel**
- Check build logs: Deployments → click deployment → Logs
- Run `npm install && npm run build` locally to debug

**Old version showing**
- Clear browser cache: Ctrl+Shift+Del
- Hard refresh: Ctrl+F5

**404 on root page**
- Ensure root directory is set to `./client` in Vercel settings

---

## 📝 Files Modified

- `client/src/App.vue` - Updated WebSocket connection
- `client/vercel.json` - Vercel configuration
- `client/.vercelignore` - Files to ignore
- `server/main.go.prod` - Production-ready Go server

---

## 🎮 Testing the CLI

You can also test with CLI client locally:

```bash
cd cli-client
./cli-client wss://meetgor-api.onrender.com/ws TestPlayer
```

Then try commands:
- `/start` - Start game
- `/guess Paris` - Make a guess
- `/quit` - Exit
