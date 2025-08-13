export type Tag = {
  id: number
  name: string
}

export type Confession = {
  id: number
  title: string
  description: string
  language: string
  snippet: string
  tags: Tag[]
  sentiment?: string
  isFlagged?: boolean
  createdAt: string
  upvotes: number
}

export type ConfessionRequest = {
  title: string
  description: string
  snippet?: string
  language: string
  tags?: string[]
  isFlagged?: boolean
}
