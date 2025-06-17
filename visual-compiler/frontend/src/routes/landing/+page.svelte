<script>
  import { onMount } from 'svelte';
  import { fade } from 'svelte/transition';
  import Loader from '$lib/components/landing/Loader.svelte';
  import Hero from '$lib/components/landing/Hero.svelte';
  import Features from '$lib/components/landing/Features.svelte';
  import Walkthrough from '$lib/components/landing/Walkthrough.svelte';

  let is_loading = true;

  onMount(() => {
    setTimeout(() => {
      is_loading = false;
    }, 3200); 
  });
</script>


{#if is_loading}
  <Loader />
{:else}

  <main in:fade={{ duration: 500 }}>
    <Hero />
    <Features />
    <Walkthrough />
  </main>
{/if}


<style>
  :global(body) {
    margin: 0;
    font-family: 'Inter', sans-serif; 
    box-sizing: border-box;
    /* This prevents a scrollbar from flashing during the animation */
    overflow: hidden; 
  }

  /* When the loader is gone, restore scrolling */
  :global(body):has(main) {
    overflow: auto;
  }
</style>
