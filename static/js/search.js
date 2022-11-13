import { ajaxGetHandler } from './ajax.js';

export function search(searchWrapper, playlistSearchContainer, filterList, path) {
    let searchContainer = document.querySelector(searchWrapper);
    let searchInput = searchContainer.querySelector("input");
    let inputClear = searchContainer.querySelector("i.input-clear");
    let ajaxSearchContainer = document.querySelector(playlistSearchContainer);
    let filterOptions = document.querySelector(filterList).querySelectorAll('.filter-type');
    let inputEvent = new Event("input");

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
                console.log("Ajax Search Data: ");
                console.log(data);
                //console.log(setAjaxContent());
            })
            .catch((error) => {
                console.log("Error trying to Get Ajax Search Data.");
                console.log(error);
            });
    }

    filterOptions.forEach(filter => {
        filter.addEventListener("click", () => {
            filter.classList.toggle('selected');
        });
    });

    searchInput.addEventListener("input", function (evt) {
        let searchString = evt.target.value;
        if (typeof searchString === "string" && searchString.trim().length === 0) {
            ajaxSearchContainer.classList.remove("searching");
        } else {
            ajaxSearchContainer.classList.add("searching");
            searchAjax(encodeURIComponent(searchString), parseFilters(filterOptions));
        }
    });

    inputClear.addEventListener("click", function (evt) {
        searchInput.value = "";
        searchInput.dispatchEvent(inputEvent);
    });
}