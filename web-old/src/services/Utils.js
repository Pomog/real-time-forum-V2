import Ws from "./Ws.js";

const parseJwt = (token) => {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
};

const getUser = () => {
    const user = {
        id: localStorage.getItem('sub'),
        role: localStorage.getItem('role'),
        accessToken: localStorage.getItem('accessToken'),
        refreshToken: localStorage.getItem('refreshToken')
    };

    // Log the user object to the console
    console.log("getUser output:", user);

    return user;
};

const logOut = () => {
    localStorage.removeItem('sub');
    localStorage.removeItem('role');
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    Ws.disconnect()
        .then(() => {
            console.log('WebSocket disconnected successfully.');
        })
        .catch((err) => {
            console.error('Error disconnecting WebSocket:', err);
        });
}



const fileToBase64 = (file) => {
    return new Promise(resolve => {
        let fileReader = new FileReader();

        fileReader.onload = (fileLoadedEvent) => {
            resolve(fileLoadedEvent.target.result)
        }

        fileReader.readAsDataURL(file)
    })
}

const base64isImage = (base64string) => {
    return /image\/(jpeg|png|gif)/.test(base64string)
}


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


    app.innerHTML = `
        <h1>${titles[status]}</h1><br>
        <h2>${message || ''}</h2>
    `
}

const drawErrorMessage = (err) => {
    const inputError = document.getElementById("error-message")
    if (inputError) {
        inputError.innerText = err
    }
}

const debounce = (func, wait, immediate) => {
    let timeout;
    return function () {
        const context = this, args = arguments;
        const later = function () {
            timeout = null;
            if (!immediate) func.apply(context, args);
        };
        const callNow = immediate && !timeout;
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
        if (callNow) func.apply(context, args);
    };
};

const throttle = (func, delay) => {
    let toThrottle = false;
    return function () {
        if (!toThrottle) {
            toThrottle = true;
            func.apply(this, arguments)
            setTimeout(() => {
                toThrottle = false
            }, delay);
        }
    };
};


export default { parseJwt, getUser, logOut, fileToBase64, base64isImage, showError, drawErrorMessage, debounce, throttle }