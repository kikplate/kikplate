"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Heart } from "lucide-react"
import { useSetBookmark } from "@/src/presentation/hooks/usePlates"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { toast } from "sonner"
import { ApiError } from "@/src/data/repositories/httpClient"
import { AuthService } from "@/src/domain/services/AuthService"

interface Props {
  plateId: string
  isBookmarked?: boolean
  prominent?: boolean
  className?: string
}

export function BookmarkButtonClient({ plateId, isBookmarked = false, prominent = false, className }: Props) {
  const router = useRouter()
  const [bookmarked, setBookmarked] = useState(isBookmarked)
  const mutation = useSetBookmark()

  async function handleBookmarkToggle() {
    if (!AuthService.isAuthenticated()) {
      toast.error("Sign in to bookmark", {
        description: "You need an account to save plates.",
        action: { label: "Log in", onClick: () => router.push("/login") },
      })
      return
    }
    try {
      await mutation.mutateAsync({ id: plateId, bookmarked: !bookmarked })
      setBookmarked(!bookmarked)
      toast.success(bookmarked ? "Removed from bookmarks" : "Added to bookmarks")
    } catch (e) {
      if (e instanceof ApiError && e.isUnauthorized()) {
        toast.error("Sign in required", {
          description: "Your session expired or is invalid. Sign in again to bookmark.",
          action: { label: "Log in", onClick: () => router.push("/login") },
        })
        return
      }
      if (e instanceof ApiError && e.isForbidden()) {
        toast.error("Cannot update bookmark", { description: e.message })
        return
      }
      toast.error("Failed to update bookmark", {
        description: e instanceof Error ? e.message : undefined,
      })
    }
  }

  return (
    <Button
      onClick={handleBookmarkToggle}
      disabled={mutation.isPending}
      variant={bookmarked ? "default" : "outline"}
      className={cn(
        prominent
          ? "h-11 w-full gap-2 border-primary/90 bg-primary text-primary-foreground shadow-lg hover:bg-primary/90"
          : "h-9 gap-2 border-border/80 bg-background px-3.5 text-foreground hover:border-primary/40 hover:bg-muted",
        className,
      )}
    >
      <Heart className={cn("h-3.5 w-3.5", bookmarked && "fill-current")} />
      <span className="text-sm font-semibold">{bookmarked ? "Bookmarked" : "Bookmark"}</span>
    </Button>
  )
}
