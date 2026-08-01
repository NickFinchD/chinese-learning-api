import { api } from '@/services/client'

import type { ApiResponse } from '@/types/api'
import type { AdminListParams, PagedResult } from '@/types/admin'
import type { TextItem } from '@/types/text'

export interface TextPayload {
  title: string
  hanzi: string
  pinyin: string
  translation: string
  hsk_level: number
  image_url: string
}

export async function adminListTexts(params: AdminListParams) {
  const response = await api.get<ApiResponse<PagedResult<TextItem>>>('/admin/texts', { params })

  return response.data
}

export async function createText(payload: TextPayload) {
  const response = await api.post<ApiResponse<TextItem>>('/admin/texts', payload)

  return response.data
}

export async function updateText(id: number, payload: TextPayload) {
  const response = await api.put<ApiResponse<TextItem>>(`/admin/texts/${id}`, payload)

  return response.data
}

export async function deleteText(id: number) {
  await api.delete(`/admin/texts/${id}`)
}

export async function uploadTextImage(file: File) {
  const formData = new FormData()
  formData.append('image', file)

  const response = await api.post<ApiResponse<{ url: string }>>('/admin/texts/upload-image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })

  return response.data
}
