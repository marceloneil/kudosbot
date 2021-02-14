CREATE TABLE kudos(
  id uuid primary key NOT NULL,
  timestamp timestamptz NOT NULL,
  sender_id text NOT NULL,
  receiver_id text NOT NULL,
  event_id text NOT NULL,
  reaction text NOT NULL,
  UNIQUE(sender_id, receiver_id, event_id)
);

CREATE INDEX kudos_timestamp_sender_idx ON kudos (timestamp, sender_id);
CREATE INDEX kudos_timestamp_receiver_idx ON kudos (timestamp, receiver_id);
