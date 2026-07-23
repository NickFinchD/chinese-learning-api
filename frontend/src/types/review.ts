export interface ReviewWord {
  word_id: number
  hanzi: string
  pinyin: string
  translation: string
}

export interface ReviewStatistics {
  total_words: number
  ready_for_review: number
  reviewed_words: number
}

export interface ReviewSession {
  total: number
  words: ReviewWord[]
}
