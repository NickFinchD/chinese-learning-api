import { api } from '@/services/client'

import type { ApiResponse } from '@/types/api'
import type { AdminListParams, PagedResult } from '@/types/admin'
import type { Word, WordDetail } from '@/types/word'

export interface ExamplePayload {
  hanzi: string
  pinyin: string
  translation: string
}

export interface WordPayload {
  hanzi: string
  pinyin: string
  translation: string
  part_of_speech: string
  hsk_level: number
  examples: ExamplePayload[]
}

export async function adminListWords(params: AdminListParams) {
  const response = await api.get<ApiResponse<PagedResult<Word>>>('/admin/words', { params })

  return response.data
}

export async function createWord(payload: WordPayload) {
  const response = await api.post<ApiResponse<WordDetail>>('/admin/words', payload)

  return response.data
}

export async function updateWord(id: number, payload: WordPayload) {
  const response = await api.put<ApiResponse<WordDetail>>(`/admin/words/${id}`, payload)

  return response.data
}

export async function deleteWord(id: number) {
  await api.delete(`/admin/words/${id}`)
}
