<script lang="ts">
    // Placeholder for tutorial functionality
</script>

<div class="phase-tutorial">
    <div class="tutorial-content">
        <h2>Optimizer Tutorial</h2>
        <p class="description">
            Tutorial content will be customized here...
        </p>
    </div>
</div>

<style>
    .phase-tutorial {
        padding: 1rem;
        height: 100%;
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }

    .tutorial-content {
        display: flex;
        flex-direction: column;
        height: 100%;
        min-height: 0;
    }

    h2 {
        color: #8B5CF6;
        margin-bottom: 1rem;
        font-size: 1.5rem;
    }

    h3 {
        color: #7C3AED;
        margin: 1.5rem 0 0.5rem 0;
    }

    .description {
        line-height: 1.6;
        color: #333;
        margin: 1rem 0;
        font-size: 0.95rem;
    }

    .separator {
        border-bottom: 1px solid #e5e7eb;
        margin: 1rem 0;
        width: 100%;
    }

    .tutorial-container {
        display: flex;
        flex-direction: column;
        height: 100%;
        min-height: 0;
        padding: 1rem;
    }

    .tutorial-content-area {
        flex: 1;
        min-height: 0;
        overflow-y: auto;
        padding-bottom: 2rem;
    }

    .tutorial-step {
        margin-bottom: 2rem;
    }

    .optimization-types {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        margin: 1rem 0;
    }

    .opt-type {
        background: #F3E8FF;
        padding: 0.75rem;
        border-radius: 8px;
        border-left: 4px solid #8B5CF6;
    }

    .level-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 1rem;
        margin: 1rem 0;
    }

    .level-card {
        background: #F8FAFC;
        border: 1px solid #E2E8F0;
        border-radius: 8px;
        padding: 1rem;
        text-align: center;
    }

    .level-header {
        font-weight: 600;
        color: #8B5CF6;
        margin-bottom: 0.5rem;
    }

    .level-desc {
        font-size: 0.9rem;
        color: #64748B;
    }

    .technique-list {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        margin: 1rem 0;
    }

    .technique {
        background: #F1F5F9;
        padding: 0.75rem;
        border-radius: 6px;
    }

    .impact-section {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.5rem;
        margin: 1rem 0;
    }

    .benefits, .trade-offs {
        background: #FAFAFA;
        padding: 1rem;
        border-radius: 8px;
    }

    .benefits h4 {
        color: #059669;
        margin: 0 0 0.5rem 0;
    }

    .trade-offs h4 {
        color: #DC2626;
        margin: 0 0 0.5rem 0;
    }

    .navigation-container {
        position: sticky;
        margin-top: 1rem;
        bottom: 0;
        background: #fff;
        z-index: 10;
        padding: 0rem 0;
        display: flex;
        justify-content: center;
    }

    .navigation {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;
        max-width: 500px;
    }

    .nav-button {
        padding: 0.5rem 1rem;
        background: #C8A8E9;
        color: #4C1D95;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        min-width: 100px;
        font-size: 0.9rem;
    }

    .nav-button:disabled {
        background: #cccccc;
        cursor: not-allowed;
    }

    .step-counter {
        color: #666;
        font-size: 0.9rem;
    }

    /* Dark Mode Styles */
    :global(html.dark-mode) .phase-tutorial {
        background: #1a2a4a;
    }

    :global(html.dark-mode) h2 {
        color: #C4B5FD;
    }

    :global(html.dark-mode) h3 {
        color: #A78BFA;
    }

    :global(html.dark-mode) .description {
        color: #d1d5db;
    }

    :global(html.dark-mode) .separator {
        border-bottom-color: #374151;
    }

    :global(html.dark-mode) .opt-type {
        background: #2D1B69;
        border-left-color: #8B5CF6;
        color: #E0E7FF;
    }

    :global(html.dark-mode) .level-card {
        background: #2d3748;
        border-color: #4a5568;
    }

    :global(html.dark-mode) .level-header {
        color: #C4B5FD;
    }

    :global(html.dark-mode) .level-desc {
        color: #CBD5E0;
    }

    :global(html.dark-mode) .technique {
        background: #2d3748;
        color: #E2E8F0;
    }

    :global(html.dark-mode) .benefits,
    :global(html.dark-mode) .trade-offs {
        background: #2d3748;
        color: #E2E8F0;
    }

    :global(html.dark-mode) .navigation-container {
        background: #1a2a4a;
    }

    :global(html.dark-mode) .nav-button {
        background: #8B5CF6;
        color: #F3E8FF;
    }

    :global(html.dark-mode) .nav-button:disabled {
        background: #4b5563;
        color: #9ca3af;
    }

    :global(html.dark-mode) .step-counter {
        color: #d1d5db;
    }
</style>
