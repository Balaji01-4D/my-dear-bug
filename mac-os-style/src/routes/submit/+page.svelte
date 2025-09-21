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
    let submitState: LoadingState = $state('idle');
    let errorMessage = $state('');
    
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

    async function handleSubmit() {
        if (!title.trim() || !description.trim()) {
            alert("Please fill in at least the title and description!");
            return;
        }
        
        submitState = 'loading';
        errorMessage = '';
        
        try {
            const confessionData: ConfessionRequest = {
                title: title.trim(),
                description: description.trim(),
                snippet: codeSnippet.trim() || undefined,
                language,
                tags: tags.split(',').map(t => t.trim()).filter(t => t)
            };
            
            const createdConfession = await api.confessions.create(confessionData);
            
            submitState = 'success';
            
            // Reset form
            title = "";
            description = "";
            codeSnippet = "";
            language = "JavaScript";
            tags = "";
            
            // Show success message and redirect
            alert("Thank you for sharing your confession! 🐛");
            
            // Navigate to confessions page or to the specific confession
            goto('/confessions');
            
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
        submitState = 'idle';
        errorMessage = '';
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
                <Text 
                    name="Title" 
                    bind:value={title} 
                    placeholder="e.g., Spent 3 hours debugging... forgot semicolon"
                />
                <Gap size={10} />
                <div class="textarea-container">
                    <label class="label">Description *</label>
                    <textarea 
                        bind:value={description}
                        placeholder="Tell us your story... What happened? How did you solve it? What did you learn?"
                        rows="6"
                        class="description-input"
                    ></textarea>
                </div>
            </Group>
            
            <Gap size={15} />
            
            <Group title="Technical Details">
                <Dropdown 
                    name="Programming Language" 
                    bind:value={language}
                    options={languages.map(lang => ({ value: lang, label: lang }))}
                />
                <Gap size={10} />
                <div class="textarea-container">
                    <label class="label">Code Snippet (Optional)</label>
                    <textarea 
                        bind:value={codeSnippet}
                        placeholder="Paste the problematic code here..."
                        rows="8"
                        class="code-input"
                    ></textarea>
                </div>
                <Gap size={10} />
                <Text 
                    name="Tags (Optional)" 
                    bind:value={tags} 
                    placeholder="debugging, production, css, semicolon (comma separated)"
                />
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
</style>
