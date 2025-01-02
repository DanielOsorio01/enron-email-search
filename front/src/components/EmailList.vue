<template>
  <div class="container mx-auto p-4">
    <div class="grid gap-4">
      <div
        v-for="(email, index) in emails"
        :key="index"
        class="bg-white shadow rounded-lg p-4 hover:shadow-lg transition-shadow"
        @click="showEmailDetail(email)"
      >
        <h2 class="text-xl font-semibold text-gray-800 mb-2">{{ email.subject }}</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm text-gray-600">
          <div><span class="font-medium">From:</span> {{ email.from }}</div>
          <div><span class="font-medium">To:</span> {{ email.to }}</div>
        </div>
        <p class="mt-3 text-gray-700">
          <template v-for="(part, index) in highlightedBody(email.body, searchQuery)" :key="index">
            <span v-if="part.highlight" class="bg-yellow-200">{{ part.text }}</span>
            <span v-else>{{ part.text }}</span>
          </template>
        </p>
      </div>
    </div>

    <!-- Email Detail Dialog -->
    <EmailDetail
      v-if="emailOpened && selectedEmail"
      :email="selectedEmail"
      :search-query="searchQuery"
      :handle-close="closeEmailDetail"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, type PropType } from 'vue'
import type { Email } from '../types/Email'
import EmailDetail from './EmailDetail.vue'

export default defineComponent({
  props: {
    emails: {
      type: Array as PropType<Email[]>,
      required: true,
    },
    searchQuery: {
      type: String,
      required: true,
    },
  },
  components: { EmailDetail },
  setup() {
    const highlightedBody = (
      body: string,
      searchQuery: string,
    ): { text: string; highlight: boolean }[] => {
      if (!searchQuery) {
        return [{ text: body.slice(0, 100), highlight: false }]
      }

      const query = searchQuery.toLowerCase()
      const index = body.toLowerCase().indexOf(query)

      if (index === -1) {
        return [{ text: body.slice(0, 100), highlight: false }]
      }

      const start = Math.max(0, index - 50)
      const end = Math.min(body.length, index + query.length + 50)
      const excerpt = body.slice(start, end)

      const regex = new RegExp(`(${query})`, 'gi')
      const parts = excerpt.split(regex)

      const result = parts.map((part) => ({
        text: part,
        highlight: part.toLowerCase() === query,
      }))

      // Create a new array for the ellipsis at the start if needed
      const startEllipsis = start > 0 ? [{ text: '[...] ', highlight: false }] : []
      // Add ellipsis at the end if needed
      const endEllipsis = end < body.length ? [{ text: ' [...]', highlight: false }] : []

      // Concatenate the arrays
      return startEllipsis.concat(result, endEllipsis)
    }

    const emailOpened = ref(false)
    const selectedEmail = ref<Email | null>(null)
    const showEmailDetail = (email: Email) => {
      selectedEmail.value = email
      emailOpened.value = true
    }
    const closeEmailDetail = () => {
      emailOpened.value = false
      selectedEmail.value = null
    }
    return { highlightedBody, emailOpened, selectedEmail, showEmailDetail, closeEmailDetail }
  },
})
</script>
