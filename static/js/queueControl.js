function queueControl() {
    let queueToggle = document.querySelector('.queue-show-trigger');
    let queueContainer = document.querySelector('.queue-container');

    queueToggle.addEventListener('click', function (evt) {
        queueContainer.classList.toggle('show');
    });

}

window.addEventListener('DOMContentLoaded', function (evt) {
    queueControl();
});