# Configuration Management

Hexa CLI uses a multi-level configuration system that provides maximum flexibility while maintaining security and ease of use.

## Configuration Hierarchy

Configuration values are resolved in the following order of precedence (1 = highest priority):

1. **🚩 CLI Flags** - Runtime overrides

   ```bash
   hexa config --user-email="override@cli.com"
   ```

2. **🌍 Environment Variables** - System environment with `HEXA_` prefix

   ```bash
   export HEXA_USER_EMAIL="env@company.com"
   ```

3. **📁 Project Local Secrets** - `./.hexa.local.yml` (gitignored)

   - Project-specific secrets and local settings
   - Never committed to version control

4. **📁 Project Config** - `./.hexa.yml` (committed)

   - Project-specific configuration overrides
   - Committed to version control and shared with team

5. **👤 User Global Config** - `~/.hexa.yml`

   - Personal base configuration
   - Applied across all projects as foundation

6. **💾 Code Defaults** - Hardcoded fallback values

## Configuration Files

### User Global Configuration (`~/.hexa.yml`)

Personal base configuration applied to all projects:

```yaml
user:
  email: "me@personal.com"
  name: "John Doe"

jira:
  host: "personal.jira.com"
  token: "${HEXA_JIRA_TOKEN}"  # Secret - use env var

defaults:
  debug: false
  timeout: 60s
```

### Project Configuration (`.hexa.yml`)

Project-specific overrides (only when needed):

```yaml
# Override specific values for this project
jira:
  host: "company.jira.com" # Override company Jira
  project: "MYPROJ"
  token: "${HEXA_JIRA_TOKEN}" # Use environment variable

user:
  email: "${HEXA_USER_EMAIL}" # Template for team flexibility
```

### Project Local Configuration (`.hexa.local.yml`)

Gitignored secrets and local overrides:

```yaml
jira:
  token: "secret-company-token" # Direct secret value
  host: "internal.jira.com"     # Local override

user:
  email: "dev@company.com" # Local override if needed
```

## Environment Overrides

Hexa relies on Viper to combine configuration sources. YAML files can keep placeholders such as `${HEXA_JIRA_TOKEN}`; when you ask Viper for a value (e.g. `viper.GetString("jira.token")`), any exported environment variable with the matching `HEXA_` key takes precedence over the literal placeholder. If the environment variable is missing, the placeholder string is returned unchanged so you can detect missing secrets explicitly. This behaviour follows [Viper's environment variable support](https://github.com/spf13/viper#working-with-environment-variables).

### Automatic Binding

Hexa automatically binds environment variables with the `HEXA_` prefix:

```bash
export HEXA_USER_EMAIL="user@example.com"      # May vary by environment
export HEXA_JIRA_TOKEN="secret-token"          # Secret - use env var
export HEXA_JIRA_HOST="jira.company.com"       # May vary by environment
```

These become accessible as:

- `viper.GetString("user.email")`
- `viper.GetString("jira.token")`
- `viper.GetString("jira.host")`

### .env File Support

Hexa automatically loads a `.env` file from the current working directory if present. This provides a convenient way to set environment variables for development without modifying your shell environment.

**⚠️ Security Warning**: `.env` files may contain sensitive information. Always add `.env` to your `.gitignore` to prevent accidental commits of secrets.

```bash
# .env file in project root
HEXA_JIRA_TOKEN=dev-secret-token
HEXA_USER_EMAIL=dev@company.com
HEXA_JIRA_HOST=dev.jira.company.com
```

### Manual Binding

For explicit control, specific environment variables can be bound:

```go
viper.BindEnv("user.email", "HEXA_USER_EMAIL")
viper.BindEnv("custom.value", "MY_CUSTOM_VAR")
```

## Security Tips

Ensure sensitive files stay out of version control by keeping a `.gitignore` entry for `.hexa.local.yml`, `*.local.yml`, `.env`, and similar artefacts. Project templates should include these patterns so contributors do not commit secrets by accident.

## Usage Examples

### Basic Setup

1. **Setup your global config once** (`~/.hexa.yml`):

   ```yaml
   user:
     email: "me@personal.com"
     name: "John Doe"

   jira:
     host: "personal.jira.com"
     token: "${HEXA_JIRA_TOKEN}"  # Secret - use env var
   ```

2. **Most projects work immediately** - No additional config needed!

3. **For special projects, add overrides** (`.hexa.yml`):

   ```yaml
   # Only override what's different
   jira:
     host: "company.jira.com"
   ```

4. **For project secrets** (`.hexa.local.yml`):
   ```yaml
   jira:
     token: "company-secret-token"
   ```

### Development Workflow

**First time setup:**

```bash
# Configure once in your home directory
cat > ~/.hexa.yml << 'EOF'
user:
  email: "your@email.com"
  name: "Your Name"
jira:
  host: "your.jira.com"
  token: "${HEXA_JIRA_TOKEN}"  # Secret - use env var
EOF

# Set your secret environment variables
export HEXA_JIRA_TOKEN="your-secret-jira-token"
```

**Working with company projects:**

```bash
# Add project-specific overrides only when needed
cat > .hexa.yml << 'EOF'
jira:
  host: "company.jira.com"
  project: "COMPANY_PROJECT"
  token: "${HEXA_JIRA_TOKEN}"  # Use env var for secret
EOF

# Option 1: Use .env file (gitignored)
cat > .env << 'EOF'
HEXA_JIRA_TOKEN=company-secret-token
HEXA_USER_EMAIL=dev@company.com
EOF

# Option 2: Use local config file (gitignored)
cat > .hexa.local.yml << 'EOF'
jira:
  token: "company-specific-secret-token"  # Direct secret value
EOF
```

### Configuration Debugging

Use the debug mode to see configuration resolution:

```bash
DEBUG=true hexa config
```

This shows:

- Which config files were loaded
- Template expansion results
- Final merged configuration
- Missing environment variables

## Best Practices

### 1. **User-Centric Setup**

- Configure `~/.hexa.yml` once with your base settings
- Most projects will work immediately without additional config
- Only add project configs when specific overrides are needed

### 2. **Security First**

- Use `.hexa.local.yml` for project secrets (automatically gitignored)
- Use `.env` files for development environment variables (automatically gitignored)
- Never commit API tokens or passwords in `.hexa.yml`
- Use environment variables in CI/CD
- Always add `.env` and `*.local.yml` to `.gitignore`

### 3. **Team Collaboration**

- Keep project configs (`.hexa.yml`) minimal - only what's different
- Provide templates using `${HEXA_VAR}` syntax for team flexibility
- Document required environment variables in project README

### 4. **Maintainability**

- User config holds the foundation, projects add specifics
- Use meaningful key names and consistent structure
- Add comments for complex configurations

## Migration from Simple Config

If you currently have a simple configuration setup:

1. **Move to user global** - Put your main config in `~/.hexa.yml`
2. **Project configs become overrides** - Only add `.hexa.yml` when projects need different settings
3. **Secrets go secure** - Move tokens/passwords to:
   - `.hexa.local.yml` for direct values (gitignored)
   - `.env` file for environment variables (gitignored)
   - Environment variables for CI/CD
4. **Leverage templates** - Use `${HEXA_VAR}` syntax for sensitive values

### Example Migration

**Before (insecure):**
```yaml
# .hexa.yml (committed - BAD!)
jira:
  token: "secret-token-exposed"
```

**After (secure):**
```yaml
# .hexa.yml (committed - GOOD!)
jira:
  token: "${HEXA_JIRA_TOKEN}"  # Template for env var
```

```bash
# .env (gitignored - GOOD!)
HEXA_JIRA_TOKEN=secret-token-safe
```

The configuration system is designed to be intuitive: configure once globally, override when needed, keep secrets secure.
