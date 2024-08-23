# HomeView Documentation
The HomeView class is responsible for rendering and managing the homepage of the application, specifically displaying post categories and posts within those categories. It extends the AbstractView class, handling the dynamic loading of posts based on selected categories and pagination.

## Features
- **Category Display:**
Dynamically loads and displays post categories from the server.
Users can select a category to view posts within that category.

- **Post Loading:**
Fetches and displays posts based on the selected category and page number.
Handles pagination with "Newer" and "Older" buttons to navigate through pages of posts.
Automatically updates the URL query parameters to reflect the current category and page.

- **Pagination:**
Users can navigate through pages of posts using "Newer" and "Older" buttons.
The pagination state is updated and reflected in the URL.

## Class Structure
### Constructor
```javascript
constructor(params) {
    super(params);
    this.setTitle("Home");
}
```
The constructor sets the title of the page to "Home" and initializes the view.

## Methods
- **getHtml()**
This method returns the HTML structure for the home page, including containers for categories, posts, and pagination controls.

```javascript
async getHtml() {
    return `
        <div id="categories"></div>
        <div id="category-title"></div>
        <div id="posts"></div>
        <div id="posts-msg"></div>
        <div class="navigation-buttons">
            <button id="prev-button">Newer</button>
            <p id="page-number" style="display: none"></p>
            <button id="next-button">Older</button>
        </div>
    `;
}
```

- **init()**
The init() method is responsible for initializing the view, including loading categories, setting up event listeners for category selection, and handling pagination.

```javascript
async init() {
    const urlParams = new URLSearchParams(window.location.search)
    currCategoryID = urlParams.get('category') || 1
    currPageNum = urlParams.get('page') || 1
    updateQueryParams()

    const categories = await getCategories()
    if (!categories) {
        return
    }
    await drawCategories(categories)

    const categoryEl = document.getElementById(`category-${currCategoryID}`)
    if (!categoryEl) {
        Utils.showError(404, `Cannot find category`)
        return
    } else {
        categoryEl.click()
    }

    const nextButtonEl = document.getElementById(`next-button`)
    const prevButtonEl = document.getElementById(`prev-button`)
    const pageNumber = document.getElementById(`page-number`)
    pageNumber.innerText = currPageNum

    nextButtonEl.addEventListener("click", () => {
        if (postsEnded) {
            return
        }
        currPageNum++
        pageNumber.innerText = currPageNum
        updateQueryParams()

        drawPostsByCategoryID(currCategoryID, currPageNum)
    })

    prevButtonEl.addEventListener("click", () => {
        if (currPageNum !== 1) {
            postsEnded = false
            currPageNum--
            pageNumber.innerText = currPageNum
            updateQueryParams()
            drawPostsByCategoryID(currCategoryID, currPageNum)
        }
    })
}
```
## Utility Functions
- **getCategories()**
Fetches the list of categories from the server.

- **drawCategories(categories)**
Dynamically creates and appends category buttons to the DOM, setting up click handlers to load posts by category.

- **drawPostsByCategoryID(categoryID, page)**
Fetches and displays posts for a specific category and page.

- **updateQueryParams()**
Updates the URL query parameters to reflect the current category and page number.

## External Dependencies
- [fetcher](fetcher.md): Handles HTTP requests to fetch data from the server.
- [Utils](Utils.md): Provides utility functions like error handling.
- [PostUtils](PostUtils.md): Provides helper functions for creating and managing post elements.
