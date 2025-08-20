<script lang="ts">
    import type { Token } from '$lib/types';
    import { lexerState } from '$lib/stores/lexer';
    import { onMount, onDestroy } from 'svelte';
    import { get } from 'svelte/store';

    export let phase: string;
    export let tokens: Token[] = [];
    export let unexpected_tokens: string[] = [];
    export let show_tokens = false;

    // Subscribe to lexer state
    const unsubscribe = lexerState.subscribe(state => {
        if (state.tokens) {
            tokens = state.tokens;
            show_tokens = true;
        }
        if (state.tokens_unidentified) {
            unexpected_tokens = state.tokens_unidentified || [];
        }
    });

    onMount(() => {
        // Check initial lexer state for tokens
        const state = get(lexerState);
        if (state.tokens) {
            tokens = state.tokens;
            show_tokens = true;
        }
        if (state.tokens_unidentified) {
            unexpected_tokens = state.tokens_unidentified || [];
        }
    });

    onDestroy(() => {
        unsubscribe();
    });
</script>

<div class="artifact-viewer">
    <div class="artefact-heading">
            <h2 class="artifact-title">Lexer Artefact</h2>
        </div>

        <div class="tokens-heading">
                    <h3>Tokens</h3>
		</div>
    {#if phase === 'lexer' && show_tokens}

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
                    <br>
                    Unidentified Token(s):
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
        transition: background-color 0.3s ease;
        padding: 2rem 1.5rem 1.5rem;
		background-color: #f8f9fa;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
		display: flex;
		flex-direction: column;
		gap: 1rem;
		height: 100%;
    }

    h3 {
		color: #001a6e;
		font-size: 1.5rem;
		margin: 0;
		font-family: 'Times New Roman', serif;
	}

    .token-table {
        width: 100%;
        table-layout: fixed;
        border-collapse: collapse;
         margin: 0rem 0 1rem 0;
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

    .token-list {
        display: flex;
        flex-wrap: wrap;
        gap: 1rem;
    }

    .artefact-heading {
        border-bottom: 1px solid #e5e7eb;
		padding-bottom: 0.75rem;
    }

    .tokens-heading {
		border-bottom: none;
		padding-bottom: 0;
		text-align: center;
	}

    .artifact-title {
		margin: 0;
		color: #001a6e;
		font-family: 'Times New Roman', serif;
		font-size: 1.25rem;
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

    :global(html.dark-mode) .tokens-heading {
		border-bottom-color: #4a5568;
	}
    :global(html.dark-mode) .artifact-title,
	:global(html.dark-mode) h3 {
		color: #ebeef1;
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
