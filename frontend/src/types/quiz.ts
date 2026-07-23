export interface QuizOption {
  id: number
  text: string
}

export interface Quiz {
  id: number
  question: string
  hsk_level: number
  options: QuizOption[]
}
