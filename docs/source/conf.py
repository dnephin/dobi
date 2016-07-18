# -*- coding: utf-8 -*-
#
# dobi documentation build configuration file,

import sphinx_bootstrap_theme

# -- General configuration ------------------------------------------------

needs_sphinx = '1.4.5'

extensions = []
templates_path = ['_templates']
source_suffix = '.rst'
master_doc = 'index'

# General information about the project.
project = u'dobi'
copyright = u'2016, Daniel Nephin'
author = u'Daniel Nephin'

version = release = u'0.2'
language = None
exclude_patterns = []

pygments_style = 'sphinx'

todo_include_todos = False


# -- Options for HTML output ----------------------------------------------

html_theme = 'alabaster'
html_theme_path = sphinx_bootstrap_theme.get_html_theme_path()
"""
html_theme_options = {
    'bootswatch_theme': 'yeti',
    'bootstrap_version': "3",

    'navbar_site_name': "Pages",
    'navbar_sidebarrel': False,
    'navbar_pagenav': False,


}
"""
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
