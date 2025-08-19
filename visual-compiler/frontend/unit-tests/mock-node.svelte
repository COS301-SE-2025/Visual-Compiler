<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	export let label = '';
	const dispatch = createEventDispatcher();

	function handleClick() {
		dispatch('nodeClicked');
	}
</script>

<button on:click={handleClick}>
	{label}
</button>
