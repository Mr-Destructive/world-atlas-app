<template>
  <div class="w-full max-w-2xl bg-white rounded-[2.5rem] shadow-xl p-8 border-b-[8px] border-gray-200">
      
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
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div v-for="player in players" :key="player.id" 
                   class="flex items-center gap-4 p-4 bg-gray-50 rounded-[1.5rem] border-[3px] border-transparent transition-all hover:border-gray-200 group">
                  <div class="relative">
                      <img :src="`https://api.dicebear.com/7.x/avataaars/svg?seed=${player.name}`" class="w-12 h-12 rounded-2xl bg-white shadow-sm border-2 border-gray-100 group-hover:rotate-6 transition-transform" alt="avatar" />
                      <div class="absolute -top-1 -right-1 w-4 h-4 bg-duo-green rounded-full border-2 border-white" v-if="player.id === myId"></div>
                  </div>
                  <div class="flex-1 min-w-0">
                      <div class="font-black text-gray-700 truncate">{{ player.name }}</div>
                      <div class="flex gap-1">
                          <span v-if="player.id === myId" class="text-[8px] font-black text-duo-green bg-green-100 px-2 py-0.5 rounded-full uppercase">You</span>
                          <span v-if="player.type === 1" class="text-[8px] font-black text-duo-blue bg-blue-100 px-2 py-0.5 rounded-full uppercase">Bot</span>
                      </div>
                  </div>
              </div>
          </div>
      </div>

      <!-- Controls -->
      <div class="pt-2">
          <div v-if="isHost" class="flex flex-col sm:flex-row gap-4">
               <button @click="emit('addBot')" class="flex-1 py-4 bg-white hover:bg-gray-50 text-duo-blue font-black rounded-[1.5rem] border-[3px] border-b-[6px] border-gray-200 hover:border-blue-200 active:border-b-[3px] active:translate-y-[3px] transition-all uppercase tracking-widest text-sm">
                   🤖 Add Bot
               </button>
               <button @click="startGame" class="flex-[1.5] py-4 bg-duo-green hover:bg-green-500 text-white font-black rounded-[2rem] border-b-[8px] border-green-700 active:border-b-0 active:translate-y-2 transition-all text-xl uppercase tracking-[0.2em] shadow-lg">
                   Start Game!
               </button>
          </div>
          <div v-else class="text-center p-8 bg-gray-50 rounded-[2rem] border-2 border-dashed border-gray-200">
              <div class="inline-block animate-bounce mb-2">🎈</div>
              <div class="text-gray-400 font-black uppercase tracking-widest text-xs">Waiting for the host to start the fun...</div>
          </div>
      </div>

    </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = defineProps<{
  roomId: string,
  players: any[],
  myId: string
}>()

const emit = defineEmits(['start', 'addBot', 'quit'])

const selectedMode = ref("CLASSIC")

const isHost = computed(() => {
    return props.players.length > 0 && props.players[0].id === props.myId
})

const startGame = () => {
    emit('start', { mode: selectedMode.value })
}

const copyCode = () => {
    navigator.clipboard.writeText(props.roomId)
}
</script>