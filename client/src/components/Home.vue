<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'

const props = defineProps<{
    userName?: string
}>()

const emit = defineEmits(['join'])

const name = ref("")
const roomId = ref("")
const isCreating = ref(true)
const isAnonymous = ref(false)

onMounted(() => {
    if (props.userName) {
        name.value = props.userName
    }
})

// Keep name in sync if prop changes
watch(() => props.userName, (newVal) => {
    if (newVal) {
        name.value = newVal
        isAnonymous.value = false
    }
})

const join = () => {
  const finalName = isAnonymous.value ? `Explorer-${Math.floor(Math.random() * 9999)}` : name.value
  if (!finalName && !isAnonymous.value) return
  emit('join', { name: finalName, roomId: isCreating.value ? "" : roomId.value })
}
</script>

<template>
  <div class="w-full max-w-md">
    <!-- Character/Mascot Placeholder Area -->
    <div class="text-center mb-8 animate-bounce-slow">
        <div class="inline-block p-6 bg-white rounded-full shadow-xl border-b-4 border-gray-200 text-6xl">
            🌍
        </div>
    </div>

    <div class="bg-white rounded-[2rem] shadow-xl p-8 border-b-[6px] border-gray-200">
        <div class="text-center mb-8">
            <h2 class="text-3xl font-black text-gray-800 leading-tight">Ready for an adventure?</h2>
            <p class="text-gray-500 font-bold mt-2">The world is waiting for your words!</p>
        </div>

        <div class="space-y-6">
            <!-- Anonymous Toggle -->
            <div v-if="!userName" @click="isAnonymous = !isAnonymous" 
                 class="flex items-center gap-4 p-4 rounded-2xl border-2 cursor-pointer transition-all"
                 :class="isAnonymous ? 'border-duo-green bg-green-50' : 'border-gray-100 bg-gray-50 hover:border-gray-200'">
                <div class="w-12 h-12 rounded-full flex items-center justify-center text-2xl bg-white shadow-sm">
                    {{ isAnonymous ? '👤' : '🆔' }}
                </div>
                <div class="flex-1">
                    <div class="font-black text-gray-700">Play Anonymously</div>
                    <div class="text-xs font-bold text-gray-400">Jump right in without a name</div>
                </div>
                <div class="w-6 h-6 rounded-full border-2 flex items-center justify-center"
                     :class="isAnonymous ? 'bg-duo-green border-duo-green text-white' : 'border-gray-300'">
                    <span v-if="isAnonymous">✓</span>
                </div>
            </div>

            <!-- Name Input (conditional) -->
            <div v-if="!isAnonymous" class="space-y-2 animate-fade-in">
                <label class="block text-sm font-black text-gray-400 uppercase tracking-widest px-2">Your Explorer Name</label>
                <input v-model="name" type="text" 
                       :disabled="!!userName"
                       class="w-full p-4 rounded-2xl border-2 border-gray-200 focus:border-duo-blue focus:ring-8 focus:ring-blue-100 outline-none font-black text-lg text-gray-700 transition-all placeholder-gray-300 disabled:bg-gray-50 disabled:text-gray-500" 
                       placeholder="e.g. MarcoPolo" />
            </div>

            <!-- Mode Selector -->
            <div class="flex gap-2 p-1.5 bg-gray-100 rounded-2xl">
                <button @click="isCreating = true" 
                        :class="{'bg-white shadow-md text-duo-green scale-[1.02]': isCreating, 'text-gray-500 hover:text-gray-700': !isCreating}" 
                        class="flex-1 py-3 font-black rounded-xl transition-all uppercase text-sm tracking-wider">
                    New Room
                </button>
                <button @click="isCreating = false" 
                        :class="{'bg-white shadow-md text-duo-blue scale-[1.02]': !isCreating, 'text-gray-500 hover:text-gray-700': isCreating}" 
                        class="flex-1 py-3 font-black rounded-xl transition-all uppercase text-sm tracking-wider">
                    Join Room
                </button>
            </div>

            <!-- Room ID Input -->
            <div v-if="!isCreating" class="space-y-2 animate-fade-in">
                <label class="block text-sm font-black text-gray-400 uppercase tracking-widest px-2">Secret Code</label>
                <input v-model="roomId" type="text" 
                       class="w-full p-4 rounded-2xl border-2 border-gray-200 focus:border-duo-blue focus:ring-8 focus:ring-blue-100 outline-none font-black text-lg text-gray-700 transition-all placeholder-gray-300 uppercase" 
                       placeholder="Enter Code" />
            </div>

            <button @click="join" 
                    :disabled="!isAnonymous && !name"
                    class="w-full py-5 bg-duo-green hover:bg-green-500 disabled:bg-gray-200 disabled:border-gray-300 disabled:text-gray-400 text-white font-black rounded-[1.5rem] border-b-[6px] border-green-700 active:border-b-0 active:translate-y-1 transition-all text-xl uppercase tracking-[0.2em] shadow-lg mt-4">
                LET'S GO!
            </button>
        </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in {
    animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
}

.animate-bounce-slow {
    animation: bounceSlow 3s infinite ease-in-out;
}

@keyframes bounceSlow {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-15px); }
}
</style>
