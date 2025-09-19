PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "library" (
	id integer primary key autoincrement,
	name text not null unique,
	path text not null unique,
	remote_path text null default '',
	last_scan_at datetime not null default '0000-00-00 00:00:00',
	updated_at datetime not null default current_timestamp,
	created_at datetime not null default current_timestamp
);
INSERT INTO "library" VALUES (1,'Music Library','/lib/path','','2025-09-16 13:51:59.001592505-04:00','2025-09-16 13:40:58.882881019-04:00','2025-09-16 17:40:58');
CREATE TABLE IF NOT EXISTS "media_file"
(
	id varchar(255) not null
		primary key,
	path varchar(255) default '' not null,
	title varchar(255) default '' not null,
	album varchar(255) default '' not null,
	artist varchar(255) default '' not null,
	artist_id varchar(255) default '' not null,
	album_artist varchar(255) default '' not null,
	album_id varchar(255) default '' not null,
	has_cover_art bool default FALSE not null,
	track_number integer default 0 not null,
	disc_number integer default 0 not null,
	year integer default 0 not null,
	size integer default 0 not null,
	suffix varchar(255) default '' not null,
	duration real default 0 not null,
	bit_rate integer default 0 not null,
	genre varchar(255) default '' not null,
	compilation bool default FALSE not null,
	created_at datetime,
	updated_at datetime
, full_text varchar(255) default '', album_artist_id varchar(255) default '', order_album_name varchar(255) collate nocase, order_album_artist_name varchar(255) collate nocase, order_artist_name varchar(255) collate nocase, sort_album_name varchar(255) collate nocase, sort_artist_name varchar(255) collate nocase, sort_album_artist_name varchar(255) collate nocase, sort_title varchar(255) collate nocase, disc_subtitle varchar(255), mbz_recording_id varchar(255), mbz_album_id varchar(255), mbz_artist_id varchar(255), mbz_album_artist_id varchar(255), mbz_album_type varchar(255), mbz_album_comment varchar(255), catalog_num varchar(255), comment varchar, bpm integer, channels integer, order_title varchar null collate NOCASE, mbz_release_track_id varchar(255), rg_album_gain real, rg_album_peak real, rg_track_gain real, rg_track_peak real, date varchar(255) default '' not null, original_year int default 0 not null, original_date varchar(255) default '' not null, release_year int default 0 not null, release_date varchar(255) default '' not null, lyrics JSONB default '[]', library_id INTEGER NOT NULL DEFAULT 0);
DELETE FROM "media_file";
INSERT INTO "main"."media_file" ("id", "path", "title", "album", "artist", "artist_id", "album_artist", "album_id", "has_cover_art", "track_number", "disc_number", "year", "size", "suffix", "duration", "bit_rate", "genre", "compilation", "created_at", "updated_at", "full_text", "album_artist_id", "order_album_name", "order_album_artist_name", "order_artist_name", "sort_album_name", "sort_artist_name", "sort_album_artist_name", "sort_title", "disc_subtitle", "mbz_recording_id", "mbz_album_id", "mbz_artist_id", "mbz_album_artist_id", "mbz_album_type", "mbz_album_comment", "catalog_num", "comment", "bpm", "channels", "order_title", "mbz_release_track_id", "rg_album_gain", "rg_album_peak", "rg_track_gain", "rg_track_peak", "date", "original_year", "original_date", "release_year", "release_date", "lyrics", "library_id") VALUES ('5c214deb5b2dba739e0d6af56f61d1c7', 'tests/fixtures/huey lewis - sports/hue lou.mp3', 'Hue Lou', 'Sports', 'Huey Lewis', 'b3d149f33481d7070d98724eef55b8c6', 'Huey Lewis', '9f08f5b8706718e5e129f14d88d5b3c1', '1', '1', '0', '1980', '15127935', 'mp3', '372.220001220703', '320', '', '0', '2024-01-12T13:09:51', '2024-01-12T12:12:50.4830029-05:00', 'bloopers', 'b3d149f33481d7070d98724eef55b8c6', 'foo', 'foo', 'foo', '', '', '', '', '', '', '', '', '', '', '', '', '', '0', '2', 'foo', '', '0.0', '1.0', '0.0', '1.0', '1980', '0', '', '0', '', '[]', '1');
INSERT INTO "main"."media_file" ("id", "path", "title", "album", "artist", "artist_id", "album_artist", "album_id", "has_cover_art", "track_number", "disc_number", "year", "size", "suffix", "duration", "bit_rate", "genre", "compilation", "created_at", "updated_at", "full_text", "album_artist_id", "order_album_name", "order_album_artist_name", "order_artist_name", "sort_album_name", "sort_artist_name", "sort_album_artist_name", "sort_title", "disc_subtitle", "mbz_recording_id", "mbz_album_id", "mbz_artist_id", "mbz_album_artist_id", "mbz_album_type", "mbz_album_comment", "catalog_num", "comment", "bpm", "channels", "order_title", "mbz_release_track_id", "rg_album_gain", "rg_album_peak", "rg_track_gain", "rg_track_peak", "date", "original_year", "original_date", "release_year", "release_date", "lyrics", "library_id") VALUES ('6ea5a2baa32842109925f67b3151fb80', 'tests/fixtures/mc5 - back in the usa/tutti fruitti.mp3', 'Tutti Fruitti', 'Back in the USA', 'MC5', '9bcc875883440c642594c4cb14f92832', 'MC5', '2d75e4e1738e20f7a3ce2bfb313f7e8d', '1', '1', '0', '2017', '12520563', 'mp3', '310.940002441406', '320', '', '0', '2024-01-12T13:14:51', '2024-01-14T19:48:22-05:00', 'whatev', '9bcc875883440c642594c4cb14f92832', 'foo', 'foo', 'foo', '', '', '', '', '', '', '', '', '', '', '', '', 'Visit https://xxxguyincognitoxxx.bandcamp.com', '0', '2', 'foo', '', '0.0', '1.0', '0.0', '1.0', '2017', '0', '', '0', '', '[]', '1');
COMMIT;
