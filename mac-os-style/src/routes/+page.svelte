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
                    <div class="user-name">Shit Happens<sup>β</sup></div>
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
                            onclick={() => {
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
            <h2 class="cta-title">Ready to Share Your Coding Horror Story?</h2>
            <p class="cta-subtitle">Join thousands of developers who've turned their worst coding moments into learning experiences!</p>
            
            <div class="cta-buttons">
                <a href="/submit" class="primary-cta">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M12 19l7-7 3 3-7 7-3-3z"/>
                        <path d="m18 13-1.5-7.5L2 2l3.5 14.5L13 18l5-5z"/>
                        <path d="m2 2 7.586 7.586"/>
                        <circle cx="11" cy="11" r="2"/>
                    </svg>
                    Confess Your Bug Story
                </a>
                <a href="/confessions" class="secondary-cta">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="11" cy="11" r="8"/>
                        <path d="m21 21-4.35-4.35"/>
                    </svg>
                    Browse Confessions
                </a>
            </div>
        </div>

        <!-- Why Join Section -->
        <div class="benefits-section">
            <Group title="Why Developers Love Shit Happens">
                <div class="benefits-grid">
                    <div class="benefit-card">
                        <div class="benefit-icon">😅</div>
                        <h3>Laugh at Your Mistakes</h3>
                        <p>Turn your most embarrassing bugs into hilarious stories that help others feel less alone.</p>
                    </div>
                    <div class="benefit-card">
                        <div class="benefit-icon">🧠</div>
                        <h3>Learn from Chaos</h3>
                        <p>Discover solutions, patterns, and debugging techniques from real-world disasters.</p>
                    </div>
                    <div class="benefit-card">
                        <div class="benefit-icon">🤝</div>
                        <h3>Connect with Peers</h3>
                        <p>Realize you're not the only one who's spent 6 hours debugging a missing semicolon.</p>
                    </div>
                    <div class="benefit-card">
                        <div class="benefit-icon">🚀</div>
                        <h3>Grow Together</h3>
                        <p>Build resilience, share wisdom, and celebrate the beautiful mess that is software development.</p>
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
            Remember: Every expert was once a beginner. Every pro was once an amateur. Every icon was once an unknown. 
            <strong>Your bugs are just stepping stones to greatness!</strong>
            <br><br>
            <a href="/submit" style="color: var(--color-success); font-weight: bold;">Share your story now →</a>
        </Admonition>
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

    sup {
        position: absolute;
        color: var(--color-warning);
        right: -20px;
        top: -10px;
        font-size: 1rem;
        background: rgba(255, 193, 7, 0.2);
        padding: 2px 6px;
        border-radius: 8px;
        border: 1px solid var(--color-warning);
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
        border: 2px solid var(--font-color-accent);
        background: transparent;
        cursor: pointer;
        transition: all 0.3s ease;
    }

    .indicator.active {
        background: var(--color-success);
        border-color: var(--color-success);
        transform: scale(1.2);
    }

    .indicator:hover {
        border-color: var(--color-success);
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
        background: var(--color-success);
        color: white;
        border-color: var(--color-success);
        box-shadow: 0 4px 15px rgba(40, 167, 69, 0.3);
    }

    .primary-cta:hover {
        background: #218838;
        border-color: #218838;
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(40, 167, 69, 0.4);
    }

    .secondary-cta {
        background: transparent;
        color: var(--font-color);
        border-color: var(--font-color-accent);
    }

    .secondary-cta:hover {
        background: rgba(255, 255, 255, 0.1);
        border-color: var(--font-color);
        transform: translateY(-2px);
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
        font-size: 3rem;
        margin-bottom: 15px;
        display: block;
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
    }
</style>
