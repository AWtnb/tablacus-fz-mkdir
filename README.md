# README

Command line tool for [Tablacus Explorer](https://tablacus.github.io/explorer.html).


1. Traverse the path in reverse order (first the parent folder, then its parent folder, and so on...) and display fuzzy-finder based on the first-found [`dirnames.yaml`](/dirnames.yaml).
1. If a directory with an index number (i.e., the name of the directory begins with a number; `0_plain`, `1_proofed`, ...) exists in working directory, incremented index is inserted to the beginning of the new directory name.

![img](image.png)