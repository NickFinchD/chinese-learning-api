export interface PagedResult<T> {
  items: T[]
  total: number
  page: number
  limit: number
}

export interface AdminListParams {
  search?: string
  page?: number
  limit?: number
}
