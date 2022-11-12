CREATE TABLE Users (
                        user_id bigserial UNIQUE NOT NULL,
                        username varchar UNIQUE,
                        password varchar,
                        first_name varchar,
                        last_name varchar,
                        admin_level int CHECK (admin_level >= 0 AND admin_level <= 3),
                        last_login date DEFAULT now(),
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
                        album_id int NOT NULL,
                        artist_id int,
                        song_path varchar,
                        cover_path varchar,
                        uploaded_date date DEFAULT 'now()',
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

CREATE TABLE Messages (
                        user_id int,
                        admin_level int CHECK (admin_level >= 0 AND admin_level <= 3),
                        message varchar(500),
                        created_date date DEFAULT now()
)



-- tested some and looks like it works? CREATE TABLE Likes (
--                         user_id int,
--                         song_id int,
-- 							isLike boolean,
-- 						  PRIMARY KEY (user_id, song_id)
-- );
-- ALTER TABLE Likes ADD FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;


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
		
		if album_date > new.uploaded_date THEN
			new.uploaded_date = album_date;
		END if;
		
		return new;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER checkSongDate BEFORE INSERT OR UPDATE ON Song
    FOR EACH ROW EXECUTE FUNCTION checkSongDate();

CREATE OR REPLACE FUNCTION addAlbumIfSingle() RETURNS trigger AS $$
BEGIN
		IF new.album_id IS NULL THEN
			INSERT INTO Album (name, artist_id, date_added)
			VALUES (new.title, new.artist_id, new.uploaded_date)
			RETURNING album.album_id into new.album_id;
END IF;

return new;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER addAlbumIfSingle BEFORE INSERT ON Song
    FOR EACH ROW EXECUTE FUNCTION addAlbumIfSingle();