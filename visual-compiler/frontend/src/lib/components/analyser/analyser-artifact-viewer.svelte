<script lang="ts">

    import type {SymbolTable,Symbol} from '$lib/types';

    export let phase: string;
    export let symbol_table: Symbol[] = [];
    export let show_symbol_table = false;

</script>

<div class="artifact-viewer">
    {#if phase === 'analyser' && show_symbol_table}
        <div class="symbol-heading">
            <h3>Symbols</h3>
        </div>
        {#if symbol_table && symbol_table.length > 0}
            <table class="symbol-table">
                <thead>
                    <tr>
                        <th>Type</th>
                        <th>Name</th>
                        <th>Scope</th>
                    </tr>
                </thead>
                <tbody>
                    {#each symbol_table as symbol}
                        <tr>
                            <td>{symbol.type}</td>
                            <td>{symbol.name}</td>
                            <td>{symbol.scope}</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        {:else}
            <div class="no-symbols">No symbols generated yet</div>
        {/if}

        
    {:else if phase === 'analyser'}
        <div class="empty-state">Symbols will appear here after generation</div>
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
        width: 33%;
        padding: 0.75rem 1rem;
        border-bottom: 1px solid #e0e0e0;
        text-align: left;
        transition: border-color 0.3s ease, color 0.3s ease;
    }
    
    .symbol-table tr:last-child td {
        border-bottom: none; 
    }

    .symbol-table th {
        background: #041a47;
        color: white;
        font-weight: 500;
        transition: background-color 0.3s ease, color 0.3s ease;
    }


   
    .symbol-heading {
        color: #041a47;
        transition: color 0.3s ease;
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
        background: #1a202c;
    }

   
    :global(html.dark-mode) .empty-state,
    :global(html.dark-mode) .no-symbols {
        color: #a0aec0;
    }
</style>
