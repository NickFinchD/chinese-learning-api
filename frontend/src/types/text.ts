export type TextStatus = 'not_started' | 'in_progress' | 'completed'

export interface TextItem {
  id: number
  title: string
  hanzi: string
  pinyin: string
  translation: string
  hsk_level: number
  image_url: string
  status: TextStatus
}
