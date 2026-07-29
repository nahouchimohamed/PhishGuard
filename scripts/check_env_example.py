#!/usr/bin/env python3
"""Ensure .env.example keys look sane and are referenced (stdlib, Python 3.8+)."""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
ENV_PATH = ROOT / ".env.example"
KEY_RE = re.compile(r"^([A-Za-z_][A-Za-z0-9_]*)=(.*)$")

# Keys that are intentionally documentation-only or optional placeholders.
OPTIONAL_OR_DOC_ONLY = {
    "LOCALSTACK_AUTH_TOKEN",
    "AWS_ENDPOINT_URL",
}

SCAN_GLOBS = (
    "deploy/**/*",
    "backend/**/*",
    "apps/web/**/*",
    "infra/terraform/**/*",
    "docs/**/*",
    "Makefile",
    "README.md",
    "CONTRIBUTING.md",
    "SECURITY.md",
)

def parse_env_example():
    keys = []
    for line_no, raw in enumerate(ENV_PATH.read_text(encoding="utf-8").splitlines(), 1):
        line = raw.strip()
        if not line or line.startswith("#"):
            # Allow commented optional keys to be discovered for docs.
            if line.startswith("#") and "=" in line:
                body = line.lstrip("#").strip()
                m = KEY_RE.match(body)
                if m:
                    OPTIONAL_OR_DOC_ONLY.add(m.group(1))
            continue
        m = KEY_RE.match(line)
        if not m:
            raise SystemExit("check_env_example: bad line {}: {!r}".format(line_no, raw))
        keys.append(m.group(1))
    return keys

def collect_corpus():
    chunks = []
    for pattern in SCAN_GLOBS:
        for path in ROOT.glob(pattern):
            if not path.is_file():
                continue
            if path.name == ".env.example":
                continue
            # Skip bulky / binary-ish paths.
            if any(part in {"node_modules", ".git", "bin", "build"} for part in path.parts):
                continue
            try:
                chunks.append(path.read_text(encoding="utf-8"))
            except (UnicodeDecodeError, OSError):
                continue
    return "\n".join(chunks)

def main():
    if not ENV_PATH.is_file():
        print("check_env_example: missing .env.example")
        return 1

    keys = parse_env_example()
    if len(keys) != len(set(keys)):
        dupes = sorted({k for k in keys if keys.count(k) > 1})
        print("check_env_example: duplicate keys: {}".format(", ".join(dupes)))
        return 1

    corpus = collect_corpus()
    missing = []
    for key in keys:
        if key in OPTIONAL_OR_DOC_ONLY:
            continue
        if key not in corpus:
            missing.append(key)

    if missing:
        print("check_env_example: FAILED — keys not referenced outside .env.example:")
        for key in missing:
            print("  - {}".format(key))
        return 1

    print("check_env_example: OK ({} active key(s))".format(len(keys)))
    return 0

if __name__ == "__main__":
    sys.exit(main())
