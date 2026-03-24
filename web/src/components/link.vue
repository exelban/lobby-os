<template>
  <div @mousedown="startLongPress" @mouseup="stopLongPress" class="relative z-0 w-[175px] sm:w-[210px] h-[70px]" :class="editMode ? 'animate-shaking cursor-move' : ''">
    <div class="h-full" :class="editMode ? 'pointer-events-none' : ''">
      <a :href="value.url" target="_blank">
        <div class="relative overflow-hidden h-full flex flex-row bg-white items-center gap-3 rounded-xl shadow opacity-90 hover:opacity-100 hover:shadow-xl transition-shadow ease-in-out duration-300" :style="{'backgroundColor': color, 'color': getContrast(color)}">
          <div class="flex flex-row p-4 pr-0 gap-3 items-center overflow-hidden">
            <div class="max-w-[28px] min-w-[28px] sm:max-w-[38px] sm:min-w-[38px]">
              <img v-if="icon" :src="icon" alt="Icon"/>
              <svg v-else :style="{color: getContrast(color)}" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M14.8284 12L16.2426 13.4142L19.071 10.5858C20.6331 9.02365 20.6331 6.49099 19.071 4.9289C17.509 3.3668 14.9763 3.3668 13.4142 4.9289L10.5858 7.75732L12 9.17154L14.8284 6.34311C15.6095 5.56206 16.8758 5.56206 17.6568 6.34311C18.4379 7.12416 18.4379 8.39049 17.6568 9.17154L14.8284 12Z" fill="currentColor"/>
                <path d="M12 14.8285L13.4142 16.2427L10.5858 19.0711C9.02372 20.6332 6.49106 20.6332 4.92896 19.0711C3.36686 17.509 3.36686 14.9764 4.92896 13.4143L7.75739 10.5858L9.1716 12L6.34317 14.8285C5.56212 15.6095 5.56212 16.8758 6.34317 17.6569C7.12422 18.4379 8.39055 18.4379 9.1716 17.6569L12 14.8285Z" fill="currentColor"/>
                <path d="M14.8285 10.5857C15.219 10.1952 15.219 9.56199 14.8285 9.17147C14.4379 8.78094 13.8048 8.78094 13.4142 9.17147L9.1716 13.4141C8.78107 13.8046 8.78107 14.4378 9.1716 14.8283C9.56212 15.2188 10.1953 15.2188 10.5858 14.8283L14.8285 10.5857Z" fill="currentColor"/>
              </svg>
            </div>
            <h1 class="overflow-hidden text-ellipsis whitespace-nowrap first-letter:capitalize text-sm">{{ value.name }}</h1>
          </div>
          <div class="flex justify-center items-center ml-auto h-full bg-neutral-100 bg-opacity-30 dark:bg-neutral-500 dark:bg-opacity-30">
            <svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"/></svg>
          </div>
        </div>
      </a>
    </div>

    <div v-if="editMode" class="w-5 h-5 bg-white text-neutral-900 hover:bg-red-500 hover:text-white dark:bg-black dark:text-white shadow rounded-full absolute -top-2 -left-2 z-10 cursor-pointer transition ease-in-out duration-200" title="Delete link" @click="store.dispatch('openDelete', props.value.id)">
      <svg viewBox="-5 -5 34 34"  fill="none" xmlns="http://www.w3.org/2000/svg">
        <path fill-rule="evenodd" clip-rule="evenodd" d="M17 5V4C17 2.89543 16.1046 2 15 2H9C7.89543 2 7 2.89543 7 4V5H4C3.44772 5 3 5.44772 3 6C3 6.55228 3.44772 7 4 7H5V18C5 19.6569 6.34315 21 8 21H16C17.6569 21 19 19.6569 19 18V7H20C20.5523 7 21 6.55228 21 6C21 5.44772 20.5523 5 20 5H17ZM15 4H9V5H15V4ZM17 7H7V18C7 18.5523 7.44772 19 8 19H16C16.5523 19 17 18.5523 17 18V7Z" fill="currentColor"/>
        <path d="M9 9H11V17H9V9Z" fill="currentColor" />
        <path d="M13 9H15V17H13V9Z" fill="currentColor" />
      </svg>
    </div>
    <div v-if="editMode" class="w-5 h-5 bg-white text-neutral-900 hover:bg-blue-500 hover:text-white dark:bg-black dark:text-white shadow rounded-full absolute -top-2 -right-2 z-10 cursor-pointer transition ease-in-out duration-200" title="Edit link" @click="store.dispatch('openEditLink', {...value})">
      <svg viewBox="-10 -10 44 44"  fill="none" xmlns="http://www.w3.org/2000/svg">
        <path fill-rule="evenodd" clip-rule="evenodd" d="M21.2635 2.29289C20.873 1.90237 20.2398 1.90237 19.8493 2.29289L18.9769 3.16525C17.8618 2.63254 16.4857 2.82801 15.5621 3.75165L4.95549 14.3582L10.6123 20.0151L21.2189 9.4085C22.1426 8.48486 22.338 7.1088 21.8053 5.99367L22.6777 5.12132C23.0682 4.7308 23.0682 4.09763 22.6777 3.70711L21.2635 2.29289ZM16.9955 10.8035L10.6123 17.1867L7.78392 14.3582L14.1671 7.9751L16.9955 10.8035ZM18.8138 8.98525L19.8047 7.99429C20.1953 7.60376 20.1953 6.9706 19.8047 6.58007L18.3905 5.16586C18 4.77534 17.3668 4.77534 16.9763 5.16586L15.9853 6.15683L18.8138 8.98525Z" fill="currentColor"/>
        <path d="M2 22.9502L4.12171 15.1717L9.77817 20.8289L2 22.9502Z" fill="currentColor"/>
      </svg>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue"
import presets from "../presets.js"
import {useStore} from "vuex"

const store = useStore()
const props = defineProps({
  value: {
    type: {
      id: {
        type: String,
        required: true,
      },
      url: {
        type: String,
        required: true,
      },
      name: {
        type: String,
        required: true,
      },
      color: String,
      icon: String,
      preview: Boolean,
    },
    required: true,
  },
  preview: Boolean,
})
const color = computed(() => {
  if (props.value.color) return props.value.color
  if (!props.value.preset) return "#ffffff"
  if (presets[props.value.preset] && presets[props.value.preset].color) return presets[props.value.preset].color
  return "#ffffff"
})
const icon = computed(() => {
  if (props.value.icon) return props.value.icon
  if (!props.value.preset) return null
  if (presets[props.value.preset]) return presets[props.value.preset].icon
  return null
})
const editMode = computed(() => store.state.editMode && !props.preview)
let longPressTimer = null

const getContrast = (color) => {
  if (color.slice(0, 1) === "#") {
    color = color.slice(1)
  }
  if (color.length === 3) {
    color = color.split('').map(hex => hex + hex).join('')
  }

  const r = parseInt(color.substr(0,2),16)
  const g = parseInt(color.substr(2,2),16)
  const b = parseInt(color.substr(4,2),16)
  const yiq = ((r * 299) + (g * 587) + (b * 114)) / 1000

  return (yiq >= 128) ? "black" : "white"
}

const startLongPress = (e) => {
  if (!("buttons" in e ? e.buttons === 1 : e.which === 1)) return
  longPressTimer = setTimeout(() => {
    store.dispatch("startEditMode")
    e.preventDefault()
  }, 500)
}
const stopLongPress = () => {
  clearTimeout(longPressTimer)
}
</script>
