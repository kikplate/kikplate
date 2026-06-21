"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Copy, Check, CheckCircle2, XCircle, Pencil, Trash2 } from "lucide-react"
import type { MeResult } from "@/src/domain/entities/User"
import { EditProfileModal } from "./EditProfileModal"
import { DeleteAccountModal } from "./DeleteAccountModal"
import { useLogout } from "@/src/presentation/hooks/useAuth"

function CopyButton({ value }: { value: string }) {
  const [copied, setCopied] = useState(false)
  async function handleCopy() {
    await navigator.clipboard.writeText(value)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }
  return (
    <button onClick={handleCopy} className="text-muted-foreground transition-colors hover:text-foreground">
      {copied ? <Check className="h-3.5 w-3.5 text-green-500" /> : <Copy className="h-3.5 w-3.5" />}
    </button>
  )
}





interface Row {
  label: string
  value: React.ReactNode
  copyable?: string
}

function oauthOrTrustedProvider(provider: string): boolean {
  return provider !== "local"
}

export function ProfileDetails({ me }: { me: MeResult }) {
  const [editOpen, setEditOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const router = useRouter()
  const logout = useLogout()

  const rows: Row[] = [
    {
      label: "Account ID",
      value: (
        <span className="max-w-xs truncate font-mono text-xs text-muted-foreground">
          {me.account_id}
        </span>
      ),
      copyable: me.account_id,
    },
    me.username
      ? {
          label: "Username",
          value: <span className="text-sm">{me.username}</span>,
        }
      : null,
    me.email
      ? {
          label: "Email",
          value: <span className="text-sm">{me.email}</span>,
        }
      : null,
    {
      label: "Provider",
      value: <span className="text-sm capitalize">{me.provider}</span>,
    },
    me.role
      ? {
          label: "Role",
          value: <span className="text-sm capitalize">{me.role}</span>,
        }
      : null,
    oauthOrTrustedProvider(me.provider) || me.is_active !== undefined
      ? {
          label: "Email verified",
          value:
            oauthOrTrustedProvider(me.provider) || me.is_active === true ? (
              <span className="flex items-center gap-1 text-sm text-green-600">
                <CheckCircle2 className="h-3.5 w-3.5" /> Verified
              </span>
            ) : (
              <span className="flex items-center gap-1 text-sm text-destructive">
                <XCircle className="h-3.5 w-3.5" /> Not verified
              </span>
            ),
        }
      : null,
  ].filter(Boolean) as Row[]

  return (
    <>
      <div className="max-w-lg space-y-6">
        <div>
          <p className="mb-4 text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Account details
          </p>
          <div className="divide-y divide-border border border-border">
            {rows.map((row) => (
              <div key={row.label} className="flex flex-col gap-2 bg-card px-4 py-3 sm:flex-row sm:items-center sm:justify-between">
                <span className="shrink-0 text-xs text-muted-foreground sm:w-32">{row.label}</span>
                <div className="flex min-w-0 items-center gap-2 self-start sm:self-auto">
                  {row.value}
                  {row.copyable && <CopyButton value={row.copyable} />}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
          <button
            onClick={() => setEditOpen(true)}
            className="flex h-8 w-full items-center justify-center gap-1.5 px-3 text-xs border border-border text-muted-foreground transition-colors hover:bg-muted hover:text-foreground sm:w-auto"
          >
            <Pencil className="h-3 w-3" />
            Edit profile
          </button>
          <button
            onClick={() => setDeleteOpen(true)}
            className="flex h-8 w-full items-center justify-center gap-1.5 px-3 text-xs border border-destructive/40 text-destructive/70 transition-colors hover:border-destructive hover:bg-destructive/5 hover:text-destructive sm:w-auto"
          >
            <Trash2 className="h-3 w-3" />
            Delete account
          </button>
        </div>
      </div>

      {editOpen && (
        <EditProfileModal
          me={me}
          onClose={() => setEditOpen(false)}
          onSaved={() => setEditOpen(false)}
        />
      )}
      {deleteOpen && (
        <DeleteAccountModal
          username={me.username ?? me.account_id}
          onClose={() => setDeleteOpen(false)}
          onDeleted={() => {
            setDeleteOpen(false)
            logout()
            router.push("/")
            router.refresh()
          }}
        />
      )}
    </>
  )
}