<script lang="ts">
    import { goto } from '$app/navigation';
    
    let showSearchModal = $state(false);
    let searchQuery = $state('');

    function openSearch() {
        showSearchModal = true;
        // Focus the input after a short delay to ensure modal is rendered
        setTimeout(() => {
            const input = document.querySelector('.modal-search-input') as HTMLInputElement;
            if (input) input.focus();
        }, 100);
    }

    function closeSearch() {
        showSearchModal = false;
        searchQuery = '';
    }

    function handleSearch(event: Event) {
        event.preventDefault();
        if (!searchQuery.trim()) return;
        
        closeSearch();
        goto(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Escape') {
            closeSearch();
        } else if (event.key === 'Enter') {
            handleSearch(event);
        }
    }

    // Global keyboard shortcut (Ctrl/Cmd + K)
    function handleGlobalKeydown(event: KeyboardEvent) {
        if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
            event.preventDefault();
            openSearch();
        }
    }

    // Add global listener
    if (typeof document !== 'undefined') {
        document.addEventListener('keydown', handleGlobalKeydown);
    }
</script>

<!-- Floating Search Button -->
<button class="floating-search-btn" on:click={openSearch} title="Search (Ctrl+K)">
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8"/>
        <path d="m21 21-4.35-4.35"/>
    </svg>
</button>

<!-- Search Modal -->
{#if showSearchModal}
    <div class="search-modal-overlay" on:click={closeSearch} role="dialog" aria-modal="true">
        <div class="search-modal" on:click={(e) => e.stopPropagation()}>
            <div class="search-modal-header">
                <h3>Search Bug Confessions</h3>
                <button class="close-btn" on:click={closeSearch}>
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="18" y1="6" x2="6" y2="18"/>
                        <line x1="6" y1="6" x2="18" y2="18"/>
                    </svg>
                </button>
            </div>
            
            <form on:submit={handleSearch} class="modal-search-form">
                <input
                    type="text"
                    bind:value={searchQuery}
                    placeholder="Search for debugging stories, error messages, solutions..."
                    class="modal-search-input"
                    on:keydown={handleKeydown}
                />
                <div class="search-shortcuts">
                    <span class="shortcut">Enter to search</span>
                    <span class="shortcut">Esc to close</span>
                </div>
            </form>
            
            <div class="quick-links">
                <a href="/search" class="quick-link">Advanced Search</a>
                <a href="/categories" class="quick-link">Browse Categories</a>
                <a href="/confessions?filter=top" class="quick-link">Top Stories</a>
            </div>
        </div>
    </div>
{/if}

<style>
    .floating-search-btn {
        position: fixed;
        bottom: 20px;
        right: 20px;
        width: 56px;
        height: 56px;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.9) 0%, rgba(118, 75, 162, 0.9) 100%);
        border: none;
        border-radius: 50%;
        color: white;
        cursor: pointer;
        box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
        transition: all 0.3s ease;
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
    }

    .floating-search-btn:hover {
        transform: scale(1.1);
        box-shadow: 0 6px 25px rgba(102, 126, 234, 0.5);
    }

    .search-modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.7);
        backdrop-filter: blur(10px);
        display: flex;
        align-items: flex-start;
        justify-content: center;
        padding: 10vh 20px 20px;
        z-index: 2000;
        animation: fadeIn 0.2s ease;
    }

    @keyframes fadeIn {
        from { opacity: 0; }
        to { opacity: 1; }
    }

    .search-modal {
        background: rgba(44, 39, 51, 0.95);
        border: 1px solid rgba(255, 255, 255, 0.2);
        border-radius: 16px;
        padding: 30px;
        width: 100%;
        max-width: 600px;
        box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
        animation: slideUp 0.3s ease;
    }

    @keyframes slideUp {
        from { 
            opacity: 0;
            transform: translateY(30px);
        }
        to { 
            opacity: 1;
            transform: translateY(0);
        }
    }

    .search-modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 25px;
    }

    .search-modal-header h3 {
        margin: 0;
        color: var(--font-color);
        font-size: 1.5rem;
        font-weight: 600;
    }

    .close-btn {
        background: none;
        border: none;
        color: var(--font-color-accent);
        cursor: pointer;
        padding: 5px;
        border-radius: 6px;
        transition: all 0.3s ease;
    }

    .close-btn:hover {
        background: rgba(255, 255, 255, 0.1);
        color: var(--font-color);
    }

    .modal-search-form {
        margin-bottom: 25px;
    }

    .modal-search-input {
        width: 100%;
        padding: 16px 20px;
        border: 2px solid rgba(255, 255, 255, 0.1);
        border-radius: 12px;
        background: rgba(44, 39, 51, 0.3);
        color: var(--font-color);
        font-size: 1.1rem;
        transition: all 0.3s ease;
        box-sizing: border-box;
    }

    .modal-search-input:focus {
        outline: none;
        border-color: rgba(102, 126, 234, 0.5);
        box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
    }

    .modal-search-input::placeholder {
        color: var(--font-color-accent);
    }

    .search-shortcuts {
        display: flex;
        justify-content: space-between;
        margin-top: 12px;
    }

    .shortcut {
        font-size: 0.85rem;
        color: var(--font-color-accent);
        background: rgba(255, 255, 255, 0.05);
        padding: 4px 8px;
        border-radius: 6px;
        border: 1px solid rgba(255, 255, 255, 0.1);
    }

    .quick-links {
        display: flex;
        gap: 15px;
        flex-wrap: wrap;
    }

    .quick-link {
        color: rgba(102, 126, 234, 0.8);
        text-decoration: none;
        font-size: 0.95rem;
        padding: 8px 12px;
        border-radius: 8px;
        background: rgba(102, 126, 234, 0.1);
        border: 1px solid rgba(102, 126, 234, 0.2);
        transition: all 0.3s ease;
    }

    .quick-link:hover {
        background: rgba(102, 126, 234, 0.2);
        border-color: rgba(102, 126, 234, 0.4);
        color: rgba(102, 126, 234, 0.9);
    }

    /* Hide on mobile to avoid overlap with navigation */
    @media (max-width: 768px) {
        .floating-search-btn {
            display: none;
        }
        
        .search-modal-overlay {
            padding: 5vh 10px 10px;
        }
        
        .search-modal {
            padding: 20px;
        }
        
        .quick-links {
            flex-direction: column;
        }
    }
</style>
