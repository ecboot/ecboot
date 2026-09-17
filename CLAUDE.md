# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ECBOOT is an e-commerce platform monorepo (`org.juling.ecboot`) in early scaffold stage: a Java 25 / Spring Boot 4.1.1 backend plus three independent frontend apps. Commit messages are Chinese conventional commits (`feat:`, `chore:`) — match that style. Most backend modules are empty placeholder POMs describing a planned layering; only `start/` contains runnable code.

## Repository Layout

- `start/` — the single runnable Spring Boot application (`EcbootApplication`). Its parent is `spring-boot-starter-parent` directly, so it builds standalone.
- `infrastructure/` — `ecboot-common`, `ecboot-infra-core`: shared libraries (placeholder POMs).
- `services/` — domain service modules: `ecboot-service-user`, `ecboot-service-shop` (placeholder POMs).
- `apps/` — channel API modules: `ecboot-api-user`, `ecboot-api-shop`, `ecboot-api-admin`, `ecboot-api-common` (placeholder POMs).
- `dependencies/` — `ecboot-dependencies` BOM placeholder.
- `frontend/` — `ecboot-web` (storefront), `ecboot-admin` (admin console), `ecboot-mobile` (cross-platform mini-program app).
- `.specify/` — Spec Kit scaffolding (spec-driven development via the `speckit-*` skills); the constitution at `.specify/memory/constitution.md` is still an unfilled template.
- `scripts/codegen.sh` — empty placeholder for future code generation.

## Commands

### Backend

The root `pom.xml` (`ecboot-parent`) is the aggregator + parent (inherits `spring-boot-starter-parent`). **All builds run from the repo root** — single-module builds fail because the BOM in `dependencies/` is resolved from the reactor:

```bash
./mvnw clean package -DskipTests      # full reactor build
./mvnw test -pl start -am             # tests for the runnable app (+ its module deps)
./mvnw spring-boot:run -pl start -am  # run the app (auto-starts compose.yaml services; needs Docker)
```

Dependency rules are enforced by maven-enforcer at the `validate` phase — violations fail the build:

- Layer direction: `start → apps → services → infrastructure`; `ecboot-api-common` must not depend on `services/*`; cycles fail at build time.
- Version arbitration: third-party/framework versions are declared ONLY in `dependencies/pom.xml` (`spring-boot.version` property + BOM imports); plugin versions ONLY in the root `pluginManagement`. Business module POMs carry zero version numbers.
- Per-module ban lists: slot properties `enforcer.banned.1..5` (+ `enforcer.allowed.1` exception), defined per module, defaults in the root POM. A single `<exclude>` does NOT support comma-separated lists.

Local infrastructure is defined in `compose.yaml`: MySQL (db `mydatabase`, user `myuser`/`secret`), Redis, Elasticsearch 9.3.3 (security disabled). Start manually with `docker compose up -d`.

### Frontend

`ecboot-web` and `ecboot-admin` (TanStack Start + React 19 + Tailwind CSS 4 + Vite 8):

```bash
cd frontend/ecboot-web   # or ecboot-admin
npm install
npm run dev              # vite dev server on port 3000 — BOTH projects use 3000; change one if running both
npm run generate-routes  # regenerate src/routeTree.gen.ts (tsr generate)
npm run build
```

`ecboot-mobile` (uni-app, Vue 3):

```bash
cd frontend/ecboot-mobile
npm run dev:h5           # H5 dev server
npm run dev:mp-weixin    # WeChat mini-program dev build
npm run build:mp-weixin  # WeChat mini-program production build
npm run type-check       # vue-tsc --noEmit
```

npm is the package manager (`package-lock.json` committed, no workspace config — each frontend is independent).

## Architecture

### Backend (layering, enforced)

The modules form a modular monolith: `infrastructure/*` (shared libs) ← `services/*` (domain: user, shop) ← `apps/*` (API per channel: user, shop, admin, plus `api-common`) — assembled by `start/` into one deployable. The dependency wiring IS implemented and enforced (see Backend commands above).

The `start/` stack: WebMVC, Security, JPA + Flyway (MySQL), Redis, Elasticsearch, Quartz, Mail, WebSocket, RestClient, Validation, Actuator. Build-time extras: Lombok + `spring-boot-configuration-processor` annotation processors, Hibernate bytecode enhancement, GraalVM native plugin. Spring Boot 4 uses per-tech test starters (e.g. `spring-boot-starter-webmvc-test`).

`start/src/main/resources/application.yaml` is nearly empty (`spring.application.name: ecboot`) — runtime configuration is still to be defined.

### Frontends

- **web / admin**: TanStack Start (full-stack React framework with SSR) using file-based routing — routes live in `src/routes/`, and the router config in `src/router.tsx` consumes the generated `src/routeTree.gen.ts`. Never hand-edit `routeTree.gen.ts`; add route files and regenerate. Imports use the `#/*` alias which maps to `./src/*` (defined in each `package.json` `imports` field).
- **mobile**: uni-app — one Vue 3 codebase targeting WeChat/Alipay/etc. mini-programs, H5, and native apps. Pages are registered in `src/pages.json`; app-level config in `src/manifest.json`. Platform selection happens entirely through npm scripts (`dev:mp-*` / `build:mp-*`).
