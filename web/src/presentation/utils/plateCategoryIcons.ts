import {
  BookOpen,
  Bot,
  Cloud,
  Cpu,
  Database,
  Gamepad2,
  Globe,
  Layers,
  MoreHorizontal,
  Package,
  Server,
  Shield,
  Smartphone,
  Terminal,
  Wrench,
  type LucideIcon,
} from "lucide-react"

const plateCategoryIconMap: Record<string, LucideIcon> = {
  "server": Server,
  "globe": Globe,
  "layers": Layers,
  "smartphone": Smartphone,
  "terminal": Terminal,
  "wrench": Wrench,
  "package": Package,
  "database": Database,
  "cloud": Cloud,
  "shield": Shield,
  "cpu": Cpu,
  "gamepad-2": Gamepad2,
  "book-open": BookOpen,
  "bot": Bot,
  "more-horizontal": MoreHorizontal,
}

export function getPlateCategoryIcon(iconKey: string): LucideIcon {
  return plateCategoryIconMap[iconKey] ?? MoreHorizontal
}
