import type { PlateBadge } from "@/src/domain/entities/Plate"
import { cn } from "@/lib/utils"
import { getBadgeIcon } from "@/src/presentation/utils/badgeIcons"
import { tierColour } from "@/src/presentation/utils/plateUtils"

interface Props {
  badges?: PlateBadge[]
  max?: number
  className?: string
}

export function PlateBadgeChips({ badges, max = 3, className = "" }: Props) {
  const withBadge = badges?.filter((pb) => pb.badge) ?? []
  if (withBadge.length === 0) return null

  const shown = withBadge.slice(0, max)

  return (
    <div className={cn("flex flex-wrap gap-1", className)}>
      {shown.map((pb) => {
        const b = pb.badge!
        const Icon = getBadgeIcon(b.slug)
        return (
          <span
            key={pb.id}
            className={`inline-flex items-center gap-1 px-1.5 py-0.5 text-[10px] border font-medium ${tierColour(b.tier)}`}
          >
            <Icon className="h-2.5 w-2.5 shrink-0" />
            {b.name}
          </span>
        )
      })}
      {withBadge.length > max ? (
        <span className="self-center text-[10px] text-muted-foreground/60 tabular-nums">
          +{withBadge.length - max}
        </span>
      ) : null}
    </div>
  )
}
