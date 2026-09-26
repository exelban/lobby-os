<template>
  <div @mousedown="startLongPress" @mouseup="stopLongPress" @mouseleave="stopLongPress" class="link-item relative z-0" :class="editMode ? 'animate-shaking cursor-move' : ''">
    <div class="h-full" :class="editMode ? 'pointer-events-none' : ''">
      <a class="link-card block h-full rounded-2xl" :href="preview || editMode ? undefined : linkURL" :tabindex="preview || editMode ? -1 : 0" @click="preventNavigation" :title="value.name" :target="store.state.openNewTab ? '_blank' : '_self'" rel="noopener noreferrer">
        <div class="link-surface relative h-full flex items-center gap-2 rounded-2xl px-3" :style="store.state.fullCardColor ? { background: color, color: getContrast(color), borderColor: 'transparent' } : {}">
          <div class="flex flex-1 min-w-0 gap-3 items-center">
            <div class="link-badge shrink-0 rounded-xl flex items-center justify-center" :style="{ backgroundColor: color, color: getContrast(color) }">
              <img v-if="icon && !iconFailed" :src="icon" alt="" class="link-icon object-contain" @error="iconFailed = true"/>
              <svg v-else class="link-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M14.8284 12L16.2426 13.4142L19.071 10.5858C20.6331 9.02365 20.6331 6.49099 19.071 4.9289C17.509 3.3668 14.9763 3.3668 13.4142 4.9289L10.5858 7.75732L12 9.17154L14.8284 6.34311C15.6095 5.56206 16.8758 5.56206 17.6568 6.34311C18.4379 7.12416 18.4379 8.39049 17.6568 9.17154L14.8284 12Z" fill="currentColor"/>
                <path d="M12 14.8285L13.4142 16.2427L10.5858 19.0711C9.02372 20.6332 6.49106 20.6332 4.92896 19.0711C3.36686 17.509 3.36686 14.9764 4.92896 13.4143L7.75739 10.5858L9.1716 12L6.34317 14.8285C5.56212 15.6095 5.56212 16.8758 6.34317 17.6569C7.12422 18.4379 8.39055 18.4379 9.1716 17.6569L12 14.8285Z" fill="currentColor"/>
                <path d="M14.8285 10.5857C15.219 10.1952 15.219 9.56199 14.8285 9.17147C14.4379 8.78094 13.8048 8.78094 13.4142 9.17147L9.1716 13.4141C8.78107 13.8046 8.78107 14.4378 9.1716 14.8283C9.56212 15.2188 10.1953 15.2188 10.5858 14.8283L14.8285 10.5857Z" fill="currentColor"/>
              </svg>
            </div>
            <span class="link-name truncate first-letter:capitalize font-normal">{{ value.name }}</span>
          </div>
          <div class="shrink-0 opacity-40">
            <svg xmlns="http://www.w3.org/2000/svg" height="14" viewBox="0 0 24 24" width="14"><path d="M0 0h24v24H0z" fill="none"/><path fill="currentColor" d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"/></svg>
          </div>
        </div>
      </a>
    </div>

      <button v-if="editMode" type="button" class="link-action edit-action" title="Edit link" aria-label="Edit link" @pointerdown.stop @mousedown.stop @touchstart.stop @click.stop="store.openEditLink(value)">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="m16 3 5 5M4 15 16.5 2.5a2.1 2.1 0 0 1 3 0l2 2a2.1 2.1 0 0 1 0 3L9 20l-6 1 1-6Z"/>
        </svg>
      </button>
      <button v-if="editMode" type="button" class="link-action delete-action" title="Delete link" aria-label="Delete link" @pointerdown.stop @mousedown.stop @touchstart.stop @click.stop="store.openDelete(value.id)">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="M3 6h18M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2M5 6l1 14a1 1 0 0 0 1 1h10a1 1 0 0 0 1-1l1-14M10 10v7M14 10v7"/>
        </svg>
      </button>
  </div>
</template>

<script setup>
import { computed, ref, watch, onBeforeUnmount } from "vue"
import presets from "../presets.js"
import { normalizeURL } from "../utils/urls.js"
import store from "../utils/store.js"

const props = defineProps({
  value: {
    type: Object,
    required: true,
  },
  preview: Boolean,
})
const linkURL = computed(() => {
  try { return normalizeURL(props.value.url) } catch { return undefined }
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
const iconFailed = ref(false)
watch(icon, () => { iconFailed.value = false })

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
  if (props.preview || editMode.value) return
  if (!("buttons" in e ? e.buttons === 1 : e.which === 1)) return
  longPressTimer = setTimeout(() => {
    store.startEditMode()
  }, 500)
}
const stopLongPress = () => {
  clearTimeout(longPressTimer)
}
const preventNavigation = (event) => {
  if (props.preview || editMode.value) event.preventDefault()
}
onBeforeUnmount(stopLongPress)
</script>

<style scoped>
.link-item {
  width: var(--card-width);
  height: calc(70px * var(--card-scale));
}

.link-action {
  position: absolute;
  top: -8px;
  z-index: 10;
  display: grid;
  place-items: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 1px solid light-dark(rgb(255 255 255 / 75%), rgb(255 255 255 / 15%));
  background: light-dark(rgb(255 255 255 / 90%), rgb(25 29 37 / 90%));
  color: light-dark(#555e6d, #c4cad4);
  box-shadow: 0 2px 6px rgb(0 0 0 / 12%);
  backdrop-filter: blur(12px);
  cursor: pointer;
  transition: background 150ms ease, color 150ms ease;
}

.edit-action { right: -8px; }
.delete-action { left: -8px; }
.link-action:focus-visible { outline: 2px solid currentColor; outline-offset: 2px; }
.edit-action:is(:hover, :focus-visible) { color: light-dark(#2563eb, #93b4ff); background: light-dark(#dbeafe, #263857); }
.delete-action:is(:hover, :focus-visible) { color: light-dark(#dc2626, #fca5a5); background: light-dark(#fee2e2, #512a30); }

.link-badge {
  width: var(--card-badge);
  height: var(--card-badge);
}

.link-icon {
  width: var(--card-icon);
  height: var(--card-icon);
}

.link-name {
  font-size: calc(14px * var(--card-scale));
}

.link-surface {
  color: light-dark(var(--text-color), var(--dark-text-color));
  background: light-dark(var(--card-bg), var(--dark-card-bg));
  border: 1px solid light-dark(var(--card-border), var(--dark-card-border));
  box-shadow: 0 2px 8px rgb(0 0 0 / 4%);
  backdrop-filter: blur(12px);
  transition: background 180ms ease, box-shadow 180ms ease, transform 180ms ease;
}

@media (hover: hover) {
  .link-card[href]:hover .link-surface {
    background: light-dark(rgb(255 255 255 / 90%), rgb(255 255 255 / 10%));
    box-shadow: 0 4px 14px rgb(0 0 0 / 8%);
    transform: translateY(-2px);
  }
}
</style>
