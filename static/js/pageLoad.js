import {ajaxGetHandler} from './ajax.js';
import {updateViewDiscover, updateViewAdminControl, updateViewHomeControl} from "./updateView.js";
import {filterControl} from './filter.js';

function getDiscoverTemplate(){
    let tableHeaders = `<thead><tr>
    <th>Title</th>
    <th>Album</th>
    <th>Date Added</th>
    <th></th>
    <th>Time</th>
  </tr></thead>`;
    let template = `<div class="current-banner">
                                <img src="/static/img/abstract-2-sized.jpg">
                                   <div class="playlist-view-wrapper">
                                        <div class="playlist-view-info">
                                            <div class="playlist-view-left">
                                                <div class="playlist-view-cover"></div>
                                            </div>               
                                            <div class="playlist-view-right">
                                                <div class="playlist-view-name alt">{{viewName}}</div>
                                                <div class="playlist-view-extra"><span class="playlist-song-count">{{song-count}}</span><span class="playlist-total-time">{{total-time}}</span></div>
                                           </div>                 
                                        </div>
                                   </div>
                            </div>
                            <div class="filter-container">
                                <div class="filter-item"><span>plays</span></div>
                                <div class="filter-item"><span>likes</span></div>
                                <div class="filter-item"><span class="date-filter-span">min</span><input type="input"></div>
                                <div class="filter-item"><span class="date-filter-span">max</span><input type="input"></div>
                                <div class="filter-item"><span class="date-filter-span">start</span><input type="date"></div>
                                <div class="filter-item"><span class="date-filter-span">end</span><input type="date"></div>
                                <div class="submit-filter">Submit</div>
                            </div>
                           </div>
                            <table class="playlist-table-container">${tableHeaders}{{table-content}}</table>
                            `;
    return template;
}

function getAdminControlTemplate() {
        let tableHeaders = `<thead><tr>
    <th>UserID</th>
    <th>Username</th>
    <th>First Name</th>
    <th>Last Name</th>
    <th>Admin Level</th>
    <th>Joined Date</th>
    <th>Playlist Count</th>
    <th></th>
  </tr></thead>`;
        let template = `<div class="current-banner">
                                <img src="/static/img/abstract-2-sized.jpg">
                                   <div class="playlist-view-wrapper">
                                        <div class="playlist-view-info">
                                            <div class="playlist-view-left">
                                                <div class="playlist-view-cover"></div>
                                            </div>               
                                            <div class="playlist-view-right">
                                                <div class="playlist-view-name alt">{{viewName}}</div>
                                                <div class="playlist-view-extra"><span class="playlist-total-time"></span></div>
                                           </div>                 
                                        </div>
                                   </div>
                            </div>
                            <div class="filter-container">
                                <div class="filter-item"><span>artists</span></div>
                                <div class="filter-item"><span>users</span></div>
<!--                                <div class="filter-item"><span class="date-filter-span">min</span><input type="input"></div>-->
<!--                                <div class="filter-item"><span class="date-filter-span">max</span><input type="input"></div>-->
                                <div class="filter-item"><span class="date-filter-span">start</span><input type="date"></div>
                                <div class="filter-item"><span class="date-filter-span">end</span><input type="date"></div>
                                <div class="submit-filter">Submit</div>
                            </div>
                           </div>
                            <table class="playlist-table-container">${tableHeaders}{{table-content}}</table>
                            `;
        return template;
}

function getHomeControlTemplate() {
    let template = `<div class="current-banner"><img src="/static/img/wakanda4ever.png"></div>
                        <div class="playlist-wrapper" data-playlist-id="" data-playlist-name="">
                            <div class="row">   
                                <div class="playlist-title">Top songs<i class="bi bi-layers-fill playlist-icon"></i></div>
                            </div>
                            <div class="row cards-wrapper">
                                {{song-content}}
                            </div>
                        </div>`;
    return template;
}

function pageLoadManager(){
    let mainView = document.querySelector('.music-manager-container');
    let bodyContainer = document.body;

    document.addEventListener('click', function (evt) {
        let target = evt.target;
        if(target.parentElement.classList.contains('page-load-trigger')){
            let parent = target.parentElement;
            let targetPage = parent.getAttribute('data-page-index');
            if(mainView.classList.contains('show-animation')){
                mainView.classList.remove('show-animation');
            }
            loadPage(targetPage, parent.getAttribute('data-view-name'));
        }
    });
    function loadPage(targetPage, viewName, path = '/pageLoad?', ){
        let data = new URLSearchParams({
            index : targetPage
        });
        ajaxGetHandler(path + data)
            .then((data) => {
                console.log("Retrieving page data...");
                const classes = mainView.className.split(" ").filter(c => !c.startsWith('index-'));
                const viewSelector = 'index-' + viewName.toLowerCase();
                mainView.className = classes.join(" ").trim();
                mainView.classList.add(viewSelector);
                if(targetPage === 'discover'){
                    updateViewDiscover(data, mainView, viewName, bodyContainer, getDiscoverTemplate());
                    filterControl.filterControlInit();
                    filterControl.filterListen();
                }
                else if (targetPage === 'admin') {
                    updateViewAdminControl(data, mainView, viewName, bodyContainer, getAdminControlTemplate())
                    filterControl.filterControlInit();
                    filterControl.filterListen();

                }
                else if (targetPage === 'home'){
                    updateViewHomeControl(data, mainView, viewName, bodyContainer, getHomeControlTemplate())
                }
            })
            .catch((error) => {
                console.log("Retrieving page error...");
                console.log(error);
            });
    }
}

window.addEventListener('DOMContentLoaded', function () {
    pageLoadManager();
    // * Simulate Init Load
    document.querySelector('.page-load-trigger.home span').click();
});