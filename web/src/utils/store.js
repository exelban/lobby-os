import { reactive } from "vue"
import { readPreference, writePreference } from "./preferences.js"
import { normalizeURL } from "./urls.js"
import { normalizeBackground } from "./backgrounds.js"

const cardSizes = ["small", "medium", "large"]

export const createGroupID = () => `group-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 9)}`

const groupLinks = links => links.reduce((groups, link) => {
  const group = link.group ?? ""
  if (!groups[group]) groups[group] = []
  groups[group].push(link)
  return groups
}, Object.create(null))

const groupNames = links => Object.fromEntries(Object.entries(links).map(([id, items]) =>
  [id, items[0]?.groupName ?? ""]
))
const flattenLinks = (links, names) => Object.entries(links).flatMap(([group, items]) =>
  items.map(link => ({...link, group, groupName: names[group] ?? ""}))
)

const randomColor = () => {
  let letters = "0123456789ABCDEF"
  let color = "#"
  for (let i = 0; i < 6; i++) { color += letters[Math.floor(Math.random() * 16)] }
  return color
}
const demo = [
  {
    id: "ubiquiti",
    name: "UniFi",
    url: "https://ui.com",
    preset: "unifi"
  },
  {
    id: "proxmox",
    name: "Proxmox",
    url: "https://www.proxmox.com/en/",
    preset: "proxmox"
  },
  {
    id: "EndPoll",
    name: "EndPoll",
    url: "https://endpoll.com",
    preset: "endpoll",
    group: "dev",
    groupName: "",
  },
  {
    id: "hass",
    name: "Home Assistant",
    url: "https://demo.home-assistant.io/",
    preset: "hass"
  },
  {
    id: "jellyfin",
    name: "Jellyfin",
    url: "https://jellyfin.org",
    preset: "jellyfin"
  },
  {
    id: "gitea",
    name: "Gitea",
    url: "https://gitea.io",
    preset: "gitea",
    group: "dev",
    groupName: "",
  }
]

const state = reactive({
  links: {},
  groupNames: {},
  window: {
    settings: false,
    edit: false,
    delete: false,
  },
  editMode: false,
  editObject: null,
  deleteID: null,
  version: "",
  openNewTab: readPreference("openNewTab") !== "false",
  fullCardColor: readPreference("fullCardColor") !== "false",
  cardSize: cardSizes.includes(readPreference("cardSize")) ? readPreference("cardSize") : "medium",
  linksPerRow: Math.max(0, parseInt(readPreference("linksPerRow")) || 0),

  theme: ["light", "dark"].includes(readPreference("theme")) ? readPreference("theme") : "system",
  background: normalizeBackground(readPreference("background")),
})

const store = {
  state,
  get exportLinks() { return flattenLinks(state.links, state.groupNames) },
  get groupOptions() {
    return Object.keys(state.links).map((id, index) => ({id, label: state.groupNames[id] || `Group ${index + 1}`}))
  },
  setTheme(value) {
    state.theme = ["light", "dark"].includes(value) ? value : "system"
    document.documentElement.classList.remove("light", "dark")
    if (state.theme !== "system") document.documentElement.classList.add(state.theme)
    writePreference("theme", state.theme)
  },
  setBackground(value) {
    state.background = normalizeBackground(value)
    document.body.dataset.background = state.background
    writePreference("background", state.background)
  },
  initialize(data) {
    state.version = data.version
    state.links = groupLinks(state.version === "demo" ? demo : (data.links ?? []))
    state.groupNames = groupNames(state.links)
  },
  addGroup() {
    const id = createGroupID()
    state.links[id] = []
    state.groupNames[id] = ""
  },
  setOpenNewTab(value) {
    state.openNewTab = Boolean(value)
    writePreference("openNewTab", state.openNewTab)
  },
  setFullCardColor(value) {
    state.fullCardColor = Boolean(value)
    writePreference("fullCardColor", state.fullCardColor)
  },
  setCardSize(value) {
    if (!cardSizes.includes(value)) return
    state.cardSize = value
    // Keep the choice usable for this session if browser storage is disabled.
    writePreference("cardSize", value)
  },
  setLinksPerRow(value) {
    value = parseInt(value) || 0
    if (value < 0) value = 0
    state.linksPerRow = value
    writePreference("linksPerRow", value)
  },

  async save({links = state.links, names = state.groupNames} = {}) {
    if (state.version === "demo") return
    const array = flattenLinks(links, names)
    const response = await fetch("/api", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(array)
    })
    if (!response.ok) throw new Error("Could not save links. Please try again.")
  },
  async upsert(link) {
    link = {...link, name: link.name.trim(), url: normalizeURL(link.url), group: link.group ?? ""}
    if (!link.name) throw new Error("Enter a link name.")
    const names = {...state.groupNames}
    if (!Object.hasOwn(state.links, link.group)) names[link.group] = (link.groupName ?? "").trim()
    const shadow = Object.assign(Object.create(null), state.links)
    let oldGroup
    let oldIndex = -1
    for (const group of Object.keys(shadow)) {
      const index = shadow[group].findIndex(item => item.id === link.id)
      if (index !== -1) { oldGroup = group; oldIndex = index; break }
    }
    if (oldGroup === link.group) {
      shadow[oldGroup] = [...shadow[oldGroup]]
      shadow[oldGroup][oldIndex] = link
    } else {
      if (oldGroup !== undefined) shadow[oldGroup] = shadow[oldGroup].filter(item => item.id !== link.id)
      shadow[link.group] = [...(shadow[link.group] ?? []), link]
    }
    for (const group of Object.keys(shadow)) if (!shadow[group].length) delete shadow[group]
    await store.save({links: shadow, names})
    state.links = shadow
    state.groupNames = names
  },
  async renameGroup({group, name}) {
    name = name.trim()
    if (name === state.groupNames[group]) return
    const names = {...state.groupNames, [group]: name}
    await store.save({names})
    state.groupNames = names
  },

  async deleteLink() {
    const shadow = {...state.links}
    Object.keys(shadow).filter(k => shadow[k].findIndex(l => l.id === state.deleteID) !== -1).forEach(k => {
      shadow[k] = shadow[k].filter(l => l.id !== state.deleteID)
    })
    Object.keys(shadow).filter(k => shadow[k].length === 0).forEach(k => {
      delete shadow[k]
    })
    await store.save({links: shadow})
    state.links = shadow
    store.closeDelete()
  },

  startEditMode() {
    state.editMode = true
  },
  stopEditMode() {
    state.links = Object.fromEntries(
      Object.entries(state.links).filter(([k, v]) => v.length > 0)
    )
    state.editMode = false
  },

  openNewLink() {
    state.window.edit = true
    state.editObject = {
      color: randomColor(),
      group: Object.keys(state.links)[0] ?? ""
    }
  },
  openEditLink(obj) {
    state.window.edit = true
    state.editObject = {...obj, group: Object.keys(state.links).find(group => state.links[group].some(link => link.id === obj.id)) ?? ""}
  },
  closeEditLink() {
    state.window.edit = false
    state.editObject = null
  },

  openDelete(id) {
    state.window.delete = true
    state.deleteID = id
  },
  closeDelete() {
    state.window.delete = false
    state.deleteID = null
  },

  async importLinks(links) {
    if (!Array.isArray(links) || links.some(link =>
      !link || typeof link !== "object" ||
      typeof link.name !== "string" || !link.name.trim() ||
      typeof link.url !== "string" || !link.url.trim() ||
      ["id", "group", "groupName", "color", "icon", "preset"].some(key => link[key] != null && typeof link[key] !== "string")
    )) throw new Error("Choose a JSON array of links with a name and URL.")
    const ids = new Set()
    const grouped = groupLinks(links.map(link => {
      let id = link.id
      if (!id || ids.has(id)) id = Math.random().toString(36).slice(2) + Date.now().toString(36)
      ids.add(id)
      return {...link, id, url: normalizeURL(link.url)}
    }))
    const names = groupNames(grouped)
    await store.save({links: grouped, names})
    state.links = grouped
    state.groupNames = names
  },

  async loadFaviconURL(url) {
    const response = await fetch("/api/favicon?url=" + encodeURIComponent(normalizeURL(url)))
    if (!response.ok) throw new Error("Could not load the favicon. Please try again.")
    return response.text()
  }

}

export default store
