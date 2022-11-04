let subMenu = document.getElementById("subMenu");
function toggleMenu() {
    subMenu.classList.toggle("open-menu")
}

const music = new Audio('/MusicLibrary/audio/adele.mp3');
let masterPlay = document.getElementById("masterPLay");
if (masterPlay) {
    masterPlay.addEventListener('click', () => {
        if (music.paused || music.currentTime <= 0) {
            music.play();
            masterPlay.classList.remove('bi-play-fill')
            masterPlay.classList.add('bi-pause-fill')
        } else {
            music.pause();
            masterPlay.classList.add('bi-play-fill')
            masterPlay.classList.remove('bi-pause-fill')
        }
    });
}
//login
let x = document.getElementById("login");
let y = document.getElementById("register");
let z = document.getElementById("btn");
let password = document.getElementById("password");
let confirm_password = document.getElementById("confirm_password");
let form_box = document.getElementById("form_box")

function register() {
    x.style.left = "-400px";
    y.style.left = "50px";
    z.style.left = "110px";
    form_box.style.height = "460px"
}
function login() {
    x.style.left = "50px";
    y.style.left = "450px";
    z.style.left = "0"
    form_box.style.height = "400px"
}

function validatePassword() {
    if (password.value !== confirm_password.value) {
        return confirm_password.setCustomValidity("Passwords don't match")
    }
    return confirm_password.setCustomValidity("")
}
password.onchange = validatePassword
confirm_password.onkeyup = validatePassword
