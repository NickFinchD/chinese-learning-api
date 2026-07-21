export interface Course {
  id: number
  title: string
  description: string
  hsk_level: number
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Lesson {
  id: number
  title: string
  lesson_number: number
}

export interface CourseDetails {
  id: number
  title: string
  description: string
  hsk_level: number
  lessons: Lesson[]
}