import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { ProperName } from '@/types/name'

export async function getNames() {
  const response = await api.get<ApiResponse<ProperName[]>>('/names/')

  return response.data
}
