-- This migration only appends rows after migrations 000026/000029 (ids 1-105);
-- everything it added comes strictly after them.
DELETE FROM sentence_exercises WHERE id > 105;
