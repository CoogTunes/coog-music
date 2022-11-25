function toggleMenu() {
    let subMenu = document.getElementById("subMenu");
    let userControl = document.querySelector('.user-control');

    userControl.addEventListener('click', function (evt) {
        subMenu.classList.toggle("open-menu");
        evt.preventDefault();
    });
}

window.addEventListener('DOMContentLoaded', function () {
    toggleMenu();
});
