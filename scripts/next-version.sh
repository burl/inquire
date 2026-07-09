#!/usr/bin/env bash
# Print the next semantic version tag from conventional commits since the last tag.
#
# Bump rules (highest wins across the range):
#   - breaking change  -> major  (type! / type(scope)! / BREAKING CHANGE: in body)
#   - feat             -> minor
#   - anything else    -> patch
#
# Type normalization (scope ignored):
#   fix(foo): bar              -> fix
#   fix(bar)!: changed ...     -> fix!  (breaking)
#   fix!: ...                  -> fix!  (breaking)
#
# Usage: scripts/next-version.sh [options] [ref]
#   ref defaults to HEAD (CI on master uses HEAD == master).
#
# Options:
#   -n, --dry-run, -v, --verbose  Explain each commit's bump on stderr.
#                                 The version is still printed alone on stdout.
#   -h, --help                    Show this help.
set -euo pipefail

dry_run=0
ref=HEAD

usage() {
  sed -n '2,20p' "$0" | sed 's/^# \?//'
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -n|--dry-run|-v|--verbose)
      dry_run=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    -*)
      printf 'unknown option: %s\n' "$1" >&2
      exit 2
      ;;
    *)
      ref="$1"
      shift
      ;;
  esac
done

level_name() {
  case "$1" in
    2) printf 'major' ;;
    1) printf 'minor' ;;
    *) printf 'patch' ;;
  esac
}

log() {
  if [[ "$dry_run" -eq 1 ]]; then
    printf '%s\n' "$*" >&2
  fi
}

# Prefer numeric v* tags (this module ships as inquire/v2).
latest="$(git tag -l 'v[0-9]*' --sort=-v:refname | head -1 || true)"

if [[ -z "$latest" ]]; then
  log "base:    (none)"
  log "range:   (no prior tag)"
  log "result:  default first release"
  printf 'v2.0.0'
  exit 0
fi

ver="${latest#v}"
# Strip optional pre-release / build metadata for arithmetic.
ver="${ver%%-*}"
ver="${ver%%+*}"
IFS=. read -r major minor patch _ <<< "$ver"
major="${major:-0}"
minor="${minor:-0}"
patch="${patch:-0}"

# bump level: 0=patch, 1=minor, 2=major
level=0
commit_count=0

log "base:    ${latest}"
log "range:   ${latest}..${ref}"
log "commits:"

# %h = short hash; %B = subject + body; NUL separates fields/commits.
while true; do
  IFS= read -r -d '' hash || break
  IFS= read -r -d '' msg || break

  # Trim leading blank lines; first non-empty line is the subject.
  subject=""
  body=""
  while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ -z "$subject" ]]; then
      [[ -z "$line" ]] && continue
      subject="$line"
    else
      body+="$line"$'\n'
    fi
  done <<< "$msg"

  [[ -z "$subject" ]] && continue
  commit_count=$((commit_count + 1))

  # Per-commit classification (independent of running max).
  commit_level=0
  reason="non-conventional (patch)"
  norm=""

  # Regex must live in a variable: unescaped ')' inside [[ =~ ]] ends the expression.
  conv_re='^([a-zA-Z]+)(\([^)]*\))?(!)?:'
  if [[ "$subject" =~ $conv_re ]]; then
    ctype="${BASH_REMATCH[1]}"
    bang="${BASH_REMATCH[3]:-}"
    if [[ -n "$bang" ]]; then
      norm="${ctype}!"
      commit_level=2
      reason="breaking (${norm})"
    elif [[ "$ctype" == "feat" ]]; then
      norm="feat"
      commit_level=1
      reason="feat (minor)"
    else
      norm="$ctype"
      commit_level=0
      reason="${ctype} (patch)"
    fi
  fi

  # Footer / body marker wins over subject type (Conventional Commits).
  if [[ "$body" =~ BREAKING[[:space:]]CHANGE: ]]; then
    commit_level=2
    if [[ -n "$norm" ]]; then
      reason="BREAKING CHANGE: footer (${norm})"
    else
      reason="BREAKING CHANGE: footer"
    fi
  fi

  if [[ "$commit_level" -gt "$level" ]]; then
    level=$commit_level
  fi

  # Truncate long subjects for the report.
  display="$subject"
  if [[ ${#display} -gt 72 ]]; then
    display="${display:0:69}..."
  fi
  log "  ${hash}  $(level_name "$commit_level")  ${reason}"
  log "            ${display}"
done < <(git log "${latest}..${ref}" --format='%h%x00%B%x00' 2>/dev/null || true)

if [[ "$commit_count" -eq 0 ]]; then
  log "  (none — defaulting to patch)"
fi

base_tag="v${major}.${minor}.${patch}"
case "$level" in
  2)
    major=$((major + 1))
    minor=0
    patch=0
    ;;
  1)
    minor=$((minor + 1))
    patch=0
    ;;
  *)
    patch=$((patch + 1))
    ;;
esac
next="v${major}.${minor}.${patch}"

log "bump:    $(level_name "$level")  (${base_tag} -> ${next})"
log "result:  ${next}"

printf '%s' "$next"
