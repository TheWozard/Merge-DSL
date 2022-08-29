# Merge-DSL
[![Test](https://github.com/TheWozard/Merge-DSL/actions/workflows/test.yml/badge.svg)](https://github.com/TheWozard/Merge-DSL/actions/workflows/test.yml)
[![Lint](https://github.com/TheWozard/Merge-DSL/actions/workflows/lint.yml/badge.svg)](https://github.com/TheWozard/Merge-DSL/actions/workflows/lint.yml)

A YAML based DSL for merging multiple partial documents together.

WIP Names:
- Construct (Constructing a final document)
- Conscript (Construct + Script)
- Zipper (Zips together documents)
- Weaver (Weaves together documents)

## Table of Contents

- [Merge-DSL](#merge-dsl)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
    - [Definition](#definition)

## Overview
DSL is comprised of 3 main components
| Name       | Description                                                       |
| ---------- | ----------------------------------------------------------------- |
| Definition | Declaration of the final output document and how to handle merges |
| Rules      | Declaration of special rules to be applied during traversal       |
| Documents  | Raw input data to be merged                                       |

### Definition
The basis for any merge operation. Defines the schema of the final document and any merge properties.