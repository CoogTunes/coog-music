function modalManager(){
    let currentModal = null;
    document.addEventListener('click', function (evt) {
        let target = evt.target;
        if(target.parentElement.hasAttribute('data-target-modal')){
            showModal(target.parentElement.getAttribute('data-target-modal'));
        }
        else if(target.classList.contains('modal-container') || target.classList.contains('close-modal')){
            hideModal();
        }
    });

    function showModal(targetModal){
        let modal = document.getElementById(targetModal);
        if(modal){
            modal.classList.add('show');
            currentModal = modal;
        }
    }

    function hideModal(){
        currentModal.classList.remove('show');
    }
}

window.addEventListener('DOMContentLoaded', function () {
    modalManager();
});