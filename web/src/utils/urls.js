export const normalizeURL = value => {
  let text = String(value ?? "").trim()
  if (!text || /\s/.test(text)) throw new Error("Enter a valid address, such as nas.local:5000.")
  const explicitScheme = /^[a-z][a-z\d+.-]*:\/\//i.test(text)
  if (!explicitScheme && !text.startsWith("//") && /^[a-z][a-z\d+.-]*:/i.test(text) && !/^[^/:]+:\d+(?:[/?#]|$)/.test(text)) {
    throw new Error("Use an HTTP or HTTPS address.")
  }
  if (!explicitScheme) text = "https://" + text.replace(/^\/\//, "")
  let url
  try { url = new URL(text) } catch { throw new Error("Enter a valid address, such as nas.local:5000.") }
  if (!["http:", "https:"].includes(url.protocol) || !url.hostname) throw new Error("Use an HTTP or HTTPS address.")
  if (!explicitScheme) {
    const host = url.hostname
    const local = !host.includes(".") || /\.(local|lan|home|localhost)$/.test(host) || /^\d+\.\d+\.\d+\.\d+$/.test(host) || host.startsWith("[")
    if ((local && !/:443(?:[/?#]|$)/.test(text)) || (url.port && url.port !== "443")) url.protocol = "http:"
  }
  return url.href
}
