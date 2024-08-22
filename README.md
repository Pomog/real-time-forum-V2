# Real-Time Forum V2

Real-Time Forum V2 is a single-page application built with Golang, SQLite, JavaScript, HTML, and CSS. It allows users to register, log in, create posts, comment on posts, and send private messages in real-time. The forum leverages WebSockets for real-time communication and provides a modern, interactive user experience without relying on frontend frameworks like React or Angular.

- [SERVER instructions](server.md)
- [CLIENT instructions](client.md)

- [View Management](viewmanagement.md)

## Features

### Registration and Login
- Users can register with the following details: Nickname, Age, Gender, First Name, Last Name, E-mail, and Password.
- Login is possible using either the nickname or e-mail combined with the password.
- Users can log out from any page on the forum.

### Posts and Comments
- Users can create posts categorized into different sections.
- Posts are displayed in a feed view.
- Users can comment on posts, with comments displayed only when a post is selected.

### Private Messaging
- Real-time private messaging with WebSocket support.
- A user-friendly chat interface that displays online/offline status.
- Messages are displayed with timestamps and sender information.
- Chat history is loaded dynamically as users scroll.

## Project Structure

- **Backend**: Golang handles server-side logic, WebSocket connections, and data management.
- **Database**: SQLite is used to store user data, posts, comments, and messages.
- **Frontend**: JavaScript manages the client-side WebSocket connections and DOM manipulation. HTML and CSS provide the structure and styling for the application.
- **Single Page Application**: The forum is a single-page application, with all page changes handled dynamically via JavaScript.

## Technologies Used

- **Golang**: Backend development, WebSocket handling, and data processing.
- **SQLite**: Database management.
- **JavaScript**: Frontend interactivity, WebSocket client implementation.
- **HTML/CSS**: Page structure and styling.

## Getting Started

### Prerequisites
- Go 1.15+
- SQLite
- Git
