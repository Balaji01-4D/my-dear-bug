<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import Text from "$lib/components/settings/Text.svelte";
    import Dropdown from "$lib/components/settings/Dropdown.svelte";
    import Gap from "$lib/components/Gap.svelte";
    import api from "$lib/api";
    import type { ConfessionRequest, LoadingState } from "$lib/types";
    import { goto } from "$app/navigation";

    let title = $state("");
    let description = $state("");
    let codeSnippet = $state("");
    let language = $state("JavaScript");
    let tags = $state("");
    let isFlagged = $state(false);
    let submitState: LoadingState = $state('idle');
    let errorMessage = $state('');
    
    // Tag suggestions
    let tagSuggestions: string[] = $state([]);
    let showSuggestions = $state(false);
    let isSearchingTags = $state(false);
    
    // Code editor
    let codeTextarea: HTMLTextAreaElement;
    let highlightedCode: HTMLElement;
    
    // Notification system
    let showNotification = $state(false);
    let notificationMessage = $state('');
    let notificationIcon = $state('🐛');
    let notificationType: 'success' | 'error' | 'info' = $state('success');
    
    // Code placeholder text
    const codePlaceholder = `// Paste your problematic code here...
// Example:
function buggyFunction() {
    var result = 0
    for (var i = 0; i <= array.length; i++) {
        result += array[i]; // Off-by-one error!
    }
    return result;
}`;
    
    const languages = [
        "JavaScript",
        "TypeScript", 
        "Python",
        "Java",
        "C++",
        "C#",
        "PHP",
        "Ruby",
        "Go",
        "Rust",
        "CSS",
        "HTML",
        "SQL",
        "DevOps",
        "Other"
    ];

    // Notification functions
    function showMacNotification(message: string, type: 'success' | 'error' | 'info' = 'success', icon: string = '🐛') {
        notificationMessage = message;
        notificationType = type;
        notificationIcon = icon;
        showNotification = true;

        // Auto-hide after 4 seconds
        setTimeout(() => {
            hideNotification();
        }, 4000);
    }

    function hideNotification() {
        showNotification = false;
        setTimeout(() => {
            notificationMessage = '';
        }, 300); // Wait for animation to complete
    }

    // Code editor functionality
    function handleCodeInput(event: Event) {
        const textarea = event.target as HTMLTextAreaElement;
        highlightCode(textarea.value);
        
        // Enforce character and line limits
        const lines = textarea.value.split('\n');
        if (lines.length > 100) {
            const truncatedLines = lines.slice(0, 100);
            codeSnippet = truncatedLines.join('\n');
            textarea.value = codeSnippet;
        }
    }

    function handleCodeKeydown(event: KeyboardEvent) {
        const textarea = event.target as HTMLTextAreaElement;
        const { selectionStart, selectionEnd, value } = textarea;

        // Auto-closing brackets and quotes
        const pairs: Record<string, string> = {
            '(': ')',
            '[': ']',
            '{': '}',
            '"': '"',
            "'": "'",
            '`': '`'
        };

        if (pairs[event.key]) {
            event.preventDefault();
            const before = value.substring(0, selectionStart);
            const after = value.substring(selectionEnd);
            const newValue = before + event.key + pairs[event.key] + after;
            
            codeSnippet = newValue;
            textarea.value = newValue;
            
            // Set cursor position between the pair
            setTimeout(() => {
                textarea.selectionStart = textarea.selectionEnd = selectionStart + 1;
            }, 0);
            
            highlightCode(newValue);
            return;
        }

        // Auto-indent on Enter
        if (event.key === 'Enter') {
            event.preventDefault();
            const before = value.substring(0, selectionStart);
            const after = value.substring(selectionEnd);
            const lines = before.split('\n');
            const currentLine = lines[lines.length - 1];
            
            // Calculate indentation
            const indentMatch = currentLine.match(/^(\s*)/);
            let indent = indentMatch ? indentMatch[1] : '';
            
            // Add extra indent if line ends with opening bracket
            if (/[{(\[]$/.test(currentLine.trim())) {
                indent += '    '; // 4 spaces
            }
            
            const newValue = before + '\n' + indent + after;
            codeSnippet = newValue;
            textarea.value = newValue;
            
            setTimeout(() => {
                textarea.selectionStart = textarea.selectionEnd = selectionStart + 1 + indent.length;
            }, 0);
            
            highlightCode(newValue);
            return;
        }

        // Tab handling
        if (event.key === 'Tab') {
            event.preventDefault();
            const before = value.substring(0, selectionStart);
            const after = value.substring(selectionEnd);
            const newValue = before + '    ' + after; // 4 spaces
            
            codeSnippet = newValue;
            textarea.value = newValue;
            
            setTimeout(() => {
                textarea.selectionStart = textarea.selectionEnd = selectionStart + 4;
            }, 0);
            
            highlightCode(newValue);
            return;
        }
    }

    function syncScroll(event: Event) {
        const textarea = event.target as HTMLTextAreaElement;
        if (highlightedCode) {
            highlightedCode.scrollTop = textarea.scrollTop;
            highlightedCode.scrollLeft = textarea.scrollLeft;
        }
    }

    function highlightCode(code: string) {
        if (!highlightedCode) return;
        
        // Simple syntax highlighting for common patterns
        let highlighted = code
            // Escape HTML
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            
            // Comments
            .replace(/\/\/.*$/gm, '<span class="comment">$&</span>')
            .replace(/\/\*[\s\S]*?\*\//g, '<span class="comment">$&</span>')
            .replace(/#.*$/gm, '<span class="comment">$&</span>')
            
            // Strings
            .replace(/"([^"\\]|\\.)*"/g, '<span class="string">$&</span>')
            .replace(/'([^'\\]|\\.)*'/g, '<span class="string">$&</span>')
            .replace(/`([^`\\]|\\.)*`/g, '<span class="string">$&</span>')
            
            // Keywords (common across languages)
            .replace(/\b(function|const|let|var|if|else|for|while|return|class|import|export|from|try|catch|finally|throw|new|this|null|undefined|true|false|async|await|def|print|len|range|int|str|float|bool|list|dict|tuple|set|public|private|protected|static|void|string|int|double|char|boolean|array)\b/g, '<span class="keyword">$&</span>')
            
            // Numbers
            .replace(/\b\d+\.?\d*\b/g, '<span class="number">$&</span>')
            
            // Operators
            .replace(/[+\-*/%=<>!&|^~]/g, '<span class="operator">$&</span>');

        highlightedCode.innerHTML = highlighted + '\n'; // Add newline to prevent layout shift
    }

    // Initialize highlighting when component mounts
    $effect(() => {
        if (codeSnippet) {
            highlightCode(codeSnippet);
        }
    });

    async function handleSubmit() {
        // Validate required fields according to Go backend validation
        const trimmedTitle = title.trim();
        const trimmedDescription = description.trim();
        const trimmedCodeSnippet = codeSnippet.trim();
        
        if (!trimmedTitle || trimmedTitle.length < 5 || trimmedTitle.length > 100) {
            errorMessage = "Title must be between 5 and 100 characters";
            submitState = 'error';
            return;
        }
        
        if (!trimmedDescription || trimmedDescription.length < 10) {
            errorMessage = "Description must be at least 10 characters";
            submitState = 'error';
            return;
        }
        
        if (!language.trim()) {
            errorMessage = "Programming language is required";
            submitState = 'error';
            return;
        }

        // Validate code snippet constraints
        if (trimmedCodeSnippet) {
            if (trimmedCodeSnippet.length > 2000) {
                errorMessage = "Code snippet must be less than 2000 characters";
                submitState = 'error';
                return;
            }
            
            const lines = trimmedCodeSnippet.split('\n');
            if (lines.length > 100) {
                errorMessage = "Code snippet must have fewer than 100 lines";
                submitState = 'error';
                return;
            }

            // Basic security check - prevent potentially malicious content
            const suspiciousPatterns = [
                /DROP\s+TABLE/i,
                /DELETE\s+FROM/i,
                /INSERT\s+INTO/i,
                /UPDATE\s+.*SET/i,
                /EXEC\s*\(/i,
                /eval\s*\(/i,
                /document\.write/i,
                /innerHTML\s*=/i,
                /<script[\s\S]*?>[\s\S]*?<\/script>/i,
                /javascript\s*:/i,
                /on\w+\s*=/i // onclick, onload, etc.
            ];

            const containsSuspiciousContent = suspiciousPatterns.some(pattern => 
                pattern.test(trimmedCodeSnippet)
            );

            if (containsSuspiciousContent) {
                errorMessage = "Code snippet contains potentially unsafe content. Please remove any executable code, SQL commands, or HTML/JavaScript injections.";
                submitState = 'error';
                return;
            }
        }
        
        submitState = 'loading';
        errorMessage = '';
        
        try {
            // Parse tags and filter out empty ones
            const tagList = tags.split(',')
                .map(t => t.trim())
                .filter(t => t.length > 0);
            
            const confessionData: ConfessionRequest = {
                title: trimmedTitle,
                description: trimmedDescription,
                snippet: trimmedCodeSnippet || undefined,
                language: language.trim(),
                tags: tagList.length > 0 ? tagList : undefined,
                isFlagged: isFlagged
            };
            
            const createdConfession = await api.confessions.create(confessionData);
            
            submitState = 'success';
            
            // Reset form
            title = "";
            description = "";
            codeSnippet = "";
            language = "JavaScript";
            tags = "";
            isFlagged = false;
            
            // Show success message and redirect
            showMacNotification("Thank you for sharing your confession! 🐛", 'success', '✅');
            
            // Navigate to confessions page after a brief delay
            setTimeout(() => {
                goto('/confessions');
            }, 1500);
            
        } catch (error) {
            submitState = 'error';
            errorMessage = api.handleError(error as Error);
        }
    }
    
    function resetForm() {
        title = "";
        description = "";
        codeSnippet = "";
        language = "JavaScript";
        tags = "";
        isFlagged = false;
        submitState = 'idle';
        errorMessage = '';
        tagSuggestions = [];
        showSuggestions = false;
    }
    
    async function searchTagSuggestions(query: string) {
        if (!query || query.length < 1) {
            tagSuggestions = [];
            showSuggestions = false;
            return;
        }
        
        isSearchingTags = true;
        
        try {
            const results = await api.tags.suggest(query);
            tagSuggestions = results.map(tag => tag.name);
            showSuggestions = tagSuggestions.length > 0;
        } catch (error) {
            console.error('Tag search failed:', error);
            tagSuggestions = [];
            showSuggestions = false;
        } finally {
            isSearchingTags = false;
        }
    }
    
    function handleTagInput(event: Event) {
        const target = event.target as HTMLInputElement;
        const value = target.value;
        
        // Get the last tag being typed (after the last comma)
        const lastCommaIndex = value.lastIndexOf(',');
        const currentTag = lastCommaIndex === -1 
            ? value.trim() 
            : value.substring(lastCommaIndex + 1).trim();
        
        if (currentTag) {
            searchTagSuggestions(currentTag);
        } else {
            tagSuggestions = [];
            showSuggestions = false;
        }
    }
    
    function selectTagSuggestion(suggestion: string) {
        const lastCommaIndex = tags.lastIndexOf(',');
        
        if (lastCommaIndex === -1) {
            // No comma, replace entire value
            tags = suggestion;
        } else {
            // Replace the part after the last comma
            tags = tags.substring(0, lastCommaIndex + 1) + ' ' + suggestion;
        }
        
        tagSuggestions = [];
        showSuggestions = false;
    }
    
    function handleTagInputBlur() {
        // Delay hiding suggestions to allow for clicks
        setTimeout(() => {
            showSuggestions = false;
        }, 200);
    }
</script>

<Page title="Submit Your Story">
    <section class="submit-container">
        <div class="hero-section">
            <h2>Share Your Coding Confession</h2>
            <p>Help others learn from your debugging adventures and coding mishaps</p>
        </div>
        
        <Gap size={20} />
        
        <div class="form-container">
            <Group title="Basic Information">
                <div class="input-container">
                    <label class="label">Title *</label>
                    <input 
                        type="text"
                        bind:value={title} 
                        placeholder="e.g., Spent 3 hours debugging... forgot semicolon"
                        class="title-input"
                        minlength="5"
                        maxlength="100"
                        required
                    />
                    <div class="field-info">
                        <span class="char-count">{title.length}/100 characters</span>
                        {#if title.length < 5 && title.length > 0}
                            <span class="validation-error">Minimum 5 characters required</span>
                        {/if}
                    </div>
                </div>
                <Gap size={10} />
                <div class="textarea-container">
                    <label class="label">Description *</label>
                    <textarea 
                        bind:value={description}
                        placeholder="Tell us your story... What happened? How did you solve it? What did you learn?"
                        rows="6"
                        class="description-input"
                        minlength="10"
                        required
                    ></textarea>
                    <div class="field-info">
                        <span class="char-count">{description.length} characters</span>
                        {#if description.length < 10 && description.length > 0}
                            <span class="validation-error">Minimum 10 characters required</span>
                        {/if}
                    </div>
                </div>
            </Group>
            
            <Gap size={15} />
            
            <Group title="Technical Details">
                <div class="input-container">
                    <label class="label">Programming Language *</label>
                    <select bind:value={language} class="language-dropdown" required>
                        {#each languages as lang}
                            <option value={lang}>{lang}</option>
                        {/each}
                    </select>
                </div>
                <Gap size={10} />
                <div class="textarea-container">
                    <label class="label">Code Snippet (Optional)</label>
                    <div class="code-editor-container">
                        <div class="code-editor-wrapper">
                            <div class="code-editor-header">
                                <span class="code-editor-title">Code Editor</span>
                                <div class="code-editor-info">
                                    <span class="char-count">{codeSnippet.length}/2000 characters</span>
                                    <span class="lines-count">{codeSnippet.split('\n').length} lines</span>
                                </div>
                            </div>
                            <div class="code-editor" class:has-content={codeSnippet.length > 0}>
                                <!-- Syntax highlighted background -->
                                <pre class="code-highlight" aria-hidden="true"><code bind:this={highlightedCode}></code></pre>
                                <!-- Actual textarea for input -->
                                <textarea 
                                    bind:value={codeSnippet}
                                    placeholder={codePlaceholder}
                                    class="code-input"
                                    spellcheck="false"
                                    oninput={handleCodeInput}
                                    onkeydown={handleCodeKeydown}
                                    onscroll={syncScroll}
                                    maxlength="2000"
                                    bind:this={codeTextarea}
                                ></textarea>
                            </div>
                            {#if codeSnippet.length >= 1900}
                                <div class="warning-message">
                                    ⚠️ Approaching character limit. Consider shortening your code snippet.
                                </div>
                            {/if}
                            {#if codeSnippet.length >= 2000}
                                <div class="error-message">
                                    🚫 Character limit reached. Please shorten your code snippet.
                                </div>
                            {/if}
                        </div>
                    </div>
                </div>
                <Gap size={10} />
                <div class="input-container">
                    <label class="label">Tags (Optional)</label>
                    <div class="tags-input-container">
                        <input 
                            type="text"
                            bind:value={tags} 
                            placeholder="debugging, production, css, semicolon (comma separated)"
                            class="tags-input"
                            oninput={handleTagInput}
                            onblur={handleTagInputBlur}
                            onfocus={(e) => handleTagInput(e)}
                        />
                        {#if isSearchingTags}
                            <div class="search-indicator">🔄</div>
                        {/if}
                        
                        {#if showSuggestions && tagSuggestions.length > 0}
                            <div class="tag-suggestions">
                                {#each tagSuggestions as suggestion}
                                    <button 
                                        type="button"
                                        class="tag-suggestion"
                                        onclick={() => selectTagSuggestion(suggestion)}
                                    >
                                        {suggestion}
                                    </button>
                                {/each}
                            </div>
                        {/if}
                    </div>
                    <div class="field-info">
                        <span class="help-text">Separate tags with commas • Start typing for suggestions</span>
                    </div>
                </div>
                <Gap size={10} />
                <div class="checkbox-container">
                    <label class="checkbox-label">
                        <input 
                            type="checkbox" 
                            bind:checked={isFlagged}
                            class="flag-checkbox"
                        />
                        <span class="checkbox-text">Flag this confession for review</span>
                    </label>
                    <div class="field-info">
                        <span class="help-text">Check this if your confession contains sensitive content</span>
                    </div>
                </div>
            </Group>
            
            <Gap size={20} />
            
            <!-- Error Message -->
            {#if submitState === 'error' && errorMessage}
                <div class="error-message">
                    <p>❌ {errorMessage}</p>
                </div>
                <Gap size={15} />
            {/if}
            
            <div class="submit-section">
                <button 
                    type="button"
                    onclick={handleSubmit}
                    class="submit-button"
                    disabled={submitState === 'loading'}
                >
                    {#if submitState === 'loading'}
                        Submitting...
                    {:else}
                        Submit Confession
                    {/if}
                </button>
                
                {#if submitState === 'error'}
                    <button 
                        type="button"
                        onclick={resetForm}
                        class="reset-button"
                    >
                        Reset Form
                    </button>
                {/if}
                
                <p class="disclaimer">
                    By submitting, you agree to share your story with the community to help other developers learn and grow.
                </p>
            </div>
        </div>
    </section>
</Page>

<!-- macOS-style Notification -->
{#if showNotification}
    <div class="mac-notification" class:show={showNotification} class:success={notificationType === 'success'} class:error={notificationType === 'error'} class:info={notificationType === 'info'}>
        <div class="notification-content">
            <div class="notification-icon">{notificationIcon}</div>
            <div class="notification-text">
                <div class="notification-title">Shit Happens</div>
                <div class="notification-message">{notificationMessage}</div>
            </div>
            <button class="notification-close" onclick={hideNotification} aria-label="Close notification">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
            </button>
        </div>
    </div>
{/if}

<style>
    .submit-container {
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
    }
    
    .form-container {
        width: 100%;
    }
    
    .input-container {
        display: flex;
        flex-direction: column;
        gap: 8px;
    }
    
    .textarea-container {
        display: flex;
        flex-direction: column;
        gap: 8px;
    }
    
    .label {
        color: var(--font-color);
        font-weight: 500;
        font-size: 0.95rem;
    }
    
    .tags-input {
        width: 100%;
        padding: 12px;
        border: 1px solid var(--border-input);
        border-radius: 8px;
        background: var(--bg-input);
        color: var(--font-color);
        font-family: inherit;
        font-size: 0.95rem;
    }
    
    .tags-input-container {
        position: relative;
        width: 100%;
    }
    
    .search-indicator {
        position: absolute;
        right: 15px;
        top: 50%;
        transform: translateY(-50%);
        font-size: 1rem;
        animation: spin 1s linear infinite;
    }
    
    @keyframes spin {
        from { transform: translateY(-50%) rotate(0deg); }
        to { transform: translateY(-50%) rotate(360deg); }
    }
    
    .tag-suggestions {
        position: absolute;
        top: 100%;
        left: 0;
        right: 0;
        background: var(--bg-modal);
        border: 1px solid var(--border-input);
        border-top: none;
        border-radius: 0 0 8px 8px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
        z-index: 1000;
        max-height: 200px;
        overflow-y: auto;
    }
    
    .tag-suggestion {
        display: block;
        width: 100%;
        padding: 10px 15px;
        background: none;
        border: none;
        color: var(--font-color);
        font-family: inherit;
        font-size: 0.9rem;
        text-align: left;
        cursor: pointer;
        transition: background-color 0.2s ease;
        border-bottom: 1px solid var(--border-card);
    }
    
    .tag-suggestion:last-child {
        border-bottom: none;
    }
    
    .tag-suggestion:hover {
        background: var(--bg-input-focus);
        color: #667eea;
    }

    /* Code Editor Styles */
    .code-editor-container {
        width: 100%;
    }

    .code-editor-wrapper {
        border: 2px solid var(--border-input);
        border-radius: 12px;
        overflow: hidden;
        background: var(--bg-input);
        transition: border-color 0.3s ease;
    }

    .code-editor-wrapper:focus-within {
        border-color: #667eea;
        box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }

    .code-editor-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 16px;
        background: rgba(44, 39, 51, 0.5);
        border-bottom: 1px solid var(--border-input);
        font-family: var(--font-family-mono);
    }

    .code-editor-title {
        color: var(--font-color-accent);
        font-size: 0.9rem;
        font-weight: 500;
    }

    .code-editor-info {
        display: flex;
        gap: 15px;
        font-size: 0.8rem;
        color: var(--font-color-muted);
    }

    .code-editor {
        position: relative;
        font-family: var(--font-family-mono);
        font-size: 14px;
        line-height: 1.5;
        height: 300px;
    }

    .code-editor.has-content {
        background: rgba(0, 0, 0, 0.2);
    }

    .code-highlight {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        margin: 0;
        padding: 20px;
        background: transparent;
        border: none;
        overflow: auto;
        white-space: pre-wrap;
        word-wrap: break-word;
        color: transparent;
        caret-color: transparent;
        pointer-events: none;
        z-index: 1;
    }

    .code-input {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        width: 100%;
        height: 100%;
        margin: 0;
        padding: 20px;
        background: transparent;
        border: none;
        outline: none;
        font-family: inherit;
        font-size: inherit;
        line-height: inherit;
        color: var(--font-color);
        resize: none;
        white-space: pre-wrap;
        word-wrap: break-word;
        overflow: auto;
        z-index: 2;
    }

    .code-input::placeholder {
        color: var(--font-color-muted);
        opacity: 0.6;
    }

    /* Syntax highlighting colors */
    .code-highlight :global(.comment) {
        color: #6a994e;
        font-style: italic;
    }

    .code-highlight :global(.string) {
        color: #f4a261;
    }

    .code-highlight :global(.keyword) {
        color: #e76f51;
        font-weight: bold;
    }

    .code-highlight :global(.number) {
        color: #2a9d8f;
    }

    .code-highlight :global(.operator) {
        color: #e9c46a;
    }

    /* Warning and error messages */
    .warning-message {
        padding: 10px 16px;
        background: rgba(255, 193, 7, 0.1);
        border-top: 1px solid rgba(255, 193, 7, 0.3);
        color: var(--color-warning);
        font-size: 0.9rem;
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .error-message {
        padding: 10px 16px;
        background: rgba(220, 38, 38, 0.1);
        border-top: 1px solid rgba(220, 38, 38, 0.3);
        color: var(--color-danger);
        font-size: 0.9rem;
        display: flex;
        align-items: center;
        gap: 8px;
    }

    /* Responsive adjustments for code editor */
    @media (max-width: 768px) {
        .code-editor {
            height: 250px;
        }
        
        .code-input, .code-highlight {
            padding: 15px;
            font-size: 13px;
        }
        
        .code-editor-info {
            flex-direction: column;
            gap: 5px;
            align-items: flex-end;
        }
    }
    
    .title-input,
    .language-dropdown {
        width: 100%;
        padding: 12px;
        border: 1px solid var(--border-input);
        border-radius: 8px;
        background: var(--bg-input);
        color: var(--font-color);
        font-family: inherit;
        font-size: 0.95rem;
    }
    
    .title-input:focus,
    .tags-input:focus,
    .language-dropdown:focus {
        outline: none;
        border-color: var(--border-input-focus);
        background: var(--bg-input-focus);
    }
    
    .language-dropdown {
        width: 100%;
        padding: 12px;
        border: 1px solid var(--border-input);
        border-radius: 8px;
        background: var(--bg-input);
        color: var(--font-color);
        font-family: inherit;
        font-size: 0.95rem;
        cursor: pointer;
    }
    
    .language-dropdown:focus {
        outline: none;
        border-color: var(--border-input-focus);
        background: var(--bg-input-focus);
    }
    
    .checkbox-container {
        display: flex;
        flex-direction: column;
        gap: 8px;
    }
    
    .checkbox-label {
        display: flex;
        align-items: center;
        gap: 10px;
        cursor: pointer;
        font-size: 0.95rem;
        color: var(--font-color);
    }
    
    .flag-checkbox {
        width: 18px;
        height: 18px;
        cursor: pointer;
    }
    
    .checkbox-text {
        user-select: none;
    }
    
    .field-info {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 0.8rem;
        margin-top: 4px;
    }
    
    .char-count {
        color: var(--font-color-accent);
    }
    
    .help-text {
        color: var(--font-color-accent);
        font-style: italic;
    }
    
    .validation-error {
        color: #ff6b6b;
        font-weight: 500;
    }
    
    .description-input,
    .code-input {
        width: 100%;
        padding: 12px;
        border: 1px solid var(--border-input);
        border-radius: 8px;
        background: var(--bg-input);
        color: var(--font-color);
        font-family: inherit;
        font-size: 0.95rem;
        line-height: 1.5;
        resize: vertical;
        min-height: 120px;
    }
    
    .code-input {
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 0.9rem;
        background: #1e1e1e;
        color: #f8f8f2;
        border-color: #444;
        min-height: 160px;
    }
    
    .description-input:focus,
    .code-input:focus {
        outline: none;
        border-color: var(--border-input-focus);
        background: var(--bg-input-focus);
    }
    
    .submit-section {
        text-align: center;
    }
    
    .error-message {
        text-align: center;
        padding: 15px;
        background: rgba(255, 0, 0, 0.1);
        border: 1px solid rgba(255, 0, 0, 0.2);
        border-radius: 8px;
        color: var(--font-color);
    }
    
    .submit-button {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        border: none;
        padding: 14px 32px;
        border-radius: 8px;
        font-size: 1.1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.3s ease;
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        margin-right: 10px;
    }
    
    .submit-button:disabled {
        opacity: 0.6;
        cursor: not-allowed;
        transform: none !important;
    }
    
    .submit-button:hover:not(:disabled) {
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
    }
    
    .submit-button:active:not(:disabled) {
        transform: translateY(0);
    }
    
    .reset-button {
        background: var(--bg-input);
        border: 1px solid var(--border-input);
        color: var(--font-color);
        padding: 14px 24px;
        border-radius: 8px;
        font-size: 1rem;
        cursor: pointer;
        transition: all 0.2s ease;
        margin-left: 10px;
    }
    
    .reset-button:hover {
        background: var(--bg-input-focus);
        border-color: var(--border-input-focus);
    }
    
    .disclaimer {
        margin-top: 15px;
        color: var(--font-color-accent);
        font-size: 0.9rem;
        line-height: 1.4;
    }

    /* macOS-style Notification */
    .mac-notification {
        position: fixed;
        top: 20px;
        right: 20px;
        width: 380px;
        background: rgba(60, 60, 60, 0.95);
        backdrop-filter: blur(40px);
        border-radius: 12px;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3), 0 0 0 1px rgba(255, 255, 255, 0.1);
        z-index: 9999;
        transform: translateX(420px);
        opacity: 0;
        transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
        border: 1px solid rgba(255, 255, 255, 0.1);
    }

    .mac-notification.show {
        transform: translateX(0);
        opacity: 1;
    }

    .mac-notification.success {
        border-left: 4px solid #34d399;
    }

    .mac-notification.error {
        border-left: 4px solid #ef4444;
    }

    .mac-notification.info {
        border-left: 4px solid #3b82f6;
    }

    .notification-content {
        display: flex;
        align-items: flex-start;
        gap: 12px;
        padding: 16px;
    }

    .notification-icon {
        font-size: 24px;
        line-height: 1;
        flex-shrink: 0;
        margin-top: 2px;
    }

    .notification-text {
        flex: 1;
        min-width: 0;
    }

    .notification-title {
        font-size: 14px;
        font-weight: 600;
        color: white;
        margin-bottom: 2px;
        line-height: 1.2;
    }

    .notification-message {
        font-size: 13px;
        color: rgba(255, 255, 255, 0.85);
        line-height: 1.3;
        word-wrap: break-word;
    }

    .notification-close {
        flex-shrink: 0;
        background: none;
        border: none;
        color: rgba(255, 255, 255, 0.6);
        cursor: pointer;
        padding: 4px;
        border-radius: 6px;
        transition: all 0.2s ease;
        margin-top: -2px;
        margin-right: -4px;
    }

    .notification-close:hover {
        color: white;
        background: rgba(255, 255, 255, 0.1);
    }

    .notification-close:active {
        transform: scale(0.95);
    }

    /* Responsive adjustments for notification */
    @media (max-width: 768px) {
        .mac-notification {
            width: calc(100vw - 40px);
            right: 20px;
            left: 20px;
            transform: translateY(-120px);
        }

        .mac-notification.show {
            transform: translateY(0);
        }
    }

    @media (max-width: 480px) {
        .notification-content {
            padding: 12px;
            gap: 10px;
        }

        .notification-icon {
            font-size: 20px;
        }

        .notification-title {
            font-size: 13px;
        }

        .notification-message {
            font-size: 12px;
        }
    }
</style>
