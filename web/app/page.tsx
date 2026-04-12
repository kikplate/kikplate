import { HeroSearch } from "@/src/presentation/components/common/HeroSearch"
import { FeaturedPlates } from "@/src/presentation/components/plates/FeaturedPlates"
import { CategoriesGrid } from "@/src/presentation/components/common/CategoriesGrid"
import { StatsBanner } from "@/src/presentation/components/common/StatsBanner"
import type { AppConfig, PlateCategory } from "@/src/domain/entities/Config"
import { getServerApiBaseUrl } from "@/src/lib/api"

export default async function HomePage() {
  const base = await getServerApiBaseUrl()
  let plateCategories: PlateCategory[] = []
  try {
    const res = await fetch(`${base}/config`, { cache: "no-store" })
    if (res.ok) {
      const cfg = (await res.json()) as AppConfig
      plateCategories = cfg.plate_categories ?? []
    }
  } catch {
    plateCategories = []
  }

  return (
    <div>
      <HeroSearch />
      <FeaturedPlates />
      <CategoriesGrid categories={plateCategories} />
      <StatsBanner />
    </div>
  )
}
