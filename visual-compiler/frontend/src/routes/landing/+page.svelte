<script>
  // Import the Hero component
  import Hero from '$lib/components/Hero.svelte';
</script>

<main>
  <Hero />
  <!-- We will add more sections here later -->
</main>

<style>
  :global(body) {
    /* Removes default margin so the hero section fits perfectly */
    margin: 0;
    font-family: 'Inter', sans-serif; /* A clean, modern font */
    box-sizing: border-box;
  }

  
</style>
