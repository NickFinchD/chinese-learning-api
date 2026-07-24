import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { MockExamAttempt, MockExamPaper, MockExamResult } from '@/types/mockExam'

export interface SubmitMockExamPayload {
  quiz_answers: { quiz_id: number, option_id: number }[]
  sentence_answers: { exercise_id: number, chunks: string[] }[]
  duration_seconds: number
}

export async function getMockExamPaper(hsk: number) {
  const response = await api.get<ApiResponse<MockExamPaper>>(`/mock-exams/${hsk}`)

  return response.data
}

export async function submitMockExam(hsk: number, payload: SubmitMockExamPayload) {
  const response = await api.post<ApiResponse<MockExamResult>>(`/mock-exams/${hsk}/submit`, payload)

  return response.data
}

export async function getMockExamHistory() {
  const response = await api.get<ApiResponse<MockExamAttempt[]>>('/mock-exams/history')

  return response.data
}
