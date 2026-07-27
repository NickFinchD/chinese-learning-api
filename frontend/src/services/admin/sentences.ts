import { api } from '@/services/client'

import type { ApiResponse } from '@/types/api'
import type { AdminListParams, PagedResult } from '@/types/admin'
import type { SentenceExercise } from '@/types/lesson'

export interface SentenceExercisePayload {
  translation: string
  chunks: string[]
  pinyin: string
  hsk_level: number
}

export async function adminListSentences(params: AdminListParams) {
  const response = await api.get<ApiResponse<PagedResult<SentenceExercise>>>('/admin/sentences', { params })

  return response.data
}

export async function createSentence(payload: SentenceExercisePayload) {
  const response = await api.post<ApiResponse<SentenceExercise>>('/admin/sentences', payload)

  return response.data
}

export async function updateSentence(id: number, payload: SentenceExercisePayload) {
  const response = await api.put<ApiResponse<SentenceExercise>>(`/admin/sentences/${id}`, payload)

  return response.data
}

export async function deleteSentence(id: number) {
  await api.delete(`/admin/sentences/${id}`)
}
