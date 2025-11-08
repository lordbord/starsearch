# AGENTS.md - Development Guidelines for starsearch

## Build/Test Commands

- **Build**: `go build -o starsearch ./cmd/starsearch`
- **Run**: `go run ./cmd/starsearch`
- **Test**: `go test ./...` (no test files exist yet)
- **Test single package**: `go test ./internal/app` (when tests exist)
- **Lint**: `go vet ./...` (built-in Go vet)
- **Format**: `go fmt ./...`

## Code Style Guidelines

### Imports
- Group imports: stdlib, third-party, internal packages
- Use absolute imports for internal packages (e.g., `"starsearch/internal/types"`)

### Naming Conventions
- **Packages**: lowercase, single word when possible
- **Types**: PascalCase (e.g., `Model`, `Document`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase, descriptive names
- **Constants**: PascalCase for exported, camelCase for unexported

### Error Handling
- Always handle errors explicitly
- Use `fmt.Errorf` with `%w` for error wrapping
- Return errors from functions, don't panic unless unrecoverable

### Types & Interfaces
- Define shared types in `internal/types/`
- Use concrete types unless abstraction is needed
- Prefer composition over inheritance

### File Organization
- One main type per file where practical
- Keep related functions together
- Use `internal/` for non-exported packages

### Bubble Tea Patterns
- Model struct should contain all state
- Update method handles messages and returns commands
- View method returns the UI representation
- Use tea.Cmd for async operations

### Gemini Protocol
- Follow Project Gemini specification
- Handle all status codes appropriately
- Use TOFU (Trust On First Use) for certificates