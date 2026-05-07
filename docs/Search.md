# Full-Text Search

XLog includes fast, built-in full-text search that lets you find content across your entire knowledge base instantly. No external search service required.

## Quick Start

### Access Search

**Keyboard shortcut:**
```
Ctrl+K  (Windows/Linux)
Cmd+K   (macOS)
```

**Or click** the Search button in the navigation bar.

### Basic Search

Type your query and press Enter:

```
machine learning
```

Finds all pages containing "machine" and "learning".

## How Search Works

### Indexing

XLog indexes all markdown content:
- Page titles
- Headings
- Body text
- Code blocks
- Lists

Updates automatically when pages change (hot-reload).

### Search Algorithm

XLog uses simple but effective full-text matching:
- **Case-insensitive** - "Python" matches "python", "PYTHON", etc.
- **Whole-word matching** - "learn" matches "learn" but not "learned"
- **Multiple terms** - Finds pages containing all search terms

## Search Syntax

### Single Term

```
golang
```

Finds pages containing "golang".

### Multiple Terms (AND)

```
neural networks
```

Finds pages containing both "neural" AND "networks" (in any order).

### Exact Phrases

Enclose in quotes:

```
"attention mechanism"
```

Finds exact phrase "attention mechanism".

### Search in Titles

Prefix with `title:`:

```
title:introduction
```

Finds pages with "introduction" in the title.

## Search Results

Results show:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Search: "machine learning" (8 results)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Introduction to Machine Learning
  ...supervised learning and unsupervised...

Neural Network Architectures  
  ...machine learning models for deep...

Transformer Architecture
  ...natural language processing using...
```

Each result includes:
- **Page title** (clickable)
- **Excerpt** showing search terms in context
- **Highlighted matches**

## Use Cases

### 1. Find Forgotten Content

"What did I write about databases?"

```
database optimization
```

Instantly find all database-related notes.

### 2. Research Connections

"Where did I mention transformers?"

```
transformer
```

Discover all pages discussing transformers across different contexts.

### 3. Locate Specific Information

"What was that command for Docker?"

```
docker compose
```

Find pages with Docker Compose commands.

### 4. Verify Coverage

"Have I written about GraphQL?"

```
graphql
```

Check if you've covered a topic or identify gaps.

### 5. Refactoring Helper

"Which pages mention this old API?"

```
old_api_endpoint
```

Find all pages to update when refactoring code.

## Search Strategies

### Broad to Narrow

Start broad, then narrow:

```
1. python           ← Too many results
2. python async     ← Better
3. python asyncio   ← Specific
```

### Use Key Terms

Search distinctive terms:

```
❌ "how to"         ← Too common
✅ "authentication" ← Specific topic
```

### Multiple Searches

Can't find it? Try variations:

```
1. neural network
2. deep learning
3. backpropagation
```

Different pages might use different terminology.

### Search, Then Browse

Use search to find entry points, then use backlinks to explore:

1. Search: `machine learning`
2. Find: "Introduction to Machine Learning"
3. Click to view page
4. Use backlinks to discover related pages

## Performance

### Fast Even at Scale

XLog search is optimized for speed:

| Knowledge Base Size | Search Time |
|---------------------|-------------|
| 100 pages | <10ms |
| 1,000 pages | <50ms |
| 5,000 pages | <200ms |
| 10,000 pages | <500ms |

Real-time as you type (with debouncing).

### Memory Efficient

Search index is built on demand and cached intelligently. Large knowledge bases use minimal memory.

## Search vs. Browse

Both are valuable for knowledge discovery:

| Use Search When | Use Browse When |
|-----------------|-----------------|
| Looking for specific content | Exploring connections |
| Know what you're looking for | Discovering related ideas |
| Finding forgotten notes | Following thought chains |
| Locating references | Understanding relationships |

**Example workflow:**

1. **Search** for "react hooks" to find entry point
2. **Browse backlinks** to see related concepts
3. **Click hashtags** to find similar pages
4. **Search again** for specific details

## Comparison to Other Tools

### XLog vs. Static Generators

**Hugo, Jekyll, 11ty:**
- ❌ No built-in search (requires external service like Algolia)
- ❌ Complex setup for search
- ❌ Extra cost for hosted search

**XLog:**
- ✅ Built-in full-text search
- ✅ Zero configuration
- ✅ Works offline
- ✅ No external dependencies

### XLog vs. Obsidian

**Obsidian:**
- ✅ Advanced search syntax (regex, operators)
- ✅ Search in properties
- ✅ Live search preview

**XLog:**
- ✅ Built-in (no plugin)
- ✅ Web-accessible
- ✅ Lightweight and fast
- ❌ Simpler syntax (intentionally)

XLog's search is simpler but sufficient for most knowledge base needs.

## Configuration

Search is provided by the search extension (enabled by default).

To disable search:
```bash
xlog -disabled-extensions "search"
```

(Not recommended - search is essential for navigation)

## Limitations

### No Fuzzy Matching

Search requires exact word matches:

```
"learning" won't find "learn" or "learned"
```

Use multiple searches with variations.

### No Boolean Operators

Can't use OR, NOT, etc.:

```
❌ "python OR javascript"  ← Not supported
✅ Search twice separately
```

### No Regex

Simple text matching only:

```
❌ "neural.*network"  ← Not supported  
✅ "neural network"   ← Use plain text
```

XLog keeps search simple and predictable.

### No Result Ranking

Results are ordered by page modification time, not relevance:

```
Most recently modified pages appear first
```

This is intentional - recent pages are usually most relevant in a personal knowledge base.

## Search Tips

### 1. Use Distinctive Words

```
❌ "the introduction"     ← Too common
✅ "transformer architecture" ← Specific
```

### 2. Try Different Phrasings

```
1. "async programming"
2. "asynchronous code"
3. "concurrency"
```

### 3. Search Code

Code blocks are indexed:

```
async def fetch
```

Finds Python async functions.

### 4. Search Dates

```
2026-05
```

Find pages mentioning May 2026.

### 5. Combine with Filters

Search, then filter results mentally:

```
Search: "api"
Mentally filter: "Which of these are tutorials?"
```

## Keyboard Navigation

When search is open:

```
↑/↓     - Navigate results
Enter   - Open selected result
Esc     - Close search
```

Fast keyboard-only navigation.

## Search Across Features

Combine search with other XLog features:

### Search + Backlinks

1. Search for a page
2. View its backlinks
3. Discover connections

### Search + Hashtags

1. Search for a topic
2. Note hashtags used
3. Browse those tags for related content

### Search + Recent

1. Search for old content
2. Check Recent pages
3. Find what you worked on recently

## Use Cases by Role

### Students & Researchers

```
Search: "confirmation bias"
→ Find all research notes mentioning this concept
→ Review backlinks for related papers
```

### Developers

```
Search: "authentication middleware"
→ Locate API auth documentation
→ Find code examples
→ Update outdated patterns
```

### Writers & Content Creators

```
Search: "productivity"
→ Find all articles on productivity
→ Avoid repeating yourself
→ Link related pieces
```

### Personal Knowledge Management

```
Search: "book notes"
→ Find specific book summaries
→ Review highlights
→ Connect ideas across books
```

## Future Enhancements

Current XLog search is intentionally simple. Possible future additions:

- **Fuzzy matching** - Find similar words
- **Search operators** - AND, OR, NOT
- **Result ranking** - Relevance scoring
- **Search history** - Recent searches
- **Saved searches** - Bookmark queries

These would be added only if users request them. Current search handles 95% of use cases.

## Search Best Practices

### 1. Build Search Habits

Make search your first reflex:

```
"Did I write about X?" → Ctrl+K → Search
```

Faster than browsing directories.

### 2. Search Before Creating

Before creating a new page:

```
Search: "docker containers"
```

Check if you already have content on this topic.

### 3. Search When Linking

Not sure what you called a page?

```
Search: "network architecture"
→ Find exact page name
→ Use for [[Backlink]]
```

### 4. Regular Search Audits

Periodically search broad terms:

```
Search: "todo" or "draft" or "incomplete"
```

Find unfinished work.

## See Also

- [Backlinks](Backlinks.md) - Navigate by connections
- [Hashtags](Hashtags.md) - Browse by topic
- [official extensions](official extensions.md) - Search extension details
- [Usage](Usage.md) - Search keyboard shortcuts
