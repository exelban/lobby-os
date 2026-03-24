<template>
  <div class="fixed z-30 top-0 left-0 w-full h-full p-3 sm:p-0 flex justify-center items-center bg-black/30 dark:bg-black/50 backdrop-blur-sm">
    <div @click.stop class="relative w-full max-w-sm max-h-full">
      <div class="relative bg-white/85 dark:bg-neutral-900/90 backdrop-blur-xl rounded-2xl shadow-2xl shadow-black/10 dark:shadow-black/40 border border-white/60 dark:border-white/10 overflow-hidden">
        <div class="p-6 text-center">
          <div class="w-12 h-12 mx-auto mb-4 flex items-center justify-center rounded-full bg-red-50 dark:bg-red-500/10">
            <svg class="w-6 h-6 text-red-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"/>
            </svg>
          </div>
          <h3 class="mb-2 text-base font-semibold text-neutral-800 dark:text-neutral-200">Delete this link?</h3>
          <p class="mb-5 text-sm text-neutral-500 dark:text-neutral-400">This action cannot be undone.</p>
          <div class="flex gap-3">
            <button @click="store.dispatch('closeDelete')" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 dark:text-neutral-300 bg-white/60 dark:bg-white/5 border border-neutral-200/80 dark:border-white/10 rounded-xl hover:bg-neutral-100/80 dark:hover:bg-white/10 transition-all duration-200">Cancel</button>
            <button @click="store.dispatch('delete')" class="flex-1 py-2.5 text-sm font-medium text-white bg-red-500 hover:bg-red-600 rounded-xl transition-all duration-200 shadow-sm shadow-red-500/20">Delete</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useStore } from "vuex"
const store = useStore()
</script>
