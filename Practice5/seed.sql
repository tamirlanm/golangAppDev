
TRUNCATE TABLE user_friends, users RESTART IDENTITY CASCADE;

INSERT INTO users (id, name, email, gender, birth_date) VALUES
(1, 'Ali', 'ali@mail.com', 'male', '2002-01-10'),
(2, 'Aruzhan', 'aruzhan@mail.com', 'female', '2001-05-15'),
(3, 'Timur', 'timur@mail.com', 'male', '2000-03-12'),
(4, 'Aida', 'aida@mail.com', 'female', '2003-07-19'),
(5, 'Dias', 'dias@mail.com', 'male', '2002-11-01'),
(6, 'Dana', 'dana@mail.com', 'female', '2001-08-22'),
(7, 'Nursultan', 'nursultan@mail.com', 'male', '1999-04-17'),
(8, 'Madina', 'madina@mail.com', 'female', '2002-12-09'),
(9, 'Sanzhar', 'sanzhar@mail.com', 'male', '2000-06-06'),
(10, 'Kamila', 'kamila@mail.com', 'female', '2001-10-10'),
(11, 'Adil', 'adil@mail.com', 'male', '2003-01-25'),
(12, 'Amina', 'amina@mail.com', 'female', '2002-02-14'),
(13, 'Eldar', 'eldar@mail.com', 'male', '1998-09-30'),
(14, 'Zarina', 'zarina@mail.com', 'female', '2004-03-05'),
(15, 'Miras', 'miras@mail.com', 'male', '2001-07-21'),
(16, 'Indira', 'indira@mail.com', 'female', '2000-11-11'),
(17, 'Askar', 'askar@mail.com', 'male', '2002-04-08'),
(18, 'Laura', 'laura@mail.com', 'female', '2003-06-18'),
(19, 'Ruslan', 'ruslan@mail.com', 'male', '1999-12-24'),
(20, 'Tomiris', 'tomiris@mail.com', 'female', '2004-08-13');

INSERT INTO user_friends (user_id, friend_id) VALUES
-- user 1 friends
(1, 3),
(1, 4),
(1, 5),
(1, 6),
(1, 8),

-- user 2 friends
(2, 3),
(2, 4),
(2, 5),
(2, 7),
(2, 8),

-- reverse directions for friendship
(3, 1),
(4, 1),
(5, 1),
(6, 1),
(8, 1),

(3, 2),
(4, 2),
(5, 2),
(7, 2),
(8, 2),

-- extra friendships
(3, 9),
(9, 3),

(4, 10),
(10, 4),

(5, 11),
(11, 5),

(6, 12),
(12, 6),

(7, 13),
(13, 7),

(8, 14),
(14, 8),

(9, 15),
(15, 9),

(10, 16),
(16, 10),

(11, 17),
(17, 11),

(12, 18),
(18, 12),

(13, 19),
(19, 13),

(14, 20),
(20, 14);

SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));