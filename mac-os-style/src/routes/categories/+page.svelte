<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import Gap from "$lib/components/Gap.svelte";
    
    const categories = [
        { name: "JavaScript Nightmares", count: 234, color: "#f7df1e" },
        { name: "CSS Horror Stories", count: 189, color: "#1572b6" },
        { name: "Python Mishaps", count: 156, color: "#3776ab" },
        { name: "Database Disasters", count: 145, color: "#336791" },
        { name: "Production Failures", count: 123, color: "#ff6b6b" },
        { name: "Debugging Adventures", count: 98, color: "#4ecdc4" },
        { name: "Git Conflicts", count: 87, color: "#f05032" },
        { name: "API Nightmares", count: 76, color: "#61dafb" },
        { name: "Memory Leaks", count: 65, color: "#ff9800" },
        { name: "Semicolon Sins", count: 54, color: "#9c27b0" }
    ];
</script>

<Page title="Categories">
    <section class="categories-container">
        <div class="hero-section">
            <h2>Browse by Category</h2>
            <p>Explore confessions organized by programming language, type of mistake, or technology</p>
        </div>
        
        <Gap size={20} />
        
        <div class="categories-grid">
            {#each categories as category}
                <a href="/category/{category.name.toLowerCase().replace(/ /g, '-')}" class="category-card">
                    <div class="category-icon" style="background-color: {category.color}">
                        {category.name.charAt(0)}
                    </div>
                    <div class="category-info">
                        <h3 class="category-name">{category.name}</h3>
                        <p class="category-count">{category.count} confessions</p>
                    </div>
                </a>
            {/each}
        </div>
        
        <Gap size={30} />
        
        <Group title="Popular Tags">
            <div class="tags-container">
                {#each ["debugging", "production", "javascript", "css", "python", "database", "git", "api", "performance", "security", "testing", "deployment"] as tag}
                    <a href="/tag/{tag}" class="tag-link">#{tag}</a>
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
    
    .categories-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
        gap: 15px;
        margin-bottom: 20px;
    }
    
    .category-card {
        display: flex;
        align-items: center;
        gap: 15px;
        padding: 20px;
        background: var(--bg-card);
        border: 1px solid var(--border-card);
        border-radius: 12px;
        text-decoration: none;
        color: var(--font-color);
        transition: all 0.2s ease;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
    
    .category-card:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
        text-decoration: none;
        color: var(--font-color);
    }
    
    .category-icon {
        width: 50px;
        height: 50px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-weight: bold;
        font-size: 1.5rem;
        flex-shrink: 0;
    }
    
    .category-info {
        flex: 1;
    }
    
    .category-name {
        font-size: 1.2rem;
        font-weight: 600;
        margin: 0 0 5px 0;
        color: var(--font-color);
    }
    
    .category-count {
        color: var(--font-color-accent);
        margin: 0;
        font-size: 0.9rem;
    }
    
    .tags-container {
        display: flex;
        flex-wrap: wrap;
        gap: 12px;
        padding: 10px 0;
    }
    
    .tag-link {
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color-accent);
        padding: 8px 14px;
        border-radius: 20px;
        text-decoration: none;
        font-size: 0.9rem;
        transition: all 0.2s ease;
    }
    
    .tag-link:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
        color: var(--font-color);
        transform: translateY(-1px);
        text-decoration: none;
    }
</style>
