// * Update View with Playlist
import {templateReplace} from "./razer.js";
import {dateParse} from "./date.js";
import {songCount} from "./helpers.js";

function tableFilterHeaders(filterList){
    let tableHeaders = "";
    let filterOption = "";
    if(filterList.has('users')){
        tableHeaders = `<thead><tr>
                            <th>Title</th>
                            <th>Album</th>
                            <th>Date Added</th>
                            <th></th>
                            <th>Time</th>
                            </tr></thead>`;
        filterOption = 'users';
    }
    else if(filterList.has('artists')){
        tableHeaders = `<thead><tr>
                            <th>UserID</th>
                            <th>Username</th>
                            <th>First Name</th>
                            <th>Last Name</th>
                            <th>Admin Level</th>
                            <th>Joined Date</th>
                            <th>Playlist Count</th>
                            <th></th>
                          </tr></thead>`;
        filterOption = 'artists';
    }
    else {
        tableHeaders = `<thead><tr>
                            <th>Title</th>
                            <th>Album</th>
                            <th>Date Added</th>
                            <th></th>
                            <th>Time</th>
                            </tr></thead>`;
    }

    return {tableHeaders, filterOption};
}

function tableItemTemplates(filterValue){
    switch (filterValue){
        case 'users':
            return "";
        case 'artist':
            return "";
        default:
            return `<tr class="table-song-item">
    <td><div class="playlist-item-flex content-wrapper" data-audio-path="{{audio}}" data-music-state="paused"><div class="playlist-img-contain audio-cover"><img src="{{cover}}"></div><div class="song-info-item"><div class="song-info-title">{{song}}</div><div class="song-info-artist">{{artist}}</div></div><div class="buttons playlist">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div></div></td>
    <td>{{album}}</td>
    <td>{{date}}</td>
    <td>{{likes}}</td>
    <td>Time</td>
  </tr>`;

    }
}

export function updateViewPlaylist(data, viewContainer, viewName, bodyContainer, htmltemplate){
    viewContainer.innerHTML = '';

    // Fill the Table with Songs
    let tableHTML = '';
    data.forEach((entry) => {
        const mapObj = {
            "{{song}}": entry.Title,
            "{{album}}": entry.Album,
            "{{artist}}": entry.Artist ?? entry.Artist_name,
            "{{date}}": dateParse((entry.UploadedDate ?? entry.Uploaded_date)),
            "{{dislikes}}": entry.Dislikes ?? "",
            "{{likes}}": entry.Likes ?? entry.Total_likes,
            "{{cover}}": entry.CoverPath,
            "{{audio}}": entry.SongPath,
            "{{songID}}" : entry.SongID,
        }
        let songItem = `<tr class="table-song-item" data-song-id="{{songID}}">
    <td><div class="playlist-item-flex content-wrapper" data-audio-path="{{audio}}" data-music-state="paused"><div class="playlist-img-contain audio-cover"><img src="{{cover}}"></div><div class="song-info-item"><div class="song-info-title">{{song}}</div><div class="song-info-artist">{{artist}}</div></div><div class="buttons playlist">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div></div></td>
    <td>{{album}}</td>
    <td>{{date}}</td>
    <td><div class="like-container"><i class="bi bi-heart"></i><div class="current-likes">{{likes}}</div></div></td>
    <td>Time</td>
  </tr>`;
        songItem = templateReplace(songItem, mapObj);
        tableHTML += songItem;
    });

    htmltemplate = templateReplace(htmltemplate,     {"{{viewName}}": viewName})
    htmltemplate = templateReplace(htmltemplate,     {"{{table-content}}": tableHTML});
    htmltemplate = templateReplace(htmltemplate, {"{{song-count}}": songCount(data)})
    viewContainer.insertAdjacentHTML('beforeend', htmltemplate);

    // Add Elements To Main View & Update Body View
    if(bodyContainer.classList.contains('index-home')){
        bodyContainer.classList.remove('index-home');
        bodyContainer.classList.add('index-view');
    }
    viewContainer.classList.add('show-animation');
}

export function updateViewDiscover(data, viewContainer, viewName, bodyContainer, htmltemplate){
    viewContainer.innerHTML = '';

    console.log(data);
    // Fill the Table with Songs
    let tableHTML = '';
    data.forEach((entry) => {
        const mapObj = {
            "{{song}}": entry.Title,
            "{{album}}": entry.Album,
            "{{artist}}": entry.Artist ?? entry.Artist_name,
            "{{date}}": dateParse((entry.UploadedDate ?? entry.Uploaded_date)),
            "{{dislikes}}": entry.Dislikes ?? "",
            "{{likes}}": entry.Likes ?? entry.Total_likes,
            "{{cover}}": entry.CoverPath,
            "{{audio}}": entry.SongPath,
        }
        let songItem = `<tr class="table-song-item">
    <td><div class="playlist-item-flex content-wrapper" data-audio-path="{{audio}}" data-music-state="paused"><div class="playlist-img-contain audio-cover"><img src="{{cover}}"></div><div class="song-info-item"><div class="song-info-title">{{song}}</div><div class="song-info-artist">{{artist}}</div></div><div class="buttons playlist">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div></div></td>
    <td>{{album}}</td>
    <td>{{date}}</td>
    <td>{{likes}}</td>
    <td>Time</td>
  </tr>`;
        songItem = templateReplace(songItem, mapObj);
        tableHTML += songItem;
    });

    htmltemplate = templateReplace(htmltemplate,     {"{{viewName}}": viewName});
    htmltemplate = templateReplace(htmltemplate,     {"{{table-content}}": tableHTML});
    htmltemplate = templateReplace(htmltemplate, {"{{song-count}}": songCount(data)});
    viewContainer.insertAdjacentHTML('beforeend', htmltemplate);

    // Add Elements To Main View & Update Body View
    if(bodyContainer.classList.contains('index-home')){
        bodyContainer.classList.remove('index-home');
        bodyContainer.classList.add('index-view');
    }

    viewContainer.classList.add('show-animation');
}

export function updateTableView(data, tableContainer, songTotal, filterList, mainView){
    if(songTotal)
        songTotal.innerHTML = songCount(data);
    tableContainer.innerHTML = '';

    let tableFilterHeader = tableFilterHeaders(filterList);
    let tableHead = tableFilterHeader.tableHeaders;
    let filterValue = tableFilterHeader.filterOption;
    let htmltemplate = `${tableHead}{{table-content}}`;
    let tableHTML = '';

    data.forEach((entry) => {
        const mapObj = {
            "{{song}}": entry.Title ?? entry.Song_title ?? "",
            "{{album}}": entry.Album ?? entry.Album_name ?? "",
            "{{artist}}": entry.Artist ?? entry.Artist_name ?? "",
            "{{date}}": dateParse((entry.UploadedDate ?? entry.Uploaded_date ?? "")),
            "{{dislikes}}": entry.Dislikes ?? "",
            "{{likes}}": entry.Likes ?? entry.Total_likes ?? "",
            "{{cover}}": entry.CoverPath ?? entry.Cover_path ?? "",
            "{{audio}}": entry.SongPath ?? entry.Song_path ?? "",
            "{{time}}": entry.duration ?? "",
        }
        let songItem = tableItemTemplates(filterValue);
        songItem = templateReplace(songItem, mapObj);
        tableHTML += songItem;
    });

    htmltemplate = templateReplace(htmltemplate,     {"{{table-content}}": tableHTML});
    // htmltemplate = templateReplace(htmltemplate, {"{{song-count}}": songCount(Object.keys(data).length)});
    tableContainer.insertAdjacentHTML('beforeend', htmltemplate);
}

export function updateViewAdminControl(data, viewContainer, viewName, bodyContainer, htmltemplate){
    viewContainer.innerHTML = '';

    console.log(data);
    // Fill the Table with Users
    let tableHTML = '';
    data.forEach((entry) => {
        const mapObj = {
            "{{song}}": entry.Title ?? entry.Song_title ?? "",
            "{{album}}": entry.Album ?? entry.Album_name ?? "",
            "{{artist}}": entry.Artist ?? entry.Artist_name ?? "",
            "{{date}}": dateParse((entry.UploadedDate ?? entry.Uploaded_date ?? "")),
            "{{dislikes}}": entry.Dislikes ?? "",
            "{{likes}}": entry.Likes ?? entry.Total_likes ?? "",
            "{{cover}}": entry.CoverPath ?? entry.Cover_path ?? "",
            "{{audio}}": entry.SongPath ?? entry.Song_path ?? "",
            "{{time}}": entry.duration ?? "",
            "{{artistID}}": entry.Album_ID ?? "",
            "{{albumID}}": entry.Album_ID ?? "",
            "{{totalPlays}}": entry.Total_Plays ?? "",
            "{{username}}": entry.Username ?? "",
            "{{firstname}}": entry.First_name ?? "",
            "{{lastname}}": entry.Last_name ?? "",
            "{{adminlevel}}": entry.Admin_level ?? "",
            "{{joinedDate}}": entry.JoinedDate ?? "",
            "{{playlistCount}}": entry.Playlist_count ?? "",
            "{{userid}}": entry.User_id,
        }
        let songItem = `<tr class="table-item">
                        <td>{{userid}}</td>  
                         <td>{{username}}</td>                     
                        <td>{{firstname}}</td>
                        <td>{{lastname}}</td>
                        <td>{{adminlevel}}</td>
                        <td>{{joinedDate}}</td>
                         <td>{{playlistCount}}</td> 
                      </tr>`;
        songItem = templateReplace(songItem, mapObj);
        tableHTML += songItem;
    });

    htmltemplate = templateReplace(htmltemplate,     {"{{viewName}}": viewName});
    htmltemplate = templateReplace(htmltemplate,     {"{{table-content}}": tableHTML});
    htmltemplate = templateReplace(htmltemplate, {"{{song-count}}": songCount(data)});
    viewContainer.insertAdjacentHTML('beforeend', htmltemplate);

    // Add Elements To Main View & Update Body View
    if(bodyContainer.classList.contains('index-home')){
        bodyContainer.classList.remove('index-home');
        bodyContainer.classList.add('index-view');
    }
    viewContainer.classList.add('show-animation');
}

export function updateViewHomeControl(data, viewContainer, viewName, bodyContainer, htmltemplate){
    viewContainer.innerHTML = '';

    console.log(data);
    // Fill the Table with Users
    let homeSongs = '';
    data.forEach((entry) => {
        const mapObj = {
            "{{song}}": entry.Title,
            "{{album}}": entry.Album,
            "{{artist}}": entry.Artist ?? entry.Artist_name,
            "{{date}}": dateParse((entry.UploadedDate ?? entry.Uploaded_date)),
            "{{dislikes}}": entry.Dislikes ?? "",
            "{{likes}}": entry.Likes ?? entry.Total_likes,
            "{{cover}}": entry.CoverPath,
            "{{audio}}": entry.SongPath,
        }
        let songItem = `<div class="content">
                  <div class="content-wrapper" data-audio-path="{{audio}}" data-music-state="paused">
                    <div class="content-img audio-cover">
                      <img src="{{cover}}">
                      <div class="buttons">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div>
                    </div>
                    <div class="song-info-title">{{song}}</div>
                    <div class="song-info-artist">{{artist}}</div>
                  </div>
                </div>`;
        songItem = templateReplace(songItem, mapObj);
        homeSongs += songItem;
    });

    htmltemplate = templateReplace(htmltemplate,     {"{{song-content}}": homeSongs});
    viewContainer.insertAdjacentHTML('beforeend', htmltemplate);

    // Add Elements To Main View & Update Body View
    if(bodyContainer.classList.contains('index-view')){
        bodyContainer.classList.remove('index-view');
        bodyContainer.classList.add('index-home');
    }
    viewContainer.classList.add('show-animation');
}