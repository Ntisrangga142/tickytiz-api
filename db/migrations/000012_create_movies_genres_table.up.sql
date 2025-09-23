CREATE TABLE public.movies_genres (
  id_movie INTEGER NOT NULL,
  id_genre INTEGER NOT NULL,
  CONSTRAINT movies_genres_pk PRIMARY KEY (id_movie, id_genre),
  CONSTRAINT fk_id_movie_genres FOREIGN KEY (id_movie) REFERENCES public.movies (id),
  CONSTRAINT fk_id_genre_movie  FOREIGN KEY (id_genre) REFERENCES public.genres (id)
);