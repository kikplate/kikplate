import { NextResponse } from "next/server"

export const revalidate = 600

function parseRepo(spec: string): { owner: string; repo: string } | null {
  const parts = spec.split("/").filter(Boolean)
  if (parts.length !== 2) return null
  return { owner: parts[0], repo: parts[1] }
}

async function fetchGitHubStars(): Promise<number | null> {
  const spec = process.env.STATS_GITHUB_REPO?.trim() || "kikplate/kikplate"
  const parsed = parseRepo(spec)
  if (!parsed) return null

  const headers: Record<string, string> = {
    Accept: "application/vnd.github+json",
    "X-GitHub-Api-Version": "2022-11-28",
  }
  const pat = process.env.GITHUB_STATS_TOKEN?.trim()
  if (pat && pat !== "unset") {
    headers.Authorization = `Bearer ${pat}`
  }

  const res = await fetch(`https://api.github.com/repos/${parsed.owner}/${parsed.repo}`, {
    headers,
    next: { revalidate: 600 },
  })

  if (!res.ok) return null
  const data = (await res.json()) as { stargazers_count?: number }
  return typeof data.stargazers_count === "number" ? data.stargazers_count : null
}

export async function GET() {
  try {
    const githubStars = await fetchGitHubStars()
    return NextResponse.json({ githubStars })
  } catch (e) {
    console.error("external stats:", e)
    return NextResponse.json({ githubStars: null }, { status: 200 })
  }
}
