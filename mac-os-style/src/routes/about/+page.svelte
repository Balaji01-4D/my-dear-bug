<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import Gap from "$lib/components/Gap.svelte";
</script>

<Page title="About">
    <section class="about-container">
        <div class="hero-section">
            <h2>About Shit Happens</h2>
            <p>Where developers share their most embarrassing coding moments and help others learn from their mistakes</p>
        </div>
        
        <Gap size={20} />
        
        <Group title="Our Mission">
            <div class="content">
                <p>
                    <strong>Shit Happens</strong> is a safe space for developers to share their coding confessions, debugging nightmares, and production disasters. 
                    We believe that sharing our mistakes helps the entire developer community learn and grow.
                </p>
                <p>
                    Every developer has that one story - the time they deleted the production database, the semicolon that took 6 hours to find, 
                    or the CSS that worked perfectly until it didn't. These stories are valuable learning experiences disguised as embarrassing moments.
                </p>
            </div>
        </Group>
        
        <Gap size={15} />
        
        <Group title="Why Share Your Confessions?">
            <div class="benefits-grid">
                <div class="benefit">
                    <h4>🤝 Help Others</h4>
                    <p>Your mistake might save another developer hours of debugging</p>
                </div>
                <div class="benefit">
                    <h4>📚 Learn Together</h4>
                    <p>Read about common pitfalls and how others solved them</p>
                </div>
                <div class="benefit">
                    <h4>😌 Feel Better</h4>
                    <p>Realize you're not alone - we all make mistakes</p>
                </div>
                <div class="benefit">
                    <h4>🎯 Share Knowledge</h4>
                    <p>Turn your painful experience into valuable wisdom</p>
                </div>
            </div>
        </Group>
        
        <Gap size={15} />
        
        <Group title="Community Guidelines">
            <div class="content">
                <ul class="guidelines">
                    <li><strong>Be Honest:</strong> Share real experiences, not made-up stories</li>
                    <li><strong>Be Helpful:</strong> Include how you solved the problem (if you did)</li>
                    <li><strong>Be Kind:</strong> We're all learning - no judgment here</li>
                    <li><strong>Be Specific:</strong> Technical details help others learn</li>
                    <li><strong>Be Safe:</strong> Don't share sensitive company information</li>
                </ul>
            </div>
        </Group>
        
        <Gap size={15} />
        
        <Group title="Statistics">
            <div class="stats-grid">
                <div class="stat">
                    <div class="stat-number">1,247</div>
                    <div class="stat-label">Confessions Shared</div>
                </div>
                <div class="stat">
                    <div class="stat-number">589</div>
                    <div class="stat-label">Developers Helped</div>
                </div>
                <div class="stat">
                    <div class="stat-number">∞</div>
                    <div class="stat-label">Bugs Squashed</div>
                </div>
            </div>
        </Group>
        
        <Gap size={20} />
        
        <div class="call-to-action">
            <h3>Ready to Share Your Story?</h3>
            <p>Join thousands of developers who've found healing through confession</p>
            <a href="/submit" class="cta-button">Submit Your Confession</a>
        </div>
    </section>
</Page>

<style>
    .about-container {
        width: 100%;
        max-width: 700px;
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
        line-height: 1.5;
    }
    
    .content p {
        color: var(--font-color);
        line-height: 1.6;
        margin-bottom: 15px;
    }
    
    .benefits-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 20px;
        margin: 15px 0;
    }
    
    .benefit {
        text-align: center;
    }
    
    .benefit h4 {
        color: var(--font-color);
        font-size: 1.1rem;
        margin-bottom: 8px;
    }
    
    .benefit p {
        color: var(--font-color-accent);
        font-size: 0.9rem;
        line-height: 1.4;
        margin: 0;
    }
    
    .guidelines {
        list-style: none;
        padding: 0;
    }
    
    .guidelines li {
        color: var(--font-color);
        padding: 8px 0;
        border-bottom: 1px solid var(--border-card);
    }
    
    .guidelines li:last-child {
        border-bottom: none;
    }
    
    .stats-grid {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 20px;
        margin: 20px 0;
    }
    
    .stat {
        text-align: center;
    }
    
    .stat-number {
        font-size: 2.5rem;
        font-weight: bold;
        color: var(--font-color);
        line-height: 1;
    }
    
    .stat-label {
        color: var(--font-color-accent);
        font-size: 0.9rem;
        margin-top: 5px;
    }
    
    .call-to-action {
        text-align: center;
        padding: 30px 20px;
        background: var(--bg-card);
        border: 1px solid var(--border-card);
        border-radius: 12px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
    
    .call-to-action h3 {
        color: var(--font-color);
        font-size: 1.5rem;
        margin-bottom: 10px;
    }
    
    .call-to-action p {
        color: var(--font-color-accent);
        margin-bottom: 20px;
    }
    
    .cta-button {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        padding: 12px 24px;
        border-radius: 8px;
        text-decoration: none;
        font-weight: 600;
        transition: all 0.3s ease;
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        display: inline-block;
    }
    
    .cta-button:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
        text-decoration: none;
        color: white;
    }
</style>
