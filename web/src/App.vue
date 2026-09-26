<template>
  <div class="card-layout contents" :style="{ '--card-scale': { small: 0.85, medium: 1, large: 1.2 }[store.state.cardSize] }">
  <div v-if="searchOpen" class="fixed top-4 left-1/2 -translate-x-1/2 z-20 w-[min(90vw,400px)] flex items-center gap-2 rounded-xl border border-white/60 dark:border-white/10 bg-white/90 dark:bg-neutral-900/90 p-3 shadow-lg backdrop-blur-xl">
    <input ref="searchInput" v-model="query" type="search" aria-label="Search links" placeholder="Search links or hostnames…" class="min-w-0 flex-1 bg-transparent text-sm outline-none" @keydown.escape.stop="closeSearch" @keydown.enter.prevent="openFirstResult">
    <button @click="closeSearch" aria-label="Close search" class="text-sm text-neutral-500">Esc</button>
  </div>
  <main @click="onDashboardClick" class="min-h-dvh p-1 sm:p-2 flex flex-row justify-center items-center content-center">
    <div v-if="Object.keys(store.state.links).length" class="container flex flex-col items-center gap-4 sm:gap-6">
      <div v-for="g in visibleGroups" :key="g" class="max-w-full">
        <div v-if="store.state.groupNames[g] || store.state.editMode" class="px-2 mb-2 text-center">
          <input v-if="store.state.editMode" :value="store.state.groupNames[g]" :disabled="savingGroup || savingOrder" @change="renameGroup(g, $event)" @keydown.enter="$event.target.blur()" @keydown.escape.stop.prevent="cancelGroupEdit(g, $event)" aria-label="Group name" placeholder="Group name (optional)" class="max-w-full w-48 text-center text-sm font-medium bg-transparent rounded-lg border border-neutral-400/30 py-1 px-2">
          <h2 v-else class="text-sm font-medium text-neutral-600 dark:text-neutral-400">{{ store.state.groupNames[g] }}</h2>
        </div>
        <ul v-if="query.trim()" class="flex flex-wrap justify-center gap-x-4 sm:gap-x-5 gap-y-2 sm:gap-y-3 p-1 sm:p-2 border border-transparent mx-auto" :class="store.state.linksPerRow ? 'limit-row' : ''" :style="{ '--lpr': store.state.linksPerRow }">
          <li v-for="link in store.state.links[g].filter(matchesSearch)" :key="link.id"><CLink :value="link"/></li>
        </ul>
        <draggable v-else v-model="store.state.links[g]" @start="beforeDrag = JSON.parse(JSON.stringify(store.state.links))" @end="saveOrder" tag="ul" group="links" :disabled="!store.state.editMode || savingOrder || savingGroup" item-key="id" ghost-class="opacity-30" class="flex flex-wrap items-center justify-center flex-row gap-x-4 sm:gap-x-5 gap-y-2 sm:gap-y-3 p-1 sm:p-2 rounded-lg mx-auto" :style="store.state.linksPerRow ? { '--lpr': store.state.linksPerRow } : {}" :class="[store.state.editMode ? 'min-w-[190px] min-h-[86px] bg-neutral-500/10 border border-dashed border-neutral-400/40 dark:border-neutral-600/40' : 'border border-transparent', store.state.linksPerRow ? 'limit-row' : '']">
          <template #item="{element}">
            <li><CLink @click.stop :value="element"></CLink></li>
          </template>
        </draggable>
      </div>

      <p v-if="!visibleGroups.length" role="status" class="text-sm text-neutral-500 dark:text-neutral-400">No links match “{{ query }}”.</p>
      <div v-if="store.state.editMode" class="h-10">
        <button @click.stop="query = ''; store.addGroup()" title="Add group" aria-label="Add group" class="w-9 h-9 flex items-center justify-center bg-neutral-200/60 hover:bg-neutral-300/60 dark:bg-neutral-700/50 dark:hover:bg-neutral-600/50 text-neutral-500 dark:text-neutral-400 rounded-full backdrop-blur-xs transition-all duration-200 text-lg">+</button>
      </div>
    </div>
    <button v-else @click="store.openNewLink()" class="shadow-xs text-neutral-600 bg-white/70 backdrop-blur-xs border border-white/80 hover:bg-white/90 font-medium rounded-2xl text-sm px-8 py-4 dark:bg-white/5 dark:text-neutral-300 dark:border-white/10 dark:hover:bg-white/10 transition-all duration-300">Create first link</button>
  </main>

  <p v-if="saveError" role="alert" class="fixed top-4 left-1/2 -translate-x-1/2 z-50 max-w-[90vw] rounded-xl bg-red-50 text-red-700 p-4 shadow-lg">
    {{ saveError }} <button @click="saveError = ''" class="ml-2 underline">Dismiss</button>
  </p>

  <!-- Floating toolbar - always visible on desktop, toggle button on mobile -->
  <div class="fixed bottom-5 left-1/2 -translate-x-1/2 sm:flex hidden items-center gap-1 px-2 py-2 rounded-full bg-white/70 dark:bg-neutral-900/70 backdrop-blur-xl border border-white/60 dark:border-white/10 shadow-lg shadow-black/5 dark:shadow-black/30 transition-all duration-300 z-20">
    <Toolbar @search="openSearch"/>
  </div>

  <!-- Mobile toolbar toggle + expandable panel -->
  <div class="sm:hidden fixed bottom-4 right-4 z-20 flex flex-col items-end gap-2">
    <Transition name="toolbar">
      <div v-if="mobileToolbar" class="flex items-center gap-1 px-2 py-2 rounded-full bg-white/70 dark:bg-neutral-900/70 backdrop-blur-xl border border-white/60 dark:border-white/10 shadow-lg shadow-black/5 dark:shadow-black/30">
        <Toolbar @search="openSearch" @select="mobileToolbar = false"/>
      </div>
    </Transition>
    <button @click.stop="mobileToolbar = !mobileToolbar" aria-label="Toggle toolbar" :aria-expanded="mobileToolbar" class="w-11 h-11 flex items-center justify-center rounded-full bg-white/70 dark:bg-neutral-900/70 backdrop-blur-xl border border-white/60 dark:border-white/10 shadow-lg shadow-black/5 dark:shadow-black/30 text-neutral-500 dark:text-neutral-400 transition-all duration-200">
      <svg v-if="!mobileToolbar" xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z"/></svg>
      <svg v-else xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 0 24 24" width="20"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/></svg>
    </button>
  </div>

  <Transition>
    <settings v-if="store.state.window.settings"/>
  </Transition>
  <Transition>
    <delete v-if="store.state.window.delete"/>
  </Transition>
  <Transition>
    <edit v-if="store.state.window.edit" :value="store.state.editObject"/>
  </Transition>
  </div>
</template>

<script setup>
import draggable from "vuedraggable"
import Toolbar from "./components/toolbar.vue"
import CLink from "./components/link.vue"
import {onBeforeMount, onMounted, onBeforeUnmount, computed, ref, nextTick} from "vue"
import store from "./utils/store.js"
import Settings from "./components/settings.vue"
import Edit from "./components/edit.vue"
import Delete from "./components/delete.vue"

const mobileToolbar = ref(false)
const savingOrder = ref(false)
const saveError = ref("")
const savingGroup = ref(false)
const query = ref("")
const searchOpen = ref(false)
const searchInput = ref(null)

let beforeDrag = null

const matchesSearch = link => `${link.name} ${link.url}`.toLowerCase().includes(query.value.trim().toLowerCase())
const visibleGroups = computed(() => Object.keys(store.state.links).filter(group =>
  !query.value.trim() || store.state.links[group].some(matchesSearch)
))

const openSearch = async () => {
  searchOpen.value = true
  mobileToolbar.value = false
  await nextTick()
  searchInput.value?.focus()
}
const closeSearch = () => { query.value = ""; searchOpen.value = false }
const openFirstResult = () => {
  if (!store.state.editMode) document.querySelector("main .link-card[href]")?.click()
}
const finishEditing = () => {
  if (store.state.editMode && !savingOrder.value && !Object.values(store.state.window).some(Boolean)) {
    store.stopEditMode()
  }
}
const onDashboardClick = event => {
  if (event.target.closest(".link-item, input, textarea, select, button, a, [contenteditable]")) return
  finishEditing()
}
const cancelGroupEdit = (group, event) => {
  event.target.value = store.state.groupNames[group] ?? ""
  event.target.blur()
  finishEditing()
}
const onKeydown = event => {
  if (Object.values(store.state.window).some(Boolean)) return
  if (event.key === "Escape" && store.state.editMode) {
    event.preventDefault()
    finishEditing()
    return
  }
  if (event.key !== "/" || event.ctrlKey || event.metaKey || event.altKey) return
  if (event.target.closest("input, textarea, select") || event.target.isContentEditable) return
  event.preventDefault()
  openSearch()
}

onMounted(() => document.addEventListener("keydown", onKeydown))
onBeforeUnmount(() => document.removeEventListener("keydown", onKeydown))

const renameGroup = async (group, event) => {
  savingGroup.value = true
  try {
    await store.renameGroup({group, name: event.target.value})
    saveError.value = ""
  } catch (error) {
    event.target.value = store.state.groupNames[group] ?? ""
    saveError.value = error.message || "Could not rename the group."
  } finally {
    savingGroup.value = false
  }
}

const saveOrder = async () => {
  savingOrder.value = true
  await nextTick()
  try {
    await store.save()
    saveError.value = ""
  } catch {
    store.state.links = beforeDrag
    saveError.value = "Could not save the new order. Please try again."
  } finally {
    savingOrder.value = false
  }
}

store.setTheme(store.state.theme)
store.setBackground(store.state.background)

onBeforeMount(() => {
  store.initialize(window.config ?? {version: "demo"})
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

.card-layout {
  --card-width: calc(175px * var(--card-scale));
  --card-badge: calc(40px * var(--card-scale));
  --card-icon: calc(28px * var(--card-scale));
}

/* Include the gap, padding and border in the row width. */
.limit-row {
  max-width: calc(var(--lpr) * var(--card-width) + (var(--lpr) - 1) * 16px + 10px);
}
@media (min-width: 640px) {
  .card-layout {
    --card-width: calc(210px * var(--card-scale));
    --card-badge: calc(44px * var(--card-scale));
    --card-icon: calc(32px * var(--card-scale));
  }
  .limit-row {
    max-width: calc(var(--lpr) * var(--card-width) + (var(--lpr) - 1) * 20px + 18px);
  }
}
</style>
