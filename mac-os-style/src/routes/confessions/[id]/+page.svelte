<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import ConfessionModal from "$lib/components/ConfessionModal.svelte";
    import Gap from "$lib/components/Gap.svelte";
    import api from "$lib/api";
    import { goto } from "$app/navigation";
    import type { PageData } from './$types';
    import type { Confession } from '$lib/types';

    let { data }: { data: PageData } = $props();
    let confession: Confession = $state(data.confession);

    async function handleUpvote(confessionId: number) {
        try {
            await api.confessions.upvote(confessionId);
            // Update the upvote count locally
            confession = { ...confession, upvotes: confession.upvotes + 1 };
        } catch (error) {
            alert(api.handleError(error as Error));
        }
    }

    function handleClose() {
        // Navigate back to confessions page
        goto('/confessions');
    }
</script>

<svelte:head>
    <title>{confession.title} | Shit Happens - Coding Confessions</title>
    <meta name="description" content={confession.description.substring(0, 160)} />
    <meta property="og:title" content={confession.title} />
    <meta property="og:description" content={confession.description.substring(0, 160)} />
    <meta property="og:type" content="article" />
</svelte:head>

<!-- Display the confession in a modal-like view but as a full page -->
<div class="confession-page">
    <div class="page-container">
        <div class="breadcrumb">
            <button onclick={() => goto('/confessions')} class="back-button">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M19 12H5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M12 19l-7-7 7-7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                Back to Confessions
            </button>
        </div>
        
        <Gap size={20} />
        
        <ConfessionModal 
            {confession} 
            onClose={handleClose}
            onUpvote={handleUpvote}
        />
    </div>
</div>

<style>
    .confession-page {
        min-height: 100vh;
        background: var(--bg-main);
        padding: 20px;
    }
    
    .page-container {
        max-width: 800px;
        margin: 0 auto;
    }
    
    .breadcrumb {
        margin-bottom: 20px;
    }
    
    .back-button {
        display: flex;
        align-items: center;
        gap: 8px;
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color);
        padding: 10px 16px;
        border-radius: 8px;
        cursor: pointer;
        font-size: 0.95rem;
        transition: all 0.2s ease;
    }
    
    .back-button:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
        color: #667eea;
    }
    
    /* Override modal styles for full page display */
    :global(.modal-backdrop) {
        position: static !important;
        background: transparent !important;
        backdrop-filter: none !important;
        padding: 0 !important;
    }
    
    :global(.modal-content) {
        max-height: none !important;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1) !important;
        animation: none !important;
    }
</style>
