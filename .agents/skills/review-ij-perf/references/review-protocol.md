# Review lifecycle and delivery

## Evidence ledger

Keep a private, lightweight ledger:

- Repository/PR, base/head SHAs, eligibility, preview/posting mode.
- Latest completed review by this workflow and its reviewed SHA.
- Issue/doc links and materially unavailable evidence.
- Prior finding → thread ID, evidence, maintainer decision, current status.
- Candidate → changed cause, trigger, consequence, fix, anchor, proof, duplicate check.
- Actual validation results.

Do not post this ledger.

## Follow-ups

1. Load reviews, conversation, inline replies, and thread resolution state when exposed. Read ambiguous threads; outdated anchors alone do not mean fixed.
2. Identify this workflow's last completed review from reviewer identity, review commit ID, and discussion. An arbitrary human comment, pending review, or failed submission is not that baseline.
3. If current head equals the reviewed SHA, stop without posting unless explicitly asked to reconsider.
4. If baseline metadata is unavailable, inspect commit history and review links. Never guess a reviewed SHA or assert no new commits. If it cannot be established, explain the limit and withhold a duplicate submission.
5. If substantive findings exist but this workflow has never completed a review, assess their status and review the full effective diff. Someone else's review must not suppress this one.
6. Compare last-reviewed head to current head and inspect the effective PR diff against its current base. For rebased/force-pushed history, reconcile findings against current code; do not assume old positions/ancestry.
7. Mark prior findings internally: fixed, still open, accepted/declined, superseded, unverifiable. Verify fixes in callers/tests too.
8. Do not duplicate human/agent findings or reopen intentional/won't-fix decisions. If new evidence proves a materially different failure outside an accepted trade-off, explain that distinction before reporting it.
9. Reuse existing threads for necessary updates to this workflow's own findings. Reply once with changed status/evidence; do not repeat "still open" each run. Do not mark human threads resolved or dismiss human reviews.
10. If an independently confirmed earlier blocker from any reviewer still applies, use REQUEST_CHANGES without duplicating its comment. A human finding is not automatically proven; verify it. If all blockers are fixed or explicitly accepted, approve unless new blockers remain.
11. Put the fixed/still-open/new delta first in the user-facing completion, with thread links when useful. Keep the posted review body to its one-sentence outcome.

## Inline anchors

- Fetch all changed filenames before individual patches. Handle pagination/truncation; read full base/head files for context.
- Validate anchors at the reviewed head: exact repository path, line, RIGHT for new code, LEFT for deleted lines when supported.
- Use multi-line ranges only when each line contributes. Prefer 1–5 lines.
- For an unchanged consumer broken by this PR, anchor on the changed cause and cite the consumer with a full-SHA permalink.
- File line numbers are not diff positions. Prefer line/side; derive positions from the exact patch if required.
- One root cause gets one finding. Independent causes needing different fixes get separate findings.
- Never invent anchors. If the provider cannot attach a valid finding, preserve the proposed review and explain the delivery limitation.

## Submit once and verify

Prefer the GitHub connector; discover capabilities. A single review submission should carry:

- Repository and PR number.
- Full reviewed head SHA as commit_id.
- APPROVE or REQUEST_CHANGES.
- The one-sentence review body.
- Validated inline file comments, if any.

Use add_review_to_pr when available, otherwise an authorized equivalent API/CLI. Use structured arguments or a body file to preserve Markdown/newlines.

Immediately before submission:

1. Re-read metadata and recent discussion.
2. Stop if the PR became draft, WIP, closed, or merged.
3. If base/head moved during this run, refresh affected context and anchors; never attach stale approval.
4. Remove newly duplicated findings.
5. Recompute the action, including independently confirmed still-open blockers from the history.
6. Submit once at that head and verify ID/URL, action, commit, and comments via response/read-back.

After a timeout/ambiguous result, read reviews/comments at the exact head before retrying. Never blindly resubmit. Do not substitute a COMMENT review to work around permission or self-review restrictions. Say "review prepared, submission failed" and provide the proposed outcome; never claim completion.

If essential history/action capabilities are unavailable, use an already authorized equivalent only if it fills the gap. Otherwise explain the limitation and withhold approval.

## Preview and skipped cases

Preview: show `Proposed action: APPROVE` or `Proposed action: REQUEST_CHANGES`, the exact one-sentence body, and inline locations/comments. Mark investigation gaps `Review incomplete`, identifying missing evidence and the blocked decision; do not assign a completed-review outcome.

Draft/closed/merged/unchanged-follow-up cases cause no remote writes. A brief explanation in the conversation suffices. A preview is not a submitted-review baseline.
