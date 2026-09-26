<template>
  <dialog ref="dialog" class="modal-overlay" :aria-labelledby="labelledby" @cancel.prevent="close" @click.self="close" @keydown.escape.stop>
    <slot/>
  </dialog>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from "vue"

const props = defineProps({ labelledby: String, busy: Boolean })
const emit = defineEmits(["close"])
const dialog = ref(null)
const close = () => { if (!props.busy) emit("close") }
onMounted(() => dialog.value.showModal())
onBeforeUnmount(() => dialog.value.close())
</script>

<style>
.modal-overlay {
  position: fixed;
  inset: 0;
  margin: 0;
  width: 100%;
  height: 100dvh;
  max-width: none;
  max-height: none;
  padding: 16px;
  border: 0;
  background: rgb(0 0 0 / 25%);
  backdrop-filter: blur(4px);
  color: light-dark(#252932, #e5e7eb);
}
.modal-overlay[open] { display: flex; align-items: center; justify-content: center; }
.modal-overlay::backdrop { background: transparent; }
</style>
