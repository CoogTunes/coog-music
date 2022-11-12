import { ajaxGetHandler, ajaxPostHandler } from './ajax.js';

class playListAddManager {
    constructor() {
        this.playlistSongList = new Map();
    }

    // * Add a check
    add(id, elem) {
        this.playlistSongList.set(id, elem);
    }

    // * Add a check
    remove(id) {
        this.playlistSongList.delete(id);
    }

    // * Gather Data
    getData(){
        return [...this.playlistSongList.keys()];
    }

};


// * Playlist Add Control
function playListControl() {
    let playListSearch = document.querySelector('#playListSearch');
    let searchHelper = document.querySelector('.search-modal-helper');
    const addClassList = ['control-item', 'add'];
    let addPlaylistManager = new playListAddManager();
    let addPlaylist = document.querySelector('.create-btn');
    let playListInput = document.querySelector('.playlist-create-name');
    let playListName = null;

    playListInput.addEventListener('input', function (evt){
        playListName = evt.target.value;
    });

    playListSearch.addEventListener('click', function (evt) {
        searchHelper.classList.add('show');
    });

    document.addEventListener('click', function (evt) {
        let target = evt.target;
        if(addClassList.some(className => target.classList.contains(className))) {
            let targetElement = target.closest('.search-item');
            let audioID = targetElement.getAttribute('data-audio-id');
            addPlaylistManager.add(audioID, targetElement);
            console.log(targetElement);
            console.log(audioID);
        }
    });

    addPlaylist.addEventListener('click', function (evt) {
        let data = {
            playlistName : playListName,
            playList : addPlaylistManager.getData(),
        }
        ajaxPostHandler('/addPlaylist', data);
        console.log(data);
        console.log("Sending playlist creation...");
    });
}

window.addEventListener('DOMContentLoaded', function (evt) {
    playListControl();
});