<script lang="ts">
    import Gap from "$lib/components/Gap.svelte";
    import Tab from "$lib/components/Tab.svelte";
    import User from "$lib/components/User.svelte";
    import "../app.css";

    import application from "$lib/images/tabs/application.webp";
    import clipboard from "$lib/images/tabs/upload.webp";
    import window from "$lib/images/tabs/categories.webp";

    import colors from "$lib/images/tabs/high-quality.webp";
    import fonts from "$lib/images/tabs/recent.webp";

    import keybinds from "$lib/images/tabs/trending.webp";
    import mouse from "$lib/images/tabs/mouse.webp";
    import sync from "$lib/images/tabs/sync.webp";
    import ghost from "$lib/images/ghost.png";



    import github from "$lib/images/tabs/github.svg";
    import ghostty from "$lib/images/tabs/podcast.webp";

    import config from "$lib/stores/config.svelte";
    import app from "$lib/stores/state.svelte";
    import MobileCardLayout from "$lib/components/MobileCardLayout.svelte";

    const cssConfigVars = $derived.by(() => {
        let str = "";

        const add = (key: string, val: string) => str += `--config-${key}: ${val};`;

        // Add the base colors
        add("bg", config.background);
        add("fg", config.foreground);
        add("selection-bg", config.selectionInvertFgBg ? config.foreground : config.selectionBackground || config.foreground);
        add("selection-fg", config.selectionInvertFgBg ? config.background : config.selectionForeground || config.background);

        // Add the palette colors
        const paletteSize = 16; // config.palette.length;
        for (let c = 0; c < paletteSize; c++) add(`palette-${c}`, config.palette[c]);

        // TODO: consider honoring separate fonts for bold/italic and such in previews
        // Add font settings
        add("font-family", config.fontFamily || "monospace");
        add("font-size", `${config.fontSize}px`);

        return str;
    });

    const {children} = $props();

    // Window management state
    let isDragging = $state(false);
    let isResizing = $state(false);
    let resizeDirection = $state('');
    let dragStart = $state({ x: 0, y: 0 });
    let windowStart = $state({ x: 0, y: 0, width: 0, height: 0 });
    let windowElement: HTMLDivElement;
    let isMinimized = $state(false);
    let isMaximized = $state(false);
    let previousWindowState = $state({ x: 0, y: 0, width: 0, height: 0 });
    let sidebarVisible = $state(true);

    // Window dragging functionality
    function startDrag(e: MouseEvent) {
        if (isMaximized) return;
        isDragging = true;
        const rect = windowElement.getBoundingClientRect();
        dragStart = { x: e.clientX, y: e.clientY };
        windowStart = { x: rect.left, y: rect.top, width: rect.width, height: rect.height };
        document.addEventListener('mousemove', handleDrag);
        document.addEventListener('mouseup', stopDrag);
        e.preventDefault();
    }

    function handleDrag(e: MouseEvent) {
        if (!isDragging) return;
        const deltaX = e.clientX - dragStart.x;
        const deltaY = e.clientY - dragStart.y;
        const newX = Math.max(0, Math.min(window.innerWidth - windowStart.width, windowStart.x + deltaX));
        const newY = Math.max(0, Math.min(window.innerHeight - windowStart.height, windowStart.y + deltaY));
        
        windowElement.style.left = `${newX}px`;
        windowElement.style.top = `${newY}px`;
    }

    function stopDrag() {
        isDragging = false;
        document.removeEventListener('mousemove', handleDrag);
        document.removeEventListener('mouseup', stopDrag);
    }

    // Window resizing functionality
    function startResize(e: MouseEvent, direction: string) {
        if (isMaximized) return;
        isResizing = true;
        resizeDirection = direction;
        const rect = windowElement.getBoundingClientRect();
        dragStart = { x: e.clientX, y: e.clientY };
        windowStart = { x: rect.left, y: rect.top, width: rect.width, height: rect.height };
        document.addEventListener('mousemove', handleResize);
        document.addEventListener('mouseup', stopResize);
        e.preventDefault();
        e.stopPropagation();
    }

    function handleResize(e: MouseEvent) {
        if (!isResizing) return;
        const deltaX = e.clientX - dragStart.x;
        const deltaY = e.clientY - dragStart.y;
        
        let newX = windowStart.x;
        let newY = windowStart.y;
        let newWidth = windowStart.width;
        let newHeight = windowStart.height;

        const minWidth = 400;
        const minHeight = 300;
        const maxWidth = 1400;
        const maxHeight = 900;

        if (resizeDirection.includes('right')) {
            newWidth = Math.max(minWidth, Math.min(maxWidth, windowStart.width + deltaX));
        }
        if (resizeDirection.includes('left')) {
            const proposedWidth = windowStart.width - deltaX;
            if (proposedWidth >= minWidth && proposedWidth <= maxWidth) {
                newWidth = proposedWidth;
                newX = windowStart.x + deltaX;
            }
        }
        if (resizeDirection.includes('bottom')) {
            newHeight = Math.max(minHeight, Math.min(maxHeight, windowStart.height + deltaY));
        }
        if (resizeDirection.includes('top')) {
            const proposedHeight = windowStart.height - deltaY;
            if (proposedHeight >= minHeight && proposedHeight <= maxHeight) {
                newHeight = proposedHeight;
                newY = windowStart.y + deltaY;
            }
        }

        // Ensure window stays within viewport
        newX = Math.max(0, Math.min(window.innerWidth - newWidth, newX));
        newY = Math.max(0, Math.min(window.innerHeight - newHeight, newY));

        windowElement.style.left = `${newX}px`;
        windowElement.style.top = `${newY}px`;
        windowElement.style.width = `${newWidth}px`;
        windowElement.style.height = `${newHeight}px`;
    }

    function stopResize() {
        isResizing = false;
        resizeDirection = '';
        document.removeEventListener('mousemove', handleResize);
        document.removeEventListener('mouseup', stopResize);
    }

    // Window control functions
    function minimizeWindow() {
        if (isMinimized) {
            // Restore from taskbar
            windowElement.style.display = 'flex';
            windowElement.style.transform = 'scale(1)';
            windowElement.style.opacity = '1';
            isMinimized = false;
        } else {
            // Minimize to taskbar
            windowElement.style.transform = 'scale(0.1)';
            windowElement.style.opacity = '0';
            setTimeout(() => {
                windowElement.style.display = 'none';
                isMinimized = true;
            }, 300);
        }
    }

    function maximizeWindow() {
        if (isMaximized) {
            // Restore
            windowElement.style.left = `${previousWindowState.x}px`;
            windowElement.style.top = `${previousWindowState.y}px`;
            windowElement.style.width = `${previousWindowState.width}px`;
            windowElement.style.height = `${previousWindowState.height}px`;
            isMaximized = false;
        } else {
            // Maximize
            const rect = windowElement.getBoundingClientRect();
            previousWindowState = { x: rect.left, y: rect.top, width: rect.width, height: rect.height };
            windowElement.style.left = '0px';
            windowElement.style.top = '0px';
            windowElement.style.width = '100vw';
            windowElement.style.height = '100vh';
            isMaximized = true;
        }
    }

    function closeWindow() {
        // Close with animation to taskbar
        windowElement.style.transform = 'scale(0.05)';
        windowElement.style.opacity = '0';
        setTimeout(() => {
            windowElement.style.display = 'none';
            isMinimized = true;
        }, 400);
    }

    function restoreFromTaskbar() {
        // Restore window from taskbar icon
        if (isMinimized) {
            windowElement.style.display = 'flex';
            windowElement.style.transform = 'scale(0.1)';
            windowElement.style.opacity = '0';
            setTimeout(() => {
                windowElement.style.transform = 'scale(1)';
                windowElement.style.opacity = '1';
                isMinimized = false;
                centerWindow();
            }, 50);
        }
    }

    function centerWindow() {
        if (!windowElement) return;
        const rect = windowElement.getBoundingClientRect();
        const centerX = (window.innerWidth - rect.width) / 2;
        const centerY = (window.innerHeight - rect.height) / 2;
        windowElement.style.left = `${Math.max(0, centerX)}px`;
        windowElement.style.top = `${Math.max(0, centerY)}px`;
    }

    function toggleSidebar() {
        sidebarVisible = !sidebarVisible;
    }

    function handleSidebarToggle(e: MouseEvent) {
        e.stopPropagation();
        toggleSidebar();
    }

    // Initialize window position on mount
    $effect(() => {
        if (windowElement) {
            centerWindow();
        }
        
        // Add keyboard shortcut for sidebar toggle
        const handleKeydown = (e: KeyboardEvent) => {
            if ((e.ctrlKey || e.metaKey) && e.key === 'b') {
                e.preventDefault();
                toggleSidebar();
            }
        };
        
        document.addEventListener('keydown', handleKeydown);
        
        return () => {
            document.removeEventListener('keydown', handleKeydown);
        };
    });




    const htmlTitle = $derived.by(() => {
        const name = app.title === "Shit Happens" ? "" : app.title;
        let title = "Shit Happens";
        if (name) title = `${title} - ${name}`;
        return title;
    });
</script>

<svelte:head>
    <title>{htmlTitle}</title>
</svelte:head>

<!-- eslint-disable-next-line svelte/require-optimized-style-attribute -->
<div class="app-window" style={cssConfigVars} bind:this={windowElement}>
    <!-- Resize handles -->
    <div class="resize-handle resize-top" onmousedown={(e) => startResize(e, 'top')}></div>
    <div class="resize-handle resize-right" onmousedown={(e) => startResize(e, 'right')}></div>
    <div class="resize-handle resize-bottom" onmousedown={(e) => startResize(e, 'bottom')}></div>
    <div class="resize-handle resize-left" onmousedown={(e) => startResize(e, 'left')}></div>
    <div class="resize-handle resize-top-left" onmousedown={(e) => startResize(e, 'top-left')}></div>
    <div class="resize-handle resize-top-right" onmousedown={(e) => startResize(e, 'top-right')}></div>
    <div class="resize-handle resize-bottom-left" onmousedown={(e) => startResize(e, 'bottom-left')}></div>
    <div class="resize-handle resize-bottom-right" onmousedown={(e) => startResize(e, 'bottom-right')}></div>
    
    <div id="sidebar" class:hidden={!sidebarVisible}>
        <div class="sidebar-header">
            <div class="window-actions-container">
                <div class="window-actions">
                    <div class="window-dot close-btn" onclick={closeWindow}><span>×</span></div>
                    <div class="window-dot minimize-btn" onclick={minimizeWindow}><span>-</span></div>
                    <div class="window-dot maximize-btn" onclick={maximizeWindow}><span>{isMaximized ? '⧉' : '+'}</span></div>
                </div>
            </div>
            <!-- Mobile sidebar close button -->
            <button class="mobile-sidebar-close" onclick={handleSidebarToggle} aria-label="Close sidebar" style="display: none;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
            </button>
            <!-- Draggable title bar -->
            <div class="title-bar" onmousedown={startDrag}>
                <button class="sidebar-toggle" onclick={handleSidebarToggle} aria-label="Toggle sidebar">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="3" y1="6" x2="21" y2="6"></line>
                        <line x1="3" y1="12" x2="21" y2="12"></line>
                        <line x1="3" y1="18" x2="21" y2="18"></line>
                    </svg>
                </button>
                <span class="window-title">{app.title}</span>
                <div class="spacer"></div>
            </div>
        </div>
        
        {#if sidebarVisible}
        <nav id="categories">
            <User route="/" />
            <Gap />
            <Tab route="/submit">
                {#snippet icon()}<img src={clipboard} alt="Submit Story" />{/snippet}
                Submit Your Story
            </Tab>
            <Tab route="/categories">
                {#snippet icon()}<img src={window} alt="Categories" />{/snippet}
                Categories
            </Tab>
            <Tab route="/search">
                {#snippet icon()}<img src={sync} alt="Search" />{/snippet}
                Search
            </Tab>
            <Gap />
            <Tab route="/confessions?filter=top">
                {#snippet icon()}<img src={colors} alt="Top Confessions" />{/snippet}
                Top Confessions
            </Tab>
            <Tab route="/confessions/trending/weekly">
                {#snippet icon()}<img src={keybinds} alt="Trending This Week" />{/snippet}
                Trending This Week
            </Tab>
            <Tab route="/confessions">
                {#snippet icon()}<img src={fonts} alt="Recent Posts" />{/snippet}
                Recent Posts
            </Tab>
            <Gap expand={true} />
            <Tab route="https://github.com/Balaji01-4D/shit-happens">
                {#snippet icon()}<div class="icon-wrapper github"><img src={github} alt="GitHub Repository" /></div>{/snippet}
                GitHub
            </Tab>
            <Tab route="/about">
                {#snippet icon()}<img src={ghostty} alt="About Shit Happens" />{/snippet}
                About
            </Tab>
        </nav>
        {/if}
    </div>
    <div id="content-view" class:full-width={!sidebarVisible}>
        {#if !sidebarVisible}
        <div class="floating-sidebar-toggle">
            <button onclick={handleSidebarToggle} aria-label="Show sidebar">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="3" y1="6" x2="21" y2="6"></line>
                    <line x1="3" y1="12" x2="21" y2="12"></line>
                    <line x1="3" y1="18" x2="21" y2="18"></line>
                </svg>
            </button>
        </div>
        {/if}
        {@render children()}
    </div>

</div>

<!-- Taskbar icon for minimized window -->
<div class="taskbar-icon" class:visible={isMinimized} onclick={restoreFromTaskbar} role="button" tabindex="0" aria-label="Restore window">
    <div class="taskbar-icon-content">
        <img src={ghost} alt="Shit Happens" class="taskbar-ghost-icon" />
        <div class="taskbar-tooltip">Shit Happens</div>
    </div>
</div>

<!-- <svelte:window onmouseup={onMouseUp} onmousemove={onMouseMove} /> -->

<style>
.app-window {
    user-select: none;
    display: flex;
    position: fixed;
    flex-direction: row;
    min-width: 400px;
    min-height: 300px;
    max-width: 1400px;
    max-height: 900px;
    width: var(--app-width);
    height: var(--app-height);
    border: 1px solid var(--border-level-1);
    box-shadow: 0 0 20px -1px rgba(0,0,0,0.7);
    border-radius: var(--radius-level-1);
    overflow: hidden;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    z-index: 1000;
    transform-origin: bottom right;
}

.app-window:hover {
    box-shadow: 0 0 30px -1px rgba(0,0,0,0.8);
}

/* Taskbar icon */
.taskbar-icon {
    position: fixed;
    bottom: 20px;
    right: 20px;
    width: 60px;
    height: 60px;
    cursor: pointer;
    z-index: 2000;
    opacity: 0;
    transform: scale(0.8);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    pointer-events: none;
}

.taskbar-icon.visible {
    opacity: 1;
    transform: scale(1);
    pointer-events: all;
}

.taskbar-icon:hover {
    transform: scale(1.1);
}

.taskbar-icon:active {
    transform: scale(0.95);
}

.taskbar-icon-content {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(44, 39, 51, 0.9);
    backdrop-filter: blur(20px);
    border-radius: 50%;
    border: 2px solid rgba(255, 255, 255, 0.1);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.taskbar-ghost-icon {
    width: 32px;
    height: 32px;
    filter: drop-shadow(0 2px 8px rgba(0, 0, 0, 0.3));
    transition: all 0.3s ease;
}

.taskbar-icon:hover .taskbar-ghost-icon {
    filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.5)) brightness(1.1);
}

.taskbar-tooltip {
    position: absolute;
    bottom: 70px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(0, 0, 0, 0.8);
    color: white;
    padding: 6px 12px;
    border-radius: 6px;
    font-size: 12px;
    white-space: nowrap;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s ease;
}

.taskbar-icon:hover .taskbar-tooltip {
    opacity: 1;
}

/* Resize handles */
.resize-handle {
    position: absolute;
    z-index: 1001;
}

.resize-top, .resize-bottom {
    left: 8px;
    right: 8px;
    height: 8px;
    cursor: ns-resize;
}

.resize-left, .resize-right {
    top: 8px;
    bottom: 8px;
    width: 8px;
    cursor: ew-resize;
}

.resize-top { top: -4px; }
.resize-bottom { bottom: -4px; }
.resize-left { left: -4px; }
.resize-right { right: -4px; }

.resize-top-left, .resize-top-right, .resize-bottom-left, .resize-bottom-right {
    width: 12px;
    height: 12px;
}

.resize-top-left { top: -4px; left: -4px; cursor: nw-resize; }
.resize-top-right { top: -4px; right: -4px; cursor: ne-resize; }
.resize-bottom-left { bottom: -4px; left: -4px; cursor: sw-resize; }
.resize-bottom-right { bottom: -4px; right: -4px; cursor: se-resize; }

/* Hover effect for resize handles */
.resize-handle:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.title-bar {
    flex: 1;
    height: 30px;
    cursor: move;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 10px;
    user-select: none;
    position: relative;
}

.title-bar:hover {
    background-color: rgba(255, 255, 255, 0.05);
}

.sidebar-toggle {
    position: absolute;
    left: 10px;
    background: none;
    border: none;
    color: var(--font-color);
    cursor: pointer;
    padding: 4px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
    z-index: 10;
}

.sidebar-toggle:hover {
    background-color: rgba(255, 255, 255, 0.1);
    color: var(--font-color-accent);
}

.spacer {
    position: absolute;
    right: 10px;
    width: 24px;
}

.window-title {
    color: var(--font-color);
    font-size: 13px;
    font-weight: 500;
    opacity: 0.8;
}

.app-window {
    /* Rest of existing styles remain the same */
    /* margin: 20px auto; removed for positioning */
}

/* .app-window .draggable {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 50px;
    cursor: move;
    z-index: 10;
} */

/* TODO: try this without pseudoelement using outline */
.app-window::before {
    content: "";
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    box-shadow: 0 0 1px rgba(255, 255, 255, 0.3) inset;
    border-radius: inherit;
    z-index: 2;
    pointer-events: none;
}

#sidebar {
    width: var(--sidebar-width);
    /* black: #272329; white: #544F57; */
    background: rgba(50, 46, 52, 0.7);
    backdrop-filter: blur(10px);
    padding: 5px;
    border-right: 2px solid var(--border-level-1);
    display: flex;
    flex-direction: column;
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar-header {
    display: flex;
    flex-direction: column;
}

.sidebar-header .window-actions-container {
    display: flex;
    padding: 15px 0 0 15px;
    margin-bottom: 20px;
}

.sidebar-header .window-actions {
    display: flex;
    gap: 8px;
}

.app-window .window-dot {
    border-radius: 50%;
    height: 12px;
    width: 12px;
    border: 0;
    display: inline-flex;
    justify-content: center;
    align-items: center;
    color: rgba(0, 0, 0, 0);
    cursor: pointer;
    transition: all 0.2s ease;
}

.window-dot:hover {
    transform: scale(1.1);
}

.window-dot:hover span {
    color: rgba(0, 0, 0, 0.8);
}

.window-dot:active {
    transform: scale(0.95);
}

.window-dot span {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    margin-top: -3px;
    margin-right: -1px;
    font-size: 10px;
    font-weight: bold;
}

.app-window .window-actions:hover .window-dot {
    /* background: white!important; */
    /* cursor: pointer; */
    color: rgba(0, 0, 0, 0.5);
}

.window-dot {
    background-color: var(--color-warning);
}

.window-dot.close-btn {
    background-color: var(--color-danger);
}

.window-dot.minimize-btn {
    background-color: var(--color-warning);
}

.window-dot.maximize-btn {
    background-color: var(--color-success);
}

#content-view {
    background: rgba(44, 39, 51, 0.8);
    backdrop-filter: blur(10px);
    flex: 1;
    display: flex;
    min-width: 0;
    position: relative;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

#content-view::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, 
        rgba(255, 255, 255, 0.02) 0%, 
        rgba(255, 255, 255, 0.01) 50%, 
        rgba(0, 0, 0, 0.02) 100%);
    pointer-events: none;
    z-index: -1;
}

#content-view.full-width {
    width: 100%;
}

.floating-sidebar-toggle {
    position: absolute;
    top: 20px;
    left: 20px;
    z-index: 1000;
}

.floating-sidebar-toggle button {
    background: rgba(44, 39, 51, 0.9);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    color: var(--font-color);
    cursor: pointer;
    padding: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.floating-sidebar-toggle button:hover {
    background: rgba(44, 39, 51, 0.95);
    border-color: rgba(255, 255, 255, 0.2);
    color: var(--font-color-accent);
    transform: scale(1.05);
}

#sidebar {
    width: var(--sidebar-width);
    min-width: var(--sidebar-width);
    max-width: var(--sidebar-width);
    background: rgba(44, 39, 51, 0.9);
    backdrop-filter: blur(15px);
    padding: 5px;
    border-right: 2px solid var(--border-level-1);
    display: flex;
    flex-direction: column;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    flex-shrink: 0;
}

#sidebar.hidden {
    width: 0;
    min-width: 0;
    max-width: 0;
    padding: 0;
    border-right: none;
    overflow: hidden;
}



#categories {
    display: flex;
    flex-direction: column;
    flex: 1;
    padding: 10px;
}

#categories img {
    width: 100%;
}

#categories .icon-wrapper {
    background: linear-gradient(#D3E3E9, #908F8C);
    width: 20px;
    height: 20px;
    border-radius: var(--radius-level-4);
    display: inline-flex;
    justify-content: center;
    align-items: center;
}

#categories .icon-wrapper img {
    height: 14px;
    width: 14px;
}

#categories .icon-wrapper.github {
    background: linear-gradient(#9C45A9, #3B1E68);
}

#categories .icon-wrapper.github img {
    filter: invert(100%);
    height: 18px;
    width: 18px;
}

/* Mobile responsive overrides without touching desktop look */
@media (max-width: 900px) {
    :root { --app-width: 100vw; --app-height: 100dvh; --sidebar-width: 240px; }
    .app-window {
        position: fixed;
        top: 0 !important; left: 0 !important; right: 0; bottom: 0;
        width: 100vw !important; height: 100dvh !important;
        max-width: none; max-height: none; min-width: 0; min-height: 0;
        border-radius: 0; box-shadow: none; border: none;
        flex-direction: column;
    }
    #sidebar { position: fixed; top: 0; left: 0; bottom: 0; height: auto; z-index: 3000; transform: translateX(0); width: 75vw; max-width: 320px; }
    #sidebar.hidden { transform: translateX(-100%); width: 75vw; min-width: 0; max-width: 320px; }
    #content-view { flex: 1; width: 100%; overflow-y: auto; -webkit-overflow-scrolling: touch; }
    .floating-sidebar-toggle { top: 12px; left: 12px; }
    .title-bar { display: none; }
    .sidebar-header { padding-top: 4px; position: relative; }
    .sidebar-header .window-actions-container { display: none; }
    
    /* Mobile close button for sidebar */
    .mobile-sidebar-close {
        display: block !important;
        position: absolute;
        top: 8px;
        right: 8px;
        background: rgba(255,255,255,0.1);
        border: none;
        border-radius: 8px;
        padding: 8px;
        color: var(--font-color);
        cursor: pointer;
        z-index: 3001;
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background-color 0.2s ease;
    }
    
    .mobile-sidebar-close:hover {
        background-color: rgba(255,255,255,0.2);
    }
    .resize-handle, .taskbar-icon { display: none; }
}

@media (max-width: 900px) and (prefers-color-scheme: light) {
    .app-window { background: rgba(250,250,250,0.75); }
}
</style>

<!-- Mobile optimized layout overlay -->
<MobileCardLayout />
