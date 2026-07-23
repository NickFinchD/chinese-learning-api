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

export interface Lesson {
  id: number
  course_id: number
  title: string
  description: string
  lesson_number: number
  steps: LessonStep[]
}