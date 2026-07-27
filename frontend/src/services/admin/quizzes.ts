import { api } from '@/services/client'

import type { ApiResponse } from '@/types/api'
import type { AdminListParams, PagedResult } from '@/types/admin'
import type { QuizDirection } from '@/types/lesson'

export interface AdminQuizOption {
  id?: number
  text: string
  pinyin: string
  is_correct: boolean
  sort_order?: number
}

export interface AdminQuiz {
  id: number
  question: string
  hsk_level: number
  direction: QuizDirection
  hanzi: string
  pinyin: string
  options: AdminQuizOption[]
  created_at: string
  updated_at: string
}

export interface QuizPayload {
  question: string
  hsk_level: number
  direction: QuizDirection
  hanzi: string
  pinyin: string
  options: { text: string, pinyin: string, is_correct: boolean }[]
}

export async function adminListQuizzes(params: AdminListParams) {
  const response = await api.get<ApiResponse<PagedResult<AdminQuiz>>>('/admin/quizzes', { params })

  return response.data
}

export async function createQuiz(payload: QuizPayload) {
  const response = await api.post<ApiResponse<AdminQuiz>>('/admin/quizzes', payload)

  return response.data
}

export async function updateQuiz(id: number, payload: QuizPayload) {
  const response = await api.put<ApiResponse<AdminQuiz>>(`/admin/quizzes/${id}`, payload)

  return response.data
}

export async function deleteQuiz(id: number) {
  await api.delete(`/admin/quizzes/${id}`)
}
