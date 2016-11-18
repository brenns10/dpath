Report Structure
----------------

- Abstract: summarize the topic of the paper.
- Introduction:
  - short introduction to XML: what it is, why it's used, and an example
  - short introduction to XPath: why we want it, what it does, an example, where
    it's used
  - a full description of what this project does, based on these overviews
  - layout of the remainder of the paper
- Full description of the XPath language:
  - Syntax is the least important, but I should include the basic types of
    literals, the QName syntax, how white space is treated.
  - Describe the context object and the concept of axes
  - Give a solid definition of the semantics of a path expression
  - Try to give a reasonable example if I can.
  - Description of approaches for efficient query evaluation (*)
- Full description of the DPath implementation:
  - First, describe the parts of DPath that aren't identical to XPath.
    - Predicates
    - Pound shorthand
    - descendant-or-self::node() vs descendant-or-self::*
  - Subsection on the lexical analyzer and parser. Keep this short. While
    interesting, the part of this that is "implementing a language" is really
    just implementation details. The project is querying.
  - Describe context object briefly
  - An overview of the arithmetic expression language (but again, this isn't
    that important)
  - Describe the Sequence-based approach
    - lazy loading, etc
  - Lots of discussion on the implementation of path expressions via sequences.
    - Should give an overview of each sequence defined
    - Should describe how each axis is implemented
    - Give an example of a path and show a diagram of all the sequences
  - Included functions
- Future Work
- Related Work
  - glob
  - find (there's potential to compare runtimes here)
    - I also believe find has an XPath implementation
    - There's plenty to compare, cite find manual, etc.
  - ls (recursive) + grep, also can do runtime comparisons here
- Conclusion
