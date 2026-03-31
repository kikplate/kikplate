"use client"

import Link from "next/link"
import { useEffect, useState } from "react"
import {
  BarChart, Bar, LineChart, Line, PieChart, Pie, Cell,
  XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
} from "recharts"
import { BarChart3, ShieldCheck } from "lucide-react"
import { useStats } from "@/src/presentation/hooks/usePlates"
import { useBadges } from "@/src/presentation/hooks/useBadges"
import { buildPlateHref, formatCount } from "@/src/presentation/utils/plateUtils"

const BASE = "/api"

/* ─── resolve CSS custom-property values for recharts (SVG attrs) ─── */
function useChartColors() {
  const [colors, setColors] = useState({
    chart1: "#3b82f6", chart2: "#14b8a6", chart3: "#f59e0b",
    chart4: "#a855f7", chart5: "#06b6d4", fg: "#888",
    card: "#fff", border: "#e5e5e5",
  })
  useEffect(() => {
    function resolve() {
      const s = getComputedStyle(document.documentElement)
      const get = (v: string) => s.getPropertyValue(v).trim() || undefined
      setColors({
        chart1: get("--chart-1") ?? "#3b82f6",
        chart2: get("--chart-2") ?? "#14b8a6",
        chart3: get("--chart-3") ?? "#f59e0b",
        chart4: get("--chart-4") ?? "#a855f7",
        chart5: get("--chart-5") ?? "#06b6d4",
        fg: get("--foreground") ?? "#888",
        card: get("--card") ?? "#fff",
        border: get("--border") ?? "#e5e5e5",
      })
    }
    resolve()
    // re-resolve when theme toggles (class change on <html>)
    const observer = new MutationObserver(resolve)
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ["class"] })
    return () => observer.disconnect()
  }, [])
  return colors
}

/* ─── types ─── */
interface MonthlyCount { month: string; count: number }
interface CategoryCount { category: string; count: number }
interface PlateRanked {
  id: string; slug: string; name: string
  bookmark_count: number; avg_rating: number; category: string
}

/* ─── independent fetcher hook ─── */
function useFetch<T>(path: string) {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  useEffect(() => {
    let cancelled = false
    fetch(`${BASE}${path}`)
      .then((r) => (r.ok ? r.json() : null))
      .then((d) => { if (!cancelled) setData(d) })
      .catch(() => {})
      .finally(() => { if (!cancelled) setLoading(false) })
    return () => { cancelled = true }
  }, [path])
  return { data, loading }
}

/* ─── skeleton loader ─── */
function PanelLoader() {
  return (
    <div className="flex items-center justify-center py-20">
      <div className="h-5 w-5 animate-spin rounded-full border-2 border-muted-foreground border-t-transparent" />
    </div>
  )
}

/* ─── panel wrapper ─── */
function Panel({ title, subtitle, children }: { title: string; subtitle: string; children: React.ReactNode }) {
  return (
    <section className="border border-border bg-card p-5 sm:p-6">
      <div className="mb-5">
        <h2 className="text-xl font-bold text-foreground">{title}</h2>
        <p className="mt-1 text-sm text-muted-foreground">{subtitle}</p>
      </div>
      {children}
    </section>
  )
}

/* ─── KPI row ─── */
function KpiPanel() {
  const { data: stats, isLoading } = useStats()
  const kpis = [
    { label: "Templates", value: stats?.total_plates ?? 0 },
    { label: "Contributors", value: stats?.total_contributors ?? 0 },
    { label: "Categories", value: stats?.total_categories ?? 0 },
    { label: "Bookmarks", value: stats?.total_bookmarks ?? 0 },
  ]
  return (
    <div className="grid grid-cols-2 gap-3 sm:grid-cols-4 sm:gap-4">
      {kpis.map((k) => (
        <div key={k.label} className="border border-border bg-card p-4 sm:p-5">
          <p className="text-xs uppercase tracking-widest text-muted-foreground">{k.label}</p>
          <p className="mt-2 text-3xl font-black tabular-nums text-foreground">
            {isLoading ? <span className="inline-block h-5 w-5 animate-spin rounded-full border-2 border-muted-foreground border-t-transparent" /> : formatCount(k.value)}
          </p>
        </div>
      ))}
    </div>
  )
}

/* ─── Monthly growth line chart ─── */
function GrowthPanel() {
  const c = useChartColors()
  const { data, loading } = useFetch<MonthlyCount[]>("/plates/stats/growth?months=12")
  const rows = (data ?? []).map((r) => ({
    ...r,
    label: r.month.slice(5),
  }))
  return (
    <Panel title="New Plates Added Monthly" subtitle="Monthly new templates over the last 12 months">
      {loading ? <PanelLoader /> : rows.length === 0 ? (
        <p className="text-sm text-muted-foreground">No data yet.</p>
      ) : (
        <ResponsiveContainer width="100%" height={280}>
          <BarChart data={rows} margin={{ top: 8, right: 8, bottom: 0, left: -12 }}>
            <CartesianGrid strokeDasharray="3 3" stroke={c.fg} strokeOpacity={0.1} />
            <XAxis dataKey="label" tick={{ fontSize: 11, fill: c.fg, fillOpacity: 0.6 }} stroke={c.fg} strokeOpacity={0.1} />
            <YAxis allowDecimals={false} tick={{ fontSize: 11, fill: c.fg, fillOpacity: 0.6 }} stroke={c.fg} strokeOpacity={0.1} />
            <Tooltip
              contentStyle={{ backgroundColor: c.card, border: `1px solid ${c.border}`, borderRadius: 0, fontSize: 12, color: c.fg }}
              labelStyle={{ fontWeight: 700 }}
            />
            <Bar dataKey="count" name="Plates" fill={c.chart5} radius={[2, 2, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      )}
    </Panel>
  )
}

/* ─── Cumulative growth area ─── */
function CumulativePanel() {
  const c = useChartColors()
  const { data, loading } = useFetch<MonthlyCount[]>("/plates/stats/growth?months=24")
  const rows = (data ?? []).reduce<{ label: string; total: number }[]>((acc, r) => {
    const prev = acc.length > 0 ? acc[acc.length - 1].total : 0
    acc.push({ label: r.month.slice(5), total: prev + r.count })
    return acc
  }, [])
  return (
    <Panel title="Total Plates Over Time" subtitle="Cumulative plate count growth">
      {loading ? <PanelLoader /> : rows.length === 0 ? (
        <p className="text-sm text-muted-foreground">No data yet.</p>
      ) : (
        <ResponsiveContainer width="100%" height={280}>
          <LineChart data={rows} margin={{ top: 8, right: 8, bottom: 0, left: -12 }}>
            <CartesianGrid strokeDasharray="3 3" stroke={c.fg} strokeOpacity={0.1} />
            <XAxis dataKey="label" tick={{ fontSize: 11, fill: c.fg, fillOpacity: 0.6 }} stroke={c.fg} strokeOpacity={0.1} />
            <YAxis allowDecimals={false} tick={{ fontSize: 11, fill: c.fg, fillOpacity: 0.6 }} stroke={c.fg} strokeOpacity={0.1} />
            <Tooltip
              contentStyle={{ backgroundColor: c.card, border: `1px solid ${c.border}`, borderRadius: 0, fontSize: 12, color: c.fg }}
              labelStyle={{ fontWeight: 700 }}
            />
            <Line type="monotone" dataKey="total" name="Total" stroke={c.chart5} strokeWidth={2.5} dot={{ r: 3 }} />
          </LineChart>
        </ResponsiveContainer>
      )}
    </Panel>
  )
}

/* ─── Category distribution ─── */
function CategoryPanel() {
  const c = useChartColors()
  const { data, loading } = useFetch<CategoryCount[]>("/plates/stats/categories")
  const rows = (data ?? []).slice(0, 10)
  return (
    <Panel title="Plates by Category" subtitle="Distribution of templates across categories">
      {loading ? <PanelLoader /> : rows.length === 0 ? (
        <p className="text-sm text-muted-foreground">No data yet.</p>
      ) : (
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={rows} layout="vertical" margin={{ top: 0, right: 8, bottom: 0, left: 0 }}>
            <CartesianGrid strokeDasharray="3 3" stroke={c.fg} strokeOpacity={0.1} horizontal={false} />
            <XAxis type="number" allowDecimals={false} tick={{ fontSize: 11, fill: c.fg, fillOpacity: 0.6 }} stroke={c.fg} strokeOpacity={0.1} />
            <YAxis dataKey="category" type="category" width={100} tick={{ fontSize: 11, fill: c.fg, fillOpacity: 0.6 }} stroke={c.fg} strokeOpacity={0.1} />
            <Tooltip
              contentStyle={{ backgroundColor: c.card, border: `1px solid ${c.border}`, borderRadius: 0, fontSize: 12, color: c.fg }}
            />
            <Bar dataKey="count" name="Plates" fill={c.chart5} radius={[0, 2, 2, 0]} />
          </BarChart>
        </ResponsiveContainer>
      )}
    </Panel>
  )
}

/* ─── Badge tier donut ─── */
function BadgeTierPanel() {
  const c = useChartColors()
  const pieShades = [c.chart5, c.chart4, c.chart3, c.chart2, c.chart1]
  const { data: badges, isLoading } = useBadges()
  const rows = (() => {
    if (!badges?.length) return []
    const map = new Map<string, number>()
    for (const b of badges) map.set(b.tier, (map.get(b.tier) ?? 0) + 1)
    return Array.from(map.entries())
      .sort((a, b) => b[1] - a[1])
      .map(([tier, count]) => ({ tier, count }))
  })()
  const total = rows.reduce((s, r) => s + r.count, 0)

  return (
    <Panel title="Badge Tier Distribution" subtitle="Share of defined badges by tier">
      {isLoading ? <PanelLoader /> : rows.length === 0 || total === 0 ? (
        <p className="text-sm text-muted-foreground">No data yet.</p>
      ) : (
        <div className="grid grid-cols-1 items-center gap-6 sm:grid-cols-[auto_1fr]">
          <ResponsiveContainer width={220} height={220}>
            <PieChart>
              <Pie
                data={rows}
                dataKey="count"
                nameKey="tier"
                cx="50%"
                cy="50%"
                innerRadius={60}
                outerRadius={90}
                paddingAngle={2}
                strokeWidth={0}
              >
                {rows.map((_, i) => (
                  <Cell key={i} fill={pieShades[i % pieShades.length]} />
                ))}
              </Pie>
              <Tooltip
                contentStyle={{ backgroundColor: c.card, border: `1px solid ${c.border}`, borderRadius: 0, fontSize: 12, color: c.fg }}
              />
            </PieChart>
          </ResponsiveContainer>
          <div className="space-y-2">
            {rows.map((r, i) => (
              <div key={r.tier} className="flex items-center justify-between border-b border-border pb-2 text-sm last:border-b-0">
                <div className="flex items-center gap-2">
                  <span className="h-2.5 w-2.5" style={{ backgroundColor: pieShades[i % pieShades.length] }} />
                  <span className="capitalize text-foreground">{r.tier}</span>
                </div>
                <span className="tabular-nums text-muted-foreground">
                  {r.count} ({((r.count / total) * 100).toFixed(0)}%)
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </Panel>
  )
}

/* ─── Top bookmarked ─── */
function TopBookmarkedPanel() {
  const { data, loading } = useFetch<PlateRanked[]>("/plates/stats/top-bookmarked?limit=10")
  const rows = data ?? []
  return (
    <Panel title="Most Bookmarked Plates" subtitle="Top templates by total bookmarks">
      {loading ? <PanelLoader /> : rows.length === 0 ? (
        <p className="text-sm text-muted-foreground">No data yet.</p>
      ) : (
        <div className="space-y-3">
          {rows.map((r, i) => {
            const max = rows[0]?.bookmark_count ?? 1
            const w = `${Math.max((r.bookmark_count / max) * 100, 8)}%`
            return (
              <div key={r.id} className="space-y-1.5">
                <div className="flex items-center justify-between gap-4">
                  <div className="flex items-center gap-2 min-w-0">
                    <span className="shrink-0 w-5 text-xs tabular-nums text-muted-foreground text-right">{i + 1}.</span>
                    <Link href={buildPlateHref(r.slug)} className="text-sm font-medium text-foreground truncate hover:underline">
                      {r.name}
                    </Link>
                    <span className="shrink-0 text-[10px] text-muted-foreground capitalize">{r.category}</span>
                  </div>
                  <p className="shrink-0 text-sm tabular-nums text-muted-foreground">{formatCount(r.bookmark_count)}</p>
                </div>
                <div className="h-1.5 w-full bg-muted/60">
                  <div className="h-1.5 bg-chart-5" style={{ width: w }} />
                </div>
              </div>
            )
          })}
        </div>
      )}
    </Panel>
  )
}

/* ─── Top rated ─── */
function TopRatedPanel() {
  const { data, loading } = useFetch<PlateRanked[]>("/plates/stats/top-rated?limit=10")
  const rows = data ?? []
  return (
    <Panel title="Highest Rated Plates" subtitle="Top templates by average rating">
      {loading ? <PanelLoader /> : rows.length === 0 ? (
        <p className="text-sm text-muted-foreground">No data yet.</p>
      ) : (
        <div className="space-y-3">
          {rows.map((r, i) => {
            const max = 5
            const w = `${Math.max((r.avg_rating / max) * 100, 8)}%`
            return (
              <div key={r.id} className="space-y-1.5">
                <div className="flex items-center justify-between gap-4">
                  <div className="flex items-center gap-2 min-w-0">
                    <span className="shrink-0 w-5 text-xs tabular-nums text-muted-foreground text-right">{i + 1}.</span>
                    <Link href={buildPlateHref(r.slug)} className="text-sm font-medium text-foreground truncate hover:underline">
                      {r.name}
                    </Link>
                    <span className="shrink-0 text-[10px] text-muted-foreground capitalize">{r.category}</span>
                  </div>
                  <p className="shrink-0 text-sm tabular-nums text-muted-foreground">{r.avg_rating.toFixed(1)} ★</p>
                </div>
                <div className="h-1.5 w-full bg-muted/60">
                  <div className="h-1.5 bg-chart-4" style={{ width: w }} />
                </div>
              </div>
            )
          })}
        </div>
      )}
    </Panel>
  )
}

/* ─── main page ─── */
export default function StatsPage() {
  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto px-4 py-10 sm:py-12">
        <div className="border border-border bg-muted/20 p-5 sm:p-7">
          <p className="mb-2 text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
            Public Dashboard
          </p>
          <h1 className="text-3xl font-black tracking-tight text-foreground sm:text-4xl">
            KikPlate Stats
          </h1>
          <p className="mt-3 max-w-3xl text-sm leading-relaxed text-muted-foreground sm:text-base">
            Live platform metrics and distribution charts from public registry data.
            Each panel loads independently and updates automatically.
          </p>
        </div>

        <div className="mt-6">
          <KpiPanel />
        </div>

        <div className="mt-8 grid grid-cols-1 gap-4 xl:grid-cols-2">
          <GrowthPanel />
          <CumulativePanel />
        </div>

        <div className="mt-4 grid grid-cols-1 gap-4 xl:grid-cols-2">
          <CategoryPanel />
          <BadgeTierPanel />
        </div>

        <div className="mt-4 grid grid-cols-1 gap-4 xl:grid-cols-2">
          <TopBookmarkedPanel />
          <TopRatedPanel />
        </div>

        <div className="mt-6 flex items-center gap-2 text-xs text-muted-foreground">
          <BarChart3 className="h-3.5 w-3.5" />
          <ShieldCheck className="h-3.5 w-3.5" />
          <span>Data source: public API endpoints</span>
        </div>
      </div>
    </div>
  )
}