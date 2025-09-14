<script lang="ts">
    import { lexerState } from '$lib/stores/lexer';
    import type { Symbol } from '$lib/types';

    export let phase: string;
    export let symbol_table: Symbol[] = [];
    export let show_symbol_table = false;
    export let analyser_error: boolean;
    export let analyser_error_details: string;

    // Add reactive statement to handle loaded symbol table
    $: if ($lexerState?.symbol_table) {
        symbol_table = $lexerState.symbol_table;
        show_symbol_table = true;
        analyser_error = false;
        analyser_error_details = '';
    }
</script>

<div class="artifact-container">
    <div class="artifact-header">
        <h2 class="artifact-title">Analyser Artefact</h2>
    </div>

    <div class="artifact-viewer">
        {#if phase === 'analyser' && analyser_error}
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
                    <h4>Semantic Error Found</h4>
                    <p class="error-message">
                        The source code could not be analysed with the provided scope rules and type rules. Please check your input again.<br>
                    </p>
                    <pre class="error-details">{analyser_error_details}</pre>
                </div>
            {/if}
        {#if phase === 'analyser' && show_symbol_table}

            {#if symbol_table && symbol_table.length > 0}

            <div class="artifact-header">
                    <h3>Symbols</h3>
                </div>
                
                <table class="symbol-table">
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>Type</th>
                            <th>Scope</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each symbol_table as symbol}
                            <tr>
                                <td>{symbol.name}</td>
                                <td>{symbol.type}</td>
                                <td>{symbol.scope}</td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            {:else}
                <div class="no-symbols">No symbols generated</div>
            {/if}
        {/if}
    </div>
</div>

<style>
    .artifact-container {
		padding: 1.5rem;
		background-color: #f8f9fa;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
		display: flex;
		flex-direction: column;
		gap: 1rem;
		height: 100%;
	}

    .artifact-header {
		border-bottom: 1px solid #e5e7eb;
		padding-bottom: 0.75rem;
	}

	.artifact-viewer .artifact-header {
		border-bottom: none;
		padding-bottom: 0;
		margin-bottom: 1.5rem;
		text-align: center;
	}

	.artifact-title {
		margin: 0;
		color: #001a6e;
		font-family: 'Times New Roman', serif;
		font-size: 1.25rem;
	}

	h3 {
		color: #001a6e;
		font-size: 1.5rem;
		margin: 0;
		font-family: 'Times New Roman', serif;
	}

    .symbol-table {
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

    .symbol-table th,
    .symbol-table td {
        width: 33%;  /* Updated to distribute space evenly among 3 columns */
        padding: 0.75rem 1rem;
        border-bottom: 1px solid #e0e0e0;
        text-align: left;
        transition: border-color 0.3s ease, color 0.3s ease;
    }
    
    .symbol-table tr:last-child td {
        border-bottom: none; 
    }

    .symbol-table th {
        background: #BED2E6;
        color: 000000;
        font-weight: 500;
        transition: background-color 0.3s ease, color 0.3s ease;
    }


   
    .symbol-heading {
        color: #041a47;
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


    .empty-state,
    .no-symbols {
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
    :global(html.dark-mode) .symbol-heading {
        color: #e2e8f0;
        border-bottom-color: #4a5568;
    }

    :global(html.dark-mode) .symbol-table {
        background: #2d3748;
    }

    :global(html.dark-mode) .symbol-table th,
    :global(html.dark-mode) .symbol-table td {
        border-color: #4a5568;
        color: #e2e8f0;
    }

    :global(html.dark-mode) .symbol-table th {
        background: #001A6E;
        color: #ffffff;
    }

    :global(html.dark-mode) .symbol-table {
        background: #2d3748;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
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

    :global(html.dark-mode) .empty-state,
    :global(html.dark-mode) .no-symbols {
        color: #a0aec0;
    }
</style>
