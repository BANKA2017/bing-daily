BEGIN TRANSACTION;
DROP TABLE IF EXISTS "bing";
CREATE TABLE IF NOT EXISTS "bing" (
	"startdate"	integer DEFAULT 0,
	"url"	text,
	"urlbase"	text,
	"copyright"	text,
	"copyrightlink"	text,
	"title"	text,
	"quiz"	text,
	"wp"	integer DEFAULT 0,
	"hsh"	text,
	"drk"	integer DEFAULT 0,
	"top"	integer DEFAULT 0,
	"bot"	integer DEFAULT 0,
	"hs"	text,
	"width"	INTEGER DEFAULT 0,
	"height"	INTEGER DEFAULT 0,
	"blurhash"	TEXT,
	"color"	TEXT DEFAULT 000000,
	PRIMARY KEY("startdate")
);
DROP TABLE IF EXISTS "d1_kv";
CREATE TABLE IF NOT EXISTS "d1_kv" (
	"key"	TEXT,
	"value"	TEXT NOT NULL,
	PRIMARY KEY("key")
);
COMMIT;
