// Read legacy cookies once, without losing preferences if storage is blocked.
export const readPreference = key => {
  try {
    const value = localStorage.getItem(key)
    if (value !== null) return value
  } catch {}
  if (!["theme", "linksPerRow"].includes(key)) return null
  try {
    const cookie = document.cookie.split("; ").find(item => item.startsWith(`${key}=`))
    if (!cookie) return null
    const value = decodeURIComponent(cookie.slice(key.length + 1))
    writePreference(key, value)
    return value
  } catch { return null }
}

export const writePreference = (key, value) => {
  try {
    localStorage.setItem(key, String(value))
    if (["theme", "linksPerRow"].includes(key)) document.cookie = `${key}=; Max-Age=0; Path=/`
  } catch {}
}
