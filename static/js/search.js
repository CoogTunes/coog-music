import { ajaxGetHandler } from './ajax.js';
import { templateReplace } from './razer.js';
import { isEmpty } from './helpers.js';

export function search(searchWrapper, playlistSearchContainer, filterList, path) {
    let searchContainer = document.querySelector(searchWrapper);
    let searchInput = searchContainer.querySelector("input");
    let inputClear = searchContainer.querySelector("i.input-clear");
    let ajaxSearchContainer = document.querySelector(playlistSearchContainer);
    let filterOptions = document.querySelector(filterList).querySelectorAll('.filter-type');
    let inputEvent = new Event("input");
    let searchString = "";
    let currentFilter = [...filterOptions].filter(filter => filter.classList.contains('selected'))[0];

    function getFilterValue(){
        return currentFilter.querySelector('input').value;
    }

    let ajaxSearchControl = {
        setAjaxContent : function(data = null){
            switch(getFilterValue(currentFilter)){
                case "artist":
                    ajaxSearchControl.setArtistContent(data);
                    break;
                case "album":
                    ajaxSearchControl.setAlbumContent(data);
                    break;
                case "song":
                    ajaxSearchControl.setSongContent(data);
                    break;
            }
        },

        setArtistContent : function(data){

        },

        setSongContent : function(data){
            data.forEach((entry) =>{
                const mapObj = {
                    "{{audio-id}}": entry.Song_id,
                    "{{title}}": entry.Title,
                    "{{sub-title}}": entry.Artist_name,
                    "{{img}}" : entry.CoverPath,
                }
                let songItemTemplate = `<div class="search-item" data-search-item="Song" data-audio-id="{{audio-id}}">
          <div class="search-item-img">
            <img src="{{img}}">
          </div>
          <div class="search-item-info">
            <div class="search-item-title">{{title}}</div>
            <div class="search-item-sub-title">{{sub-title}}</div>
          </div>
          <div class="search-item-control">
            <div class="control-item add">Add</div>
          </div>
        </div>`;
                songItemTemplate = templateReplace(songItemTemplate, mapObj);
                ajaxSearchContainer.insertAdjacentHTML('beforeend', songItemTemplate);
            });
        },

        setAlbumContent : function(data){
            data.forEach((entry) =>{
                const mapObj = {
                    "{{audio-id}}": entry.Song_id,
                    "{{title}}": entry.Title,
                    "{{sub-title}}": entry.Artist_name,
                    "{{img}}" : entry.CoverPath,
                }
                let songItemTemplate = `<div class="search-item" data-search-item="Song" data-audio-id="{{audio-id}}">
          <div class="search-item-img">
            <img src="{{img}}">
          </div>
          <div class="search-item-info">
            <div class="search-item-title">{{title}}</div>
            <div class="search-item-sub-title">{{sub-title}}</div>
          </div>
          <div class="search-item-control">
            <div class="control-item add">Add</div>
          </div>
        </div>`;
                songItemTemplate = templateReplace(songItemTemplate, mapObj);
                ajaxSearchContainer.insertAdjacentHTML('beforeend', songItemTemplate);
            });

        }


    }

    // * Allows for Multiple Filters => Future support
    function parseFilters(optionList){
        let selectedOptions = [...optionList].filter(filter => filter.classList.contains('selected'));
        let filterValues = [];

        selectedOptions.forEach(filter => {
            filterValues.push(filter.querySelector('input').value);
        });

        return filterValues;
    }

    function searchAjax(stringToFind, filterList = null) {
        let data = new URLSearchParams({
            strTarget: stringToFind,
            filters: filterList
        });

        console.log(data);

        // * Init Ajax Handlers
        ajaxGetHandler(path + data)
            .then((data) => {
                ajaxSearchContainer.innerHTML = "";
                console.log("Retrieving results...");
                if(!isEmpty(data))
                    ajaxSearchControl.setAjaxContent(data);
            })
            .catch((error) => {
                console.log("Error trying to Get Ajax Search Data.");
                console.log(error);
            });
    }

    filterOptions.forEach(filter => {
        filter.addEventListener("click", () => {
            if(currentFilter !== filter){
                currentFilter.classList.remove('selected');
            }
            filter.classList.add('selected');
            currentFilter = filter;
            if(searchString.trim().length > 0 && typeof searchString === "string"){
                searchAjax(encodeURIComponent(searchString), getFilterValue(currentFilter));
            }
        });
    });

    searchInput.addEventListener("input", function (evt) {
        searchString = evt.target.value;
        if (typeof searchString === "string" && searchString.trim().length === 0) {
            ajaxSearchContainer.classList.remove("searching");
        } else {
            ajaxSearchContainer.classList.add("searching");
            searchAjax(encodeURIComponent(searchString), currentFilter.querySelector('input').value);
        }
    });

    inputClear.addEventListener("click", function (evt) {
        searchInput.value = "";
        searchInput.dispatchEvent(inputEvent);
    });
}
