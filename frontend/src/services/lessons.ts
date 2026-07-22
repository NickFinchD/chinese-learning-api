import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { Lesson } from '@/types/lesson'

export async function getLesson(id: number) {
  const response = await api.get<ApiResponse<Lesson>>(`/lessons/${id}`)

  return response.data
}