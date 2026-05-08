#extension

# Automatic Page Linking

The autolink pages extension is XLog's killer feature: it automatically converts page name mentions
into links without requiring markdown link syntax.

## How It Works

Simply write naturally and mention any page name. XLog automatically detects it and creates a link.

**Example:**

```markdown
I'm working on the API documentation and need to reference the Authentication page.
```

If you have pages named `API documentation.md` and `Authentication.md`, XLog automatically converts
those mentions to clickable links. No `[Page Name](/link)` syntax required.

## Benefits

1. **Natural Writing** - Write like you think; XLog handles the linking
2. **Retroactive Linking** - Create a new page and all existing mentions automatically become links
3. **Knowledge Graph** - Pages interconnect naturally as you write
4. **Zero Maintenance** - No need to manually add links across old pages

## Automatic Backlinks

The extension also maintains backlinks for every page. At the bottom of each page, you'll see a
list of all pages that link to the current page. This creates bidirectional relationships
automatically.

## Technical Details

- Case-insensitive matching: "Machine Learning" matches `machine learning.md`
- Longest match first: "machine learning model" matches before "machine learning"
- Word boundary detection: prevents "test.mdx" from matching "test.md"
