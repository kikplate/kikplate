package lib

import "strings"

type PlateCategoryConfig struct {
	Slug        string `mapstructure:"slug" json:"slug"`
	Label       string `mapstructure:"label" json:"label"`
	Description string `mapstructure:"description" json:"description"`
	Icon        string `mapstructure:"icon" json:"icon"`
}

const plateCategoryFallbackSlug = "other"

func DefaultPlateCategories() []PlateCategoryConfig {
	return []PlateCategoryConfig{
		{Slug: "backend", Label: "Backend", Description: "APIs, services, microservices", Icon: "server"},
		{Slug: "frontend", Label: "Frontend", Description: "Web UIs, SPAs, SSR apps", Icon: "globe"},
		{Slug: "fullstack", Label: "Full Stack", Description: "End-to-end project starters", Icon: "layers"},
		{Slug: "mobile", Label: "Mobile", Description: "iOS, Android, cross-platform", Icon: "smartphone"},
		{Slug: "cli", Label: "CLI", Description: "Command line tools & scripts", Icon: "terminal"},
		{Slug: "devops", Label: "DevOps", Description: "Docker, CI/CD, infrastructure", Icon: "wrench"},
		{Slug: "library", Label: "Library", Description: "Reusable packages and SDKs", Icon: "package"},
		{Slug: "database", Label: "Database", Description: "Schemas, migrations, seeds", Icon: "database"},
		{Slug: "cloud", Label: "Cloud", Description: "AWS, GCP, Azure starters", Icon: "cloud"},
		{Slug: "security", Label: "Security", Description: "Auth, encryption, compliance", Icon: "shield"},
		{Slug: "iot", Label: "IoT", Description: "Embedded, edge, hardware", Icon: "cpu"},
		{Slug: "game", Label: "Game Dev", Description: "Game engines, frameworks", Icon: "gamepad-2"},
		{Slug: "documentation", Label: "Documentation", Description: "Docs sites, wikis, guides", Icon: "book-open"},
		{Slug: "ai-ml", Label: "AI / ML", Description: "Machine learning, LLMs, data", Icon: "bot"},
		{Slug: "other", Label: "Other", Description: "Everything else", Icon: "more-horizontal"},
	}
}

func EffectivePlateCategories(env Env) []PlateCategoryConfig {
	var base []PlateCategoryConfig
	if len(env.PlateCategories) > 0 {
		base = append([]PlateCategoryConfig{}, env.PlateCategories...)
	} else {
		base = DefaultPlateCategories()
	}
	if !plateCategoriesContainSlugCI(base, plateCategoryFallbackSlug) {
		for _, d := range DefaultPlateCategories() {
			if strings.EqualFold(strings.TrimSpace(d.Slug), plateCategoryFallbackSlug) {
				base = append(base, d)
				break
			}
		}
	}
	return base
}

func plateCategoriesContainSlugCI(list []PlateCategoryConfig, slug string) bool {
	want := strings.ToLower(strings.TrimSpace(slug))
	for _, c := range list {
		if strings.ToLower(strings.TrimSpace(c.Slug)) == want {
			return true
		}
	}
	return false
}

func plateCategoryCanonByLowerKey(env Env) map[string]string {
	out := make(map[string]string)
	for _, c := range EffectivePlateCategories(env) {
		slug := strings.TrimSpace(c.Slug)
		if slug == "" {
			continue
		}
		out[strings.ToLower(slug)] = slug
	}
	return out
}

func NormalizePlateCategory(env Env, raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return resolveFallbackCanonSlug(env)
	}
	if canon, ok := plateCategoryCanonByLowerKey(env)[strings.ToLower(s)]; ok {
		return canon
	}
	return resolveFallbackCanonSlug(env)
}

func resolveFallbackCanonSlug(env Env) string {
	idx := plateCategoryCanonByLowerKey(env)
	if canon, ok := idx[strings.ToLower(plateCategoryFallbackSlug)]; ok {
		return canon
	}
	return plateCategoryFallbackSlug
}

func PlateCategorySlugs(env Env) []string {
	cats := EffectivePlateCategories(env)
	out := make([]string, 0, len(cats))
	for _, c := range cats {
		if s := strings.TrimSpace(c.Slug); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func NormalizePlateCategoryFilter(env Env, in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := make(map[string]struct{})
	out := make([]string, 0, len(in))
	for _, c := range in {
		n := NormalizePlateCategory(env, c)
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out
}
