import {ajaxGetHandler, ajaxPostHandler} from './ajax.js';
import {search} from './search.js';
import {templateReplace} from './razer.js';
import {updateViewPlaylist} from "./updateView.js";

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
    getData() {
        return [...this.playlistSongList.keys()];
    }

    // * Get Element Object
    getElement(id) {
        return this.playlistSongList.get(id);
    }

    // * Key Exists 
    keyExists(id) {
        return this.playlistSongList.has(id);
    }

    // * Create Element Item Preview
    playListPreviewItem(audioID) {
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

}

function getPlayListTemplate() {
    let tableHeaders = `<thead><tr>
    <th>Title</th>
    <th>Album</th>
    <th>Date Added</th>
    <th></th>
    <th>Time</th>
  </tr></thead>`;
    let playListTemplate = `<div class="current-banner">
                                <img src="/static/img/abstract-1.jpg">
                                   <div class="playlist-view-wrapper">
                                        <div class="playlist-view-info">
                                            <div class="playlist-view-left">
                                                <div class="playlist-view-cover"></div>
                                            </div>               
                                            <div class="playlist-view-right">
                                                <div class="playlist-view-type">PUBLIC PLAYLIST</div>
                                                <div class="playlist-view-name">{{viewName}}</div>
                                                <div class="playlist-view-extra"><span class="playlist-user"></span><span class="playlist-song-count">{{song-count}}</span><span class="playlist-total-time"></span></div>
                                           </div>                 
                                        </div>
                                   </div>
                            </div>
                            <table class="playlist-table-container">${tableHeaders}{{table-content}}</table>
                            `;
    return playListTemplate;
}

// * Playlist Add Control
function playListControl() {
    let searchHelper = document.querySelector('.search-modal-helper');
    let playListContainer = document.querySelector('.search-container.playlist');
    let addPlaylistManager = new playListAddManager();
    let addPlaylist = document.querySelector('.create-btn');
    let tempListWrapper = document.querySelector('.temp-playlist-wrapper');
    let playListSidePanel = document.querySelector('.my-playlist-container');
    let playListInput = document.querySelector('.playlist-create-name');
    let mainView = document.querySelector('.music-manager-container');
    let bodyContainer = document.body;

    playListContainer.addEventListener('click', function (evt) {
        searchHelper.classList.add('show');
    });

    document.addEventListener("click", function (evt) {
        let target = evt.target;
        if (target.matches(".control-item.add")) {
            let targetElement = target.closest(".search-item");
            let audioID = targetElement.getAttribute("data-audio-id");
            if (!addPlaylistManager.keyExists(audioID)) {
                addPlaylistManager.add(audioID, targetElement);
                tempListWrapper.append(addPlaylistManager.playListPreviewItem(audioID));
                console.log(targetElement);
                console.log(audioID);
            }
        } else if (target.matches(".control-item.remove")) {
            let targetElement = target.closest(".search-item");
            let audioID = targetElement.getAttribute("data-audio-id");
            addPlaylistManager.remove(audioID);
            targetElement.remove(audioID);
        } else if (target.parentElement.classList.contains('playlist-view-trigger')) {
            let targetElement = target.parentElement;
            let playlistID = targetElement.getAttribute("data-playlist-id");
            if(mainView.classList.contains('show-animation')){
                mainView.classList.remove('show-animation');
            }
            loadPlaylistIntoView(playlistID, targetElement.getAttribute('data-view-name'));
        }
    });

    addPlaylist.addEventListener('click', function (evt) {
        console.log("Sending playlist creation...");
        let data = {
            playlistName : playListInput.value,
            playListItems : addPlaylistManager.getData()
        }
        ajaxPostHandler('/addPlaylist', data).then((data) => {
            console.log(data);
            updatePlaylistWrapper(data, playListSidePanel);
        })
            .catch((error) => {
                console.log("Error trying to Create Playlist...");
                console.log(error);
            });
    });

    search('.search-container.playlist', '#playlistSearchFound', ".search-filter.playlist", "/playlist/search/?");
    loadPlaylists();

    // * Load Users Playlists
    function loadPlaylists() {
        ajaxGetHandler('/loadPlaylists').then((data) => {
            console.log('Loading User Playlists');
            console.log(data);
            updatePlaylistWrapper(data, playListSidePanel);
        })
            .catch((error) => {
                console.log("Error trying to Load Playlists...");
                console.log(error);
            });
    }

    // * PlayList Side-Wrapper Append To
    function updatePlaylistWrapper(data, appendToElem) {
        if(!(Symbol.iterator in Object(data))){
            data = new Array(data);
        }
        data.forEach((entry) => {
            const mapObj = {
                "{{playlist-id}}": entry.Playlist_id,
                "{{playlist-name}}": entry.Name
            }
            let playListTemplate = `<div class="playlist-item playlist-view-trigger" data-playlist-id="{{playlist-id}}" data-view-name="{{playlist-name}}">
            <span class="playlist-item-title">{{playlist-name}}</span>
            <i class="bi bi-chevron-double-right"></i>
        </div>`;

            playListTemplate = templateReplace(playListTemplate, mapObj);
                appendToElem.insertAdjacentHTML('beforeend', playListTemplate);
        });
    }

    // * Load User Playlist Into View
    function loadPlaylistIntoView(playlistID, viewName) {
        console.log('Loading Playlist Into View...');
        let data = new URLSearchParams({
            id: playlistID,
        });
        ajaxGetHandler('/loadPlaylist?' + data).then((data) => {
            console.log(data);
            updateViewPlaylist(data, mainView, viewName, bodyContainer, getPlayListTemplate());
        }).catch((error) => {
            console.log("Error trying to Load Playlist Into View...");
            console.log(error);
        });
    }
};

window.addEventListener('DOMContentLoaded', function (evt) {
    playListControl();


});
