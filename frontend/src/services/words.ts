import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { Word } from '@/types/word'

export interface ListWordsParams {
  search?: string
  hsk?: number
}

export async function getWords(params: ListWordsParams = {}) {
  const response = await api.get<ApiResponse<Word[]>>('/words/', {
    params: {
      search: params.search || undefined,
      hsk: params.hsk || undefined,
    },
  })

  return response.data
}
