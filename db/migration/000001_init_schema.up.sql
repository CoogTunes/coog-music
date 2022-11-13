CREATE TABLE Users (
                        user_id bigserial UNIQUE NOT NULL,
                        username varchar UNIQUE,
                        password varchar,
                        first_name varchar,
                        last_name varchar,
                        admin_level int CHECK (admin_level >= 0 AND admin_level <= 3),
                        join_date date DEFAULT 'now()',
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

CREATE TABLE Album (
                        name varchar,
                        artist_id integer,
                        album_id bigserial UNIQUE PRIMARY KEY,
                        date_added date DEFAULT 'now()'
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
                        total_likes bigint default 0,
                        PRIMARY KEY (song_id, artist_id)
);

CREATE TABLE Playlist (
                        user_id integer,
                        name varchar,
                        playlist_id bigserial UNIQUE PRIMARY KEY
);

CREATE TABLE SongPlaylist(
                        song_id integer,
                        playlist_id integer,
                        PRIMARY KEY (playlist_id, song_id)
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


-- NOT USING?
-- CREATE TABLE AlbumSong (
--                         album_id integer,
--                         song_id integer,
--                         PRIMARY KEY (album_id, song_id)
-- );

-- CREATE TABLE Songplay (
--                         songplay_id bigserial UNIQUE,
--                         session_id bigserial UNIQUE,
--                         location varchar,
--                         level varchar,
--                         song_id integer,
--                         artist_id integer,
--                         user_id integer,
--                         PRIMARY KEY (songplay_id, session_id, song_id, artist_id, user_id)
-- );

-- CREATE TABLE Followers (
--                         user_id integer,
--                         artist_id integer,
--   						          PRIMARY KEY (user_id, artist_id)
-- );

ALTER TABLE Likes ADD FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Song ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Song ADD FOREIGN KEY (album_id) REFERENCES Album (album_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE ARTIST ADD FOREIGN KEY (artist_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Playlist ADD FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Album ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE SongPlaylist ADD FOREIGN KEY (playlist_id) REFERENCES Playlist (playlist_id) ON DELETE CASCADE ON UPDATE CASCADE;

-- NOT USING?
-- ALTER TABLE AlbumSong ADD FOREIGN KEY (album_id) REFERENCES Album (album_id) ON DELETE CASCADE ON UPDATE CASCADE;

-- ALTER TABLE AlbumSong ADD FOREIGN KEY (song_id) REFERENCES Song (song_id) ON DELETE CASCADE ON UPDATE CASCADE;

-- ALTER TABLE Songplay ADD FOREIGN KEY (song_id) REFERENCES Song (song_id) ON DELETE CASCADE ON UPDATE CASCADE;

-- ALTER TABLE Songplay ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id) ON DELETE CASCADE ON UPDATE CASCADE;

-- ALTER TABLE Songplay ADD FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

-- ALTER TABLE Followers ADD FOREIGN KEY (user_id) REFERENCES Users (user_id)ON DELETE CASCADE ON UPDATE CASCADE;

-- ALTER TABLE Followers ADD FOREIGN KEY (artist_id) REFERENCES Artist (artist_id) ON DELETE CASCADE ON UPDATE CASCADE;


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


-- for query 2. maybe adjust return column names
create or replace view likes_view as 
select 
	sum(case when likes.islike is true then 1 else 0 end) as likes,
	sum(case when likes.islike is false then 1 else 0 end) as dislikes,
	song.song_id, song.title as song_title, artist.name as artist_name, album.name as album_name, song.uploaded_date
from likes, song, users, artist, album
where 
	song.song_id = likes.song_id  
	and users.user_id = likes.user_id
	and artist.artist_id = song.artist_id
	and song.album_id = album.album_id
group by song.song_id, song.title, artist.name, album.name, song.uploaded_date;



-- These artist/user reports should work. Might add a couple more columns (users add num playlists, num songs liked, )
-- create or replace view usersReport as
-- 	select users.user_id, users.username, users.first_name, users.last_name, users.admin_level, users.join_date 
-- 	from users
-- 	order by users.last_name;

--   create or replace view artistsReport2 as
--     select artist.name, artist.artist_id, artist.join_date, 
--     (select count(song.song_id) from song where song.artist_id = artist.artist_id) as numSongs,
--     (select count(album.album_id) from Album where album.artist_id = artist.artist_id) as numAlbums	
-- 	from artist
-- 	order by artist."name";