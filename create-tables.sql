DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS account;
DROP TABLE IF EXISTS transaction;

CREATE TABLE user (
  user_id  INT AUTO_INCREMENT NOT NULL,
  name VARCHAR(128) NOT NULL,
  email VARCHAR(128) NOT NULL,
  password VARCHAR(128) NOT NULL,
  PRIMARY KEY (`user_id`) 
);


CREATE TABLE account(
  account_id INT AUTO_INCREMENT NOT NULL,
  user_id INT NOT NULL,
  amount DOUBLE NOT NULL,
  PRIMARY KEY (`account_id`)
);

CREATE TABLE transaction(
  transaction_id INT AUTO_INCREMENT NOT NULL,
  account_id INT NOT NULL,
  to_account_id INT NOT NULL,
  amount DOUBLE NOT NULL,
  PRIMARY KEY (`transaction_id`)
);

INSERT INTO user (name, email, password) VALUES
('Viktor', 'vallmark.viktor@gmail.com', '12345');

INSERT INTO account (user_id, amount) VALUES
(1, 1000);
