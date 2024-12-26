export interface Email {
  id: string
  from: string
  to: string
  subject: string
  body: string
  date: string
}

export interface SearchResponse {
  success: boolean
  data: {
    emails: Email[]
    total: number
    from: number
    size: number
  }
  error?: string
}
