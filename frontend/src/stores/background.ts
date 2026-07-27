import { acceptHMRUpdate, defineStore } from 'pinia'

import mountainsDay from '@/assets/backgrounds/mountains-day.png'
import mountainsNight from '@/assets/backgrounds/mountains-night.png'
import pagodaDay from '@/assets/backgrounds/pagoda-day.png'
import pagodaNight from '@/assets/backgrounds/pagoda-night.png'
import koiDay from '@/assets/backgrounds/koi-day.png'
import koiNight from '@/assets/backgrounds/koi-night.png'

export type BackgroundId = 'mountains' | 'pagoda' | 'koi'

export const BACKGROUNDS: Record<BackgroundId, { label: string, day: string, night: string }> = {
  mountains: {
    label: 'Горы',
    day: mountainsDay,
    night: mountainsNight,
  },
  pagoda: {
    label: 'Пагода',
    day: pagodaDay,
    night: pagodaNight,
  },
  koi: {
    label: 'Кои',
    day: koiDay,
    night: koiNight,
  },
}

const STORAGE_KEY = 'wojiao-background'

export const useBackgroundStore = defineStore('background', {
  state: () => ({
    backgroundId: (localStorage.getItem(STORAGE_KEY) as BackgroundId | null) ?? 'mountains',
  }),

  actions: {
    setBackground(id: BackgroundId) {
      this.backgroundId = id
      localStorage.setItem(STORAGE_KEY, id)
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useBackgroundStore, import.meta.hot))
}
