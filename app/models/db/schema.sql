CREATE TABLE agent  (
  id int UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name varchar(255) NOT NULL,
  enable tinyint NOT NULL,
  type_id int NOT NULL,
  event_max_age bigint NOT NULL,
  prop_json_str text NOT NULL,
  create_time datetime NOT NULL,
  deleted tinyint NOT NULL
);

CREATE TABLE agent_relation  (
  id int UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  src_agent_id int NOT NULL,
  dst_agent_id int NOT NULL,
  deleted tinyint NOT NULL
);

CREATE TABLE agent_type  (
  id int UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name varchar(255) NOT NULL,
  allow_input tinyint NOT NULL,
  allow_output tinyint NOT NULL,
  deleted tinyint NOT NULL
);

CREATE TABLE event  (
  id int UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  json_str text NOT NULL,
  src_agent_id int NOT NULL,
  created_time datetime NOT NULL,
  deleted tinyint NOT NULL
);