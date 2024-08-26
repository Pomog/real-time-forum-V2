import Ws from "./Ws.js";

// Function to parse JWT and extract the payload
const parseJwt = (token) => {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join('')); // Decode Base64 and URI components

    return JSON.parse(jsonPayload);
};

// Function to get the current user's details from localStorage
const getUser = () => {
    return {
        id: localStorage.getItem('sub'),
        role: localStorage.getItem('role'),
        accessToken: localStorage.getItem('accessToken'),
        refreshToken: localStorage.getItem('refreshToken')
    }
}

// Function to log the user out and clear user data from localStorage
const logOut = async () => {
    localStorage.removeItem('sub')
    localStorage.removeItem('role')
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    await Ws.disconnect()
}

// Function to convert a file to Base64 format
const fileToBase64 = (file) => {
    return new Promise(resolve => {
        let fileReader = new FileReader();

        fileReader.onload = (fileLoadedEvent) => {
            resolve(fileLoadedEvent.target.result)
        }

        fileReader.readAsDataURL(file)
    })
}

// Function to check if a Base64 string represents an image
const base64isImage = (base64string) => {
    return /image\/(jpeg|png|gif)/.test(base64string)
}

// Function to display an error message in the app
const showError = (status, message) => {
    const app = document.querySelector("#app")

    const titles = {
        400: "400 Bad Request",
        401: "401 Unauthorized",
        403: "403 Forbidden",
        404: "404 Not Found",
        500: "500 Internal Server Error",
        503: "503 Service Unavailable"
    }

    // Set the inner HTML of the app element to display the error
    app.innerHTML = `
        <h1>${titles[status]}</h1><br>
        <h2>${message || ''}</h2>
    `
}

// Function to display an error message in a specific element with ID "error-message"
const drawErrorMessage = (err) => {
    const inputError = document.getElementById("error-message")
    if (inputError) {
        inputError.innerText = err
    }
}

// Function to debounce another function, delaying its execution
const debounce = (func, wait, immediate) => {
    let timeout;
    return function () {
        const context = this, args = arguments;
        const later = function () {
            timeout = null;
            if (!immediate) func.apply(context, args); // Call the function after delay
        };
        const callNow = immediate && !timeout;
        clearTimeout(timeout);
        timeout = setTimeout(later, wait); // Set timeout for the delay
        if (callNow) func.apply(context, args);
    };
};


// Function to throttle another function, limiting its execution rate
const throttle = (func, delay) => {
    let toThrottle = false;
    return function () {
        if (!toThrottle) {
            toThrottle = true;
            func.apply(this, arguments) // Execute the function
            setTimeout(() => {
                toThrottle = false // Reset throttle after delay
            }, delay);
        }
    };
};


export default { parseJwt, getUser, logOut, fileToBase64, base64isImage, showError, drawErrorMessage, debounce, throttle }