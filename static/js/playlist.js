import { ajaxGetHandler, ajaxPostHandler } from './ajax.js';
import { search } from './search.js';

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

    getElement(id){
        return this.playlistSongList.get(id);
    }

};


// * Playlist Add Control
function playListControl() {
    const addClassList = ['control-item', 'add'];
    let searchHelper = document.querySelector('.search-modal-helper');
    let addPlaylistManager = new playListAddManager();
    let addPlaylist = document.querySelector('.create-btn');
    let tempListWrapper = document.querySelector('.temp-playlist-wrapper');

    playListSearch.addEventListener('click', function (evt) {
        searchHelper.classList.add('show');
    });

    document.addEventListener('click', function (evt) {
        let target = evt.target;
        if(addClassList.some(className => target.classList.contains(className))) {
            let targetElement = target.closest('.search-item');
            let audioID = targetElement.getAttribute('data-audio-id');
            addPlaylistManager.add(audioID, targetElement);
            // tempListWrapper.append(addPlaylistManager.getElement(audioID));
            console.log(targetElement);
            console.log(audioID);
        }
    });

    addPlaylist.addEventListener('click', function (evt) {
        let data = addPlaylistManager.getData();
        ajaxPostHandler('/addPlaylist', data);
        console.log(data);
        console.log("Sending playlist creation...");
    });

    search('.search-container.playlist', '#playlistSearchFound', ".search-filter.playlist","/playlist/search/?");
}

window.addEventListener('DOMContentLoaded', function (evt) {
    playListControl();
});