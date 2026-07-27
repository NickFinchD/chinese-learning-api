import { api } from '@/services/client'

import type { ApiResponse } from '@/types/api'
import type { AdminListParams, PagedResult } from '@/types/admin'
import type { GrammarNote } from '@/types/lesson'

export interface GrammarNotePayload {
  title: string
  explanation: string
  example_hanzi: string
  example_pinyin: string
  example_translation: string
  example2_hanzi: string
  example2_pinyin: string
  example2_translation: string
  example3_hanzi: string
  example3_pinyin: string
  example3_translation: string
  hsk_level: number
}

export async function adminListGrammarNotes(params: AdminListParams) {
  const response = await api.get<ApiResponse<PagedResult<GrammarNote>>>('/admin/grammar', { params })

  return response.data
}

export async function createGrammarNote(payload: GrammarNotePayload) {
  const response = await api.post<ApiResponse<GrammarNote>>('/admin/grammar', payload)

  return response.data
}

export async function updateGrammarNote(id: number, payload: GrammarNotePayload) {
  const response = await api.put<ApiResponse<GrammarNote>>(`/admin/grammar/${id}`, payload)

  return response.data
}

export async function deleteGrammarNote(id: number) {
  await api.delete(`/admin/grammar/${id}`)
}
