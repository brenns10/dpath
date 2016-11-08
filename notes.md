Notes on XPath Syntax
=====================

So far I have stripped several features from the XPath syntax. Likely, more will
be stripped as I go.

First, I have removed "for", "if", quantifiers. These are major modifications to
the expression syntax but they seem to be XQuery features, not XPath.

Next, I have removed range expressions (5 to 10) because they will likely not be
useful for directory traversal, and they would waste implementation time.
