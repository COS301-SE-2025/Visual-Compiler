<script lang="ts">
  let activeTab: 'login' | 'register' = 'login';
  let showPassword = false;
  let showConfirmPassword = false;

  // Registration fields
  let regUsername = '';
  let regEmail = '';
  let regPassword = '';
  let regConfirmPassword = '';

  // Login fields
  let loginEmail = '';
  let loginPassword = '';

  function handleRegister(event: Event) {
    event.preventDefault();
    if (regPassword !== regConfirmPassword) {
      alert("Passwords don't match!");
      return;
    }
    console.log('Register', { regUsername, regEmail, regPassword });
  }

  function handleLogin(event: Event) {
    event.preventDefault();
    console.log('Login', { loginEmail, loginPassword });
  }

  function togglePasswordVisibility() {
    showPassword = !showPassword;
  }

  function toggleConfirmPasswordVisibility() {
    showConfirmPassword = !showConfirmPassword;
  }
</script>

<div class="video-background">
  <video autoplay muted loop playsinline>
    <source src="/videos/auth-bg.mp4" type="video/mp4">
    <img src="/images/auth-bg-fallback.jpg" alt="Background">
  </video>

  <div class="container">
    <aside class="side-panel">
      <div class="logo">
        <svg viewBox="0 0 24 24" width="40" height="40">
          <path fill="#3b82f6" d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
        </svg>
        <h2>Your Brand</h2>
      </div>

      <nav class="tab-nav">
        <button class:active={activeTab === 'login'} on:click={() => (activeTab = 'login')}>
          Login
        </button>
        <button class:active={activeTab === 'register'} on:click={() => (activeTab = 'register')}>
          Register
        </button>
      </nav>

      <section class="tab-content">
        {#if activeTab === 'login'}
          <form on:submit={handleLogin}>
            <div class="input-group">
              <input
                type="email"
                id="loginEmail"
                bind:value={loginEmail}
                placeholder=" "
                required
              />
              <label for="loginEmail">Email</label>
              <svg class="input-icon" viewBox="0 0 24 24">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                <polyline points="22,6 12,13 2,6"/>
              </svg>
            </div>

            <div class="input-group">
              <input
                type={showPassword ? 'text' : 'password'}
                id="loginPassword"
                bind:value={loginPassword}
                placeholder=" "
                required
              />
              <label for="loginPassword">Password</label>
              <svg class="input-icon" viewBox="0 0 24 24">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <button type="button" class="show-password" on:click={togglePasswordVisibility}>
                <svg viewBox="0 0 24 24" width="20" height="20">
                  {#if showPassword}
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  {:else}
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  {/if}
                </svg>
              </button>
            </div>

            <button type="submit" class="primary-btn">Login</button>
          </form>
        {:else}
          <form on:submit={handleRegister}>
            <div class="input-group">
              <input
                type="text"
                id="regUsername"
                bind:value={regUsername}
                placeholder=" "
                required
              />
              <label for="regUsername">Username</label>
              <svg class="input-icon" viewBox="0 0 24 24">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
            </div>

            <div class="input-group">
              <input
                type="email"
                id="regEmail"
                bind:value={regEmail}
                placeholder=" "
                required
              />
              <label for="regEmail">Email</label>
              <svg class="input-icon" viewBox="0 0 24 24">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                <polyline points="22,6 12,13 2,6"/>
              </svg>
            </div>

            <div class="input-group">
              <input
                type={showPassword ? 'text' : 'password'}
                id="regPassword"
                bind:value={regPassword}
                placeholder=" "
                required
              />
              <label for="regPassword">Password</label>
              <svg class="input-icon" viewBox="0 0 24 24">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <button type="button" class="show-password" on:click={togglePasswordVisibility}>
                <svg viewBox="0 0 24 24" width="20" height="20">
                  {#if showPassword}
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  {:else}
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  {/if}
                </svg>
              </button>
              <div class="input-hint">At least 8 characters</div>
            </div>

            <div class="input-group">
              <input
                type={showConfirmPassword ? 'text' : 'password'}
                id="regConfirmPassword"
                bind:value={regConfirmPassword}
                placeholder=" "
                required
              />
              <label for="regConfirmPassword">Confirm Password</label>
              <svg class="input-icon" viewBox="0 0 24 24">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <button type="button" class="show-password" on:click={toggleConfirmPasswordVisibility}>
                <svg viewBox="0 0 24 24" width="20" height="20">
                  {#if showConfirmPassword}
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  {:else}
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  {/if}
                </svg>
              </button>
            </div>

            <button type="submit" class="primary-btn">Create Account</button>
          </form>
        {/if}
      </section>
    </aside>

    <main class="main-content">
      <div class="feature-quote">
        <h1>Welcome to Our Platform</h1>
        <p>"The best way to predict the future is to create it."</p>
        <p>- Peter Drucker</p>
      </div>
    </main>
  </div>
</div>

<style>
  :global(body) {
    margin: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    color: #111827;
  }

  .video-background {
    position: relative;
    width: 100%;
    height: 100vh;
    overflow: hidden;
    color: #fff;
  }

  .video-background video {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    min-width: 100%;
    min-height: 100%;
    z-index: -1;
    object-fit: cover;
    opacity: 0.9;
  }

  .container {
    display: flex;
    height: 100vh;
    backdrop-filter: blur(2px);
  }

  .side-panel {
    width: 450px;
    background: rgba(255, 255, 255, 0.97);
    padding: 3rem 2.5rem;
    box-shadow: 0 0 30px rgba(0, 0, 0, 0.1);
    overflow-y: auto;
  }

  .logo {
    text-align: center;
    margin-bottom: 2.5rem;
    color: #111827;
  }

  .logo h2 {
    margin-top: 0.5rem;
    font-size: 1.5rem;
    font-weight: 600;
  }

  .tab-nav {
    display: flex;
    border-bottom: 1px solid #e5e7eb;
    margin-bottom: 2rem;
  }

  .tab-nav button {
    flex: 1;
    padding: 1rem;
    background: transparent;
    border: none;
    border-bottom: 3px solid transparent;
    cursor: pointer;
    font-size: 1rem;
    font-weight: 500;
    color: #6b7280;
    transition: all 0.2s ease;
  }

  .tab-nav button.active {
    border-bottom-color: #3b82f6;
    color: #3b82f6;
    font-weight: 600;
  }

  .tab-nav button:hover {
    color: #3b82f6;
  }

  .tab-content form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .input-group {
    position: relative;
    margin-bottom: 0.5rem;
  }

  .input-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    width: 18px;
    height: 18px;
    stroke: #9ca3af;
    stroke-width: 2;
    stroke-linecap: round;
    stroke-linejoin: round;
    fill: none;
    transition: all 0.2s ease;
    z-index: 2;
  }

  .input-group input {
    width: 100%;
    padding: 1.5rem 1rem 0.5rem 40px;
    font-size: 1rem;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    background: #fff;
    transition: all 0.2s ease;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    height: 56px;
    box-sizing: border-box;
    position: relative;
    z-index: 1;
  }

  .input-group input:hover {
    border-color: #9ca3af;
  }

  .input-group label {
    position: absolute;
    left: 40px;
    top: 50%;
    transform: translateY(-50%);
    color: #9ca3af;
    font-size: 1rem;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    pointer-events: none;
    z-index: 1;
    background: transparent;
  }

  .input-group input:focus {
    border-color: #3b82f6;
    outline: none;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
  }

  .input-group input:focus ~ .input-icon {
    stroke: #3b82f6;
  }

  .input-group input:focus ~ label,
  .input-group input:not(:placeholder-shown) ~ label {
    transform: translateY(-20px) scale(0.75); /* Increased from -12px to -20px and scale from 0.85 to 0.75 */
    left: 42px;
    color: #3b82f6;
    background: #fff;
    padding: 0 4px;
    z-index: 3;
  }

  .input-group input::placeholder {
    color: transparent;
    transition: all 0.2s ease;
  }

  .input-group input:focus::placeholder {
    color: #9ca3af;
  }

  .show-password {
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
    z-index: 2;
    color: #9ca3af;
  }

  .show-password:hover {
    color: #3b82f6;
  }

  .show-password svg {
    stroke: currentColor;
    stroke-width: 2;
    stroke-linecap: round;
    stroke-linejoin: round;
    fill: none;
  }

  .input-hint {
    position: absolute;
    bottom: -18px;
    left: 40px;
    font-size: 0.75rem;
    color: #6b7280;
    opacity: 0;
    transform: translateY(5px);
    transition: all 0.2s ease;
  }

  .input-group:hover .input-hint {
    opacity: 1;
    transform: translateY(0);
  }

  .primary-btn {
    padding: 0.875rem;
    font-size: 1rem;
    font-weight: 500;
    background: #3b82f6;
    color: #fff;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    margin-top: 0.5rem;
  }

  .primary-btn:hover {
    background: #2563eb;
    transform: translateY(-1px);
    box-shadow: 0 4px 6px -1px rgba(59, 130, 246, 0.2);
  }

  .main-content {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
  }

  .feature-quote {
    max-width: 600px;
    text-align: center;
    color: white;
    text-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
  }

  .feature-quote h1 {
    font-size: 2.5rem;
    margin-bottom: 1.5rem;
    font-weight: 700;
    line-height: 1.2;
  }

  .feature-quote p {
    font-size: 1.25rem;
    margin: 0.5rem 0;
    opacity: 0.9;
    line-height: 1.6;
  }

  .feature-quote p:last-child {
    font-style: italic;
    opacity: 0.8;
  }

  @media (max-width: 768px) {
    .container {
      flex-direction: column;
    }

    .side-panel {
      width: 100%;
      padding: 2rem 1.5rem;
    }

    .main-content {
      display: none;
    }

    .input-group input {
      padding: 1.25rem 1rem 0.5rem 40px;
    }

    .input-group input:focus ~ label,
    .input-group input:not(:placeholder-shown) ~ label {
      transform: translateY(-16px) scale(0.75); /* Adjusted for mobile */
    }
  }
</style>