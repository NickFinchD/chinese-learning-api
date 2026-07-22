import { api } from './client'

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