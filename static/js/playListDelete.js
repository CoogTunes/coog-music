import { ajaxPostHandler } from "./ajax.js";
import { redirectHome } from "./redirects.js";

function playlistDeleteControl() {
  document.addEventListener("click", function (evt) {
    let target = evt.target;

    if (target.matches(".control-playlist-item.delete-playlist.table-view")) {
      let parentElement = target.closest(".main-playlist-control");
      removeSong(parentElement.getAttribute("data-playlist-id"));
    }
  });

  function removeSong(playlistId) {
    console.log("Attempting delete playlist...");
    let data = {
      playlistID: playlistId,
    };
    ajaxPostHandler("/deletePlaylist", data)
      .then((data) => {
        redirectHome();
        deletePlayListFrontEnd(playlistId);
        console.log("Deleting playlist...");
      })
      .catch((error) => {
        console.log("Error deleting playlist...");
        console.log(error);
      });
  }

  function deletePlayListFrontEnd(playListID){
    let playListContainer = document.querySelector('.my-playlist-container');
    let playListTarget = playListContainer.querySelector(`.playlist-item[data-playlist-id="${playListID}"]`);
    console.log(playListContainer);
    console.log(playListTarget);
    playListContainer.removeChild(playListTarget);
  }
}

window.addEventListener("DOMContentLoaded", function (evt) {
  playlistDeleteControl();
});
