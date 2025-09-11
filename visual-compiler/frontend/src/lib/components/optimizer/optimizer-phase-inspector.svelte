<script lang="ts">
    let selectedLanguage = 'Java';
    let selectedTechniques: string[] = [];
    let inputCode = '';

    const languages = ['Java', 'Python', 'Go'];
    const techniques = ['Constant Folding', 'Dead Code Elimination', 'Loop Unrolling'];

    function toggleTechnique(technique: string) {
        if (selectedTechniques.includes(technique)) {
            selectedTechniques = selectedTechniques.filter(t => t !== technique);
        } else {
            selectedTechniques = [...selectedTechniques, technique];
        }
    }
</script>

<div class="inspector-container">
    <div class="inspector-header">
        <h2 class="inspector-title">Optimizing</h2>
    </div>

    <div class="language-section">
        <h3 class="section-heading">Language</h3>
        <div class="language-buttons">
            {#each languages as language}
                <button 
                    class="language-btn {selectedLanguage === language ? 'selected' : ''}"
                    on:click={() => selectedLanguage = language}
                >
                    {language}
                </button>
            {/each}
        </div>
    </div>

    <div class="technique-section">
        <h3 class="section-heading">Technique</h3>
        <div class="technique-buttons">
            {#each techniques as technique}
                <button 
                    class="technique-btn {selectedTechniques.includes(technique) ? 'selected' : ''}"
                    on:click={() => toggleTechnique(technique)}
                >
                    {technique}
                </button>
            {/each}
        </div>
    </div>

    <div class="input-section">
        <h3 class="section-heading">Input</h3>
        <div class="input-area">
            <textarea 
                bind:value={inputCode}
                placeholder="Enter your code here..."
                rows="8"
            ></textarea>
        </div>
    </div>
</div>

<style>
    .inspector-container {
        padding: 1.5rem;
        background-color: #F8FAFC;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        display: flex;
        flex-direction: column;
        gap: 2rem;
        height: 100%;
        font-family: Arial, sans-serif;
    }

    .inspector-header {
        padding-bottom: 0.75rem;
        margin-bottom: -1rem;
    }

    .inspector-title {
        margin: 0;
        color: #8B5CF6;
        font-family: 'Times New Roman', serif;
        font-size: 1.5rem;
        text-align: center;
        font-weight: bold;
    }

    .section-heading {
        margin: 0 0 1rem 0;
        color: #000;
        font-size: 1.5rem;
        font-weight: bold;
    }

    /* Language Section */
    .language-section {
        display: flex;
        flex-direction: column;
    }

    .language-buttons {
        display: flex;
        gap: 0.75rem;
        flex-wrap: wrap;
    }

    .language-btn {
        padding: 0.75rem 1.5rem;
        border: 2px solid #000;
        background: #fff;
        color: #000;
        font-size: 1rem;
        font-weight: 500;
        cursor: pointer;
        border-radius: 4px;
        transition: all 0.2s ease;
        min-width: 80px;
    }

    .language-btn:hover {
        background: #f0f0f0;
    }

    .language-btn.selected {
        background: #8B5CF6;
        color: #fff;
        border-color: #8B5CF6;
    }

    /* Technique Section */
    .technique-section {
        display: flex;
        flex-direction: column;
    }

    .technique-buttons {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .technique-btn {
        padding: 0.75rem 1rem;
        border: 2px solid #000;
        background: #fff;
        color: #000;
        font-size: 1rem;
        font-weight: 500;
        cursor: pointer;
        border-radius: 4px;
        transition: all 0.2s ease;
        text-align: left;
        width: 100%;
    }

    .technique-btn:hover {
        background: #f0f0f0;
    }

    .technique-btn.selected {
        background: #8B5CF6;
        color: #fff;
        border-color: #8B5CF6;
    }

    /* Input Section */
    .input-section {
        display: flex;
        flex-direction: column;
        flex: 1;
    }

    .input-area {
        flex: 1;
        display: flex;
        flex-direction: column;
    }

    .input-area textarea {
        width: 100%;
        min-height: 120px;
        padding: 0.75rem;
        border: 2px solid #000;
        border-radius: 4px;
        font-family: 'Courier New', monospace;
        font-size: 0.9rem;
        resize: vertical;
        background: #fff;
        color: #000;
        flex: 1;
    }

    .input-area textarea:focus {
        outline: none;
        border-color: #8B5CF6;
        box-shadow: 0 0 0 2px rgba(139, 92, 246, 0.2);
    }

    .input-area textarea::placeholder {
        color: #999;
        font-style: italic;
    }

    /* Dark Mode Styles */
    :global(html.dark-mode) .inspector-container {
        background: #1a2a4a;
    }

    :global(html.dark-mode) .inspector-title {
        color: #C4B5FD;
    }

    :global(html.dark-mode) .section-heading {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .language-btn,
    :global(html.dark-mode) .technique-btn {
        background: #2d3748;
        color: #E5E7EB;
        border-color: #4a5568;
    }

    :global(html.dark-mode) .language-btn:hover,
    :global(html.dark-mode) .technique-btn:hover {
        background: #374151;
    }

    :global(html.dark-mode) .language-btn.selected,
    :global(html.dark-mode) .technique-btn.selected {
        background: #8B5CF6;
        color: #fff;
        border-color: #8B5CF6;
    }

    :global(html.dark-mode) .input-area textarea {
        background: #2d3748;
        color: #E5E7EB;
        border-color: #4a5568;
    }

    :global(html.dark-mode) .input-area textarea:focus {
        border-color: #8B5CF6;
    }

    :global(html.dark-mode) .input-area textarea::placeholder {
        color: #9CA3AF;
    }
</style>
