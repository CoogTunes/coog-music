// * Update View with Playlist
import { templateReplace } from "./razer.js";
import { dateParse } from "./date.js";
import { songCount } from "./helpers.js";

function tableFilterHeaders(filterList) {
  let tableHeaders = `<thead><tr>
                            <th>Title</th>
                            <th>Album</th>
                            <th>Date Added</th>
                            <th>Plays</th>
                            <th>Time</th>
                            </tr></thead>`;;
  let filterOption = "";

  if (filterList.has("users")) {
    tableHeaders = `<thead><tr>
    <th>UserID</th>
    <th>Username</th>
    <th>First Name</th>
    <th>Last Name</th>
    <th>Admin Level</th>
    <th>Joined Date</th>
    <th>Playlist Count</th>
    <th>Liked Songs Count</th>
    <th>Favorite Artist</th>
    <th></th>`;
    filterOption = "users";
  } else if (filterList.has("artists")) {
    tableHeaders = `<thead><tr>
                            <th>Name</th>
                            <th>Artist ID</th>
                            <th>Join Date</th>
                            <th>Num. of Songs</th>
                            <th>Num. of Albums</th>
                            <th>Total Plays</th>
                            <th>Avg Plays</th>
                            <th>Most Liked Song</th>
                            <th></th>
                          </tr></thead>`;
    filterOption = "artists";
  } else if (filterList.has("plays")) {
    filterOption = 'plays';
    tableHeaders = `<thead><tr>
                            <th>Title</th>
                            <th>Album</th>
                            <th>Date Added</th>
                            <th>Plays</th>
                            <th>Time</th>
                            </tr></thead>`;
  } else if (filterList.has("likes")) {
    filterOption = 'likes';
    tableHeaders = `<thead><tr>
                            <th>Title</th>
                            <th>Album</th>
                            <th>Date Added</th>
                            <th>Likes</th>
                            <th>Time</th>
                            </tr></thead>`;
  }

  return { tableHeaders, filterOption };
}

function tableItemTemplates(filterValue) {
  switch (filterValue) {
    case "artists":
      return `<tr class="table-item">
                        <td><div class="table-item-flex">{{name}}</div></td>  
                        <td>{{artistId}}</td>                  
                        <td>{{joinedDate}}</td>
                        <td>{{numSongs}}</td>
                        <td>{{numAlbums}}</td>
                        <td>{{totalPlays}}</td>  
                        <td>{{avgPlays}}</td>
                        <td>{{mostLikedSongs}}</td>
                      </tr>`;
    case "users":
      return `<tr class="table-item">
                        <td><div class="table-item-flex">{{userid}}</div></td>  
                        <td>{{username}}</td>                     
                        <td>{{firstname}}</td>
                        <td>{{lastname}}</td>
                        <td>{{adminlevel}}</td>
                        <td>{{joinedDate}}</td>
                        <td>{{playlistCount}}</td>
                        <td>{{likedSongsCount}}</td>  
                        <td>{{commonArtist}}</td>
                        <td><i class="bi bi-trash-fill admin-remove-user" data-user-id="{{userid}}"></i></td>
                      </tr>`;
    case "likes":
      return `<tr class="table-song-item" data-audio-id="{{songID}}">
    <td><div class="playlist-item-flex content-wrapper" data-audio-path="{{audio}}" data-music-state="paused" data-audio-id="{{songID}}"><div class="playlist-img-contain audio-cover"><img src="{{cover}}"></div><div class="song-info-item"><div class="song-info-title">{{song}}</div><div class="song-info-artist">{{artist}}</div></div><div class="buttons playlist">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div></div></td>
    <td>{{album}}</td>
    <td>{{date}}</td>
    <td class="song-likes-count">{{likes}}</td>
    <td>{{duration}}</td>
  </tr>`;
    case "plays":
      return `<tr class="table-song-item" data-audio-id="{{songID}}">
    <td><div class="playlist-item-flex content-wrapper" data-audio-path="{{audio}}" data-music-state="paused" data-audio-id="{{songID}}"><div class="playlist-img-contain audio-cover"><img src="{{cover}}"></div><div class="song-info-item"><div class="song-info-title">{{song}}</div><div class="song-info-artist">{{artist}}</div></div><div class="buttons playlist">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div></div></td>
    <td>{{album}}</td>
    <td>{{date}}</td>
    <td class="song-play-count">{{totalPlays}}</td>
    <td>{{duration}}</td>
  </tr>`;
    default:
      return `<tr class="table-song-item" data-audio-id="{{songID}}">
    <td><div class="playlist-item-flex content-wrapper" data-audio-path="{{audio}}" data-music-state="paused" data-audio-id="{{songID}}"><div class="playlist-img-contain audio-cover"><img src="{{cover}}"></div><div class="song-info-item"><div class="song-info-title">{{song}}</div><div class="song-info-artist">{{artist}}</div></div><div class="buttons playlist">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div></div></td>
    <td>{{album}}</td>
    <td>{{date}}</td>
    <td class="song-play-count">{{totalPlays}}</td>
    <td>{{duration}}</td>
  </tr>`;
  }
}

export function updateTableView(data, tableContainer, songTotal, filterList) {
  if (songTotal) songTotal.innerHTML = songCount(data);
  tableContainer.innerHTML = "";

  let tableFilterHeader = tableFilterHeaders(filterList);
  let tableHead = tableFilterHeader.tableHeaders;
  let filterValue = tableFilterHeader.filterOption;
  console.log(filterValue);
  let htmltemplate = `${tableHead}{{table-content}}`;
  let tableHTML = "";

  data.forEach((entry) => {
    if (entry && entry.Admin_level) {
      if (entry.Admin_level == 1) {
        entry.Admin_level = 'User'
      } else if (entry.Admin_level == 2) {
        entry.Admin_level = 'Artist'
      }
    }
    const mapObj = {
      "{{likes}}": entry.Likes ?? entry.Total_likes ?? "",
      "{{dislikes}}": entry.Dislikes ?? "",
      "{{songID}}": entry.Song_id ?? "",
      "{{song}}": entry.Title ?? "",
      "{{artistId}}": entry.Artist_id ?? "",
      "{{audio}}": entry.SongPath ?? "",
      "{{cover}}": entry.CoverPath ?? "",
      "{{date}}": dateParse(entry.UploadedDate ?? entry.Uploaded_date) ?? "",
      "{{album}}": entry.Album ?? "",
      "{{albumId}}": entry.Album_id ?? "",
      "{{totalPlays}}": entry.Total_plays ?? entry.Total_Plays ?? "",
      "{{artist}}": entry.Artist ?? entry.Artist_name ?? "",
      "{{duration}}": entry.Duration ?? "",
      "{{username}}": entry.Username ?? "",
      "{{firstname}}": entry.First_name ?? "",
      "{{lastname}}": entry.Last_name ?? "",
      "{{adminlevel}}": entry.Admin_level ?? "",
      "{{joinedDate}}": dateParse(entry.JoinedDate ?? entry.Join_date) ?? "",
      "{{playlistCount}}": entry.Playlist_count ?? "",
      "{{userid}}": entry.User_id,
      "{{likedSongsCount}}": entry.Liked_songs_count ?? "",
      "{{commonArtist}}": entry.Common_artist ?? "",
      "{{avgPlays}}": entry.Avg_Plays ?? "",
      "{{mostLikedSongs}}": entry.Most_liked_song ?? "",
      "{{numSongs}}": entry.Num_songs ?? "",
      "{{numAlbums}}": entry.Num_Albums ?? "",
      "{{name}}": entry.Name ?? "",
    };
    let songItem = tableItemTemplates(filterValue);
    songItem = templateReplace(songItem, mapObj);
    tableHTML += songItem;
  });

  htmltemplate = templateReplace(htmltemplate, {
    "{{table-content}}": tableHTML,
  });
  // htmltemplate = templateReplace(htmltemplate, {"{{song-count}}": songCount(Object.keys(data).length)});
  tableContainer.insertAdjacentHTML("beforeend", htmltemplate);
}

export function updateArtistPage(data, viewContainer, artistID, artistName, bodyContainer, htmltemplate) {
  viewContainer.innerHTML = "";

  let keys = Object.keys(data);
  let albumList = "";
  let singleList = "";
  let tableHTML = "";
  let totalPlays = 0;

  keys.forEach((entry) => {
    let albumName = entry;
    let albumSongs = data[entry];
    let albumCover = data[entry][0].CoverPath;
    let audioPath = data[entry][0].SongPath;
    let albumID = data[entry][0].Album_id;
    let albumSongList = new Array();
    console.log(albumSongs);

    const mapObj = {
      "{{artistId}}": artistID,
      "{{cover}}": albumCover,
      "{{album}}": albumName,
      "{{albumId}}": albumID,
      "{{artist}}": artistName,
    };

    let contentItem = `<div class="content">
                  <div class="content-wrapper" data-artist-name="{{artist}}" data-album-name="{{album}}" data-music-state="paused" data-album-id="{{albumId}}" data-album-songs={}>
                    <div class="content-img audio-cover">
                      <img src="{{cover}}">
                      <div class="buttons">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div>
                    </div>
                    <div class="song-info-title">{{album}}</div>
                    <div class="song-info-artist">{{artist}}</div>
                  </div>
                </div>`;
    

    albumSongs.forEach((entry) => {
      const mapObj = {
        "{{likes}}": entry.Likes ?? entry.Total_likes,
        "{{dislikes}}": entry.Dislikes ?? "",
        "{{songID}}": entry.Song_id,
        "{{song}}": entry.Title,
        "{{artistId}}": entry.Artist_id,
        "{{audio}}": entry.SongPath,
        "{{cover}}": entry.CoverPath,
        "{{date}}": dateParse(entry.UploadedDate ?? entry.Uploaded_date),
        "{{album}}": entry.Album,
        "{{albumId}}": entry.Album_id,
        "{{totalPlays}}": entry.Total_plays,
        "{{artist}}": entry.Artist ?? entry.Artist_name,
        "{{duration}}": entry.Duration,
      };

      const albumSongInfo = {
        songID : entry.Song_id,
        audioPath : entry.SongPath,
        cover : entry.CoverPath,
        artistName : entry.Artist_name,
        songTitle : entry.Title,
      }
      
      albumSongList.push(albumSongInfo);

      totalPlays += entry.Total_plays;

      let songItem = `<tr class="table-song-item" data-song-id="{{songID}}">
        <td><div class="playlist-item-flex content-wrapper" data-audio-path="{{audio}}" data-music-state="paused" data-song-id="{{songID}}"><div class="playlist-img-contain audio-cover"><img src="{{cover}}"></div><div class="song-info-item"><div class="song-info-title">{{song}}</div><div class="song-info-artist">{{artist}}</div></div><div class="buttons playlist">
                            <button><i class="bi bi-play-fill play-btn"></i></button>
                          </div></div></td>
        <td>{{album}}</td>
        <td>{{date}}</td>
        <td><div class="like-container"><i class="bi bi-heart"></i><div class="current-likes">{{likes}}</div></div></td>
        <td>{{duration}}</td>
        <td>
        <div class="control-wrapper">
          <i class="bi bi-three-dots control-playlist-trigger"></i>
          <div class="control-playlist-container">
<!--            <div class="control-playlist-item add-to-queue table-view">Add to queue</div>-->
            <div class="control-playlist-item go-to-artist table-view" data-artist-id="{{artistId}}" data-artist-name="{{artist}}">Go to artist</div>
            <!--<div class="control-playlist-item go-to-album table-view">Go to album</div>-->
            <div class="control-playlist-item add-to-a-playlist table-view">Add to playlist</div>
<!--            <div class="control-playlist-item share-song table-view">Share</div>-->
          </div>
        </div>
        </td>
      </tr>`;
      songItem = templateReplace(songItem, mapObj);
      tableHTML += songItem;
    });

    if(albumSongs.length > 1){
      let albumSongsJSON = JSON.stringify(albumSongList);
      console.log(albumSongsJSON);
      contentItem = `<div class="content">
      <div class="content-wrapper" data-artist-name="{{artist}}" data-album-name="{{album}}" data-music-state="paused" data-album-id="{{albumId}}" data-album-songs='${albumSongsJSON}'>
        <div class="content-img audio-cover">
          <img src="{{cover}}">
          <div class="buttons">
            <button><i class="bi bi-play-fill play-btn"></i></button>
          </div>
        </div>
        <div class="song-info-title">{{album}}</div>
        <div class="song-info-artist">{{artist}}</div>
      </div>
    </div>`;
      contentItem = templateReplace(contentItem, mapObj);
      albumList += contentItem;
    }
    else {
      contentItem = `<div class="content">
      <div class="content-wrapper" data-artist-name="{{artist}}" data-song-name="{{album}}" data-music-state="paused" data-song-id="{{albumId}}" data-audio-path="${audioPath}">
        <div class="content-img audio-cover">
          <img src="{{cover}}">
          <div class="buttons">
            <button><i class="bi bi-play-fill play-btn"></i></button>
          </div>
        </div>
        <div class="song-info-title">{{album}}</div>
        <div class="song-info-artist">{{artist}}</div>
      </div>
    </div>`;
      contentItem = templateReplace(contentItem, mapObj);
      singleList += contentItem;
    }

  });

  htmltemplate = templateReplace(htmltemplate, {
    "{{artistAlbums}}": albumList,
  });

  htmltemplate = templateReplace(htmltemplate, {
    "{{artistSingles}}": singleList,
  });

  htmltemplate = templateReplace(htmltemplate, { "{{allPlaysCount}}": (totalPlays * 107407).toLocaleString() });
  htmltemplate = templateReplace(htmltemplate, { "{{artistID}}": artistID });
  htmltemplate = templateReplace(htmltemplate, { "{{artist}}": artistName });
  htmltemplate = templateReplace(htmltemplate, {
    "{{table-content}}": tableHTML,
  });
  // htmltemplate = templateReplace(htmltemplate, { "{{artist-cover}}": artistName.replace(/ /g, '').toLowerCase()});

  // * Add Elements To Main View & Update Body View
  if (bodyContainer.classList.contains("index-home")) {
    bodyContainer.classList.remove("index-home");
    bodyContainer.classList.add("index-view");
  }
  viewContainer.classList.add("show-animation");
  viewContainer.insertAdjacentHTML("beforeend", htmltemplate);
}

export function updateViewPlaylist(data, viewContainer, playlistID, viewName, bodyContainer, htmltemplate) {
  viewContainer.innerHTML = "";

  // * Fill the Table with Songs
  let tableHTML = "";
  data.forEach((entry) => {
    const mapObj = {
      "{{likes}}": entry.Likes ?? entry.Total_likes,
      "{{dislikes}}": entry.Dislikes ?? "",
      "{{songID}}": entry.Song_id,
      "{{song}}": entry.Title,
      "{{artistId}}": entry.Artist_id,
      "{{audio}}": entry.SongPath,
      "{{cover}}": entry.CoverPath,
      "{{date}}": dateParse(entry.UploadedDate ?? entry.Uploaded_date),
      "{{album}}": entry.Album,
      "{{albumId}}": entry.Album_id,
      "{{totalPlays}}": entry.Total_plays,
      "{{artist}}": entry.Artist ?? entry.Artist_name,
      "{{duration}}": entry.Duration,
      "{{playlistID}}": playlistID,
    };
    let songItem = `<tr class="table-song-item" data-song-id="{{songID}}" data-playlist-id="{{playlistID}}">
    <td><div class="playlist-item-flex content-wrapper" data-audio-path="{{audio}}" data-music-state="paused" data-song-id="{{songID}}"><div class="playlist-img-contain audio-cover"><img src="{{cover}}"></div><div class="song-info-item"><div class="song-info-title">{{song}}</div><div class="song-info-artist">{{artist}}</div></div><div class="buttons playlist">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div></div></td>
    <td>{{album}}</td>
    <td>{{date}}</td>
    <td><div class="like-container"><i class="bi bi-heart"></i><div class="current-likes">{{likes}}</div></div></td>
    <td>{{duration}}</td>
    <td>
    <div class="control-wrapper">
      <i class="bi bi-three-dots control-playlist-trigger"></i>
      <div class="control-playlist-container">
<!--        <div class="control-playlist-item add-to-queue table-view">Add to queue</div>-->
        <div class="control-playlist-item go-to-artist table-view" data-artist-id="{{artistId}}" data-artist-name="{{artist}}">Go to artist</div>
            <!--<div class="control-playlist-item go-to-album table-view">Go to album</div>-->
        <div class="control-playlist-item remove-from-playlist table-view">Remove from playlist</div>
        <div class="control-playlist-item add-to-a-playlist table-view">Add to playlist</div>
<!--        <div class="control-playlist-item share-song table-view">Share</div>-->
      </div>
    </div>
    </td>
  </tr>`;
    songItem = templateReplace(songItem, mapObj);
    tableHTML += songItem;
  });

  htmltemplate = templateReplace(htmltemplate, {
    "{{playlistID}}": playlistID,
  });
  htmltemplate = templateReplace(htmltemplate, { "{{viewName}}": viewName });
  htmltemplate = templateReplace(htmltemplate, {
    "{{table-content}}": tableHTML,
  });
  htmltemplate = templateReplace(htmltemplate, {
    "{{song-count}}": songCount(data),
  });
  viewContainer.insertAdjacentHTML("beforeend", htmltemplate);

  // * Add Elements To Main View & Update Body View
  if (bodyContainer.classList.contains("index-home")) {
    bodyContainer.classList.remove("index-home");
    bodyContainer.classList.add("index-view");
  }
  viewContainer.classList.add("show-animation");
}

export function updateViewDiscover(data, viewContainer, viewName, bodyContainer, htmltemplate) {
  viewContainer.innerHTML = "";

  console.log(data);
  // Fill the Table with Songs
  let tableHTML = "";
  data.forEach((entry) => {
    const mapObj = {
      "{{likes}}": entry.Likes ?? entry.Total_likes,
      "{{dislikes}}": entry.Dislikes ?? "",
      "{{songID}}": entry.Song_id,
      "{{song}}": entry.Title,
      "{{artistId}}": entry.Artist_id,
      "{{audio}}": entry.SongPath,
      "{{cover}}": entry.CoverPath,
      "{{date}}": dateParse(entry.UploadedDate ?? entry.Uploaded_date),
      "{{album}}": entry.Album,
      "{{albumId}}": entry.Album_id,
      "{{totalPlays}}": entry.Total_plays,
      "{{artist}}": entry.Artist ?? entry.Artist_name,
      "{{duration}}": entry.Duration,
    };
    let songItem = `<tr class="table-song-item" data-song-id="{{songID}}">
    <td><div class="playlist-item-flex content-wrapper" data-audio-path="{{audio}}" data-music-state="paused" data-song-id="{{songID}}"><div class="playlist-img-contain audio-cover"><img src="{{cover}}"></div><div class="song-info-item"><div class="song-info-title">{{song}}</div><div class="song-info-artist">{{artist}}</div></div><div class="buttons playlist">
                        <button><i class="bi bi-play-fill play-btn"></i></button>
                      </div></div></td>
    <td>{{album}}</td>
    <td>{{date}}</td>
    <td class="song-play-count">{{totalPlays}}</td>
    <td>{{duration}}</td>
    <td>
    <div class="control-wrapper">
      <i class="bi bi-three-dots control-playlist-trigger"></i>
      <div class="control-playlist-container">
<!--        <div class="control-playlist-item add-to-queue table-view">Add to queue</div>-->
        <div class="control-playlist-item go-to-artist table-view" data-artist-id="{{artistId}}" data-artist-name="{{artist}}">Go to artist</div>
            <!--<div class="control-playlist-item go-to-album table-view">Go to album</div>-->
        <div class="control-playlist-item add-to-a-playlist table-view">Add to playlist</div>
<!--        <div class="control-playlist-item share-song table-view">Share</div>-->
      </div>
    </div>
    </td>
  </tr>`;
    songItem = templateReplace(songItem, mapObj);1
    tableHTML += songItem;
  });

  htmltemplate = templateReplace(htmltemplate, { "{{viewName}}": viewName });
  htmltemplate = templateReplace(htmltemplate, {
    "{{table-content}}": tableHTML,
  });
  htmltemplate = templateReplace(htmltemplate, {
    "{{song-count}}": songCount(data),
  });
  viewContainer.insertAdjacentHTML("beforeend", htmltemplate);

  // Add Elements To Main View & Update Body View
  if (bodyContainer.classList.contains("index-home")) {
    bodyContainer.classList.remove("index-home");
    bodyContainer.classList.add("index-view");
  }

  viewContainer.classList.add("show-animation");
}

export function updateViewAdminControl(data, viewContainer, viewName, bodyContainer, htmltemplate) {
  viewContainer.innerHTML = "";

  console.log(data);
  // Fill the Table with Users
  let tableHTML = "";
  data.forEach((entry) => {
    if (entry && entry.Admin_level) {
      if (entry.Admin_level == 1) {
        entry.Admin_level = 'User'
      } else if (entry.Admin_level == 2) {
        entry.Admin_level = 'Artist'
      }
    }
    const mapObj = {
      "{{song}}": entry.Title ?? entry.Song_title ?? "",
      "{{album}}": entry.Album ?? entry.Album_name ?? "",
      "{{artist}}": entry.Artist ?? entry.Artist_name ?? "",
      "{{date}}": dateParse(entry.UploadedDate ?? entry.Uploaded_date ?? ""),
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
      "{{joinedDate}}": dateParse(entry.JoinedDate ?? entry.Join_date) ?? "",
      "{{playlistCount}}": entry.Playlist_count ?? "",
      "{{userid}}": entry.User_id,
      "{{likedSongsCount}}": entry.Liked_songs_count ?? "",
      "{{commonArtist}}": entry.Common_artist ?? "",
      "{{avgPlays}}": entry.Avg_Plays ?? "",
      "{{mostLikedSongs}}": entry.Most_liked_song ?? "",
      "{{numSongs}}": entry.Num_songs ?? "",
      "{{numAlbums}}": entry.Num_Albums ?? "",
      "{{name}}": entry.Name ?? "",
    }
    let songItem = `<tr class="table-item">
                        <td><div class="table-item-flex">{{userid}}</div></td>  
                         <td>{{username}}</td>                     
                        <td>{{firstname}}</td>
                        <td>{{lastname}}</td>
                        <td>{{adminlevel}}</td>
                        <td>{{joinedDate}}</td>
                        <td>{{playlistCount}}</td>
                        <td>{{likedSongsCount}}</td>  
                        <td>{{commonArtist}}</td>
                        <td><i class="bi bi-trash-fill admin-remove-user" data-user-id="{{userid}}"></i></td>
                      </tr>`;
    songItem = templateReplace(songItem, mapObj);
    tableHTML += songItem;
  });

  htmltemplate = templateReplace(htmltemplate, { "{{viewName}}": viewName });
  htmltemplate = templateReplace(htmltemplate, {
    "{{table-content}}": tableHTML,
  });
  htmltemplate = templateReplace(htmltemplate, {
    "{{song-count}}": songCount(data),
  });
  viewContainer.insertAdjacentHTML("beforeend", htmltemplate);

  // Add Elements To Main View & Update Body View
  if (bodyContainer.classList.contains("index-home")) {
    bodyContainer.classList.remove("index-home");
    bodyContainer.classList.add("index-view");
  }
  viewContainer.classList.add("show-animation");
}

export function updateViewHomeControl(data, viewContainer, viewName, bodyContainer, htmltemplate) {
  viewContainer.innerHTML = "";

  console.log(data);
  // Fill the Table with Users
  let homeSongs = "";
  data.forEach((entry) => {
    const mapObj = {
      "{{song}}": entry.Title,
      "{{album}}": entry.Album,
      "{{artist}}": entry.Artist ?? entry.Artist_name,
      "{{date}}": dateParse(entry.UploadedDate ?? entry.Uploaded_date),
      "{{dislikes}}": entry.Dislikes ?? "",
      "{{likes}}": entry.Likes ?? entry.Total_likes,
      "{{cover}}": entry.CoverPath,
      "{{audio}}": entry.SongPath,
      "{{songID}}": entry.Song_id,
    };
    let songItem = `<div class="content">
                  <div class="content-wrapper" data-audio-path="{{audio}}" data-music-state="paused" data-song-id="{{songID}}">
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

  htmltemplate = templateReplace(htmltemplate, {
    "{{song-content}}": homeSongs,
  });
  viewContainer.insertAdjacentHTML("beforeend", htmltemplate);

  // Add Elements To Main View & Update Body View
  if (bodyContainer.classList.contains("index-view")) {
    bodyContainer.classList.remove("index-view");
    bodyContainer.classList.add("index-home");
  }
  viewContainer.classList.add("show-animation");
}
