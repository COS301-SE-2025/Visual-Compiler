<script lang="ts">
  import { goto } from '$app/navigation';
  import { theme, ToggleTheme } from '../../stores/theme';

  const light_logo_url = '/half_stack_phoenix_only.png';
  const dark_logo_url = '/half_stack_phoenix_grey.png';

  const light_mode_icon_url = '/lightmode.png';
  const dark_mode_icon_url = '/darkmode.png';

  $: current_logo_url = $theme === 'dark' ? dark_logo_url : light_logo_url;

  // logout
  // Return type: void
  // Parameter type(s): none
  // Navigates the user to the login page.
  function logout() {
    goto('/auth-page');
  }
</script>

<header class="navbar">
  <img src={current_logo_url} alt="Visual Compiler logo" class="logo" />
  <h1 class="title">Visual Compiler</h1>

  <button class="theme-toggle" on:click={ToggleTheme} aria-label="Toggle theme">
    <img
      src={$theme === 'light' ? light_mode_icon_url : dark_mode_icon_url}
      alt={$theme === 'light' ? 'Light mode' : 'Dark mode'}
      class="theme-icon"
    />
  </button>

  <button class="logout-btn" on:click={logout} aria-label="Logout">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
      <polyline points="16 17 21 12 16 7" />
      <line x1="21" y1="12" x2="9" y2="12" />
    </svg>
  </button>
</header>

<style>
  .navbar {
    position: relative;
    height: 3.5rem;
    background: #fff;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
    display: flex;
    justify-content: center;
    align-items: center;
    transition: background-color 0.3s ease, color 0.3s ease;
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

  .theme-toggle {
    position: absolute;
    right: 4rem;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.25rem;
    line-height: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .theme-icon {
    width: 40px;
    height: 40px;
    transition: opacity 0.3s ease;
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
    right: 4rem;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1.25rem;
    color: #041a47;
    padding: 0.25rem;
    line-height: 1;
  }

 
  :global(html.dark-mode) .navbar {
    background: #041a47;
  }

  :global(html.dark-mode) .title {
    color: #d3d3d3;
  }

  :global(html.dark-mode) .logout-btn {
    color: #d3d3d3;
  }

  :global(html.dark-mode) .logout-btn:hover {
    color: #b0b0b0;
  }

  :global(html.dark-mode) .theme-toggle {
    color: #d3d3d3;
  }
</style>