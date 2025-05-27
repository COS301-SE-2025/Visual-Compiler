<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  const logoUrl = '/half_stack_phoenix_only.png';
  let theme: 'light' | 'dark' = 'light';

  function logout() {
    goto('/login');
  }

  function toggleTheme() {
    theme = theme === 'light' ? 'dark' : 'light';
    document.documentElement.setAttribute('svelvet-theme', theme);
  }

  onMount(() => {
    // Always default to light on load
    theme = 'light';
    document.documentElement.setAttribute('svelvet-theme', 'light');
  });
</script>

<header class="navbar">
  <img src={logoUrl} alt="Visual Compiler logo" class="logo" />
  <h1 class="title">Visual Compiler</h1>

  <button class="theme-toggle" on:click={toggleTheme}>
    {theme === 'light' ? '☀️ ' : '🌙'}
  </button>

  <button class="logout-btn" on:click={logout} aria-label="Logout">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
         stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
      <polyline points="16 17 21 12 16 7"/>
      <line x1="21" y1="12" x2="9" y2="12"/>
    </svg>
  </button>
</header>

<style>
  .navbar {
    position: relative;
    height: 3.5rem;
    background: #fff;
    box-shadow: 0 1px 4px rgba(0,0,0,0.1);
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .logo {
    position: absolute;
    left: 1rem;
    width: 2.5rem;
    height: auto;
  }

  .title {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
    color: #041a47;
  }

  .logout-btn {
    position: absolute;
    right: 1rem;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
    display: flex;
    align-items: center;
    color: #041a47;
  }

  .logout-btn:hover {
    color: #030f45;
  }

  .logout-btn svg {
    width: 1.5rem;
    height: 1.5rem;
  }

  .theme-toggle {
    position: absolute;
    right: 3.25rem;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    color: #041a47;
  }
</style>
