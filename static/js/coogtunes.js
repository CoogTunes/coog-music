let subMenu = document.getElementById("subMenu");
function toggleMenu() {
    subMenu.classList.toggle("open-menu")
}

let x = document.getElementById("login");
let y = document.getElementById("register");
let z = document.getElementById("btn");
let password = document.getElementById("password");
let confirm_password = document.getElementById("confirm_password");
let form_box = document.getElementById("form_box")
let currentTab = document.querySelector('.toggle_btn.active');

function register(evt) {
    // x.style.left = "-400px";
    // y.style.left = "50px";
    y.style.transform = "translate(" + 0 + "%)";
    x.style.transform = "translate(" + 200 + "%)";
    z.style.left = "98px";
    form_box.style.height = "640px";
    if(currentTab !== evt.target && currentTab != null)
        currentTab.classList.remove('active');
    evt.target.classList.add('active');
    currentTab = evt.target;
}
let togglePassword = document.getElementById("togglePassword");
let toggleCPassword = document.getElementById("toggleCPassword");

togglePassword.addEventListener("click", function () {
    const type = password.getAttribute("type") === "password" ? "text" : "password";
    password.setAttribute("type", type);
    this.classList.toggle("bi-eye");
})
toggleCPassword.addEventListener("click", function () {
    const type = confirm_password.getAttribute("type") === "password" ? "text" : "password";
    confirm_password.setAttribute("type", type);
    this.classList.toggle("bi-eye");
})

function login(evt) {
    // x.style.left = "50px";
    // y.style.left = "450px";
    y.style.transform = "translate(" + -200 + "%)";
    x.style.transform = "translate(" + 0 + "%)";
    z.style.left = "0"
    form_box.style.height = "400px";
    if(currentTab !== evt.target && currentTab != null)
        currentTab.classList.remove('active');
    evt.target.classList.add('active');
    currentTab = evt.target;
}
let textOverflowManager = function(selector, maxLength){
    let truncateTxt = document.querySelectorAll(selector)

    truncateTxt.forEach(text => {
        text.innerHTML = truncateText(text.innerHTML, maxLength);
    });
}

function truncateText(text, maxLength){
    let returnTxt = text;
    if(returnTxt.length > maxLength){
        return returnTxt.substring(0, maxLength) + "...";
    }
    return returnTxt;
}

function validatePassword() {
    const isNonWhiteSpace = /^\S*$/;
    if (!isNonWhiteSpace.test(password.value)) {
        return "Password must not contain whitespaces."
    }

    const containsUpper = /^(?=.*[A-Z]).*$/;
    if (!containsUpper.test(password.value)) {
        return "Password must contain one uppercase."
    }

    const containsLower = /^(?=.*[a-z]).*$/;
    if (!containsLower.test(password.value)) {
        return "Password must contain one lowercase."
    }

    const isContainsNumber = /^(?=.*[0-9]).*$/;
    if (!isContainsNumber.test(password.value)) {
        return "Password must contain one digit."
    }

    const containsSpecialCharacter = /^(?=.*[~`!@#$%^&*()--+={}\[\]|\\:;"'<>,.?/_₹]).*$/;
    if (!containsSpecialCharacter.test(password.value)) {
        return "Password must contain one special character"
    }

    const isValidLength = /^.{8,15}$/
    if (!isValidLength.test(password.value)) {
        return "Password must be 8-15 characters long"
    }

    if (password.value !== confirm_password.value) {
        return "Passwords don't match"
    }

}

function checkPassword() {
    let message = validatePassword()
    if (!message) {
        password.setCustomValidity("")
    }
    else {
        password.setCustomValidity(message)
    }
}
password.onchange = checkPassword
confirm_password.onkeyup = checkPassword



window.addEventListener('DOMContentLoaded', function () {
    textOverflowManager('.truncate-txt', 75);
});