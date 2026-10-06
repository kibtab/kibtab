# Sphinx configuration for the Kibtab documentation.
#
# Read the Docs runs this configuration; see .readthedocs.yaml. The source
# is the Markdown project under docs/. The MyST parser reads each .md page
# unchanged, so no page moves to reStructuredText. The Furo theme ships in
# requirements.txt beside sphinx and myst-parser.

import os

# -- Project -----------------------------------------------------------------

project = "Kibtab"
author = "The Kibtab authors"
copyright = "The Kibtab authors"
release = "latest"

# -- General -----------------------------------------------------------------

# The documentation index lives at docs/index.md. root_doc names the master
# document that hosts the toctree.
root_doc = "index"
source_suffix = {".md": "markdown"}

# MyST turns each .md page into a Sphinx document.
extensions = ["myst_parser"]

# Each heading gets a stable "#slug" anchor at depth 3.
myst_heading_anchors = 3

# The build output never counts as a source page.
exclude_patterns = ["_build"]

# Every page sits outside a toctree at the moment. The warning is expected
# for now, and each page still builds and stays reachable by link and search.
suppress_warnings = ["toc.not_included"]

# -- HTML output (Furo theme) ------------------------------------------------

html_theme = "furo"

# Read the Docs sets the canonical URL. A local build uses "/".
html_baseurl = os.environ.get("READTHEDOCS_CANONICAL_URL", "/")

# Sphinx builds search into the HTML builder, so no extension is needed.
html_search_language = "en"
