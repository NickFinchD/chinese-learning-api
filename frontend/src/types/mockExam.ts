import type { SentenceExercise } from './lesson'
import type { Quiz } from './quiz'

export type MockExamQuestion =
  | { type: 'quiz'; quiz: Quiz }
  | { type: 'sentence'; sentence: SentenceExercise }

export interface MockExamPaper {
  hsk_level: number
  time_limit_seconds: number
  questions: MockExamQuestion[]
}

export interface MockExamAttempt {
  id: number
  hsk_level: number
  total_questions: number
  correct_count: number
  score_percent: number
  passed: boolean
  duration_seconds: number
  created_at: string
}

export interface MockExamResult {
  attempt: MockExamAttempt
  xp_awarded: number
}
