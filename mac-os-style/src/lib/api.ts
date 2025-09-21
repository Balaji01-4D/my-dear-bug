import type { Confession, ConfessionRequest, Tag } from './types';

// Base URL for your hosted backend
const API_BASE = 'https://my-dear-bug.onrender.com';

// Simple fetch wrapper with error handling
async function apiFetch(endpoint: string, options?: RequestInit): Promise<Response> {
  const url = `${API_BASE}${endpoint}`;
  
  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    ...options,
  });

  if (!response.ok) {
    let errorMessage = `HTTP ${response.status}`;
    try {
      const errorData = await response.json();
      errorMessage = errorData?.error || errorData?.message || errorMessage;
    } catch {
      // Use default error message
    }
    throw new Error(errorMessage);
  }

  return response;
}

// Build URL with query parameters
function buildUrl(endpoint: string, params?: Record<string, string | number | undefined>): string {
  if (!params) return endpoint;
  
  const url = new URL(endpoint, API_BASE);
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== '') {
      url.searchParams.set(key, String(value));
    }
  });
  
  return url.pathname + url.search;
}

// Confession API endpoints
export const confessionApi = {
  // GET /confessions - List with pagination
  async getAll(params?: { limit?: number; offset?: number }): Promise<Confession[]> {
    const endpoint = buildUrl('/confessions', params);
    const response = await apiFetch(endpoint);
    return response.json();
  },

  // GET /confessions/:id - Get confession details
  async getById(id: number): Promise<Confession> {
    const response = await apiFetch(`/confessions/${id}`);
    return response.json();
  },

  // POST /confessions - Create a confession
  async create(confession: ConfessionRequest): Promise<Confession> {
    const response = await apiFetch('/confessions', {
      method: 'POST',
      body: JSON.stringify(confession),
    });
    return response.json();
  },

  // GET /confessions/language/:language - Filter by language
  async getByLanguage(language: string, params?: { limit?: number; offset?: number }): Promise<Confession[]> {
    const endpoint = buildUrl(`/confessions/language/${encodeURIComponent(language)}`, params);
    const response = await apiFetch(endpoint);
    return response.json();
  },

  // GET /confessions/top - Highest upvoted confessions
  async getTop(params?: { limit?: number }): Promise<Confession[]> {
    const endpoint = buildUrl('/confessions/top', params);
    const response = await apiFetch(endpoint);
    return response.json();
  },

  // GET /confessions/trending/weekly - Trending confessions for the last 7 days
  async getTrendingWeekly(params?: { limit?: number }): Promise<Confession[]> {
    const endpoint = buildUrl('/confessions/trending/weekly', params);
    const response = await apiFetch(endpoint);
    return response.json();
  },

  // GET /confessions/trending/monthly - Trending confessions for the last 30 days  
  async getTrendingMonthly(params?: { limit?: number }): Promise<Confession[]> {
    const endpoint = buildUrl('/confessions/trending/monthly', params);
    const response = await apiFetch(endpoint);
    return response.json();
  },

  // GET /confessions/hall-of-fame - All-time notable confessions
  async getHallOfFame(params?: { limit?: number }): Promise<Confession[]> {
    const endpoint = buildUrl('/confessions/hall-of-fame', params);
    const response = await apiFetch(endpoint);
    return response.json();
  },

  // GET /confessions/random - Random confession
  async getRandom(): Promise<Confession> {
    const response = await apiFetch('/confessions/random');
    return response.json();
  },

  // POST /confessions/:id/upvote - Upvote a confession
  async upvote(id: number): Promise<void> {
    await apiFetch(`/confessions/${id}/upvote`, {
      method: 'POST',
    });
  },
};

// Tag API endpoints
export const tagApi = {
  // GET /tags - List all tags
  async getAll(): Promise<Tag[]> {
    const response = await apiFetch('/tags');
    return response.json();
  },

  // GET /tags/suggest?query=<prefix> - Autocomplete suggestions
  async suggest(query: string): Promise<Tag[]> {
    const endpoint = buildUrl('/tags/suggest', { query });
    const response = await apiFetch(endpoint);
    return response.json();
  },
};

// Utility functions for common operations
export const api = {
  confessions: confessionApi,
  tags: tagApi,
  
  // Helper to handle common error scenarios
  handleError(error: Error): string {
    console.error('API Error:', error);
    
    if (error.message.includes('Failed to fetch')) {
      return 'Unable to connect to the server. Please check your internet connection.';
    }
    
    if (error.message.includes('404')) {
      return 'The requested resource was not found.';
    }
    
    if (error.message.includes('500')) {
      return 'Server error. Please try again later.';
    }
    
    return error.message || 'An unexpected error occurred.';
  },
};

export default api;
