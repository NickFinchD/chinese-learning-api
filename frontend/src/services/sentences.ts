import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { SentenceExercise } from '@/types/lesson'

export async function getSentencesByHSK(hsk: number) {
  const response = await api.get<ApiResponse<SentenceExercise[]>>('/sentences/', {
    params: { hsk },
  })

  return response.data
}
