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
                        Duration   int NOT NULL,
                        -- total_likes bigint default 0, Getting this through count()
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
                        message varchar(500),
                        created_date date DEFAULT now(),
                        isRead boolean default false,
                        message_id bigserial UNIQUE PRIMARY KEY
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

-- This Trigger adds a 'single' album if no album is available to add to.
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

-- This trigger delete the value from the likes table if it exists, when being added
CREATE OR REPLACE FUNCTION onLikeInsert() RETURNS trigger AS $$
BEGIN
		IF EXISTS (select likes.user_id from likes where likes.user_id = new.user_id AND likes.song_id = new.song_id) THEN
			DELETE FROM likes WHERE likes.user_id = new.user_id AND likes.song_id = new.song_id; 
END IF;

return new;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER onLikeInsert BEFORE INSERT ON Likes
    FOR EACH ROW EXECUTE FUNCTION onLikeInsert();


-- alerts all admin if a bad album date is added
CREATE OR REPLACE FUNCTION CheckAlbumDate() RETURNS TRIGGER AS $$
DECLARE 
    artist_join_date date;
BEGIN
	SELECT join_date into artist_join_date FROM ARTIST, ALBUM WHERE ARTIST.artist_id = new.artist_id AND album.album_id = new.album_id;
    IF (artist_join_date) > new.date_added 
		THEN 
-- 		DELETE FROM ALBUM WHERE album.Album_id = new.album_id; (maybe use?)
		INSERT INTO Messages select users.user_id, 
		CONCAT('ALBUM ID ', new.album_id, ' of date ', new.date_added, ' is before artist join date of ', artist_join_date, '.') 
		 from users where users.admin_level = 3;
    END IF;
    return new;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER CheckAlbumDate AFTER INSERT ON Album
    FOR EACH ROW EXECUTE FUNCTION CheckAlbumDate();


CREATE OR REPLACE FUNCTION CheckRatings() RETURNS TRIGGER AS $$
    DECLARE
        messageArtistId int;
        messageToSend varchar(500);
        totalLikes int;
        songTitle varchar;
    BEGIN
        SELECT sum(
        CASE
        WHEN likes.islike IS TRUE THEN 1
        ELSE 0
        END) into totalLikes from likes where likes.song_id=new.song_id;

        IF (totalLikes%2 = 0) THEN 

            select distinct song.artist_id into messageArtistId from likes, song 
            where likes.song_id = new.song_id and song.song_id = likes.song_id;

            select distinct song.title into songTitle from likes, song 
            where likes.song_id = new.song_id and song.song_id = likes.song_id;

            messageToSend = CONCAT('Your song ', songTitle, ' has reached ', totalLikes, ' likes!' );
                    RAISE NOTICE 'messageToSend %', messageToSend;

            IF new.islike = true
                AND NOT EXISTS (SELECT * FROM messages WHERE messages.user_id = messageArtistId and messages.message = messageToSend)
                THEN

                INSERT INTO Messages (user_id, message) values (messageArtistId, 
                messageToSend);
            END IF;

        END IF;
    return new;
    END;
$$ language plpgsql;

CREATE OR REPLACE TRIGGER CheckRatings AFTER INSERT ON Likes
    FOR EACH ROW EXECUTE FUNCTION CheckRatings();


-- for query 2. maybe adjust return column names
CREATE OR REPLACE VIEW likes_view AS
  SELECT COALESCE(( SELECT sum(
                CASE
                    WHEN likes.islike IS TRUE THEN 1
                    ELSE 0
                END) AS sum
           FROM likes
          WHERE likes.song_id = song.song_id), 0::bigint) AS likes,
    COALESCE(( SELECT sum(
                CASE
                    WHEN likes.islike IS FALSE THEN 1
                    ELSE 0
                END) AS sum
           FROM likes
          WHERE likes.song_id = song.song_id), 0::bigint) AS dislikes,
    song.*,
    artist.name AS artist_name,
    album.name AS album_name
    
   FROM song
     LEFT JOIN artist ON artist.artist_id = song.artist_id
     LEFT JOIN album ON album.album_id = song.album_id;



-- These artist/user reports should work. Might add a couple more columns (users add num playlists, num songs liked, )
create or replace view usersReport as
	select users.user_id, users.username, users.first_name, users.last_name, users.admin_level, users.join_date, count(playlist.playlist_id) playlist_count
	from users
	left join playlist on playlist.user_id = users.user_id
	group by users.user_id,users.username, users.first_name, users.last_name, users.admin_level, users.join_date
	order by users.user_id;

  create or replace view artistsReport as
    select artist.name, artist.artist_id, artist.join_date, 
    (select count(song.song_id) from song where song.artist_id = artist.artist_id) as num_songs,
    (select count(album.album_id) from Album where album.artist_id = artist.artist_id) as num_albums,
	sum(total_plays) total_plays,
	round(avg(total_plays),0) average_plays
	from artist, song
	where song.artist_id = artist.artist_id group by artist.name, artist.artist_id, artist.join_date  order by artist.artist_id;

CREATE OR REPLACE VIEW songReport as
    SELECT DISTINCT SONG.*, ARTIST.name as artist_name, ALBUM.NAME as album_name FROM SONG, ARTIST, ALBUM
    WHERE  ARTIST.name = (SELECT name from ARTIST WHERE ARTIST.Artist_id = Song.artist_id)
    AND ALBUM.name = (SELECT name from ALBUM WHERE ALBUM.ALBUM_ID = Song.ALBUM_ID)
    ORDER BY SONG.total_plays desc;