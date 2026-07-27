import { api } from '@/services/client'

import type { ApiResponse } from '@/types/api'
import type { AdminListParams, PagedResult } from '@/types/admin'

export interface AdminUser {
  id: number
  username: string
  email: string
  is_admin: boolean
  created_at: string
}

export async function adminListUsers(params: AdminListParams) {
  const response = await api.get<ApiResponse<PagedResult<AdminUser>>>('/admin/users', { params })

  return response.data
}

export async function setUserAdmin(id: number, isAdmin: boolean) {
  const response = await api.patch<ApiResponse<AdminUser>>(`/admin/users/${id}/admin`, { is_admin: isAdmin })

  return response.data
}
