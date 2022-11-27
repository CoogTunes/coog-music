import {ajaxGetHandler} from './ajax.js';
import { updateArtistPage } from "./updateView.js";

function getArtistTemplate(){
  let tableHeaders = `<thead>
  <tr>
    <th></th>
    <th></th>
    <th></th>
    <th></th>
    <th></th>
    <th></th>
  </tr>
  </thead>`;
  let artistPageTemplate = `<div class="current-banner">
                                     <img src="/static/img/abtract-6-1.webp">
                                     <div class="artist-view-wrapper">
                                          <div class="artist-view-info" data-artist-id="{{artist}}">
                                              <div class="verified-badge"><div class="white-check"></div><svg role="img" height="24" width="24" class="verified-svg" viewBox="0 0 24 24"><path d="M10.814.5a1.658 1.658 0 012.372 0l2.512 2.572 3.595-.043a1.658 1.658 0 011.678 1.678l-.043 3.595 2.572 2.512c.667.65.667 1.722 0 2.372l-2.572 2.512.043 3.595a1.658 1.658 0 01-1.678 1.678l-3.595-.043-2.512 2.572a1.658 1.658 0 01-2.372 0l-2.512-2.572-3.595.043a1.658 1.658 0 01-1.678-1.678l.043-3.595L.5 13.186a1.658 1.658 0 010-2.372l2.572-2.512-.043-3.595a1.658 1.658 0 011.678-1.678l3.595.043L10.814.5zm6.584 9.12a1 1 0 00-1.414-1.413l-6.011 6.01-1.894-1.893a1 1 0 00-1.414 1.414l3.308 3.308 7.425-7.425z"></path></svg></div>      
                                              <div class="artist-view-name">{{artist}}</div>
                                              <div class="spacer-2"></div>                     
                                          </div>
                                     </div>
                              </div>
                              <div class="artist-page-control">
                              </div>
                              <div class="playlist-wrapper artist-data-container">
                                <div class="row">
                                  <div class="category-title">Popular</div>
                                  <div class="divider alt"></div>
                                </div>
                                <div class="songs-wrapper">
                                  <table class="playlist-table-container artist-song-list">${tableHeaders}{{table-content}}</table>
                                  <div class="see-more-trigger">SEE MORE</div>
                                </div>
                                <div class="row">
                                  <div class="category-title">Discography</div>
                                </div>
                                <div class="tab-container">
                                    <div class="tab-item selected" data-container-target="#artist-albums">Albums</div>
                                    <div class="tab-item" data-container-target="#artist-singles">Singles and EPs</div>
                                </div>
                                <div class="cards-wrapper selected" id="artist-albums">{{artistAlbums}}</div>
                                <div class="cards-wrapper" id="artist-singles">{{artistSingles}}</div>
                              </div>
                              `;
    return artistPageTemplate;
}

function artistPageManager(){
  let mainView = document.querySelector(".music-manager-container");
  let bodyContainer = document.body;

  document.addEventListener('click', function (evt) {
    let target = evt.target;

    if(target.matches('.control-playlist-item.go-to-artist.table-view')){
      if (mainView.classList.contains("show-animation")) {
        mainView.classList.remove("show-animation");
      }
      goToArtistPage(target.getAttribute('data-artist-id'), target.getAttribute('data-artist-name'), '/artistPage?');
    }
  });

  function goToArtistPage(artistID, artistName, path){
    let data = new URLSearchParams({
      artistID : artistID,
      artistName : artistName,
    });
  ajaxGetHandler(path + data)
      .then((data) => {
          console.log("Retrieving artist page...");
          updateArtistPage(data, mainView, artistID, artistName, bodyContainer, getArtistTemplate());
          console.log(data);
      })
      .catch((error) => {
          console.log("Retrieving artist page error...");
          console.log(error);
      });
  }

}

window.addEventListener('DOMContentLoaded', function (evt) {
  artistPageManager();
});