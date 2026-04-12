"use client"

import { useState } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { useSubmitRepository } from "@/src/presentation/hooks/usePlates"
import { useMyOrganizations } from "@/src/presentation/hooks/useOrganizations"
import { useMe } from "@/src/presentation/hooks/useAuth"
import { useConfig } from "@/src/presentation/hooks/useConfig"
import { Loader2, AlertCircle, Copy, Check, FileCode2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"

export function SubmitRepositoryForm() {
  const router = useRouter()
  const submit = useSubmitRepository()
  const { data: me } = useMe()
  const { data: organizations } = useMyOrganizations()
  const { data: appConfig, isLoading: configLoading } = useConfig()
  const plateCategories = appConfig?.plate_categories ?? []
  const [repoUrl, setRepoUrl] = useState("")
  const [branch, setBranch] = useState("main")
  const [organizationId, setOrganizationId] = useState("")
  const selectedOrganization = organizations?.find((org) => org.id === organizationId)
  const isPersonalSubmission = !organizationId
  const ownerHint = selectedOrganization?.name ?? me?.username ?? "your-username"
  const [ownerCopied, setOwnerCopied] = useState(false)
  const kikplateFileUrl = repoUrl
    ? `${repoUrl.replace(/\.git$/, "").replace(/\/$/, "")}/blob/${branch || "main"}/kikplate.yaml`
    : null

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    try {
      await submit.mutateAsync({
        repo_url: repoUrl,
        branch,
        organization_id: organizationId || undefined,
      })
      toast.success("Plate submitted. Complete verification from your account.")
      router.replace("/account?tab=plates")
      router.refresh()
    } catch {
    }
  }

  const errorMsg = submit.error instanceof Error ? submit.error.message : null

  return (
    <form
      onSubmit={handleSubmit}
      className="grid grid-cols-1 gap-8 lg:grid-cols-2 lg:gap-10 lg:items-start"
    >
      <div className="border border-border bg-card lg:sticky lg:top-8 self-start">
        <div className="border-b border-border bg-muted/20 px-4 py-3 flex items-start gap-2.5">
          <FileCode2 className="h-4 w-4 text-muted-foreground shrink-0 mt-0.5" />
          <div>
            <h2 className="text-sm font-semibold text-foreground">Before you submit</h2>
            <p className="mt-0.5 text-xs text-muted-foreground leading-relaxed">
              Add <code className="font-mono bg-muted px-1 py-0.5">kikplate.yaml</code> at the repository root.{" "}
              <Link
                href="/docs?doc=how-it-works"
                className="font-medium text-foreground underline underline-offset-2 hover:text-muted-foreground"
              >
                How it works
              </Link>
            </p>
          </div>
        </div>

        <div className="px-4 py-4 space-y-5 text-xs text-muted-foreground leading-relaxed">
          <div>
            <p className="font-semibold text-foreground mb-1.5">Repository</p>
            <ul className="list-disc pl-4 space-y-1 mb-3">
              <li>The repository must be public.</li>
              <li>
                It must include <code className="font-mono bg-muted px-1 py-0.5">kikplate.yaml</code> at the root on
                the branch you enter below.
              </li>
            </ul>

            <p className="font-semibold text-foreground mb-2 text-[11px] uppercase tracking-wide">
              Personal or organization
            </p>
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="rounded border border-border/80 bg-muted/10 px-3 py-2.5 space-y-1">
                <p className="text-xs font-semibold text-foreground">Personal</p>
                <p className="text-[11px] text-muted-foreground leading-relaxed">
                  Under Repository details, leave <span className="text-foreground font-medium">Personal</span>. The
                  plate is registered on <span className="text-foreground font-medium">your user account</span> only.
                  In <code className="font-mono bg-muted px-1 py-0.5">kikplate.yaml</code>, set{" "}
                  <code className="font-mono bg-muted px-1 py-0.5">owner</code> to your Kikplate username
                  {me?.username ? (
                    <> (<span className="font-mono text-foreground">{me.username}</span>)</>
                  ) : (
                    <> (the username on your profile)</>
                  )}.
                </p>
              </div>
              <div className="rounded border border-border/80 bg-muted/10 px-3 py-2.5 space-y-1">
                <p className="text-xs font-semibold text-foreground">Organization</p>
                <p className="text-[11px] text-muted-foreground leading-relaxed">
                  Under Repository details, choose an <span className="text-foreground font-medium">organization you own</span>.
                  The plate is registered <span className="text-foreground font-medium">under that organization</span> (not
                  only on your user). In <code className="font-mono bg-muted px-1 py-0.5">kikplate.yaml</code>, set{" "}
                  <code className="font-mono bg-muted px-1 py-0.5">owner</code> to the{" "}
                  <span className="text-foreground font-medium">organization&apos;s exact name</span>—not your personal
                  username.
                </p>
              </div>
            </div>

            <p className="mt-3 rounded border border-border/70 bg-muted/15 px-3 py-2 text-[11px] leading-relaxed">
              <span className="font-medium text-foreground">For this submit: </span>
              {isPersonalSubmission ? (
                <>
                  You have <span className="text-foreground font-medium">Personal</span> selected →{" "}
                  <code className="font-mono bg-muted px-1 py-0.5">owner</code> must be{" "}
                  <span className="font-mono text-foreground font-medium">{ownerHint}</span>.
                </>
              ) : (
                <>
                  You have <span className="text-foreground font-medium">{selectedOrganization?.name}</span> selected →
                  the plate goes to that org → <code className="font-mono bg-muted px-1 py-0.5">owner</code> must be{" "}
                  <span className="font-mono text-foreground font-medium">{ownerHint}</span> (org name, not your username).
                </>
              )}
            </p>
          </div>

          <div>
            <p className="font-semibold text-foreground mb-1.5">Manifest fields</p>
            <ul className="list-disc pl-4 space-y-1">
              <li>
                <span className="text-foreground">Required:</span>{" "}
                <code className="font-mono bg-muted px-1 py-0.5">name</code> (plate title, used for the slug),{" "}
                <code className="font-mono bg-muted px-1 py-0.5">owner</code> (must match the value above).
              </li>
              <li>
                <span className="text-foreground">Optional:</span>{" "}
                <code className="font-mono bg-muted px-1 py-0.5">description</code>,{" "}
                <code className="font-mono bg-muted px-1 py-0.5">category</code> (must be one of the slugs below;
                matching is case-insensitive; anything else or empty becomes{" "}
                <code className="font-mono bg-muted px-1 py-0.5">other</code>),{" "}
                <code className="font-mono bg-muted px-1 py-0.5">tags</code>.
              </li>
              <li>
                <span className="text-foreground">After submit:</span> add{" "}
                <code className="font-mono bg-muted px-1 py-0.5">verification_token</code> from your account, push,
                then verify so the plate can go public. Sync keeps re-reading this file.
              </li>
            </ul>

            <div className="mt-3 rounded border border-border/70 bg-muted/10 px-3 py-2.5">
              <p className="text-[11px] font-semibold text-foreground mb-2">Allowed category slugs</p>
              {configLoading ? (
                <p className="text-[11px] text-muted-foreground italic">Loading…</p>
              ) : plateCategories.length === 0 ? (
                <p className="text-[11px] text-muted-foreground italic">None returned.</p>
              ) : (
                <ul className="flex flex-wrap gap-1.5">
                  {plateCategories.map((c) => (
                    <li
                      key={c.slug}
                      title={c.label ? `${c.label}: ${c.description}` : c.description}
                      className="border border-border bg-background px-2 py-0.5 font-mono text-[11px] text-foreground"
                    >
                      {c.slug}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>

          <div>
            <p className="font-semibold text-foreground mb-1.5">Minimal example</p>
            <pre className="border border-border bg-muted/20 p-3 font-mono text-[11px] text-foreground overflow-x-auto leading-snug">
{`name: My Starter
owner: your-username
description: Short summary
category: backend
tags:
  - golang
  - docker`}
            </pre>
          </div>
        </div>
      </div>

      <div className="min-w-0 space-y-6">
        <div className="space-y-4">
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
          Repository details
        </p>

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-foreground">
            GitHub URL <span className="text-destructive">*</span>
          </label>
          <input
            required
            type="url"
            placeholder="https://github.com/username/repo"
            value={repoUrl}
            onChange={(e) => setRepoUrl(e.target.value)}
            className="w-full border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors"
          />
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-foreground">Branch</label>
          <input
            type="text"
            placeholder="main"
            value={branch}
            onChange={(e) => setBranch(e.target.value)}
            className="w-full border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors"
          />
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-foreground">Organization</label>
          <select
            value={organizationId}
            onChange={(e) => setOrganizationId(e.target.value)}
            className="w-full border border-input bg-background px-3 py-2 text-sm outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors"
          >
            <option value="">Personal (no organization)</option>
            {organizations?.map((org) => (
              <option key={org.id} value={org.id}>{org.name}</option>
            ))}
          </select>
          {!organizations || organizations.length === 0 ? (
            <p className="text-xs text-muted-foreground">
              No organizations found. You can still submit personally using your username as owner.
            </p>
          ) : (
            <p className="text-xs text-muted-foreground">
              Choose Personal to submit under your account, or pick an organization.
            </p>
          )}
        </div>
        </div>

      {errorMsg && (
        <div className="border border-destructive/40 bg-destructive/5 px-4 py-3 space-y-2">
          <div className="flex items-start gap-2 text-sm text-destructive">
            <AlertCircle className="h-4 w-4 shrink-0 mt-0.5" />
            <span className="font-medium">{errorMsg}</span>
          </div>
          {errorMsg.includes("owner field") && (
            <div className="pl-6 space-y-2 text-xs text-muted-foreground">
              <p>
                Open{" "}
                {kikplateFileUrl ? (
                  <a
                    href={kikplateFileUrl}
                    target="_blank"
                    rel="noreferrer"
                    className="font-mono underline underline-offset-2 hover:text-foreground"
                  >
                    kikplate.yaml
                  </a>
                ) : (
                  <code className="font-mono bg-muted px-1 py-0.5">kikplate.yaml</code>
                )}{" "}
                and update the <code className="font-mono bg-muted px-1 py-0.5">owner</code> field, then push and try again.
              </p>
              <div className="flex items-center gap-0 border border-border">
                <code className="flex-1 truncate bg-muted/20 px-3 py-2.5 text-xs font-mono text-foreground">
                  owner: {ownerHint}
                </code>
                <button
                  type="button"
                  onClick={async () => {
                    await navigator.clipboard.writeText(`owner: ${ownerHint}`)
                    setOwnerCopied(true)
                    toast.success("Copied to clipboard")
                    setTimeout(() => setOwnerCopied(false), 2000)
                  }}
                  className="border-l border-border px-3 py-2.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                >
                  {ownerCopied
                    ? <Check className="h-3.5 w-3.5 text-green-500" />
                    : <Copy className="h-3.5 w-3.5" />}
                </button>
              </div>
            </div>
          )}
          {errorMsg.includes("Account settings to set one") && (
            <p className="text-xs text-muted-foreground pl-6">
              Go to <strong>Account → Profile</strong> and set a username before submitting a plate.
            </p>
          )}
          {(errorMsg.includes("not found") || errorMsg.includes("fetch")) && (
            <p className="text-xs text-muted-foreground pl-6">
              Make sure the repository is public, the URL is correct, and
              the <code className="font-mono bg-muted px-1 py-0.5">kikplate.yaml</code> exists
              on the <code className="font-mono bg-muted px-1 py-0.5">{branch}</code> branch.
            </p>
          )}
          {errorMsg.includes("conflict") && (
            <p className="text-xs text-muted-foreground pl-6">
              A plate with this name already exists. Rename your plate in{" "}
              <code className="font-mono bg-muted px-1 py-0.5">kikplate.yaml</code> and try again.
            </p>
          )}
        </div>
      )}

        <div className="border-t border-border pt-5">
          <Button
            type="submit"
            disabled={submit.isPending || !repoUrl}
            className="gap-2"
          >
            {submit.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            Submit repository plate
          </Button>
        </div>
      </div>
    </form>
  )
}
