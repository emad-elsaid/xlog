import {EditorView} from "@codemirror/view"
import {EditorState} from "@codemirror/state"
import {markdownLanguage} from "@codemirror/lang-markdown"
import {lintKeymap} from "@codemirror/lint"
import {keymap, highlightSpecialChars, drawSelection, dropCursor, rectangularSelection, crosshairCursor, highlightActiveLineGutter} from "@codemirror/view"
import {defaultHighlightStyle, syntaxHighlighting, indentOnInput, bracketMatching, foldGutter, foldKeymap} from "@codemirror/language"
import {defaultKeymap, history, historyKeymap, indentWithTab} from "@codemirror/commands"
import {searchKeymap, highlightSelectionMatches} from "@codemirror/search"
import {autocompletion, completionKeymap, closeBrackets, closeBracketsKeymap} from "@codemirror/autocomplete"

const textarea = document.getElementById("content")
textarea.setAttribute("hidden", "true");

const editor = document.getElementById("editor")

const startState = EditorState.create({
  doc: textarea.value,
  extensions: [
    markdownLanguage,
    highlightActiveLineGutter(),
    highlightSpecialChars(),
    history(),
    drawSelection(),
    dropCursor(),
    EditorState.allowMultipleSelections.of(true),
    EditorView.lineWrapping,
    EditorView.perLineTextDirection.of(true),
    indentOnInput(),
    syntaxHighlighting(defaultHighlightStyle, {fallback: true}),
    bracketMatching(),
    closeBrackets(),
    autocompletion(),
    rectangularSelection(),
    crosshairCursor(),
    highlightSelectionMatches(),
    keymap.of([
      ...closeBracketsKeymap,
      ...defaultKeymap,
      ...searchKeymap,
      ...historyKeymap,
      ...completionKeymap,
      ...lintKeymap,
      indentWithTab
    ])  ],
})

const view = new EditorView({
  state: startState,
  parent: editor,
})

view.focus()
