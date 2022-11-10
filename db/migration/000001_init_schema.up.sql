CREATE TABLE Users (
                        user_id bigserial UNIQUE NOT NULL,
                        username varchar UNIQUE,
                        password varchar,
                        first_name varchar,
                        last_name varchar,
                        admin_level int CHECK (admin_level >= 1 AND admin_level <= 3),
                        -- 1 = basic
                        -- 2 = artist
                        -- 3 = admin
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

CREATE TABLE Album (
                        name varchar,
                        artist_id integer,
                        album_id bigserial UNIQUE PRIMARY KEY,
                        date_added date DEFAULT 'now()'
);

CREATE TABLE AlbumSong (
                            album_id integer,
                            song_id integer,
                           	PRIMARY KEY (album_id, song_id)
);

CREATE TABLE SongPlaylist(
                            song_id integer,
                            playlist_id integer,
                            PRIMARY KEY (playlist_id, song_id)

);

CREATE TABLE Followers (
                        user_id integer,
                        artist_id integer,
  						PRIMARY KEY (user_id, artist_id)
);


ALTER TABLE Song ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Song ADD FOREIGN KEY (album_id) REFERENCES Album (album_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE ARTIST ADD FOREIGN KEY (artist_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Songplay ADD FOREIGN KEY (song_id) REFERENCES Song (song_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Songplay ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Songplay ADD FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Playlist ADD FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Album ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE AlbumSong ADD FOREIGN KEY (album_id) REFERENCES Album (album_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE AlbumSong ADD FOREIGN KEY (song_id) REFERENCES Song (song_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE SongPlaylist ADD FOREIGN KEY (playlist_id) REFERENCES Playlist (playlist_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE SongPlaylist ADD FOREIGN KEY (song_id) REFERENCES Song (song_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Followers ADD FOREIGN KEY (user_id) REFERENCES Users (user_id)ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Followers ADD FOREIGN KEY (artist_id) REFERENCES Artist (artist_id) ON DELETE CASCADE ON UPDATE CASCADE;

CREATE OR REPLACE FUNCTION checkAlbumDate() RETURNS trigger AS $$
	DECLARE
		artist_date date;
    BEGIN
		select ar.join_date into artist_date from artist as ar where new.artist_id = ar.artist_id;
		
		if artist_date > new.date_added THEN
			new.date_added = artist_date;
		END if;
		
		return new;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER checkAlbumDate BEFORE INSERT OR UPDATE ON Album
    FOR EACH ROW EXECUTE FUNCTION checkAlbumDate();
	
	
CREATE OR REPLACE FUNCTION checkSongDate() RETURNS trigger AS $$
	DECLARE
		album_date date;
    BEGIN
		select al.date_added into album_date from album as al where new.album_id = al.album_id;
		
		if album_date > new.release_date THEN
			new.release_date = album_date;
		END if;
		
		return new;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER checkSongDate BEFORE INSERT OR UPDATE ON Song
    FOR EACH ROW EXECUTE FUNCTION checkSongDate();
