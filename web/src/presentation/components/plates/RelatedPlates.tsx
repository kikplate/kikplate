import Link from "next/link"
import { Heart, Star } from "lucide-react"
import type { Plate } from "@/src/domain/entities/Plate"
import { formatCount } from "@/src/presentation/utils/plateUtils"
import { PlateBadgeChips } from "@/src/presentation/components/plates/PlateBadgeChips"

export function RelatedPlates({
  plates,
  exploreHref,
  exploreLabel,
}: {
  plates: Plate[]
  exploreHref: string
  exploreLabel: string
}) {
  if (plates.length === 0) return null

  return (
    <div className="border border-border bg-card p-5">
      <p className="mb-3 text-xs font-semibold uppercase tracking-[0.14em] text-muted-foreground">
        Related plates
      </p>
      <ul className="space-y-3">
        {plates.map((p) => (
          <li key={p.id}>
            <Link
              href={`/plates/${encodeURIComponent(p.slug)}`}
              className="group block rounded-sm py-0.5 transition-colors hover:bg-muted/40"
            >
              <span className="line-clamp-2 text-sm font-semibold leading-snug text-foreground group-hover:text-primary">
                {p.name}
              </span>
              {p.description ? (
                <span className="mt-0.5 line-clamp-2 block text-xs text-muted-foreground">
                  {p.description}
                </span>
              ) : null}
              <PlateBadgeChips badges={p.badges} max={2} className="mt-1.5" />
              <span className="mt-1.5 flex flex-wrap items-center gap-3 text-[11px] text-muted-foreground">
                {p.avg_rating > 0 ? (
                  <span className="inline-flex items-center gap-1">
                    <Star className="h-3 w-3" />
                    {p.avg_rating.toFixed(1)}
                  </span>
                ) : null}
                <span className="inline-flex items-center gap-1">
                  <Heart className="h-3 w-3" />
                  {formatCount(p.bookmark_count)}
                </span>
                <span className="capitalize">{p.category}</span>
              </span>
            </Link>
          </li>
        ))}
      </ul>
      <Link
        href={exploreHref}
        className="mt-4 inline-block text-xs font-medium text-muted-foreground transition-colors hover:text-foreground"
      >
        {exploreLabel}
        {" →"}
      </Link>
    </div>
  )
}
