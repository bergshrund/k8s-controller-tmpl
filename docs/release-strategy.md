## Release Strategy

- Releases are triggered by pushing a Git tag that follows semantic versioning (e.g., `v1.2.3`).
- Each release builds and publishes:
  - A Helm chart with the version taken from the `Chart.yaml`, matching the Git tag.
- The chart is committed to a GitOps repository watched by **Flux**.
- **Flux** handles the deployment to Kubernetes clusters:
  - A **staging** environment receives updates automatically on each release.
  - The **production** environment is updated manually through a pull request and approval.
- Configuration values are separated per environment using dedicated Helm `values` files.
- **Rollbacks** are supported by reverting the GitOps state (e.g., using Git history).
- All releases are documented via **GitHub Releases** and `CHANGELOG.md`.
