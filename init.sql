CREATE TABLE style (
  id SERIAL PRIMARY KEY ,
  name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE artist (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE music (
  id SERIAL PRIMARY KEY,
  artist_id INT NOT NULL,
  name VARCHAR(50) UNIQUE NOT NULL,
  link VARCHAR(200) UNIQUE NOT NULL,
  cover VARCHAR(300) NOT NULL,
  CONSTRAINT fk_music_artist_id
  FOREIGN KEY (artist_id) 
  REFERENCES artist (id)
);

CREATE TABLE music_style ( 
  id SERIAL PRIMARY KEY,
  music_id INT NOT NULL,
  style_id INT NOT NULL,
  CONSTRAINT fk_music_style_music_id
  FOREIGN KEY (music_id)
  REFERENCES music (id),
  CONSTRAINT fk_music_style_id
  FOREIGN KEY (style_id)
  REFERENCES style (id)
);
