<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import ConfessionModal from "$lib/components/ConfessionModal.svelte";
    import type { Confession } from "$lib/types";
    import api from "$lib/api";

    // Data from page load
    let { data } = $props();
    let confessions: Confession[] = $state(data.confessions || []);
    let error = $state(data.error || '');

    // Modal state
    let selectedConfession: Confession | null = $state(null);
    let showModal = $state(false);

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

    function closeModal() {
        showModal = false;
        selectedConfession = null;
    }

    // Handle keyboard navigation
    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Escape' && showModal) {
            closeModal();
        }
    }
</script>

<svelte:window onkeydown={handleKeydown} />

<Page title="Trending This Week - Shit Happens">
    <section>
        <div class="header">
            <div class="title-section">
                <h1>🔥 Trending This Week</h1>
                <p class="subtitle">
                    The most discussed coding confessions from the past 7 days. These are the bugs and stories that have captured the developer community's attention!
                </p>
            </div>
            
            <div class="back-nav">
                <a href="/confessions" class="back-link">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="m15 18-6-6 6-6"/>
                    </svg>
                    All Confessions
                </a>
            </div>
        </div>

        {#if error}
            <Group title="Error">
                <div class="error-message">
                    <p>{error}</p>
                    <button onclick={() => window.location.reload()} class="retry-btn">
                        Try Again
                    </button>
                </div>
            </Group>
        {:else}
            <Group title="This Week's Hottest Coding Disasters ({confessions.length} confessions)">
                {#if confessions.length === 0}
                    <div class="empty-state">
                        <div class="empty-icon">📊</div>
                        <h3>No trending confessions this week yet</h3>
                        <p>Be the first to share a confession that gets the community talking!</p>
                        <a href="/submit" class="submit-link">Share Your Story</a>
                    </div>
                {:else}
                    <div class="confessions-grid">
                        {#each confessions as confession (confession.id)}
                            <article class="confession-card" onclick={(e) => openConfessionModal(confession, e)}>
                                <div class="card-header">
                                    <h2 class="confession-title">{confession.title}</h2>
                                    <div class="language-badge">{confession.language}</div>
                                </div>
                                
                                <p class="confession-description">
                                    {confession.description.length > 150 
                                        ? confession.description.substring(0, 150) + '...' 
                                        : confession.description}
                                </p>
                                
                                {#if confession.tags && confession.tags.length > 0}
                                    <div class="tags">
                                        {#each confession.tags.slice(0, 3) as tag}
                                            <span class="tag">#{tag.name}</span>
                                        {/each}
                                        {#if confession.tags.length > 3}
                                            <span class="tag-more">+{confession.tags.length - 3} more</span>
                                        {/if}
                                    </div>
                                {/if}
                                
                                <div class="card-footer">
                                    <div class="metadata">
                                        <span class="date">
                                            {new Date(confession.created_at).toLocaleDateString()}
                                        </span>
                                        {#if confession.sentiment}
                                            <span class="sentiment sentiment-{confession.sentiment}">{confession.sentiment}</span>
                                        {/if}
                                    </div>
                                    
                                    <button 
                                        class="heart-btn"
                                        onclick={(e) => {
                                            e.stopPropagation();
                                            handleUpvote(confession.id);
                                        }}
                                        aria-label="Upvote confession"
                                    >
                                        <svg class="heart-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                                        </svg>
                                        <span class="upvote-count">{confession.upvotes}</span>
                                    </button>
                                </div>
                            </article>
                        {/each}
                    </div>
                {/if}
            </Group>
        {/if}
    </section>
</Page>

<!-- Modal for detailed confession view -->
{#if showModal && selectedConfession}
    <ConfessionModal 
        {selectedConfession} 
        onClose={closeModal}
        onUpvote={handleUpvote}
    />
{/if}

<style>
    section {
        padding: 20px;
        max-width: 1200px;
        margin: 0 auto;
    }

    .header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 30px;
        gap: 20px;
    }

    .title-section h1 {
        font-size: 2.5rem;
        font-weight: 700;
        color: var(--font-color);
        margin: 0 0 10px 0;
        background: linear-gradient(135deg, #ff6b6b, #ffa500);
        background-clip: text;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
    }

    .subtitle {
        color: var(--font-color-accent);
        font-size: 1.1rem;
        line-height: 1.5;
        margin: 0;
        max-width: 600px;
    }

    .back-nav {
        flex-shrink: 0;
    }

    .back-link {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        padding: 10px 16px;
        background: rgba(44, 39, 51, 0.3);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        color: var(--font-color-accent);
        text-decoration: none;
        transition: all 0.2s ease;
        font-size: 0.95rem;
    }

    .back-link:hover {
        background: rgba(44, 39, 51, 0.5);
        border-color: rgba(255, 255, 255, 0.2);
        color: var(--font-color);
        transform: translateX(-2px);
    }

    .error-message {
        text-align: center;
        padding: 40px 20px;
    }

    .error-message p {
        color: var(--color-danger);
        font-size: 1.1rem;
        margin-bottom: 20px;
    }

    .retry-btn {
        background: var(--color-danger);
        color: white;
        border: none;
        padding: 10px 20px;
        border-radius: 8px;
        cursor: pointer;
        font-size: 1rem;
        transition: background 0.2s ease;
    }

    .retry-btn:hover {
        background: #dc2626;
    }

    .empty-state {
        text-align: center;
        padding: 60px 20px;
    }

    .empty-icon {
        font-size: 4rem;
        margin-bottom: 20px;
        opacity: 0.7;
    }

    .empty-state h3 {
        color: var(--font-color);
        font-size: 1.5rem;
        margin: 0 0 10px 0;
    }

    .empty-state p {
        color: var(--font-color-accent);
        font-size: 1.1rem;
        margin: 0 0 30px 0;
    }

    .submit-link {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        background: var(--color-success);
        color: white;
        text-decoration: none;
        padding: 12px 24px;
        border-radius: 8px;
        font-weight: 600;
        transition: all 0.2s ease;
    }

    .submit-link:hover {
        background: #218838;
        transform: translateY(-1px);
    }

    .confessions-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
        gap: 25px;
        padding: 20px 0;
    }

    .confession-card {
        background: rgba(44, 39, 51, 0.4);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 15px;
        padding: 25px;
        cursor: pointer;
        transition: all 0.3s ease;
        backdrop-filter: blur(10px);
    }

    .confession-card:hover {
        transform: translateY(-5px);
        border-color: rgba(255, 255, 255, 0.2);
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
    }

    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 15px;
        margin-bottom: 15px;
    }

    .confession-title {
        color: var(--font-color);
        font-size: 1.3rem;
        font-weight: 600;
        margin: 0;
        line-height: 1.3;
        flex: 1;
    }

    .language-badge {
        background: linear-gradient(135deg, #667eea, #764ba2);
        color: white;
        padding: 4px 10px;
        border-radius: 12px;
        font-size: 0.8rem;
        font-weight: 500;
        flex-shrink: 0;
    }

    .confession-description {
        color: var(--font-color-accent);
        line-height: 1.6;
        margin: 0 0 15px 0;
    }

    .tags {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        margin-bottom: 20px;
    }

    .tag {
        background: rgba(102, 126, 234, 0.2);
        color: var(--color-success);
        padding: 4px 10px;
        border-radius: 12px;
        font-size: 0.8rem;
        font-weight: 500;
    }

    .tag-more {
        background: rgba(255, 255, 255, 0.1);
        color: var(--font-color-accent);
        padding: 4px 10px;
        border-radius: 12px;
        font-size: 0.8rem;
    }

    .card-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-top: 15px;
        border-top: 1px solid rgba(255, 255, 255, 0.1);
    }

    .metadata {
        display: flex;
        gap: 15px;
        align-items: center;
    }

    .date {
        color: var(--font-color-accent);
        font-size: 0.9rem;
    }

    .sentiment {
        padding: 3px 8px;
        border-radius: 10px;
        font-size: 0.8rem;
        font-weight: 500;
        text-transform: capitalize;
    }

    .sentiment-positive { background: rgba(34, 197, 94, 0.2); color: var(--color-success); }
    .sentiment-negative { background: rgba(239, 68, 68, 0.2); color: var(--color-danger); }
    .sentiment-neutral { background: rgba(156, 163, 175, 0.2); color: var(--font-color-accent); }

    .heart-btn {
        display: flex;
        align-items: center;
        gap: 6px;
        background: transparent;
        border: 2px solid rgba(255, 255, 255, 0.1);
        border-radius: 20px;
        padding: 6px 12px;
        color: var(--font-color-accent);
        cursor: pointer;
        transition: all 0.2s ease;
        font-size: 0.9rem;
    }

    .heart-btn:hover {
        border-color: #e91e63;
        color: #e91e63;
        transform: scale(1.05);
    }

    .heart-btn:active {
        transform: scale(0.95);
    }

    .heart-icon {
        width: 16px;
        height: 16px;
        transition: all 0.2s ease;
    }

    .heart-btn:hover .heart-icon {
        fill: #e91e63;
        stroke: #e91e63;
        animation: heartPulse 0.6s ease;
    }

    @keyframes heartPulse {
        0% { transform: scale(1); }
        50% { transform: scale(1.2); }
        100% { transform: scale(1); }
    }

    .upvote-count {
        font-weight: 600;
        min-width: 20px;
        text-align: center;
    }

    /* Responsive design */
    @media (max-width: 768px) {
        .header {
            flex-direction: column;
            align-items: stretch;
        }

        .title-section h1 {
            font-size: 2rem;
        }

        .confessions-grid {
            grid-template-columns: 1fr;
            gap: 20px;
        }

        .card-header {
            flex-direction: column;
            align-items: flex-start;
            gap: 10px;
        }
    }
</style>
