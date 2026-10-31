# Architectural Decision Record (ADR) & Design Specifications

## 1. Context & Business Problem
Enterprises operating across multiple cloud environments face three critical friction points:
1. **Cloud-Provider Lock-in:** Fragmented deployment configurations tailored to single providers.
2. **Security Vulnerabilities Slipping into Production:** Containers deployed without shift-left static analysis and vulnerability admission gating.
3. **Configuration Drift:** Manual operator tweaks (`kubectl edit`, Azure Portal, AWS Console) creating cluster inconsistencies.

## 2. Decision Matrix
| Area | Technology Selected | Alternatives Considered | Rationale |
| :--- | :--- | :--- | :--- |
| **GitOps Engine** | Argo CD (ApplicationSet) | Flux v2, Jenkins X | Industry standard, declarative multi-cluster ApplicationSet matrix generator, rich UI for enterprise visualization. |
| **Container Runtime Security** | Aquasecurity Trivy | Clair, Anchore | Fastest vulnerability database updates, native GitHub Actions SARIF integration, non-root Distroless compatibility. |
| **Microservice Image Base** | Google Distroless Debian 12 | Alpine, Ubuntu Minimal | Complete absence of package managers and shells eliminates 100% of standard OS CVEs and prevents privilege escalation attacks. |
| **Secret Management** | Secrets Store CSI Driver + Cloud Key Vaults (IRSA / Workload Identity / Azure Key Vault) | Sealed Secrets, HashiCorp Vault | Native integration with cloud providers' Managed Identities eliminates long-lived static API tokens. |

## 3. Threat Modeling & Zero-Trust Defense
* **Shift-Left SAST Gate:** SonarQube analyzes code smells, OWASP Top 10 vulnerabilities, and enforces Quality Gate before container compilation.
* **Trivy Vulnerability Gate:** Fails the CI pipeline with exit code 1 if any `CRITICAL` or `HIGH` unfixed CVEs are identified.
* **Non-Root & Read-Only Root Filesystem:** Pods enforce `runAsNonRoot: true`, user `65532:65532`, and `readOnlyRootFilesystem: true`.
* **Least-Privilege IRSA / Workload Identity:** Kubernetes ServiceAccounts bind dynamically to IAM roles via OIDC tokens, avoiding stored cloud credentials.
