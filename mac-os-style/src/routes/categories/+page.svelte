<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import Gap from "$lib/components/Gap.svelte";
    import api from "$lib/api";
    import type { Tag, LoadingState } from "$lib/types";
    import { onMount } from "svelte";
    
    let tags: Tag[] = $state([]);
    let loadingState: LoadingState = $state('idle');
    let errorMessage = $state('');
    
    // Common programming languages with their colors
    const languageColors: Record<string, string> = {
        'javascript': '#f7df1e',
        'typescript': '#3178c6',
        'python': '#3776ab',
        'java': '#ed8b00',
        'css': '#1572b6',
        'html': '#e34f26',
        'go': '#00add8',
        'rust': '#ce422b',
        'php': '#777bb4',
        'ruby': '#cc342d',
        'cpp': '#00599c',
        'csharp': '#239120',
        'sql': '#336791',
        'devops': '#ff6b6b'
    };
    
    function getLanguageColor(name: string): string {
        const key = name.toLowerCase().replace(/[+#]/g, '');
        return languageColors[key] || '#6b7280';
    }
    
    async function loadTags() {
        loadingState = 'loading';
        errorMessage = '';
        
        try {
            tags = await api.tags.getAll();
            loadingState = 'success';
        } catch (error) {
            loadingState = 'error';
            errorMessage = api.handleError(error as Error);
        }
    }
    
    onMount(() => {
        loadTags();
    });
</script>

<Page title="Categories">
    <section class="categories-container">
        <div class="hero-section">
            <h2>Browse by Tags</h2>
            <p>Explore confessions organized by technology, problem type, or category</p>
        </div>
        
        <Gap size={20} />
        
        <!-- Popular Programming Languages -->
        <Group title="Popular Programming Languages">
            <div class="languages-container">
                {#each ["javascript", "python", "css", "html", "java", "typescript", "go", "rust", "php", "ruby", "sql", "devops"] as language}
                    <a href="/confessions?language={language}" class="language-link">{language}</a>
                {/each}
            </div>
        </Group>
        
        <Gap size={5} />
        
        <!-- Loading State -->
        {#if loadingState === 'loading'}
            <div class="loading">
                <p>Loading tags...</p>
            </div>
        {/if}
        
        <!-- Error State -->
        {#if loadingState === 'error'}
            <div class="error-message">
                <p>❌ {errorMessage}</p>
                <button onclick={loadTags} class="retry-button">
                    Try Again
                </button>
            </div>
        {/if}
        
        <!-- Tags Grid -->
        {#if loadingState === 'success' && tags.length > 0}
            <div class="tags-grid">
                {#each tags as tag}
                    <a href="/confessions?tag={tag.name}" class="tag-card">
                        <div class="tag-icon" style="background-color: {getLanguageColor(tag.name)}">
                            {tag.name.charAt(0).toUpperCase()}
                        </div>
                        <div class="tag-info">
                            <h3 class="tag-name">#{tag.name}</h3>
                            <p class="tag-description">Click to see confessions</p>
                        </div>
                    </a>
                {/each}
            </div>
        {/if}
        
        <!-- Empty State -->
        {#if loadingState === 'success' && tags.length === 0}
            <div class="empty-state">
                <p>No tags found. Be the first to create confessions with tags!</p>
            </div>
        {/if}
    </section>
</Page>

<style>
    .categories-container {
        width: 100%;
        max-width: 1200px;
        margin: 0 auto;
        padding: 0 20px;
    }
    
    .hero-section {
        text-align: center;
        margin-bottom: 40px;
        padding: 20px 0;
    }
    
    .hero-section h2 {
        font-size: 2.5rem;
        font-weight: bold;
        color: var(--font-color);
        margin-bottom: 15px;
        line-height: 1.2;
    }
    
    .hero-section p {
        color: var(--font-color-accent);
        font-size: 1.1rem;
        max-width: 600px;
        margin: 0 auto;
        line-height: 1.5;
    }
    
    .loading, .empty-state {
        text-align: center;
        padding: 60px 20px;
        color: var(--font-color-accent);
    }
    
    .loading p, .empty-state p {
        font-size: 1.1rem;
        margin: 0;
    }
    
    .error-message {
        text-align: center;
        padding: 30px 20px;
        background: rgba(255, 0, 0, 0.1);
        border: 1px solid rgba(255, 0, 0, 0.2);
        border-radius: 12px;
        color: var(--font-color);
        margin: 20px 0;
        max-width: 500px;
        margin-left: auto;
        margin-right: auto;
    }
    
    .retry-button {
        margin-top: 15px;
        padding: 10px 20px;
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color);
        border-radius: 8px;
        cursor: pointer;
        font-size: 0.9rem;
        transition: all 0.2s ease;
    }
    
    .retry-button:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
    }
    
    .tags-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 24px;
        margin: 30px 0;
        align-items: start;
    }
    
    .tag-card {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 24px;
        background: var(--bg-card);
        border: 1px solid var(--border-card);
        border-radius: 16px;
        text-decoration: none;
        transition: all 0.3s ease;
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
        height: fit-content;
    }
    
    .tag-card:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
        border-color: var(--border-card-hover, var(--border-card));
    }
    
    .tag-icon {
        width: 56px;
        height: 56px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-weight: bold;
        font-size: 1.4rem;
        flex-shrink: 0;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }
    
    .tag-info {
        flex: 1;
        min-width: 0;
    }
    
    .tag-name {
        margin: 0 0 8px 0;
        font-size: 1.3rem;
        font-weight: 600;
        color: var(--font-color);
        line-height: 1.3;
    }
    
    .tag-description {
        margin: 0;
        color: var(--font-color-accent);
        font-size: 0.95rem;
        line-height: 1.4;
    }
    
    .languages-container {
        display: flex;
        flex-wrap: wrap;
        gap: 12px;
        justify-content: center;
        margin: 20px 0;
    }
    
    .language-link {
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color);
        padding: 10px 18px;
        border-radius: 24px;
        text-decoration: none;
        font-size: 0.9rem;
        font-weight: 500;
        transition: all 0.2s ease;
        text-transform: capitalize;
        white-space: nowrap;
    }
    
    .language-link:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
        color: var(--font-color);
        transform: translateY(-1px);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }

    /* Mobile responsiveness */
    @media (max-width: 768px) {
        .categories-container {
            padding: 0 16px;
        }
        
        .hero-section {
            padding: 10px 0;
            margin-bottom: 30px;
        }
        
        .hero-section h2 {
            font-size: 2rem;
        }
        
        .hero-section p {
            font-size: 1rem;
        }
        
        .tags-grid {
            grid-template-columns: 1fr;
            gap: 16px;
            margin: 20px 0;
        }
        
        .tag-card {
            padding: 20px;
        }
        
        .tag-icon {
            width: 48px;
            height: 48px;
            font-size: 1.2rem;
        }
        
        .tag-name {
            font-size: 1.1rem;
        }
        
        .languages-container {
            justify-content: flex-start;
            gap: 10px;
        }
        
        .language-link {
            padding: 8px 14px;
            font-size: 0.85rem;
        }
    }

    @media (max-width: 480px) {
        .hero-section h2 {
            font-size: 1.75rem;
        }
        
        .tag-card {
            padding: 16px;
            gap: 12px;
        }
        
        .tag-icon {
            width: 40px;
            height: 40px;
            font-size: 1rem;
        }
    }
</style>
