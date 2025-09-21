<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import Gap from "$lib/components/Gap.svelte";
    import api from "$lib/api";
    import type { Confession, LoadingState } from "$lib/types";
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { browser } from "$app/environment";

    let confessions: Confession[] = $state([]);
    let loadingState: LoadingState = $state('idle');
    let errorMessage = $state('');
    let selectedFilter = $state('all');
    let currentPage = $state(0);
    const pageSize = 12;

    // Get query parameters
    let languageFilter = $derived(browser ? $page.url.searchParams.get('language') : null);
    let tagFilter = $derived(browser ? $page.url.searchParams.get('tag') : null);

    const filters = [
        { value: 'all', label: 'All Confessions' },
        { value: 'top', label: 'Top Voted' },
        { value: 'trending-weekly', label: 'Trending This Week' },
        { value: 'trending-monthly', label: 'Trending This Month' },
        { value: 'hall-of-fame', label: 'Hall of Fame' },
        { value: 'random', label: 'Random' }
    ];

    async function loadConfessions(filter: string = 'all', page: number = 0, language?: string | null, tag?: string | null) {
        loadingState = 'loading';
        errorMessage = '';
        
        try {
            const params = { limit: pageSize, offset: page * pageSize };
            
            // If filtering by language, use the language endpoint
            if (language) {
                confessions = await api.confessions.getByLanguage(language, params);
            } else {
                // Handle other filters
                switch (filter) {
                    case 'top':
                        confessions = await api.confessions.getTop(params);
                        break;
                    case 'trending-weekly':
                        confessions = await api.confessions.getTrendingWeekly(params);
                        break;
                    case 'trending-monthly':
                        confessions = await api.confessions.getTrendingMonthly(params);
                        break;
                    case 'hall-of-fame':
                        confessions = await api.confessions.getHallOfFame(params);
                        break;
                    case 'random':
                        // For random, just get one confession
                        const randomConfession = await api.confessions.getRandom();
                        confessions = [randomConfession];
                        break;
                    default:
                        confessions = await api.confessions.getAll(params);
                }
            }
            
            // Client-side filtering for tags since we don't have a tag filter endpoint
            if (tag && !language) {
                confessions = confessions.filter(c => 
                    c.tags.some(t => t.name.toLowerCase() === tag.toLowerCase())
                );
            }
            
            loadingState = 'success';
        } catch (error) {
            loadingState = 'error';
            errorMessage = api.handleError(error as Error);
        }
    }

    async function handleUpvote(confessionId: number) {
        try {
            await api.confessions.upvote(confessionId);
            // Update the upvote count locally
            confessions = confessions.map(c => 
                c.id === confessionId 
                    ? { ...c, upvotes: c.upvotes + 1 }
                    : c
            );
        } catch (error) {
            alert(api.handleError(error as Error));
        }
    }

    function handleFilterChange(filter: string) {
        selectedFilter = filter;
        currentPage = 0;
        loadConfessions(filter, 0, languageFilter, tagFilter);
    }

    function loadMore() {
        if (selectedFilter !== 'random') {
            currentPage += 1;
            loadConfessions(selectedFilter, currentPage, languageFilter, tagFilter);
        }
    }

    // Watch for URL parameter changes and reload data
    $effect(() => {
        if (browser) {
            currentPage = 0;
            loadConfessions(selectedFilter, 0, languageFilter, tagFilter);
        }
    });

    onMount(() => {
        loadConfessions(selectedFilter, 0, languageFilter, tagFilter);
    });
</script>

<Page title="Browse Confessions">
    <section class="confessions-container">
        <div class="hero-section">
            <h2>Coding Confessions</h2>
            <p>Learn from others' debugging adventures and coding mishaps</p>
        </div>
        
        <Gap size={20} />
        
        <!-- Active Filters Info -->
        {#if languageFilter || tagFilter}
            <div class="active-filters">
                <span>Filtering by:</span>
                {#if languageFilter}
                    <span class="filter-tag">Language: {languageFilter}</span>
                {/if}
                {#if tagFilter}
                    <span class="filter-tag">Tag: #{tagFilter}</span>
                {/if}
                <a href="/confessions" class="clear-filters">Clear filters</a>
            </div>
            <Gap size={15} />
        {/if}
        
        <!-- Filter Options -->
        <div class="filter-section">
            {#each filters as filter}
                <button 
                    type="button"
                    class="filter-button"
                    class:active={selectedFilter === filter.value}
                    onclick={() => handleFilterChange(filter.value)}
                >
                    {filter.label}
                </button>
            {/each}
        </div>
        
        <Gap size={20} />
        
        <!-- Loading State -->
        {#if loadingState === 'loading'}
            <div class="loading">
                <p>Loading confessions...</p>
            </div>
        {/if}
        
        <!-- Error State -->
        {#if loadingState === 'error'}
            <div class="error-message">
                <p>❌ {errorMessage}</p>
                <button onclick={() => loadConfessions(selectedFilter, currentPage)} class="retry-button">
                    Try Again
                </button>
            </div>
        {/if}
        
        <!-- Confessions List -->
        {#if loadingState === 'success' && confessions.length > 0}
            {#each confessions as confession}
                <div class="confession-card">
                    <div class="confession-header">
                        <h3 class="confession-title">{confession.title}</h3>
                        <div class="confession-meta">
                            <span class="language-tag">{confession.language}</span>
                            <button 
                                class="upvote-button"
                                onclick={() => handleUpvote(confession.id)}
                            >
                                ⭐ {confession.upvotes}
                            </button>
                        </div>
                    </div>
                    
                    <p class="confession-description">{confession.description}</p>
                    
                    {#if confession.snippet}
                        <div class="code-snippet">
                            <pre><code>{confession.snippet}</code></pre>
                        </div>
                    {/if}
                    
                    <div class="confession-tags">
                        {#each confession.tags as tag}
                            <span class="tag">{tag.name}</span>
                        {/each}
                    </div>
                    
                    <div class="confession-footer">
                        <span class="created-at">
                            {new Date(confession.createdAt).toLocaleDateString()}
                        </span>
                    </div>
                </div>
                
                <Gap size={15} />
            {/each}
            
            <!-- Load More Button -->
            {#if selectedFilter !== 'random' && confessions.length >= pageSize}
                <div class="load-more-section">
                    <button onclick={loadMore} class="load-more-button">
                        Load More
                    </button>
                </div>
            {/if}
        {/if}
        
        <!-- Empty State -->
        {#if loadingState === 'success' && confessions.length === 0}
            <div class="empty-state">
                <p>No confessions found. Be the first to share your story!</p>
            </div>
        {/if}
    </section>
</Page>

<style>
    .confessions-container {
        width: 100%;
        max-width: 800px;
        margin: 0 auto;
    }
    
    .hero-section {
        text-align: center;
        margin-bottom: 30px;
    }
    
    .hero-section h2 {
        font-size: 2.5rem;
        font-weight: bold;
        color: var(--font-color);
        margin-bottom: 10px;
    }
    
    .hero-section p {
        color: var(--font-color-accent);
        font-size: 1.1rem;
    }
    
    .active-filters {
        display: flex;
        align-items: center;
        gap: 10px;
        flex-wrap: wrap;
        justify-content: center;
        padding: 10px 20px;
        background: var(--bg-card);
        border: 1px solid var(--border-card);
        border-radius: 8px;
        font-size: 0.9rem;
    }
    
    .filter-tag {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        padding: 4px 8px;
        border-radius: 12px;
        font-size: 0.8rem;
    }
    
    .clear-filters {
        color: var(--font-color-accent);
        text-decoration: none;
        font-size: 0.8rem;
    }
    
    .clear-filters:hover {
        color: var(--font-color);
    }
    
    .filter-section {
        display: flex;
        flex-wrap: wrap;
        gap: 10px;
        justify-content: center;
        margin-bottom: 20px;
    }
    
    .filter-button {
        padding: 8px 16px;
        border: 1px solid var(--border-input);
        background: var(--bg-input);
        color: var(--font-color);
        border-radius: 20px;
        cursor: pointer;
        font-size: 0.9rem;
        transition: all 0.2s ease;
    }
    
    .filter-button:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
    }
    
    .filter-button.active {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        border-color: transparent;
    }
    
    .loading, .empty-state {
        text-align: center;
        padding: 40px 20px;
        color: var(--font-color-accent);
    }
    
    .error-message {
        text-align: center;
        padding: 20px;
        background: rgba(255, 0, 0, 0.1);
        border: 1px solid rgba(255, 0, 0, 0.2);
        border-radius: 8px;
        color: var(--font-color);
    }
    
    .retry-button {
        margin-top: 10px;
        padding: 8px 16px;
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color);
        border-radius: 4px;
        cursor: pointer;
    }
    
    .confession-card {
        background: var(--bg-card);
        border: 1px solid var(--border-card);
        border-radius: 12px;
        padding: 20px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        transition: transform 0.2s ease, box-shadow 0.2s ease;
    }
    
    .confession-card:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    }
    
    .confession-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 15px;
        gap: 15px;
    }
    
    .confession-title {
        font-size: 1.4rem;
        font-weight: 600;
        color: var(--font-color);
        margin: 0;
        flex: 1;
    }
    
    .confession-meta {
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        gap: 5px;
        flex-shrink: 0;
    }
    
    .language-tag {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        padding: 4px 12px;
        border-radius: 20px;
        font-size: 0.8rem;
        font-weight: 500;
    }
    
    .upvote-button {
        background: none;
        border: none;
        color: var(--font-color-accent);
        font-size: 0.9rem;
        font-weight: 500;
        cursor: pointer;
        transition: color 0.2s ease;
    }
    
    .upvote-button:hover {
        color: var(--font-color);
    }
    
    .confession-description {
        color: var(--font-color);
        line-height: 1.6;
        margin-bottom: 15px;
    }
    
    .code-snippet {
        background: #1e1e1e;
        border: 1px solid var(--border-input);
        border-radius: 8px;
        padding: 15px;
        margin-bottom: 15px;
        overflow-x: auto;
    }
    
    .code-snippet pre {
        margin: 0;
        color: #f8f8f2;
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 0.9rem;
        line-height: 1.4;
    }
    
    .confession-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        margin-bottom: 10px;
    }
    
    .tag {
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color-accent);
        padding: 4px 10px;
        border-radius: 16px;
        font-size: 0.8rem;
    }
    
    .confession-footer {
        text-align: right;
        margin-top: 10px;
    }
    
    .created-at {
        color: var(--font-color-accent);
        font-size: 0.8rem;
    }
    
    .load-more-section {
        text-align: center;
        margin-top: 20px;
    }
    
    .load-more-button {
        padding: 10px 24px;
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color);
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.2s ease;
    }
    
    .load-more-button:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
    }
</style>
