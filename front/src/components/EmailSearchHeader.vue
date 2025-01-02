<template>
  <header class="fixed inset-x-0 top-0 z-50 bg-white/80 backdrop-blur-sm">
    <div class="w-full mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header Content -->
      <div class="flex flex-col sm:flex-row items-center justify-between py-4 gap-4">
        <!-- Logo/Title -->
        <div class="flex-shrink-0">
          <h1 class="text-xl sm:text-2xl font-bold text-indigo-600">Enron Email Search</h1>
        </div>

        <!-- Search Bar -->
        <div class="w-full sm:max-w-2xl">
          <div class="relative">
            <input
              v-model="searchQuery"
              @keyup.enter="handleSearch"
              type="text"
              class="w-full rounded-lg border border-gray-300 py-2 pl-4 pr-20 text-gray-900 placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:border-transparent"
              placeholder="Search emails..."
            />
            <button
              @click="handleSearch"
              class="absolute right-2 top-1/2 -translate-y-1/2 px-4 py-1 rounded-md bg-indigo-600 text-white text-sm hover:bg-indigo-700 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-600"
            >
              Search
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>

  <!-- Spacer to prevent content from hiding behind fixed header -->
  <div class="h-24 sm:h-20"></div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const searchQuery = ref('')
const emit = defineEmits(['search'])

const handleSearch = () => {
  emit('search', searchQuery.value)
}
</script>
