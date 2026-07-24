import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { User } from '@/types/auth'

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export async function login(request: LoginRequest) {
  const response = await api.post<ApiResponse<User>>(
    '/auth/login',
    request,
  )

  return response.data
}

export async function register(request: RegisterRequest) {
  const response = await api.post<ApiResponse<User>>(
    '/auth/register',
    request,
  )

  return response.data
}

export async function logout() {
  await api.post('/auth/logout')
}