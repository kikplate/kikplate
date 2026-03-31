import { NextRequest, NextResponse } from "next/server"

const API_INTERNAL_URL = process.env.API_INTERNAL_URL ?? "http://localhost:3001"

function buildTargetUrl(path: string[], request: NextRequest) {
  const target = new URL(`${API_INTERNAL_URL.replace(/\/$/, "")}/${path.join("/")}`)
  target.search = request.nextUrl.search
  return target
}

function buildForwardHeaders(request: NextRequest) {
  const headers = new Headers(request.headers)
  headers.delete("host")
  headers.delete("connection")
  headers.delete("content-length")
  headers.delete("transfer-encoding")
  return headers
}

async function proxy(request: NextRequest, context: { params: Promise<{ path: string[] }> }) {
  const { path } = await context.params
  const targetUrl = buildTargetUrl(path, request)

  const upstream = await fetch(targetUrl, {
    method: request.method,
    headers: buildForwardHeaders(request),
    body: request.method === "GET" || request.method === "HEAD" ? undefined : await request.arrayBuffer(),
    redirect: "manual",
    cache: "no-store",
  })

  const responseHeaders = new Headers(upstream.headers)
  responseHeaders.delete("content-encoding")
  responseHeaders.delete("content-length")
  responseHeaders.delete("transfer-encoding")

  return new NextResponse(await upstream.arrayBuffer(), {
    status: upstream.status,
    statusText: upstream.statusText,
    headers: responseHeaders,
  })
}

export { proxy as GET, proxy as POST, proxy as PATCH, proxy as PUT, proxy as DELETE, proxy as OPTIONS, proxy as HEAD }