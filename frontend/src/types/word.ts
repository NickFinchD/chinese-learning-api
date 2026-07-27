export interface Word {
  id: number
  hanzi: string
  pinyin: string
  translation: string
  part_of_speech: string
  hsk_level: number
}

export interface WordExample {
  id: number
  hanzi: string
  pinyin: string
  translation: string
}

export interface WordDetail extends Word {
  examples?: WordExample[]
}
