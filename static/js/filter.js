import {ajaxGetHandler} from './ajax.js';
import {updateTableView} from './updateView.js';

class filterList {
    constructor() {
        this.filterList = new Map();
    }
    // * Add a check
    add(id, elem) {
        this.filterList.set(id, elem);
    }

    // * Get Filters -> Re-work to get the current selected filters
    getFilters(){
        return this.filterList;
    }

    // * Add a check
    remove(id) {
        this.filterList.delete(id);
    }
};

export let filterControl = {
    filterManager :  null,
    mainView : null,
    tableContainer : null,
    total : null,

    filterControlInit : function filterControlInit(){
        filterControl.filterManager =  new filterList();
        filterControl.mainView = document.querySelector('.music-manager-container');
    },

    filterListen: function filterListen(){
        document.addEventListener('click', function (evt){
            let target = evt.target;

            if(target.classList.contains('filter-item')){
                let parent = target;
                let filterType = parent.querySelector('span').innerHTML;
                filterControl.toggleSelection(parent);
                if(filterControl.parentContains(parent)){
                    filterControl.filterManager.add(filterType, parent);
                }
                else
                    filterControl.filterManager.remove(filterType);
            }
            else if(target.parentElement.classList.contains('filter-item') && target.tagName.toLowerCase() != 'input'){
                let parent = target.parentElement;
                let filterType = parent.querySelector('span').innerHTML;
                filterControl.toggleSelection(parent);
                if(filterControl.parentContains(parent)){
                    filterControl.filterManager.add(filterType, parent);
                }
                else
                    filterControl.filterManager.remove(filterType);
            }
            else if(target.classList.contains('submit-filter')){
                filterControl.getFilterData();
            }
        });
    },

    getFilterData : function getFilterData(){
        console.log('Loading Filter Data Into View...');
        let data = filterControl.parseFilterData(filterControl.filterManager.getFilters());
        console.log(data);
        ajaxGetHandler('/filters?' + new URLSearchParams(data)).then((data) => {
            console.log(data);
            filterControl.setTableContainer();
            updateTableView(data, filterControl.tableContainer, filterControl.total, filterControl.filterManager.getFilters());
        }).catch((error) => {
            console.log("Error trying to Load Playlist Into View...");
            console.log(error);
        });
    },

    parseFilterData : function parseFilterData(filterData){
        let data = {};
        filterData.forEach((value, key) => {
            let input = value.querySelector('input');
            if(input)
                data[key] = input.value;
            else
                data[key] = "true";

        });
        return data;
    },

    toggleSelection: function toggleSelection(target) {
        target.classList.toggle('selected');
    },

    parentContains: function parentContains(parent){
        if(parent.classList.contains('selected')){
            return true;
        }
        return false;
    },

    setTableContainer: function setTableContainer(){
        filterControl.tableContainer = this.mainView.querySelector('.playlist-table-container');
    },

    setTotalCount: function setTotalCount(){
        filterControl.total = this.mainView.querySelector('.playlist-song-count');
    },
};

