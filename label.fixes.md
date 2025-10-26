DONE. [HIGH] cmd/gh/label/sync.go:54 + internal/gh/labels.go:26 – ReadGithubLabelsFromConfig bubbles up os.ReadFile/ReadYAMLField errors when .hexa.yml or github.labels is absent, so the new hw gh label sync command exits before the confirm prompt and never fetches labels on a fresh repo. Guard the “missing file/field” case and treat it as “no existing data” instead of an error.

DONE [HIGH] cmd/gh/label/sync.go:85 – config.GetProjectConfigPath still relies on os.Getwd(), so running the sync from a subdirectory will write .hexa.yml into that subdirectory rather than the repo root that git rev-parse --show-toplevel reports. That’s a regression now that the command persists data; please resolve through the git-root helper you already added (internal/git.GetRepoRootPath).

[HIGH] alexandria (new symlink) – this introduces an absolute path to /Users/maximilien/... into the repo. It won’t exist on other machines and will break packaging; looks accidental and should be removed from the commit.

[MEDIUM] cmd/gh/label/sync.go:98 – the command shells out to jq for pretty-printing but never checks that the binary exists, so on hosts without jq the command fails after mutating the YAML. Prefer json.MarshalIndent (or at least detect the missing binary and fall back).

[MEDIUM] cmd/gh/label/sync.go:65 – huh.NewConfirm().Run() errors (for example when run in a non-TTY, CI, or when bubbletea cannot initialize) are ignored, leaving confirm at its previous value. Please handle the error and surface a helpful message.

[MEDIUM] go.mod – adding github.com/charmbracelet/huh pulls ~20 transitive UI dependencies just to ask a yes/no question. That’s a hefty footprint for a CLI helper and complicates distribution; consider implementing confirms with stdlib (bufio.Scanner, fmt) or reuse an existing lightweight helper if you already vendor Bubble Tea elsewhere.

[LOW] cmd/gh/label/sync.go:19 – the package-level confirm bool carries state across command invocations and makes tests harder; an idiomatic approach keeps the variable local to the function that needs it.
[LOW] cmd/gh/label/sync.go:43/79 – leftover debug/commented lines and the "c'est partiii" message stray from the project's English-facing output style; please clean them up.
[LOW] internal/gh/utils.go:10/18 – the guard helpers discard the CLI output, so users only see "not authenticated/no remote detected". Capturing and returning the trimmed stderr would make the failure actionable.
Open questions / assumptions

Should the sync command also guard against the absence of gh/jq by checking them up front (similar to VerifyGhAuthenticated)? That would align with the existing pre-flight pattern.
Do we want to surface a non-interactive mode (flag) for CI usage, or is this intentionally interactive-only?
Change summary

Adds a new gh label sync Cobra command with supporting Git/GitHub helpers, YAML read/write enhancements, and Charmbracelet-based prompts.
Introduces a label-sync shell script and registers the gh command tree in main.go.
Updates configuration helpers to expose project-config paths and nested YAML lookups, and bumps Go module deps accordingly.
Next steps: once the blockers above are resolved, I'd rerun npm run test and npm run build to ensure the new command wires cleanly before opening the PR.

---

## Additional Critical Analysis

### Structural Issues

[HIGH] internal/config/yaml.go:164 – `splitKey()` function is undefined. The nested key parsing loop uses `splitKey(key)` but this function doesn't exist in the diff. Either it's missing or the code doesn't compile.

[HIGH] internal/gh/labels.go:33-72 – `normalizeKeys()` is completely unnecessary. YAML unmarshaling with struct tags already handles case normalization. This function adds 40 lines of complexity for zero benefit. The recursive traversal of maps/slices is pure noise when `yaml.Unmarshal` already does this work.

[MEDIUM] cmd/gh/label/sync.go:54 + 82 – double fetch pattern creates confusion. `ReadGithubLabelsFromConfig()` is called twice in different code paths (initial check + after confirm). This wastes I/O and makes the control flow harder to reason about.

[MEDIUM] cmd/gh/label/sync.go:19 – global `confirm` variable is worse than LOW priority. It's never reset between invocations, so if the command runs twice in the same process (tests, REPL, etc.), the second run inherits stale state. Should be function-local or explicitly reset.

### Missing Validation

[MEDIUM] internal/gh/labels.go:73-86 – no validation after unmarshal. `FetchLabels()` blindly trusts GitHub API response and YAML config. Labels without `name` or `color` will silently corrupt state. Add validation:
```go
for i, label := range labels {
    if label.Name == "" || label.Color == "" {
        return nil, fmt.Errorf("label[%d] missing required fields", i)
    }
}
```

[LOW] cmd/gh/label/sync.go:97-119 – jq pretty-print happens after writing to disk. If `jq` fails, the YAML is already mutated. Either check `jq` availability up front or don't shell out at all (use `json.MarshalIndent` directly to stdout without disk write).

### Architecture Questions

[INFO] cmd/gh/label/label_sync.sh – this 84-line shell script is dead code. It implements the same sync logic as `sync.go` but is never invoked. If this is migration scaffolding for incremental porting, document it explicitly or remove it to avoid confusion.

[INFO] Migration path unclear – the shell script writes to `repo.config.json` while the Go command writes to `.hexa.yml` under `github.labels`. Are these meant to coexist? If the shell version is deprecated, mark it clearly or delete it.

### Code Quality

[LOW] internal/gh/labels.go:12-16 – struct tags include `yaml:"name"` but `json:"name"` is also present. If this only ever comes from YAML config, the json tags are dead weight. If it's also serialized to JSON elsewhere, document why both are needed.

[LOW] cmd/gh/label/sync.go:78 – hardcoded message "c'est partiii\n" should be removed or replaced with proper English output like "Fetching labels from GitHub...". Inconsistent with the rest of the CLI's language.

### Recommendations

1. **Fix `splitKey()` immediately** – either implement it or use an existing helper. This is a compilation blocker.

2. **Delete `normalizeKeys()`** – YAML already handles this. If there's a real case-sensitivity issue, show a failing test case first.

3. **Consolidate or document migration** – either finish migrating from shell to Go and delete the bash script, or explicitly mark `label_sync.sh` as deprecated with a comment pointing to the Go implementation.

4. **Add pre-flight check for jq** – align with existing `VerifyGhAuthenticated` pattern. Better yet, stop using `jq` entirely and use stdlib.

5. **Add label validation** – don't trust API/config blindly. Enforce required fields.

6. **Make confirm local** – move `var confirm bool` inside `writeLabelsInProject()` to avoid state leakage.