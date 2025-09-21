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
        
        <Gap size={30} />
        
        <Group title="Popular Programming Languages">
            <div class="languages-container">
                {#each ["javascript", "python", "css", "html", "java", "typescript", "go", "rust", "php", "ruby", "sql", "devops"] as language}
                    <a href="/confessions?language={language}" class="language-link">{language}</a>
                {/each}
            </div>
        </Group>
    </section>
</Page>

<style>
    .categories-container {
        width: 100%;
        max-width: 900px;
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
    
    .tags-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
        gap: 20px;
        margin-bottom: 30px;
    }
    
    .tag-card {
        display: flex;
        align-items: center;
        gap: 15px;
        padding: 20px;
        background: var(--bg-card);
        border: 1px solid var(--border-card);
        border-radius: 12px;
        text-decoration: none;
        transition: all 0.3s ease;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
    
    .tag-card:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
        border-color: var(--border-card-hover, var(--border-card));
    }
    
    .tag-icon {
        width: 50px;
        height: 50px;
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-weight: bold;
        font-size: 1.2rem;
        flex-shrink: 0;
    }
    
    .tag-info {
        flex: 1;
    }
    
    .tag-name {
        margin: 0 0 5px 0;
        font-size: 1.2rem;
        font-weight: 600;
        color: var(--font-color);
    }
    
    .tag-description {
        margin: 0;
        color: var(--font-color-accent);
        font-size: 0.9rem;
    }
    
    .languages-container {
        display: flex;
        flex-wrap: wrap;
        gap: 12px;
    }
    
    .language-link {
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color);
        padding: 8px 16px;
        border-radius: 20px;
        text-decoration: none;
        font-size: 0.9rem;
        transition: all 0.2s ease;
        text-transform: capitalize;
    }
    
    .language-link:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
        color: var(--font-color);
    }
</style>
