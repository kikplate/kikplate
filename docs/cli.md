# CLI Reference

The Kikplate CLI (`kik`) is a standalone Go binary that lets you browse, submit, manage, and scaffold from plates on any Kikplate server.

## Installation

### macOS/Linux via Homebrew tap

```
brew tap kikplate/homebrew-kikplate
brew install --cask kik
```

### Windows via Scoop (alternative)

```
scoop bucket add kikplate https://github.com/kikplate/scoop-bucket.git
scoop install kik
```

### From release archives

Download the archive for your platform from [GitHub Releases](https://github.com/kikplate/kikplate/releases).

```
tar -xzf kik-<version>-linux-amd64.tar.gz
sudo mv kik-<version>-linux-amd64 /usr/local/bin/kik
```

### From Source

```
git clone https://github.com/kikplate/kikplate.git
cd kikplate/cli
go build -o kik .
sudo mv kik /usr/local/bin/kik
```

## Configuration

The CLI stores its configuration at `~/.kik/config.yaml`.

### Initialize

```
kik config init
```

This creates a default config file:

```yaml
server:
  address: http://localhost:3001
auth:
  token: ""
```

### Use a Custom Config Path

Every command accepts a global `--config` flag:

```
kik --config /path/to/config.yaml search --name golang
```

### View Active Config

```
kik config view
```

### Config File Reference

| Key | Description | Default |
|-----|-------------|----------|
| `server.address` | Base URL of the Kikplate API server | `http://localhost:3001` |
| `auth.token` | JWT stored after login. Set automatically by `kik login`. | `""` |

## Authentication

### Login

```
kik login --email you@example.com --password yourpassword
```

On success, the JWT is saved to `~/.kik/config.yaml`. Subsequent commands use it automatically.

### Who Am I

```
kik whoami
```

Prints the username, display name, email, and account ID of the authenticated account.

### Logout

```
kik logout
```

Removes the stored token from the config file.

## Searching and Browsing

### Search Plates

```
kik search [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--name` | Search by name or keyword | |
| `--category` | Filter by category | |
| `--tag` | Filter by tag | |
| `--limit` | Results per page | 20 |
| `--page` | Page number | 1 |

Examples:

```
kik search --name "golang rest"
kik search --category backend --tag gin
kik search --name next.js --limit 10 --page 2
```

Output:

```
SLUG                   NAME                 CATEGORY   RATING   VERIFIED
myorg/go-chi-rest      Go Chi REST API      backend    4.8      yes
myorg/gin-starter      Gin Starter          backend    4.5      yes

Total: 2
```

### Describe a Plate

```
kik describe myorg/go-chi-rest
```

Output includes name, category, status, visibility, repository URL, branch, rating, star count, verification status, tags, owner, organization, and timestamps.

## Managing Your Plates

### View Your Plates

```
kik my plates
```

### View Your Bookmarks

```
kik my bookmarks
```

### View Your Organizations

```
kik my orgs
```

### View Sync Status

```
kik my sync
```

Shows each plate's sync status, last sync time, and next scheduled sync.

## Managing Local Plates

The CLI maintains a local list of plates you use frequently at `~/.kik/plates.json`. This is separate from the server.

### Add a Plate Locally

```
kik plates add myorg/go-chi-rest
```

Fetches the plate details from the server and stores them locally.

### List Local Plates

```
kik plates list
```

### Remove a Local Plate

```
kik plates remove myorg/go-chi-rest
```

## Submitting a Plate

Submission requires authentication.

```
kik submit https://github.com/myorg/my-template
```

Optional flags:

| Flag | Description |
|------|-------------|
| `--branch` | Git branch to use (default: `main`) |
| `--org` | Organization ID to submit under |

Example for an organization-scoped plate:

```
kik submit https://github.com/myorg/my-template --org <org-uuid> --branch main
```

After submission, the CLI prints the assigned slug and the `verification_token` you need to add to your `plate.yaml`.

## Verifying a Plate

After adding the verification token to `plate.yaml` and pushing the commit:

```
kik verify myorg/my-template
```

On success the plate becomes approved, public, and verified. The CLI confirms the new status.

## Scaffolding a Project

Scaffolding creates a new project from a plate. The source repository is cloned, the `plate.yaml` manifest is stripped, and `.kikplate.origin` is written to record the provenance.

### Scaffold to a Local Directory

```
kik scaffold myorg/go-chi-rest my-new-project
```

### Scaffold to a Remote Repository

```
kik scaffold myorg/go-chi-rest https://github.com/you/new-repo.git
```

This clones the template, creates an initial commit, and pushes to the remote. If the remote already contains a scaffolded project, the command refuses to push unless you pass `--force`.

Optional flags:

| Flag | Description |
|------|-------------|
| `--local` | Force scaffold to local directory even if a URL-shaped argument is given |
| `--force` | Overwrite a remote that was already scaffolded from a kik plate |

## Generating a Project

Generation renders files from a plate schema and values, then writes output locally or pushes to a remote repository.

### Generate from Registry Plate

```
kik generate myorg/go-http-server --set projectName=my-app --set modulePath=github.com/you/my-app
```

### Generate from Local Plate Directory

```
kik generate --template ./example-plate --set projectName=my-app --set modulePath=github.com/you/my-app
```

### Generate with Values File

```
kik generate myorg/go-http-server -f values.yaml
```

### Generate and Push to Remote Repository

```
kik generate myorg/go-http-server -f values.yaml --repo https://github.com/you/my-app.git
```

Optional flags:

| Flag | Description |
|------|-------------|
| `-f, --file` | YAML values file |
| `--set` | Inline value override (`key=value`), repeatable |
| `--output-dir` | Output directory (default: `./<slug>`) |
| `--repo` | Push generated output to a remote repository |
| `--template` | Use local plate directory instead of server |
| `--force` | Force push when used with `--repo` |

## Global Flags

| Flag | Description |
|------|-------------|
| `--config` | Path to config file |
| `--help`, `-h` | Show help for any command |

## Command Reference Summary

| Command | Description | Auth Required |
|---------|-------------|--------------|
| `config init` | Create default config file | No |
| `config view` | Print active config | No |
| `login` | Authenticate and save token | No |
| `logout` | Remove stored token | No |
| `whoami` | Show authenticated user | Yes |
| `search` | Search plates on the server | No |
| `describe [slug]` | Show plate details | No |
| `submit [repo-url]` | Submit repository as plate | Yes |
| `verify [slug]` | Verify submitted plate | Yes |
| `scaffold [slug] [target]` | Scaffold project from plate | No |
| `generate [slug]` | Generate project from schema-based plate | No |
| `plates add [slug]` | Add plate to local list | No |
| `plates list` | List locally added plates | No |
| `plates remove [slug]` | Remove plate from local list | No |
| `my plates` | List your plates on server | Yes |
| `my bookmarks` | List bookmarked plates | Yes |
| `my orgs` | List your organizations | Yes |
| `my sync` | View sync status of your plates | Yes |
