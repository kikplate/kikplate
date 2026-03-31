"use client"

import { useRouter } from "next/navigation"
import {
  Server, Globe, Layers, Smartphone,
  Terminal, Package, Wrench, MoreHorizontal,
  Database, Cloud, Shield, Cpu, Gamepad2, BookOpen, Bot
} from "lucide-react"

const CATEGORIES = [
  { slug: "backend",        label: "Backend",        icon: Server,         description: "APIs, services, microservices" },
  { slug: "frontend",       label: "Frontend",       icon: Globe,          description: "Web UIs, SPAs, SSR apps" },
  { slug: "fullstack",      label: "Full Stack",     icon: Layers,         description: "End-to-end project starters" },
  { slug: "mobile",         label: "Mobile",         icon: Smartphone,     description: "iOS, Android, cross-platform" },
  { slug: "cli",            label: "CLI",            icon: Terminal,       description: "Command line tools & scripts" },
  { slug: "devops",         label: "DevOps",         icon: Wrench,         description: "Docker, CI/CD, infrastructure" },
  { slug: "library",        label: "Library",        icon: Package,        description: "Reusable packages and SDKs" },
  { slug: "database",       label: "Database",       icon: Database,       description: "Schemas, migrations, seeds" },
  { slug: "cloud",          label: "Cloud",          icon: Cloud,          description: "AWS, GCP, Azure starters" },
  { slug: "security",       label: "Security",       icon: Shield,         description: "Auth, encryption, compliance" },
  { slug: "iot",            label: "IoT",            icon: Cpu,            description: "Embedded, edge, hardware" },
  { slug: "game",           label: "Game Dev",       icon: Gamepad2,       description: "Game engines, frameworks" },
  { slug: "documentation",  label: "Documentation",  icon: BookOpen,       description: "Docs sites, wikis, guides" },
  { slug: "ai-ml",          label: "AI / ML",        icon: Bot,            description: "Machine learning, LLMs, data" },
  { slug: "other",          label: "Other",          icon: MoreHorizontal, description: "Everything else" },
]

export function CategoriesGrid() {
  const router = useRouter()

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
          <p className="text-xs tabular-nums text-muted-foreground border border-border bg-card px-3 py-1.5">{CATEGORIES.length} domains</p>
        </div>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
          {CATEGORIES.map(({ slug, label, icon: Icon, description }) => (
            <button
              key={slug}
              onClick={() => router.push(`/explore?category=${slug}`)}
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
          ))}
        </div>

      </div>
    </section>
  )
}