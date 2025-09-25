<script lang="ts">
	import { AddToast } from '$lib/stores/toast';
	import { goto } from '$app/navigation';
	import { projectName, deleteProject } from '$lib/stores/project';
	import { resetPipeline } from '$lib/stores/pipeline';
	import { resetSourceCode } from '$lib/stores/source-code';
	import { resetLexerState } from '$lib/stores/lexer';
	import { phase_completion_status } from '$lib/stores/pipeline';

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
	$: is_register_button_disabled =
		!reg_username.trim() ||
		!reg_email.trim() ||
		!reg_password.trim() ||
		!reg_confirm_password.trim();

	const disabled_icon_path = '/disabled_state.png';
	const enabled_icon_path = '/enabled_state.png';

	// handleRegister
	// Return type: Promise<void>
	// Parameter type(s): Event
	// Handles the user registration form submission, validates input, and sends data to the API.
	async function handleRegister(event: Event) {
		event.preventDefault();

		if (reg_password !== reg_confirm_password) {
			AddToast('Passwords don\'t match. Please make sure both password fields are identical', 'error');
			return;
		}

		try {
			const response = await fetch('http://localhost:8080/api/users/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					email: reg_email,
					username: reg_username,
					password: reg_password
				})
			});

			const data = await response.json();

			if (!response.ok) {
				// Handle specific backend error messages
				if (data.error.includes("Password")) {
					AddToast(`Registration error: Password must be at least 8 characters`, 'error');
				} else if (data.error.includes("Username")) {
					AddToast(`Registration error: Username must be at least 6 characters`, 'error');
				} else if (data.error.includes("Email already exists")) {
					AddToast(`Registration error: This email is already registered. Please use a different email address`, 'error');
				} else if (data.error.includes("Username is already taken")) {
					AddToast(`Registration error: This username is already taken. Please choose a different username`, 'error');
				} else if (data.error.includes("Input is invalid")) {
					AddToast(`Registration error: Please check your input format and try again`, 'error');
				} else if (data.error.includes("Database error")) {
					AddToast(`Registration error: Database connection issue. Please try again later`, 'error');
				} else if (data.error.includes("Error in hashing password")) {
					AddToast(`Registration error: Server error processing password. Please try again`, 'error');
				} else if (data.error.includes("Error in registering user")) {
					AddToast(`Registration error: Could not complete registration. Please try again`, 'error');
				} else {
					AddToast(`Registration failed: ${data.error || 'Please check your information and try again'}`, 'error');
				}
				return;
			}

			// Success case - backend returns 201 with message "Successfully registered user"
			if (response.status === 201 && data.message === "Successfully registered user") {
				AddToast('Account successfully created! Please check your email and verify your account before logging in.', 'success');

				// Reset the form
				reg_email = '';
				reg_username = '';
				reg_password = '';
				reg_confirm_password = '';

				// Switch to login tab
				active_tab = 'login';
			}
		} catch (error) {
			AddToast(`Registration error: ${(error as Error).message}. Please check your connection and try again`, 'error');
		}
	}

	// handleLogin
	// Return type: Promise<void>
	// Parameter type(s): Event
	// Handles the user login form submission, sends credentials to the API, and navigates on success.
	async function handleLogin(event: Event) {
		event.preventDefault();
		if (is_login_button_disabled) return;

		try {
			const response = await fetch('http://localhost:8080/api/users/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					login: login_username,
					password: login_password
				})
			});

			const data = await response.json();

			if (!response.ok) {
				// Handle specific backend error messages
				if (data.error.includes("Input is invalid")) {
					AddToast('Login error: Please check your input format and try again', 'error');
				} else if (data.error.includes("Invalid credentials")) {
					AddToast('Login failed: Invalid username/email or password. Please check your credentials', 'error');
				} else if (data.error.includes("Password is incorrect")) {
					AddToast('Login failed: Incorrect password. Please try again', 'error');
				} else if (data.error.includes("Database error")) {
					AddToast('Login error: Database connection issue. Please try again later', 'error');
				} else {
					AddToast(`Login failed: ${data.error || 'Please check your username and password'}`, 'error');
				}
				return;
			}

			// Success case - backend returns 200 with message, id, is_admin, and projects
			if (response.status === 200 && data.message && data.id) {
				// Store user data in localStorage
				localStorage.setItem('user_id', data.id);
				localStorage.setItem('users_id', data.id); // Also store as users_id for optimiser compatibility
				localStorage.setItem('is_admin', data.is_admin ? 'true' : 'false');
				
				// Store the Auth0 access token in sessionStorage (as requested)
				if (data.auth_token) {
					sessionStorage.setItem('access_token', data.auth_token);
					sessionStorage.setItem('authToken', data.auth_token); // Alternative key for compatibility
				}
				
				// Store projects if available
				if (data.projects) {
					localStorage.setItem('user_projects', JSON.stringify(data.projects));
				}

				console.log('Login successful:', {
					user_id: data.id,
					is_admin: data.is_admin,
					projects: data.projects,
					has_auth_token: !!data.auth_token
				});

				// Extract username from the welcome message or use the login input
				const welcomeMessage = data.message.includes('Welcome') ? 
					data.message : 
					`Welcome back! Redirecting to your workspace...`;
				
				AddToast(welcomeMessage, 'success');

				sessionStorage.setItem('showWelcomeOverlay', 'true');

				await new Promise((res) => setTimeout(res, 2000));

				await goto('/main-workspace');
			}
		} catch (error) {
			AddToast(`Login error: ${(error as Error).message}. Please check your connection and try again`, 'error');
		}
	}

	// togglePasswordVisibility
	// Return type: void
	// Parameter type(s): none
	// Toggles the visibility of the password field.
	function togglePasswordVisibility() {
		show_password = !show_password;
	}

	// toggleConfirmPasswordVisibility
	// Return type: void
	// Parameter type(s): none
	// Toggles the visibility of the confirm password field.
	function toggleConfirmPasswordVisibility() {
		show_confirm_password = !show_confirm_password;
	}

	// handleGuestLogin
	// Return type: Promise<void>
	// Parameter type(s): none
	// Handles guest login by setting guest access token, creating a guest project, and navigating to workspace.
	async function handleGuestLogin() {
		try {
			// Clear all workspace state first to ensure clean slate for guest
			resetPipeline(); // Clear pipeline/canvas data
			resetSourceCode(); // Clear source code
			resetLexerState(); // Clear lexer state
			
			// Reset phase completion status
			phase_completion_status.set({
				source: false,
				lexer: false,
				parser: false,
				analyser: false,
				translator: false
			});

			// Clear all workspace-specific state that might persist
			// Note: These will be cleared when the workspace loads, but we clear them here too for safety
			if (typeof window !== 'undefined') {
				// Clear any workspace session data
				sessionStorage.removeItem('workspace_tokens');
				sessionStorage.removeItem('workspace_syntax_tree');
				sessionStorage.removeItem('workspace_symbol_table');
				sessionStorage.removeItem('workspace_translated_code');
			}

			// Set guest access token
			sessionStorage.setItem('access_token', 'guestuser');
			sessionStorage.setItem('authToken', 'guestuser'); // Alternative key for compatibility
			
			// Set guest user data in localStorage
			localStorage.setItem('user_id', '68d32088d29390ec2c897f35');
			localStorage.setItem('users_id', '68d32088d29390ec2c897f35'); // Also store as users_id for optimiser compatibility
			localStorage.setItem('is_admin', 'false');
			
			// Clear any existing projects for guest user
			localStorage.removeItem('user_projects');

			// Generate random 5-character string for project name
			const randomChars = Math.random().toString(36).substring(2, 7).toUpperCase();
			const guestProjectName = `Guest Project ${randomChars}`;

			// Create a guest project
			try {
				const response = await fetch('http://localhost:8080/api/users/save', {
					method: 'POST',
					headers: {
						'accept': 'application/json',
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({
						project_name: guestProjectName,
						users_id: localStorage.getItem('user_id') || '68d32088d29390ec2c897f35'
					})
				});

				if (response.ok) {
					// Store the created project name for the guest session
					localStorage.setItem('guest_project_name', guestProjectName);
					
					// Set the project name in the store so it displays in the workspace
					projectName.set(guestProjectName);
					
					console.log('Guest project created:', guestProjectName);

					// Set up cleanup for guest project on logout/close
					// setupGuestProjectCleanup(guestProjectName);
				} else {
					console.warn('Failed to create guest project, continuing with guest login');
					// Still set the project name even if creation fails
					projectName.set(guestProjectName);
				}
			} catch (projectError) {
				console.warn('Error creating guest project:', projectError);
				// Continue with guest login even if project creation fails
				// Still set the project name for display purposes
				projectName.set(guestProjectName);
			}

			console.log('Guest login initiated');

			AddToast('Welcome! You\'re now using the Visual Compiler as a guest', 'success');

			sessionStorage.setItem('showWelcomeOverlay', 'false');

			await new Promise((res) => setTimeout(res, 1500));

			await goto('/main-workspace');
		} catch (error) {
			AddToast(`Guest login error: ${(error as Error).message}. Please try again`, 'error');
		}
	}

	// deleteGuestProject
	// Return type: Promise<void>
	// Parameter type(s): string (projectName)
	// Deletes a guest project from the backend to ensure temporary projects are cleaned up
	async function deleteGuestProject(projectName: string): Promise<void> {
		try {
			console.log(`Attempting to delete guest project: ${projectName}`);
			await deleteProject(projectName, localStorage.getItem('user_id') || '68d32088d29390ec2c897f35');
			console.log(`Guest project deleted successfully: ${projectName}`);
		} catch (error) {
			console.error(`Error deleting guest project: ${projectName}`, error);
		}
	}

	// setupGuestProjectCleanup
	// Return type: void  
	// Parameter type(s): string (projectName)
	// Sets up event listeners to clean up guest projects when the user logs out or closes the application
	function setupGuestProjectCleanup(projectName: string): void {
		if (typeof window === 'undefined') return;

		// Function to handle cleanup
		const cleanup = async () => {
			const userId = localStorage.getItem('user_id');
			if (userId === (localStorage.getItem('user_id') || '68d32088d29390ec2c897f35')) {
				await deleteGuestProject(projectName);
				// Clear guest project data
				localStorage.removeItem('guest_project_name');
			}
		};

		// Handle page unload (browser close, refresh, navigate away)
		const handleBeforeUnload = (event: BeforeUnloadEvent) => {
			const userId = localStorage.getItem('user_id') || '68d32088d29390ec2c897f35';
			const guestId = '68d32088d29390ec2c897f35';
			if (userId === guestId) {
				// Use sendBeacon for more reliable cleanup on page unload
				navigator.sendBeacon('http://localhost:8080/api/users/deleteProject', JSON.stringify({
					project_name: projectName,
					users_id: userId
				}));
			}
		};

		// Handle visibility change (tab close, window minimize)
		const handleVisibilityChange = async () => {
			if (document.hidden) {
				const userId = localStorage.getItem('user_id') || '68d32088d29390ec2c897f35';
				const guestId = '68d32088d29390ec2c897f35';
				if (userId === guestId) {
					await cleanup();
				}
			}
		};

		// Add event listeners
		window.addEventListener('beforeunload', handleBeforeUnload);
		document.addEventListener('visibilitychange', handleVisibilityChange);

		// Store cleanup function reference for potential manual cleanup
		(window as any).guestProjectCleanup = cleanup;
	}
</script>

<div class="video-background">
	<video autoplay muted loop playsinline>
		<source src="/backgrounds.mp4" type="video/mp4" />
	</video>

	<div class="container">
		<aside class="side-panel">
			<div class="logo">
				<img
					src="/half_stack_blue.png"
					alt="Brand Logo"
					width="80"
					height="80"
					class="brand-logo"
				/>
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
					<form on:submit={handleLogin}>
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
								<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
								<circle cx="12" cy="7" r="4" />
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
								<rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
								<path d="M7 11V7a5 5 0 0 1 10 0v4" />
							</svg>
							<button type="button" class="show-password" on:click={togglePasswordVisibility}>
								<svg viewBox="0 0 24 24" width="20" height="20">
									{#if show_password}
										<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
										<circle cx="12" cy="12" r="3" />
									{:else}
										<path
											d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"
										/>
										<line x1="1" y1="1" x2="23" y2="23" />
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

					<!-- Guest Login Section -->
					<div class="guest-login-section">
						<div class="divider">
							<span>or</span>
						</div>
						<button type="button" class="guest-login-btn" on:click={handleGuestLogin} aria-label="Continue as Guest">
							<svg class="guest-icon" viewBox="0 0 24 24" width="20" height="20">
								<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
								<circle cx="12" cy="7" r="4" />
								<path d="M12 14v7" opacity="0.5" />
								<path d="M8 18h8" opacity="0.5" />
							</svg>
							<span class="guest-text">
								<span class="guest-title">Continue as Guest</span>
								<span class="guest-subtitle">Try the Visual Compiler without an account</span>
							</span>
						</button>
					</div>
				{:else}
					<form on:submit={handleRegister}>
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
								<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
								<circle cx="12" cy="7" r="4" />
							</svg>
						</div>

						<div class="input-group">
							<input type="email" id="regEmail" bind:value={reg_email} placeholder=" " required />
							<label for="regEmail">Email</label>
							<svg class="input-icon" viewBox="0 0 24 24">
								<path
									d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"
								/>
								<polyline points="22,6 12,13 2,6" />
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
								<rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
								<path d="M7 11V7a5 5 0 0 1 10 0v4" />
							</svg>
							<button type="button" class="show-password" on:click={togglePasswordVisibility}>
								<svg viewBox="0 0 24 24" width="20" height="20">
									{#if show_password}
										<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
										<circle cx="12" cy="12" r="3" />
									{:else}
										<path
											d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"
										/>
										<line x1="1" y1="1" x2="23" y2="23" />
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
								<rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
								<path d="M7 11V7a5 5 0 0 1 10 0v4" />
							</svg>
							<button
								type="button"
								class="show-password"
								on:click={toggleConfirmPasswordVisibility}
							>
								<svg viewBox="0 0 24 24" width="20" height="20">
									{#if show_confirm_password}
										<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
										<circle cx="12" cy="12" r="3" />
									{:else}
										<path
											d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"
										/>
										<line x1="1" y1="1" x2="23" y2="23" />
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
				<h1 class="app-title">Visual Compiler</h1>
			</div>
		</main>
	</div>
</div>

<style>
	:global(body) {
		margin: 0;
		font-family:
			'Inter',
			-apple-system,
			BlinkMacSystemFont,
			sans-serif;
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
		border-bottom-color: #041a47;
		color: #041a47;
		font-weight: 600;
	}

	.tab-nav button:hover:not(.active) {
		color: #041a47;
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
		border-color: #041a47;
		outline: none;
		box-shadow: 0 0 0 3px rgba(4, 21, 98, 0.2);
	}

	.input-group input:focus ~ .input-icon {
		stroke: #041a47;
	}

	.input-group input:focus ~ label,
	.input-group input:not(:placeholder-shown) ~ label {
		color: #041a47;
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
		color: #041a47;
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
		transition:
			opacity 0.2s ease,
			transform 0.2s ease;
	}

	.icon-submit-btn:hover:not(:disabled) .login-action-icon {
		transform: scale(1.1);
	}

	.icon-submit-btn:disabled .login-action-icon {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.icon-submit-btn:disabled {
		cursor: not-allowed;
	}

	/* Guest Login Styles */
	.guest-login-section {
		width: 280px;
		margin: 1.5rem auto 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}

	.divider {
		position: relative;
		width: 100%;
		text-align: center;
		margin: 0.5rem 0;
	}

	.divider::before {
		content: '';
		position: absolute;
		top: 50%;
		left: 0;
		right: 0;
		height: 1px;
		background: linear-gradient(to right, transparent, #e5e7eb 20%, #e5e7eb 80%, transparent);
		z-index: 1;
	}

	.divider span {
		position: relative;
		background: rgba(255, 255, 255, 0.95);
		padding: 0 1rem;
		font-size: 0.875rem;
		color: #6b7280;
		font-weight: 500;
		z-index: 2;
		border-radius: 4px;
	}

	.guest-login-btn {
		width: 100%;
		padding: 0.875rem 1rem;
		background: rgba(255, 255, 255, 0.9);
		border: 2px solid #e5e7eb;
		border-radius: 8px;
		cursor: pointer;
		display: flex;
		align-items: center;
		gap: 0.75rem;
		transition: all 0.2s ease;
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
	}

	.guest-login-btn:hover {
		border-color: #041a47;
		box-shadow: 0 2px 4px rgba(4, 26, 71, 0.1);
		transform: translateY(-1px);
	}

	.guest-icon {
		flex-shrink: 0;
		stroke: #6b7280;
		stroke-width: 2;
		stroke-linecap: round;
		stroke-linejoin: round;
		fill: none;
		transition: stroke 0.2s ease;
	}

	.guest-login-btn:hover .guest-icon {
		stroke: #041a47;
	}

	.guest-text {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		flex: 1;
	}

	.guest-title {
		font-size: 0.9rem;
		font-weight: 600;
		color: #374151;
		line-height: 1.2;
		transition: color 0.2s ease;
	}

	.guest-subtitle {
		font-size: 0.75rem;
		color: #6b7280;
		line-height: 1.3;
		margin-top: 0.125rem;
		transition: color 0.2s ease;
	}

	.guest-login-btn:hover .guest-title {
		color: #041a47;
	}

	.guest-login-btn:hover .guest-subtitle {
		color: #374151;
	}

	.main-content {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 2rem;
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
		.icon-submit-btn {
			max-width: 280px;
		}
		.tab-content form {
			align-items: center;
		}
	}

	.app-title {
		font-size: 3rem;
		font-weight: 900;
		margin: 0 auto;
		line-height: 1.1;
		color: #1d3159;
	}
</style>