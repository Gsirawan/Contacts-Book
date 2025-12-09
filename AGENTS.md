# Agent Guidelines for contact-book

## Build/Run/Test Commands
- **Run main app**: `go run contacts.go`
- **Run TUI app**: `go run bubble-tea/contacts-tui.go`
- **Build**: `go build contacts.go` or `go build -o contacts-tui bubble-tea/contacts-tui.go`
- **Install deps**: `go mod tidy`
- **No tests exist yet** - test files would use `go test ./...` or `go test -run TestFunctionName`

## Code Style

**Imports**: Standard library first, blank line, then third-party (e.g., `golang.org/x/text`, `github.com/charmbracelet/bubbletea`)

**Formatting**: Use `gofmt` - tabs for indentation, organize imports, octal file permissions as `0o644`

**Types**: Exported types use PascalCase (e.g., `Contact`), unexported use camelCase. Struct fields capitalized for export.

**Naming**: Clear descriptive names - `addContact()`, `listContact()`, `deleteContact()`. Use single-letter vars only in short scopes.

**Error Handling**: Check all errors. Use `log.Fatalf()` for fatal errors with format `"Error <action> %v\n"`. Defer file closes immediately after opening.

**Constants**: Group related constants in `const` blocks. File paths and UI text as constants (e.g., `filePath = "contacts.txt"`).

**Comments**: Function comments start with function name (e.g., `// create new contacts`). Keep concise.

**Validation**: Regex for emails (`^[a-zA-Z0-9._%+-]+@...`), mobile must be 10 digits starting with 0, convert 0 to +971.

**Data**: CSV format `Name,Email,Mobile` in contacts.txt. Use `bufio.Scanner` for reading, `os.OpenFile` with explicit flags.
