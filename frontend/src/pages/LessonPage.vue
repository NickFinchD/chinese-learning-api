<template>
  <div v-if="lessons.loading">
    Loading...
  </div>

  <div v-else-if="lessons.current">
    <h1 class="mb-2 text-3xl font-bold">
      {{ lessons.current.title }}
    </h1>

    <p class="mb-2">
      {{ lessons.current.description }}
    </p>

    <div class="mb-2 text-sm text-gray-500">
      Step {{ currentStepIndex + 1 }} of {{ lessons.current.steps.length }}
    </div>

    <div class="mb-6 h-2 overflow-hidden rounded-full bg-gray-200">
      <div
        class="h-full bg-blue-600 transition-all duration-300"
        :style="{
          width: `${((currentStepIndex + 1) / lessons.current.steps.length) * 100}%`
        }"
      />
    </div>

    <div class="mb-6 flex justify-between">
      <button
        class="rounded bg-gray-200 px-4 py-2 disabled:opacity-50"
        :disabled="currentStepIndex === 0"
        @click="previousStep"
      >
        Previous
      </button>

      <button
        class="rounded bg-blue-600 px-4 py-2 text-white"
        @click="nextStep"
      >
        {{ isLastStep ? 'Finish lesson' : 'Next' }}
      </button>
    </div>

    <div
  v-if="currentStep"
  class="mb-4"
>
  <LessonStepRenderer
    :step="currentStep"
  />
</div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'


import { useLessonsStore } from '@/stores/lessons'
import LessonStepRenderer from '@/components/lesson/LessonStepRenderer.vue'


const lessons = useLessonsStore()
const route = useRoute()

const currentStepIndex = ref(0)

const currentStep = computed(() => {
  return lessons.current?.steps[currentStepIndex.value] ?? null
})

const isLastStep = computed(() => {
  if (!lessons.current) {
    return false
  }

  return currentStepIndex.value === lessons.current.steps.length - 1
})

function nextStep() {
  if (!lessons.current) {
    return
  }

  if (currentStepIndex.value < lessons.current.steps.length - 1) {
    currentStepIndex.value++
    return
  }

  console.log('Lesson finished')
}

function previousStep() {
  if (currentStepIndex.value > 0) {
    currentStepIndex.value--
  }
}

onMounted(() => {
  lessons.loadLesson(Number(route.params.id))
})
</script>