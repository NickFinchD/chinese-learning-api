<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Достижения
    </h1>

    <div
      v-if="gamification.progress"
      class="mb-8 max-w-md rounded-2xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
    >
      <div class="mb-2 flex items-center justify-between">
        <span class="text-lg font-semibold text-gray-900 dark:text-white">Уровень {{ gamification.progress.level }}</span>
        <span class="text-sm text-gray-500 dark:text-gray-400">{{ gamification.progress.hours_active }} ч. в приложении</span>
      </div>

      <div class="mb-1 h-3 overflow-hidden rounded-full bg-gray-200/50 dark:bg-white/10">
        <div
          class="h-full bg-[var(--color-primary)] transition-all duration-300"
          :style="{ width: `${xpPercent}%` }"
        />
      </div>

      <div class="text-sm text-gray-500 dark:text-gray-400">
        {{ gamification.progress.current_level_xp }} / {{ gamification.progress.xp_for_next_level }} XP до следующего уровня
      </div>
    </div>

    <div
      v-if="gamification.loading"
      class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
    >
      <BaseSpinner />
      Загрузка...
    </div>

    <div
      v-else
      class="space-y-8"
    >
      <div
        v-for="group in groupedAchievements"
        :key="group.metric"
      >
        <h2 class="mb-4 text-xl font-semibold text-gray-900 dark:text-white">
          {{ group.label }}
        </h2>

        <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          <div
            v-for="(achievement, index) in group.items"
            :key="achievement.code"
            class="animate-fade-in-up rounded-xl border p-6 shadow-sm backdrop-blur-xl transition hover:-translate-y-0.5 hover:shadow-md"
            :class="achievement.unlocked
              ? 'border-[var(--color-accent)]/50 bg-[var(--color-accent)]/10 dark:border-[var(--color-accent)]/30 dark:bg-[var(--color-accent)]/10'
              : 'border-white/50 bg-white/20 opacity-60 dark:border-white/10 dark:bg-white/5'"
            :style="{ animationDelay: `${Math.min(index * 40, 300)}ms` }"
          >
            <div class="mb-2 flex items-center justify-between gap-2">
              <h3 class="font-semibold text-gray-900 dark:text-white">
                {{ achievement.title }}
              </h3>

              <AppIcon
                :name="achievement.unlocked ? 'trophy' : 'lock'"
                :size="22"
                :filled="achievement.unlocked"
                :class="achievement.unlocked ? 'text-[var(--color-accent)]' : 'text-gray-400 dark:text-gray-500'"
              />
            </div>

            <p class="mb-3 text-sm text-gray-600 dark:text-gray-400">
              {{ achievement.description }}
            </p>

            <div
              v-if="!achievement.unlocked"
              class="mb-3"
            >
              <div class="mb-1 h-1.5 overflow-hidden rounded-full bg-gray-200/50 dark:bg-white/10">
                <div
                  class="h-full bg-[var(--color-primary)] transition-all duration-300"
                  :style="{ width: `${progressPercent(achievement)}%` }"
                />
              </div>

              <div class="text-xs text-gray-500 dark:text-gray-400">
                {{ Math.min(achievement.current_value, achievement.threshold) }} / {{ achievement.threshold }}
              </div>
            </div>

            <div class="text-xs text-gray-500 dark:text-gray-400">
              <span v-if="achievement.unlocked">
                Получено {{ formatDate(achievement.unlocked_at) }} · +{{ achievement.xp_reward }} XP
              </span>
              <span v-else>
                Награда: +{{ achievement.xp_reward }} XP
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'

import { useGamificationStore } from '@/stores/gamification'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

import type { Achievement } from '@/types/gamification'

const gamification = useGamificationStore()

const metricLabels: Record<string, string> = {
  hours_active: 'Усидчивость',
  words_learned: 'Эрудит',
  lessons_completed: 'Отличник',
}

const xpPercent = computed(() => {
  if (!gamification.progress || gamification.progress.xp_for_next_level === 0) {
    return 0
  }

  return Math.min(100, (gamification.progress.current_level_xp / gamification.progress.xp_for_next_level) * 100)
})

const groupedAchievements = computed(() => {
  const groups = new Map<string, Achievement[]>()

  for (const achievement of gamification.achievements) {
    const list = groups.get(achievement.metric) ?? []
    list.push(achievement)
    groups.set(achievement.metric, list)
  }

  return Array.from(groups.entries()).map(([metric, items]) => ({
    metric,
    label: metricLabels[metric] ?? metric,
    items: items.sort((a, b) => a.tier - b.tier),
  }))
})

function progressPercent(achievement: Achievement) {
  if (achievement.threshold <= 0) {
    return 0
  }

  return Math.min(100, (achievement.current_value / achievement.threshold) * 100)
}

function formatDate(iso?: string) {
  if (!iso) {
    return ''
  }

  return new Date(iso).toLocaleDateString('ru-RU')
}

onMounted(() => {
  gamification.loadProgress()
  gamification.loadAchievements()
})
</script>
