<script>
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';

	import Loader from '$lib/components/landing/loader-animation.svelte';
	import Hero from '$lib/components/landing/hero.svelte';
	import Features from '$lib/components/landing/features.svelte';
	import Walkthrough from '$lib/components/landing/walkthrough.svelte';
	import InfoSection from '$lib/components/landing/info-section.svelte';

	let is_loading = true;

	onMount(() => {
		setTimeout(() => {
			is_loading = false;
		}, 3000);
	});
</script>

{#if is_loading}
	<Loader />
{:else}
	<main in:fade={{ duration: 500 }}>
		<Hero />
		<Features />
		<Walkthrough />
		<InfoSection />
	</main>
{/if}

<style>
	:global(body) {
		margin: 0;
		font-family: 'Inter', sans-serif;
		box-sizing: border-box;
		overflow: hidden;
	}

	:global(body):has(main) {
		overflow: auto;
	}
</style>
