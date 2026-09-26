<template>
  <Modal labelledby="settings-title" :busy="importing" @close="close">
    <div class="settings-panel" @click.stop>
      <header class="settings-header">
        <h2 id="settings-title" tabindex="-1" autofocus>Settings</h2>
        <button type="button" class="icon-button" aria-label="Close settings" @click="close">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path stroke-linecap="round" d="m6 6 12 12M6 18 18 6"/></svg>
        </button>
      </header>

      <div class="settings-body">
        <section class="settings-section" aria-labelledby="appearance-title">
          <h3 id="appearance-title">Appearance</h3>
          <fieldset>
            <legend>Color mode</legend>
            <div class="segmented">
              <label v-for="opt in [{val: 'light', label: 'Light'}, {val: 'dark', label: 'Dark'}, {val: 'system', label: 'System'}]" :key="opt.val" class="choice" :class="{ selected: theme === opt.val }">
                <input type="radio" name="theme" :value="opt.val" :checked="theme === opt.val" @change="setTheme(opt.val)" class="sr-only">
                {{ opt.label }}
              </label>
            </div>
          </fieldset>
          <div class="setting-row">
            <label for="background-theme">Background theme</label>
            <div class="select-control">
              <select id="background-theme" class="background-select" :value="store.state.background" @change="store.setBackground($event.target.value)">
                <option v-for="background in backgrounds" :key="background.id" :value="background.id">{{ background.label }}</option>
              </select>
              <svg class="select-chevron" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><path d="m7 10 5 5 5-5" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </div>
          </div>
          <div class="setting-row">
            <label for="full-card-color">Full card color</label>
            <button id="full-card-color" type="button" role="switch" :aria-checked="store.state.fullCardColor" @click="store.setFullCardColor(!store.state.fullCardColor)" class="switch" :class="{ enabled: store.state.fullCardColor }">
              <span></span>
            </button>
          </div>
        </section>

        <section class="settings-section" aria-labelledby="layout-title">
          <h3 id="layout-title">Layout</h3>
          <fieldset>
            <legend>Card size</legend>
            <div class="segmented">
              <label v-for="size in ['small', 'medium', 'large']" :key="size" class="choice capitalize" :class="{ selected: store.state.cardSize === size }">
                <input type="radio" name="card-size" :value="size" :checked="store.state.cardSize === size" @change="store.setCardSize(size)" class="sr-only">
                {{ size }}
              </label>
            </div>
          </fieldset>
          <div class="setting-row">
            <span id="links-per-row-label">Links per row</span>
            <div class="stepper" role="group" aria-labelledby="links-per-row-label">
              <button type="button" aria-label="Fewer links per row" :disabled="linksPerRow <= 0" @click="setLinksPerRow(linksPerRow - 1)">−</button>
              <output aria-live="polite">{{ linksPerRow === 0 ? "Auto" : linksPerRow }}</output>
              <button type="button" aria-label="More links per row" @click="setLinksPerRow(linksPerRow + 1)">+</button>
            </div>
          </div>
        </section>

        <section class="settings-section" aria-labelledby="links-title">
          <div class="setting-row">
            <h3 id="links-title">Your links</h3>
            <div class="flex gap-2">
              <button type="button" class="secondary-button" @click="download">Export</button>
              <button type="button" class="secondary-button" @click="fileInput.click()" :disabled="importing">{{ importing ? "Importing…" : "Import" }}</button>
            </div>
          </div>
          <p class="hint">Import replaces your current links.</p>
          <div class="setting-row">
            <label for="open-new-tab">Open links in a new tab</label>
            <button id="open-new-tab" type="button" role="switch" :aria-checked="store.state.openNewTab" @click="store.setOpenNewTab(!store.state.openNewTab)" class="switch" :class="{ enabled: store.state.openNewTab }"><span></span></button>
          </div>
          <p v-if="importStatus" role="status" class="text-sm">{{ importStatus }}</p>
          <input ref="fileInput" type="file" accept=".json" class="hidden" @change="importFile">
        </section>
      </div>

      <footer class="settings-footer">
        <a href="https://lobby-os.org" target="_blank" rel="noopener noreferrer">
          <svg class="footer-logo" width="100" height="26" viewBox="0 0 250 64" role="img" aria-label="LobbyOS">
            <g class="footer-logo-mark" fill="currentColor">
              <rect x="13" y="13" width="17" height="17" rx="4"/>
              <rect x="13" y="34" width="17" height="17" rx="4"/>
              <rect x="34" y="34" width="17" height="17" rx="4"/>
              <path d="m36.5 27.5 12-12m-10 0h10v10" fill="none" stroke="currentColor" stroke-width="4.5" stroke-linecap="round" stroke-linejoin="round"/>
            </g>
            <text x="76" y="46" fill="currentColor" font-size="40" letter-spacing="-1.5"><tspan font-weight="600">Lobby</tspan><tspan font-weight="400">OS</tspan></text>
          </svg>
        </a>
        <a href="https://github.com/exelban/lobby-os" target="_blank" rel="noopener noreferrer" class="github-link" aria-label="GitHub" title="GitHub">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path d="M12 .297C5.37.297 0 5.67 0 12.297c0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.043-1.61-4.043-1.61-.546-1.387-1.333-1.756-1.333-1.756-1.09-.745.083-.729.083-.729 1.205.084 1.838 1.237 1.838 1.237 1.07 1.835 2.809 1.305 3.495.998.108-.776.418-1.305.762-1.605-2.665-.305-5.467-1.334-5.467-5.93 0-1.31.469-2.381 1.236-3.221-.124-.303-.536-1.524.117-3.176 0 0 1.008-.322 3.301 1.23a11.52 11.52 0 0 1 3.003-.404c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.655 1.652.243 2.873.12 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222 0 1.606-.015 2.898-.015 3.293 0 .322.216.694.825.576C20.565 22.092 24 17.597 24 12.297c0-6.627-5.373-12-12-12"/>
          </svg>
        </a>
      </footer>
    </div>
  </Modal>
</template>

<script setup>
import Modal from "./modal.vue"
import store from "../utils/store.js"
import { backgrounds } from "../utils/backgrounds.js"
import {ref, computed} from "vue"

const close = () => { if (!importing.value) store.state.window.settings = false }
const theme = computed(() => store.state.theme)
const fileInput = ref(null)
const importing = ref(false)
const importStatus = ref("")
const linksPerRow = computed(() => store.state.linksPerRow)
const setLinksPerRow = (val) => {
  store.setLinksPerRow(val < 0 ? 0 : val)
}

const importFile = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  importing.value = true
  importStatus.value = ""
  try {
    await store.importLinks(JSON.parse(await file.text()))
    importStatus.value = "Links imported."
  } catch (error) {
    importStatus.value = error instanceof SyntaxError ? "This file is not valid JSON." : error.message
  } finally {
    importing.value = false
    e.target.value = ""
  }
}

const setTheme = store.setTheme

const download = () => {
  const array = store.exportLinks
  const elem = document.createElement("a")
  elem.setAttribute("href", `data:text/json;charset=utf-8,${encodeURIComponent(JSON.stringify(array))}`)
  elem.setAttribute("download", "links.json")
  document.body.appendChild(elem)
  elem.click()
  elem.remove()
}
</script>

<style scoped>
.settings-panel {
  width: 440px;
  max-width: 100%;
  max-height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border-radius: 22px;
  background: light-dark(#fafafa, #1b1e25);
  border: 1px solid light-dark(rgb(255 255 255 / 80%), rgb(255 255 255 / 10%));
  box-shadow: 0 24px 80px rgb(0 0 0 / 22%);
}
.settings-header, .settings-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
  padding: 18px 24px;
}
.settings-header { border-bottom: 1px solid light-dark(#e8e9ed, #30343c); }
.settings-header h2 { font-size: 20px; font-weight: 600; letter-spacing: -0.025em; outline: none; }
.settings-body { overflow-y: auto; overscroll-behavior: contain; padding: 0 24px; }
.settings-section { display: grid; gap: 16px; padding: 20px 0; }
.settings-section + .settings-section { border-top: 1px solid light-dark(#e8e9ed, #30343c); }
.settings-section h3 { font-size: 12px; font-weight: 600; color: light-dark(#737985, #969eac); letter-spacing: 0.04em; text-transform: uppercase; }
fieldset { min-width: 0; }
legend { margin-bottom: 10px; font-size: 14px; }
.segmented { display: flex; gap: 3px; padding: 4px; border-radius: 11px; background: light-dark(#eceef1, #11151c); }
.choice { flex: 1; min-width: 0; padding: 8px 4px; border-radius: 8px; text-align: center; font-size: 13px; font-weight: 500; cursor: pointer; color: light-dark(#707680, #9ca3af); }
.choice.selected { background: light-dark(#fff, #343a46); color: light-dark(#242a35, #fff); box-shadow: 0 1px 4px rgb(0 0 0 / 10%); }
.choice:has(input:focus-visible) { outline: 2px solid #3b82f6; outline-offset: 2px; }
.select-control { position: relative; width: 156px; max-width: 55%; flex-shrink: 0; }
.background-select {
  appearance: none;
  -webkit-appearance: none;
  display: block;
  width: 100%;
  min-height: 36px;
  padding: 8px 32px 8px 12px;
  border: 1px solid light-dark(#dde1e7, #3b414d);
  border-radius: 10px;
  background: light-dark(#f0f1f3, #242a35);
  color: inherit;
  font: inherit;
  font-size: 13px;
  line-height: 18px;
  cursor: pointer;
}
.background-select:focus { outline: none; }
.background-select:focus-visible { border-color: #3b82f6; outline: 2px solid light-dark(#3b82f640, #3b82f680); outline-offset: 2px; }
.select-chevron { position: absolute; right: 11px; top: 50%; transform: translateY(-50%); color: light-dark(#737985, #a0a7b3); pointer-events: none; }
.setting-row { display: flex; align-items: center; justify-content: space-between; gap: 12px; font-size: 14px; }
.setting-row label { cursor: pointer; }
.switch { width: 40px; height: 24px; padding: 3px; flex-shrink: 0; border-radius: 99px; background: light-dark(#c7ccd4, #4b5361); transition: background 150ms ease; }
.switch span { display: block; width: 18px; height: 18px; border-radius: 50%; background: #fff; box-shadow: 0 1px 3px rgb(0 0 0 / 15%); transition: transform 150ms ease; }
.switch.enabled { background: #3b82f6; }
.switch.enabled span { transform: translateX(16px); }
.stepper { display: flex; align-items: center; padding: 3px; border: 1px solid light-dark(#e0e3e8, #3b414d); border-radius: 10px; }
.stepper button { width: 30px; height: 30px; border-radius: 6px; font-size: 19px; }
.stepper output { min-width: 48px; text-align: center; font-size: 13px; font-variant-numeric: tabular-nums; }
.icon-button { display: grid; place-items: center; width: 32px; height: 32px; border-radius: 8px; color: light-dark(#808691, #a0a7b3); }
.secondary-button { padding: 7px 12px; border: 1px solid light-dark(#dde1e7, #3b414d); border-radius: 8px; font-size: 13px; font-weight: 500; }
button:disabled { opacity: 0.35; cursor: not-allowed; }
.hint { margin-top: -8px; font-size: 12px; color: light-dark(#808691, #969eac); }
.settings-footer { background: light-dark(#f0f1f3, #171a21); }
.settings-footer a { display: flex; align-items: center; gap: 8px; font-size: 12px; font-weight: 600; }
.footer-logo { flex-shrink: 0; color: light-dark(#202839, #eef2f9); }
.footer-logo-mark { color: light-dark(#335cff, #eef2f9); }
.github-link { padding: 10px; border-radius: 9px; color: light-dark(#282e39, #e2e6ed); }
@media (hover: hover) {
  .background-select:hover { background: light-dark(#e9ebef, #303641); border-color: light-dark(#c7ccd4, #505967); }
  .icon-button:hover, .stepper button:not(:disabled):hover, .secondary-button:not(:disabled):hover { background: light-dark(#e9ebef, #303641); }
  .choice:not(.selected):hover { color: light-dark(#242a35, #fff); }
  .github-link:hover { opacity: 0.85; }
}
@media (max-width: 380px) {
  .settings-header, .settings-footer { padding-inline: 16px; }
  .settings-body { padding-inline: 16px; }
}
</style>
