function playlistControl() {
  let currentControl = null;

  document.addEventListener("click", function (evt) {
    let target = evt.target;

    if (target.matches(".control-playlist-trigger") || target.matches('.main-playlist-control-trigger')) {
      togglePlaylistControl(target);
    }
    else if(!target.matches('.control-playlist-container') && !target.closest('.control-playlist-container') && currentControl != null) {
      currentControl.classList.remove('show');
    }
  });

  function togglePlaylistControl(target){
    let playListControl = target.nextElementSibling;
    if(currentControl !== null && currentControl !== playListControl){
      currentControl.classList.remove('show');
    }
    playListControl.classList.toggle('show');
    currentControl = playListControl;
  }
}

window.addEventListener("DOMContentLoaded", function (evt) {
  playlistControl();
});
