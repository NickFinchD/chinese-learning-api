export interface GamificationProgress {
  xp: number
  level: number
  current_level_xp: number
  xp_for_next_level: number
  total_seconds_active: number
  hours_active: number
}

export interface Achievement {
  code: string
  title: string
  description: string
  tier: number
  metric: string
  threshold: number
  current_value: number
  xp_reward: number
  unlocked: boolean
  unlocked_at?: string
}
