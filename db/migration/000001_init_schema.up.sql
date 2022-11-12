CREATE TABLE Users (
                        user_id bigserial UNIQUE NOT NULL,
                        username varchar UNIQUE,
                        password varchar,
                        first_name varchar,
                        last_name varchar,
                        admin_level int CHECK (admin_level >= 0 AND admin_level <= 3),
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
);

CREATE TABLE Likes (
              user_id int,
              song_id int,
							isLike boolean,
						  PRIMARY KEY (user_id, song_id)
);

ALTER TABLE Likes ADD FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

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

-- CREATE OR REPLACE FUNCTION CheckAlbumDate() RETURNS TRIGGER AS $$
-- DECLARE 
--     uid int;
-- BEGIN
--     IF (SELECT join_date FROM ARTIST, ALBUM WHERE artist_id = old.artist_id) > old.date_added 
--     THEN DELETE FROM ALBUM WHERE Album_id = old.album_id;
--     SELECT user_id INTO uid FROM USERS WHERE admin_level >= 2;
--     INSERT INTO Messages (user_id, admin_level, message) values(uid, 2, 'Album aborted');
    
--     END IF;
--     return new;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE OR REPLACE TRIGGER CheckAlbumDate AFTER INSERT ON Album
--     FOR EACH ROW EXECUTE FUNCTION CheckAlbumDate();


-- CREATE OR REPLACE FUNCTION CheckRatings() RETURNS TRIGGER AS $$
-- BEGIN
--     IF new.total_likes > 10 
--     THEN
--     INSERT INTO MESSAGES (user_id, admin_level, message) VALUES (new.artist_id, 2, 'Your song has reached 10 likes!')
--     END IF
--     return new;
-- END;
-- $$ language plpgsql

-- CREATE OR REPLACE TRIGGER CheckRatings AFTER UPDATE ON Song
--     FOR EACH ROW EXECUTE FUCTION CheckRatings();