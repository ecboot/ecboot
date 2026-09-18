<!--VITE PLUS START-->

# Using Vite+, the Unified Toolchain for the Web

This project is using Vite+, a unified toolchain built on top of Vite, Rolldown, Vitest, tsdown, Oxlint, Oxfmt, and Vite Task. Vite+ wraps runtime management, package management, and frontend tooling in a single global CLI called `vp`. Vite+ is distinct from Vite, and it invokes Vite through `vp dev` and `vp build`. Run `vp help` to print a list of commands and `vp <command> --help` for information about a specific command.

Docs are local at `node_modules/vite-plus/docs` or online at https://viteplus.dev/guide/.

## Built-in Commands vs Scripts

`vp <name>` runs a built-in command. `vp run <name>` runs a `package.json` script or a `vite.config.ts` task. Scripts cannot overwrite built-ins, so `vp dev` and `vp run dev` may do different things. Check `package.json` and `vite.config.ts` first, and run `vp run <name>` when the project defines a script or task with that name.

## Tool Versions

Run `vp toolchain` to show versions and relationships in the active Vite+
release. Add a tool name to select part of the graph. For example, run
`vp toolchain vite`. Use `--global` to ignore the local `vite-plus` package. Use
`vp why <package>` to show the package-manager dependency graph.

## Review Checklist

- [ ] Run `vp install` after pulling remote changes and before getting started.
- [ ] Run `vp check` and `vp test` to format, lint, type check and test changes.
- [ ] Check if there are `vite.config.ts` tasks or `package.json` scripts necessary for validation, run via `vp run <script>`.
- [ ] If setup, runtime, or package-manager behavior looks wrong, run `vp env doctor` and include its output when asking for help.

<!--VITE PLUS END-->

# ECBOOT

Guidance for AI coding agents working in this repository.

## Agent Rules

- Generate commit messages in Chinese, following Conventional Commits style, e.g. `feat: 新增用户登录`, `chore: 升级依赖`.
- 金额（Money）一律使用 `DECIMAL(10,2)` 存储与 `CHAR(3)` ISO 4217 币种（默认 `CNY`）；应用层统一 `BigDecimal`，比较必须用 `compareTo`。禁止 `float`/`double` 存金额，禁止 `equals`/`==` 比较金额。

## Project Overview

ECBOOT is an e-commerce platform monorepo (`org.juling.ecboot`) in early scaffold stage: a Java 25 / Spring Boot 4.1.1 backend at `apps/api` plus four frontend apps, all in a Vite+ / pnpm workspace (Node >= 22.18, pnpm). Most backend modules are empty placeholder POMs describing a planned layering; only `ecboot-start` contains runnable code. Spec Kit scaffolding lives in `.specify/` (spec-driven development via the `speckit-*` skills); the ratified constitution at `.specify/memory/constitution.md` (v1.0.0) defines five principles: 模块化单体 / 统一技术栈 / 中文优先 / 可验证交付 / 简单优先.

## Repository Layout

- `apps/api/` — the Java backend. `apps/api/pom.xml` (`ecboot-parent`) is the Maven aggregator + parent (inherits `spring-boot-starter-parent`); modules are flat beneath it:
  - `ecboot-start` — the single runnable Spring Boot application (`EcbootApplication`); depends on the three channel API modules.
  - `ecboot-api-user` / `ecboot-api-shop` / `ecboot-api-admin` / `ecboot-api-common` — channel API modules (empty `src/` skeletons).
  - `ecboot-service-user` / `ecboot-service-shop` — domain services (placeholder POMs).
  - `ecboot-common` / `ecboot-infra-core` — shared libraries (placeholder POMs).
  - `ecboot-dependencies` — BOM; the single version-arbitration point.
- `apps/web` — `ecboot-web`, storefront (TanStack Start + React 19 + Tailwind CSS 4 + Vite 8).
- `apps/admin` — `ecboot-admin`, admin console (same stack as `apps/web`).
- `apps/mobile` — `ecboot-mobile`, uni-app (Vue 3) targeting mini-programs / H5 / native apps.
- `apps/website` — plain Vite+ site (`vp dev` / `vp build`).
- `packages/utils` — shared TypeScript package (`vp pack` build, `vp test`).
- `deployments/deploy.sh` — deployment script.
- `scripts/codegen.sh` — empty placeholder for future code generation.
- `docs/` — project documentation; `specs/` — feature specs (e.g. `001-maven-module-deps`).

## Commands

### Backend

**All backend builds run from `apps/api/`** — single-module builds fail because `ecboot-dependencies` is resolved from the reactor:

```bash
cd apps/api
./mvnw clean package -DskipTests             # full reactor build
./mvnw test -pl ecboot-start -am             # tests for the runnable app (+ its module deps)
./mvnw spring-boot:run -pl ecboot-start -am  # run the app (auto-starts compose.yaml services; needs Docker)
```

Dependency rules are enforced by maven-enforcer at the `validate` phase — violations fail the build:

- Layer direction: `ecboot-start → ecboot-api-* → ecboot-service-* → ecboot-common / ecboot-infra-core`; `ecboot-api-common` must not depend on `ecboot-service-*`; cycles fail at build time.
- Version arbitration: third-party/framework versions are declared ONLY in `ecboot-dependencies/pom.xml` (`spring-boot.version` property + BOM imports); plugin versions ONLY in `apps/api/pom.xml` `pluginManagement`. Business module POMs carry zero version numbers.
- Per-module ban lists: slot properties `enforcer.banned.1..5` (+ `enforcer.allowed.1` exception), defined per module, defaults in `apps/api/pom.xml`. A single `<exclude>` does NOT support comma-separated lists.

Local infrastructure is defined in `apps/api/compose.yaml`: MySQL 8.4 LTS (db `mydatabase`, user `myuser`/`secret`), Redis, Elasticsearch 9.3.3 (security disabled). Start manually with `docker compose up -d` from `apps/api/`.

Database migrations live in `apps/api/ecboot-start/src/main/resources/db/migration/` (Flyway's default location — they run automatically on app startup). Schema design rationale: `docs/schema-design.md`.

### Frontend

JS dependencies are managed by pnpm from the repo root (workspace globs `apps/*`, `packages/*`, `tools/*`; versions centralized in the `pnpm-workspace.yaml` catalog). Run `pnpm install` after pulling changes; the Vite+ CLI (`vp`) wraps the toolchain (see top of this file).

`ecboot-web` and `ecboot-admin` (TanStack Start + React 19 + Tailwind CSS 4 + Vite 8):

```bash
cd apps/web              # or apps/admin
pnpm run dev             # vite dev server on port 3000 — BOTH projects use 3000; change one if running both
pnpm run generate-routes # regenerate src/routeTree.gen.ts (tsr generate)
pnpm run build
```

`ecboot-mobile` (uni-app, Vue 3):

```bash
cd apps/mobile
pnpm run dev:h5           # H5 dev server
pnpm run dev:mp-weixin    # WeChat mini-program dev build
pnpm run build:mp-weixin  # WeChat mini-program production build
pnpm run type-check       # vue-tsc --noEmit
```

Repo-wide validation: `pnpm run ready` (= `vp check && vp run -r test && vp run -r build`).

## Architecture

### Backend (layering, enforced)

The `apps/api` modules form a modular monolith: `ecboot-common` / `ecboot-infra-core` (shared libs) ← `ecboot-service-user` / `ecboot-service-shop` (domain) ← `ecboot-api-user` / `ecboot-api-shop` / `ecboot-api-admin` (+ `ecboot-api-common`) (API per channel) — assembled by `ecboot-start` into one deployable. The dependency wiring IS implemented and enforced (see Backend commands above).

The `ecboot-start` stack: WebMVC, Security, JPA + Flyway (MySQL), Redis, Elasticsearch, Quartz, Mail, WebSocket, RestClient, Validation, Actuator; plus runtime-optional devtools and spring-boot-docker-compose. Build-time extras: Lombok + `spring-boot-configuration-processor` annotation processors, Hibernate bytecode enhancement, GraalVM native plugin. Spring Boot 4 uses per-tech test starters (e.g. `spring-boot-starter-webmvc-test`).

`apps/api/ecboot-start/src/main/resources/application.yaml` is nearly empty (`spring.application.name: ecboot`) — runtime configuration is still to be defined.

### Frontends

- **web / admin**: TanStack Start (full-stack React framework with SSR) using file-based routing — routes live in `src/routes/`, and the router config in `src/router.tsx` consumes the generated `src/routeTree.gen.ts`. Never hand-edit `routeTree.gen.ts`; add route files and regenerate. Imports use the `#/*` alias which maps to `./src/*` (defined in each `package.json` `imports` field).
- **mobile**: uni-app — one Vue 3 codebase targeting WeChat/Alipay/etc. mini-programs, H5, and native apps. Pages are registered in `src/pages.json`; app-level config in `src/manifest.json`. Platform selection happens entirely through package scripts (`dev:mp-*` / `build:mp-*`).
