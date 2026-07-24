DELETE FROM achievements WHERE code = 'excellent_student_4';

UPDATE achievements SET threshold = 5, description = 'Завершили 5 уроков' WHERE code = 'excellent_student_1';
UPDATE achievements SET threshold = 20, description = 'Завершили 20 уроков' WHERE code = 'excellent_student_2';
UPDATE achievements SET threshold = 50, description = 'Завершили 50 уроков' WHERE code = 'excellent_student_3';
