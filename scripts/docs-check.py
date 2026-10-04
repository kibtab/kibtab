#!/usr/bin/env python3
"""Check the Kibtab documentation for style, links, and section references.

The script checks four things:

* each `section N` reference resolves to a heading in the file that it names
* each relative Markdown link points at a file that exists
* each prose sentence stays under the word limit
* each prose sentence avoids the banned words

The canonical documents are excluded from every check. Copy them as written.
The vendored skill is excluded too. The skill is pinned to a commit.

Usage:

    python3 scripts/docs-check.py plan.md AGENTS.md README.md docs
"""

import os
import re
import sys

CANONICAL = {
    "LICENSE",
    "CODE_OF_CONDUCT.md",
}

SKIP_DIRS = {".git", "node_modules", "vendor", "skills"}

BANNED = (
    "should",
    "would",
    "may",
    "might",
    "can not",
    "don't",
    "doesn't",
    "isn't",
    "it's",
    "we're",
    "you're",
    "can't",
    "won't",
)

MAX_WORDS = 20

FENCE = re.compile(r"^\s*(```|~~~)")
LINK = re.compile(r"\[([^\]]*)\]\(([^)]+)\)")
SECTION = re.compile(r"([A-Za-z0-9_./-]+\.md)[^.\n]{0,40}?section\s+(\d+)")
CODE_SPAN = re.compile(r"`[^`]*`")
COMMENT = re.compile(r"<!--.*?-->", re.DOTALL)
SENTENCE = re.compile(r"(?<=[.!?:])\s+")


def markdown_files(roots):
    """Yield each first-party Markdown file at or below each root."""
    for root in roots:
        if os.path.isfile(root):
            yield root
            continue
        if not os.path.isdir(root):
            print(f"missing path: {root}")
            continue
        for base, dirs, names in os.walk(root):
            dirs[:] = [d for d in dirs if d not in SKIP_DIRS]
            for name in sorted(names):
                if name.endswith(".md"):
                    yield os.path.join(base, name)


def headings(text):
    """Return each numbered section number in the text."""
    found = set()
    for line in text.splitlines():
        match = re.match(r"^#{1,6}\s+(\d+)\.", line)
        if match:
            found.add(match.group(1))
    return found


def read(path):
    """Return the text of the file."""
    with open(path, encoding="utf-8") as handle:
        return handle.read()


def is_checked(path):
    """Return True when the path is a canonical document or the skill."""
    parts = set(os.path.normpath(path).split(os.sep))
    return path in CANONICAL or "skills" in parts


def check_sections(files, failures):
    """Check each `section N` reference against the headings of its file."""
    cache = {}

    def load(path):
        if path not in cache:
            if not os.path.isfile(path):
                cache[path] = None
            else:
                cache[path] = read(path)
        return cache[path]

    for path in files:
        text = COMMENT.sub("", read(path))
        for line in text.splitlines():
            for target, number in SECTION.findall(CODE_SPAN.sub("", line)):
                resolved = os.path.normpath(
                    os.path.join(os.path.dirname(path), target)
                )
                body = load(resolved)
                if body is None:
                    failures.append(f"{path}: no file for ref: {target}")
                elif number not in headings(body):
                    failures.append(f"{path}: {target} has no section {number}")


def check_links(files, failures):
    """Check each relative Markdown link points at a file that exists."""
    for path in files:
        in_fence = False
        in_comment = False
        for number, line in enumerate(read(path).splitlines(), start=1):
            if FENCE.match(line):
                in_fence = not in_fence
                continue
            if in_fence:
                continue
            if "<!--" in line:
                in_comment = True
            if in_comment:
                if "-->" in line:
                    in_comment = False
                continue
            for _, target in LINK.findall(line):
                if target.startswith(("http://", "https://", "#", "mailto:")):
                    continue
                clean = target.split("#", 1)[0]
                if not clean:
                    continue
                resolved = os.path.normpath(
                    os.path.join(os.path.dirname(path), clean)
                )
                if not os.path.exists(resolved):
                    failures.append(f"{path}:{number}: broken link: {target}")


def prose_lines(path):
    """Yield each prose unit. Skip tables, headings, code, and comments.

    A bullet and its wrapped lines form one unit. The rule counts the words of
    the whole bullet, not the words of each line.
    """
    in_fence = False
    in_comment = False
    buffer = []

    def flush():
        if buffer:
            text = " ".join(buffer).strip()
            del buffer[:]
            return text or None
        return None

    for raw in read(path).splitlines():
        if FENCE.match(raw):
            in_fence = not in_fence
            pending = flush()
            if pending:
                yield pending
            continue
        if in_fence:
            continue
        if "<!--" in raw:
            in_comment = True
        if in_comment:
            if "-->" in raw:
                in_comment = False
            continue
        stripped = raw.strip()
        if not stripped:
            pending = flush()
            if pending:
                yield pending
            continue
        if stripped.startswith(("|", "#", ">", "[")):
            pending = flush()
            if pending:
                yield pending
            continue
        if raw.lstrip().startswith(("*", "-", "+")) or re.match(
            r"^\d+\.\s", stripped
        ):
            pending = flush()
            if pending:
                yield pending
            buffer.append(stripped.lstrip("*-+ ").strip())
            continue
        if buffer and raw.startswith((" ", "\t")):
            buffer.append(stripped)
            continue
        pending = flush()
        if pending:
            yield pending
        buffer.append(stripped)

    pending = flush()
    if pending:
        yield pending


def check_style(files, failures):
    """Check each prose unit for the banned words and for sentence length."""
    for path in files:
        for unit in prose_lines(path):
            bare = CODE_SPAN.sub("", unit)
            for word in BANNED:
                if re.search(rf"\b{re.escape(word)}\b", bare, re.IGNORECASE):
                    failures.append(f"{path}: banned word: {word}")
            for sentence in SENTENCE.split(unit):
                sentence = sentence.strip()
                count = len(sentence.split())
                if count > MAX_WORDS:
                    failures.append(
                        f"{path}: {count} words: {sentence[:60]}"
                    )


def main():
    roots = sys.argv[1:]
    if not roots:
        print(__doc__)
        return 2
    files = [p for p in markdown_files(roots) if not is_checked(p)]
    failures = []
    check_sections(files, failures)
    check_links(files, failures)
    check_style(files, failures)
    for line in failures:
        print(line)
    print(f"checked {len(files)} file(s), {len(failures)} problem(s)")
    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main())