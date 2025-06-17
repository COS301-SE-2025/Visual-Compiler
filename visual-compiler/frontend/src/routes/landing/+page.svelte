<script>
	import Features from '$lib/components/landing/Features.svelte';
  import Hero from '$lib/components/landing/Hero.svelte';
	import Walkthrough from '$lib/components/landing/Walkthrough.svelte';
</script>

<main>
  <Hero />
  <Features />
  <Walkthrough />
</main>

<style>
  :global(body) {

    margin: 0;
    font-family: 'Inter', sans-serif; 
    box-sizing: border-box;
  }

  
</style>
