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
CREATE TABLE IF NOT EXISTS "media_file" (
	"id"	varchar(255) NOT NULL,
	"path"	varchar(255) NOT NULL DEFAULT '',
	"title"	varchar(255) NOT NULL DEFAULT '',
	"album"	varchar(255) NOT NULL DEFAULT '',
	"artist"	varchar(255) NOT NULL DEFAULT '',
	"artist_id"	varchar(255) NOT NULL DEFAULT '',
	"album_artist"	varchar(255) NOT NULL DEFAULT '',
	"album_id"	varchar(255) NOT NULL DEFAULT '',
	"has_cover_art"	bool NOT NULL DEFAULT FALSE,
	"track_number"	integer NOT NULL DEFAULT 0,
	"disc_number"	integer NOT NULL DEFAULT 0,
	"year"	integer NOT NULL DEFAULT 0,
	"size"	integer NOT NULL DEFAULT 0,
	"suffix"	varchar(255) NOT NULL DEFAULT '',
	"duration"	real NOT NULL DEFAULT 0,
	"bit_rate"	integer NOT NULL DEFAULT 0,
	"genre"	varchar(255) NOT NULL DEFAULT '',
	"compilation"	bool NOT NULL DEFAULT FALSE,
	"created_at"	datetime,
	"updated_at"	datetime,
	"full_text"	varchar(255) DEFAULT '',
	"album_artist_id"	varchar(255) DEFAULT '',
	"order_album_name"	varchar(255) COLLATE nocase,
	"order_album_artist_name"	varchar(255) COLLATE nocase,
	"order_artist_name"	varchar(255) COLLATE nocase,
	"sort_album_name"	varchar(255) COLLATE nocase,
	"sort_artist_name"	varchar(255) COLLATE nocase,
	"sort_album_artist_name"	varchar(255) COLLATE nocase,
	"sort_title"	varchar(255) COLLATE nocase,
	"disc_subtitle"	varchar(255),
	"mbz_recording_id"	varchar(255),
	"mbz_album_id"	varchar(255),
	"mbz_artist_id"	varchar(255),
	"mbz_album_artist_id"	varchar(255),
	"mbz_album_type"	varchar(255),
	"mbz_album_comment"	varchar(255),
	"catalog_num"	varchar(255),
	"comment"	varchar,
	"bpm"	integer,
	"channels"	integer,
	"order_title"	varchar COLLATE NOCASE,
	"mbz_release_track_id"	varchar(255),
	"rg_album_gain"	real,
	"rg_album_peak"	real,
	"rg_track_gain"	real,
	"rg_track_peak"	real,
	"date"	varchar(255) NOT NULL DEFAULT '',
	"original_year"	int NOT NULL DEFAULT 0,
	"original_date"	varchar(255) NOT NULL DEFAULT '',
	"release_year"	int NOT NULL DEFAULT 0,
	"release_date"	varchar(255) NOT NULL DEFAULT '',
	"lyrics"	JSONB DEFAULT '[]',
	"library_id"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id")
);
INSERT INTO media_file VALUES('5c214deb5b2dba739e0d6af56f61d1c7','music/Crazy Rhythms/feelies, the - crazy rhythms - 09 - crazy rhythms.mp3','Crazy Rhythms','Crazy Rhythms','The Feelies','b3d149f33481d7070d98724eef55b8c6','The Feelies','9f08f5b8706718e5e129f14d88d5b3c1',1,9,0,1980,15127935,'mp3',372.22000122070301132,320,'',0,'2024-01-23T13:57:19.360759986-05:00','2024-01-12T12:12:50.4830029-05:00',' crazy feelies rhythms the','b3d149f33481d7070d98724eef55b8c6','Crazy Rhythms','Feelies','Feelies','','','','','','','','','','','','','',0,2,'Crazy Rhythms','',0.0,1.0,0.0,1.0,'1980',0,'',0,'','[]',1);
INSERT INTO media_file VALUES('37141ae2932c8e06cc3716c3b9c55a48','music/Crazy Rhythms/feelies, the - crazy rhythms - 15 - i wanna sleep in your arms (live).mp3','I Wanna Sleep in Your Arms (live)','Crazy Rhythms','The Feelies','b3d149f33481d7070d98724eef55b8c6','The Feelies','9f08f5b8706718e5e129f14d88d5b3c1',1,15,0,1980,4949641,'mp3',117.76000213622999978,320,'',0,'2024-01-23T13:57:19.360841218-05:00','2024-01-12T12:12:56.6038184-05:00',' arms crazy feelies i in live rhythms sleep the wanna your','b3d149f33481d7070d98724eef55b8c6','Crazy Rhythms','Feelies','Feelies','','','','','','','','','','','','','',0,2,'I Wanna Sleep in Your Arms (live)','',0.0,1.0,0.0,1.0,'1980',0,'',0,'','[]',1);
INSERT INTO media_file VALUES('54c5999927b56e2887c3a5cfd21bdfbf','music/Crazy Rhythms/feelies, the - crazy rhythms - 13 - moscow nights (demo).mp3','Moscow Nights (demo)','Crazy Rhythms','The Feelies','b3d149f33481d7070d98724eef55b8c6','The Feelies','9f08f5b8706718e5e129f14d88d5b3c1',1,13,0,1980,10601468,'mp3',259.05999755859397738,320,'',0,'2024-01-23T13:57:19.360927352-05:00','2024-01-12T12:12:54.3347391-05:00',' crazy demo feelies moscow nights rhythms the','b3d149f33481d7070d98724eef55b8c6','Crazy Rhythms','Feelies','Feelies','','','','','','','','','','','','','',0,2,'Moscow Nights (demo)','',0.0,1.0,0.0,1.0,'1980',0,'',0,'','[]',1);
INSERT INTO media_file VALUES('1fde8382304840139358a101b081db9c','music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 02 Can We Please Go Back To How It Was Before I messed Up-.mp3','Can We Please Go Back To How It Was Before I messed Up?','Lovedrug','Guy Incognito','9bcc875883440c642594c4cb14f92832','Guy Incognito','2d75e4e1738e20f7a3ce2bfb313f7e8d',1,2,0,2017,12192231,'mp3',302.73001098632801131,320,'',0,'2024-01-23T14:57:19.456103508-05:00','2024-01-14T19:48:22-05:00',' back before can messed go guy how i incognito it lovedrug please to up? was we','9bcc875883440c642594c4cb14f92832','Lovedrug','Guy Incognito','Guy Incognito','','','','','','','','','','','','','Visit https://xxxguyincognitoxxx.bandcamp.com',0,2,'Can We Please Go Back To How It Was Before I messed Up?','',0.0,1.0,0.0,1.0,'2017',0,'',0,'','[]',1);
INSERT INTO media_file VALUES('07bb5aec148d087b0192c538721d0627','music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 03 Seven.mp3','Seven','Lovedrug','Guy Incognito','9bcc875883440c642594c4cb14f92832','Guy Incognito','2d75e4e1738e20f7a3ce2bfb313f7e8d',1,3,0,2017,10491428,'mp3',260.260009765625,320,'',0,'2024-01-23T14:57:19.456103508-05:00','2024-01-14T19:48:22-05:00',' guy incognito lovedrug seven','9bcc875883440c642594c4cb14f92832','Lovedrug','Guy Incognito','Guy Incognito','','','','','','','','','','','','','Visit https://xxxguyincognitoxxx.bandcamp.com',0,2,'Seven','',0.0,1.0,0.0,1.0,'2017',0,'',0,'','[]',1);
INSERT INTO media_file VALUES('6ea5a2baa32842109925f67b3151fb80','music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 01 Lovedrug.mp3','Lovedrug','Lovedrug','Guy Incognito','9bcc875883440c642594c4cb14f92832','Guy Incognito','2d75e4e1738e20f7a3ce2bfb313f7e8d',1,1,0,2017,12520563,'mp3',310.94000244140602262,320,'',0,'2024-01-23T14:57:19.456103508-05:00','2024-01-14T19:48:22-05:00',' guy incognito lovedrug','9bcc875883440c642594c4cb14f92832','Lovedrug','Guy Incognito','Guy Incognito','','','','','','','','','','','','','Visit https://xxxguyincognitoxxx.bandcamp.com',0,2,'Lovedrug','',0.0,1.0,0.0,1.0,'2017',0,'',0,'','[]',1);
COMMIT;
