import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { Quiz } from '@/types/quiz'

export interface CheckAnswerResponse {
  correct: boolean
}

export async function checkAnswer(
  quizId: number,
  optionId: number,
): Promise<CheckAnswerResponse> {
  const { data } = await api.post<CheckAnswerResponse>(
    `/quizzes/${quizId}/check`,
    {
      option_id: optionId,
    },
  )

  return data
}

export async function getQuizzesByHSK(hsk: number) {
  const response = await api.get<ApiResponse<Quiz[]>>('/quizzes/', {
    params: { hsk },
  })

  return response.data
}