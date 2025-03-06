# MessageMesh Development Guide

## Build/Run Commands
- `wails dev`: Start development environment with hot-reload
- `wails build`: Build the application
- `go run main.go`: Run the Go backend directly
- `pnpm run dev`: Start frontend development server
- `pnpm run build`: Build the frontend
- `pnpm run check`: Run TypeScript type checking

## Code Style
### Go
- Naming: PascalCase for exported types/functions, camelCase for private
- Error handling: Use debug.Log* functions for logging errors
- Imports: Group standard, external, then internal packages
- Formatting: Use `gofmt` defaults

### Frontend (Svelte/TypeScript)
- Component naming: PascalCase with .svelte extension
- Styling: Tailwind CSS utility classes
- TypeScript: Enable strict mode and type all props
- State management: Prefer Svelte stores for shared state

## Project Structure
- `backend/`: Go backend services (p2p, consensus, database)
- `backend/models/`: Data structures shared with frontend
- `frontend/src/`: Svelte application and components