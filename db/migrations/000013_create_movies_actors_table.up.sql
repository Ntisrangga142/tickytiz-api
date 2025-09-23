CREATE TABLE public.movies_actors (
  id_movie INTEGER NOT NULL,
  id_actor INTEGER NOT NULL,
  CONSTRAINT movies_actors_pk PRIMARY KEY (id_movie, id_actor),
  CONSTRAINT fk_id_movie_actor FOREIGN KEY (id_movie) REFERENCES public.movies (id),
  CONSTRAINT fk_id_actor_movie FOREIGN KEY (id_actor) REFERENCES public.actors (id)
);