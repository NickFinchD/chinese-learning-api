export interface Word {
  id: number
  hanzi: string
  pinyin: string
  translation: string
}

export interface QuizOption {
  id: number
  text: string
  sort_order: number
}

export interface Quiz {
  id: number
  question: string
  options: QuizOption[]
}

export interface GrammarNote {
  id: number
  title: string
  explanation: string
  example_hanzi?: string
  example_pinyin?: string
  example_translation?: string
  hsk_level: number
}

export interface SentenceExercise {
  id: number
  translation: string
  chunks: string[]
  pinyin: string
  hsk_level: number
}

export type LessonStep =
  | {
      id: number
      step_type: 'word'
      sort_order: number
      data: Word
    }
  | {
      id: number
      step_type: 'quiz'
      sort_order: number
      data: Quiz
    }
  | {
      id: number
      step_type: 'grammar'
      sort_order: number
      data: GrammarNote
    }
  | {
      id: number
      step_type: 'sentence_builder'
      sort_order: number
      data: SentenceExercise
    }

export interface Lesson {
  id: number
  course_id: number
  title: string
  description: string
  lesson_number: number
  steps: LessonStep[]
}