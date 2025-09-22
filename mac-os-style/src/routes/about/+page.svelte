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
        
        <div class="sections-container">
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
            
            <Group title="Why 'Shit Happens'?">
                <div class="content">
                    <p>
                        The name "Shit Happens" was inspired while watching the classic movie <em>Forrest Gump</em>. 
                        In the film, there's a memorable scene where this phrase perfectly captures life's unpredictable moments - 
                        and that's exactly what coding is like!
                    </p>
                    <p>
                        Just like Forrest's journey, our coding adventures are filled with unexpected bumps, 
                        surprising victories, and moments where we realize that sometimes, well... shit happens. 
                        And that's perfectly okay - it's all part of the learning process.
                    </p>
                </div>
            </Group>
            
            
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
            <Group title="Thank You for Using Shit Happens">
                <div class="content">
                    <p>
                        Thank you for being a part of our community! Your stories, support, and willingness to share help make this a better place for everyone. 
                        Whether you're here to laugh, learn, or lend a hand, we're glad you're with us.
                    </p>
                    <p>
                        Happy coding, and remember: every bug is just another story waiting to be told!
                    </p>
                </div>
            </Group>
        </div>
    </section>
</Page>

<style>
    .about-container {
        width: 100%;
        max-width: 800px;
        margin: 0 auto;
        padding: 0 20px;
    }
    
    .hero-section {
        text-align: center;
        margin-bottom: 40px;
        padding: 40px 20px;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
        border-radius: 20px;
        border: 1px solid rgba(102, 126, 234, 0.2);
    }
    
    .hero-section h2 {
        font-size: 2.5rem;
        font-weight: bold;
        color: var(--font-color);
        margin-bottom: 15px;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        background-clip: text;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
    }
    
    .hero-section p {
        color: var(--font-color-accent);
        font-size: 1.2rem;
        line-height: 1.6;
        max-width: 600px;
        margin: 0 auto;
    }

    .sections-container {
        display: flex;
        flex-direction: column;
        gap: 25px;
        margin-bottom: 40px;
    }

    .sections-container :global(.group) {
        background: rgba(44, 39, 51, 0.3);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 15px;
        padding: 25px;
        transition: transform 0.3s ease, box-shadow 0.3s ease;
        backdrop-filter: blur(10px);
    }

    .sections-container :global(.group):hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
        border-color: rgba(255, 255, 255, 0.15);
    }
    
    .content p {
        color: var(--font-color);
        line-height: 1.7;
        margin-bottom: 16px;
        font-size: 1rem;
    }

    .content p:last-child {
        margin-bottom: 0;
    }

    .content strong {
        color: var(--font-color-accent);
    }

    .content em {
        color: var(--color-success);
        font-style: italic;
    }
    
    .benefits-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
        gap: 20px;
        margin: 0;
    }
    
    .benefit {
        text-align: center;
        padding: 20px 15px;
        background: rgba(44, 39, 51, 0.2);
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.05);
        transition: all 0.3s ease;
    }

    .benefit:hover {
        background: rgba(44, 39, 51, 0.4);
        border-color: rgba(255, 255, 255, 0.1);
        transform: translateY(-3px);
    }
    
    .benefit h4 {
        color: var(--font-color);
        font-size: 1.1rem;
        margin-bottom: 10px;
        font-weight: 600;
    }
    
    .benefit p {
        color: var(--font-color-accent);
        font-size: 0.95rem;
        line-height: 1.5;
        margin: 0;
    }
    
    .guidelines {
        list-style: none;
        padding: 0;
        margin: 0;
    }
    
    .guidelines li {
        color: var(--font-color);
        padding: 12px 0;
        border-bottom: 1px solid rgba(255, 255, 255, 0.08);
        transition: all 0.3s ease;
        position: relative;
        padding-left: 20px;
    }

    .guidelines li::before {
        content: "•";
        position: absolute;
        left: 0;
        color: var(--font-color-accent);
        font-weight: bold;
    }
    
    .guidelines li:last-child {
        border-bottom: none;
    }

    .guidelines li:hover {
        color: var(--font-color-accent);
        padding-left: 25px;
    }

    .guidelines strong {
        color: var(--color-success);
    }

    /* Responsive Design */
    @media (max-width: 768px) {
        .about-container {
            padding: 0 15px;
        }

        .hero-section {
            padding: 30px 20px;
            margin-bottom: 30px;
        }

        .hero-section h2 {
            font-size: 2rem;
        }

        .hero-section p {
            font-size: 1.1rem;
        }

        .sections-container {
            gap: 20px;
            margin-bottom: 30px;
        }

        .sections-container :global(.group) {
            padding: 20px;
        }

        .benefits-grid {
            grid-template-columns: 1fr;
            gap: 15px;
        }
    }

    @media (max-width: 480px) {
        .benefits-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
