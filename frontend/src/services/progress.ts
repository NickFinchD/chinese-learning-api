import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { LessonProgress } from '@/types/progress'

export async function startLesson(lessonId: number) {
  const response = await api.post<ApiResponse<{ status: string }>>(`/lessons/${lessonId}/start`)

  return response.data
}

export async function getLessonProgress(lessonId: number) {
  const response = await api.get<ApiResponse<LessonProgress>>(`/lessons/${lessonId}/progress`)

  return response.data
}

export async function updateLessonStep(lessonId: number, currentStep: number) {
  const response = await api.post<ApiResponse<{ status: string; current_step: number }>>(
    `/lessons/${lessonId}/step`,
    { current_step: currentStep },
  )

  return response.data
}

export async function completeLesson(lessonId: number, score: number) {
  const response = await api.post<ApiResponse<{ status: string; score: number }>>(
    `/lessons/${lessonId}/complete`,
    { score },
  )

  return response.data
}
