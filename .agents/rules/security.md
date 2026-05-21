# Security and Development Rules

## 1. Access Restrictions (Sensitive Files)
* **DO NOT** attempt to read, write, modify, or list actual `.env` files or any other files containing active API keys, tokens, or credentials (e.g., `.env`, `.env.local`, `.env.production`). Note that `.env.example` is **NOT** blocked and must be updated to track configuration schemas.
* If any configuration change requires adding or editing a secret, **ask the user to do it manually**. Do not attempt to bypass this via shell commands (e.g., `cat`, `echo`, `cp`).

## 2. Environment Variable Configuration
* When updating or adding environment variables to any application, you **MUST** document them in the root `.env.example` file with sensible placeholder values and a brief description of their purpose.

## 3. Git Workflow Rules
* **No Direct Commits to Main**: Do not commit or push changes directly to the `main` branch.
* **Feature Branching**: Always create a new descriptive branch (e.g. `feat/feature-name`, `fix/bug-name`, `docs/doc-name`) off of the latest `main` branch for any work.
* **Pull Request Flow**:
  * Commit all code changes onto the feature branch.
  * Push the branch to the remote repository.
  * Suggest or instruct the user to open a Pull Request (PR) to merge the feature branch into `main`.
  * Never merge the branch into `main` directly.
  * **Resolving Review Comments**: After pushing commits that address PR review comments, leave a summary comment on the Pull Request using the GitHub CLI:
    `gh pr review <pr_number> --comment --body "Resolved feedback: <summary of fixes>"`
