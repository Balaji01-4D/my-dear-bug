<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import Gap from "$lib/components/Gap.svelte";
    import ConfessionModal from "$lib/components/ConfessionModal.svelte";
    import TagCloud from "$lib/components/TagCloud.svelte";
    import api from "$lib/api";
    import type { Confession, LoadingState } from "$lib/types";
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { browser } from "$app/environment";

    let confessions: Confession[] = $state([]);
    let loadingState: LoadingState = $state('idle');
    let errorMessage = $state('');
    let currentPage = $state(0);
    const pageSize = 12;

    // Get query parameters
    let languageFilter = $derived(browser ? $page.url.searchParams.get('language') : null);
    let tagFilter = $derived(browser ? $page.url.searchParams.get('tag') : null);
    let filterParam = $derived(browser ? $page.url.searchParams.get('filter') : null);
    let selectedFilter = $state('all');

    // Modal state
    let selectedConfession: Confession | null = $state(null);
    let showModal = $state(false);

    // Update selectedFilter when URL changes
    $effect(() => {
        if (browser && filterParam) {
            selectedFilter = filterParam;
        }
    });

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
            
            // Also update the modal confession if it's the same one
            if (selectedConfession && selectedConfession.id === confessionId) {
                selectedConfession = { ...selectedConfession, upvotes: selectedConfession.upvotes + 1 };
            }
        } catch (error) {
            alert(api.handleError(error as Error));
        }
    }

    function openConfessionModal(confession: Confession, event?: MouseEvent) {
        // If Ctrl/Cmd key is pressed, open in new tab
        if (event && (event.ctrlKey || event.metaKey)) {
            window.open(`/confessions/${confession.id}`, '_blank');
            return;
        }
        
        // Otherwise open modal
        selectedConfession = confession;
        showModal = true;
    }

    function closeConfessionModal() {
        showModal = false;
        selectedConfession = null;
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
        
        <!-- Popular Tags Section -->
        {#if !languageFilter && !tagFilter && selectedFilter === 'all'}
            <div class="tags-section">
                <TagCloud title="Browse by Tags" limit={12} />
            </div>
            <Gap size={25} />
        {/if}
        
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
                <div class="confession-card" onclick={(e) => openConfessionModal(confession, e)}>
                    <div class="confession-header">
                        <h3 class="confession-title">{confession.title}</h3>
                        <div class="confession-meta">
                            <span class="language-tag">{confession.language}</span>
                            <button 
                                class="heart-button"
                                onclick={(e) => {
                                    e.stopPropagation();
                                    handleUpvote(confession.id);
                                }}
                            >
                                <svg class="heart-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" 
                                          stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                </svg>
                                <span class="upvote-count">{confession.upvotes}</span>
                            </button>
                        </div>
                    </div>
                    
                    <p class="confession-description">
                        {confession.description.length > 200 
                            ? confession.description.substring(0, 200) + '...' 
                            : confession.description}
                    </p>
                    
                    {#if confession.snippet}
                        <div class="code-preview">
                            <div class="code-preview-header">
                                <span class="code-language">{confession.language}</span>
                                <span class="view-full">Click to view full code</span>
                            </div>
                        </div>
                    {/if}
                    
                    <div class="confession-tags">
                        {#each confession.tags.slice(0, 4) as tag}
                            <span class="tag">{tag.name}</span>
                        {/each}
                        {#if confession.tags.length > 4}
                            <span class="tag more-tags">+{confession.tags.length - 4} more</span>
                        {/if}
                    </div>
                    
                    <div class="confession-footer">
                        <span class="created-at">
                            {new Date(confession.createdAt).toLocaleDateString()}
                        </span>
                        <span class="click-hint">Click to view • Ctrl+Click for new tab</span>
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

<!-- Confession Modal -->
{#if showModal && selectedConfession}
    <ConfessionModal 
        confession={selectedConfession} 
        onClose={closeConfessionModal}
        onUpvote={handleUpvote}
    />
{/if}

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
    
    .tags-section {
        background: var(--bg-card);
        border: 1px solid var(--border-card);
        border-radius: 12px;
        padding: 20px;
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
        transition: all 0.3s ease;
        cursor: pointer;
        position: relative;
    }
    
    .confession-card:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
        border-color: var(--border-input-focus);
    }
    
    .confession-card:hover .click-hint {
        opacity: 1;
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
        line-height: 1.3;
        transition: color 0.2s ease;
    }
    
    .confession-card:hover .confession-title {
        color: #667eea;
    }
    
    .confession-meta {
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        gap: 8px;
        flex-shrink: 0;
    }
    
    .language-tag {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        padding: 4px 12px;
        border-radius: 20px;
        font-size: 0.8rem;
        font-weight: 500;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }
    
    .heart-button {
        background: none;
        border: 1px solid var(--border-input);
        color: var(--font-color-accent);
        padding: 6px 10px;
        border-radius: 20px;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 0.9rem;
        font-weight: 600;
        transition: all 0.3s ease;
        position: relative;
        z-index: 10;
    }
    
    .heart-button:hover {
        border-color: #e91e63;
        color: #e91e63;
        background: rgba(233, 30, 99, 0.1);
        transform: scale(1.05);
    }
    
    .heart-icon {
        transition: all 0.3s ease;
    }
    
    .heart-button:hover .heart-icon path {
        fill: rgba(233, 30, 99, 0.2);
    }
    
    .upvote-count {
        font-weight: 700;
        min-width: 15px;
        text-align: left;
    }
    
    .confession-description {
        color: var(--font-color);
        line-height: 1.6;
        margin-bottom: 15px;
        display: -webkit-box;
        -webkit-line-clamp: 4;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }
    
    .code-preview {
        background: linear-gradient(135deg, #1e1e1e 0%, #2d2d2d 100%);
        border: 1px solid #444;
        border-radius: 8px;
        padding: 12px 16px;
        margin-bottom: 15px;
        position: relative;
    }
    
    .code-preview-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
    
    .code-language {
        color: #f8f8f2;
        font-size: 0.8rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }
    
    .view-full {
        color: #667eea;
        font-size: 0.8rem;
        font-style: italic;
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
        margin-bottom: 15px;
    }
    
    .tag {
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color-accent);
        padding: 4px 10px;
        border-radius: 16px;
        font-size: 0.8rem;
        transition: all 0.2s ease;
    }
    
    .tag:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
        color: var(--font-color);
    }
    
    .more-tags {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        border-color: transparent;
    }
    
    .confession-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 10px;
    }
    
    .created-at {
        color: var(--font-color-accent);
        font-size: 0.8rem;
    }
    
    .click-hint {
        color: var(--font-color-accent);
        font-size: 0.8rem;
        font-style: italic;
        opacity: 0;
        transition: opacity 0.3s ease;
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
    
    /* Mobile responsiveness */
    @media (max-width: 768px) {
        .hero-section h2 {
            font-size: 2rem;
        }
        
        .confession-card {
            padding: 16px;
        }
        
        .confession-header {
            flex-direction: column;
            gap: 12px;
            align-items: stretch;
        }
        
        .confession-meta {
            flex-direction: row;
            justify-content: space-between;
            align-items: center;
        }
        
        .confession-title {
            font-size: 1.2rem;
        }
        
        .confession-description {
            -webkit-line-clamp: 3;
        }
        
        .filter-section {
            flex-direction: column;
            gap: 8px;
        }
        
        .filter-button {
            padding: 10px 16px;
        }
        
        .click-hint {
            font-size: 0.75rem;
        }
        
        .heart-button {
            padding: 8px 12px;
        }
    }
    
    @media (max-width: 480px) {
        .confessions-container {
            padding: 0 10px;
        }
        
        .confession-card {
            padding: 12px;
        }
        
        .confession-tags {
            gap: 6px;
        }
        
        .tag {
            font-size: 0.75rem;
            padding: 3px 8px;
        }
        
        .language-tag {
            font-size: 0.75rem;
            padding: 3px 10px;
        }
    }
</style>
