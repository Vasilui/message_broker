CREATE TYPE task_type AS ENUM
    ('replenishment',
    'write-off');

CREATE TYPE task_status AS ENUM
    ('open',
    'close',
    'rejected');

CREATE TABLE users(
      id serial PRIMARY KEY,
      username character varying(50) NOT NULL,
      balance bigint NOT NULL
);

CREATE TABLE tasks(
      id bigserial PRIMARY KEY,
      user_id serial NOT NULL,
      type_task task_type NOT NULL,
      amount int8 NOT NULL,
      status task_status DEFAULT 'open',
      created timestamp with time zone DEFAULT Now()
);
