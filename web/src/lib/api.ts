import { headers } from "next/headers"

const API_PROXY_BASE = "/api"

export async function getServerApiBaseUrl() {
  const requestHeaders = await headers()
  const host = requestHeaders.get("x-forwarded-host") ?? requestHeaders.get("host")
  const proto = requestHeaders.get("x-forwarded-proto") ?? "http"

  if (!host) {
    return API_PROXY_BASE
  }

  return `${proto}://${host}${API_PROXY_BASE}`
}