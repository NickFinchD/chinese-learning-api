-- A final small batch of high-frequency words (与, 美, 色, 语言, 幸福...)
-- turned up by a third pass of the same tokenizer check as migrations
-- 000041/000042 — 与 alone appears bare 22 times across the 52 texts.
INSERT INTO words (hanzi, pinyin, translation, part_of_speech, hsk_level) VALUES
('与', 'yǔ', 'и, с (книжный союз)', 'conjunction', 4),
('美', 'měi', 'красивый', 'adjective', 2),
('色', 'sè', 'цвет', 'noun', 2),
('语言', 'yǔyán', 'язык (речь)', 'noun', 3),
('幸福', 'xìngfú', 'счастье; счастливый', 'noun', 3),
('座', 'zuò', 'счётное слово (для гор, зданий)', 'measure_word', 3),
('顿', 'dùn', 'счётное слово (для приёмов пищи)', 'measure_word', 3),
('交', 'jiāo', 'передавать, сдавать', 'verb', 4),
('充实', 'chōngshí', 'насыщенный, содержательный', 'adjective', 5),
('洗', 'xǐ', 'мыть', 'verb', 2),
('香', 'xiāng', 'ароматный', 'adjective', 3),
('景色', 'jǐngsè', 'пейзаж', 'noun', 3),
('节', 'jié', 'праздник', 'noun', 3),
('亲近', 'qīnjìn', 'близкий, родной', 'adjective', 5),
('平', 'píng', 'ровный, плоский', 'adjective', 4);
