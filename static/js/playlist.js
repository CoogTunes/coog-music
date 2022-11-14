import {ajaxGetHandler, ajaxPostHandler} from './ajax.js';
import {search} from './search.js';
import {templateReplace} from './razer.js';

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


// * Playlist Add Control
function playListControl() {
    let searchHelper = document.querySelector('.search-modal-helper');
    let addPlaylistManager = new playListAddManager();
    let addPlaylist = document.querySelector('.create-btn');
    let tempListWrapper = document.querySelector('.temp-playlist-wrapper');
    let playListSidePanel = document.querySelector('.my-playlist-container');
    let playListInput = document.querySelector('.playlist-create-name');
    let mainView = document.querySelector('.music-manager-container');
    let bodyContainer = document.body;

    playListSearch.addEventListener('click', function (evt) {
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
            loadPlaylistIntoView(playlistID);
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
            let playListTemplate = `<div class="playlist-item playlist-view-trigger" data-playlist-id="{{playlist-id}}">
            <span class="playlist-item-title">{{playlist-name}}</span>
            <i class="bi bi-chevron-double-right"></i>
        </div>`;

            playListTemplate = templateReplace(playListTemplate, mapObj);
            appendToElem.insertAdjacentHTML('beforeend', playListTemplate);
        });
    }

    // * Update View with Playlist
    function updateView(data, viewContainer){
        viewContainer.innerHTML = '';

        // Create elements
        let currentBanner = document.createElement('div');
        let songTable = document.createElement('table');


        // Assign needed classes and components
        currentBanner.classList.add('current-banner');
        currentBanner.innerHTML = '<img src="/static/img/abstract-1.jpg">';
        songTable.classList.add('playlist-table-container');
        let tableHeaders = `<tr>
    <th>Title</th>
    <th>Album</th>
    <th>Date Added</th>
    <th></th>
    <th>Time</th>
  </tr>`;
        songTable.insertAdjacentHTML('beforeend', tableHeaders);

        // Fill the Table with Songs
        data.forEach((entry) => {
            let songItem = `<tr class="table-song-item">
    <td><div class="playlist-item-flex content-wrapper" data-audio-path="{{audio}}" data-music-state="paused"><div class="playlist-img-contain audio-cover"><img src="{{cover}}"></div><div class="song-info-item"><div class="song-info-title">{{song}}</div><div class="song-info-artist">{{artist}}</div></div><div class="buttons playlist">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div></div></td>
    <td>{{album}}</td>
    <td>{{date}}</td>
    <td>{{likes}}</td>
  </tr>`;
            const mapObj = {
                "{{song}}": entry.Title,
                "{{album}}": entry.Album,
                "{{artist}}": entry.Artist,
                "{{date}}": entry.UploadedDate,
                "{{dislikes}}": entry.Dislikes,
                "{{likes}}": entry.Likes,
                "{{cover}}": entry.CoverPath,
                "{{audio}}": entry.SongPath,
            }
            songItem = templateReplace(songItem, mapObj);
            songTable.insertAdjacentHTML('beforeend', songItem);
        });

        // Add Elements To Main View & Update Body View
        if(bodyContainer.classList.contains('index-home')){
            bodyContainer.classList.remove('index-home');
            bodyContainer.classList.add('index-view');
        }

        mainView.append(currentBanner);
        mainView.append(songTable);
    }

    // * Load User Playlist Into View
    function loadPlaylistIntoView(playlistID) {
        console.log('Loading Playlist Into View...');
        let data = new URLSearchParams({
            id: playlistID,
        });
        ajaxGetHandler('/loadPlaylist?' + data).then((data) => {
            console.log(data);
            updateView(data, mainView);
        }).catch((error) => {
            console.log("Error trying to Load Playlist Into View...");
            console.log(error);
        });
    }
};

window.addEventListener('DOMContentLoaded', function (evt) {
    playListControl();
});
