# Vercel Deployment Guide (Step-by-Step)

Deploy Vue frontend to Vercel and connect to your Render backend.

## Prerequisites

- Vercel account: https://vercel.com/signup
- GitHub account
- Code pushed to GitHub
- Backend deployed at: `https://meetgor-api.onrender.com`

---

## Step 1: Update WebSocket Connection in Vue

The frontend needs to connect to your Render backend API.

### 1.1 Find where WebSocket is initialized

Search your Vue components for WebSocket connection code. Usually in:
- `client/src/main.ts`
- `client/src/components/Game.vue`
- Or any component that connects to `/ws`

### 1.2 Update the connection to use environment variable

Replace hardcoded `localhost:8080` with:

```typescript
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
const wsUrl = apiUrl.replace('http', 'ws');
const ws = new WebSocket(`${wsUrl}/ws?name=${playerName}&room=${roomId}`);
```

**Example in a Vue component:**

```vue
<script setup lang="ts">
import { ref } from 'vue';

const playerName = ref('Guest');
const roomId = ref('');

const connectToGame = () => {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
  const wsUrl = apiUrl.replace('http', 'ws').replace('https', 'wss');
  
  const ws = new WebSocket(
    `${wsUrl}/ws?name=${playerName.value}&room=${roomId.value}`
  );
  
  ws.onopen = () => console.log('Connected!');
  ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    console.log('Message:', msg);
  };
};
</script>
```

---

## Step 2: Build Frontend Locally

Verify the build works before deploying:

```bash
cd client
npm install
npm run build
```

You should see `client/dist/` directory created with built files.

---

## Step 3: Commit & Push to GitHub

```bash
cd /home/meet/code/playground/go/wa-1
git add client/src/
git commit -m "Update WebSocket to use environment variable"
git push origin main
```

---

## Step 4: Deploy on Vercel Dashboard

### Option A: Using Vercel Dashboard (Recommended)

**Step 4.1: Connect Repository**

1. Go to https://vercel.com/new
2. Click **"Import Project"**
3. Paste GitHub repo URL or select from list
4. Click **"Import"**

**Step 4.2: Configure Project**

1. **Framework Preset**: Select `Vite`
2. **Build and Output Settings**:
   - **Build Command**: `npm run build`
   - **Output Directory**: `dist`
   - **Install Command**: `npm install`
3. **Root Directory**: Select `./client` (important!)
4. Click **"Deploy"**

**Step 4.3: Add Environment Variable**

After deployment starts, go to **Settings** → **Environment Variables**:

1. Click **"Add"**
2. **Name**: `VITE_API_URL`
3. **Value**: `https://meetgor-api.onrender.com`
4. **Environments**: Select `Production`, `Preview`, `Development`
5. Click **"Save"**

**Step 4.4: Redeploy with Environment Variable**

Go to **Deployments** → Latest deployment → **Redeploy**

---

## Step 5: Verify Deployment

### 5.1 Check Vercel URL

After redeployment, Vercel provides a URL like:
```
https://world-atlas-app.vercel.app
```

### 5.2 Test the App

Open https://world-atlas-app.vercel.app and:

1. Try to start a game
2. Make a guess
3. Check browser console (F12) for WebSocket connection
4. Should see: `wss://meetgor-api.onrender.com/ws` in Network tab

### 5.3 Troubleshoot Connection Issues

**If WebSocket fails:**

Check browser console:
```
Failed to connect to wss://meetgor-api.onrender.com/ws
```

**Solutions:**

1. Verify `VITE_API_URL` environment variable is set
2. Check Render backend is still running: https://meetgor-api.onrender.com/
3. Try in incognito mode (clear cache)
4. Update CORS in backend if needed:

```go
// In server/main.go
w.Header().Set("Access-Control-Allow-Origin", "https://world-atlas-app.vercel.app")
```

---

## Step 6: Set Up Custom Domain (Optional)

In Vercel Dashboard:

1. Go to **Settings** → **Domains**
2. Add your custom domain (e.g., `worldatlas.com`)
3. Follow DNS setup instructions

---

## Option B: Using Vercel CLI

If you prefer command line:

```bash
# Install Vercel CLI
npm install -g vercel

# Login
vercel login

# Deploy from client folder
cd client
vercel --prod

# Add environment variable when prompted
# Or via dashboard after
```

---

## Complete URLs Reference

| Service | URL |
|---------|-----|
| **Backend API** | `https://meetgor-api.onrender.com` |
| **WebSocket** | `wss://meetgor-api.onrender.com/ws` |
| **Frontend** | `https://world-atlas-app.vercel.app` |
| **Vercel Dashboard** | https://vercel.com/dashboard |
| **Render Dashboard** | https://dashboard.render.com |

---

## Testing Checklist

- [ ] Frontend loads at Vercel URL
- [ ] Can start a new game
- [ ] Can join a game room
- [ ] Can make guesses
- [ ] WebSocket connects (check Network tab → WS)
- [ ] Real-time messages appear
- [ ] No CORS errors in console

---

## Rollback / Updates

**To update frontend:**

```bash
cd client
git commit -am "Update feature X"
git push origin main
```

Vercel automatically redeploys on push to `main` branch.

**To disable auto-deploy:**

1. Vercel Dashboard → Settings → Git
2. Toggle "Auto-deploy" off

---

## Common Issues

| Issue | Solution |
|-------|----------|
| `Cannot find environment variable VITE_API_URL` | Make sure to add it in Vercel Settings → Environment Variables |
| WebSocket connection timeout | Verify Render backend is running |
| `404 Not Found` | Check root directory is set to `./client` |
| Build fails | Run `npm install` locally and check for errors |
| Old version still showing | Clear browser cache (Ctrl+Shift+Del) |

---

## Next Steps

1. Share your Vercel URL: `https://world-atlas-app.vercel.app`
2. Share the CLI client for testing: `./cli-client wss://meetgor-api.onrender.com/ws`
3. Monitor logs: Vercel Dashboard → Deployments → Logs
