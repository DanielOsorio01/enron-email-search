<template>
  <EmailSearchHeader @search="handleSearch" />

  <main class="container mx-auto px-4 sm:px-6 lg:px-8">
    <div v-if="errorMessage" class="mt-4 text-red-500">
      {{ errorMessage }}
    </div>

    <!-- Welcome State -->
    <div v-if="!hasSearched" class="mt-20 text-center">
      <h1 class="text-3xl font-bold text-gray-900 mb-4">Welcome to Enron Email Search</h1>
      <p class="text-gray-600 text-lg">
        Start by typing a search term above to explore the email database.
      </p>
    </div>

    <!-- No Results State -->
    <div v-else-if="hasSearched && hits === 0" class="mt-20 text-center">
      <h2 class="text-xl font-semibold text-gray-900 mb-2">Oops! No results found</h2>
      <p class="text-gray-600">Try searching with another term</p>
    </div>

    <!-- Results State -->
    <template v-else-if="hits > 0">
      <div class="mt-6 mb-4">
        <h2 class="text-xl font-bold text-gray-900">
          Results <span class="text-gray-500">({{ hits }})</span>
        </h2>
      </div>
      <EmailList :emails="emails" :search-query="query" />
    </template>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Email } from './types/Email'
import EmailList from './components/EmailList.vue'
import EmailSearchHeader from './components/EmailSearchHeader.vue'

const query = ref('')
const emails = ref<Email[]>([])
const errorMessage = ref('')
const hits = ref(0)
const hasSearched = ref(false)

const handleSearch = async (searchQuery: string) => {
  hasSearched.value = true
  try {
    const response = await fetch(`http://localhost:3000/emails?term=${searchQuery}`)
    if (!response.ok) {
      throw new Error(`Failed to fetch emails: ${response.statusText}`)
    }
    const jsonObj = await response.json()
    if (jsonObj.success) {
      query.value = searchQuery
      errorMessage.value = ''
      emails.value = jsonObj.data.emails
      hits.value = jsonObj.data.total
    } else {
      errorMessage.value = jsonObj.message
      hits.value = 0
    }
  } catch (error) {
    errorMessage.value = 'Failed to connect to the server. Please try again later.'
    hits.value = 0
    console.error(error)
  }
}
</script>
