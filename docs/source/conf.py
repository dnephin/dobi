# -*- coding: utf-8 -*-
#
# dobi documentation build configuration file,

# -- General configuration ------------------------------------------------

needs_sphinx = '1.4.5'

templates_path = ['_templates']
source_suffix = '.rst'
master_doc = 'index'

# General information about the project.
project = u'dobi'
copyright = u'2016, Daniel Nephin'
author = u'Daniel Nephin'

version = release = u'0.3-dev'
language = None
exclude_patterns = []

pygments_style = 'sphinx'

todo_include_todos = False


# -- Options for HTML output ----------------------------------------------

import alabaster

html_theme = 'alabaster'
html_theme_path = [alabaster.get_path()]
extensions = ['alabaster', 'sphinx.ext.githubpages']

html_sidebars = {
    '**': [
        'about.html',
        'navigation.html',
        'relations.html',
        'searchbox.html',
    ]
}

html_theme_options = {
    'description': "A build automation tool for Docker applications",
    'github_user': 'dnephin',
    'github_repo': 'dobi',
    'github_button': True,
    'github_type': 'star',
    'github_banner': True,
}

# html_title = u'dobi v0.2'
# html_short_title = None
# html_logo = None
# html_favicon = None
html_static_path = ['_static']
# html_last_updated_fmt = None

# html_domain_indices = True
# html_use_index = True
# html_split_index = False
# html_show_sourcelink = True
# html_show_sphinx = True
html_show_copyright = False 
# html_use_opensearch = ''
# html_file_suffix = None
# html_search_language = 'en'
# html_search_options = {'type': 'default'}
# html_search_scorer = 'scorer.js'
htmlhelp_basename = 'dobidoc'
