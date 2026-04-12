"use client"

import { useRouter } from "next/navigation"
import { getPlateCategoryIcon } from "@/src/presentation/utils/plateCategoryIcons"
import type { PlateCategory } from "@/src/domain/entities/Config"

interface Props {
  categories: PlateCategory[]
}

export function CategoriesGrid({ categories }: Props) {
  const router = useRouter()

  if (categories.length === 0) {
    return null
  }

  return (
    <section className="bg-background py-20">
      <div className="container mx-auto px-4">

        <div className="mb-10 flex flex-wrap items-end justify-between gap-4">
          <div>
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Taxonomy
            </p>
            <h2 className="mt-2 text-3xl font-bold tracking-tight text-foreground">Browse by category</h2>
            <p className="mt-2 max-w-md text-sm leading-relaxed text-muted-foreground">
              Find templates by domain, from backend services to AI pipelines and infrastructure.
            </p>
          </div>
          <p className="text-xs tabular-nums text-muted-foreground border border-border bg-card px-3 py-1.5">{categories.length} domains</p>
        </div>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
          {categories.map(({ slug, label, description, icon }) => {
            const Icon = getPlateCategoryIcon(icon)
            return (
              <button
                key={slug}
                onClick={() => router.push(`/explore?category=${encodeURIComponent(slug)}`)}
                className="group flex items-start gap-3 border border-border bg-card p-4 text-left transition-all hover:border-foreground/20 hover:bg-background hover:-translate-y-0.5"
              >
                <div className="mt-0.5 shrink-0 text-muted-foreground transition-colors group-hover:text-foreground">
                  <Icon className="h-4 w-4" />
                </div>
                <div className="min-w-0">
                  <p className="text-sm font-semibold text-foreground">{label}</p>
                  <p className="mt-0.5 text-xs text-muted-foreground">{description}</p>
                </div>
              </button>
            )
          })}
        </div>

      </div>
    </section>
  )
}
