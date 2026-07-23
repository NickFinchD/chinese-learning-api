import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { Word } from '@/types/word'

export async function getSavedWords() {
  const response = await api.get<ApiResponse<Word[]>>('/words/saved')

  return response.data
}

export async function saveWord(wordId: number) {
  const response = await api.post(`/words/${wordId}/save`)

  return response.data
}

export async function unsaveWord(wordId: number) {
  const response = await api.delete(`/words/${wordId}/save`)

  return response.data
}
