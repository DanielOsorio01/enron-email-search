import axios from 'axios'
import type { SearchResponse } from '../types/email'

const api = axios.create({
  baseURL: 'http://localhost:3000',
})

export const searchEmails = async (term: string) => {
  try {
    // Add metadata to the request
    const response = await api.get<SearchResponse>(`/emails`, {
      params: {
        term,
        max_results: 10,
      },
    })
    return response.data
  } catch (error) {
    console.error('Error searching emails:', error)
    throw error
  }
}
