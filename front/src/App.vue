<script setup lang="ts">
import { ref } from 'vue'
import type { SearchResponse } from './types/email'
import { searchEmails } from './services/api'

const searchTerm = ref('')
const searchResponse = ref<SearchResponse | null>(null)
const isLoading = ref(false)
const error = ref('')

const handleSearch = async () => {
  if (!searchTerm.value.trim()) return

  isLoading.value = true
  error.value = ''

  try {
    const response = await searchEmails(searchTerm.value)
    searchResponse.value = response
  } catch (e) {
    error.value = 'Failed to search emails' + (e.message ? `: ${e.message}` : '')
    searchResponse.value = null
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-100 py-8 px-4">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-3xl font-bold text-gray-900 mb-8">Email Search</h1>

      <!-- Search Form -->
      <div class="mb-6">
        <div class="flex gap-4">
          <input
            v-model="searchTerm"
            type="text"
            placeholder="Search emails..."
            class="flex-1 p-3 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            @keyup.enter="handleSearch"
          />
          <button
            @click="handleSearch"
            class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:outline-none disabled:opacity-50"
            :disabled="isLoading"
          >
            {{ isLoading ? 'Searching...' : 'Search' }}
          </button>
        </div>

        <p v-if="searchResponse?.data.total" class="mt-2 text-sm text-gray-600">
          Found {{ searchResponse.data.total }} results
        </p>
      </div>

      <!-- Error Message -->
      <div v-if="error" class="mb-6 p-4 bg-red-100 text-red-700 rounded-lg">
        {{ error }}
      </div>

      <!-- Results -->
      <div v-if="searchResponse?.data.emails.length" class="space-y-4">
        <div
          v-for="email in searchResponse.data.emails"
          :key="email.id"
          class="bg-white p-4 rounded-lg shadow"
        >
          <div class="flex justify-between items-start mb-2">
            <div>
              <h2 class="font-semibold text-lg">{{ email.subject }}</h2>
              <p class="text-sm text-gray-600">From: {{ email.from }}</p>
              <p class="text-sm text-gray-600">To: {{ email.to }}</p>
            </div>
            <span class="text-sm text-gray-500">
              {{ new Date(email.date).toLocaleDateString() }}
            </span>
          </div>
          <p class="text-gray-700 whitespace-pre-line">
            {{ email.body }}
          </p>
        </div>
      </div>

      <!-- No Results -->
      <div
        v-else-if="!isLoading && searchTerm && searchResponse"
        class="text-center py-8 text-gray-600"
      >
        No emails found matching your search.
      </div>
    </div>
  </div>
</template>
