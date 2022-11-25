import { ajaxPostHandler } from "./ajax.js";

function adminRemoveControl() {
    document.addEventListener("click", function (evt) {
        let target = evt.target;

        if (target.matches(".bi.bi-trash-fill.admin-remove-user")) {
            removeUser(target.getAttribute("data-user-id"), target);
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