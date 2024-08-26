# Utils Utility Documentation
The Utils.js module provides a set of utility functions that are commonly used throughout the application.
These functions handle tasks such as JWT parsing, user management, file conversions, error handling, and function optimizations like debouncing and throttling.

## Features
JWT Parsing: Decodes JSON Web Tokens (JWT) to extract the payload.
User Management: Retrieves user details from localStorage and handles user logout.
File Handling: Converts files to Base64 format and checks if a Base64 string represents an image.
Error Handling: Displays error messages based on HTTP status codes and custom messages.
Function Optimizations: Includes debouncing and throttling functions to control the frequency of function execution.
## Module Structure
  - **parseJwt(token)**
Description: Parses a JSON Web Token (JWT) to decode its payload and returns it as a JSON object.
Parameters: token (string) - The JWT string.
Returns: A JSON object representing the decoded payload.
```javascript
const tokenPayload = Utils.parseJwt(accessToken);
console.log(tokenPayload);
```


  - **getUser()**
Description: Retrieves the current user's details (ID, role, access token, refresh token) from localStorage.
Returns: An object containing the user's id, role, accessToken, and refreshToken.
```javascript
const user = Utils.getUser();
console.log(user.id, user.role);
```

  - **logOut()**
Description: Logs the user out by removing user-related information from localStorage and disconnecting the WebSocket.
```javascript
Utils.logOut();
```

  - **fileToBase64(file)**
Description: Converts a file object to a Base64-encoded string.
Parameters: file (File) - The file to be converted.
```javascript
const base64String = await Utils.fileToBase64(file);
console.log(base64String);
```

  - **base64isImage(base64string)**
Description: Checks if a Base64 string represents an image (JPEG, PNG, GIF).
Parameters: base64string (string) - The Base64 string to be checked.
Returns: A boolean (true if the string is an image, false otherwise).
```javascript
const isImage = Utils.base64isImage(base64String);
console.log(isImage); // true or false
```

  - **showError(status, message)**
Description: Displays an error message on the application based on the HTTP status code and an optional custom message.
Parameters:
status (number) - The HTTP status code (e.g., 404, 500).
message (string, optional) - A custom error message.
```javascript
Utils.showError(404, "Page not found");
```

  - **drawErrorMessage(err)**
Description: Displays an error message in a specific DOM element with the ID error-message.
Parameters: err (string) - The error message to display.
```javascript
Utils.drawErrorMessage("Invalid username or password.");
```

  - **debounce(func, wait, immediate)**
Description: Creates a debounced function that delays invoking func until after wait milliseconds have elapsed since the last time the debounced function was invoked.
Parameters:
func (function) - The function to debounce.
wait (number) - The number of milliseconds to delay.
immediate (boolean, optional) - If true, trigger the function on the leading edge, instead of the trailing.
Returns: A debounced version of the func.
```javascript
const debouncedFunction = Utils.debounce(() => console.log('Debounced!'), 300);
window.addEventListener('resize', debouncedFunction);
```

  - **throttle(func, delay)**
Description: Creates a throttled function that only invokes func at most once per every delay milliseconds.
Parameters:
func (function) - The function to throttle.
delay (number) - The number of milliseconds to wait before allowing the next call.
Returns: A throttled version of the func.
```javascript
const throttledFunction = Utils.throttle(() => console.log('Throttled!'), 1000);
window.addEventListener('scroll', throttledFunction);
```

## Usage Example
```javascript
import Utils from './path/to/Utils.js';

// Parse JWT
const tokenPayload = Utils.parseJwt(accessToken);

// Get user details
const user = Utils.getUser();

// Log out user
Utils.logOut();

// Convert a file to Base64
const base64String = await Utils.fileToBase64(file);

// Check if a Base64 string is an image
const isImage = Utils.base64isImage(base64String);

// Show error message
Utils.showError(404, "Page not found");

// Debounce a function
const debouncedResize = Utils.debounce(() => console.log('Resized!'), 300);
window.addEventListener('resize', debouncedResize);

// Throttle a function
const throttledScroll = Utils.throttle(() => console.log('Scrolled!'), 1000);
window.addEventListener('scroll', throttledScroll);
```