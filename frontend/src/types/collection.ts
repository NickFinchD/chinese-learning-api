import type { Word } from './word'

export interface Collection {
  id: number
  name: string
  word_count: number
  is_curated: boolean
  source_collection_id: number | null
  created_at: string
  updated_at: string
}

export interface CollectionDetail extends Collection {
  words: Word[]
}
