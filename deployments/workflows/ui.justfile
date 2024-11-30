set working-directory := "ui"

[no-cd]
[doc("Install dependencies")]
install-deps:
    @bun i

[no-cd]
[doc("Start development server")]
dev:
    @bun --bun run dev

[no-cd]
[doc("Build production bundle")]
build url="/query": install-deps
    @VITE_GRAPGQL_API_URL={{ url }} bun --bun run build --outDir ../changescout/pkg/ui/dist/dist
