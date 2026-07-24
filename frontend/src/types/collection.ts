import type { Word } from './word'

export interface Collection {
  id: number
  name: string
  word_count: number
  created_at: string
  updated_at: string
}

export interface CollectionDetail extends Collection {
  words: Word[]
}
