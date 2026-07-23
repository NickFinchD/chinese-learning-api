import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { Achievement, GamificationProgress } from '@/types/gamification'

export async function sendHeartbeat() {
  const response = await api.post<ApiResponse<GamificationProgress>>('/gamification/heartbeat')

  return response.data
}

export async function getProgress() {
  const response = await api.get<ApiResponse<GamificationProgress>>('/gamification/progress')

  return response.data
}

export async function getAchievements() {
  const response = await api.get<ApiResponse<Achievement[]>>('/gamification/achievements')

  return response.data
}
