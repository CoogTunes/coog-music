import { ajaxPostHandler } from "./ajax.js";

function playlistAddToControl() {
  document.addEventListener("mouseover", function (evt) {
    let target = evt.target;

    if (target.matches(".control-playlist-item.add-to-a-playlist.table-view")) {
      let parentElement =
        target.closest(".content-wrapper") ??
        target.closest(".table-song-item");
      console.log(parentElement);
      console.log(parentElement.getAttribute("data-song-id"));
      populatePlaylists(parentElement.getAttribute("data-song-id"), target);
    }
  });

  function populatePlaylists(songID, target){
    let myPlayLists = document.querySelector('.my-playlist-container').querySelectorAll('.playlist-item');

    if(!target.querySelector('.control-playlist-list')){
      let dynamicPlayList = document.createElement('div');
      dynamicPlayList.classList.add('control-playlist-list');
      myPlayLists.forEach((playlist) => {
        let playListItem = document.createElement('div');
        let playListName = playlist.getAttribute('data-view-name');
        let playListID = playlist.getAttribute('data-playlist-id');
        playListItem.setAttribute('data-view-name', playListName);
        playListItem.setAttribute('data-playlist-id', playListID);
        playListItem.innerHTML = playListName;
        playListItem.addEventListener('click', function () {
          addSong(songID, playListID, target);
        });
        dynamicPlayList.append(playListItem);
      });
  
      target.append(dynamicPlayList);
    }

    logs();
    function logs(){
      console.log(songID);
      console.log(target);
      console.log(myPlayLists);
    }

    function addSong(songId = null, playlistId = null, parentElement = null) {
      console.log("Attempting to add song...");
      console.log("Song ID: " + songId);
      console.log("PlayList ID: " + playlistId);
      let data = {
        playlistID: playlistId,
        songID: songId,
      };
      ajaxPostHandler("/addSong", data)
        .then((data) => {
          console.log("Adding song to playlist...");
          parentElement.removeChild(dynamicPlayList);
        })
        .catch((error) => {
          console.log("Error adding song...");
          console.log(error);
        });
    }
  }
}

window.addEventListener("DOMContentLoaded", function (evt) {
  playlistAddToControl();
});
