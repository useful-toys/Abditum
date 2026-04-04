---
description: Commit uncommitted TUI documentation files individually with conventional commits
---

<objective>
Find all uncommitted markdown documentation files related to TUI (specifications, design system, architecture) and commit each individually with a conventional commit message that explains the modifications.

This command:
- Detects uncommitted `.md` files in the project root and relevant directories
- Analyzes the diff of each file to understand what changed
- Generates an appropriate conventional commit message based on the content changes
- Commits each file separately for clean git history
</objective>

<execution_context>
No external workflow files needed. This command is self-contained.
</execution_context>

<context>
No arguments needed. The command discovers files automatically.
</context>

<process>
1. **Discover uncommitted markdown files:**
   - Run `git status --porcelain` to find all untracked and modified files
   - Filter for `.md` files that are TUI-related:
     - Files matching `tui-*.md` in project root
     - Ingnore all files in `docs/` directory
     - Ignore any other `.md` documentation files
   - Exclude files that should not be committed (e.g., files in `.gitignore`)

2. **If no uncommitted markdown files found:**
   - Report: "No uncommitted TUI documentation files found."
   - Exit

3. **For each uncommitted markdown file:**
   a. **Analyze the changes:**
      - If file is new (untracked): note it's a new file
      - If file is modified: Run `git diff -- <file>` to see the changes
      - Read the diff output to understand what changed
   
   b. **Determine conventional commit type:**
      - `docs:` for documentation additions/updates (most common)
      - `style:` if only formatting changes
   
   c. **Generate commit message:**
      - Format: `<type>(tui): <brief description>`
      - Body should summarize key changes (2-3 bullet points max)
      - Example:
        ```
        docs(tui): add component specifications for modal dialogs
        
        - Define modal dialog states and transitions
        - Add keyboard navigation patterns
        - Specify accessibility requirements
        ```
   
   d. **Commit the file:**
      - Run: `git add <file>`
      - Run: `git commit -m "<commit message>"`
      - Report success with the commit hash

4. **Summary:**
   - List all files committed
   - Show commit messages for review
</process>