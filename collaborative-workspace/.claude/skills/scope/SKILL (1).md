<!--
Adapted from mattpocock/skills (MIT License)
https://github.com/mattpocock/skills/blob/main/skills/productivity/grilling/SKILL.md
-->

---
name: grilling
description: Interview the user relentlessly about a plan or design. Use when the user wants to stress-test a plan before building, or uses any 'grill' trigger phrases.
---

Interview me relentlessly about every aspect of this plan until we reach a shared understanding. 
Walk down each branch of the design tree, resolving dependencies between decisions one by one.
Ask for as much detail as possible so that both the happy path and edge cases are covered.  
For each question:
- provide at least 2 alternative answers
- provide your recommended answer.

Ask the questions one at a time, waiting for feedback on each question before continuing. 
Asking multiple questions at once is bewildering.

Start with functionality (non-technical) questions first. When we reach a shared understanding on non-technical, offer me:
- Continue with the technical grilling session (shared understanding reached when both functional and technical sessions are concluded)
- Stop (shared understanding reached only by functional session)

When we reach a shared understanding, offer to:
- Create /use-case 
- Create /user-story
- Create a standalone Markdown file containing the raw grilling session Q&A