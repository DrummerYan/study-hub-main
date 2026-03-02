# Repository Guidelines

## Project Structure & Module Organization
- `server/`: Go API and admin backend built on `gin-vue-admin`. Key layers mirror each other by domain: `api/v1`, `service`, `router`, and `model`. Shared setup lives in `core/`, `initialize/`, `global/`, and `config/`.
- `web/`: Vue 3 + Vite admin UI. Application code is under `web/src`, with feature views in `src/view`, shared UI in `src/components`, API calls in `src/api`, and state in `src/pinia`.
- `wx/`: uni-app WeChat client. Pages and business code live in `wx/src`; static assets are in `wx/public` and `wx/src/static`.

Deployment manifests are in `deploy/`, and top-level packaging helpers live in `Makefile`.

## Build, Test, and Development Commands
- `cd server && go run main.go`: start the backend with `server/config.yaml`.
- `cd server && go test ./...`: run Go tests.
- `cd web && npm install && npm run serve`: start the admin UI locally.
- `cd web && npm run build`: build the admin UI for production.
- `cd wx && yarn install && yarn dev:mp-weixin`: run the mini-program in watch mode.
- `cd wx && yarn build:mp-weixin`: build the WeChat client.
- `cd wx && yarn test:mp-weixin`: run uni-app Jest tests for the WeChat target.
- `make build-local`: package the web bundle and Go binary together.

## Coding Style & Naming Conventions
Use `gofmt` for Go code; keep packages lowercase and match existing domain file names such as `edu_organization` and `sys_user`. In `web/` and `wx/`, follow the established 2-space indentation and the ESLint rules in [`web/.eslintrc.js`](/Users/yandrummer/Desktop/study-hub-main/web/.eslintrc.js). Use `PascalCase` for Vue components and `camelCase` for functions, stores, and utilities.

## Testing Guidelines
Backend tests live next to source files as `*_test.go`, mainly in `server/utils`. Add or update tests when touching shared logic, validators, timers, or parsing utilities. Frontend test coverage is limited; when changing `wx/`, run the closest Jest target and document manual verification for `web/` or mini-program UI flows.

## Commit & Pull Request Guidelines
Recent history uses both concise Chinese summaries and Conventional Commit prefixes such as `feat:` and `fix:`. Prefer short, imperative subjects and include the affected area when useful, for example `feat: add organization payment export`. PRs should describe the user-facing change, list touched apps, note config or schema changes, and include screenshots for UI updates.

## Security & Configuration Tips
Do not commit real secrets in [`server/config.yaml`](/Users/yandrummer/Desktop/study-hub-main/server/config.yaml) or Docker config variants. Keep local credentials and WeChat app settings out of screenshots and sample data.
