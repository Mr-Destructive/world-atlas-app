<script setup lang="ts">
import { ref, computed, watch } from 'vue'

interface Player {
  id: string
  name: string
  score: number
  lives: number
  isTurn: boolean
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

const props = defineProps<{
  gameState: GameState,
  myId: string
}>()

const emit = defineEmits(['submit', 'quit'])

const inputWord = ref("")
const historyEnd = ref<HTMLElement | null>(null)

const isMyTurn = computed(() => {
  return props.gameState.currentTurn === props.myId
})

const lastLetter = computed(() => {
    if (!props.gameState?.lastWord) return ""
    return props.gameState.lastWord.slice(-1)
})

const sortedPlayers = computed((): Player[] => {
  if (!props.gameState) return []
  return props.gameState.turnOrder
    .map((id: string) => props.gameState!.players[id])
    .filter((p): p is Player => !!p)
})

const submitWord = () => {
  if (!inputWord.value) return
  emit('submit', inputWord.value)
  inputWord.value = ""
}

// Scroll to bottom of history
watch(() => props.gameState.history.length, () => {
  setTimeout(() => {
    historyEnd.value?.scrollIntoView({ behavior: 'smooth' })
  }, 100)
})

const formatTime = (ts: number) => {
    return new Date(ts * 1000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const getPlayerStats = (player: Player) => {
    if (!player.mostUsedPlaces) return []
    return Object.entries(player.mostUsedPlaces)
        .sort(([,a]: any, [,b]: any) => b - a)
        .slice(0, 3)
}
</script>

<template>
    <div class="w-full h-full grid grid-cols-1 lg:grid-cols-3 gap-6 overflow-y-auto lg:overflow-hidden p-1 lg:p-0">
      
      <!-- Left Column: Players & Stats (Order 2 on Mobile) -->
      <div class="lg:col-span-1 space-y-4 flex flex-col min-h-0 order-2 lg:order-1 h-fit lg:h-auto">
        <div class="bg-white rounded-[2rem] shadow-xl p-6 border-b-[6px] border-gray-200 flex flex-col min-h-0 flex-1 max-h-60 lg:max-h-none">
           <div class="flex items-center gap-2 mb-3 px-2 sticky top-0 bg-white z-10">
               <span class="text-xl">👥</span>
               <h3 class="font-black text-gray-400 uppercase text-xs tracking-[0.2em]">Travelers</h3>
           </div>
           <div class="space-y-3 overflow-y-auto pr-1 flex-1 custom-scrollbar">
              <div v-for="player in sortedPlayers" :key="player.id" 
                   class="flex items-center gap-3 p-3 rounded-[1.5rem] transition-all border-2"
                   :class="player.isTurn ? 'bg-green-50 border-duo-green shadow-md translate-x-1' : 'bg-gray-50 border-transparent opacity-70 grayscale-[0.5]'">
                
                <div class="relative">
                    <img :src="`https://api.dicebear.com/7.x/avataaars/svg?seed=${player.name}`" class="w-10 h-10 rounded-2xl bg-white shadow-sm border-2" :class="player.isTurn ? 'border-duo-green' : 'border-gray-100'" />
                    <div v-if="player.isTurn" class="absolute -top-2 -right-2 w-5 h-5 bg-duo-green rounded-full border-2 border-white flex items-center justify-center text-[10px] animate-bounce">
                        ⭐️
                    </div>
                </div>

                <div class="flex-1 min-w-0">
                    <div class="font-black text-gray-700 truncate text-xs">{{ player.name }}</div>
                    <div class="flex gap-0.5 mt-0.5">
                        <span v-for="i in 3" :key="i" class="text-[10px]" :class="i > player.lives ? 'opacity-20' : ''">❤️</span>
                    </div>
                </div>

                <div v-if="player.score" class="text-duo-yellow font-black text-sm">{{ player.score }}</div>
              </div>
           </div>
        </div>

        <!-- Mini Stats Card -->
        <div class="bg-white rounded-[2rem] shadow-xl p-6 border-b-[6px] border-gray-200 h-1/3 shrink-0">
            <div class="flex items-center gap-2 mb-2 px-2">
                <span class="text-lg">🏅</span>
                <h3 class="font-black text-gray-400 uppercase text-[10px] tracking-[0.2em]">Discovery Points</h3>
            </div>
            <div class="space-y-3 overflow-y-auto h-[calc(100%-2rem)] pr-1 custom-scrollbar">
                <div v-for="player in sortedPlayers" :key="'stats-'+player.id" class="mb-3 last:mb-0">
                    <div class="flex items-center gap-2 mb-1 px-2">
                        <div class="text-[8px] font-black text-gray-500 uppercase tracking-widest">{{ player.name }}</div>
                    </div>
                    <div class="flex flex-wrap gap-1 px-1">
                        <span v-for="[place, count] in getPlayerStats(player)" :key="place" class="text-[8px] bg-blue-50 text-duo-blue px-2 py-0.5 rounded-lg font-black uppercase tracking-wider">
                            {{ place }} <span class="opacity-40 ml-1">{{ count }}</span>
                        </span>
                    </div>
                </div>
            </div>
        </div>
      </div>

      <!-- Center & Right: Game & History (Order 1 on Mobile) -->
      <div class="lg:col-span-2 space-y-4 flex flex-col min-h-0 order-1 lg:order-2 h-[600px] lg:h-auto">
          
          <!-- Main Arena -->
          <div class="bg-white rounded-[2.5rem] shadow-xl p-6 border-b-[8px] border-gray-200 relative overflow-hidden transition-all shrink-0"
               :class="{'ring-8 ring-duo-green/20': isMyTurn}">
              <div class="absolute top-4 left-6">
                   <button @click="emit('quit')" class="text-[10px] font-black uppercase text-gray-300 hover:text-red-400 transition-colors">Quit Game</button>
              </div>
              <div class="absolute top-4 right-6">
                  <div class="flex items-center gap-2 bg-duo-yellow text-white px-3 py-1 rounded-2xl shadow-sm">
                      <span class="text-[10px] font-black uppercase tracking-widest">Round {{ gameState.round }}</span>
                  </div>
              </div>

              <div class="text-center mb-4">
                <div class="text-gray-400 text-[9px] uppercase font-black tracking-[0.3em] mb-2">Current Target Letter</div>
                
                <div class="relative inline-flex flex-col items-center group">
                    <div class="text-8xl font-black text-duo-green drop-shadow-[0_4px_0_rgba(22,163,74,0.2)] animate-bounce-slow flex flex-col items-center leading-none">
                        {{ lastLetter || '?' }}
                    </div>
                    
                    <div v-if="gameState.lastWord" class="mt-2 px-3 py-1 bg-gray-50 rounded-xl border-2 border-gray-100">
                        <p class="text-gray-400 text-[9px] font-black uppercase tracking-widest">
                            Continue from: <span class="text-gray-600 font-bold">{{ gameState.lastWord }}</span>
                        </p>
                    </div>
                </div>
              </div>

              <!-- Input Area -->
              <div class="relative group max-w-xl mx-auto">
                   <input 
                     v-model="inputWord"
                     @keyup.enter="submitWord"
                     :disabled="!isMyTurn"
                     type="text" 
                     class="relative w-full p-4 text-xl font-black text-center rounded-[1.5rem] border-[3px] border-gray-200 focus:outline-none focus:border-duo-blue focus:ring-0 disabled:bg-gray-50 disabled:text-gray-300 transition-all placeholder-gray-200 tracking-[0.05em] shadow-inner"
                     :placeholder="isMyTurn ? 'Type a place...' : 'Waiting for turn...'"
                   />
                   <button 
                     v-if="isMyTurn"
                     @click="submitWord"
                     class="absolute right-2 top-2 bottom-2 px-6 bg-duo-green hover:bg-green-500 text-white font-black rounded-[1rem] border-b-[4px] border-green-700 active:border-b-0 active:translate-y-1 transition-all text-base shadow-lg">
                     GO!
                   </button>
              </div>
          </div>

          <!-- Chat History -->
          <div class="bg-white rounded-[2rem] shadow-xl border-b-[6px] border-gray-200 flex flex-col min-h-0 flex-1 overflow-hidden">
              <div class="px-6 py-3 border-b-2 border-gray-100 flex justify-between items-center bg-white z-10 shrink-0">
                  <div class="flex items-center gap-2">
                      <div class="w-2 h-2 bg-duo-green rounded-full animate-pulse"></div>
                      <span class="font-black text-gray-500 uppercase tracking-widest text-[10px]">Adventure Log</span>
                  </div>
                  <span class="bg-gray-100 text-gray-400 px-3 py-0.5 rounded-full text-[9px] font-black">{{ gameState.history.length }} MOVES</span>
              </div>
              
              <div class="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50/30 custom-scrollbar">
                  <div v-if="gameState.history.length === 0" class="h-full flex flex-col items-center justify-center text-gray-300 animate-pulse">
                      <div class="text-4xl mb-2">✨</div>
                      <div class="font-black uppercase tracking-[0.2em] text-[10px]">Waiting for the first spark...</div>
                  </div>

                  <div v-for="(move, index) in gameState.history" :key="index" 
                       class="flex flex-col animate-slide-up"
                       :class="move.playerId === myId ? 'items-end' : 'items-start'">
                      
                      <!-- Name and Time -->
                      <div class="flex items-center gap-2 mb-1 px-2" :class="move.playerId === myId ? 'flex-row-reverse' : ''">
                          <img :src="`https://api.dicebear.com/7.x/avataaars/svg?seed=${move.playerName}`" class="w-5 h-5 rounded-full bg-white border shadow-sm" />
                          <span class="text-[9px] font-black text-gray-400 uppercase tracking-wide">{{ move.playerId === myId ? 'You' : move.playerName }}</span>
                          <span class="text-[8px] text-gray-300 font-bold">{{ formatTime(move.timestamp) }}</span>
                      </div>

                      <!-- Word Bubble -->
                      <div class="relative group max-w-[85%]">
                          <div class="px-5 py-3 rounded-[1.2rem] shadow-sm transition-all group-hover:shadow-md border-b-[3px]"
                               :class="move.playerId === myId ? 'bg-duo-blue text-white border-blue-600 rounded-tr-none' : 'bg-white text-gray-800 border-gray-200 rounded-tl-none'">
                              
                              <div class="flex items-center justify-between gap-3">
                                  <div class="flex flex-col">
                                      <div class="text-lg font-black tracking-tight leading-none">
                                          <span>{{ move.word.slice(0, -1) }}</span><span class="underline decoration-[2px] underline-offset-[2px] decoration-duo-green">{{ move.word.slice(-1) }}</span>
                                      </div>
                                      <div class="text-[8px] font-black uppercase tracking-widest opacity-40 mt-1">
                                          {{ move.type }}
                                      </div>
                                  </div>
                                  
                                  <div class="flex flex-col items-center opacity-30">
                                      <div class="text-[7px] font-black uppercase">Next</div>
                                      <div class="font-black text-sm">{{ move.word.slice(-1) }}</div>
                                  </div>
                              </div>
                          </div>
                          <!-- Success Indicator -->
                          <div class="absolute -top-1 -right-1 w-4 h-4 bg-duo-green rounded-full border-2 border-white flex items-center justify-center text-[8px] shadow-sm transform scale-0 group-hover:scale-100 transition-transform">
                              ✅
                          </div>
                      </div>
                  </div>
                  <div ref="historyEnd"></div>
              </div>
          </div>

      </div>
    </div>
</template>

<style scoped>
.animate-slide-up {
    animation: slideUp 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

@keyframes slideUp {
    from { opacity: 0; transform: translateY(20px) scale(0.95); }
    to { opacity: 1; transform: translateY(0) scale(1); }
}

@keyframes bounceSlow {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-10px); }
}

/* Custom scrollbar for history */
::-webkit-scrollbar {
  width: 8px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 20px;
  border: 2px solid #f8fafc;
}
::-webkit-scrollbar-thumb:hover {
  background: #cbd5e1;
}
</style>
