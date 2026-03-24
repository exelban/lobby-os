import Vuex from "vuex"

const randomColor = () => {
  let letters = "0123456789ABCDEF"
  let color = "#"
  for (let i = 0; i < 6; i++) { color += letters[Math.floor(Math.random() * 16)] }
  return color
}
const demo = [
  {
    id: "ubiquiti",
    name: "Ubiquiti",
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
    id: "jam",
    name: "Monitoring",
    url: "https://jam.serhiy.io",
    preset: "jam",
    group: "dev",
  },
  {
    id: "hass",
    name: "Home Assistant",
    url: "https://www.home-assistant.io",
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
  }
]

export default new Vuex.Store({
  state: {
    links: [],
    window: {
      settings: false,
      edit: false,
      delete: false,
    },
    editMode: false,
    editObject: null,
    deleteID: null,
    version: ""
  },
  mutations: {
    initialize(state, data) {
      let links = demo
      state.version = data.version

      if (state.version === "demo") links = demo
      else links = data.links

      state.links = links.reduce((acc, obj) => {
        const group = obj.group ?? ""
        if (!acc[group]) {
          acc[group] = []
        }
        acc[group].push(obj)
        return acc
      }, {})
    },
    setLinks(state, links) {
      state.links = links
    },
    setSettingsWindow(state, newValue) {
      state.window.settings = newValue
    },
    setEditWindow(state, newValue) {
      state.window.edit = newValue
    },
    setDeleteWindow(state, newValue) {
      state.window.delete = newValue
    },
    setEditMode(state, editMode) {
      state.editMode = editMode
    },
    setEditObject(state, editObject) {
      state.editObject = editObject
    },
    setDeleteID(state, id) {
      state.deleteID = id
    },
    addGroup(state) {
      state.links[Date.now()] = []
    }
  },
  actions: {
    async save({ state }, list) {
      if (state.version === "demo") return
      if (!list) list = state.links
      const array = Object.keys(list).flatMap(k =>
        list[k].map(l => ({ ...l, group: k }))
      )
      return fetch("/api", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(array)
      })
    },
    async upsert({ commit, state, dispatch }, link) {
      const shadow = {...state.links}
      const updated = Object.keys(shadow).some(k => {
        const index = shadow[k].findIndex(l => l.id === link.id)
        if (index !== -1) {
          shadow[k][index] = link
          return true
        }
        return false
      })
      if (!updated) {
        if (!shadow[""]) shadow[""] = []
        shadow[""].push(link)
      }
      return new Promise((resolve, reject) => {
        dispatch("save", shadow).then(() => {
          commit("setLinks", shadow)
          resolve()
        }).catch(() => {
          reject()
        })
      })
    },
    async delete({ commit, state, dispatch, getters }) {
      const shadow = {...state.links}
      Object.keys(shadow).filter(k => shadow[k].findIndex(l => l.id === state.deleteID) !== -1).forEach(k => {
        shadow[k] = shadow[k].filter(l => l.id !== state.deleteID)
      })
      Object.keys(shadow).filter(k => shadow[k].length === 0).forEach(k => {
        delete shadow[k]
      })
      return new Promise((resolve, reject) => {
        dispatch("save", shadow).then(() => {
          commit("setLinks", shadow)
          dispatch("closeDelete")
          resolve()
        }).catch(() => {
          dispatch("closeDelete")
          reject()
        })
      })
    },

    startEditMode({ commit }) {
      commit("setEditMode", true)
    },
    stopEditMode({ commit }) {
      commit("setLinks", Object.fromEntries(
        Object.entries(this.state.links).filter(([k, v]) => v.length > 0)
      ))
      commit("setEditMode", false)
    },

    openNewLink({ commit }) {
      commit("setEditWindow", true)
      commit("setEditObject", {
        color: randomColor()
      })
    },
    openEditLink({ commit }, obj) {
      commit("setEditWindow", true)
      commit("setEditObject", obj)
    },
    closeEditLink({ commit }) {
      commit("setEditWindow", false)
      commit("setEditObject", null)
    },

    openDelete({ commit }, id) {
      commit("setDeleteWindow", true)
      commit("setDeleteID", id)
    },
    closeDelete({ commit }) {
      commit("setDeleteWindow", false)
      commit("setDeleteID", null)
    },

    async importLinks({ commit, dispatch }, links) {
      const grouped = links.reduce((acc, obj) => {
        const group = obj.group ?? ""
        if (!acc[group]) acc[group] = []
        acc[group].push(obj)
        return acc
      }, {})
      return dispatch("save", grouped).then(() => {
        commit("setLinks", grouped)
      })
    },

    loadFaviconURL(_, url) {
      return fetch("/api/favicon?url=" + encodeURIComponent(url)).then(res => res.text())
    }
  }
})