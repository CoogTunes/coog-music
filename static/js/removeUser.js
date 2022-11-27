import { ajaxPostHandler } from "./ajax.js";

function adminRemoveControl() {
    document.addEventListener("click", function (evt) {
        let target = evt.target;

        if (target.matches(".bi.bi-trash-fill.admin-remove-user")) {
            let indexToRemove = Array
                        .from(document.getElementsByClassName('table-item-flex'))
                        .findIndex(x => x.innerHTML == target.getAttribute("data-user-id"));
            let username = document.getElementsByClassName('table-item')[indexToRemove].getElementsByTagName('td')[1].innerHTML;
            if (confirm('Are you sure you want to delete the user, ' + username + ', permanently? This will delete all information related to them.')) {
                    removeUser(target.getAttribute("data-user-id"), target);
                }
        }
    });

    function removeUser(userId, target) {
        console.log("Attempting delete user...");
        let data = {
            userID: userId,
        };
        ajaxPostHandler("/deleteUser", data)
            .then((data) => {
                deleteUserFrontEnd(userId, target);
                console.log("Deleting user...");
            })
            .catch((error) => {
                console.log("Error deleting user...");
                console.log(error);
            });
    }

    function deleteUserFrontEnd(userId, target){
         let tableItem = target.closest('.table-item');
         let parent = tableItem.parentElement;
         parent.removeChild(tableItem);
    }
}

window.addEventListener("DOMContentLoaded", function (evt) {
    adminRemoveControl();
});