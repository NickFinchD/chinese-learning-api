import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { InProgressWord, WordProgress } from '@/types/learning'
import type { Word } from '@/types/word'

export async function getLearningProgress() {
  const response = await api.get<ApiResponse<WordProgress[]>>('/learning/')

  return response.data
}

export async function getLearnedWords() {
  const response = await api.get<ApiResponse<Word[]>>('/learning/learned')

  return response.data
}

export async function getInProgressWords() {
  const response = await api.get<ApiResponse<InProgressWord[]>>('/learning/in-progress')

  return response.data
}

export async function recordWordAnswer(wordId: number, correct: boolean) {
  const response = await api.post<ApiResponse<WordProgress>>(
    `/learning/${wordId}/answer`,
    { correct },
  )

  return response.data
}
