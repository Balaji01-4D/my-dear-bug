<script lang="ts">
    import type { Confession } from "$lib/types";
    import api from "$lib/api";
    import Gap from "$lib/components/Gap.svelte";
    
    interface Props {
        confession: Confession;
        onClose: () => void;
        onUpvote?: (id: number) => void;
    }
    
    let { confession = $bindable(), onClose, onUpvote }: Props = $props();
    let isUpvoting = $state(false);
    let hasUpvoted = $state(false);
    
    async function handleUpvote() {
        if (isUpvoting || hasUpvoted) return;
        
        isUpvoting = true;
        try {
            await api.confessions.upvote(confession.id);
            confession.upvotes += 1;
            hasUpvoted = true;
            onUpvote?.(confession.id);
        } catch (error) {
            alert(api.handleError(error as Error));
        } finally {
            isUpvoting = false;
        }
    }
    
    function handleBackdropClick(event: MouseEvent) {
        if (event.target === event.currentTarget) {
            onClose();
        }
    }
    
    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Escape') {
            onClose();
        }
    }
</script>

<svelte:window onkeydown={handleKeyDown} />

<div class="modal-backdrop" onclick={handleBackdropClick}>
    <div class="modal-content">
        <div class="modal-header">
            <div class="confession-header">
                <h2 class="confession-title">{confession.title}</h2>
                <div class="confession-actions">
                    <span class="language-tag">{confession.language}</span>
                    <button class="close-button" onclick={onClose}>
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M18 6L6 18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            <path d="M6 6L18 18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                        </svg>
                    </button>
                </div>
            </div>
        </div>
        
        <div class="modal-body">
            <p class="confession-description">{confession.description}</p>
            
            {#if confession.snippet}
                <Gap size={15} />
                <div class="code-snippet">
                    <div class="code-header">
                        <span class="code-language">{confession.language}</span>
                        <button class="copy-button" onclick={() => navigator.clipboard.writeText(confession.snippet)}>
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <rect x="9" y="9" width="13" height="13" rx="2" ry="2" stroke="currentColor" stroke-width="2" fill="none"/>
                                <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" stroke="currentColor" stroke-width="2" fill="none"/>
                            </svg>
                            Copy
                        </button>
                    </div>
                    <pre><code>{confession.snippet}</code></pre>
                </div>
            {/if}
            
            {#if confession.tags && confession.tags.length > 0}
                <Gap size={15} />
                <div class="confession-tags">
                    {#each confession.tags as tag}
                        <span class="tag">{tag.name}</span>
                    {/each}
                </div>
            {/if}
            
            <Gap size={20} />
            
            <div class="confession-footer">
                <div class="confession-meta">
                    <span class="created-at">
                        Posted on {new Date(confession.createdAt).toLocaleDateString('en-US', { 
                            year: 'numeric', 
                            month: 'long', 
                            day: 'numeric' 
                        })}
                    </span>
                    {#if confession.sentiment}
                        <span class="sentiment-badge sentiment-{confession.sentiment}">
                            {confession.sentiment}
                        </span>
                    {/if}
                </div>
                
                <div class="actions">
                    <button 
                        class="heart-button"
                        class:liked={hasUpvoted}
                        class:loading={isUpvoting}
                        onclick={handleUpvote}
                        disabled={isUpvoting || hasUpvoted}
                    >
                        <svg class="heart-icon" width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" 
                                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                        </svg>
                        <span class="upvote-count">{confession.upvotes}</span>
                    </button>
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.75);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
        padding: 20px;
        backdrop-filter: blur(4px);
    }
    
    .modal-content {
        background: var(--bg-card);
        border: 1px solid var(--border-card);
        border-radius: 16px;
        width: 100%;
        max-width: 700px;
        max-height: 90vh;
        overflow: hidden;
        box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
        animation: modalSlideIn 0.3s ease-out;
    }
    
    @keyframes modalSlideIn {
        from {
            opacity: 0;
            transform: translateY(-20px) scale(0.95);
        }
        to {
            opacity: 1;
            transform: translateY(0) scale(1);
        }
    }
    
    .modal-header {
        padding: 24px 24px 0 24px;
        border-bottom: 1px solid var(--border-card);
    }
    
    .confession-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 20px;
        padding-bottom: 20px;
    }
    
    .confession-title {
        font-size: 1.8rem;
        font-weight: 700;
        color: var(--font-color);
        margin: 0;
        flex: 1;
        line-height: 1.3;
    }
    
    .confession-actions {
        display: flex;
        align-items: center;
        gap: 15px;
        flex-shrink: 0;
    }
    
    .language-tag {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        padding: 6px 16px;
        border-radius: 20px;
        font-size: 0.9rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }
    
    .close-button {
        background: none;
        border: none;
        color: var(--font-color-accent);
        cursor: pointer;
        padding: 8px;
        border-radius: 8px;
        transition: all 0.2s ease;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    
    .close-button:hover {
        background: var(--bg-input);
        color: var(--font-color);
    }
    
    .modal-body {
        padding: 24px;
        overflow-y: auto;
        max-height: calc(90vh - 120px);
    }
    
    .confession-description {
        color: var(--font-color);
        line-height: 1.7;
        font-size: 1.1rem;
        margin: 0 0 20px 0;
        white-space: pre-wrap;
    }
    
    .code-snippet {
        background: #1e1e1e;
        border: 1px solid #444;
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }
    
    .code-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 16px;
        background: #2d2d2d;
        border-bottom: 1px solid #444;
    }
    
    .code-language {
        color: #f8f8f2;
        font-size: 0.9rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }
    
    .copy-button {
        background: #444;
        border: none;
        color: #f8f8f2;
        padding: 6px 10px;
        border-radius: 6px;
        font-size: 0.8rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 6px;
        transition: background 0.2s ease;
    }
    
    .copy-button:hover {
        background: #555;
    }
    
    .code-snippet pre {
        margin: 0;
        padding: 20px;
        color: #f8f8f2;
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 0.95rem;
        line-height: 1.5;
        overflow-x: auto;
    }
    
    .confession-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 10px;
    }
    
    .tag {
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color-accent);
        padding: 8px 14px;
        border-radius: 20px;
        font-size: 0.9rem;
        font-weight: 500;
        transition: all 0.2s ease;
    }
    
    .tag:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
        color: var(--font-color);
    }
    
    .confession-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-top: 20px;
        border-top: 1px solid var(--border-card);
    }
    
    .confession-meta {
        display: flex;
        align-items: center;
        gap: 15px;
    }
    
    .created-at {
        color: var(--font-color-accent);
        font-size: 0.95rem;
    }
    
    .sentiment-badge {
        padding: 4px 12px;
        border-radius: 16px;
        font-size: 0.8rem;
        font-weight: 600;
        text-transform: capitalize;
    }
    
    .sentiment-happy {
        background: rgba(34, 197, 94, 0.2);
        color: #22c55e;
        border: 1px solid rgba(34, 197, 94, 0.3);
    }
    
    .sentiment-sad {
        background: rgba(239, 68, 68, 0.2);
        color: #ef4444;
        border: 1px solid rgba(239, 68, 68, 0.3);
    }
    
    .sentiment-neutral {
        background: rgba(156, 163, 175, 0.2);
        color: #9ca3af;
        border: 1px solid rgba(156, 163, 175, 0.3);
    }
    
    .actions {
        display: flex;
        align-items: center;
        gap: 15px;
    }
    
    .heart-button {
        background: none;
        border: 2px solid var(--border-input);
        color: var(--font-color-accent);
        padding: 10px 16px;
        border-radius: 50px;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 1rem;
        font-weight: 600;
        transition: all 0.3s ease;
        position: relative;
        overflow: hidden;
    }
    
    .heart-button:hover:not(:disabled) {
        border-color: #e91e63;
        color: #e91e63;
        transform: scale(1.05);
    }
    
    .heart-button.liked {
        background: linear-gradient(135deg, #e91e63, #f06292);
        border-color: #e91e63;
        color: white;
        animation: heartPulse 0.6s ease-out;
    }
    
    .heart-button.liked .heart-icon path {
        fill: currentColor;
        stroke: currentColor;
    }
    
    .heart-button:disabled {
        cursor: not-allowed;
        opacity: 0.7;
    }
    
    .heart-button.loading {
        animation: pulse 1.5s infinite;
    }
    
    @keyframes heartPulse {
        0% { transform: scale(1); }
        50% { transform: scale(1.2); }
        100% { transform: scale(1); }
    }
    
    @keyframes pulse {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.7; }
    }
    
    .heart-icon {
        transition: all 0.3s ease;
    }
    
    .upvote-count {
        font-weight: 700;
        min-width: 20px;
        text-align: left;
    }
    
    /* Mobile responsiveness */
    @media (max-width: 768px) {
        .modal-backdrop {
            padding: 10px;
        }
        
        .modal-content {
            max-height: 95vh;
        }
        
        .confession-title {
            font-size: 1.5rem;
        }
        
        .confession-header {
            flex-direction: column;
            gap: 15px;
            align-items: stretch;
        }
        
        .confession-actions {
            justify-content: space-between;
        }
        
        .confession-footer {
            flex-direction: column;
            gap: 15px;
            align-items: stretch;
        }
        
        .actions {
            justify-content: center;
        }
    }
</style>
