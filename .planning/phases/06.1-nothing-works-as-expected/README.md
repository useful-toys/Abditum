# Phase 06.1 — Quick Navigation Guide

**Location:** `.planning/phases/06.1-nothing-works-as-expected/`

---

## 📚 Documents Overview

### 1. **RESEARCH.md** — Technical Deep Dive
**For:** Developers, Architects  
**Size:** 19 KB | **Time to Read:** 30-40 minutes

**Contains:**
- Forensic analysis of all 10 gaps
- Specification references with line numbers
- Code snippets showing current implementation
- Problem + Impact for each gap
- Root cause analysis
- Effort estimates by gap

**Sections:**
```
├─ Executive Summary (10 gaps overview)
├─ Fluxo 1 Analysis (2 gaps: error paths)
├─ Fluxo 2 Analysis (3 gaps: password gate, overwrite dialog)
├─ Fluxos 3-5 Analysis (5 gaps: save-and-exit broken)
├─ Summary Table (all gaps at a glance)
├─ Effort Estimation (by severity)
├─ Prioritization Framework (Tier 1/2/3)
└─ Investigation Summary
```

**Key Takeaway:** Complete technical specification of what's broken and why.

---

### 2. **PLAN.md** — Execution Roadmap
**For:** Implementation Teams  
**Size:** 28 KB | **Time to Read:** 45-60 minutes

**Contains:**
- Phase goal (what success looks like)
- Execution strategy (3 tiers, 2 sub-phases)
- Detailed task breakdown (8 main tasks)
- Code examples for each fix
- Testing strategy (20+ test cases)
- Risk assessment with mitigations
- Timeline and effort estimates
- Success criteria and artifacts

**Task Structure:**
```
TIER 1.A: Fluxo 1 Fixes (0.5 hours)
├─ Task 1.A.1: Invalid magic error path fix
└─ Task 1.A.2: File picker retry verification

TIER 1.B: Fluxo 2 Fix #1 (1 hour)
└─ Task 1.B.1: Overwrite dialog handling

TIER 2.A: Fluxo 2 Fix #2 (2-3 hours)
└─ Task 2.A.1: Password strength validation gate

TIER 3: Fluxos 3-5 Exit Flows (4-5 hours)
├─ Task 3.1: Message types definition
├─ Task 3.2: Fluxo 3 & 4 confirmation dialogs
├─ Task 3.3: Wire decision results to handler
├─ Task 3.4: Create saveAndExitFlow (NEW)
├─ Task 3.5: External modification modal
└─ Task 3.6: Integration testing

TOTAL EFFORT: 9-12 hours
```

**Key Takeaway:** Step-by-step execution guide with code examples and test cases.

---

### 3. **SUMMARY.md** — Executive Overview
**For:** Project Managers, Stakeholders  
**Size:** 8 KB | **Time to Read:** 10-15 minutes

**Contains:**
- What was discovered (10 gaps)
- The critical issue (save-and-exit broken)
- Gaps by severity
- Root causes
- Fix strategy
- Timeline (9-12 hours)
- Risk summary
- Success criteria

**Visual Summaries:**
- Severity breakdown table
- Gaps by flow comparison
- File changes list
- Risk matrix

**Key Takeaway:** High-level understanding of problems, scope, and timeline.

---

### 4. **REVIEW.md** — Quality Assurance
**For:** QA Teams, Technical Reviewers  
**Size:** 10 KB | **Time to Read:** 15-20 minutes

**Contains:**
- Documentation completeness checklist
- Technical accuracy verification
- Code reference spot-checks
- Effort estimation validation
- Spec compliance matrix (before/after)
- Quality assessment
- Final sign-off

**Checklists:**
- [x] All 10 gaps documented
- [x] Code locations verified
- [x] Effort estimates realistic
- [x] Compliance matrices complete
- [x] Testing strategy defined

**Key Takeaway:** Confidence that documentation is accurate and complete.

---

## 🎯 Quick Reference by Role

### I'm a Developer and Need to Implement This

**Start here:** PLAN.md → Your assigned task section → Code examples  
**Reference:** RESEARCH.md → Your specific gap section → Root cause  
**Verify:** REVIEW.md → Spec compliance matrix → After-fix compliance

**Timeline:** 
- Read PLAN.md: 45 min
- Implement task: Depends on tier (0.5-3 hours per task)
- Write tests: Included in task time
- Total: 9-12 hours for full phase

---

### I'm a Tech Lead Planning This Work

**Start here:** SUMMARY.md → Quick overview  
**Then:** PLAN.md → Task breakdown → Dependencies  
**Reference:** RESEARCH.md → Gap analysis for discussions  
**Verify:** REVIEW.md → Sign-off section

**Key Decisions:**
1. Execution order: Tier 1 → 2 → 3 (or Tier 1 CRITICAL first)
2. Parallelization: 1.A + 1.B + 2.A can be parallel
3. Resources: ~1 developer, 9-12 hours total
4. Dependencies: Phase 06 must be complete (✓ Done)

---

### I'm a Project Manager

**Start here:** SUMMARY.md → "What We Discovered" + "Timeline Estimado"  
**Understand:** The 2 CRITICAL gaps blocking vault operations  
**Track:** Success criteria in SUMMARY.md  
**Plan:** 9-12 hour sprint (or 2-3 days with breaks)

**Key Metrics:**
- 10 gaps identified and planned
- 2 CRITICAL, 2 HIGH, 6 MEDIUM
- Estimated effort: 9-12 hours
- Expected outcome: 100% spec compliance

---

### I'm Reviewing Quality & Specs

**Start here:** REVIEW.md → Checklists + Sign-off  
**Verify:** RESEARCH.md → Spec references (fluxos.md line numbers)  
**Validate:** PLAN.md → Task definitions + Test cases  
**Spot-check:** Code examples for syntax

**All verified:** ✅ APPROVED FOR EXECUTION

---

## 📋 The 10 Gaps at a Glance

### Fluxo 1: Abrir Cofre (2 gaps)
```
Gap 1.1: Invalid magic file → endFlow() instead of loop back
          File: flow_open_vault.go:120-128
          Fix: 3 lines

Gap 1.2: Corrupted payload → endFlow() instead of loop back
          File: flow_open_vault.go:124-128
          Fix: 3 lines
```

### Fluxo 2: Criar Cofre (3 gaps)
```
Gap 2.1: No password strength validation gate (HIGH)
          File: flow_create_vault.go:62-115 (missing state)
          Fix: ~60 lines (new state machine)

Gap 2.2: Overwrite dialog "Voltar" exits instead of returning
          File: flow_create_vault.go:71-85
          Fix: ~15 lines

Gap 2.3: Extension auto-addition not visible (INFO)
          File: dialogs.go (investigation needed)
          Fix: Verify implementation
```

### Fluxos 3-5: Exit Flows (5 gaps — CRITICAL!)
```
Gap 3.1: No confirmation when exiting with no vault (Fluxo 3)
          File: root.go:243-257
          Fix: ~15 lines (add confirmation dialog)

Gap 3.2: No confirmation when exiting with clean vault (Fluxo 4)
          File: root.go:243-257
          Fix: ~15 lines (add confirmation dialog)

Gap 3.3: "Salvar" button has Cmd: nil (CRITICAL)
          File: root.go:248-253
          Fix: ~110 lines (wire to saveAndExitFlow)

Gap 3.4: No external file modification check (HIGH)
          File: New saveAndExitFlow
          Fix: ~45 lines (check + modal)

Gap 3.5: No atomic save execution (CRITICAL)
          File: New saveAndExitFlow
          Fix: ~60 lines (atomic save logic)

Gap 3.6: No backup communication on failure
          File: New saveAndExitFlow
          Fix: ~15 lines (error messaging)
```

---

## 🔄 Task Dependencies

```
Start
  ↓
[Tier 1.A] Fluxo 1 fixes (0.5 hrs)
  ↓
[Tier 1.B] Fluxo 2 overwrite (1 hr) ← Can parallel with 1.A
  ↓
[Tier 2.A] Password strength (2-3 hrs) ← Can parallel with 1.B
  ↓
[Tier 3.1] Message types (0.5 hrs)
  ↓
[Tier 3.2 + 3.3] Dialogs + wiring (1.5 hrs) ← Can parallel
  ↓
[Tier 3.4] saveAndExitFlow (2-3 hrs)
  ↓
[Tier 3.5] External mod modal (0.5 hrs)
  ↓
[Tier 3.6] Integration tests (1-2 hrs)
  ↓
Done: All 10 gaps fixed, 100% spec compliance

CRITICAL PATH: 9-12 hours sequential
PARALLELIZED:  7-8 hours wall-clock time
```

---

## 📊 Success Metrics

| Metric | Target | How to Verify |
|--------|--------|---------------|
| **Gaps Fixed** | All 10 | RESEARCH.md gaps section |
| **Tests Passing** | 141 + 20+ new | Run test suite |
| **Spec Compliance** | 100% all flows | REVIEW.md compliance matrix |
| **Code Quality** | No errors/warnings | Compiler + linter |
| **Documentation** | Complete | All tasks have DoD |

---

## 🚀 Next Steps

1. **Immediate:** 
   - [ ] Review SUMMARY.md (10 min)
   - [ ] Review PLAN.md task breakdown (30 min)
   
2. **Before Starting:**
   - [ ] Assign developers to tasks
   - [ ] Set up testing environment
   - [ ] Review RESEARCH.md for context (optional)
   
3. **During Execution:**
   - [ ] Follow PLAN.md task sequence
   - [ ] Use RESEARCH.md for technical reference
   - [ ] Track progress against timeline
   - [ ] Run tests after each task
   
4. **After Completion:**
   - [ ] Verify all 10 gaps fixed
   - [ ] Confirm 141+ tests passing
   - [ ] Validate spec compliance (REVIEW.md matrix)
   - [ ] Update ROADMAP.md with Phase 06.1 status

---

## 📞 Document Contents by Question

**"What's broken?"** → SUMMARY.md: "Gaps by Severity"  
**"Why is it broken?"** → RESEARCH.md: "Gap [X] - Root Cause"  
**"How do I fix it?"** → PLAN.md: "Task [Y]" + Code examples  
**"Is it accurate?"** → REVIEW.md: "Spot-Check Results"  
**"How long will it take?"** → SUMMARY.md: "Timeline" or PLAN.md: "Task Timeline"  
**"What's spec-compliant?"** → REVIEW.md: "Spec Compliance Matrix"

---

## 📝 File Sizes & Complexity

| Document | Size | Complexity | Best For |
|----------|------|-----------|----------|
| SUMMARY.md | 8 KB | Low | Quick overview |
| RESEARCH.md | 19 KB | Medium | Technical deep-dive |
| PLAN.md | 28 KB | High | Implementation guide |
| REVIEW.md | 10 KB | Medium | Quality assurance |
| **Total** | **65 KB** | **Medium** | **Complete package** |

---

**Ready to begin Phase 06.1 execution?**

Choose your entry point above and start reading! 🚀
