import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { TextItem } from '@/types/text'

export async function getTexts(hsk?: number) {
  const response = await api.get<ApiResponse<TextItem[]>>('/texts/', {
    params: { hsk: hsk || undefined },
  })

  return response.data
}

export async function getText(id: number) {
  const response = await api.get<ApiResponse<TextItem>>(`/texts/${id}`)

  return response.data
}

export async function markTextRead(id: number) {
  const response = await api.post<ApiResponse<{ status: string }>>(`/texts/${id}/read`)

  return response.data
}

export async function markTextUnread(id: number) {
  const response = await api.delete<ApiResponse<{ status: string }>>(`/texts/${id}/read`)

  return response.data
}
