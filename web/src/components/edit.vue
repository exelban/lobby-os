<template>
  <Modal labelledby="edit-title" :busy="Boolean(status)" @close="store.closeEditLink()">
    <div @click.stop class="relative w-full max-w-md max-h-full overflow-y-auto">
      <div class="relative bg-white/85 dark:bg-neutral-900/90 backdrop-blur-xl rounded-2xl shadow-2xl shadow-black/10 dark:shadow-black/40 border border-white/60 dark:border-white/10 overflow-hidden">
        <div class="flex items-center justify-between px-5 py-4 border-b border-neutral-200/50 dark:border-neutral-700/50">
          <h3 id="edit-title" tabindex="-1" autofocus class="text-lg font-semibold text-neutral-800 dark:text-white">{{ value.id ? "Edit" : "New" }} link</h3>
          <button @click="store.closeEditLink()" :disabled="Boolean(status)" aria-label="Close editor" class="w-8 h-8 flex items-center justify-center text-neutral-400 hover:text-neutral-600 dark:hover:text-neutral-200 hover:bg-neutral-100/80 dark:hover:bg-white/10 rounded-lg transition-all duration-200">
            <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <form class="p-5 flex flex-col gap-4" @submit.prevent="addOrSaveLink">
          <div class="flex flex-col gap-3 text-sm">
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-neutral-500 dark:text-neutral-400 uppercase tracking-wider" for="name">Name</label>
              <input v-model="value.name" id="name" autocapitalize="off" placeholder="Link name" type="text" class="form-field" required>
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-neutral-500 dark:text-neutral-400 uppercase tracking-wider" for="url">URL</label>
              <input v-model="value.url" @blur="normalizeAddress" id="url" autocapitalize="off" placeholder="nas.local:5000 or https://…" type="text" class="form-field" required>
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-neutral-500 dark:text-neutral-400 uppercase tracking-wider" for="group">Group</label>
              <select v-if="!newGroup" v-model="value.group" id="group" class="form-field select-field">
                <option v-for="group in store.groupOptions" :key="group.id" :value="group.id">{{ group.label }}</option>
              </select>
              <input v-else v-model="newGroupName" id="group" placeholder="Group name (optional)" class="form-field">
              <button v-if="store.groupOptions.length" type="button" @click="newGroup = !newGroup" class="self-start text-xs text-blue-600 dark:text-blue-400">{{ newGroup ? "Choose existing group" : "Create new group" }}</button>
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-neutral-500 dark:text-neutral-400 uppercase tracking-wider" for="preset">Preset</label>
              <select id="preset" v-model="value.preset" @change="onPresetChange" class="form-field select-field">
                <option v-bind:value="undefined">None</option>
                <option v-for="(p, key) in presets" v-bind:value="key">{{ p.name }}</option>
              </select>
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-neutral-500 dark:text-neutral-400 uppercase tracking-wider" for="color">Color</label>
              <div class="flex flex-row items-center bg-white/60 dark:bg-white/5 border border-neutral-200/80 dark:border-white/10 rounded-xl overflow-hidden focus-within:ring-2 focus-within:ring-blue-500/30 focus-within:border-blue-400/50 transition-all duration-200">
                <input v-model="value.color" id="color" type="text" placeholder="#FFFFFF" class="w-full py-2.5 px-4 bg-transparent text-neutral-800 dark:text-neutral-200 placeholder-neutral-400 dark:placeholder-neutral-500 leading-tight focus:outline-hidden">
                <input v-model="value.color" type="color" class="w-[50px] h-[38px] bg-transparent border-none cursor-pointer mr-1">
              </div>
            </div>
            <div v-if="store.state.version !== 'demo'">
              <button type="button" @click="loadIcon" :disabled="!isUrlValid" class="w-full py-2 text-xs font-medium rounded-xl transition-all duration-200" :class="isUrlValid ? 'bg-blue-500/10 text-blue-600 dark:text-blue-400 hover:bg-blue-500/20' : 'bg-neutral-100/60 dark:bg-white/5 text-neutral-400 dark:text-neutral-600 cursor-not-allowed'">
                Load Site Favicon
              </button>
            </div>
          </div>

          <div class="w-full h-[90px] flex items-center justify-center">
            <CLink :value="value" :preview="true"/>
          </div>

          <p v-if="error" role="alert" class="text-sm text-red-600 dark:text-red-400">{{ error }}</p>

          <button type="submit" class="text-white bg-blue-500 hover:bg-blue-600 font-medium rounded-xl text-sm px-5 py-3 text-center transition-all duration-200 shadow-xs shadow-blue-500/20">
            {{ value.id ? "Save" : "Add" }}
          </button>
        </form>
      </div>

      <div v-if="status" class="absolute inset-0 bg-white/80 dark:bg-neutral-900/80 backdrop-blur-xs flex justify-center items-center rounded-2xl z-50">
        <svg v-if="status === 'in_progress'" class="w-10 h-10 text-neutral-200 animate-spin dark:text-neutral-700 fill-blue-500" viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" fill="currentColor"/><path d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" fill="currentFill"/></svg>
        <svg v-else-if="status === 'ok'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-12 h-12 text-emerald-500"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M9 16.2L4.8 12l-1.4 1.4L9 19 21 7l-1.4-1.4L9 16.2z"/></svg>
        <div v-else class="flex flex-col justify-center items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-12 h-12 text-red-500"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/></svg>
          <p class="text-sm text-neutral-600 dark:text-neutral-400">{{ status }}</p>
        </div>
      </div>
    </div>
  </Modal>
</template>

<script setup>
import Modal from "./modal.vue"
import CLink from "./link.vue"
import store from "../utils/store.js"
import presets from "../presets.js"
import { normalizeURL } from "../utils/urls.js"
import { createGroupID } from "../utils/store.js"
import {ref, computed, watch} from "vue"

const props = defineProps(["value"])
const status = ref(undefined)
const error = ref("")
const newGroup = ref(!store.groupOptions.length)
const newGroupName = ref("")

watch(() => props.value, value => {
  if (!value.color && presets[value.preset]?.color) {
    value.color = presets[value.preset].color
  }
}, { immediate: true })

const isUrlValid = computed(() => {
  try {
    normalizeURL(props.value.url)
    return true
  } catch {
    return false
  }
})

let suggestedName = ""
watch(() => props.value.url, url => {
  if (props.value.id) return
  try {
    const name = new URL(normalizeURL(url)).hostname.replace(/^www\./, "")
    if (!props.value.name || props.value.name === suggestedName) props.value.name = name
    suggestedName = name
  } catch {}
})
const normalizeAddress = () => {
  if (isUrlValid.value) props.value.url = normalizeURL(props.value.url)
}

const onPresetChange = () => {
  props.value.icon = null
  const color = presets[props.value.preset]?.color
  if (color) {
    props.value.color = color
  }
}

const addOrSaveLink = (e) => {
  e.preventDefault()
  if (status.value) return
  try { props.value.url = normalizeURL(props.value.url) } catch (e) { error.value = e.message; return }
  if (!props.value.id) props.value.id = Math.random().toString(36).substring(2, 9)
  error.value = ""
  status.value = "in_progress"
  const link = newGroup.value ? {...props.value, group: createGroupID(), groupName: newGroupName.value} : props.value
  store.upsert(link).then(() => {
    status.value = "ok"
    setTimeout(() => store.closeEditLink(), 500)
  }).catch((e) => {
    status.value = undefined
    error.value = e.message || "Could not save the link. Please try again."
  })
}

const loadIcon = () => {
  error.value = ""
  status.value = "in_progress"
  store.loadFaviconURL(props.value.url).then((icon) => {
    props.value.icon = icon
  }).catch((e) => {
    error.value = e.message
  }).finally(() => {
    status.value = undefined
  })
}
</script>

<style scoped>
@reference "../style.css";
.form-field {
  @apply bg-white/60 dark:bg-white/5 border border-neutral-200/80 dark:border-white/10 rounded-xl w-full py-2.5 px-4 text-neutral-800 dark:text-neutral-200 placeholder-neutral-400 dark:placeholder-neutral-500 leading-tight focus:outline-hidden focus:ring-2 focus:ring-blue-500/30 focus:border-blue-400/50 transition-all duration-200;
}
.select-field {
  appearance: none;
  -webkit-appearance: none;
  padding-right: 40px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%23888' stroke-width='1.8' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 14px center;
}
</style>
