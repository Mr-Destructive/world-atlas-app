<template>
  <div class="w-full max-w-sm bg-white rounded-[2rem] shadow-xl p-6 border-b-[8px] border-gray-200">
    <div class="flex mb-5 bg-gray-100 p-1 rounded-xl">
      <button 
        @click="isLogin = true"
        class="flex-1 py-1.5 rounded-lg font-black text-xs uppercase tracking-wider transition-all"
        :class="isLogin ? 'bg-white text-duo-green shadow-sm' : 'text-gray-400 hover:text-gray-600'"
      >
        Login
      </button>
      <button 
        @click="isLogin = false"
        class="flex-1 py-1.5 rounded-lg font-black text-xs uppercase tracking-wider transition-all"
        :class="!isLogin ? 'bg-white text-duo-blue shadow-sm' : 'text-gray-400 hover:text-gray-600'"
      >
        Signup
      </button>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-2.5">
      <div>
        <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Username</label>
        <input 
          v-model="username" 
          type="text" 
          required
          class="w-full px-3 py-2 bg-gray-50 rounded-lg font-bold border-2 border-gray-100 focus:border-duo-green focus:outline-none transition-colors text-sm"
          placeholder="Explorer Name"
        />
      </div>
      
      <div>
        <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Password</label>
        <input 
          v-model="password" 
          type="password" 
          required
          class="w-full px-3 py-2 bg-gray-50 rounded-lg font-bold border-2 border-gray-100 focus:border-duo-green focus:outline-none transition-colors text-sm"
          placeholder="••••••••"
        />
      </div>

      <div v-if="error" class="text-red-500 text-xs font-black text-center bg-red-50 p-2 rounded-lg">
        {{ error }}
      </div>

      <button 
        type="submit" 
        :disabled="loading"
        class="w-full py-2.5 mt-1 bg-duo-green hover:bg-green-500 text-white font-black rounded-[1.2rem] border-b-[5px] border-green-700 active:border-b-0 active:translate-y-1 transition-all uppercase tracking-wider shadow-lg disabled:opacity-50 text-xs"
      >
        {{ loading ? 'Wait...' : (isLogin ? 'Log In' : 'Sign Up') }}
      </button>
      
      <div class="text-center pt-1">
          <button type="button" @click="$emit('guest')" class="text-gray-400 text-xs font-black uppercase hover:text-gray-600 underline">
              Continue as Guest
          </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits(['login', 'guest'])

const isLogin = ref(true)
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleSubmit = async () => {
  error.value = ''
  loading.value = true
  
  const endpoint = isLogin.value ? '/api/login' : '/api/register'
  
  try {
    // Get API URL from environment variable or default to current host for dev
    const apiUrl = import.meta.env.VITE_API_URL || window.location.origin
    const res = await fetch(`${apiUrl}${endpoint}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: username.value, password: password.value })
    })

    const data = await res.json()
    
    if (!res.ok) {
      throw new Error(data.error || 'Something went wrong')
    }

    emit('login', data)
  } catch (err: any) {
    error.value = typeof err.message === 'string' ? err.message : 'Failed to connect'
  } finally {
    loading.value = false
  }
}
</script>
