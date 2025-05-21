<script lang="ts">
  let active_tab: 'login' | 'register' = 'login';
  let show_password = false;
  let show_confirm_password = false;

  let reg_username = '';
  let reg_email = '';
  let reg_password = '';
  let reg_confirm_password = '';

  let login_username = '';
  let login_password = '';

  $: is_login_button_disabled = !login_username.trim() || !login_password.trim();
  $: is_register_button_disabled = !reg_username.trim() || !reg_email.trim() || 
      !reg_password.trim() || !reg_confirm_password.trim();

  const disabled_icon_path = '/disabled_state.png';
  const enabled_icon_path = '/enabled_state.png';

  function handle_register(event: Event) {
    event.preventDefault();
    if (reg_password !== reg_confirm_password) {
      alert("Passwords don't match!");
      return;
    }
    console.log('Register', { reg_username, reg_email, reg_password });
  }

  function handle_login(event: Event) {
    event.preventDefault();
    if (is_login_button_disabled) return;
    console.log('Login', { login_username, login_password });
  }

  function toggle_password_visibility() {
    show_password = !show_password;
  }

  function toggle_confirm_password_visibility() {
    show_confirm_password = !show_confirm_password;
  }
</script>

<div class="video-background">
  <video autoplay muted loop playsinline>
    <source src="/backgrounds.mp4" type="video/mp4">
  </video>

  <div class="container">
    <aside class="side-panel">
      <div class="logo">
      <img src="/half_stack.png" alt="Brand Logo" width="80" height="80" class="brand-logo"/>
      </div>

      <nav class="tab-nav">
        <button class:active={active_tab === 'login'} on:click={() => (active_tab = 'login')}>
          Login
        </button>
        <button class:active={active_tab === 'register'} on:click={() => (active_tab = 'register')}>
          Register
        </button>
      </nav>

      <section class="tab-content">
        {#if active_tab === 'login'}
          <form on:submit={handle_login}>
            <div class="input-group">
              <input
                type="text"
                id="loginUsername"
                bind:value={login_username}
                placeholder=" "
                required
              />
              <label for="loginUsername">Username</label>
              <svg class="input-icon" viewBox="0 0 24 24">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
            </div>

            <div class="input-group">
              <input
                type={show_password ? 'text' : 'password'}
                id="loginPassword"
                bind:value={login_password}
                placeholder=" "
                required
              />
              <label for="loginPassword">Password</label>
              <svg class="input-icon" viewBox="0 0 24 24">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <button type="button" class="show-password" on:click={toggle_password_visibility}>
                <svg viewBox="0 0 24 24" width="20" height="20">
                  {#if show_password}
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  {:else}
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  {/if}
                </svg>
              </button>
            </div>

            <button
              type="submit"
              class="icon-submit-btn"
              disabled={is_login_button_disabled}
              aria-label="Login"
            >
              <img
                src={is_login_button_disabled ? disabled_icon_path : enabled_icon_path}
                alt={is_login_button_disabled ? 'Login disabled' : 'Login enabled'}
                class="login-action-icon"
              />
            </button>
          </form>
        {:else}
          <form on:submit={handle_register}>
            <div class="input-group">
              <input
                type="text"
                id="regUsername"
                bind:value={reg_username}
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
                bind:value={reg_email}
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
                type={show_password ? 'text' : 'password'}
                id="regPassword"
                bind:value={reg_password}
                placeholder=" "
                required
              />
              <label for="regPassword">Password</label>
              <svg class="input-icon" viewBox="0 0 24 24">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <button type="button" class="show-password" on:click={toggle_password_visibility}>
                <svg viewBox="0 0 24 24" width="20" height="20">
                  {#if show_password}
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
                type={show_confirm_password ? 'text' : 'password'}
                id="regConfirmPassword"
                bind:value={reg_confirm_password}
                placeholder=" "
                required
              />
              <label for="regConfirmPassword">Confirm Password</label>
              <svg class="input-icon" viewBox="0 0 24 24">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <button type="button" class="show-password" on:click={toggle_confirm_password_visibility}>
                <svg viewBox="0 0 24 24" width="20" height="20">
                  {#if show_confirm_password}
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  {:else}
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  {/if}
                </svg>
              </button>
            </div>

            <button
              type="submit"
              class="icon-submit-btn"
              disabled={is_register_button_disabled}
              aria-label="Register"
            >
              <img
                src={is_register_button_disabled ? disabled_icon_path : enabled_icon_path}
                alt={is_register_button_disabled ? 'Register disabled' : 'Register enabled'}
                class="login-action-icon"
              />
            </button>
          </form>
        {/if}
      </section>
    </aside>

    <main class="main-content">
      <div class="feature-quote">
        <h1>Visual Compiler</h1>
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
    object-fit: contain;
    opacity: 0.9;
  }

  .container {
    display: flex;
    height: 100vh;
    backdrop-filter: blur(2px);
  }

  .side-panel {
    width: 340px;
    background: rgba(255, 255, 255, 0.4);
    padding: 2rem 1.75rem;
    box-shadow: 0 0 30px rgba(0, 0, 0, 0.1);
    overflow-y: auto;
    backdrop-filter: blur(8px);
  }

  .logo {
    text-align: center;
    margin-bottom: 2.5rem;
    color: #111827;
  }

  .brand-logo {
    width: 275px;
    height: 150px;
    object-fit: contain;
    transition: transform 0.2s ease;
  }

  .brand-logo:hover {
    transform: scale(1.05);
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
    border-bottom-color: #041562;
    color: #041562;
    font-weight: 600;
  }

  .tab-nav button:hover:not(.active) {
    color: #041562;
  }

  .tab-content form {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    align-items: center;
  }

  .input-group {
    position: relative;
    margin-bottom: 0.5rem;
    width: 280px;
    margin-left: auto;
    margin-right: auto;
  }

  .input-icon {
    position: absolute;
    left: 10px;
    top: 50%;
    transform: translateY(-50%);
    width: 16px;
    height: 16px;
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
    padding: 1rem 1rem 0.25rem 36px;
    font-size: 0.875rem;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    background: rgba(255, 255, 255, 0.9);
    transition: all 0.2s ease;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    height: 42px;
    box-sizing: border-box;
    position: relative;
    z-index: 1;
  }

  .input-group input:hover {
    border-color: #9ca3af;
  }

  .input-group label {
    position: absolute;
    left: 36px;
    top: 50%;
    transform: translateY(-50%);
    color: #9ca3af;
    font-size: 0.875rem;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    pointer-events: none;
    z-index: 1;
    background: transparent;
  }

  .input-group input:focus {
    border-color: #041562;
    outline: none;
    box-shadow: 0 0 0 3px rgba(4, 21, 98, 0.2);
  }

  .input-group input:focus ~ .input-icon {
    stroke: #041562;
  }

  .input-group input:focus ~ label,
  .input-group input:not(:placeholder-shown) ~ label {
    color: #041562;
    transform: translateY(-20px) scale(0.75);
    left: 38px;
    background: rgba(255, 255, 255, 0.9);
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
    color: #041562;
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
    width: 280px;
    padding: 0.75rem;
    font-size: 0.875rem;
    font-weight: 500;
    background: #041562;
    color: #fff;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    margin-top: 0.5rem;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .primary-btn:hover:not(:disabled) {
    background: #030f45;
    transform: translateY(-1px);
    box-shadow: 0 4px 6px -1px rgba(4, 21, 98, 0.2);
  }

  .primary-btn:disabled {
    background-color: #a0aec0;
    color: #e2e8f0;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
  }

  .icon-submit-btn {
    background: transparent;
    border: none;
    padding: 0;
    cursor: pointer;
    margin-top: 0.5rem;
    width: auto;
    height: auto;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
  }

  .login-action-icon {
    width: 64px;
    height: 64px;
    object-fit: contain;
    display: block;
    transition: opacity 0.2s ease, transform 0.2s ease;
  }

  .icon-submit-btn:hover:not(:disabled) .login-action-icon {
    transform: scale(1.1);
  }

  .icon-submit-btn:focus:not(:disabled) {
    outline: 2px solid #041562;
    outline-offset: 2px;
  }

  .icon-submit-btn:disabled .login-action-icon {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .icon-submit-btn:disabled {
    cursor: not-allowed;
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
      padding: 1.5rem 1rem;
    }

    .input-group,
    .primary-btn,
    .icon-submit-btn {
      max-width: 280px;
    }
     .tab-content form {
      align-items: center;
    }
  }
</style>