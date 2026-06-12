#!/usr/bin/env python3
"""Verify relative Markdown links resolve on disk (stdlib, Python 3.8+)."""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
LINK_RE = re.compile(r"\[([^\]]+)\]\(([^)]+)\)")
MD_GLOBS = ("*.md", "docs/**/*.md", ".github/**/*.md")

def iter_markdown_files():
    seen = set()
    for pattern in MD_GLOBS:
        for path in ROOT.glob(pattern):
            if not path.is_file():
                continue
            resolved = path.resolve()
            if resolved in seen:
                continue
            seen.add(resolved)
            yield path

def link_target(raw):
    target = raw.strip()
    if not target or target.startswith("#"):
        return None
    if "://" in target or target.startswith("mailto:"):
        return None
    if target.startswith("<") and target.endswith(">"):
        target = target[1:-1].strip()
    if "#" in target:
        target = target.split("#", 1)[0]
    return target or None

def main():
    errors = []
    checked = 0
    for md in sorted(iter_markdown_files()):
        text = md.read_text(encoding="utf-8")
        for match in LINK_RE.finditer(text):
            rel = link_target(match.group(2))
            if rel is None:
                continue
            checked += 1
            dest = (md.parent / rel).resolve()
            try:
                dest.relative_to(ROOT)
            except ValueError:
                errors.append("{}: link escapes repo: {}".format(md, rel))
                continue
            if not dest.exists():
                errors.append("{}: missing {}".format(md.relative_to(ROOT), rel))

    if errors:
        print("check_links: FAILED ({} link(s) checked)".format(checked))
        for err in errors:
            print("  - {}".format(err))
        return 1

    print("check_links: OK ({} Markdown file link(s) checked)".format(checked))
    return 0

if __name__ == "__main__":
    sys.exit(main())
