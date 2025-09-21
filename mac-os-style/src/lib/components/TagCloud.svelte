<script lang="ts">
    import api from "$lib/api";
    import type { Tag, LoadingState } from "$lib/types";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";

    interface Props {
        title?: string;
        limit?: number;
    }
    
    let { title = "Popular Tags", limit = 20 }: Props = $props();
    
    let tags: Tag[] = $state([]);
    let loadingState: LoadingState = $state('idle');

    async function loadTags() {
        loadingState = 'loading';
        
        try {
            const allTags = await api.tags.getAll();
            tags = allTags.slice(0, limit);
            loadingState = 'success';
        } catch (error) {
            console.error('Failed to load tags:', error);
            loadingState = 'error';
        }
    }

    function handleTagClick(tag: Tag) {
        goto(`/confessions?tag=${encodeURIComponent(tag.name)}`);
    }

    onMount(() => {
        loadTags();
    });
</script>

<div class="tag-cloud-container">
    <div class="tag-cloud-header">
        <h3 class="tag-cloud-title">{title}</h3>
    </div>
    
    {#if loadingState === 'loading'}
        <div class="loading">
            <p>Loading tags...</p>
        </div>
    {:else if loadingState === 'error'}
        <div class="error">
            <p>Failed to load tags</p>
            <button onclick={loadTags} class="retry-button">Retry</button>
        </div>
    {:else if tags.length > 0}
        <div class="tag-cloud">
            {#each tags as tag}
                <button 
                    class="tag-item" 
                    onclick={() => handleTagClick(tag)}
                    title="View confessions tagged with '{tag.name}'"
                >
                    {tag.name}
                </button>
            {/each}
        </div>
    {:else}
        <div class="empty">
            <p>No tags available</p>
        </div>
    {/if}
</div>

<style>
    .tag-cloud-container {
        width: 100%;
    }
    
    .tag-cloud-header {
        display: flex;
        justify-content: flex-start;
        align-items: center;
        margin-bottom: 15px;
    }
    
    .tag-cloud-title {
        font-size: 1.2rem;
        font-weight: 600;
        color: var(--font-color);
        margin: 0;
    }
    
    .tag-cloud {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        align-items: center;
    }
    
    .tag-item {
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color-accent);
        padding: 6px 12px;
        border-radius: 16px;
        font-size: 0.85rem;
        cursor: pointer;
        transition: all 0.2s ease;
        font-weight: 500;
    }
    
    .tag-item:hover {
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
        border-color: rgba(102, 126, 234, 0.3);
        color: var(--font-color);
        transform: translateY(-1px);
    }
    
    .loading,
    .error,
    .empty {
        text-align: center;
        padding: 20px;
        color: var(--font-color-accent);
        font-size: 0.9rem;
    }
    
    .retry-button {
        margin-top: 8px;
        padding: 6px 12px;
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color);
        border-radius: 4px;
        cursor: pointer;
        font-size: 0.8rem;
        transition: all 0.2s ease;
    }
    
    .retry-button:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
    }
    
    /* Mobile responsiveness */
    @media (max-width: 768px) {
        .tag-item {
            font-size: 0.8rem;
            padding: 5px 10px;
        }
    }
</style>
