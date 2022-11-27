import { ajaxPostHandler } from "./ajax.js";

function playlistDeleteSongControl() {
  document.addEventListener("click", function (evt) {
    let target = evt.target;

    if (target.matches(".control-playlist-item.remove-from-playlist.table-view")) {
      let parentElement =
        target.closest(".content-wrapper") ??
        target.closest(".table-song-item");
        let songName = parentElement.getElementsByClassName('song-info-title')[0].innerHTML
        if (confirm('Are you sure you want to delete ' + songName + '?')) {
          removeSong(
            parentElement.getAttribute("data-song-id"),
            parentElement.getAttribute("data-playlist-id"),
            parentElement
          );
        }
    }
  });

  function removeSong(songId, playlistId, parentElement) {
    console.log("Attempting to remove song...");
    let data = {
      playlistID: playlistId,
      songID: songId,
    };
    ajaxPostHandler("/deleteSong", data)
      .then((data) => {
        removeTableElement(parentElement);
        console.log("Deleting song from playlist...");
      })
      .catch((error) => {
        console.log("Error deleting song...");
        console.log(error);
      });
  }

  function removeTableElement(tableElement){
    let parentNode = tableElement.parentElement;
    parentNode.removeChild(tableElement);
  }
}

window.addEventListener("DOMContentLoaded", function (evt) {
  playlistDeleteSongControl();
});
