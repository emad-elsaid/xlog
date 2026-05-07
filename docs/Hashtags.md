# Hashtags

Hashtags provide flexible, non-hierarchical organization for your knowledge base. Tag pages with topics, themes, or categories, and XLog automatically creates browsable tag collections.

## What Are Hashtags?

Hashtags in XLog work like social media tags - prefix any word with `#` to create a tag:

```markdown
This is about #machine-learning and #neural-networks.

I'm exploring #golang and #web-development today.
```

XLog automatically:
1. **Makes hashtags clickable** - Click to see all pages with that tag
2. **Creates tag pages** - Aggregate pages by topic
3. **Shows tag counts** - See how many pages use each tag

## Syntax

### Basic Hashtags

Simple one-word tags:

```markdown
#python #golang #javascript
#learning #research #tutorial
```

### Multi-Word Hashtags

Use hyphens for multi-word tags:

```markdown
#machine-learning
#digital-gardening
#personal-knowledge-management
#web-development
```

Spaces don't work in hashtags - use hyphens instead:

```markdown
#knowledge-base  ✅ Works
#knowledge base  ❌ Only tags "knowledge"
```

### Case Sensitivity

Hashtags are case-insensitive for matching:

```markdown
#MachineLearning
#machine-learning
#MACHINE-LEARNING
```

All three are treated as the same tag.

## Using Hashtags

### In Page Content

Add hashtags anywhere in your markdown:

```markdown
# My Research Notes

Today I learned about #transformers in #natural-language-processing.

The #attention-mechanism is key to understanding #bert and #gpt.

Related: #machine-learning #deep-learning
```

### In Frontmatter

You can also define tags in YAML frontmatter:

```markdown
---
tags:
  - machine-learning
  - neural-networks  
  - research
---

# Neural Network Architectures

Content here...
```

Both inline hashtags and frontmatter tags work together.

### Multiple Tags

Tag with as many topics as relevant:

```markdown
# Building a Personal Wiki with XLog

Notes on #xlog #digital-gardening #knowledge-management
#markdown #static-site-generator #personal-wiki

Tags: #golang #web-development #open-source
```

## Browsing by Tag

### Click Any Hashtag

Clicking a hashtag (like `#machine-learning`) takes you to a tag page showing:

- **All pages** containing that tag
- **Page titles** and excerpts
- **Tag count** (number of tagged pages)

Example tag page:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
#machine-learning (12 pages)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Neural Network Architectures
  Deep learning models for...

Introduction to Transformers
  Attention mechanisms and...

Supervised Learning Basics
  Classification and regression...
```

### Tag Index

XLog can show all tags used across your knowledge base (if configured).

## Hashtags vs. Backlinks

Both organize knowledge but serve different purposes:

| Aspect | Hashtags | Backlinks |
|--------|----------|-----------|
| **Purpose** | Categorization by topic | Connection between specific ideas |
| **Granularity** | Broad themes | Specific concepts |
| **Creation** | Manual with # | Automatic when mentioning pages |
| **Structure** | Flat categories | Network graph |
| **Best for** | Finding similar content | Discovering relationships |

**Example:**

```markdown
# Attention Mechanisms in Transformers

The [[Transformer Architecture]] uses #attention-mechanisms.

Introduced in [[Attention Is All You Need]].

Topics: #machine-learning #nlp #deep-learning
```

- **Hashtags** (#machine-learning, #nlp) categorize broadly
- **Backlinks** ([[Transformer Architecture]]) connect specific concepts

## Use Cases

### 1. Topic Organization

Organize notes by subject:

```markdown
# Linear Regression Notes

Basic #statistics and #machine-learning algorithm.

Used for #regression tasks, not #classification.

Learn: #supervised-learning #algorithms
```

Find all #statistics pages, all #machine-learning pages, etc.

### 2. Project Tracking

Track work across projects:

```markdown
# API Endpoint Refactoring

Refactor auth endpoints for #project-apollo #backend.

Related work: #api-design #security

Status: #in-progress #Q2-2026
```

See all #project-apollo tasks or all #in-progress items.

### 3. Learning Themes

Track learning areas:

```markdown
# React Hooks Deep Dive

Learning about #react #javascript #frontend.

Part of #learning-path-frontend #2026-goals

Difficulty: #intermediate
```

Browse #learning-path-frontend to see your curriculum.

### 4. Content Type

Mark content types:

```markdown
# How to Build a RESTful API

Complete guide to REST APIs.

Type: #tutorial #guide #how-to
Topics: #api-design #backend #http
```

Find all #tutorial pages or all #how-to guides.

### 5. Status and Workflow

Track page maturity:

```markdown
# Neural Network Theory (Draft)

Status: #draft #needs-review #work-in-progress

Once complete: #published #ready
```

See all #draft pages that need finishing.

## Organizing Strategies

### Multi-Dimensional Tagging

Tag from multiple angles:

```markdown
#topic-machine-learning     ← What it's about
#type-research-note         ← What kind of page
#status-published           ← Completion status
#project-thesis             ← Related project
#difficulty-advanced        ← Complexity level
```

This creates flexible, multi-faceted organization.

### Tag Prefixes

Use prefixes for clarity:

```markdown
#lang-python
#lang-javascript
#lang-go

#topic-algorithms
#topic-databases
#topic-networking

#status-draft
#status-published
#status-archived
```

Prefixes group related tags and make browsing easier.

### Hierarchical Tags (Optional)

While hashtags are flat, you can simulate hierarchy with naming:

```markdown
#programming
#programming-languages
#programming-languages-python
#programming-languages-python-async
```

Not enforced by XLog, but useful for personal organization.

## Best Practices

### 1. Be Consistent

Pick tag names and stick with them:

```markdown
#machine-learning    ← Consistent
#ml                  ← Don't mix abbreviations
#MachineLearning     ← Different case okay, but pick one style
```

### 2. Not Too Many, Not Too Few

**Too few tags:**
```markdown
#notes  ← Too broad, not useful
```

**Too many tags:**
```markdown
#machine #learning #neural #networks #deep #learning #ai #ml #dl
← Too granular, redundant
```

**Just right:**
```markdown
#machine-learning #neural-networks #deep-learning
← Specific enough to be useful
```

### 3. Use Tags for Discovery

Tag content you'll want to find again:

```markdown
# Solution to CORS Issues

Quick fix for #troubleshooting #cors #web-development.

Bookmark: #reference #solutions
```

Later, browse #troubleshooting for debugging help.

### 4. Combine with Backlinks

Use both for powerful organization:

```markdown
# React Context API

Discusses [[React]] state management via #context-api.

Related patterns: [[Redux]], [[MobX]]

Topics: #react #state-management #frontend
```

- Backlinks connect to specific libraries
- Hashtags categorize by theme

### 5. Review and Refactor Tags

Periodically review your tags:
- Merge similar tags (#ml → #machine-learning)
- Rename for clarity (#proj1 → #project-apollo)
- Delete unused tags

## Tag Performance

Hashtags are lightweight:
- Indexed during page rendering
- Fast lookup even with thousands of tags
- No performance impact on large knowledge bases

## Configuration

Hashtags are provided by the hashtags extension (enabled by default).

To disable hashtags:
```bash
xlog -disabled-extensions "hashtags"
```

(Not recommended - hashtags are essential for organization)

## Examples

### Personal Wiki

```markdown
# Database Indexing Strategies

Exploring #database #performance #optimization techniques.

Types: #b-tree #hash-index #bitmap-index

Related: [[SQL Optimization]] [[Query Planning]]

Topics: #postgresql #mysql #database-design
```

### Research Notes

```markdown
# Paper: Attention Is All You Need

Seminal #research-paper on #transformers.

Authors: [[Vaswani et al.]]

Topics: #nlp #deep-learning #attention-mechanisms
Year: #2017 #papers-2017
Impact: #highly-cited #foundational
```

### Learning Journal

```markdown
# 2026-05-07 - Learning Log

Today: #golang #concurrency #goroutines

Completed: #tutorial #effective-go

Next: #channels #context #patterns

Status: #learning-in-public #100-days-of-code
```

### Project Documentation

```markdown
# Authentication Service Architecture

Component: #backend #microservices #authentication

Stack: #golang #jwt #oauth2

Project: #main-platform #api-gateway

Status: #production #v2.0
```

## Tag Naming Conventions

### Technical Content

```markdown
#lang-python        ← Programming languages
#tool-docker        ← Tools and software
#concept-oop        ← Concepts and theory
#pattern-singleton  ← Design patterns
```

### Personal Knowledge

```markdown
#book-atomic-habits     ← Books
#course-cs50            ← Courses  
#video-mit-ocw          ← Video series
#podcast-lex-fridman    ← Podcasts
```

### Status and Workflow

```markdown
#todo              ← Action items
#doing             ← In progress
#done              ← Completed
#blocked           ← Waiting on something
#someday-maybe     ← Backlog
```

## Limitations

### No Tag Hierarchy

XLog treats all tags as flat. These aren't parent-child:

```markdown
#programming
#programming-languages
#programming-languages-python
```

You can name them hierarchically, but XLog doesn't enforce relationships.

### No Tag Aliases

Each tag is independent:

```markdown
#ml                  ← Different tag
#machine-learning    ← Different tag
```

Pick one and be consistent.

### No Tag Descriptions

Tags are just labels - no descriptions or metadata. Document tag meanings in an index page:

```markdown
# Tag Index

## Project Tags
- **#project-apollo** - Main web platform
- **#project-beta** - Mobile app project

## Status Tags
- **#draft** - Work in progress
- **#published** - Ready for sharing
```

## Advanced Usage

### Tag-Based Navigation

Create tag navigation pages:

```markdown
# Topic Index

## Programming
- [#golang](/tags/golang)
- [#python](/tags/python)  
- [#javascript](/tags/javascript)

## Concepts
- [#algorithms](/tags/algorithms)
- [#data-structures](/tags/data-structures)
```

### Tag Clouds (Manual)

List your most-used tags:

```markdown
# My Knowledge Base

Most active topics:
#machine-learning (45) #golang (32) #web-development (28)
#algorithms (24) #database-design (18) #react (15)
```

### Dynamic Tag Pages

Tag pages are generated automatically. Bookmark frequently-used ones:

```markdown
# Quick Links

- [Machine Learning Notes](/tags/machine-learning)
- [Ongoing Projects](/tags/in-progress)
- [Reference Material](/tags/reference)
```

## See Also

- [Backlinks](Backlinks.md) - Page-to-page connections
- [Search](Search.md) - Find content across tags
- [official extensions](official extensions.md) - Hashtags extension details
- [Digital Gardening](digital gardening.md) - Organizational philosophy
