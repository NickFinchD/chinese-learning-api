import { api } from './client'

export interface LoginRequest {
  email: string
  password: string
}

export interface User {
  id: number
  username: string
  email: string
}

export interface LoginResponse {
  success: boolean
  data: {
    token: string
    user: User
  }
}

export async function login(request: LoginRequest) {
  const response = await api.post<LoginResponse>(
    '/auth/login',
    request,
  )

  return response.data
}