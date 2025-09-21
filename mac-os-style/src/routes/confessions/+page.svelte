<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import Gap from "$lib/components/Gap.svelte";

    // Mock data for now - in real app this would come from API
    const confessions = [
        {
            id: 1,
            title: "Spent 3 hours debugging... forgot semicolon",
            description: "I was debugging a complex JavaScript issue for 3 hours straight. Turns out I had a missing semicolon in a completely unrelated file that was causing the whole thing to break.",
            language: "JavaScript",
            snippet: `console.log("Hello World") // missing semicolon here
const data = fetchData();`,
            upvotes: 142,
            tags: ["JavaScript", "Debugging", "Semicolon"]
        },
        {
            id: 2,
            title: "Pushed to production on Friday... at 5 PM",
            description: "I know, I know. Never deploy on Friday. But the fix was 'simple' and the client was breathing down our necks. Spoiler alert: it wasn't simple.",
            language: "DevOps",
            snippet: `git push origin main
// Famous last words: "What could go wrong?"`,
            upvotes: 98,
            tags: ["Production", "Friday Deploy", "Murphy's Law"]
        },
        {
            id: 3,
            title: "CSS: It works! *Browser refresh* It's gone...",
            description: "Spent 2 hours getting the perfect CSS layout. Everything looked perfect. Refreshed the browser and... nothing. Turns out I was editing the wrong CSS file the entire time.",
            language: "CSS",
            snippet: `.perfect-layout {
  /* This was in the wrong file */
  display: flex;
  justify-content: center;
  align-items: center;
}`,
            upvotes: 76,
            tags: ["CSS", "Wrong File", "Frustration"]
        }
    ];
</script>

<Page title="Browse Confessions">
    <section class="confessions-container">
        <div class="hero-section">
            <h2>Coding Confessions</h2>
            <p>Learn from others' debugging adventures and coding mishaps</p>
        </div>
        
        <Gap size={20} />
        
        {#each confessions as confession}
            <div class="confession-card">
                <div class="confession-header">
                    <h3 class="confession-title">{confession.title}</h3>
                    <div class="confession-meta">
                        <span class="language-tag">{confession.language}</span>
                        <span class="upvotes">⭐ {confession.upvotes}</span>
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
                        <span class="tag">{tag}</span>
                    {/each}
                </div>
            </div>
            
            <Gap size={15} />
        {/each}
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
    
    .upvotes {
        color: var(--font-color-accent);
        font-size: 0.9rem;
        font-weight: 500;
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
    }
    
    .tag {
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color-accent);
        padding: 4px 10px;
        border-radius: 16px;
        font-size: 0.8rem;
    }
</style>
