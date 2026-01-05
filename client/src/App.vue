<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import Auth from './components/Auth.vue'
import Home from './components/Home.vue'
import Lobby from './components/Lobby.vue'
import Game from './components/Game.vue'
import GameOver from './components/GameOver.vue'

interface User {
    id: string
    username: string
    totalScore: number
    wins: number
}

interface Player {
  id: string
  name: string
  lives: number
  isTurn: boolean
  type: number // 0 human, 1 bot
  score: number
  mostUsedPlaces?: Record<string, number>
}

interface Move {
  playerId: string
  playerName: string
  word: string
  type: string
  timestamp: number
}

interface GameState {
  players: Record<string, Player>
  state: string
  lastWord: string
  turnOrder: string[]
  currentTurn: string
  history: Move[]
  round: number
}

const socket = ref<WebSocket | null>(null)
const connected = ref(false)
const gameState = ref<GameState | null>(null)
const myId = ref<string>("")
const roomId = ref<string>("")
const errorMsg = ref("")

const user = ref<User | null>(null)
const isGuest = ref(false)

onMounted(() => {
    const savedUser = localStorage.getItem('wa_user')
    if (savedUser) {
        user.value = JSON.parse(savedUser)
    }
})

const handleLogin = (userData: User) => {
    user.value = userData
    localStorage.setItem('wa_user', JSON.stringify(userData))
    isGuest.value = false
}

const handleGuest = () => {
    isGuest.value = true
}

const handleLogout = () => {
    user.value = null
    localStorage.removeItem('wa_user')
    isGuest.value = false
    quit()
}

// View State: 'AUTH' | 'HOME' | 'LOBBY' | 'GAME' | 'ENDED'
const viewState = computed(() => {
  if (!user.value && !isGuest.value) return 'AUTH'
  if (!connected.value) return 'HOME'
  if (gameState.value?.state === 'PLAYING') return 'GAME'
  if (gameState.value?.state === 'ENDED') return 'ENDED'
  return 'LOBBY'
})

const sortedPlayers = computed((): Player[] => {
  if (!gameState.value) return []
  return gameState.value.turnOrder
    .map(id => gameState.value!.players[id])
    .filter((p): p is Player => !!p)
})

const connect = ({ name, roomId: rId }: { name: string, roomId: string }) => {
  const finalName = user.value ? user.value.username : name
  
  // Get API URL from environment variable or default to current host
  const apiUrl = import.meta.env.VITE_API_URL || window.location.origin
  const protocol = apiUrl.includes('https') ? 'wss:' : 'ws:'
  const host = apiUrl.replace(/^https?:\/\//, '')
  
  const wsUrl = `${protocol}//${host}/ws?name=${encodeURIComponent(finalName)}&room=${encodeURIComponent(rId)}`
  
  socket.value = new WebSocket(wsUrl)

  socket.value.onopen = () => {
    connected.value = true
    console.log("Connected")
  }

  socket.value.onmessage = (event) => {
    const data = JSON.parse(event.data)
    if (data.type === 'WELCOME') {
      myId.value = data.payload.id
      roomId.value = data.payload.roomId
    } else if (data.type === 'GAME_STATE') {
      gameState.value = data.payload
    } else if (data.type === 'ERROR') {
        errorMsg.value = data.payload.message
        setTimeout(() => errorMsg.value = "", 3000)
    }
  }

  socket.value.onclose = () => {
    connected.value = false
    gameState.value = null
    roomId.value = ""
  }
}

const startGame = (opts?: any) => {
    // Pass mode and settings from Lobby
    const payload = opts || {}
    socket.value?.send(JSON.stringify({ type: 'START_GAME', payload }))
}

const addBot = () => {
  socket.value?.send(JSON.stringify({ type: 'ADD_BOT', payload: {} }))
}

const submitWord = (word: string) => {
  socket.value?.send(JSON.stringify({ 
    type: 'SUBMIT_WORD', 
    payload: { word } 
  }))
}

const quit = () => {
  socket.value?.close()
}
</script>

<template>
  <div class="h-screen bg-duo-gray flex flex-col items-center p-4 overflow-hidden">
    <div class="w-full max-w-6xl flex items-center justify-between mb-4 shrink-0">
        <h1 class="text-3xl font-black text-duo-green tracking-wide drop-shadow-sm">World Atlas</h1>
        <div v-if="user" class="flex items-center gap-4">
            <div class="text-right hidden sm:block">
                <div class="font-black text-gray-700">{{ user.username }}</div>
                <div class="text-[10px] text-gray-400 font-bold uppercase tracking-wider">Score: {{ user.totalScore }} • Wins: {{ user.wins }}</div>
            </div>
            <button @click="handleLogout" class="text-xs font-black text-red-400 hover:text-red-600 uppercase tracking-widest border-2 border-red-100 rounded-lg px-3 py-1 hover:bg-red-50 transition-colors">
                Logout
            </button>
        </div>
    </div>

    <!-- Global Error -->
    <div v-if="errorMsg" class="fixed top-4 right-4 z-50 p-4 bg-red-100 text-duo-red font-bold rounded-xl border-l-4 border-duo-red shadow-md animate-bounce">
       {{ errorMsg }}
    </div>

    <!-- Views -->
    <div class="w-full max-w-6xl flex-1 flex justify-center min-h-0">
        <Auth v-if="viewState === 'AUTH'" @login="handleLogin" @guest="handleGuest" />

        <Home v-else-if="viewState === 'HOME'" :user-name="user?.username" @join="connect" />
        
        <Lobby v-else-if="viewState === 'LOBBY'" 
               :room-id="roomId" 
               :players="sortedPlayers" 
               :my-id="myId"
               @start="startGame"
               @add-bot="addBot"
               @quit="quit" />
               
        <Game v-else-if="viewState === 'GAME' && gameState" 
              :game-state="gameState" 
              :my-id="myId" 
              @submit="submitWord"
              @quit="quit" />

        <GameOver v-else-if="viewState === 'ENDED' && gameState"
                  :game-state="gameState"
                  :my-id="myId"
                  @start="startGame"
                  @quit="quit" />
    </div>
  </div>
</template>