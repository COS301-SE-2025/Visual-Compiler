<script lang="ts">
    // Placeholder for optimizer artifacts
</script>

<div class="artifact-container">
    <div class="artifact-header">
        <h3 class="artifact-title">Optimizer Artifacts</h3>
    </div>

    <div class="artifact-content">
        <p>Artifact viewer content will be customized here...</p>
    </div>
</div>

<style>
    .artifact-container {
        padding: 1.5rem;
        background-color: #F8FAFC;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        display: flex;
        flex-direction: column;
        gap: 1rem;
        height: 100%;
    }

    .artifact-header {
        border-bottom: 1px solid #E5E7EB;
        padding-bottom: 0.75rem;
    }

    .artifact-title {
        margin: 0;
        color: #8B5CF6;
        font-family: 'Times New Roman', serif;
        font-size: 1.25rem;
        text-align: center;
    }

    .optimization-results {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
        flex: 1;
    }

    .section-title {
        color: #7C3AED;
        margin: 0 0 1rem 0;
        font-size: 1.1rem;
        font-weight: 600;
    }

    .performance-section {
        background: #F3E8FF;
        padding: 1rem;
        border-radius: 8px;
        border-left: 4px solid #8B5CF6;
    }

    .metrics-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
        gap: 1rem;
    }

    .metric-card {
        background: white;
        padding: 1rem;
        border-radius: 6px;
        text-align: center;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .metric-label {
        font-size: 0.9rem;
        color: #6B7280;
        margin-bottom: 0.5rem;
    }

    .metric-value {
        font-size: 1.2rem;
        font-weight: 600;
        color: #059669;
    }

    .code-comparison {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }

    .code-section {
        background: white;
        border-radius: 8px;
        overflow: hidden;
        border: 1px solid #E5E7EB;
    }

    .code-block {
        padding: 1rem;
        font-family: 'Fira Code', monospace;
        font-size: 0.85rem;
        line-height: 1.4;
        overflow-x: auto;
        max-height: 300px;
        overflow-y: auto;
    }

    .code-block.original {
        background: #FEF2F2;
        border-top: 3px solid #EF4444;
    }

    .code-block.optimized {
        background: #F0FDF4;
        border-top: 3px solid #10B981;
    }

    .applied-optimizations {
        background: white;
        padding: 1rem;
        border-radius: 8px;
        border: 1px solid #E5E7EB;
    }

    .optimization-list {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .optimization-item {
        background: #F8FAFC;
        padding: 0.75rem;
        border-radius: 6px;
        border-left: 3px solid #8B5CF6;
    }

    .optimization-name {
        display: block;
        font-weight: 600;
        color: #374151;
        margin-bottom: 0.25rem;
    }

    .optimization-desc {
        display: block;
        font-size: 0.9rem;
        color: #6B7280;
    }

    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 3rem 1rem;
        color: #8B5CF6;
        height: 100%;
        flex-grow: 1;
        text-align: center;
        font-style: italic;
    }

    .empty-state svg {
        margin-bottom: 1rem;
        opacity: 0.6;
    }

    .empty-state p {
        margin: 0;
        color: #6B7280;
    }

    /* Dark Mode Styles */
    :global(html.dark-mode) .artifact-container {
        background: #1a2a4a;
    }

    :global(html.dark-mode) .artifact-header {
        border-bottom-color: #374151;
    }

    :global(html.dark-mode) .artifact-title {
        color: #C4B5FD;
    }

    :global(html.dark-mode) .section-title {
        color: #A78BFA;
    }

    :global(html.dark-mode) .performance-section {
        background: #2D1B69;
        border-left-color: #8B5CF6;
    }

    :global(html.dark-mode) .metric-card {
        background: #374151;
        color: #E5E7EB;
    }

    :global(html.dark-mode) .metric-label {
        color: #9CA3AF;
    }

    :global(html.dark-mode) .metric-value {
        color: #34D399;
    }

    :global(html.dark-mode) .code-section {
        background: #2d3748;
        border-color: #4a5568;
    }

    :global(html.dark-mode) .code-block {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .code-block.original {
        background: #2D1B1B;
        border-top-color: #F87171;
    }

    :global(html.dark-mode) .code-block.optimized {
        background: #1B2D1B;
        border-top-color: #34D399;
    }

    :global(html.dark-mode) .applied-optimizations {
        background: #2d3748;
        border-color: #4a5568;
    }

    :global(html.dark-mode) .optimization-item {
        background: #374151;
    }

    :global(html.dark-mode) .optimization-name {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .optimization-desc {
        color: #9CA3AF;
    }

    :global(html.dark-mode) .empty-state {
        color: #A78BFA;
    }

    :global(html.dark-mode) .empty-state p {
        color: #9CA3AF;
    }
</style>
