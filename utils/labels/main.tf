provider "github" {
  token        = ""
  organization = "petomalina"
}

module "service_account" {
  source = "git::git@github.com:flowup/terraform.git//github_labels?ref=master"
  repo   = "krane"

  general = {
    "bug" = "d73a4a"
    "regression" = "c11325"
    "feature" = "6cc138"
    "enhancement" = "b0ea8c"
  }
  layers = {
    "cli" = "1d76db"
    "operator" = "f97625"
    "kubernetes" = "FFD700"
  }
  feats = [
    "tester",
    "canary-manager",
    "metric-collector",
  ]
  utils = [
    "documentation",
    "tests",
  ]
  estimates = []
  frequencies = [
    "low",
    "medium",
    "high",
  ]
  severities = [
    "trivial",
    "minor",
    "major",
    "critical",
  ]
}
