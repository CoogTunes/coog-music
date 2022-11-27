package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CoogTunes/coog-music/internal/forms"
	"github.com/CoogTunes/coog-music/internal/render"
	"golang.org/x/crypto/bcrypt"

	"github.com/CoogTunes/coog-music/internal/config"
	"github.com/CoogTunes/coog-music/internal/driver"
	"github.com/CoogTunes/coog-music/internal/models"
	"github.com/go-chi/chi/v5"

	"github.com/CoogTunes/coog-music/internal/repository"
	"github.com/CoogTunes/coog-music/internal/repository/dbrepo"

	"github.com/tcolgate/mp3"
)

// Repo the repository used by the handlers
var Repo *Repository

var UserCache models.Users

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// LOGIN/SIGNUP

func (m *Repository) GetLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	loginOrRegister := r.Form.Get("submit_button")
	fmt.Println(loginOrRegister)
	if loginOrRegister == "register" {
		fmt.Println("Passes to function PostRegistration ")
		m.PostRegistration(w, r)
		return
	}
	//_ = m.App.Session.RenewToken(r.Context())

	email := r.Form.Get("email")
	pwd := r.Form.Get("password")

	userInfo, err := m.DB.Authenticate(email, pwd)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	UserCache = userInfo
	UserCache.Concatenated_name = getArtistName(userInfo.First_name, userInfo.Last_name)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) PostRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	pwd := []byte(r.Form.Get("password"))
	hashedPwd, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return
	}

	adminLevel := r.Form.Get("user-level")
	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	lvl := 0
	if adminLevel == "user" {
		lvl = 1
	} else if adminLevel == "artist" {
		lvl = 2
	}

	registrationModel := models.Users{
		First_name:  firstName,
		Last_name:   lastName,
		Username:    r.Form.Get("email"),
		Password:    string(hashedPwd),
		Admin_level: lvl,
	}
	users, err := m.DB.AddUser(registrationModel)
	if false {
		log.Println(users)
	}
	if err != nil {
		log.Fatal(err)
	}

	UserCache = users

	UserCache.Concatenated_name = getArtistName(users.First_name, users.Last_name)

	if lvl == 2 {
		m.AddArtist(firstName, lastName)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) PageLoad(w http.ResponseWriter, r *http.Request) {
	idx := r.URL.Query().Get("index")
	if idx == "discover" {
		m.GetTopSongs(w, r)
		return
	} else if idx == "admin" {
		m.GetTopUserReport(w, r)
		return
	} else if idx == "home" {
		m.GetTopSongs(w, r)
		return
	}

}

func (m *Repository) GetTopUserReport(w http.ResponseWriter, r *http.Request) {
	userReport, err := m.DB.GetInitialUsersReport()
	if err != nil {
		log.Println(err)
	}
	if len(userReport) == 0 {
		uReport := []models.UserReport{}
		returnAsJSON(uReport, w, err)
		return
	}
	returnAsJSON(userReport, w, err)
}

// END LOGIN/SIGNUP--------------------------------------------------------------------------------------

// ADD ARTIST -------------------------------------------------------------------------------------------

func (m *Repository) AddArtist(firstName string, lastName string) {
	artistName := getArtistName(firstName, lastName)

	artistInfo := models.Artist{
		Name:      artistName,
		Artist_id: UserCache.User_id,
		Location:  "no_location",
	}

	err := m.DB.AddArtist(artistInfo)
	if err != nil {
		log.Println("Cannot add artist information")
		return
	}
}

// END ADD ARTIST ----------------------------------------------------------------------------

// HOME PAGE ---------------------------------------------------------------------------------

func (m *Repository) GetHome(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "coogtunes.page.gohtml", &models.TemplateData{
		Form:     forms.New(nil),
		UserData: UserCache,
	})

	m.GetPlaylistsByID(w, r)
}
func (m *Repository) GetTopSongs(w http.ResponseWriter, r *http.Request) {
	topSongs, err := m.DB.GetTopSongs()
	if err != nil {
		log.Println("Cannot get the top 14 songs")
	}
	if len(topSongs) == 0 {
		tSongs := []models.Song{}
		returnAsJSON(tSongs, w, err)
		return
	}
	returnAsJSON(topSongs, w, err)
}

func (m *Repository) GetPlaylistsByID(w http.ResponseWriter, r *http.Request) {
	playlists, err := m.DB.GetPlaylists(UserCache.User_id)
	if err != nil {
		log.Println("Cannot get the top 14 songs")
	}
	if len(playlists) == 0 {
		emptyPlaylist := models.Playlist{}
		returnAsJSON(emptyPlaylist, w, err)
		return
	}
	returnAsJSON(playlists, w, err)
}

// LOGOUT

func (m *Repository) LogOut(w http.ResponseWriter, r *http.Request) {
	noUserData := models.Users{}
	UserCache = noUserData
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// END OF HOME PAGE ---------------------------------------------------------------------------------

// PROFILE PAGE ---------------------------------------------------------------------------------

func (m *Repository) GetProfile(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "profilepage.page.gohtml", &models.TemplateData{
		Form:     forms.New(nil),
		UserData: UserCache,
	})
}

//  END OF PROFILE PAGE ---------------------------------------------------------------------------------

// UPLOAD MUSIC ---------------------------------------------------------------------------------

func (m *Repository) UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Passing through upload file handler")
	err := r.ParseMultipartForm(32 << 20) // 32mb
	if err != nil {
		log.Fatal("Cannot parse upload Files")
	}

	songOrAlbum := r.Form.Get("uploadType")
	if err != nil {
		fmt.Println("cannot parse the image file")
	}

	if songOrAlbum == "song" {
		fmt.Println("Passing through the upload song handler")
		m.UploadSong(w, r)
		return
	}

	m.UploadAlbum(w, r)
}

func (m *Repository) UploadSong(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatal("Could not parse multipart forms")
	}

	artistName := concatenateName(getArtistName(UserCache.First_name, UserCache.Last_name))
	songName := r.Form.Get("music_name")
	date := r.Form.Get("released_date")
	fmt.Println("Passes through the songName")
	coverFile, fhCover, err := r.FormFile("music_cover")
	if err != nil {
		log.Println("Cannot read cover_file")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Passes through coverFile")

	coverPath := "./static/media/artist/" + artistName
	err = os.MkdirAll(coverPath, os.ModePerm)
	if os.IsExist(err) {
		log.Println("Cover jpeg is already in directory")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dst, err := os.Create(fmt.Sprintf(coverPath+"/"+"%s", fhCover.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(dst, coverFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fullCoverPath := coverPath + "/" + fhCover.Filename
	coverFile.Close()
	dst.Close()

	songFile, fhSong, err := r.FormFile("music_audio")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	songPath := "./static/media/artist/" + artistName
	err = os.MkdirAll(songPath, os.ModePerm)
	if os.IsExist(err) {
		log.Println("Song audio is already in directory")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dst2, err := os.Create(fmt.Sprintf(songPath+"/"+"%s", fhSong.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(dst2, songFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fullSongPath := songPath + "/" + fhSong.Filename
	coverFile.Close()
	dst.Close()
	duration := getMp3Duration(fullSongPath)

	songInfo := models.Song{
		Title:         songName,
		Artist_id:     UserCache.User_id,
		CoverPath:     fullCoverPath,
		SongPath:      fullSongPath,
		Duration:      duration,
		Uploaded_date: date,
	}
	fmt.Println(songInfo)
	err = m.DB.AddSong(songInfo)
	if err != nil {
		log.Println("Cannot add song to the database")
	}

	fmt.Fprintf(w, "Upload successful")

}

func (m *Repository) UploadAlbum(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200 << 20)
	if err != nil {
		log.Fatal("Could not parse multipart forms")
	}
	artistName := concatenateName(getArtistName(UserCache.First_name, UserCache.Last_name))
	albumName := r.Form.Get("music_name")
	date := r.Form.Get("released_date")
	pathAlbumName := concatenateName(albumName)
	coverFile, fhCover, err := r.FormFile("music_cover")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	coverPath := "./static/media/artist/" + artistName + "/" + pathAlbumName
	err = os.MkdirAll(coverPath, os.ModePerm)
	if os.IsExist(err) {
		log.Println("Cover jpeg is already in directory")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dst, err := os.Create(fmt.Sprintf(coverPath+"/"+"%s", fhCover.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(dst, coverFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fullCoverPath := coverPath + "/" + fhCover.Filename
	coverFile.Close()
	dst.Close()

	albumInfo := models.Album{
		Name:      albumName,
		Artist_id: UserCache.User_id,
	}

	albumDBInfo, err := m.DB.AddAlbum(albumInfo)

	if err != nil {
		log.Println("Cannot add album")
		return
	}
	files := r.MultipartForm.File["music_audio"]
	for _, fileHeader := range files {
		fileSize := fileHeader.Size
		fmt.Println(fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("Could not open the file")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		buff := make([]byte, fileSize)
		_, err = file.Read(buff)
		if err != nil {
			fmt.Println("Could not read the file")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "audio/mpeg" {
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG image or MP3 file", http.StatusBadRequest)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filePath := "./static/media/artist/" + artistName + "/" + pathAlbumName
		err = os.MkdirAll(filePath, os.ModePerm)
		if os.IsExist(err) {
			fmt.Println("Song already exists!")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		f, err := os.Create(fmt.Sprintf(filePath+"/"+"%s", fileHeader.Filename))
		if err != nil {
			fmt.Println("Could not create the named file")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			fmt.Println("Could not copy the uploaded files")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		title := strings.ReplaceAll(fileHeader.Filename, filepath.Ext(fileHeader.Filename), "")
		newTitle := splitName(title)
		songPath := filePath + "/" + fileHeader.Filename
		duration := getMp3Duration(songPath)
		songInfo := models.Song{
			Title:         newTitle,
			Album:         albumName,
			SongPath:      songPath,
			CoverPath:     fullCoverPath,
			Artist_id:     UserCache.User_id,
			Album_id:      albumDBInfo.Album_id,
			Duration:      duration,
			Uploaded_date: date,
		}
		err = m.DB.AddSongForAlbum(songInfo)
		if err != nil {
			log.Println("Cannot add song")
			return
		}
	}

}

// END OF UPLOAD MUSIC ---------------------------------------------------------------------------------

// SEARCH SECTION --------------------------------------------------------------------------------------
func (m *Repository) Search(w http.ResponseWriter, r *http.Request) (string, string) {
	target := r.URL.Query().Get("strTarget")
	filter := r.URL.Query().Get("filters")
	decodedValue, err := url.QueryUnescape(target)
	if err != nil {
		log.Print("Could not decode the parameter")
	}
	return decodedValue, filter
}

// PLAYLIST SECTION ---------------------------------------------------------------------------------
// TODO: Return empty json when getting an error on
func (m *Repository) PlaylistSearch(w http.ResponseWriter, r *http.Request) {
	decodedValue, filter := m.Search(w, r)
	if filter == "song" {
		songInfo, err := m.DB.GetSongsByName(decodedValue)
		if err != nil {
			returnAsJSON(songInfo, w, err)
			log.Println("Cannot get songs!")
		}
		returnAsJSON(songInfo, w, err)
		return
	} else if filter == "album" {
		albumInfo, err := m.DB.GetSongsFromAlbum(decodedValue)
		if err != nil {
			log.Println("Cannot get songs!")
		}

		returnAsJSON(albumInfo, w, err)
		return
	} else if filter == "artist" {
		albumInfo, err := m.DB.GetSongsFromArtist(decodedValue)
		if err != nil {
			log.Println("Cannot get songs!")
		}

		returnAsJSON(albumInfo, w, err)
		return
	}
}

type PlayListJson struct {
	PlayListName string   `json:"playlistName"`
	PlayList     []string `json:"playListItems"`
}

func (m *Repository) InsertPlaylist(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var playlistInfo PlayListJson
	err := decoder.Decode(&playlistInfo)

	log.Println(playlistInfo.PlayListName)

	if err != nil {
		log.Println("Cannot decode the json")
	}
	playlist := models.Playlist{
		User_id: UserCache.User_id,
		Name:    playlistInfo.PlayListName,
	}

	plylist, err := m.DB.AddPlaylist(playlist)
	if err != nil {
		log.Println("Cannot add playlist")
		return
	}
	for _, strNum := range playlistInfo.PlayList {
		songID, err := strconv.Atoi(strNum)
		if err != nil {
			log.Println("Cannot convert string to num")
		}
		err = m.DB.AddPlaylistSong(songID, plylist.Playlist_id)
		if err != nil {
			log.Println("Cannot add playlist")
		}
	}

	returnAsJSON(plylist, w, err)

}

//type GetSongsFromPlaylist struct {
//	PlayListID string `json:"playlist-id"`
//}

func (m *Repository) GetPlaylistSongs(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	playlistID, err := strconv.Atoi(id)
	fmt.Println(playlistID)
	if err != nil {
		log.Println("Cannot convert string to int")
	}
	displaySongs, err := m.DB.GetSongsFromPlaylist(playlistID)

	if len(displaySongs) == 0 {
		disSong := []models.Song{}
		returnAsJSON(disSong, w, err)
		return
	}
	returnAsJSON(displaySongs, w, err)
}

func (m *Repository) GetPlaylists(w http.ResponseWriter, r *http.Request) {
	playlists, err := m.DB.GetPlaylists(UserCache.User_id)

	returnAsJSON(playlists, w, err)
	return
}

// END PLAYLIST SECTION --------------------------------------------------------------------------------

// USERS ------------------------------------------------------------------------

func (m *Repository) AddUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get fields
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	first_name := r.Form.Get("first_name")
	last_name := r.Form.Get("last_name")
	stringAdmin := r.Form.Get("admin")

	admin, err := strconv.Atoi(stringAdmin)
	if err != nil {
		log.Println(err)
	}

	newUser := models.Users{Username: username, Password: password, First_name: first_name, Last_name: last_name, Admin_level: admin}

	addedUser, err := m.DB.AddUser(newUser)

	returnAsJSON(addedUser, w, err)
}

func (m *Repository) GetUser(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	user, err := m.DB.GetUser(id)

	returnAsJSON(user, w, err)
}

// END USERS -------------------------------------------------------------------

func (m *Repository) GetArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := m.DB.GetArtists()

	returnAsJSON(artists, w, err)
}

//
//func (m *Repository) GetArtistsAndSongs(w http.ResponseWriter, r *http.Request) {
//	artists, err := m.DB.GetArtistsAndSongs()
//	returnAsJSON(artists, w, err)
//}

func (m *Repository) UpdateArtist(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get fields
	name := r.Form.Get("name")
	// artist_idString := r.Form.Get("artist_id")
	location := r.Form.Get("location")
	// join_date := r.Form.Get("join_date")

	// artist_id, err := strconv.Atoi(artist_idString)
	// if err != nil {
	// 	log.Println(err)
	// }

	artistToUpdate := models.Artist{Name: name, Location: location}

	addedUser, err := m.DB.UpdateArtist(artistToUpdate)

	returnAsJSON(addedUser, w, err)
}

type UpdateCount struct {
	SongId string `json:"songID"`
}

func (m *Repository) UpdateSongCount(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var updateC UpdateCount
	err := decoder.Decode(&updateC)
	if err != nil {
		log.Println("Cannot decode the json object")
	}
	id, err := strconv.Atoi(updateC.SongId)
	if err != nil {
		log.Println(err)
	}
	song, err := m.DB.UpdateSongCount(id)
	returnAsJSON(song, w, err)
}

type SongToPlaylist struct {
	PlaylistID string `json:"playlistID"`
	SongID     string `json:"songID"`
}

func (m *Repository) AddSongToPlaylist(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var songToPlaylist SongToPlaylist
	err := decoder.Decode(&songToPlaylist)
	if err != nil {
		log.Println("Cannot decode SongToPlaylist json")
	}
	sID, err := strconv.Atoi(songToPlaylist.SongID)
	if err != nil {
		log.Println("Cannot change song ID to int")
	}
	pID, err := strconv.Atoi(songToPlaylist.PlaylistID)
	if err != nil {
		log.Println("Cannot change playlist ID to int")
	}
	err = m.DB.AddPlaylistSong(sID, pID)
	if err != nil {
		log.Println("Cannot add song to playlist")
	}

	var sMessage SuccessMessage

	sMessage.Message = "success"

	returnAsJSON(sMessage, w, err)

}

func (m *Repository) GetArtistInfo(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Query().Get("artistID")
	artistName := r.URL.Query().Get("artistName")
	decodedValue, err := url.QueryUnescape(artistName)
	if err != nil {
		log.Println("Cannot get artistName")
	}
	aID, err := strconv.Atoi(artistID)
	if err != nil {
		log.Println("Cannot decode the artist id")
	}

	artistAlbums, err := m.DB.GetSongsFromArtistByID(decodedValue, aID)
	if err != nil {
		log.Println("Cannot get the artist albums")
	}

	returnAsJSON(artistAlbums, w, err)
}

func (m *Repository) GetArtistsSongsAndAlbums(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("artistName")
	artistSongsAndAlbums, err := m.DB.GetSongsFromArtist(artistName)
	if err != nil {
		log.Println("Cannot get songs and albums from artist")
	}
	returnAsJSON(artistSongsAndAlbums, w, err)
}

func (m *Repository) AddSongToAlbum(w http.ResponseWriter, r *http.Request) {
	return
}

func (m *Repository) GetAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := m.DB.GetAlbums()

	returnAsJSON(albums, w, err)
}

func (m *Repository) GetAlbumInfo(w http.ResponseWriter, r *http.Request) {
	albumID := r.URL.Query().Get("albumID")
	albumName := r.URL.Query().Get("albumName")
	decodedValue, err := url.QueryUnescape(albumName)
	if err != nil {
		log.Println("Cannot get artistName")
	}
	aID, err := strconv.Atoi(albumID)
	if err != nil {
		log.Println("Cannot decode the album id")
	}
	albumSongs, err := m.DB.GetSongsFromAlbumByID(decodedValue, aID)
	if err != nil {
		log.Println("Cannot get songs from album")
	}
	returnAsJSON(albumSongs, w, err)
}

func (m *Repository) UpdateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get fields
	userIdString := r.Form.Get("user_id")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	first_name := r.Form.Get("first_name")
	last_name := r.Form.Get("last_name")
	stringAdmin := r.Form.Get("admin_level")

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		log.Println(err)
	}

	admin, err := strconv.Atoi(stringAdmin)
	if err != nil {
		log.Println(err)
	}

	newUser := models.Users{User_id: userId, Username: username, Password: password, First_name: first_name, Last_name: last_name, Admin_level: admin}

	addedUser, err := m.DB.UpdateUser(newUser)

	returnAsJSON(addedUser, w, err)
}

func (m *Repository) Follow(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get fields
	artist_idString := r.Form.Get("artist_id")
	user_idString := r.Form.Get("user_id")

	artistId, err := strconv.Atoi(artist_idString)
	if err != nil {
		log.Println(err)
	}
	userId, err := strconv.Atoi(user_idString)
	if err != nil {
		log.Println(err)
	}

	follower, err := m.DB.Follow(artistId, userId)
	returnAsJSON(follower, w, err)
}

func (m *Repository) GetSongsForLikePage(w http.ResponseWriter, r *http.Request) {
	likedSongs, err := m.DB.GetSongsForLikePage(UserCache.User_id)
	// likedSongs, err := m.DB.GetSongsForLikePage(1)
	returnAsJSON(likedSongs, w, err)
}

type AddLike struct {
	Check  string `json:"check"`
	SongId string `json:"songID"`
}

func (m *Repository) AddOrUpdateLikeValue(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var likeStruct AddLike
	err := decoder.Decode(&likeStruct)
	if err != nil {
		log.Println("Cannot decode the json")
	}

	isLike, err := strconv.ParseBool(likeStruct.Check)
	if err != nil {
		log.Println(err)
	}
	song_id, err := strconv.Atoi(likeStruct.SongId)

	if err != nil {
		log.Println(err)
	}
	err2 := m.DB.AddOrUpdateLikeValue(isLike, song_id, UserCache.User_id)
	if err2 != nil {
		log.Println(err)
	}

	songList, err3 := m.DB.SendUpdatedLikeValue(song_id)
	if err3 != nil {
		log.Println("Cannot get the updated song like")
	}
	returnAsJSON(songList, w, err3)
}

func (m *Repository) UpdateMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := m.DB.UpdateMessages(UserCache.User_id)
	returnAsJSON(messages, w, err)
}

func (m *Repository) Filter(w http.ResponseWriter, r *http.Request) {
	likes := r.URL.Query().Get("likes")
	fmt.Println(likes)
	plays := r.URL.Query().Get("plays")
	fmt.Println(plays)
	artists := r.URL.Query().Get("artists")
	log.Println(artists)
	users := r.URL.Query().Get("users")
	fmt.Println("in filter")
	if likes == "true" && plays != "true" {
		m.GetLikesReport(w, r)
		return
	} else if plays == "true" && likes != "true" {
		m.GetSongReport(w, r)
		return
	} else if users == "true" && artists != "true" {
		fmt.Println("in filter users")
		m.GetUserReport(w, r)
	} else if artists == "true" && users != "true" {
		m.GetArtistReport(w, r)
	}
}

// REPORTS
func (m *Repository) GetLikesReport(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	// get fields

	minLikes, err := strconv.Atoi(r.URL.Query().Get("min"))
	if err != nil {
		log.Println(err)
	}
	maxLikes, err := strconv.Atoi(r.Form.Get("max"))
	if err != nil {
		log.Println(err)
	}
	minDate := r.URL.Query().Get("start")
	maxDate := r.URL.Query().Get("end")
	//minDislikes, err := strconv.Atoi(r.Form.Get("min_dislikes"))
	//if err != nil {
	//	log.Println(err)
	//}
	//maxDislikes, err := strconv.Atoi(r.Form.Get("max_dislikes"))
	//if err != nil {
	//	log.Println(err)
	//}

	likesReport, err := m.DB.GetLikesReport(minLikes, maxLikes, minDate, maxDate)
	if len(likesReport) == 0 {
		lReport := []models.Song{}
		returnAsJSON(lReport, w, err)
		return
	}
	returnAsJSON(likesReport, w, err)
}

func (m *Repository) GetUserReport(w http.ResponseWriter, r *http.Request) {

	minDate := r.URL.Query().Get("start")
	maxDate := r.URL.Query().Get("end")
	usersReport, err := m.DB.GetUsersReport(minDate, maxDate)
	if len(usersReport) == 0 {
		lReport := []models.UserReport{}
		returnAsJSON(lReport, w, err)
		return
	}
	returnAsJSON(usersReport, w, err)

}

func (m *Repository) GetArtistReport(w http.ResponseWriter, r *http.Request) {
	minDate := r.URL.Query().Get("start")
	maxDate := r.URL.Query().Get("end")
	artistReport, err := m.DB.GetArtistReport(minDate, maxDate)
	if len(artistReport) == 0 {
		lReport := []models.ArtistReport{}
		returnAsJSON(lReport, w, err)
		return
	}
	returnAsJSON(artistReport, w, err)
}

func (m *Repository) GetSongReport(w http.ResponseWriter, r *http.Request) {

	// get fields

	min_plays, err := strconv.Atoi(r.URL.Query().Get("min"))
	if err != nil {
		log.Println(err)
	}
	max_plays, err := strconv.Atoi(r.URL.Query().Get("max"))
	if err != nil {
		log.Println(err)
	}
	minDate := r.URL.Query().Get("start")
	maxDate := r.URL.Query().Get("end")

	songReport, err := m.DB.GetSongReport(minDate, maxDate, min_plays, max_plays)
	if len(songReport) == 0 {
		sReport := []models.Song{}
		returnAsJSON(sReport, w, err)
		return
	}
	returnAsJSON(songReport, w, err)
}

type DeletingSong struct {
	SongID     string `json:"songID"`
	PlaylistID string `json:"playlistID"`
}

type SuccessMessage struct {
	Message string
}

func (m *Repository) DeleteSong(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var songDelete DeletingSong
	err := decoder.Decode(&songDelete)
	if err != nil {
		log.Println("Cannot decode DeletingSong json")
	}
	sID, err := strconv.Atoi(songDelete.SongID)
	if err != nil {
		log.Println("Cannot convert strig to int for songID")
	}
	pID, err := strconv.Atoi(songDelete.PlaylistID)
	if err != nil {
		log.Println("Cannot convert strig to int for playlistID")
	}

	err = m.DB.RemoveSongFromPlaylist(sID, pID)
	if err != nil {
		log.Println("Cannot remove song from playlist")
	}
	var sMessage SuccessMessage

	sMessage.Message = "success"

	returnAsJSON(sMessage, w, err)
}

type DeletingPlaylist struct {
	Playlist string `json:"playlistID"`
}

func (m *Repository) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var playlistDelete DeletingPlaylist
	err := decoder.Decode(&playlistDelete)
	if err != nil {
		log.Println("Cannot decode DeletingPlaylist json")
	}
	pID, err := strconv.Atoi(playlistDelete.Playlist)
	if err != nil {
		log.Println("Cannot convert strig to int for playlistID")
	}

	err = m.DB.RemovePlaylist(pID)
	if err != nil {
		log.Println("Cannot delete playlist")
	}
	var sMessage SuccessMessage

	sMessage.Message = "success"

	returnAsJSON(sMessage, w, err)
}

type DeletingUser struct {
	User string `json:"userID"`
}

func (m *Repository) DeleteUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var userDelete DeletingUser
	err := decoder.Decode(&userDelete)
	if err != nil {
		log.Println("Cannot decode DeletingUser json")
	}
	uID, err := strconv.Atoi(userDelete.User)
	if err != nil {
		log.Println("Cannot convert strig to int for userID")
	}
	err = m.DB.RemoveUser(uID)
	if err != nil {
		log.Println("Cannot delete user")
	}
	var sMessage SuccessMessage

	sMessage.Message = "success"

	returnAsJSON(sMessage, w, err)
}

// HELPER FUNCTIONS
// i is the models.XYZ property
func returnAsJSON(i interface{}, w http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
	}
	j, _ := json.MarshalIndent(i, "", "   ")
	log.Println(string(j))
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		log.Print(err)
	}
}

func returnArrayJSON(playlists []models.Playlist, w http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
	}

	for _, playlist := range playlists {
		j, _ := json.MarshalIndent(playlist, "", "   ")
		log.Println(string(j))
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(j)
		if err != nil {
			log.Print(err)
		}
	}

}

func getArtistName(firstName string, lastName string) string {
	artistName := ""
	if lastName != "" {
		artistName = firstName + " " + lastName
	} else {
		artistName = firstName
	}
	return artistName
}

func concatenateName(artistName string) string {
	splitString := strings.Split(artistName, " ")
	newString := strings.Join(splitString, "_")
	return newString
}

func splitName(titleName string) string {
	splitString := strings.Split(titleName, "_")
	newString := strings.Join(splitString, " ")
	return newString
}

func getMp3Duration(path string) string {

	r, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		// return  err
	}
	tFloat := 0.0

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {

		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}

		tFloat = tFloat + f.Duration().Seconds()
	}
	r.Close()
	fmt.Println(`duration`, tFloat)
	t := int(tFloat)
	// min := int(math.Floor(tFloat / 60))
	// sec := t % 60
	minString := fmt.Sprint(math.Floor(tFloat / 60))
	secString := fmt.Sprint(t % 60)
	if len(secString) == 1 {
		secString = "0" + secString
	}
	minsec := minString + ":" + secString
	fmt.Println(minsec)
	return minsec
}
