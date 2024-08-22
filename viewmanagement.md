# Client-Side Routing and View Management

This module is responsible for handling client-side routing, view rendering, and managing user roles in the application.
It dynamically loads and initializes views based on the current URL path and the user's role, providing a seamless and role-based navigation experience.

## Key Components

### 1. Module Imports
- **Views:**
    - [`NavBar`](navbar.md): Manages the navigation bar UI.
    - [`Home`](homeview.md): Renders the homepage view.
    - `SignUp`: Renders the sign-up page.
    - `SignIn`: Renders the sign-in page.
    - `Post`: Displays a single post view.
    - `NewPost`: Renders the new post creation page.
    - `Chats`: Displays the chat interface.
    - `Profile`: Displays the user profile view.

- **Services:**
    - `Ws`: Manages WebSocket connections.
    - `Utils`: Provides utility functions, including user management and error handling.

### 2. User Roles
The application defines four user roles:
- `guest`: Unauthenticated users.
- `user`: Regular users.
- `moderator`: Users with moderation privileges.
- `admins`: Administrators with full access.

### 3. Path Matching
- **`pathToRegex(path)`:** Converts a route path (e.g., `/post/:postID`) into a regular expression to match against the current URL.
- **`getParams(match)`:** Extracts URL parameters from a matched route and returns them as an object.

### 4. Navigation
- **`navigateTo(url)`:** Pushes a new entry to the browser's history stack and triggers the `router` function to load the corresponding view.

### 5. Router Function
The `router` function is the core of the client-side routing system:
- **Routes Definition:** Defines all possible routes and the minimum user role required to access each route.
- **Matching Routes:** Iterates through the defined routes to find a match with the current URL.
- **Role Validation:** Checks the user's role against the required role for the matched route. If the user's role is insufficient, an error is displayed.
- **View Rendering:** Renders the matched view, passing any URL parameters and the current user information.

### 6. Event Listeners
- **`popstate`:** Listens for back/forward navigation events (e.g., using the browser's back button) and re-runs the `router` function.
- **`storage`:** Monitors changes to `localStorage` to detect if the user logs out in a different tab. If the user logs out, the page is reloaded to update the UI.
- **`DOMContentLoaded`:** Initializes the application when the DOM is fully loaded:
    - Adds a click event listener to handle navigation through links with the `data-link` attribute.
    - If an access token is present in `localStorage`, it establishes a WebSocket connection.
    - Calls the `router` function to load the initial view.

### 8. Exported Functionality
- **`navigateTo(url)`:** Exposed for external use, allowing other modules to trigger navigation programmatically.