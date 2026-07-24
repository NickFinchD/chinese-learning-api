import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { Collection, CollectionDetail } from '@/types/collection'

export async function getCollections() {
  const response = await api.get<ApiResponse<Collection[]>>('/collections/')

  return response.data
}

export async function getCuratedCollections() {
  const response = await api.get<ApiResponse<Collection[]>>('/collections/curated')

  return response.data
}

export async function saveCuratedCollection(id: number) {
  const response = await api.post<ApiResponse<Collection>>(`/collections/${id}/save`)

  return response.data
}

export async function getCollection(id: number) {
  const response = await api.get<ApiResponse<CollectionDetail>>(`/collections/${id}`)

  return response.data
}

export async function createCollection(name: string) {
  const response = await api.post<ApiResponse<Collection>>('/collections/', { name })

  return response.data
}

export async function renameCollection(id: number, name: string) {
  const response = await api.patch(`/collections/${id}`, { name })

  return response.data
}

export async function deleteCollection(id: number) {
  const response = await api.delete(`/collections/${id}`)

  return response.data
}

export async function addWordToCollection(collectionId: number, wordId: number) {
  const response = await api.post(`/collections/${collectionId}/words/${wordId}`)

  return response.data
}

export async function removeWordFromCollection(collectionId: number, wordId: number) {
  const response = await api.delete(`/collections/${collectionId}/words/${wordId}`)

  return response.data
}
