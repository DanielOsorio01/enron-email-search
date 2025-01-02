<template>
  <TransitionRoot as="template" :show="open" @keydown.escape="closeDialog">
    <Dialog class="relative z-50" @close="closeDialog">
      <TransitionChild
        as="template"
        enter="ease-out duration-300"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="ease-in duration-200"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-gray-500/75 transition-opacity" @click="closeDialog" />
      </TransitionChild>

      <div class="fixed inset-0 z-10 w-screen overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4 text-center sm:p-0">
          <TransitionChild
            as="template"
            enter="ease-out duration-300"
            enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            enter-to="opacity-100 translate-y-0 sm:scale-100"
            leave="ease-in duration-200"
            leave-from="opacity-100 translate-y-0 sm:scale-100"
            leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          >
            <DialogPanel
              class="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all sm:my-8 w-full sm:w-full sm:max-w-screen-lg"
              @click.stop
            >
              <!-- Modal Content -->
              <div class="bg-white px-4 py-4 sm:px-6 sm:py-6 w-full">
                <div class="flex justify-between items-center border-b pb-3 w-full">
                  <DialogTitle
                    as="h2"
                    class="text-lg sm:text-xl font-semibold text-gray-900 w-full"
                  >
                    Email Details
                  </DialogTitle>
                  <button class="text-gray-500 hover:text-gray-700" @click="closeDialog">✕</button>
                </div>
                <!-- Email Details -->
                <div class="mt-4 overflow-y-auto max-h-full w-full">
                  <p class="text-sm sm:text-base text-gray-600">
                    <strong>From:</strong> {{ email.from }}
                  </p>
                  <p class="text-sm sm:text-base text-gray-600">
                    <strong>To:</strong> {{ email.to }}
                  </p>
                  <p class="text-sm sm:text-base text-gray-600">
                    <strong>Subject:</strong> {{ email.subject }}
                  </p>
                  <!-- Body Section -->
                  <div class="mt-4 text-sm sm:text-base text-gray-700 whitespace-pre-wrap w-full">
                    <template
                      v-for="(part, index) in highlightedBody(email.body, searchQuery)"
                      :key="index"
                    >
                      <span v-if="part.highlight" class="bg-yellow-200">{{ part.text }}</span>
                      <span v-else>{{ part.text }}</span>
                    </template>
                  </div>
                </div>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup lang="ts">
import { ref, type PropType } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'

import { type Email } from '../types/Email'

// Props
const props = defineProps({
  email: {
    type: Object as PropType<Email>,
    required: true,
  },
  searchQuery: {
    type: String,
    required: true,
  },
  handleClose: {
    type: Function,
    required: true,
  },
})

// Highlight search query in the email body
const highlightedBody = (
  body: string,
  searchQuery: string,
): { text: string; highlight: boolean }[] => {
  if (!searchQuery) {
    return [{ text: body, highlight: false }]
  }

  const query = searchQuery.toLowerCase()
  const index = body.toLowerCase().indexOf(query)
  const regex = new RegExp(`(${query})`, 'gi')
  const parts = body.split(regex)

  if (index === -1) {
    return [{ text: body, highlight: false }]
  }
  return parts.map((part) => ({
    text: part,
    highlight: query.toLowerCase() === part.toLowerCase(),
  }))
}

// Modal state and close handler
const open = ref(true)
const closeDialog = () => {
  open.value = false
  props.handleClose()
}
</script>
