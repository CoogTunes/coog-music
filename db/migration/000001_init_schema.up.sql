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
                      duration varchar(10) DEFAULT 'XX:YY',
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
-- CREATE TABLE Followers (
--                         user_id integer,
--                         artist_id integer,
--   					    PRIMARY KEY (user_id, artist_id)
-- );
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


ALTER TABLE Likes ADD FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Song ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Song ADD FOREIGN KEY (album_id) REFERENCES Album (album_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE ARTIST ADD FOREIGN KEY (artist_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Playlist ADD FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE Album ADD FOREIGN KEY (artist_id) REFERENCES ARTIST (artist_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE SongPlaylist ADD FOREIGN KEY (playlist_id) REFERENCES Playlist (playlist_id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE SongPlaylist ADD FOREIGN KEY (song_id) REFERENCES song (song_id) ON DELETE CASCADE ON UPDATE CASCADE;

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

DECLARE
		previousVal boolean;
BEGIN

IF EXISTS (select likes.isLike from likes
                    where likes.user_id = new.user_id
                    AND likes.song_id = new.song_id)
		 THEN
		 select likes.isLike into previousVal from likes
		 where likes.user_id = new.user_id
		 AND likes.song_id = new.song_id;
		 	if (previousVal) then new.isLike = null; end if;
	DELETE FROM likes WHERE likes.user_id = new.user_id AND likes.song_id = new.song_id;
END IF;
return new;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER onLikeInsert BEFORE INSERT ON Likes
    FOR EACH ROW EXECUTE FUNCTION onLikeInsert();


-- alerts all admin if a bad song date is added
CREATE OR REPLACE FUNCTION AlertSongDateBeforeArtistDate() RETURNS TRIGGER AS $$
DECLARE
artist_join_date date;
BEGIN
SELECT join_date into artist_join_date FROM ARTIST, SONG WHERE ARTIST.artist_id = new.artist_id AND SONG.song_id = new.song_id;
IF (artist_join_date) > new.uploaded_date
		THEN
		INSERT INTO Messages select users.user_id,
                                    CONCAT('Song ID ', new.song_id, ' of date ', new.uploaded_date, ' is before artist join date of ', artist_join_date, '.')
                             from users where users.admin_level = 3;
END IF;
return new;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER AlertSongDateBeforeArtistDate AFTER INSERT ON Song
    FOR EACH ROW EXECUTE FUNCTION AlertSongDateBeforeArtistDate();


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

-- CREATE OR REPLACE FUNCTION notifyFollowersOfNewSong() RETURNS TRIGGER AS $$
--     DECLARE
--         artistName varchar;
--         messageToSend varchar(500);
--     BEGIN
-- 		SELECT likes_view.artist_name into artistName from likes_view where likes_view.song_id = new.song_id;
-- 		messageToSend = CONCAT(artistName, ' has released a new song called ', new.title, '!' );
-- 		insert into messages select followers.user_id, messageToSend from followers where artist_id = new.artist_id; 
--     return new;
--     END;
-- $$ language plpgsql;

-- CREATE OR REPLACE TRIGGER notifyFollowersOfNewSong AFTER INSERT ON Song
--     FOR EACH ROW EXECUTE FUNCTION notifyFollowersOfNewSong();

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



-- Grabs regular user data minus password, plus total playlist count, total number songs user has liked, and most common artist in their playlists
CREATE OR REPLACE view usersReport AS
SELECT users.user_id,
       users.username,
       users.first_name,
       users.last_name,
       users.admin_level,
       users.join_date,
       count(p1.playlist_id) AS playlist_count,
       COALESCE(( SELECT sum(
                                 CASE
                                     WHEN likes.islike IS TRUE THEN 1
                                     ELSE 0
                                     END) AS sum
                FROM likes,
                song
                WHERE likes.song_id = song.song_id AND likes.user_id = users.user_id), 0) AS liked_songs_count,
       COALESCE(( SELECT artist.name
                  FROM songplaylist
                           LEFT JOIN song ON song.song_id = songplaylist.song_id
                           LEFT JOIN artist ON artist.artist_id = song.artist_id
                  WHERE (songplaylist.playlist_id IN ( SELECT p2.playlist_id
                                                       FROM playlist as p2
                                                       WHERE p2.user_id = users.user_id))
                  GROUP BY artist.name
                  ORDER BY (count(artist.name)) DESC
                LIMIT 1), 'N/A') AS common_artist
FROM users
         LEFT JOIN playlist as p1 ON p1.user_id = users.user_id
GROUP BY users.user_id, users.username, users.first_name, users.last_name, users.admin_level, users.join_date
ORDER BY users.user_id;

CREATE OR REPLACE VIEW artistsReport AS
SELECT ar1.name,
       ar1.artist_id,
       ar1.join_date,
       ( SELECT count(song_1.song_id) AS count
FROM song song_1
WHERE song_1.artist_id = ar1.artist_id) AS num_songs,
    ( SELECT count(album.album_id) AS count
FROM album
WHERE album.artist_id = ar1.artist_id) AS num_albums,
    sum(song.total_plays) AS total_plays,
    round(avg(song.total_plays), 0) AS average_plays,
    COALESCE(( SELECT song_1.title
    FROM likes
    LEFT JOIN song song_1 ON likes.song_id = song_1.song_id
    LEFT JOIN artist ar2 ON ar2.artist_id = song_1.artist_id
    WHERE song_1.artist_id = ar1.artist_id
    GROUP BY song_1.title
    ORDER BY (count(song_1.title)) DESC
    LIMIT 1), 'N/A'::character varying) AS most_liked_song
FROM artist ar1,
    song
WHERE song.artist_id = ar1.artist_id
GROUP BY ar1.name, ar1.artist_id, ar1.join_date
ORDER BY ar1.artist_id;

-- likes_view is better
-- CREATE OR REPLACE VIEW songReport as
--     SELECT DISTINCT SONG.*, ARTIST.name as artist_name, ALBUM.NAME as album_name FROM SONG, ARTIST, ALBUM
--     WHERE  ARTIST.name = (SELECT name from ARTIST WHERE ARTIST.Artist_id = Song.artist_id)
--     AND ALBUM.name = (SELECT name from ALBUM WHERE ALBUM.ALBUM_ID = Song.ALBUM_ID)
--     ORDER BY SONG.total_plays desc;


