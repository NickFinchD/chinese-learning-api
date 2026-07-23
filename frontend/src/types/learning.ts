export interface WordProgress {
  word_id: number
  stage: number
  max_stage: number
  next_eligible_at?: string
  learned: boolean
}

export interface InProgressWord {
  word_id: number
  hanzi: string
  pinyin: string
  translation: string
  hsk_level: number
  stage: number
  max_stage: number
  repetitions_left: number
  next_eligible_at?: string
}
