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
const chatMessages = ref<any[]>([])

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
    } else if (data.type === 'CHAT_MESSAGE') {
      chatMessages.value.push(data.payload)
    } else if (data.type === 'ERROR') {
        errorMsg.value = data.payload.message
        setTimeout(() => errorMsg.value = "", 3000)
    }
  }

  socket.value.onclose = () => {
    connected.value = false
    gameState.value = null
    roomId.value = ""
    chatMessages.value = []
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
  <div class="h-screen flex flex-col items-center p-2 lg:p-4 overflow-hidden">
    <div class="w-full max-w-6xl flex items-center justify-between mb-2 lg:mb-6 shrink-0 bg-white/50 backdrop-blur-md p-3 rounded-2xl border-2 border-white/50 shadow-sm z-50">
        <div class="flex items-center gap-3">
            <img src="/mascot.svg" alt="Mascot" class="w-8 h-8 lg:w-10 lg:h-10 drop-shadow-sm hover:rotate-12 transition-transform cursor-pointer" />
            <h1 class="text-xl lg:text-2xl font-black text-duo-green tracking-wide drop-shadow-sm uppercase">World Atlas</h1>
        </div>
        <div v-if="user" class="flex items-center gap-4">
            <div class="text-right hidden sm:block">
                <div class="font-black text-gray-700 leading-tight">{{ user.username }}</div>
                <div class="text-[10px] text-gray-400 font-bold uppercase tracking-wider">Score: {{ user.totalScore }} • Wins: {{ user.wins }}</div>
            </div>
            <button @click="handleLogout" class="text-xs font-black text-red-400 hover:text-white hover:bg-red-400 border-2 border-red-100 hover:border-red-400 rounded-xl px-4 py-2 transition-all">
                LOGOUT
            </button>
        </div>
    </div>

    <!-- Global Error -->
    <div v-if="errorMsg" class="fixed top-20 right-4 z-50 p-4 bg-red-100 text-duo-red font-bold rounded-xl border-l-4 border-duo-red shadow-md animate-bounce">
       {{ errorMsg }}
    </div>

    <!-- Views -->
    <div class="w-full max-w-6xl flex-1 flex justify-center min-h-0 relative">
        <Transition name="slide-fade" mode="out-in">
            <component :is="viewState === 'AUTH' ? Auth : viewState === 'HOME' ? Home : viewState === 'LOBBY' ? Lobby : viewState === 'GAME' ? Game : GameOver"
                       v-bind="viewState === 'HOME' ? { userName: user?.username } : 
                               viewState === 'LOBBY' ? { roomId, players: sortedPlayers, myId, socket, chatMessages } :
                               viewState === 'GAME' ? { gameState, myId } :
                               viewState === 'ENDED' ? { gameState, myId } : {}"
                       @login="handleLogin" 
                       @guest="handleGuest"
                       @join="connect"
                       @start="startGame"
                       @add-bot="addBot"
                       @quit="quit"
                       @submit="submitWord"
            />
        </Transition>
    </div>
  </div>
</template>