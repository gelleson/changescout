set dotenv-load

[doc("CI/CD")]
mod ci "./deployments/workflows/ci.justfile"

[doc("Generate tasks")]
mod gen "./deployments/workflows/gen.justfile"

[doc("Dev")]
mod dev "./deployments/workflows/dev.justfile"

[doc("Codemerge")]
mod cm "./deployments/workflows/codemerge.justfile"

[doc("UI")]
mod ui "./deployments/workflows/ui.justfile"