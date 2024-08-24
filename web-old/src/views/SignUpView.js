import AbstractView from "./AbstractView.js";
import fetcher from "../services/Fetcher.js";
import router from "../index.js";


const signUp = async (body) => {
    console.log("in the signUp()");
    console.log("body");
    console.log(body);
    const path = `/api/users/sign-up`;

    try {
        const data = await fetcher.post(path, body);
        console.log("response from /api/users/sign-up");
        console.log(data);

        if (data && data.error) {
            drawError(data.error);
            return;
        }
        router.navigateTo("/sign-in");
    } catch (e) {
        console.error("Error during sign-up", e);
        drawError("An unexpected error occurred");
    }
};


const drawError = (err) => {
    const errorMessage = document.getElementById("error-message");
    errorMessage.innerText = err || "An unknown error occurred";
};

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Sign up");
    }

    async getHtml() {
        console.log("Returning View for sign-up-form")
        return `
            <form id="sign-up-form" onsubmit="return false;">
                Username: <br>
 <input type="text" id="username" placeholder="Username" required minlength="2" maxlength="64" 
    pattern="^[a-zA-Z0-9](?!.*[_.-]{2})[a-zA-Z0-9._-]{0,62}[a-zA-Z0-9]$"
    title="Username should only contain alphanumeric characters, '.', '_', and '-' symbols. No symbol at the beginning or end, and no consecutive special characters.">
    <br>
<!--    
    <input type="text" id="username" placeholder="Username" required minlength="2" maxlength="64"
    pattern="^(?![_.])(?!.*[_.-]{2})[a-zA-Z0-9._-]+(?<![_.-])$"
    title="Username should only contain alphanumerical and '.', '_', '-' symbols,
    no symbol at the beginning and at the end, no alternation of special characters">
    <br> <br>
/*-->

                First name: <br>
                <input type="text" id="first-name" placeholder="First name" required minlength="2" maxlength="64" 
                    pattern="[a-zA-Z]+$" title="First name should only contain Latin letters.">  <br> <br>

                Last name: <br>
                <input type="text" id="last-name" placeholder="Last name" required minlength="2" maxlength="64"  
                    pattern="[a-zA-Z]+$" title="Last name should only contain Latin letters."> <br> <br>

                Age: <br>
                <input type="number" id="age" placeholder="Age" required min="12" max="110"> <br> <br>

                Gender: <br>
                <input type="radio" name="gender" id="gender-male" value="1" required> Male
                <input type="radio" name="gender" id="gender-female" value="2"> Female <br> <br>

                E-mail: <br>
                <input type="email" id="email" placeholder="E-mail" required maxlength="64">  <br><br>
                
                Password: <br>
 <input type="password" id="password" placeholder="Password" minlength="7" maxlength="64" required 
    pattern="^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#~$%^&*()+|_]).{7,}$" 
    title="Password must contain at least one lowercase letter, one uppercase letter, one number, and one special character. Allowed symbols: !@#~$%^&*()+|_">
    <br>

                
                Confirm password: <br>
                <input type="password" id="password-confirm" placeholder="Password" maxlength="64" required>

                <div class="error" id="error-message"></div>
                
                <button type="submit">Sign up</button>
            </form>
        `;
    }

    /*
sign-up:1  Pattern attribute value ^[a-zA-Z0-9](?!.*[_.-]{2})[a-zA-Z0-9._-]{0,62}[a-zA-Z0-9]$ is not a valid regular expression: Uncaught SyntaxError: Invalid regular expression: /^[a-zA-Z0-9](?!.*[_.-]{2})[a-zA-Z0-9._-]{0,62}[a-zA-Z0-9]$/v: Invalid character in character class
sign-up:1  Pattern attribute value ^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#~$%^&*()+|_]).{7,}$ is not a valid regular expression: Uncaught SyntaxError: Invalid regular expression: /^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#~$%^&*()+|_]).{7,}$/v: Invalid character in character class
     */

    async init() {
        const signUpForm = document.getElementById("sign-up-form");

        signUpForm.addEventListener("submit", function () {
            const password = document.getElementById("password");
            const passwordConfirm = document.getElementById("password-confirm");

            console.log(signUpForm);
            if (password.value !== passwordConfirm.value) {
                //TODO: here IS a Problem
                drawError("Passwords Don't Match");
            } else {
                const input = {
                    username: document.getElementById("username").value,
                    firstName: document.getElementById("first-name").value,
                    lastName: document.getElementById("last-name").value,
                    age: parseInt(document.getElementById("age").value),
                    gender: parseInt(document.querySelector('input[name="gender"]:checked').value),
                    email: document.getElementById("email").value,
                    password: password.value,
                };

                console.log("INPUT FOR signUp");
                console.log(input);
                signUp(input);
            }
        });
    }
}