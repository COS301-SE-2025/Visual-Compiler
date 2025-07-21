<script lang="ts">
    import { fade } from 'svelte/transition';
   // Show/hide symbol table
   export  let showSymbolTable = true;
    
    // Mock symbol table data
    const symbolTable = [
        { type: "int", name: "count", scope: "global" },
        { type: "string", name: "message", scope: "main" },
        { type: "float", name: "price", scope: "calculate" },
        { type: "bool", name: "isValid", scope: "validate" },
        { type: "char", name: "initial", scope: "main" }
    ];

 


</script>

<div class="panel-container">
    <h2>Symbol Table</h2>
    
    {#if showSymbolTable}
        <div class="symbol-table-container" transition:fade={{ duration: 200 }}>
            <table class="symbol-table">
                <thead>
                    <tr>
                        <th>Type</th>
                        <th>Name</th>
                        <th>Scope</th>
                    </tr>
                </thead>
                <tbody>
                    {#each symbolTable as item}
                        <tr>
                            <td>{item.type}</td>
                            <td>{item.name}</td>
                            <td>{item.scope}</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

<style>
    .panel-container {
        padding: 1rem;
    }
    
    .symbol-table-container {
        margin-top: 1.5rem;
        animation: fadeIn 0.3s ease-out;
    }
    
    .symbol-table {
        width: 100%;
        border-collapse: collapse;
        border-radius: 8px;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(0, 26, 110, 0.1);
    }
    
    .symbol-table th {
        background-color: #001a6e;
        color: white;
        padding: 0.75rem 1rem;
        text-align: left;
        font-weight: 600;
        font-size: 0.9rem;
    }
    
    .symbol-table td {
        padding: 0.75rem 1rem;
        border-bottom: 1px solid #eef2f7;
        font-size: 0.9rem;
    }
    
    .symbol-table tr:last-child td {
        border-bottom: none;
    }
    
    .symbol-table tr:hover {
        background-color: #eef2f7;
    }
    
    /* Dark mode styles */
    :global(html.dark-mode) .symbol-table th {
        background-color: #1a202c;
    }
    
    :global(html.dark-mode) .symbol-table td {
        border-bottom-color: #4a5568;
        color: #e0e8f0;
        background: #2d3748;
    }
    

    
    :global(html.dark-mode) .symbol-table tr:hover {
        background-color: #2a4a8a;
    }
    
    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(10px); }
        to { opacity: 1; transform: translateY(0); }
    }

 
</style>