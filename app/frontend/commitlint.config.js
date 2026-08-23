// Matches this repo's existing commit style (see `git log`: feat:, fix:,
// chore:, perf:, ...) — enforced repo-wide via the root .husky/commit-msg
// hook, not just for frontend changes.
export default {
  extends: ['@commitlint/config-conventional'],
}
