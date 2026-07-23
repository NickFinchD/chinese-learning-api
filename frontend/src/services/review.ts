import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { ReviewSession, ReviewStatistics } from '@/types/review'

export async function getReviewStatistics() {
  const response = await api.get<ApiResponse<ReviewStatistics>>('/reviews/statistics')

  return response.data
}

export async function startReviewSession() {
  const response = await api.get<ApiResponse<ReviewSession>>('/reviews/session')

  return response.data
}

export async function answerReview(wordId: number, correct: boolean) {
  const response = await api.post('/reviews/answer', {
    word_id: wordId,
    correct,
  })

  return response.data
}
