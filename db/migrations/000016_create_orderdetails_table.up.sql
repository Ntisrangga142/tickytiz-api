CREATE TABLE public.orderdetails (
  id_order  INTEGER       NOT NULL,
  id_seat   VARCHAR(255)  NOT NULL,
  CONSTRAINT orderdetails_pk PRIMARY KEY (id_order, id_seat),
  CONSTRAINT fk_id_order_details FOREIGN KEY (id_order) REFERENCES public.orders (id),
  CONSTRAINT fk_id_seat_details  FOREIGN KEY (id_seat)  REFERENCES public.seat (id)
);
