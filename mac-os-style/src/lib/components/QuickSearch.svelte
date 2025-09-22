<script lang="ts">
    import { goto } from '$app/navigation';
    
    let {
        placeholder = "Search confessions...",
        showButton = true,
        size = "medium",
        onSearch = null
    }: {
        placeholder?: string;
        showButton?: boolean;
        size?: "small" | "medium" | "large";
        onSearch?: ((query: string) => void) | null;
    } = $props();

    let searchQuery = $state('');

    function handleSearch(event: Event) {
        event.preventDefault();
        if (!searchQuery.trim()) return;
        
        if (onSearch) {
            onSearch(searchQuery.trim());
        } else {
            // Navigate to search page with query
            goto(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
        }
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            handleSearch(event);
        }
    }

    const sizeClasses = {
        small: 'search-small',
        medium: 'search-medium', 
        large: 'search-large'
    };
</script>

<form on:submit={handleSearch} class="quick-search {sizeClasses[size]}">
    <div class="search-container">
        <input
            type="text"
            bind:value={searchQuery}
            {placeholder}
            class="search-input"
            on:keydown={handleKeydown}
        />
        {#if showButton}
            <button type="submit" class="search-btn" disabled={!searchQuery.trim()}>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="11" cy="11" r="8"/>
                    <path d="m21 21-4.35-4.35"/>
                </svg>
            </button>
        {/if}
    </div>
</form>

<style>
    .quick-search {
        display: flex;
        align-items: center;
    }

    .search-container {
        display: flex;
        align-items: center;
        background: rgba(44, 39, 51, 0.4);
        border: 1px solid rgba(255, 255, 255, 0.2);
        border-radius: 8px;
        overflow: hidden;
        transition: all 0.3s ease;
        width: 100%;
    }

    .search-container:focus-within {
        border-color: rgba(102, 126, 234, 0.5);
        box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
    }

    .search-input {
        flex: 1;
        padding: 8px 12px;
        border: none;
        background: transparent;
        color: var(--font-color);
        font-size: 0.9rem;
        outline: none;
    }

    .search-input::placeholder {
        color: var(--font-color-accent);
    }

    .search-btn {
        padding: 8px 12px;
        border: none;
        background: rgba(102, 126, 234, 0.8);
        color: white;
        cursor: pointer;
        transition: background 0.3s ease;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .search-btn:hover:not(:disabled) {
        background: rgba(102, 126, 234, 0.9);
    }

    .search-btn:disabled {
        background: rgba(255, 255, 255, 0.1);
        cursor: not-allowed;
    }

    /* Size variations */
    .search-small .search-input {
        padding: 6px 10px;
        font-size: 0.8rem;
    }

    .search-small .search-btn {
        padding: 6px 10px;
    }

    .search-small .search-btn svg {
        width: 14px;
        height: 14px;
    }

    .search-medium .search-input {
        padding: 8px 12px;
        font-size: 0.9rem;
    }

    .search-medium .search-btn {
        padding: 8px 12px;
    }

    .search-large .search-input {
        padding: 12px 16px;
        font-size: 1rem;
    }

    .search-large .search-btn {
        padding: 12px 16px;
    }

    .search-large .search-btn svg {
        width: 18px;
        height: 18px;
    }
</style>
