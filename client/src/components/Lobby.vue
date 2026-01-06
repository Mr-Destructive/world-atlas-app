<template>
  <div class="w-full max-w-4xl grid grid-cols-1 lg:grid-cols-3 gap-6 h-screen lg:h-auto max-h-[calc(100vh-120px)] overflow-y-auto">
      
      <!-- Main Lobby Content (Left/Top) -->
      <div class="lg:col-span-2 bg-white rounded-[2.5rem] shadow-xl p-8 border-b-[8px] border-gray-200 flex flex-col h-fit lg:h-auto">
      
      <!-- Lobby Header -->
      <div class="relative text-center mb-10">
        <button @click="emit('quit')" class="absolute -top-2 -right-2 w-10 h-10 flex items-center justify-center text-gray-300 hover:text-red-500 transition-colors" title="Quit Room">
            <span class="text-2xl font-bold">✕</span>
        </button>
        <div class="inline-block bg-duo-yellow text-white px-4 py-1 rounded-full text-[10px] font-black uppercase tracking-[0.3em] mb-4">
            Waiting Area
        </div>
        <h2 class="text-3xl font-black text-gray-800 mb-6">Gather your team!</h2>
        
        <div class="relative group">
            <div class="text-xs font-black text-gray-400 uppercase tracking-widest mb-3">Secret Room Code</div>
            <div @click="copyCode" class="text-5xl font-black text-duo-blue tracking-[0.2em] select-all cursor-pointer bg-blue-50/50 py-6 rounded-[2rem] border-[3px] border-blue-100 border-dashed inline-block px-12 transition-all hover:scale-105 active:scale-95">
                {{ roomId }}
            </div>
            <div class="mt-4 text-gray-400 font-bold flex items-center justify-center gap-2 text-sm">
                <span>🔗</span> Share this code to invite friends
            </div>
        </div>
      </div>

      <!-- Game Mode Selection (Host Only) -->
      <div v-if="isHost" class="mb-8">
          <h3 class="font-black text-gray-400 uppercase text-xs tracking-[0.2em] mb-4 px-4">Game Mode</h3>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <button @click="selectedMode = 'CLASSIC'" 
                      class="p-4 rounded-2xl border-2 text-left transition-all"
                      :class="selectedMode === 'CLASSIC' ? 'bg-green-50 border-duo-green' : 'bg-gray-50 border-transparent hover:border-gray-200'">
                  <div class="text-xs font-black uppercase tracking-wider mb-1" :class="selectedMode === 'CLASSIC' ? 'text-duo-green' : 'text-gray-500'">Classic</div>
                  <div class="text-[10px] text-gray-400 font-bold leading-tight">Standard rules. 3 lives.</div>
              </button>
              
              <button @click="selectedMode = 'POINT_RUSH'" 
                      class="p-4 rounded-2xl border-2 text-left transition-all"
                      :class="selectedMode === 'POINT_RUSH' ? 'bg-blue-50 border-duo-blue' : 'bg-gray-50 border-transparent hover:border-gray-200'">
                  <div class="text-xs font-black uppercase tracking-wider mb-1" :class="selectedMode === 'POINT_RUSH' ? 'text-duo-blue' : 'text-gray-500'">Point Rush</div>
                  <div class="text-[10px] text-gray-400 font-bold leading-tight">Speed & length boost scores!</div>
              </button>

              <button @click="selectedMode = 'SUDDEN_DEATH'" 
                      class="p-4 rounded-2xl border-2 text-left transition-all"
                      :class="selectedMode === 'SUDDEN_DEATH' ? 'bg-red-50 border-duo-red' : 'bg-gray-50 border-transparent hover:border-gray-200'">
                  <div class="text-xs font-black uppercase tracking-wider mb-1" :class="selectedMode === 'SUDDEN_DEATH' ? 'text-duo-red' : 'text-gray-500'">Sudden Death</div>
                  <div class="text-[10px] text-gray-400 font-bold leading-tight">1 Life. No mistakes allowed.</div>
              </button>
          </div>
      </div>

      <!-- Players List -->
      <div class="mb-10">
          <div class="flex items-center justify-between mb-6 px-4">
              <h3 class="font-black text-gray-400 uppercase text-xs tracking-[0.2em]">Travelers ({{ players.length }})</h3>
              <div class="h-1 flex-1 mx-4 bg-gray-100 rounded-full"></div>
              <div class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
                  <span v-if="players.filter(p => p.type === 0).length > 0">
                      {{ readyPlayers.size }}/{{ players.filter(p => p.type === 0).length }} ready
                  </span>
              </div>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div v-for="player in players" :key="player.id" 
                   @click="player.type === 0 && player.id === myId && toggleReady(player.id)"
                   class="flex items-center gap-4 p-4 rounded-[1.5rem] border-[3px] transition-all group"
                   :class="[player.type === 0 ? 'cursor-pointer' : '', isPlayerReady(player.id) ? 'bg-green-50 border-duo-green shadow-md' : 'bg-gray-50 border-transparent hover:border-gray-200']">
                  <div class="relative">
                      <img :src="`https://api.dicebear.com/7.x/avataaars/svg?seed=${player.name}`" class="w-12 h-12 rounded-2xl bg-white shadow-sm border-2 transition-all"
                           :class="isPlayerReady(player.id) ? 'border-duo-green' : 'border-gray-100 group-hover:rotate-6'" />
                      <div class="absolute -top-1 -right-1 w-4 h-4 bg-duo-green rounded-full border-2 border-white" v-if="player.id === myId"></div>
                  </div>
                  <div class="flex-1 min-w-0">
                      <div class="font-black text-gray-700 truncate">{{ player.name }}</div>
                      <div class="flex gap-1">
                          <span v-if="player.id === myId" class="text-[8px] font-black text-duo-green bg-green-100 px-2 py-0.5 rounded-full uppercase">You</span>
                          <span v-if="player.type === 1" class="text-[8px] font-black text-duo-blue bg-blue-100 px-2 py-0.5 rounded-full uppercase">Bot</span>
                          <span v-if="player.type === 0 && isPlayerReady(player.id)" class="text-[8px] font-black text-duo-green bg-green-100 px-2 py-0.5 rounded-full uppercase">Ready</span>
                      </div>
                  </div>
                  <div v-if="isPlayerReady(player.id)" class="text-lg">✓</div>
              </div>
          </div>
      </div>

      <!-- Controls -->
      <div class="pt-2">
          <div v-if="isHost" class="flex flex-col sm:flex-row gap-4">
               <button @click="emit('addBot')" class="flex-1 py-4 bg-white hover:bg-gray-50 text-duo-blue font-black rounded-[1.5rem] border-[3px] border-b-[6px] border-gray-200 hover:border-blue-200 active:border-b-[3px] active:translate-y-[3px] transition-all uppercase tracking-widest text-sm">
                   🤖 Add Bot
               </button>
               <button @click="startGame" :disabled="!allPlayersReady && players.filter(p => p.type === 0).length > 0" class="flex-[1.5] py-4 bg-duo-green hover:bg-green-500 disabled:bg-gray-300 disabled:border-gray-400 disabled:text-gray-500 disabled:cursor-not-allowed text-white font-black rounded-[2rem] border-b-[8px] border-green-700 disabled:border-gray-400 active:border-b-0 active:translate-y-2 transition-all text-xl uppercase tracking-[0.2em] shadow-lg">
                   {{ allPlayersReady ? 'Start Game!' : 'Waiting for Players...' }}
               </button>
          </div>
          <div v-else class="flex flex-col gap-4">
              <button @click="toggleReady(myId)" class="w-full py-4 transition-all font-black rounded-[1.5rem] border-[3px] border-b-[6px] uppercase tracking-widest text-sm"
                      :class="isPlayerReady(myId) ? 'bg-duo-green hover:bg-green-500 text-white border-green-700 active:border-b-0 active:translate-y-2' : 'bg-white hover:bg-gray-50 text-duo-blue border-gray-200 hover:border-blue-200 active:border-b-[3px] active:translate-y-[3px]'">
                  {{ isPlayerReady(myId) ? '✓ Ready!' : 'Mark as Ready' }}
              </button>
              <div class="text-center p-6 bg-gray-50 rounded-[2rem] border-2 border-dashed border-gray-200">
                  <div class="inline-block animate-bounce mb-2">🎈</div>
                  <div class="text-gray-400 font-black uppercase tracking-widest text-xs">Waiting for the host to start...</div>
                  <div v-if="players.filter(p => p.type === 0).length > 1" class="text-gray-300 font-bold text-[10px] mt-2">Everyone must be ready</div>
              </div>
          </div>
      </div>

      </div>

      <!-- Chat Section (Right/Bottom) -->
      <div class="bg-white rounded-[2.5rem] shadow-xl border-b-[8px] border-gray-200 flex flex-col h-fit lg:h-auto lg:col-span-1">
        <div class="px-6 py-4 border-b-2 border-gray-100 flex items-center gap-2 shrink-0">
            <span class="text-xl">💬</span>
            <h3 class="font-black text-gray-400 uppercase text-xs tracking-[0.2em]">Lobby Chat</h3>
        </div>

        <!-- Messages -->
        <div class="flex-1 overflow-y-auto p-4 space-y-3 min-h-0 max-h-64 lg:max-h-80 custom-scrollbar">
            <div v-if="chatMessages.length === 0" class="h-full flex items-center justify-center">
                <div class="text-center">
                    <div class="text-2xl mb-2 opacity-30">🤐</div>
                    <div class="text-gray-300 text-[10px] font-black uppercase tracking-widest">No messages yet</div>
                </div>
            </div>
            <div v-for="(msg, idx) in chatMessages" :key="idx" class="animate-fade-in">
                <div class="flex items-start gap-2">
                    <img :src="`https://api.dicebear.com/7.x/avataaars/svg?seed=${msg.playerName}`" class="w-7 h-7 rounded-lg bg-white border shadow-sm shrink-0" />
                    <div class="flex-1 min-w-0">
                        <div class="flex items-baseline gap-2">
                            <span class="font-black text-gray-700 text-sm truncate">{{ msg.playerName }}</span>
                            <span class="text-[9px] text-gray-300 font-bold">{{ formatChatTime(msg.timestamp) }}</span>
                        </div>
                        <p class="text-gray-600 text-sm font-medium break-words mt-0.5">{{ msg.message }}</p>
                    </div>
                </div>
            </div>
            <div ref="chatEnd"></div>
        </div>

        <!-- Input -->
        <div class="px-4 py-3 border-t-2 border-gray-100 shrink-0">
            <div class="flex gap-2">
                <input 
                    v-model="chatInput" 
                    @keyup.enter="sendChat"
                    type="text" 
                    class="flex-1 px-3 py-2 bg-gray-50 rounded-xl border-2 border-gray-200 focus:border-duo-blue focus:outline-none font-bold text-sm placeholder-gray-300 transition-colors"
                    placeholder="Say hi..." />
                <button 
                    @click="sendChat" 
                    :disabled="!chatInput.trim()"
                    class="px-4 py-2 bg-duo-blue hover:bg-blue-500 disabled:bg-gray-200 disabled:text-gray-400 text-white font-black rounded-xl border-b-[3px] border-blue-600 disabled:border-gray-300 active:border-b-0 active:translate-y-0.5 transition-all text-sm">
                    Send
                </button>
            </div>
        </div>
      </div>

      </div>
      </template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

interface ChatMessage {
  playerId: string
  playerName: string
  message: string
  timestamp: number
}

const props = defineProps<{
  roomId: string,
  players: any[],
  myId: string,
  socket?: WebSocket,
  chatMessages?: any[]
}>()

const emit = defineEmits(['start', 'addBot', 'quit'])

const selectedMode = ref("CLASSIC")
const readyPlayers = ref<Set<string>>(new Set())
const chatInput = ref("")
const chatEnd = ref<HTMLElement | null>(null)

// Use prop's chatMessages if provided, otherwise use local
const localChatMessages = ref<ChatMessage[]>([])
const chatMessages = computed(() => props.chatMessages || localChatMessages.value)

const isHost = computed(() => {
    return props.players.length > 0 && props.players[0].id === props.myId
})

const isPlayerReady = (playerId: string) => {
    return readyPlayers.value.has(playerId)
}

const toggleReady = (playerId: string) => {
    if (readyPlayers.value.has(playerId)) {
        readyPlayers.value.delete(playerId)
    } else {
        readyPlayers.value.add(playerId)
    }
}

const allPlayersReady = computed(() => {
    if (props.players.length === 0) return false
    return props.players.every(p => isPlayerReady(p.id))
})

const startGame = () => {
    emit('start', { mode: selectedMode.value })
}

const copyCode = () => {
    navigator.clipboard.writeText(props.roomId)
}

const sendChat = () => {
    if (!chatInput.value.trim() || !props.socket) return
    
    props.socket.send(JSON.stringify({
        type: 'CHAT',
        payload: { message: chatInput.value }
    }))
    
    chatInput.value = ""
}

const formatChatTime = (ts: number) => {
    return new Date(ts * 1000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

// Auto-scroll when chat messages change
watch(() => chatMessages.value?.length, () => {
    setTimeout(() => {
        chatEnd.value?.scrollIntoView({ behavior: 'smooth' })
    }, 50)
})
</script>