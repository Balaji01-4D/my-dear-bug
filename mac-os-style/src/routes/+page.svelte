<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import logo from "$lib/images/ghost.svg"; // We'll keep the ghost as a bug icon for now
    import sync from "$lib/images/tabs/sync.webp";
    import fonts from "$lib/images/tabs/font-playground.webp";
    import Admonition from "$lib/components/Admonition.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import LinkItem from "$lib/components/settings/LinkItem.svelte";
    import Separator from "$lib/components/settings/Separator.svelte";
</script>

<Page title="Shit Happens">
    <section>
        <div class="user">
            <div class="user-avatar">
                <img src={logo} alt="Bug Icon" />
            </div>
            <div class="user-label">
                <div class="user-name">Shit Happens<sup>β</sup></div>
                <div class="user-subtext">Coding Confessions</div>
            </div>
        </div>
        <Admonition>
            Welcome to <em>Shit Happens</em> - where developers confess their coding sins, share their debugging nightmares, and help others learn from their mistakes. 
            <a
                href="https://github.com/Balaji01-4D/my-dear-bug"
                target="_blank"
                rel="noopener noreferrer">Join us on GitHub</a>!
        </Admonition>
        <Group>
            <LinkItem name="Browse Confessions" href="/confessions" icon={sync} />
            <Separator />
            <LinkItem name="Submit Your Story" href="/submit" icon={fonts} />
        </Group>
        <Group title="What's New v0.1.0">
            <ul>
                <li>
                    🐛 Share your most embarrassing coding mistakes
                </li>
                <li>🔍 Search through thousands of developer confessions</li>
                <li>💡 Learn from others' debugging experiences</li>
                <li>⭐ Vote on the most relatable stories</li>
            </ul>
        </Group>
        <div class="group-wrapper">
            <Group title="Featured Categories" flex={1}>
                <ul>
                    <li>JavaScript Nightmares</li>
                    <li>CSS Horror Stories</li>
                    <li>Database Disasters</li>
                    <li>Production Failures</li>
                    <li>Debugging Adventures</li>
                </ul>
            </Group>
            <Group title="Community Stats">
                <ul>
                    <li>Over 1000+ confessions shared</li>
                    <li>500+ developers helped</li>
                    <li>Countless bugs squashed</li>
                    <li>Endless learning opportunities</li>
                </ul>
            </Group>
        </div>
    </section>
</Page>

<style>
    section {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        flex: 0.6;
    }

    .user {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        gap: 8px;
        margin: 20px 0;
    }

    .user-avatar img {
        height: 80%;
        width: 80%;
    }

    .user-avatar {
        display: inline-flex;
        justify-content: center;
        align-items: center;
        background: linear-gradient(#03028f, #09043a);
        height: 100px;
        width: 100px;
        border-radius: 50%;
        /* border: 2px solid black; */
        /* box-shadow: 0 0 0 2px white; */
    }

    .user-label {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: center;
        gap: 5px;
    }

    .user-name {
        position: relative;
        display: flex;
        justify-content: center;
        font-size: 1.3rem;
        color: var(--font-color);
        font-weight: 600;
    }

    sup {
        position: absolute;
        color: var(--font-color-accent);
        right: -11px;
        top: -5px;
    }

    .user-subtext {
        display: flex;
        justify-content: center;
        font-size: 1.1rem;
    }

    .group-wrapper {
        display: flex;
        width: 100%;
        gap: 12px;
    }

    ul,
    .group-wrapper ul {
        /* list-style: "↪ "; */
        list-style: "🠶 ";
        padding-left: 15px;
        margin: 0;
        display: flex;
        flex-direction: column;
        height: 100%;
        justify-content: space-evenly;
    }

    /* code {
        background: var(--bg-input-focus);
        padding: 2px 4px;
        border: 1px solid var(--border-input);
        border-radius: var(--radius-level-4);
    } */
</style>
