<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import { api } from "$lib/api";
    import { onMount } from 'svelte';
    import { browser } from '$app/environment';
    import type { PageData } from './$types';
    
    let { data }: { data: PageData } = $props();
    
    // Search state
    let searchQuery = $state(data.initialQuery || '');
    let searchResults = $state([]);
    let isLoading = $state(false);
    let hasSearched = $state(false);
    let error = $state('');

    // Pagination
    let currentPage = $state(1);
    let totalPages = $state(1);
    let totalResults = $state(0);
    const itemsPerPage = 10;

    // SessionStorage key for preserving search state
    const SEARCH_STATE_KEY = 'search_state';

    // Save search state to sessionStorage
    function saveSearchState() {
        if (browser) {
            const state = {
                searchQuery,
                searchResults,
                hasSearched,
                currentPage,
                totalPages,
                totalResults,
                timestamp: Date.now()
            };
            sessionStorage.setItem(SEARCH_STATE_KEY, JSON.stringify(state));
        }
    }

    // Load search state from sessionStorage
    function loadSearchState() {
        if (browser) {
            const saved = sessionStorage.getItem(SEARCH_STATE_KEY);
            if (saved) {
                try {
                    const state = JSON.parse(saved);
                    // Only restore if timestamp is less than 30 minutes old
                    if (Date.now() - state.timestamp < 30 * 60 * 1000) {
                        searchQuery = state.searchQuery || '';
                        searchResults = state.searchResults || [];
                        hasSearched = state.hasSearched || false;
                        currentPage = state.currentPage || 1;
                        totalPages = state.totalPages || 1;
                        totalResults = state.totalResults || 0;
                    }
                } catch (e) {
                    // Invalid saved state, ignore
                }
            }
        }
    }

    // Popular languages and tags (you can customize these based on your data)
    // Removed as per user request - no quick filters needed

    // Search function
    async function performSearch(page = 1) {
        if (!searchQuery.trim()) {
            error = 'Please enter a search query';
            return;
        }

        isLoading = true;
        error = '';
        
        try {
            const offset = (page - 1) * itemsPerPage;
            
            const params: any = {
                limit: itemsPerPage,
                offset: offset
            };
            
            if (searchQuery.trim()) {
                params.q = searchQuery.trim();
            }

            const result = await api.confessions.search(params);
            
            searchResults = result.confessions || [];
            totalResults = result.total || 0;
            totalPages = Math.ceil(totalResults / itemsPerPage);
            currentPage = page;
            hasSearched = true;
            
            // Save state after successful search
            saveSearchState();
            
        } catch (err: any) {
            error = api.handleError(err);
            searchResults = [];
            totalResults = 0;
        } finally {
            isLoading = false;
        }
    }

    // Handle search form submission
    function handleSearch(event: Event) {
        event.preventDefault();
        performSearch(1);
    }

    // Format date helper
    function formatDate(dateString: string) {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
    }

    // Truncate text helper
    function truncateText(text: string, maxLength: number = 150) {
        if (text.length <= maxLength) return text;
        return text.substr(0, maxLength) + '...';
    }

    // Auto-search if initial parameters are provided
    onMount(() => {
        // First, try to load saved state
        loadSearchState();
        
        // If no saved state and we have initial query, perform search
        if (!hasSearched && data.initialQuery) {
            performSearch(1);
        }
    });
</script>

<Page title="Search Confessions">
    <section>
        <!-- Search Form -->
        <div class="search-container">
            <Group title="🔍 Search Bug Confessions">
                <form on:submit={handleSearch} class="search-form">
                    <div class="search-input-container">
                        <input
                            type="text"
                            bind:value={searchQuery}
                            placeholder="Search for debugging stories, error messages, solutions..."
                            class="search-input"
                        />
                        <button type="submit" class="search-button" disabled={isLoading}>
                            {#if isLoading}
                                <div class="loading-spinner"></div>
                            {:else}
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <circle cx="11" cy="11" r="8"/>
                                    <path d="m21 21-4.35-4.35"/>
                                </svg>
                            {/if}
                        </button>
                    </div>
                </form>
            </Group>
        </div>

        <!-- Active Filters Display -->
        {#if searchQuery}
            <div class="active-filters">
                <span class="filter-label">Searching for:</span>
                <span class="active-filter">"{searchQuery}"</span>
                <button class="clear-filter-btn" on:click={() => {
                    searchQuery = ''; 
                    searchResults = []; 
                    hasSearched = false;
                    totalResults = 0;
                    currentPage = 1;
                    totalPages = 1;
                    error = '';
                    // Clear saved state
                    if (browser) {
                        sessionStorage.removeItem(SEARCH_STATE_KEY);
                    }
                }}>
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="18" y1="6" x2="6" y2="18"/>
                        <line x1="6" y1="6" x2="18" y2="18"/>
                    </svg>
                </button>
            </div>
        {/if}

        <!-- Error Display -->
        {#if error}
            <div class="error-message">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="12" y1="8" x2="12" y2="12"/>
                    <line x1="12" y1="16" x2="12.01" y2="16"/>
                </svg>
                {error}
            </div>
        {/if}

        <!-- Search Results -->
        {#if hasSearched}
            <div class="results-section">
                <Group title={`Search Results (${totalResults} found)`}>
                    {#if searchResults.length > 0}
                        <div class="results-list">
                            {#each searchResults as result}
                                <a href="/confessions/{result.id}" class="result-item">
                                    <div class="result-header">
                                        <h3 class="result-title">{result.title}</h3>
                                        <div class="result-meta">
                                            {#if result.language}
                                                <span class="language-badge">{result.language}</span>
                                            {/if}
                                            <span class="date">{formatDate(result.created_at)}</span>
                                        </div>
                                    </div>
                                    <p class="result-description">
                                        {truncateText(result.description || result.story)}
                                    </p>
                                    {#if result.tags && result.tags.length > 0}
                                        <div class="result-tags">
                                            {#each result.tags.slice(0, 5) as tag}
                                                <span class="tag-badge">{tag}</span>
                                            {/each}
                                        </div>
                                    {/if}
                                </a>
                            {/each}
                        </div>

                        <!-- Pagination -->
                        {#if totalPages > 1}
                            <div class="pagination">
                                <button
                                    class="page-button"
                                    disabled={currentPage === 1}
                                    on:click={() => performSearch(currentPage - 1)}
                                >
                                    Previous
                                </button>
                                
                                <span class="page-info">
                                    Page {currentPage} of {totalPages}
                                </span>
                                
                                <button
                                    class="page-button"
                                    disabled={currentPage === totalPages}
                                    on:click={() => performSearch(currentPage + 1)}
                                >
                                    Next
                                </button>
                            </div>
                        {/if}
                    {:else}
                        <div class="no-results">
                            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="var(--font-color-accent)" stroke-width="2">
                                <circle cx="11" cy="11" r="8"/>
                                <path d="m21 21-4.35-4.35"/>
                            </svg>
                            <h3>No confessions found</h3>
                            <p>Try adjusting your search terms or filters</p>
                        </div>
                    {/if}
                </Group>
            </div>
        {/if}
    </section>
</Page>

<style>
    section {
        display: flex;
        flex-direction: column;
        gap: 30px;
        padding: 20px 0;
        max-width: 1200px;
        margin: 0 auto;
        width: 100%;
    }

    .search-container {
        width: 100%;
    }

    .search-form {
        padding: 0;
    }

    .search-input-container {
        display: flex;
        gap: 10px;
        align-items: center;
    }

    .search-input {
        flex: 1;
        padding: 15px 20px;
        border: 2px solid rgba(255, 255, 255, 0.1);
        border-radius: 10px;
        background: rgba(44, 39, 51, 0.3);
        color: var(--font-color);
        font-size: 1rem;
        transition: all 0.3s ease;
    }

    .search-input:focus {
        outline: none;
        border-color: rgba(102, 126, 234, 0.5);
        box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }

    .search-button {
        padding: 15px 20px;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.8) 0%, rgba(118, 75, 162, 0.8) 100%);
        border: none;
        border-radius: 10px;
        color: white;
        cursor: pointer;
        transition: all 0.3s ease;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .search-button:hover:not(:disabled) {
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.9) 0%, rgba(118, 75, 162, 0.9) 100%);
        transform: translateY(-1px);
    }

    .search-button:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .loading-spinner {
        width: 20px;
        height: 20px;
        border: 2px solid rgba(255, 255, 255, 0.3);
        border-top: 2px solid white;
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        0% { transform: rotate(0deg); }
        100% { transform: rotate(360deg); }
    }

    .active-filters {
        display: flex;
        flex-wrap: wrap;
        gap: 10px;
        align-items: center;
        padding: 15px 20px;
        background: rgba(44, 39, 51, 0.2);
        border-radius: 10px;
        border: 1px solid rgba(255, 255, 255, 0.1);
    }

    .filter-label {
        font-weight: 600;
        color: var(--font-color);
    }

    .active-filter {
        padding: 4px 12px;
        background: rgba(102, 126, 234, 0.2);
        border-radius: 15px;
        color: var(--font-color);
        font-size: 0.9rem;
    }

    .clear-filter-btn {
        padding: 4px;
        background: rgba(255, 99, 71, 0.2);
        border: none;
        border-radius: 12px;
        color: rgba(255, 99, 71, 0.8);
        cursor: pointer;
        transition: all 0.3s ease;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .clear-filter-btn:hover {
        background: rgba(255, 99, 71, 0.3);
        color: rgba(255, 99, 71, 1);
    }

    .error-message {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 15px 20px;
        background: rgba(255, 99, 71, 0.1);
        border: 1px solid rgba(255, 99, 71, 0.3);
        border-radius: 10px;
        color: rgba(255, 99, 71, 0.9);
    }

    .results-list {
        display: flex;
        flex-direction: column;
        gap: 15px;
    }

    .result-item {
        display: block;
        padding: 25px 30px;
        background: rgba(44, 39, 51, 0.3);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 12px;
        text-decoration: none;
        color: var(--font-color);
        transition: all 0.3s ease;
        width: 100%;
    }

    .result-item:hover {
        background: rgba(44, 39, 51, 0.5);
        border-color: rgba(102, 126, 234, 0.3);
        transform: translateY(-2px);
        box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
    }

    .result-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 10px;
        flex-wrap: wrap;
        gap: 10px;
    }

    .result-title {
        margin: 0;
        font-size: 1.3rem;
        font-weight: 600;
        color: var(--font-color);
        flex: 1;
        min-width: 0;
    }

    .result-meta {
        display: flex;
        gap: 10px;
        align-items: center;
        flex-shrink: 0;
    }

    .language-badge {
        padding: 4px 10px;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.8) 0%, rgba(118, 75, 162, 0.8) 100%);
        border-radius: 12px;
        color: white;
        font-size: 0.8rem;
        font-weight: 500;
    }

    .date {
        color: var(--font-color-accent);
        font-size: 0.9rem;
    }

    .result-description {
        margin: 0 0 15px 0;
        color: var(--font-color-accent);
        line-height: 1.5;
    }

    .result-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 6px;
    }

    .tag-badge {
        padding: 3px 8px;
        background: rgba(255, 255, 255, 0.1);
        border-radius: 10px;
        color: var(--font-color-accent);
        font-size: 0.8rem;
    }

    .pagination {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 20px;
        margin-top: 30px;
        padding: 20px;
    }

    .page-button {
        padding: 10px 20px;
        background: rgba(44, 39, 51, 0.6);
        border: 1px solid rgba(255, 255, 255, 0.2);
        border-radius: 8px;
        color: var(--font-color);
        cursor: pointer;
        transition: all 0.3s ease;
    }

    .page-button:hover:not(:disabled) {
        background: rgba(102, 126, 234, 0.3);
        border-color: rgba(102, 126, 234, 0.5);
    }

    .page-button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .page-info {
        color: var(--font-color-accent);
        font-weight: 500;
    }

    .no-results {
        text-align: center;
        padding: 60px 20px;
    }

    .no-results h3 {
        margin: 20px 0 10px 0;
        color: var(--font-color);
        font-size: 1.5rem;
    }

    .no-results p {
        color: var(--font-color-accent);
        font-size: 1.1rem;
    }

    /* Responsive Design */
    @media (max-width: 768px) {
        .search-input-container {
            flex-direction: column;
        }

        .search-input, .search-button {
            width: 100%;
        }

        .result-header {
            flex-direction: column;
            align-items: flex-start;
        }

        .active-filters {
            flex-direction: column;
            align-items: flex-start;
        }

        .pagination {
            flex-direction: column;
            gap: 15px;
        }
    }
</style>
