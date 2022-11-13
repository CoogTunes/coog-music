import { ajaxPostHandler } from './ajax.js';
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

    // * Get Element Object
    getElement(id){
        return this.playlistSongList.get(id);
    }

    // * Key Exists
    keyExists(id){
        return this.playlistSongList.has(id);
    }

    // * Create Element Item Preview
    playListPreviewItem(audioID){
        let previewItem = document.createElement('div');
        previewItem.classList.add('search-item');
        previewItem.setAttribute('data-audio-id', audioID);
        let searchItem = this.getElement(audioID);
        let searchImg = searchItem.querySelector('.search-item-img').cloneNode(true);
        let searchInfo = searchItem.querySelector('.search-item-info').cloneNode(true);
        let searchControl = document.createElement('div');
        searchControl.classList.add('search-item-control');
        searchControl.innerHTML = '<div class="control-item remove">Remove</div>';
        previewItem.append(searchImg, searchInfo, searchControl);

        return previewItem;
    }

};


// * Playlist Add Control
function playListControl() {
    let searchHelper = document.querySelector('.search-modal-helper');
    let addPlaylistManager = new playListAddManager();
    let addPlaylist = document.querySelector('.create-btn');
    let tempListWrapper = document.querySelector('.temp-playlist-wrapper');
    let playListName = document.querySelector('.playlist-create-name');

    playListSearch.addEventListener('click', function (evt) {
        searchHelper.classList.add('show');
    });

    document.addEventListener("click", function (evt) {
        let target = evt.target;
        if (target.matches(".control-item.add")) {
            let targetElement = target.closest(".search-item");
            let audioID = targetElement.getAttribute("data-audio-id");
            if (!addPlaylistManager.keyExists(audioID)) {
                console.log(targetElement);
                console.log(audioID);
                addPlaylistManager.add(audioID, targetElement);
                tempListWrapper.append(addPlaylistManager.playListPreviewItem(audioID));
            }
        } else if (target.matches(".control-item.remove")) {
            let targetElement = target.closest(".search-item");
            let audioID = targetElement.getAttribute("data-audio-id");
            addPlaylistManager.remove(audioID);
            targetElement.remove(audioID);
        }
    });

    addPlaylist.addEventListener('click', function (evt) {
        console.log("Sending playlist creation...");
        let data = {
            playListName : playListName.value,
            playListItems : addPlaylistManager.getData()
        }
        console.log(data);
        ajaxPostHandler('/addPlaylist', data).then((data) => {
            console.log("Successful playlist creation...");
        })
            .catch((error) => {
                console.log("Error trying to Create Playlist...");
                console.log(error);
            });
    });

    search('.search-container.playlist', '#playlistSearchFound', ".search-filter.playlist", "/playlist/search/?");
}

window.addEventListener('DOMContentLoaded', function (evt) {
    playListControl();
});