# Example Aliases

Sample aliases you can create with `alias_manager`, grouped by category.
Import the JSON block at the bottom, or copy individual `create` commands.

## Quick import

Save the JSON section as a file, then:

```sh
alias_manager import examples.json
# or overwrite existing names:
alias_manager import examples.json --overwrite
```

## Git

| Alias | Command |
|-------|---------|
| `gl` | `git log --oneline --graph --all --decorate` |
| `gco` | `git checkout` |
| `gchb` | `git checkout -` |

```sh
alias_manager create gl "git log --oneline --graph --all --decorate"
alias_manager create gco "git checkout"
alias_manager create gchb "git checkout -"
```

## Kubernetes

| Alias | Command |
|-------|---------|
| `k` | `kubectl` |
| `kgp` | `kubectl get pods` |
| `kgpn` | `kubectl get pods -n` |
| `kgd` | `kubectl get deployments` |
| `kgs` | `kubectl get services` |
| `kgn` | `kubectl get nodes` |
| `kgsec` | `kubectl get secrets` |
| `kgcm` | `kubectl get configmap` |
| `kgctx` | `kubectl config get-contexts` |
| `kctx` | `kubectl config use-context` |
| `kns` | `kubectl config set-context --current --namespace` |
| `kl` | `kubectl logs` |
| `klf` | `kubectl logs -f` |
| `kex` | `kubectl exec -it` |
| `kcp` | `kubectl cp` |

```sh
alias_manager create k "kubectl"
alias_manager create kgp "kubectl get pods"
alias_manager create kgpn "kubectl get pods -n"
alias_manager create kgd "kubectl get deployments"
alias_manager create kgs "kubectl get services"
alias_manager create kgn "kubectl get nodes"
alias_manager create kgsec "kubectl get secrets"
alias_manager create kgcm "kubectl get configmap"
alias_manager create kgctx "kubectl config get-contexts"
alias_manager create kctx "kubectl config use-context"
alias_manager create kns "kubectl config set-context --current --namespace"
alias_manager create kl "kubectl logs"
alias_manager create klf "kubectl logs -f"
alias_manager create kex "kubectl exec -it"
alias_manager create kcp "kubectl cp"
```

## Docker

| Alias | Command |
|-------|---------|
| `d` | `docker` |
| `di` | `docker images` |
| `dps` | `docker ps` |
| `dpsa` | `docker ps -a` |
| `dex` | `docker exec -it` |
| `dcup` | `docker compose up` |
| `dcdown` | `docker compose down` |

```sh
alias_manager create d "docker"
alias_manager create di "docker images"
alias_manager create dps "docker ps"
alias_manager create dpsa "docker ps -a"
alias_manager create dex "docker exec -it"
alias_manager create dcup "docker compose up"
alias_manager create dcdown "docker compose down"
```

## tmux

| Alias | Command |
|-------|---------|
| `tmuxnew` | `tmux new -s` |
| `tmuxls` | `tmux ls` |
| `tmuxa` | `tmux attach -t` |

```sh
alias_manager create tmuxnew "tmux new -s"
alias_manager create tmuxls "tmux ls"
alias_manager create tmuxa "tmux attach -t"
```

## gcloud

| Alias | Command |
|-------|---------|
| `gal` | `gcloud auth login` |
| `gkels` | `gcloud container clusters list` |

```sh
alias_manager create gal "gcloud auth login"
alias_manager create gkels "gcloud container clusters list"
```

## alias_manager helpers

| Alias | Command |
|-------|---------|
| `am` | `alias_manager` |
| `amls` | `alias_manager list` |
| `amgrep` | `alias_manager list \| grep` |
| `ra` | `source ~/.alias_manager/aliases.sh` |

```sh
alias_manager create am "alias_manager"
alias_manager create amls "alias_manager list"
alias_manager create amgrep "alias_manager list | grep"
alias_manager create ra "source ~/.alias_manager/aliases.sh"
```

Note: `ra` is for bash/zsh. On fish use:

```sh
alias_manager create ra "source ~/.alias_manager/aliases.fish"
```

## Misc

| Alias | Command |
|-------|---------|
| `cl` | `clear` |
| `sq` | `sqlite3` |
| `nrd` | `npm run dev` |
| `bashcnfg` | `vi ~/.bashrc` |
| `sshconfig` | `vi ~/.ssh/config` |

```sh
alias_manager create cl "clear"
alias_manager create sq "sqlite3"
alias_manager create nrd "npm run dev"
alias_manager create bashcnfg "vi ~/.bashrc"
alias_manager create sshconfig "vi ~/.ssh/config"
```

## Importable JSON

Copy into `examples.json` and run `alias_manager import examples.json`:

```json
{
  "am": "alias_manager",
  "amgrep": "alias_manager list | grep",
  "amls": "alias_manager list",
  "bashcnfg": "vi ~/.bashrc",
  "cl": "clear",
  "d": "docker",
  "dcdown": "docker compose down",
  "dcup": "docker compose up",
  "dex": "docker exec -it",
  "di": "docker images",
  "dps": "docker ps",
  "dpsa": "docker ps -a",
  "gal": "gcloud auth login",
  "gchb": "git checkout -",
  "gco": "git checkout",
  "gkels": "gcloud container clusters list",
  "gl": "git log --oneline --graph --all --decorate",
  "k": "kubectl",
  "kcp": "kubectl cp",
  "kctx": "kubectl config use-context",
  "kex": "kubectl exec -it",
  "kgcm": "kubectl get configmap",
  "kgctx": "kubectl config get-contexts",
  "kgd": "kubectl get deployments",
  "kgn": "kubectl get nodes",
  "kgp": "kubectl get pods",
  "kgpn": "kubectl get pods -n",
  "kgs": "kubectl get services",
  "kgsec": "kubectl get secrets",
  "kl": "kubectl logs",
  "klf": "kubectl logs -f",
  "kns": "kubectl config set-context --current --namespace",
  "nrd": "npm run dev",
  "ra": "source ~/.alias_manager/aliases.sh",
  "sq": "sqlite3",
  "sshconfig": "vi ~/.ssh/config",
  "tmuxa": "tmux attach -t",
  "tmuxls": "tmux ls",
  "tmuxnew": "tmux new -s"
}
```

## CLI usage examples

```sh
# Setup
alias_manager init
alias_manager add-to-shell && source ~/.zshrc   # or ~/.bashrc / fish config

# Create / list / edit / delete
alias_manager create gs "git status"
alias_manager list
alias_manager list git
alias_manager edit gs "git status -sb"
alias_manager delete gs --yes

# Interactive TUI (TTY, missing args)
alias_manager create
alias_manager edit gs

# Import / export
alias_manager export ~/aliases-backup.json
alias_manager import ~/aliases-backup.json --overwrite

# Suggestions and maintenance
alias_manager suggest git
alias_manager rebuild-shell-file
alias_manager version
```
