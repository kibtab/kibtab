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

# The root document holds the toctree for the guides and the reference
# pages. The changelog release files, the notices, and the credits sit
# outside it, so Sphinx warns for each of those. The warning is silenced;
# each page still builds and stays reachable by link and search.
suppress_warnings = ["toc.not_included"]

# -- HTML output (Furo theme) ------------------------------------------------

html_theme = "furo"

# Read the Docs sets the canonical URL. A local build uses "/".
html_baseurl = os.environ.get("READTHEDOCS_CANONICAL_URL", "/")

# Sphinx builds search into the HTML builder, so no extension is needed.
html_search_language = "en"

# -- Static assets -----------------------------------------------------------

# The theme patch and the sidebar script ship beside this configuration.
# Sphinx resolves html_css_files against the output _static directory but
# only copies the files listed in html_static_path, so the path must list
# the source folder.
html_static_path = ["_static"]
html_css_files = ["rtd-linkpreviews.css"]
html_js_files = ["sidebar-reveal.js"]

# -- Show the documentation index in the sidebar -----------------------------
#
# The global navigation tree comes from the root document's toctrees, so the
# index page itself is never listed. Furo's own hook cannot add it: Furo
# forbids a toctree that names its master document. A small hook prepends a
# labelled link to the top of the tree on every page instead. It runs after
# Furo computes furo_navigation_tree (priority 800 over Furo's 500), and it
# reuses Furo's caption and toctree markup, so it needs no custom CSS.


def _prepend_index_to_sidebar(app, pagename, templatename, context, doctree):
    tree = context.get("furo_navigation_tree")
    if tree is None:
        return
    root = app.config.root_doc or "index"
    href = context["pathto"](root)
    index_entry = (
        '<p class="caption" role="heading">'
        '<span class="caption-text">Kibtab Documentation</span></p>'
        "<ul><li class=\"toctree-l1\">"
        f'<a class="reference internal" href="{href}">'
        "Documentation index</a></li></ul>"
    )
    context["furo_navigation_tree"] = index_entry + tree


def setup(app):
    app.connect("html-page-context", _prepend_index_to_sidebar, 800)
