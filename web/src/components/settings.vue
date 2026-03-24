<template>
  <div class="fixed z-30 top-0 left-0 w-full h-full p-3 sm:p-0 flex justify-center items-center bg-black/30 dark:bg-black/50 backdrop-blur-sm">
    <div @click.stop class="relative w-full max-w-md max-h-full">
      <div class="relative bg-white/85 dark:bg-neutral-900/90 backdrop-blur-xl rounded-2xl shadow-2xl shadow-black/10 dark:shadow-black/40 border border-white/60 dark:border-white/10 overflow-hidden">
        <div class="flex items-center justify-between px-5 py-4 border-b border-neutral-200/50 dark:border-neutral-700/50">
          <h3 class="text-lg font-semibold text-neutral-800 dark:text-white">Settings</h3>
          <button @click="store.commit('setSettingsWindow', false)" class="w-8 h-8 flex items-center justify-center text-neutral-400 hover:text-neutral-600 dark:hover:text-neutral-200 hover:bg-neutral-100/80 dark:hover:bg-white/10 rounded-lg transition-all duration-200">
            <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div class="p-5 flex flex-col gap-4">
          <div class="bg-white/50 dark:bg-white/5 rounded-xl border border-neutral-200/50 dark:border-white/5 p-4">
            <p class="text-xs font-medium text-neutral-500 dark:text-neutral-400 uppercase tracking-wider mb-3">Appearance</p>
            <div class="flex gap-2">
              <button v-for="opt in [{val: 'light', label: 'Light'}, {val: 'dark', label: 'Dark'}, {val: 'system', label: 'System'}]" :key="opt.val" @click="setTheme(opt.val)" type="button" class="flex-1 text-center py-2 px-3 rounded-lg text-sm font-medium transition-all duration-200 border" :class="theme === opt.val ? 'bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20' : 'border-transparent text-neutral-500 dark:text-neutral-400 hover:bg-neutral-100/60 dark:hover:bg-white/5'">
                {{ opt.label }}
              </button>
            </div>
          </div>

          <div class="bg-white/50 dark:bg-white/5 rounded-xl border border-neutral-200/50 dark:border-white/5 p-4">
            <p class="text-xs font-medium text-neutral-500 dark:text-neutral-400 uppercase tracking-wider mb-3">Links</p>
            <div class="flex gap-2">
              <button @click="download" class="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 text-sm font-medium text-blue-600 dark:text-blue-400 bg-blue-500/10 hover:bg-blue-500/20 rounded-xl transition-all duration-200">
                <svg xmlns="http://www.w3.org/2000/svg" height="16" viewBox="0 0 24 24" width="16" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 5v14m0 0l-6-6m6 6l6-6"/></svg>
                Export
              </button>
              <button @click="fileInput.click()" class="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 text-sm font-medium text-blue-600 dark:text-blue-400 bg-blue-500/10 hover:bg-blue-500/20 rounded-xl transition-all duration-200">
                <svg xmlns="http://www.w3.org/2000/svg" height="16" viewBox="0 0 24 24" width="16" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 19V5m0 0l-6 6m6-6l6 6"/></svg>
                Import
              </button>
            </div>
            <input ref="fileInput" type="file" accept=".json" class="hidden" @change="importFile">
          </div>
        </div>

        <div class="flex items-center justify-center px-5 py-3 border-t border-neutral-200/50 dark:border-neutral-700/50 text-xs text-neutral-400 dark:text-neutral-500">
          <a href="https://github.com/exelban/jad" target="_blank" class="hover:text-neutral-600 dark:hover:text-neutral-300 transition-colors duration-200">JAD {{ store.state.version }}</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useStore } from "vuex"
import {ref} from "vue"
import Cookies from "js-cookie"

const store = useStore()
const currentClass = Cookies.get("theme")
const theme = ref(currentClass || "system")
const fileInput = ref(null)

const importFile = (e) => {
  const file = e.target.files[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = (evt) => {
    try {
      const links = JSON.parse(evt.target.result)
      if (Array.isArray(links)) {
        store.dispatch("importLinks", links)
      }
    } catch {}
  }
  reader.readAsText(file)
  e.target.value = ""
}

const setTheme = (val) => {
  theme.value = val
  document.documentElement.classList.remove("dark", "light")
  if (val === "system") {
    Cookies.remove("theme")
  } else {
    document.documentElement.classList.add(val)
    Cookies.set("theme", val, { expires: 7*365 })
  }
}

const download = () => {
  const array = Object.keys(store.state.links).flatMap(k =>
      store.state.links[k].map(l => ({ ...l, group: k }))
  )
  const elem = document.createElement("a")
  elem.setAttribute("href", `data:text/json;charset=utf-8,${encodeURIComponent(JSON.stringify(array))}`)
  elem.setAttribute("download", "links.json")
  document.body.appendChild(elem)
  elem.click()
  elem.remove()
}
</script>
