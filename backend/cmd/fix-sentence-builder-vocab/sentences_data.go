package main

// newSentence describes a sentence_exercises row this command seeds before
// reassigning lessons. Every chunk is a word actually taught somewhere in
// HSK1 or HSK2 (verified against the words table), unlike most of the
// original hand-picked pool, so the cumulative-vocab matching in main() has
// something valid to assign even early in each course.
type newSentence struct {
	Translation string
	Chunks      []string
	Pinyin      string
	HSKLevel    int16
}

var newSentences = []newSentence{
	{Translation: "Я тоже.", Chunks: []string{"我", "也", "是"}, Pinyin: "wǒ yě shì", HSKLevel: 1},
	{Translation: "А ты?", Chunks: []string{"你", "呢"}, Pinyin: "nǐ ne", HSKLevel: 1},
	{Translation: "У него нет.", Chunks: []string{"他", "没", "有"}, Pinyin: "tā méi yǒu", HSKLevel: 1},
	{Translation: "Это моё.", Chunks: []string{"这", "是", "我", "的"}, Pinyin: "zhè shì wǒ de", HSKLevel: 1},
	{Translation: "Что это?", Chunks: []string{"这", "是", "什么"}, Pinyin: "zhè shì shén me", HSKLevel: 1},
	{Translation: "Кто он?", Chunks: []string{"他", "是", "谁"}, Pinyin: "tā shì shéi", HSKLevel: 1},
	{Translation: "Как насчёт того?", Chunks: []string{"那", "个", "怎么样"}, Pinyin: "nà gè zěn me yàng", HSKLevel: 1},
	{Translation: "Кто твой папа?", Chunks: []string{"你", "爸爸", "是", "谁"}, Pinyin: "nǐ bà ba shì shéi", HSKLevel: 1},
	{Translation: "Я хочу пойти.", Chunks: []string{"我", "想", "去"}, Pinyin: "wǒ xiǎng qù", HSKLevel: 1},
	{Translation: "Ты умеешь это делать?", Chunks: []string{"你", "会", "做", "吗"}, Pinyin: "nǐ huì zuò ma", HSKLevel: 1},
	{Translation: "Сегодня очень жарко.", Chunks: []string{"今天", "很", "热"}, Pinyin: "jīn tiān hěn rè", HSKLevel: 1},
	{Translation: "Это слишком маленькое.", Chunks: []string{"这", "个", "太", "小", "了"}, Pinyin: "zhè gè tài xiǎo le", HSKLevel: 1},
	{Translation: "Мой друг — учитель.", Chunks: []string{"我", "朋友", "是", "老师"}, Pinyin: "wǒ péng you shì lǎo shī", HSKLevel: 1},
	{Translation: "Эта книга очень хорошая.", Chunks: []string{"这", "本", "书", "很", "好"}, Pinyin: "zhè běn shū hěn hǎo", HSKLevel: 1},
	{Translation: "Мы идём в магазин.", Chunks: []string{"我们", "去", "商店"}, Pinyin: "wǒ men qù shāng diàn", HSKLevel: 1},
	{Translation: "Он живёт в Пекине.", Chunks: []string{"他", "住", "在", "北京"}, Pinyin: "tā zhù zài běi jīng", HSKLevel: 1},
	{Translation: "Как тебя зовут?", Chunks: []string{"你", "叫", "什么", "名字"}, Pinyin: "nǐ jiào shén me míng zi", HSKLevel: 1},
	{Translation: "Мне очень нравятся кошки.", Chunks: []string{"我", "很", "喜欢", "猫"}, Pinyin: "wǒ hěn xǐ huan māo", HSKLevel: 1},
	{Translation: "Почему ты не рад?", Chunks: []string{"你", "为什么", "不", "高兴"}, Pinyin: "nǐ wèi shén me bù gāo xìng", HSKLevel: 2},
	{Translation: "Кроме тебя, мы все идём.", Chunks: []string{"除了", "你", "我们", "都", "去"}, Pinyin: "chú le nǐ wǒ men dōu qù", HSKLevel: 2},
	{Translation: "Из-за дождя мы не идём.", Chunks: []string{"因为", "下雨", "所以", "我们", "不", "去"}, Pinyin: "yīn wèi xià yǔ suǒ yǐ wǒ men bù qù", HSKLevel: 2},
	{Translation: "Ты хочешь пить чай или воду?", Chunks: []string{"你", "想", "喝", "茶", "或者", "水"}, Pinyin: "nǐ xiǎng hē chá huò zhě shuǐ", HSKLevel: 2},
	{Translation: "Он уже учитель.", Chunks: []string{"他", "已经", "是", "老师", "了"}, Pinyin: "tā yǐ jīng shì lǎo shī le", HSKLevel: 2},
	{Translation: "Ты можешь мне помочь?", Chunks: []string{"你", "能", "帮助", "我", "吗"}, Pinyin: "nǐ néng bāng zhù wǒ ma", HSKLevel: 2},
	{Translation: "Я думаю, это очень хорошо.", Chunks: []string{"我", "觉得", "很", "好"}, Pinyin: "wǒ jué de hěn hǎo", HSKLevel: 2},
	{Translation: "Мы можем начинать.", Chunks: []string{"我们", "可以", "开始", "了"}, Pinyin: "wǒ men kě yǐ kāi shǐ le", HSKLevel: 2},
	{Translation: "Мне нравится заниматься спортом.", Chunks: []string{"我", "喜欢", "运动"}, Pinyin: "wǒ xǐ huan yùn dòng", HSKLevel: 2},
	{Translation: "Китайский язык не трудный.", Chunks: []string{"汉语", "不", "难"}, Pinyin: "hàn yǔ bù nán", HSKLevel: 2},
	{Translation: "Это новое.", Chunks: []string{"这", "个", "是", "新", "的"}, Pinyin: "zhè gè shì xīn de", HSKLevel: 2},
	{Translation: "Ты слишком медленный.", Chunks: []string{"你", "太", "慢", "了"}, Pinyin: "nǐ tài màn le", HSKLevel: 2},
	{Translation: "Сейчас не очень удобно.", Chunks: []string{"现在", "不", "太", "方便"}, Pinyin: "xiàn zài bù tài fāng biàn", HSKLevel: 2},
	{Translation: "Этот человек очень умный.", Chunks: []string{"这", "个", "人", "很", "聪明"}, Pinyin: "zhè gè rén hěn cōng míng", HSKLevel: 2},
	{Translation: "Это слишком дорого.", Chunks: []string{"这", "个", "太", "贵", "了"}, Pinyin: "zhè gè tài guì le", HSKLevel: 2},
	{Translation: "Мы идём в библиотеку.", Chunks: []string{"我们", "去", "图书馆"}, Pinyin: "wǒ men qù tú shū guǎn", HSKLevel: 2},
	{Translation: "А где твой телефон?", Chunks: []string{"你", "的", "手机", "呢"}, Pinyin: "nǐ de shǒu jī ne", HSKLevel: 2},
	{Translation: "У меня есть младший брат и младшая сестра.", Chunks: []string{"我", "有", "一", "个", "弟弟", "和", "一", "个", "妹妹"}, Pinyin: "wǒ yǒu yī gè dì di hé yī gè mèi mei", HSKLevel: 2},
	{Translation: "На горе есть снег.", Chunks: []string{"山", "上", "有", "雪"}, Pinyin: "shān shàng yǒu xuě", HSKLevel: 2},
	{Translation: "Математика не простая.", Chunks: []string{"数学", "不", "容易"}, Pinyin: "shù xué bù róng yì", HSKLevel: 2},
}
