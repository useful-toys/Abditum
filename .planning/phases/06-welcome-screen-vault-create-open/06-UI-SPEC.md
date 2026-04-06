---
status: draft
---

# UI Design Contract: Phase 06 - Welcome Screen + Vault Create/Open

## 1. Overview

This document defines the visual and interaction contracts for Phase 06, focusing on the Welcome Screen and the Vault Create/Open flows. It leverages the existing TUI Design System (`tui-design-system-novo.md`) and specific decisions from `06-CONTEXT.md` and `REQUIREMENTS.md`.

## 2. Design System

**Design System Name:** Abditum TUI Design System
**Detected from:** `tui-design-system-novo.md`

## 3. Spacing

The TUI utilizes a character and line-based spacing system.

-   **Horizontal Padding (default):** 4 character units.
-   **Vertical Padding (default):** 1 line unit.
-   **Justification for 1 line unit:** In a character-based TUI, a 1-line vertical unit is the fundamental atom of spacing. Forcing to multiples of 4 would create excessively sparse layouts, reducing information density and violating the compact nature typical of TUIs. This is a deliberate design choice for TUIs where screen real estate is at a premium.
-   **Exception:** `filePickerModal` uses 0 vertical padding (per `06-CONTEXT.md` D-03, justified by vertical space scarcity).
-   **Internal Dialog Padding:** 4 columns horizontal, 1 line vertical (per `tui-design-system-novo.md` -> Dimensionamento e Layout -> Dimensionamento de diálogos).

## 4. Typography

Typography in a TUI environment relies on ANSI attributes rather than pixel-based font sizes.

-   **Available Attributes:**
    -   **Bold:** Universal support. Used for titles, selected cursor, default actions.
    -   **Dim / Faint:** Broad support. Used for disabled items, secondary content.
    -   *Italic*: Partial support. Used for hints, virtual folders, auxiliary texts.
    -   ~~Strikethrough~~: Partial support. Used for items marked for deletion.
-   **Reliability:** `Bold` is the only universally reliable typographic highlight. `Italic` and `Strikethrough` always require visual reinforcement (e.g., symbols or secondary colors).
-   **Combinations:**
    -   `Bold + semantic color`: Title of alert/info modals.
    -   `Dim + strikethrough`: Deleted item, reinforced with `✗`.
    -   `Italic + text.secondary`: Hints and auxiliary texts.

-   **Justification for multiple attributes:** In a character-based TUI, ANSI attributes (`Bold`, `Dim`, `Italic`, `Strikethrough`) serve as the primary means of visual distinction, akin to different font weights and styles in GUIs. The TUI Design System (`tui-design-system-novo.md` -> Tipografia) explicitly defines their semantic roles and fallback behaviors. `Bold` and `Dim/Faint` are considered primary \'weights\' for hierarchical emphasis, while `Italic` and `Strikethrough` are semantic decorators that always require visual reinforcement (e.g., symbols, secondary colors) due to partial terminal support. This approach allows for richer, semantically meaningful communication within the terminal\'s constraints, without violating the principle of visual clarity or over-complicating the typographic hierarchy. The goal is to provide sufficient visual cues, not to mimic GUI font systems with arbitrary limits.

## 5. Color

Colors are defined by their functional roles within the TUI.

-   **Themes:** Tokyo Night (default for examples below) and Cyberpunk. Toggleable via `F12`.
-   **Dominant Surface:** `surface.base` (e.g., Tokyo Night: `#1a1b26`). Used for the entire screen background.
-   **Secondary Surface:** `surface.raised` (e.g., Tokyo Night: `#24283b`). Used for side panels and modal backgrounds.
-   **Accent Colors:**
    -   `accent.primary` (e.g., Tokyo Night: `#7aa2f7`):
        -   Selection bar in lists.
        -   Navigation cursor.
        -   Main action buttons/labels.
        -   Active tab highlight.
        -   Focused borders (e.g., `border.focused`).
        -   Password strength meter fill (for "Good" or "Strong").
    -   `accent.secondary` (e.g., Tokyo Night: `#bb9af7`):
        -   Favorite icon (`★`).
        -   Folder names in file navigation (`filePickerModal`).
        -   Logo gradient.
-   **Semantic Colors:**
    -   `semantic.success` (e.g., Tokyo Night: `#9ece6a`): Success messages, ON states.
    -   `semantic.warning` (e.g., Tokyo Night: `#e0af68`): Alerts, warnings (e.g., weak password, dirty state indicator `•`), destructive dialog borders.
    -   `semantic.error` (e.g., Tokyo Night: `#f7768e`): Error messages, wrong password, destructive action labels.
    -   `semantic.info` (e.g., Tokyo Night: `#7dcfff`): Informational messages.
    -   `semantic.off` (e.g., Tokyo Night: `#737aa2`): OFF states.

## 6. Copywriting

-   **Primary CTA Labels:**
    -   Welcome Screen: `n Novo cofre`, `o Abrir cofre`
    -   Password Input Modals: `Enter Confirmar`
    -   Unsaved Changes Dialog: `Salvar`, `Descartar`, `Voltar`
    -   Overwrite Confirmation: `Sobrescrever`, `Voltar`
    -   Weak Password Confirmation: `Prosseguir`, `Revisar`
    -   Error Recognition Modals: `Enter Fechar`
-   **Empty State Copy:**
    -   Welcome Screen (no vault open): `n Novo cofre    o Abrir cofre` (action hints).
-   **Error State Copy:**
    -   **Incorrect Password:** `✕ Senha incorreta` (MessageManager, 5s TTL).
    -   **Passwords Mismatch:** `✕ As senhas não conferem — digite novamente` (MessageManager, 5s TTL).
    -   **Invalid Vault File (Magic):** `O arquivo selecionado não é um cofre Abditum` (Recognition × Error modal).
    -   **Vault Version Too New:** `Este cofre foi criado por uma versão mais recente do Abditum` (Recognition × Error modal).
    -   **Corrupted Vault File:** `O arquivo está corrompido e não pode ser aberto` (Recognition × Error modal).
    -   **General Save/Create Failure:** `Não foi possível [criar/salvar] o cofre — verifique o caminho e as permissões` (Recognition × Error modal).
-   **Destructive Actions & Confirmation:**
    -   **Overwrite Existing Vault File:**
        -   Prompt: "Sobrescrever arquivo?"
        -   Type: `Confirmation × Destrutivo` dialog.
        -   Options: `Sobrescrever` (default), `Voltar`.
    -   **Discard Unsaved Changes (before Open/Create):**
        -   Prompt: "Alterações não salvas: Deseja salvar antes de sair?"
        -   Type: `Confirmation × Neutro` dialog.
        -   Options: `Salvar` (default), `Descartar`, `Voltar`.
    -   **Quit with Dirty Vault (`Ctrl+Q`):**
        -   Prompt: "Alterações não salvas: Salvar / Descartar / Voltar?"
        -   Type: `Confirmation × Neutro` dialog.
        -   Options: `Salvar` (default), `Descartar`, `Voltar`.
        -   If "Salvar" chosen and external modification detected: "Este cofre foi modificado externamente.\n\nSobrescrever / Salvar como novo arquivo / Cancelar".

## 7. Registry

**Registry:** none
**Reason:** Abditum is a Go TUI application; `shadcn/ui` is not applicable.
