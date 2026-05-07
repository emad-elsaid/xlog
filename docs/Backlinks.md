# Automatic Backlinks

Backlinks are XLog's most powerful knowledge-base feature. They automatically create bidirectional connections between pages, turning your markdown files into an interconnected knowledge graph.

## What Are Backlinks?

When you mention a page name in your markdown, XLog automatically:

1. **Creates a clickable link** to that page
2. **Shows on the target page** that you linked to it (backlink)
3. **Builds a knowledge graph** of relationships between concepts

This happens automatically - no manual link creation or management required.

## How Backlinks Work

### Automatic Page Linking

Simply mention a page name in your text (using the page filename without extension):

```markdown
I'm learning about [[Neural Networks]] and [[Machine Learning]].

This concept relates to [[Graph Theory]].
```

XLog automatically converts these to links. If the target page doesn't exist yet, XLog creates a placeholder that you can click to create the page.

### Bidirectional Relationships

Unlike traditional hyperlinks (which only go one way), backlinks create two-way connections:

**On your current page:**
```markdown
[[Neural Networks]] are inspired by biological neurons.
```

**On the Neural Networks page:**
XLog automatically shows:
```
Backlinks:
- Introduction to Machine Learning (this page links here)
```

This reveals connections you might not have consciously created, helping you discover relationships in your knowledge.

## Syntax

XLog supports multiple linking syntaxes:

### Double Bracket Syntax (Recommended)

```markdown
[[Page Name]]
[[Subdirectory/Page Name]]
```

This is the most common syntax, used by Obsidian, Roam Research, and other knowledge base tools.

### Bare Page Names

If you mention a page name in text without brackets, XLog's autolink extension can convert it to a link:

```markdown
I wrote about this in Machine Learning Notes last week.
```

Configure this via the autolink pages extension.

## Examples

### Personal Wiki

```markdown
# Database Design Patterns

Common patterns in [[Relational Databases]]:

- [[Normalization]] - Reducing redundancy
- [[Indexing Strategies]] - Performance optimization
- [[Connection Pooling]] - Resource management

Related: [[SQL Optimization]], [[NoSQL Databases]]
```

Each mentioned page automatically links, and those pages will show "Database Design Patterns" in their backlinks.

### Research Notes

```markdown
# Paper: Attention Is All You Need (2017)

Introduces the [[Transformer Architecture]] for NLP tasks.

Key innovation: [[Self-Attention Mechanism]] replaces recurrence.

Builds on: [[Sequence-to-Sequence Models]]
Leads to: [[BERT]], [[GPT]], [[T5]]

Related concepts: [[Attention Mechanisms]], [[Neural Machine Translation]]
```

This creates a research paper network showing how ideas connect across papers.

### Learning Journal

```markdown
# 2026-05-07 - Learning Log

Today I learned about [[React Hooks]].

Connected to previous learning:
- [[JavaScript Closures]] (hooks use closures internally)
- [[Functional Programming]] (hooks encourage functional style)
- [[Component Lifecycle]] (hooks replace lifecycle methods)

Questions for tomorrow: How do [[Custom Hooks]] work?
```

Over time, this builds a learning graph showing how concepts connect in your understanding.

## Backlinks Display

XLog shows backlinks at the bottom of each page:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Backlinks (3):
- Introduction to Machine Learning
- Deep Learning Fundamentals  
- AI Resources Collection
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

Click any backlink to navigate to pages that reference the current page.

## Comparison to Manual Linking

### Traditional Static Generators (Hugo, Jekyll)

**Manual linking:**
```markdown
Check out [my other post](/posts/machine-learning/).
```

**Problems:**
- You must know the exact URL path
- Links are one-directional (target page doesn't know you linked to it)
- Refactoring is painful (broken links everywhere)
- No knowledge graph emerges

### XLog Automatic Backlinks

**Automatic linking:**
```markdown
Check out [[Machine Learning]].
```

**Benefits:**
- Just use the page name, XLog handles the rest
- Bidirectional (both pages know about the connection)
- Rename a page and all links update automatically
- Knowledge graph emerges organically

## Use Cases

### 1. Research and Literature Review

Connect research papers, concepts, and citations:

```markdown
[[Smith et al. 2020]] builds on [[Jones 2018]] by adding [[Attention Mechanisms]].

Contradicts findings in [[Brown 2019]].
```

Backlinks reveal which papers are most influential (most backlinks) and how ideas connect.

### 2. Zettelkasten Note-Taking

Build a Zettelkasten-style knowledge base:

```markdown
# 202605071430 - Permanent Note on Spaced Repetition

[[Spaced Repetition]] improves long-term retention.

Evidence: [[Ebbinghaus Forgetting Curve]]
Tools: [[Anki]], [[SuperMemo]]
Related: [[Active Recall]], [[Interleaving]]

Source: [[Book - Make It Stick]]
```

Each note links to related atomic ideas, building a knowledge network.

### 3. Personal Wiki Development

Create an interconnected personal knowledge base:

```markdown
# Go Programming

[[Go]] is a statically-typed language from [[Google]].

Strengths: [[Concurrency]] via [[Goroutines]], simple syntax
Weaknesses: [[Generics]] (added in Go 1.18)

Resources: [[Effective Go]], [[Go by Example]]
Projects: [[XLog]], [[Docker]], [[Kubernetes]]
```

Backlinks show what you've written about each topic across different contexts.

### 4. Project Documentation

Connect documentation pages organically:

```markdown
# API Authentication

Our API uses [[OAuth 2.0]] for authentication.

Setup: See [[API Setup Guide]]
Security: Review [[API Security Best Practices]]
Troubleshooting: [[Common API Errors]]

Related: [[User Management]], [[Access Control]]
```

Developers can navigate documentation following their mental model, not a rigid hierarchy.

## Benefits Over Manual Linking

### 1. Lower Friction

No need to remember exact paths or URLs. Just mention the page name.

### 2. Discover Connections

Backlinks reveal unexpected relationships:
- "I didn't realize I wrote about this in three different contexts"
- "These concepts are more related than I thought"

### 3. Refactoring Safety

Rename a file? XLog updates all links automatically. No broken links to hunt down.

### 4. Knowledge Graph Emergence

Over time, your backlinks reveal:
- **Hubs** - Pages with many backlinks are central concepts
- **Clusters** - Groups of highly interconnected pages
- **Orphans** - Pages with no backlinks might need connecting

This structure emerges naturally from your writing, not from artificial categorization.

## Backlinks vs. Tags

Both backlinks and hashtags organize knowledge, but differently:

| Aspect | Backlinks | Hashtags |
|--------|-----------|----------|
| **Relationship** | Specific page-to-page | Topic categorization |
| **Creation** | Automatic when mentioning pages | Manual with #tag |
| **Hierarchy** | Network (many-to-many) | Flat categories |
| **Best for** | Connecting specific ideas | Organizing by theme |

**Example:**
```markdown
# Machine Learning Notes

Learning about [[Neural Networks]] and [[Decision Trees]].

Topics: #machine-learning #ai #python
```

- **Backlinks** connect to specific pages (Neural Networks, Decision Trees)
- **Hashtags** categorize broadly (#machine-learning, #ai)

Use both together for powerful organization.

## Configuration

Backlinks are provided by the "Autolink pages" extension (enabled by default).

To disable backlinks:
```bash
xlog -disabled-extensions "autolink_pages"
```

(Not recommended - backlinks are core to XLog's knowledge base experience)

## Performance

Backlinks are computed efficiently:
- Indexed on first page load
- Cached for subsequent requests
- Fast even with thousands of pages

For large knowledge bases (5000+ pages), backlink rendering is still sub-second.

## Best Practices

### 1. Link Liberally

Don't overthink it - link to any related concept. XLog handles the complexity.

### 2. Use Descriptive Page Names

Instead of:
```markdown
[[notes_20260507]]
```

Use:
```markdown
[[Neural Network Architectures - 2026-05-07]]
```

Descriptive names make backlinks more useful.

### 3. Review Backlinks Regularly

Check backlinks on your pages to:
- Discover forgotten connections
- Find pages that should link to each other
- Identify central concepts (many backlinks)

### 4. Create Index Pages

Build index pages that link to many related concepts:

```markdown
# Machine Learning Index

Core Concepts:
- [[Supervised Learning]]
- [[Unsupervised Learning]]
- [[Reinforcement Learning]]

Algorithms:
- [[Linear Regression]]
- [[Neural Networks]]
- [[Decision Trees]]
```

These become hub pages with many backlinks.

## Limitations

### 1. No Partial Matches

Links require exact page name matches. These won't link:
```markdown
[[Machine Learning]]  ← Links to "Machine Learning.md"
Machine Learning      ← Won't auto-link (unless autolink extension configured)
machine learning      ← Won't link (case-sensitive)
```

### 2. Cross-Directory Links Need Full Path

```markdown
[[Subdirectory/Page Name]]  ← Required for pages in subdirectories
[[Page Name]]                ← Only finds pages in current directory
```

### 3. No Link Previews

XLog shows backlinks but not link previews (hovering to see content). This is intentional - encourages clicking and exploring.

## See Also

- [Hashtags](Hashtags.md) - Complementary organization with #tags
- [Search](Search.md) - Find content across your knowledge base
- [Digital Gardening](digital gardening.md) - Philosophy behind backlinks
- [official extensions](official extensions.md) - Autolink pages extension details
