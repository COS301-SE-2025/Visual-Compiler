<script lang="ts">
    // This script defines the component's props, which receive data from the parent.
    // It expects the current compiler phase, a list of generated tokens,
    // a list of any unexpected tokens, and a boolean to control visibility.
    import type { Token } from '$lib/types';

    export let phase: string;
    export let tokens: Token[] = [];
    export let unexpected_tokens: string[] = [];
    export let show_tokens = false;
</script>

<div class="artifact-viewer">
    {#if phase === 'lexer' && show_tokens}
        <div class="tokens-heading">
            <h3>Tokens</h3>
        </div>
        {#if tokens && tokens.length > 0}
            <table class="token-table">
                <thead>
                    <tr>
                        <th>Type</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                    {#each tokens as token}
                        <tr>
                            <td>{token.type}</td>
                            <td>{token.value}</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        {:else}
            <div class="no-tokens">No tokens generated yet</div>
        {/if}

        {#if unexpected_tokens && unexpected_tokens.length > 0}
            <div class="unexpected-tokens-container">
                <h4>Unidentified Input</h4>
                <div class="token-list">
                    {#each unexpected_tokens as token}
                        <span class="unexpected-token">{token}</span>
                    {/each}
                </div>
            </div>
        {/if}
    {:else if phase === 'lexer'}
        <div class="empty-state">Tokens will appear here after generation</div>
    {/if}
</div>

<style>
    /* --- Base & Light Mode Styles --- */
    .artifact-viewer {
        padding: 2rem 1.5rem 1.5rem;
        transition: background-color 0.3s ease;
    }

    h3 {
        color: #041a47;
        font-size: 1.25rem;
        margin: 0 0 1.5rem 0;
        padding-bottom: 1rem;
        border-bottom: 2px solid #e5e7eb;
        transition: color 0.3s ease, border-bottom-color 0.3s ease;
    }

    .token-table {
        width: 100%;
        table-layout: fixed;
        border-collapse: collapse;
        margin: 1.5rem 0 2rem 0;
        background: white;
        border-radius: 10px; /* Updated border-radius */
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        overflow: hidden; /* Ensures child corners are clipped */
        transition: background-color 0.3s ease, box-shadow 0.3s ease;
    }

    .token-table th,
    .token-table td {
        width: 50%;
        padding: 0.75rem 1rem;
        border-bottom: 1px solid #e0e0e0;
        text-align: left;
        transition: border-color 0.3s ease, color 0.3s ease;
    }
    
    .token-table tr:last-child td {
        border-bottom: none; /* Removes border from last row for cleaner look */
    }

    .token-table th {
        background: #041a47;
        color: white;
        font-weight: 500;
        transition: background-color 0.3s ease, color 0.3s ease;
    }

    .unexpected-tokens-container {
        margin-top: 1.5rem;
        padding: 1.2rem;
        background: #f8f9fa;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        transition: background-color 0.3s ease, border-color 0.3s ease;
    }

    .unexpected-tokens-container h4 {
        color: #041a47;
        font-size: 1.1rem;
        margin: 0 0 1rem 0;
        padding-bottom: 0.6rem;
        border-bottom: 1px solid #e5e7eb;
        transition: color 0.3s ease, border-bottom-color 0.3s ease;
    }

    .token-list {
        display: flex;
        flex-wrap: wrap;
        gap: 1rem;
    }

    .tokens-heading {
        color: #041a47;
        transition: color 0.3s ease;
    }

    .unexpected-token {
        padding: 0.5rem 1rem;
        background: white;
        border: 1px solid #e5e7eb;
        border-radius: 6px;
        color: #041a47;
        font-family: monospace;
        font-size: 0.9rem;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
        transition: background-color 0.3s ease, border-color 0.3s ease, color 0.3s ease;
    }

    .empty-state,
    .no-tokens {
        color: #666;
        font-style: italic;
        text-align: center;
        padding: 2rem;
        transition: color 0.3s ease;
    }

    /* --- Dark Mode Styles --- */
    :global(html.dark-mode) .artifact-viewer {
        background: #1a2a4a;
    }

    :global(html.dark-mode) h3,
    :global(html.dark-mode) .tokens-heading {
        color: #e2e8f0;
        border-bottom-color: #4a5568;
    }

    :global(html.dark-mode) .token-table {
        background: #2d3748;
    }

    :global(html.dark-mode) .token-table th,
    :global(html.dark-mode) .token-table td {
        border-color: #4a5568;
        color: #e2e8f0;
    }

    :global(html.dark-mode) .token-table th {
        background: #1a202c;
    }

    :global(html.dark-mode) .unexpected-tokens-container {
        background: #2d3748;
        border-color: #4a5568;
    }

    :global(html.dark-mode) .unexpected-tokens-container h4 {
        color: #90cdf4;
        border-bottom-color: #4a5568;
    }

    :global(html.dark-mode) .unexpected-token {
        background: #4a5568;
        border-color: #718096;
        color: #e2e8f0;
    }

    :global(html.dark-mode) .empty-state,
    :global(html.dark-mode) .no-tokens {
        color: #a0aec0;
    }
</style>
