import { repoToPath } from "@/src/data/repositories/githubClient"

function encodeRepoFilePath(path: string): string {
  return path
    .split("/")
    .filter(Boolean)
    .map((s) => encodeURIComponent(s))
    .join("/")
}

export function resolveRepoMarkdownHref(
  href: string | undefined,
  repoUrl: string | undefined,
  branch: string | undefined,
  sourceFileInRepo: string
): string {
  if (href == null || href === "") return "#"
  const trimmed = href.trim()
  if (trimmed.startsWith("#")) return trimmed
  if (/^[a-z][a-z0-9+.-]*:/i.test(trimmed)) return trimmed
  if (trimmed.startsWith("//")) return `https:${trimmed}`
  if (!repoUrl || !branch) return trimmed

  const repoPath = repoToPath(repoUrl)
  const branchEnc = encodeURIComponent(branch)

  if (trimmed.startsWith("/") && !trimmed.startsWith("//")) {
    const rel = trimmed.replace(/^\/+/, "")
    if (!rel) return trimmed
    return `https://github.com/${repoPath}/blob/${branchEnc}/${encodeRepoFilePath(rel)}`
  }

  const base = `https://github.com/${repoPath}/blob/${branchEnc}/${sourceFileInRepo}`
  try {
    return new URL(trimmed, base).href
  } catch {
    return trimmed
  }
}

