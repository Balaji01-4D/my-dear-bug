<script lang="ts">
    import Page from "$lib/views/Page.svelte";
    import Group from "$lib/components/settings/Group.svelte";
    import Text from "$lib/components/settings/Text.svelte";
    import Dropdown from "$lib/components/settings/Dropdown.svelte";
    import Gap from "$lib/components/Gap.svelte";

    let title = $state("");
    let description = $state("");
    let codeSnippet = $state("");
    let language = $state("JavaScript");
    let tags = $state("");
    
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

    function handleSubmit() {
        if (!title.trim() || !description.trim()) {
            alert("Please fill in at least the title and description!");
            return;
        }
        
        // Here you would normally send to your API
        console.log("Submitting confession:", {
            title: title.trim(),
            description: description.trim(),
            snippet: codeSnippet.trim(),
            language,
            tags: tags.split(',').map(t => t.trim()).filter(t => t)
        });
        
        alert("Thank you for sharing your confession! 🐛");
        
        // Reset form
        title = "";
        description = "";
        codeSnippet = "";
        language = "JavaScript";
        tags = "";
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
            
            <div class="submit-section">
                <button 
                    type="button"
                    onclick={handleSubmit}
                    class="submit-button"
                >
                    Submit Confession
                </button>
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
    }
    
    .submit-button:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
    }
    
    .submit-button:active {
        transform: translateY(0);
    }
    
    .disclaimer {
        margin-top: 15px;
        color: var(--font-color-accent);
        font-size: 0.9rem;
        line-height: 1.4;
    }
</style>
