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
        {#if unexpected_tokens && unexpected_tokens.length > 0}
            <div class="error-state">
				<div class="error-icon">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="48"
						height="48"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.5"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<circle cx="12" cy="12" r="10" />
						<line x1="12" y1="8" x2="12" y2="12" />
						<line x1="12" y1="16" x2="12.01" y2="16" />
					</svg>
				</div>
				<h4>Unidentified Token(s) Found</h4>
				<p class="error-message">
					The source code could not be completely tokenised with the provided regex/DFA. <br>
                    Please check your input again.<br>
				</p>
                <div class="token-list">
                    {#each unexpected_tokens as token}
                        <span class="error-details">{token}</span>
                    {/each}
                </div>
			</div>
        {/if}
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

        
    {:else if phase === 'lexer'}
        <div class="empty-state">Tokens will appear here after generation</div>
    {/if}
</div>

<style>
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
        border-radius: 10px; 
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        overflow: hidden;
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
        border-bottom: none; 
    }

    .token-table th {
        background: #BED2E6;
        color: 000000;
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

    .error-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		text-align: center;
		padding: 2rem;
		flex-grow: 1;
		background-color: #fff5f5;
		border: 1px solid #e53e3e;
		border-radius: 8px;
		color: #9b2c2c;
	}
	.error-icon {
		color: #e53e3e;
		margin-bottom: 1rem;
	}
	.error-state h4 {
		margin: 0 0 0.5rem 0;
		font-size: 1.25rem;
		color: #c53030;
	}
	.error-message {
		margin: 0 0 1rem 0;
		max-width: 450px;
		line-height: 1.6;
	}
	.error-details {
		background-color: #fed7d7;
		padding: 0.75rem 1rem;
		border-radius: 6px;
		font-family: 'Fira Code', monospace;
		font-size: 0.85rem;
		white-space: pre-wrap;
		word-wrap: break-word;
		max-width: 100%;
		text-align: left;
		color: #742a2a;
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
        background: #001A6E;
        color: #ffffff; 
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

    :global(html.dark-mode) .error-state {
		background-color: #2d3748;
		border-color: #e53e3e;
		color: #fca5a5;
	}
	:global(html.dark-mode) .error-icon {
		color: #fca5a5;
	}
	:global(html.dark-mode) .error-state h4 {
		color: #fc8181;
	}
	:global(html.dark-mode) .error-details {
		background-color: #4a2d2d;
		color: #fed7d7;
	}
</style>
