CREATE TABLE Users (
                         user_id bigserial UNIQUE NOT NULL,
                         username varchar UNIQUE,
                         password varchar,
                         first_name varchar,
                         last_name varchar,
                         gender varchar,
                         admin bool Default FALSE NOT NULL,
                         PRIMARY KEY (user_id, username)
);



CREATE TABLE Artist (
                           name varchar,
                           artist_id int UNIQUE,
                           location varchar,
                           join_date date DEFAULT 'now()',
                           PRIMARY KEY (name, artist_id)
);

CREATE TABLE Song (
                        song_id bigserial UNIQUE NOT NULL,
                        title varchar,
                        artist_id int,
                        release_date date DEFAULT 'now()',
                        duration float NOT NULL,
                        album_id int,
                        total_plays bigint DEFAULT 0,
                        PRIMARY KEY (song_id, artist_id)
);

CREATE TABLE Songplay (
                            songplay_id bigserial UNIQUE,
                            session_id bigserial UNIQUE,
                            location varchar,
                            level varchar,
                            song_id integer,
                            artist_id integer,
                            user_id integer,
                            PRIMARY KEY (songplay_id, session_id, song_id, artist_id, user_id)
);


CREATE TABLE Playlist (
                            user_id integer,
                            name varchar,
                            playlist_id bigserial UNIQUE PRIMARY KEY
);

CREATE TABLE SongPlaylist(
                            song_id integer,
                            playlist_id integer,
                            name varchar,
                            PRIMARY KEY (playlist_id, song_id)

);

CREATE TABLE Album (
                         name varchar,
                         artist_id integer,
                         album_id bigserial UNIQUE PRIMARY KEY,
                         date_added date DEFAULT 'now()'
);

CREATE TABLE AlbumSong (
                            album_id integer,
                            song_id integer,
                            name varchar,
                           	PRIMARY KEY (album_id, song_id)
);


ALTER TABLE Song ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id);

ALTER TABLE Song ADD FOREIGN KEY (album_id) REFERENCES Album (album_id);

ALTER TABLE ARTIST ADD FOREIGN KEY (artist_id) REFERENCES Users (user_id);

ALTER TABLE Songplay ADD FOREIGN KEY (song_id) REFERENCES Song (song_id);

ALTER TABLE Songplay ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id);

ALTER TABLE Songplay ADD FOREIGN KEY (user_id) REFERENCES Users (user_id);

ALTER TABLE Playlist ADD FOREIGN KEY (user_id) REFERENCES Users (user_id);

ALTER TABLE Album ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id);

ALTER TABLE AlbumSong ADD FOREIGN KEY (album_id) REFERENCES Album (album_id);

ALTER TABLE AlbumSong ADD FOREIGN KEY (song_id) REFERENCES Song (song_id);

ALTER TABLE SongPlaylist ADD FOREIGN KEY (playlist_id) REFERENCES Playlist (playlist_id);

ALTER TABLE SongPlaylist ADD FOREIGN KEY (song_id) REFERENCES Song (song_id);