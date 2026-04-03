# File Operations

The `file-operations` extension provides the ability to rename and delete pages directly from the web interface.

## Deleting a Page

To delete the current page, you can use the command palette (usually accessed by pressing `Ctrl+Shift+P` or `Cmd+Shift+P`) and search for the "Delete page" command.

Alternatively, you can send a `DELETE` request to the `/-/file/delete` endpoint with the page name as a parameter.

## Renaming a Page

To rename the current page, use the command palette and search for the "Rename page" command. This will present a form where you can enter the new name for the page.

Behind the scenes, this functionality is handled by the `/-/file/rename` endpoint. A `GET` request to this endpoint shows the form, and a `POST` request with the new name performs the rename operation.
