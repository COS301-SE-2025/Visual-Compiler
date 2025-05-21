<script lang="ts">
  let activeTab: 'login' | 'register' = 'login';

  let regUsername = '';
  let regEmail = '';
  let regPassword = '';

  let loginEmail = '';
  let loginPassword = '';

  function handleRegister(event: Event) 
  {
    event.preventDefault();
    console.log('Register', { regUsername, regEmail, regPassword });
  }

  function handleLogin(event: Event) 
  {
    event.preventDefault();
    console.log('Login', { loginEmail, loginPassword });
  }
</script>

<div class="container">
  <aside class="side-panel">
    <nav class="tab-nav">
      <button
        class:active={activeTab === 'login'}
        on:click={() => (activeTab = 'login')}
      >
        Login
      </button>
      <button
        class:active={activeTab === 'register'}
        on:click={() => (activeTab = 'register')}
      >
        Register
      </button>
    </nav>

    <section class="tab-content">
      {#if activeTab === 'login'}
        <form on:submit={handleLogin}>
          <input
            type="email"
            placeholder="Email"
            bind:value={loginEmail}
            required
          />
          <input
            type="password"
            placeholder="Password"
            bind:value={loginPassword}
            required
          />
          <button type="submit">Login</button>
        </form>
      {:else}
        <form on:submit={handleRegister}>
          <input
            type="text"
            placeholder="Username"
            bind:value={regUsername}
            required
          />
          <input
            type="email"
            placeholder="Email"
            bind:value={regEmail}
            required
          />
          <input
            type="password"
            placeholder="Password"
            bind:value={regPassword}
            required
          />
          <button type="submit">Register</button>
        </form>
      {/if}
    </section>
  </aside>

  <main class="main-content">

  </main>
</div>

<style>
  .container {
    display: flex;
    height: 100vh;
  }
  .side-panel {
    width: 25%;       
    background: #f5f5f5;
    padding: 2rem;
    box-sizing: border-box;
  }
  .tab-nav {
    display: flex;
    border-bottom: 1px solid #ccc;
    margin-bottom: 1.5rem;
  }
  .tab-nav button {
    flex: 1;
    padding: 0.75rem;
    background: transparent;
    border: none;
    border-bottom: 2px solid transparent;
    cursor: pointer;
    font-size: 1rem;
  }
  .tab-nav button.active {
    border-bottom-color: #0070f3;
    font-weight: bold;
  }
  .tab-content form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .tab-content input {
    padding: 0.75rem;
    font-size: 1rem;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  .tab-content button {
    padding: 0.75rem;
    font-size: 1rem;
    background: #0070f3;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  .main-content {
    flex: 1;
  }
</style>

