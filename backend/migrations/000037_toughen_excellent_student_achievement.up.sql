-- "Отличник" (excellent_student) was too easy to max out (5/20/50 lessons).
-- Raise it to match the other three-tier achievements' curve and add a
-- fourth, harder tier: 10/50/150/300 lessons completed.
UPDATE achievements SET threshold = 10, description = 'Завершили 10 уроков' WHERE code = 'excellent_student_1';
UPDATE achievements SET threshold = 50, description = 'Завершили 50 уроков' WHERE code = 'excellent_student_2';
UPDATE achievements SET threshold = 150, description = 'Завершили 150 уроков' WHERE code = 'excellent_student_3';

INSERT INTO achievements (code, title, description, tier, metric, threshold, xp_reward) VALUES
('excellent_student_4', 'Отличник IV', 'Завершили 300 уроков', 4, 'lessons_completed', 300, 1000);
