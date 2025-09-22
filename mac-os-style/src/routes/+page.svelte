<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import logo from "$lib/images/ghost.svg";
    import sync from "$lib/images/tabs/sync.webp";
    import fonts from "$lib/images/tabs/font-playground.webp";
    import Admonition from "$lib/components/Admonition.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import LinkItem from "$lib/components/settings/LinkItem.svelte";
    import Separator from "$lib/components/settings/Separator.svelte";
    import TagCloud from "$lib/components/TagCloud.svelte";
    import QuickSearch from "$lib/components/QuickSearch.svelte";

    // Collection of funny and inspiring quotes about coding mistakes
    const quotes = [
        {
        "text": "Weeks of programming can save you hours of planning.",
        "author": "Unknown",
        "type": "funny"
    
        },
        {
            text: "Code never lies, comments sometimes do.",
            author: "Ron Jeffries",
            type: "inspiring"
        },
        {
            text: "The best error message is the one that never shows up.",
            author: "Thomas Fuchs",
            type: "inspiring"
        },
        {
            text: "It works on my machine.",
            author: "Every Developer Ever",
            type: "funny"
        },
        {
            text: "First, solve the problem. Then, write the code.",
            author: "John Johnson",
            type: "inspiring"
        },
        {
            text: "Programming is like sex: one mistake and you have to support it for the rest of your life.",
            author: "Michael Sinz",
            type: "funny"
        },
        {
            text: "I don't care if it works on your machine! We are not shipping your machine!",
            author: "Vidiu Platon",
            type: "funny"
        }
    ];

    let currentQuoteIndex = $state(0);
    let isAnimating = $state(false);

    // Rotate quotes every 5 seconds
    $effect(() => {
        const interval = setInterval(() => {
            isAnimating = true;
            setTimeout(() => {
                currentQuoteIndex = (currentQuoteIndex + 1) % quotes.length;
                isAnimating = false;
            }, 300);
        }, 5000);

        return () => clearInterval(interval);
    });

    const currentQuote = $derived(quotes[currentQuoteIndex]);
</script>

<Page title="Shit Happens">
    <section>
        <!-- Hero Section -->
        <div class="hero">
            <div class="user">
                <div class="user-avatar">
                    <img src={logo} alt="Bug Icon" />
                </div>
                <div class="user-label">
                    <div class="user-name">Shit Happens</div>
                    <div class="user-subtext">Where Coding Disasters Become Learning Treasures</div>
                </div>
            </div>

            <!-- Rotating Quote Section -->
            <div class="quote-section">
                <div class="quote-container" class:animating={isAnimating}>
                    <div class="quote-content">
                        <blockquote class="quote-text {currentQuote.type}">
                            "{currentQuote.text}"
                        </blockquote>
                        <cite class="quote-author">— {currentQuote.author}</cite>
                    </div>
                </div>
                <div class="quote-indicators">
                    {#each quotes as quote, index}
                        <button 
                            class="indicator" 
                            class:active={index === currentQuoteIndex}
                            on:click={() => {
                                isAnimating = true;
                                setTimeout(() => {
                                    currentQuoteIndex = index;
                                    isAnimating = false;
                                }, 300);
                            }}
                            aria-label="Quote {index + 1}"
                        ></button>
                    {/each}
                </div>
            </div>
        </div>

        <!-- Call to Action -->
        <div class="cta-section">
            <h2 class="cta-title">Share Your Development Journey</h2>
            <p class="cta-subtitle">Join a community of experienced developers who transform challenging debugging experiences into valuable learning resources.</p>
            
            <!-- Quick Search -->
            <div class="quick-search-container">
                <QuickSearch 
                    placeholder="Search bug confessions, solutions, and stories..."
                    size="large"
                />
            </div>
            
            <div class="cta-buttons">
                <a href="/submit" class="primary-cta">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M12 19l7-7 3 3-7 7-3-3z"/>
                        <path d="m18 13-1.5-7.5L2 2l3.5 14.5L13 18l5-5z"/>
                        <path d="m2 2 7.586 7.586"/>
                        <circle cx="11" cy="11" r="2"/>
                    </svg>
                    Submit Your Experience
                </a>
                <a href="/confessions" class="secondary-cta">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="11" cy="11" r="8"/>
                        <path d="m21 21-4.35-4.35"/>
                    </svg>
                    Explore Case Studies
                </a>
            </div>
        </div>

        <!-- Why Join Section -->
        <div class="benefits-section">
            <Group title="Why Developers Love Shit Happens">
                <div class="benefits-grid">
                    <div class="benefit-card">
                        <div class="benefit-icon">
                            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="var(--color-success)" stroke-width="2">
                                <circle cx="12" cy="12" r="10"/>
                                <path d="8 14s1.5 2 4 2 4-2 4-2"/>
                                <line x1="9" y1="9" x2="9.01" y2="9"/>
                                <line x1="15" y1="9" x2="15.01" y2="9"/>
                            </svg>
                        </div>
                        <h3>Transform Failures into Wisdom</h3>
                        <p>Convert your most challenging debugging experiences into valuable learning moments that inspire and educate fellow developers.</p>
                    </div>
                    <div class="benefit-card">
                        <div class="benefit-icon">
                            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="var(--font-color-accent)" stroke-width="2">
                                <path d="M12 2L2 7v10c0 5.55 3.84 10 10 10s10-4.45 10-10V7L12 2z"/>
                                <path d="8 11l2 2 4-4"/>
                            </svg>
                        </div>
                        <h3>Master Through Experience</h3>
                        <p>Discover proven solutions, advanced debugging methodologies, and architectural patterns through real-world case studies.</p>
                    </div>
                    <div class="benefit-card">
                        <div class="benefit-icon">
                            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="var(--color-warning)" stroke-width="2">
                                <path d="17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                                <circle cx="9" cy="7" r="4"/>
                                <path d="23 21v-2a4 4 0 0 0-3-3.87"/>
                                <path d="16 3.13a4 4 0 0 1 0 7.75"/>
                            </svg>
                        </div>
                        <h3>Build Professional Networks</h3>
                        <p>Connect with experienced developers who understand the complexities of software engineering and share similar challenges.</p>
                    </div>
                    <div class="benefit-card">
                        <div class="benefit-icon">
                            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="var(--confession-accent)" stroke-width="2">
                                <path d="M12 2L2 7v10c0 5.55 3.84 10 10 10s10-4.45 10-10V7L12 2z"/>
                                <path d="12 8v5"/>
                                <path d="12 16h.01"/>
                            </svg>
                        </div>
                        <h3>Advance Your Career</h3>
                        <p>Develop resilience, deepen technical expertise, and cultivate the problem-solving mindset that defines exceptional software engineers.</p>
                    </div>
                </div>
            </Group>
        </div>

        <!-- Quick Stats -->
        <div class="stats-section">
            <div class="stat-item">
                <div class="stat-number">1000+</div>
                <div class="stat-label">Bug Confessions</div>
            </div>
            <div class="stat-item">
                <div class="stat-number">500+</div>
                <div class="stat-label">Developers Helped</div>
            </div>
            <div class="stat-item">
                <div class="stat-number">∞</div>
                <div class="stat-label">Lessons Learned</div>
            </div>
        </div>

        <!-- Popular Tags -->
        <Group title="Popular Topics">
            <TagCloud limit={15} />
        </Group>

        <!-- Final CTA -->
        <Admonition>
            The journey from novice to expert is paved with challenges that shape exceptional developers. 
            <strong>Every debugging session, every failed deployment, every breakthrough moment contributes to your professional mastery.</strong>
            <br><br>
            <a href="/submit" style="color: var(--color-success); font-weight: bold;">Share your professional journey →</a>
        </Admonition>

        <!-- Faked for Showing Off Section -->
        <div class="showoff-section">
            <div class="showoff-header">
                <h2>🎭 Faked for Showing Off</h2>
                <p class="showoff-subtitle">Because sometimes we all need to look like coding wizards</p>
            </div>
            
            <div class="showoff-grid">
                <div class="showoff-card">
                    <div class="showoff-icon">
                        <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="var(--color-warning)" stroke-width="2">
                            <path d="M9 12l2 2 4-4"/>
                            <path d="M21 12c-1 0-3-1-3-3s2-3 3-3 3 1 3 3-2 3-3 3"/>
                            <path d="M3 12c1 0 3-1 3-3s-2-3-3-3-3 1-3 3 2 3 3 3"/>
                            <path d="M3 12h6m6 0h6"/>
                        </svg>
                    </div>
                    <h3>Instant Bug Fixes</h3>
                    <p>Pretend you solved that production bug in 5 minutes instead of 5 hours. We won't tell anyone about the 17 Stack Overflow tabs.</p>
                    <div class="fake-badge">✨ 99.9% Success Rate*</div>
                </div>

                <div class="showoff-card">
                    <div class="showoff-icon">
                        <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="var(--color-success)" stroke-width="2">
                            <polyline points="22,12 18,12 15,21 9,3 6,12 2,12"/>
                        </svg>
                    </div>
                    <h3>Heroic Debugging</h3>
                    <p>Transform "I copy-pasted from Stack Overflow" into "I architected a robust solution using industry best practices."</p>
                    <div class="fake-badge">⚡ Lightning Fast*</div>
                </div>

                <div class="showoff-card">
                    <div class="showoff-icon">
                        <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="var(--confession-accent)" stroke-width="2">
                            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                        </svg>
                    </div>
                    <h3>Legendary Status</h3>
                    <p>Become the developer who "never gets bugs in production" and "writes perfect code on the first try."</p>
                    <div class="fake-badge">🏆 Myth Status*</div>
                </div>
            </div>

            <div class="showoff-disclaimer">
                <p><small>* Results may vary. Side effects may include imposter syndrome, increased coffee consumption, and the urge to tell everyone about that one time you fixed a critical bug. Not responsible for inflated egos or disappointed managers. Please debug responsibly.</small></p>
            </div>

            <div class="showoff-cta">
                <button class="fake-button" on:click={() => alert('Just kidding! But seriously, we all fake it till we make it sometimes 😄')}>
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M21 15v4a2 2 0 0-1-2 2H5a2 2 0 0-1-2-2v-4"/>
                        <polyline points="7,10 12,15 17,10"/>
                        <line x1="12" y1="15" x2="12" y2="3"/>
                    </svg>
                    Download Fake Confidence
                </button>
            </div>
        </div>
    </section>
</Page>

<style>
    section {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        align-items: center;
        gap: 40px;
        padding: 20px 0;
        max-width: 1000px;
        margin: 0 auto;
    }

    .hero {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 30px;
        text-align: center;
        padding: 40px 20px;
    }

    .user {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        gap: 15px;
    }

    .user-avatar img {
        height: 80%;
        width: 80%;
        filter: drop-shadow(0 4px 8px rgba(0,0,0,0.3));
    }

    .user-avatar {
        display: inline-flex;
        justify-content: center;
        align-items: center;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        height: 120px;
        width: 120px;
        border-radius: 50%;
        box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);
        transition: transform 0.3s ease;
    }

    .user-avatar:hover {
        transform: scale(1.05);
    }

    .user-label {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: center;
        gap: 8px;
    }

    .user-name {
        position: relative;
        display: flex;
        justify-content: center;
        font-size: 2.5rem;
        color: var(--font-color);
        font-weight: 700;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        background-clip: text;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
    }

    .user-subtext {
        display: flex;
        justify-content: center;
        font-size: 1.2rem;
        color: var(--font-color-accent);
        font-style: italic;
    }

    /* Quote Section */
    .quote-section {
        width: 100%;
        max-width: 700px;
        padding: 30px;
        background: rgba(44, 39, 51, 0.4);
        border-radius: 20px;
        border: 1px solid rgba(255, 255, 255, 0.1);
        backdrop-filter: blur(10px);
    }

    .quote-container {
        min-height: 120px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: opacity 0.3s ease, transform 0.3s ease;
    }

    .quote-container.animating {
        opacity: 0;
        transform: translateY(10px);
    }

    .quote-content {
        text-align: center;
    }

    .quote-text {
        font-size: 1.3rem;
        font-style: italic;
        margin: 0 0 15px 0;
        line-height: 1.6;
        position: relative;
    }

    .quote-text.funny {
        color: var(--color-warning);
    }

    .quote-text.inspiring {
        color: var(--color-success);
    }

    .quote-author {
        display: block;
        font-size: 1rem;
        color: var(--font-color-accent);
        font-weight: 500;
    }

    .quote-indicators {
        display: flex;
        justify-content: center;
        gap: 8px;
        margin-top: 20px;
    }

    .indicator {
        width: 10px;
        height: 10px;
        border-radius: 50%;
        border: 2px solid rgba(118, 75, 162, 0.6);
        background: rgba(44, 39, 51, 0.4);
        cursor: pointer;
        transition: all 0.3s ease;
    }

    .indicator.active {
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.8) 0%, rgba(118, 75, 162, 0.8) 100%);
        border-color: rgba(102, 126, 234, 0.8);
        transform: scale(1.2);
        box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
    }

    .indicator:hover {
        border-color: rgba(118, 75, 162, 0.8);
        background: rgba(118, 75, 162, 0.2);
    }

    /* CTA Section */
    .cta-section {
        text-align: center;
        padding: 40px 20px;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
        border-radius: 20px;
        border: 1px solid rgba(102, 126, 234, 0.2);
        width: 100%;
        max-width: 800px;
    }

    .cta-title {
        font-size: 2rem;
        font-weight: 700;
        color: var(--font-color);
        margin: 0 0 15px 0;
    }

    .cta-subtitle {
        font-size: 1.1rem;
        color: var(--font-color-accent);
        margin: 0 0 30px 0;
        line-height: 1.5;
    }

    .quick-search-container {
        margin: 0 0 30px 0;
        max-width: 500px;
        width: 100%;
    }

    .cta-buttons {
        display: flex;
        gap: 20px;
        justify-content: center;
        flex-wrap: wrap;
    }

    .primary-cta, .secondary-cta {
        display: inline-flex;
        align-items: center;
        gap: 10px;
        padding: 15px 30px;
        border-radius: 12px;
        text-decoration: none;
        font-weight: 600;
        font-size: 1.1rem;
        transition: all 0.3s ease;
        border: 2px solid;
    }

    .primary-cta {
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.8) 0%, rgba(118, 75, 162, 0.8) 100%);
        color: white;
        border-color: rgba(102, 126, 234, 0.6);
        box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
    }

    .primary-cta:hover {
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.9) 0%, rgba(118, 75, 162, 0.9) 100%);
        border-color: rgba(102, 126, 234, 0.8);
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
    }

    .secondary-cta {
        background: rgba(44, 39, 51, 0.6);
        color: var(--font-color);
        border-color: rgba(118, 75, 162, 0.5);
        backdrop-filter: blur(10px);
    }

    .secondary-cta:hover {
        background: rgba(44, 39, 51, 0.8);
        border-color: rgba(118, 75, 162, 0.8);
        transform: translateY(-2px);
        box-shadow: 0 4px 15px rgba(118, 75, 162, 0.2);
    }

    /* Benefits Section */
    .benefits-section {
        width: 100%;
    }

    .benefits-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 25px;
        padding: 20px 0;
    }

    .benefit-card {
        background: rgba(44, 39, 51, 0.3);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 15px;
        padding: 25px;
        text-align: center;
        transition: transform 0.3s ease, box-shadow 0.3s ease;
    }

    .benefit-card:hover {
        transform: translateY(-5px);
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
        border-color: rgba(255, 255, 255, 0.2);
    }

    .benefit-icon {
        margin: 0 auto 20px auto;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 60px;
        height: 60px;
        background: rgba(255, 255, 255, 0.05);
        border-radius: 50%;
        transition: all 0.3s ease;
    }

    .benefit-card:hover .benefit-icon {
        background: rgba(255, 255, 255, 0.1);
        transform: scale(1.1);
    }

    .benefit-card h3 {
        font-size: 1.3rem;
        font-weight: 600;
        color: var(--font-color);
        margin: 0 0 10px 0;
    }

    .benefit-card p {
        color: var(--font-color-accent);
        line-height: 1.5;
        margin: 0;
    }

    /* Stats Section */
    .stats-section {
        display: flex;
        justify-content: center;
        gap: 60px;
        padding: 30px;
        background: rgba(44, 39, 51, 0.2);
        border-radius: 15px;
        border: 1px solid rgba(255, 255, 255, 0.1);
        flex-wrap: wrap;
    }

    .stat-item {
        text-align: center;
    }

    .stat-number {
        font-size: 3rem;
        font-weight: 700;
        color: var(--color-success);
        margin-bottom: 5px;
    }

    .stat-label {
        font-size: 1rem;
        color: var(--font-color-accent);
        font-weight: 500;
    }

    /* Showoff Section */
    .showoff-section {
        width: 100%;
        padding: 40px 20px;
        background: linear-gradient(135deg, rgba(255, 193, 7, 0.05) 0%, rgba(255, 152, 0, 0.05) 100%);
        border-radius: 20px;
        border: 1px solid rgba(255, 193, 7, 0.2);
        text-align: center;
        position: relative;
        overflow: hidden;
    }

    .showoff-section::before {
        content: '';
        position: absolute;
        top: -50%;
        left: -50%;
        width: 200%;
        height: 200%;
        background: linear-gradient(45deg, transparent, rgba(255, 193, 7, 0.02), transparent);
        animation: shimmer 3s ease-in-out infinite;
        pointer-events: none;
    }

    @keyframes shimmer {
        0% { transform: translateX(-100%) translateY(-100%) rotate(45deg); }
        50% { transform: translateX(0%) translateY(0%) rotate(45deg); }
        100% { transform: translateX(100%) translateY(100%) rotate(45deg); }
    }

    .showoff-header h2 {
        font-size: 2.2rem;
        font-weight: 700;
        color: var(--color-warning);
        margin: 0 0 10px 0;
        text-shadow: 0 2px 4px rgba(255, 193, 7, 0.3);
    }

    .showoff-subtitle {
        font-size: 1.1rem;
        color: var(--font-color-accent);
        margin: 0 0 30px 0;
        font-style: italic;
    }

    .showoff-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
        gap: 25px;
        margin: 30px 0;
    }

    .showoff-card {
        background: rgba(44, 39, 51, 0.4);
        border: 1px solid rgba(255, 193, 7, 0.2);
        border-radius: 15px;
        padding: 30px 20px;
        text-align: center;
        position: relative;
        transition: all 0.3s ease;
        overflow: hidden;
    }

    .showoff-card::before {
        content: '';
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: linear-gradient(90deg, transparent, rgba(255, 193, 7, 0.1), transparent);
        transition: left 0.5s ease;
    }

    .showoff-card:hover {
        transform: translateY(-5px);
        border-color: rgba(255, 193, 7, 0.4);
        box-shadow: 0 10px 30px rgba(255, 193, 7, 0.2);
    }

    .showoff-card:hover::before {
        left: 100%;
    }

    .showoff-icon {
        margin: 0 auto 20px auto;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 70px;
        height: 70px;
        background: rgba(255, 193, 7, 0.1);
        border-radius: 50%;
        border: 2px solid rgba(255, 193, 7, 0.3);
        transition: all 0.3s ease;
    }

    .showoff-card:hover .showoff-icon {
        background: rgba(255, 193, 7, 0.2);
        border-color: rgba(255, 193, 7, 0.5);
        transform: scale(1.1) rotate(5deg);
    }

    .showoff-card h3 {
        font-size: 1.4rem;
        font-weight: 600;
        color: var(--font-color);
        margin: 0 0 15px 0;
    }

    .showoff-card p {
        color: var(--font-color-accent);
        line-height: 1.6;
        margin: 0 0 20px 0;
    }

    .fake-badge {
        display: inline-block;
        background: linear-gradient(135deg, rgba(255, 193, 7, 0.8) 0%, rgba(255, 152, 0, 0.8) 100%);
        color: #1a1a1a;
        padding: 8px 15px;
        border-radius: 20px;
        font-size: 0.9rem;
        font-weight: 600;
        box-shadow: 0 2px 8px rgba(255, 193, 7, 0.3);
        animation: pulse 2s ease-in-out infinite;
    }

    @keyframes pulse {
        0%, 100% { transform: scale(1); }
        50% { transform: scale(1.05); }
    }

    .showoff-disclaimer {
        margin: 30px 0 20px 0;
        padding: 20px;
        background: rgba(44, 39, 51, 0.3);
        border-radius: 10px;
        border: 1px solid rgba(255, 255, 255, 0.1);
    }

    .showoff-disclaimer small {
        color: var(--font-color-accent);
        line-height: 1.4;
        font-style: italic;
    }

    .showoff-cta {
        margin-top: 30px;
    }

    .fake-button {
        display: inline-flex;
        align-items: center;
        gap: 10px;
        padding: 15px 30px;
        background: linear-gradient(135deg, rgba(255, 193, 7, 0.8) 0%, rgba(255, 152, 0, 0.8) 100%);
        color: #1a1a1a;
        border: none;
        border-radius: 12px;
        font-size: 1.1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.3s ease;
        box-shadow: 0 4px 15px rgba(255, 193, 7, 0.3);
    }

    .fake-button:hover {
        background: linear-gradient(135deg, rgba(255, 193, 7, 0.9) 0%, rgba(255, 152, 0, 0.9) 100%);
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(255, 193, 7, 0.4);
    }

    .fake-button:active {
        transform: translateY(0);
    }

    /* Responsive Design */
    @media (max-width: 768px) {
        .hero {
            padding: 20px 10px;
        }

        .user-name {
            font-size: 2rem;
        }

        .quote-text {
            font-size: 1.1rem;
        }

        .cta-title {
            font-size: 1.5rem;
        }

        .cta-buttons {
            flex-direction: column;
            align-items: center;
        }

        .primary-cta, .secondary-cta {
            width: 100%;
            max-width: 280px;
            justify-content: center;
        }

        .stats-section {
            gap: 30px;
        }

        .stat-number {
            font-size: 2.5rem;
        }

        .benefits-grid {
            grid-template-columns: 1fr;
        }

        .showoff-grid {
            grid-template-columns: 1fr;
        }

        .showoff-header h2 {
            font-size: 1.8rem;
        }

        .showoff-section {
            padding: 30px 15px;
        }

        .fake-button {
            width: 100%;
            max-width: 280px;
            justify-content: center;
        }
    }
</style>
