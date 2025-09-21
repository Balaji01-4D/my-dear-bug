export type Tag = {
  id: number;
  name: string;
};

export type Confession = {
  id: number;
  title: string;
  description: string;
  language: string;
  snippet: string;
  tags: Tag[];
  sentiment?: string;
  isFlagged?: boolean;
  createdAt: string;
  upvotes: number;
};

export type ConfessionRequest = {
  title: string;
  description: string;
  snippet?: string;
  language: string;
  tags?: string[];
  isFlagged?: boolean;
};

// Loading states for UI components
export type LoadingState = 'idle' | 'loading' | 'success' | 'error';

// API response wrapper for paginated results (if needed)
export type PaginatedResponse<T> = {
  data: T[];
  pagination?: {
    offset: number;
    limit: number;
    total?: number;
  };
};

// Error response structure from backend
export type ErrorResponse = {
  error: string;
  message?: string;
  statusCode?: number;
};
