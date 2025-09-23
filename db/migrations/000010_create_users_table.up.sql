CREATE TABLE public.users (
  id              INTEGER PRIMARY KEY,
  profileimg      VARCHAR(255),
  firstname       VARCHAR(255),
  lastname        VARCHAR(255), 
  phone           VARCHAR(255),
  point           INTEGER       DEFAULT 0,
  virtual_account VARCHAR(255),
  update_at       TIMESTAMP,
  CONSTRAINT fk_account_user FOREIGN KEY (id) REFERENCES public.account (id)
);
