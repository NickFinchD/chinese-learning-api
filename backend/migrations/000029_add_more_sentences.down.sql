-- The first 5 rows (ids 1-5) were seeded by migration 000026; everything
-- added by this migration comes strictly after them.
DELETE FROM sentence_exercises WHERE id > 5;
