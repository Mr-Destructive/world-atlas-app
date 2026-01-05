<script setup lang="ts">
import { computed, onMounted } from 'vue'
import confetti from 'canvas-confetti'

interface Player {
  id: string
  name: string
  score: number
  lives: number
  type: number
}

interface GameState {
  players: Record<string, Player>
  state: string
  lastWord: string
  turnOrder: string[]
  currentTurn: string
  round: number
}

const props = defineProps<{
  gameState: GameState,
  myId: string
}>()

const emit = defineEmits(['start', 'quit'])

const sortedByScore = computed((): Player[] => {
    return Object.values(props.gameState.players)
        .sort((a, b) => b.score - a.score)
})

const winner = computed((): Player => {
    const players = Object.values(props.gameState.players)
    if (players.length === 0) return { id: '', name: 'None', score: 0, lives: 0, type: 0 }
    
    const alive = players.filter((p) => p.lives > 0)
    if (alive.length === 1) return alive[0] as Player
    
    return (sortedByScore.value[0] || players[0]) as Player
})

const isWinner = computed(() => winner.value.id === props.myId)

onMounted(() => {
    if (isWinner.value) {
        const duration = 5 * 1000
        const animationEnd = Date.now() + duration
        const defaults = { startVelocity: 30, spread: 360, ticks: 60, zIndex: 0 }

        const randomInRange = (min: number, max: number) => Math.random() * (max - min) + min

        const interval: any = setInterval(function() {
            const timeLeft = animationEnd - Date.now()

            if (timeLeft <= 0) {
                return clearInterval(interval)
            }

            const particleCount = 50 * (timeLeft / duration)
            confetti({ ...defaults, particleCount, origin: { x: randomInRange(0.1, 0.3), y: Math.random() - 0.2 } })
            confetti({ ...defaults, particleCount, origin: { x: randomInRange(0.7, 0.9), y: Math.random() - 0.2 } })
        }, 250)
    }
})
</script>

<template>
  <div class="w-full max-w-2xl bg-white rounded-[3rem] shadow-xl p-10 border-b-[10px] border-gray-200 text-center animate-pop-in">
    
    <div class="mb-10">
        <div class="inline-block relative mb-6">
            <div class="text-8xl animate-bounce-slow">🏆</div>
            <div class="absolute -top-2 -right-2 text-4xl animate-pulse">✨</div>
            <div class="absolute -bottom-2 -left-2 text-4xl animate-pulse delay-700">✨</div>
        </div>
        <h2 class="text-4xl font-black text-gray-800 uppercase tracking-tight">Adventure Complete!</h2>
        <div class="flex items-center justify-center gap-2 mt-4">
            <div class="h-1 w-8 bg-gray-100 rounded-full"></div>
            <p class="text-gray-400 font-black uppercase tracking-[0.3em] text-xs">Hall of Fame</p>
            <div class="h-1 w-8 bg-gray-100 rounded-full"></div>
        </div>
    </div>

    <!-- Winner Card -->
    <div class="bg-gradient-to-b from-duo-yellow/20 to-duo-yellow/5 border-[4px] border-duo-yellow rounded-[2.5rem] p-8 mb-10 transform hover:scale-[1.02] transition-all shadow-lg relative overflow-hidden group">
        <div class="absolute top-0 left-0 w-full h-1 bg-duo-yellow/30"></div>
        <div class="text-duo-yellow font-black uppercase tracking-[0.2em] text-xs mb-4">Ultimate Explorer</div>
        <div class="flex flex-col items-center gap-4">
             <div class="relative">
                 <img :src="`https://api.dicebear.com/7.x/avataaars/svg?seed=${winner.name}`" class="w-24 h-24 rounded-[2rem] bg-white border-4 border-duo-yellow shadow-md group-hover:rotate-3 transition-transform" />
                 <div class="absolute -bottom-2 -right-2 bg-duo-yellow text-white w-8 h-8 rounded-full border-4 border-white flex items-center justify-center font-black">1</div>
             </div>
             <div class="text-center">
                 <div class="text-3xl font-black text-gray-800 tracking-tight">{{ winner.name }}</div>
                 <div class="text-duo-yellow font-black text-xl mt-1">{{ winner.score }} Discovery Points</div>
             </div>
        </div>
    </div>

    <!-- Scoreboard -->
    <div class="space-y-4 mb-10">
        <div v-for="(player, index) in sortedByScore" :key="player.id" 
             class="flex items-center justify-between p-5 rounded-[2rem] bg-gray-50 border-b-[4px] border-gray-200 transition-all hover:translate-x-1"
             :class="{'opacity-60': index > 0}">
            <div class="flex items-center gap-5">
                <span class="font-black text-gray-300 w-6 text-xl">#{{ index + 1 }}</span>
                <img :src="`https://api.dicebear.com/7.x/avataaars/svg?seed=${player.name}`" class="w-10 h-10 rounded-xl bg-white border-2 border-gray-100" />
                <span class="font-black text-gray-700 text-lg">{{ player.name }}</span>
            </div>
            <div class="flex flex-col items-end">
                <div class="font-black text-gray-800 text-xl">{{ player.score }}</div>
                <div class="text-[9px] font-black text-gray-400 uppercase tracking-widest">Points</div>
            </div>
        </div>
    </div>

    <div class="flex flex-col sm:flex-row gap-4">
        <button @click="emit('quit')" class="flex-1 py-6 bg-white hover:bg-gray-50 text-duo-blue font-black rounded-[2rem] border-[3px] border-b-[8px] border-gray-200 hover:border-blue-200 active:border-b-[3px] active:translate-y-[3px] transition-all text-xl uppercase tracking-widest">
            Home
        </button>
        <button @click="emit('start')" class="flex-[2] py-6 bg-duo-green hover:bg-green-500 text-white font-black rounded-[2rem] border-b-[8px] border-green-700 active:border-b-0 active:translate-y-2 transition-all text-2xl uppercase tracking-[0.2em] shadow-xl">
            Explore Again!
        </button>
    </div>

  </div>
</template>

<style scoped>
.animate-pop-in {
    animation: popIn 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

@keyframes popIn {
    from { opacity: 0; transform: scale(0.9); }
    to { opacity: 1; transform: scale(1); }
}

.animate-bounce-slow {
    animation: bounceSlow 3s infinite ease-in-out;
}

@keyframes bounceSlow {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-20px); }
}
</style>
