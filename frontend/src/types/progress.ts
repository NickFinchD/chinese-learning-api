export type LessonProgressStatus = 'not_started' | 'in_progress' | 'completed'

export interface LessonProgress {
  status: LessonProgressStatus
  current_step: number
  score: number
}
