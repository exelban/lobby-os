<template>
  <main class="min-h-dvh p-1 sm:p-2 flex flex-row justify-center items-center content-center">
    <div v-if="store.state.links" class="container flex flex-col items-center gap-4 sm:gap-6">
      <div v-for="g in Object.keys(store.state.links)" :key="g">
        <draggable v-model="store.state.links[g]" @change="() => store.dispatch('save')" tag="ul" group="links" :disabled="!store.state.editMode" item-key="id" ghost-class="opacity-30" class="flex flex-wrap items-center justify-center flew-row gap-x-4 sm:gap-x-5 gap-y-2 sm:gap-y-3 p-1 sm:p-2 rounded-lg" :class="store.state.editMode ? 'bg-neutral-500/10 border border-dashed border-neutral-400/40 dark:border-neutral-600/40' : 'border border-transparent'">
          <template #item="{element}">
            <li><CLink @click.stop :value="element"></CLink></li>
          </template>
        </draggable>
      </div>

      <div class="h-10">
        <button v-if="store.state.editMode" @click.stop @click="store.commit('addGroup')" class="w-9 h-9 flex items-center justify-center bg-neutral-200/60 hover:bg-neutral-300/60 dark:bg-neutral-700/50 dark:hover:bg-neutral-600/50 text-neutral-500 dark:text-neutral-400 rounded-full backdrop-blur-sm transition-all duration-200 text-lg">+</button>
      </div>
    </div>
    <button v-else @click="store.dispatch('openNewLink')" class="shadow-sm text-neutral-600 bg-white/70 backdrop-blur-sm border border-white/80 hover:bg-white/90 font-medium rounded-2xl text-sm px-8 py-4 dark:bg-white/5 dark:text-neutral-300 dark:border-white/10 dark:hover:bg-white/10 transition-all duration-300">Create first link</button>
  </main>

  <!-- Floating toolbar - always visible on desktop, toggle button on mobile -->
  <div class="fixed bottom-5 left-1/2 -translate-x-1/2 sm:flex hidden items-center gap-1 px-2 py-2 rounded-full bg-white/70 dark:bg-neutral-900/70 backdrop-blur-xl border border-white/60 dark:border-white/10 shadow-lg shadow-black/5 dark:shadow-black/30 transition-all duration-300 z-20">
    <button @click="store.commit('setSettingsWindow', true)" title="Settings" class="w-10 h-10 flex items-center justify-center rounded-full text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100/80 dark:hover:text-neutral-200 dark:hover:bg-white/10 transition-all duration-200">
      <svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0,0h24v24H0V0z" fill="none"/><path fill="currentColor" d="M19.14,12.94c0.04-0.3,0.06-0.61,0.06-0.94c0-0.32-0.02-0.64-0.07-0.94l2.03-1.58c0.18-0.14,0.23-0.41,0.12-0.61 l-1.92-3.32c-0.12-0.22-0.37-0.29-0.59-0.22l-2.39,0.96c-0.5-0.38-1.03-0.7-1.62-0.94L14.4,2.81c-0.04-0.24-0.24-0.41-0.48-0.41 h-3.84c-0.24,0-0.43,0.17-0.47,0.41L9.25,5.35C8.66,5.59,8.12,5.92,7.63,6.29L5.24,5.33c-0.22-0.08-0.47,0-0.59,0.22L2.74,8.87 C2.62,9.08,2.66,9.34,2.86,9.48l2.03,1.58C4.84,11.36,4.8,11.69,4.8,12s0.02,0.64,0.07,0.94l-2.03,1.58 c-0.18,0.14-0.23,0.41-0.12,0.61l1.92,3.32c0.12,0.22,0.37,0.29,0.59,0.22l2.39-0.96c0.5,0.38,1.03,0.7,1.62,0.94l0.36,2.54 c0.05,0.24,0.24,0.41,0.48,0.41h3.84c0.24,0,0.44-0.17,0.47-0.41l0.36-2.54c0.59-0.24,1.13-0.56,1.62-0.94l2.39,0.96 c0.22,0.08,0.47,0,0.59-0.22l1.92-3.32c0.12-0.22,0.07-0.47-0.12-0.61L19.14,12.94z M12,15.6c-1.98,0-3.6-1.62-3.6-3.6 s1.62-3.6,3.6-3.6s3.6,1.62,3.6,3.6S13.98,15.6,12,15.6z"/></svg>
    </button>
    <button @click.stop="store.dispatch(store.state.editMode ? 'stopEditMode' : 'startEditMode')" title="Edit mode" class="w-10 h-10 flex items-center justify-center rounded-full transition-all duration-200" :class="store.state.editMode ? 'text-blue-500 bg-blue-50 dark:bg-blue-500/15' : 'text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100/80 dark:hover:text-neutral-200 dark:hover:bg-white/10'">
      <svg v-if="!store.state.editMode" xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/></svg>
      <svg v-else xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M9 16.2L4.8 12l-1.4 1.4L9 19 21 7l-1.4-1.4L9 16.2z"/></svg>
    </button>
    <div class="w-px h-5 bg-neutral-200 dark:bg-neutral-700 mx-0.5"></div>
    <button @click="store.dispatch('openNewLink')" title="Add link" class="w-10 h-10 flex items-center justify-center rounded-full text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100/80 dark:hover:text-neutral-200 dark:hover:bg-white/10 transition-all duration-200">
      <svg xmlns="http://www.w3.org/2000/svg" height="22" viewBox="0 0 24 24" width="22"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/></svg>
    </button>
  </div>

  <!-- Mobile toolbar toggle + expandable panel -->
  <div class="sm:hidden fixed bottom-4 right-4 z-20 flex flex-col items-end gap-2">
    <Transition name="toolbar">
      <div v-if="mobileToolbar" class="flex items-center gap-1 px-2 py-2 rounded-full bg-white/70 dark:bg-neutral-900/70 backdrop-blur-xl border border-white/60 dark:border-white/10 shadow-lg shadow-black/5 dark:shadow-black/30">
        <button @click="store.commit('setSettingsWindow', true); mobileToolbar = false" title="Settings" class="w-10 h-10 flex items-center justify-center rounded-full text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100/80 dark:hover:text-neutral-200 dark:hover:bg-white/10 transition-all duration-200">
          <svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0,0h24v24H0V0z" fill="none"/><path fill="currentColor" d="M19.14,12.94c0.04-0.3,0.06-0.61,0.06-0.94c0-0.32-0.02-0.64-0.07-0.94l2.03-1.58c0.18-0.14,0.23-0.41,0.12-0.61 l-1.92-3.32c-0.12-0.22-0.37-0.29-0.59-0.22l-2.39,0.96c-0.5-0.38-1.03-0.7-1.62-0.94L14.4,2.81c-0.04-0.24-0.24-0.41-0.48-0.41 h-3.84c-0.24,0-0.43,0.17-0.47,0.41L9.25,5.35C8.66,5.59,8.12,5.92,7.63,6.29L5.24,5.33c-0.22-0.08-0.47,0-0.59,0.22L2.74,8.87 C2.62,9.08,2.66,9.34,2.86,9.48l2.03,1.58C4.84,11.36,4.8,11.69,4.8,12s0.02,0.64,0.07,0.94l-2.03,1.58 c-0.18,0.14-0.23,0.41-0.12,0.61l1.92,3.32c0.12,0.22,0.37,0.29,0.59,0.22l2.39-0.96c0.5,0.38,1.03,0.7,1.62,0.94l0.36,2.54 c0.05,0.24,0.24,0.41,0.48,0.41h3.84c0.24,0,0.44-0.17,0.47-0.41l0.36-2.54c0.59-0.24,1.13-0.56,1.62-0.94l2.39,0.96 c0.22,0.08,0.47,0,0.59-0.22l1.92-3.32c0.12-0.22,0.07-0.47-0.12-0.61L19.14,12.94z M12,15.6c-1.98,0-3.6-1.62-3.6-3.6 s1.62-3.6,3.6-3.6s3.6,1.62,3.6,3.6S13.98,15.6,12,15.6z"/></svg>
        </button>
        <button @click.stop="store.dispatch(store.state.editMode ? 'stopEditMode' : 'startEditMode'); mobileToolbar = false" title="Edit mode" class="w-10 h-10 flex items-center justify-center rounded-full transition-all duration-200" :class="store.state.editMode ? 'text-blue-500 bg-blue-50 dark:bg-blue-500/15' : 'text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100/80 dark:hover:text-neutral-200 dark:hover:bg-white/10'">
          <svg v-if="!store.state.editMode" xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/></svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M9 16.2L4.8 12l-1.4 1.4L9 19 21 7l-1.4-1.4L9 16.2z"/></svg>
        </button>
        <div class="w-px h-5 bg-neutral-200 dark:bg-neutral-700 mx-0.5"></div>
        <button @click="store.dispatch('openNewLink'); mobileToolbar = false" title="Add link" class="w-10 h-10 flex items-center justify-center rounded-full text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100/80 dark:hover:text-neutral-200 dark:hover:bg-white/10 transition-all duration-200">
          <svg xmlns="http://www.w3.org/2000/svg" height="22" viewBox="0 0 24 24" width="22"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/></svg>
        </button>
      </div>
    </Transition>
    <button @click.stop="mobileToolbar = !mobileToolbar" class="w-11 h-11 flex items-center justify-center rounded-full bg-white/70 dark:bg-neutral-900/70 backdrop-blur-xl border border-white/60 dark:border-white/10 shadow-lg shadow-black/5 dark:shadow-black/30 text-neutral-500 dark:text-neutral-400 transition-all duration-200">
      <svg v-if="!mobileToolbar" xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z"/></svg>
      <svg v-else xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/></svg>
    </button>
  </div>

  <Transition>
    <settings v-if="store.state.window.settings" @click="store.commit('setSettingsWindow', false)"/>
  </Transition>
  <Transition>
    <delete v-if="store.state.window.delete" @click="store.dispatch('closeDelete')"/>
  </Transition>
  <Transition>
    <edit v-if="store.state.window.edit" @click="store.dispatch('closeEditLink')" :value="store.state.editObject"/>
  </Transition>
</template>

<script setup>
import draggable from "vuedraggable"
import CLink from "./components/link.vue"
import {onBeforeMount, onMounted, ref} from "vue"
import Cookies from "js-cookie"
import {useStore} from "vuex"
import Settings from "./components/settings.vue"
import Edit from "./components/edit.vue"
import Delete from "./components/delete.vue"

const store = useStore()
const mobileToolbar = ref(false)

const currentClass = Cookies.get("theme")
if (currentClass) document.documentElement.classList.add(currentClass)

onMounted(() => {
  document.body.addEventListener("click", (e) => {
    if (store.state.editMode) store.dispatch("stopEditMode")
  })
})
onBeforeMount(() => {
  store.commit("initialize", window.config && window.config.links ? window.config : {version: "demo"})
})
</script>

<style>
.v-enter-active,
.v-leave-active {
  transition: opacity 200ms ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}

.toolbar-enter-active,
.toolbar-leave-active {
  transition: opacity 150ms ease, transform 150ms ease;
}

.toolbar-enter-from,
.toolbar-leave-to {
  opacity: 0;
  transform: translateY(8px) scale(0.95);
}
</style>
