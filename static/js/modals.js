function modalManager(){
    // * Current Triggered Modal
    let currentModal = null;

    // * Open and close the modal
    document.addEventListener('click', function (evt) {
        let target = evt.target;
        if(target.parentElement.hasAttribute('data-target-modal')){
            showModal(target.parentElement.getAttribute('data-target-modal'));
        }
        else if(target.classList.contains('modal-container') || target.classList.contains('close-modal')){
            hideModal();
        }
    });

    // * Escape key press hideModal()
    document.addEventListener('keydown', function (evt) {
        if((evt.key === 'Escape' || evt.key === 'Esc') && evt.keyCode === 27 && currentModal !== null){
            hideModal();
        }
    });

    // * Toggle Modal Open
    function showModal(targetModal){
        let modal = document.getElementById(targetModal);
        if(modal){
            modal.classList.add('show');
            currentModal = modal;
            document.body.classList.add('modal-open');
        }
    }

    // * Toggle Modal Close
    function hideModal(){
        currentModal.classList.remove('show');
        document.body.classList.remove('modal-open');
    }
}

window.addEventListener('DOMContentLoaded', function () {
    modalManager();
});