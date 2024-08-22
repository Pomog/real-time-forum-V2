# NavBarView Documentation

The `NavBarView` class is responsible for rendering and managing the navigation bar within the application.
It extends the `AbstractView` class and dynamically adjusts the navigation options based on the user's authentication status.

## Features

- **Dynamic Navigation Links:**
    - If the user is logged in, the navigation bar shows links for "Home", "New Post", "Chats", "Profile", and "Logout".
    - If the user is not logged in, the navigation bar shows links for "Sign Up" and "Sign In".

- **Logout Handling:**
    - The logout button, when clicked, will trigger the `Utils.logOut()` function to clear user session data, disconnect WebSocket connections, and clear any active intervals.

## Class Structure

### Constructor

```javascript
constructor(params, user) {
    super(params);
    this.user = user;
}
```
