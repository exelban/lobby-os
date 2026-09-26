<template>
  <div class="contents" @click="emit('select')">
    <button @click="emit('search')" title="Search links (/)" aria-label="Search links" class="toolbar-button toolbar-neutral">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="10.5" cy="10.5" r="6.5"/><path stroke-linecap="round" d="m16 16 4.5 4.5"/></svg>
    </button>
    <button @click="(store.state.window.settings = true)" title="Settings" class="toolbar-button toolbar-neutral">
      <svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0,0h24v24H0V0z" fill="none"/><path fill="currentColor" d="M19.14,12.94c0.04-0.3,0.06-0.61,0.06-0.94c0-0.32-0.02-0.64-0.07-0.94l2.03-1.58c0.18-0.14,0.23-0.41,0.12-0.61 l-1.92-3.32c-0.12-0.22-0.37-0.29-0.59-0.22l-2.39,0.96c-0.5-0.38-1.03-0.7-1.62-0.94L14.4,2.81c-0.04-0.24-0.24-0.41-0.48-0.41 h-3.84c-0.24,0-0.43,0.17-0.47,0.41L9.25,5.35C8.66,5.59,8.12,5.92,7.63,6.29L5.24,5.33c-0.22-0.08-0.47,0-0.59,0.22L2.74,8.87 C2.62,9.08,2.66,9.34,2.86,9.48l2.03,1.58C4.84,11.36,4.8,11.69,4.8,12s0.02,0.64,0.07,0.94l-2.03,1.58 c-0.18,0.14-0.23,0.41-0.12,0.61l1.92,3.32c0.12,0.22,0.37,0.29,0.59,0.22l2.39-0.96c0.5,0.38,1.03,0.7,1.62,0.94l0.36,2.54 c0.05,0.24,0.24,0.41,0.48,0.41h3.84c0.24,0,0.44-0.17,0.47-0.41l0.36-2.54c0.59-0.24,1.13-0.56,1.62-0.94l2.39,0.96 c0.22,0.08,0.47,0,0.59-0.22l1.92-3.32c0.12-0.22,0.07-0.47-0.12-0.61L19.14,12.94z M12,15.6c-1.98,0-3.6-1.62-3.6-3.6 s1.62-3.6,3.6-3.6s3.6,1.62,3.6,3.6S13.98,15.6,12,15.6z"/></svg>
    </button>
    <button @click="store.state.editMode ? store.stopEditMode() : store.startEditMode()" title="Edit mode" class="toolbar-button" :class="store.state.editMode ? 'text-blue-500 bg-blue-50 dark:bg-blue-500/15' : 'text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100/80 dark:hover:text-neutral-200 dark:hover:bg-white/10'">
      <svg v-if="!store.state.editMode" xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/></svg>
      <svg v-else xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M9 16.2L4.8 12l-1.4 1.4L9 19 21 7l-1.4-1.4L9 16.2z"/></svg>
    </button>
    <div class="w-px h-5 bg-neutral-200 dark:bg-neutral-700 mx-0.5"></div>
    <button @click="store.openNewLink()" title="Add link" class="toolbar-button toolbar-neutral">
      <svg xmlns="http://www.w3.org/2000/svg" height="22" viewBox="0 0 24 24" width="22"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/></svg>
    </button>
  </div>
</template>

<script setup>
import store from "../utils/store.js"
const emit = defineEmits(["search", "select"])
</script>

<style scoped>
@reference "../style.css";
.toolbar-button {
  @apply w-10 h-10 flex items-center justify-center rounded-full transition-all duration-200;
}
.toolbar-neutral {
  @apply text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100/80 dark:hover:text-neutral-200 dark:hover:bg-white/10;
}
</style>
