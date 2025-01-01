<script setup lang="ts">
import { ref } from 'vue'
import { Email } from './types/Email'

const search = ref('')
const results = ref<Email[]>([])
const errorMessage = ref('')

const searchEmails = async () => {
  try {
    const response = await fetch(`http://localhost:3000/emails?term=${search.value}`)
    if (!response.ok) {
      throw new Error(`Failed to fetch emails: ${response.statusText}`)
    }
    const jsonObj = await response.json()
    console.log(jsonObj)
    if (jsonObj.success) {
      errorMessage.value = ''
      results.value = jsonObj.data.emails
    } else {
      errorMessage.value = jsonObj.message
    }
  } catch (error) {
    errorMessage.value = 'Failed to connect to the server. Please try again later.'
    console.error(error)
  }
}
</script>

<template>
  <h1>ENRON EMAIL SEARCH ENGINE</h1>
  <input type="text" placeholder="Search for emails" v-model="search" />
  <button @click="searchEmails">Search</button>
  <div v-if="errorMessage">
    <p style="color: red">{{ errorMessage }}</p>
  </div>
  <div>
    <h2>Results</h2>
    <ul>
      <li v-for="email in results" :key="email.id">
        <h3>{{ email.subject }}</h3>
        <p>{{ email.body }}</p>
      </li>
    </ul>
  </div>
</template>

<style scoped></style>
