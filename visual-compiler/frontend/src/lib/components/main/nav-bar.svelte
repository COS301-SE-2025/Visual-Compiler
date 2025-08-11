<script lang="ts">
	import { goto } from '$app/navigation';
	import { theme, ToggleTheme } from '../../../lib/stores/theme';
	import HelpMenu from './help-menu.svelte';
	import AdminPanel from './admin-panel.svelte';
	import { onMount } from 'svelte';

	const light_logo_url = '/half_stack_phoenix_only.png';
	const dark_logo_url = '/half_stack_phoenix_grey.png';

	const light_mode_icon_url = '/lightmode.png';
	const dark_mode_icon_url = '/darkmode.png';

	$: current_logo_url = $theme === 'dark' ? dark_logo_url : light_logo_url;

	let show_help = false;
	let actions_container_node: HTMLElement; // A reference to the container div
	let is_admin = false;
	let show_admin_panel = false;

	// logout
	// Return type: void
	// Parameter type(s): none
	// Navigates the user to the login page.
	function logout() {
		goto('/auth-page');
	}

	// handleWindowClick
	// Return type: void
	// Parameter type(s): event (MouseEvent)
	// Closes the help menu if the click is outside its container.
	function handleWindowClick(event: MouseEvent) {
		if (
			show_help &&
			actions_container_node &&
			!actions_container_node.contains(event.target as Node)
		) {
			show_help = false;
		}
	}

	onMount(() => {
		is_admin = localStorage.getItem('is_admin') === 'true';
	});
</script>

<svelte:window on:click={handleWindowClick} />

<header class="navbar">
	<img src={current_logo_url} alt="Visual Compiler logo" class="logo" />
	<h1 class="title">Visual Compiler</h1>

	<div class="actions-container" bind:this={actions_container_node}>
		<div class="actions">
			{#if is_admin}
				<button class="action-btn admin-btn" aria-label="Admin Panel" on:click={() => show_admin_panel = !show_admin_panel}>
					<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="var(--admin-icon-color)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<path d="M12 2l3 7h7l-5.5 4.5L18 21l-6-4-6 4 1.5-7.5L2 9h7z"/>
					</svg>
				</button>
			{/if}
			<button class="action-btn help-btn" on:click={() => (show_help = !show_help)} aria-label="Help">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					><circle cx="12" cy="12" r="10"></circle><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"
					></path><line x1="12" y1="17" x2="12.01" y2="17"></line></svg
				>
			</button>

			<button class="action-btn theme-toggle" on:click={ToggleTheme} aria-label="Toggle theme">
				<img
					src={$theme === 'light' ? light_mode_icon_url : dark_mode_icon_url}
					alt={$theme === 'light' ? 'Light mode' : 'Dark mode'}
					class="theme-icon"
				/>
			</button>

			<button class="action-btn logout-btn" on:click={logout} aria-label="Logout">
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
		</div>

		{#if show_help}
			<HelpMenu on:close={() => (show_help = false)} />
		{/if}

		{#if show_admin_panel}
			<AdminPanel on:close={() => show_admin_panel = false} />
		{/if}
	</div>
</header>

<style>
	:root {
		--admin-icon-color: #041a47;
	}
	:global(html.dark-mode) {
		--admin-icon-color: #d3d3d3;
	}
	.admin-btn svg {
		stroke: var(--admin-icon-color);
		transition: stroke 0.3s;
	}

	.navbar {
		position: relative;
		height: 3.5rem;
		background: #fff;
		box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 0 1rem;
		transition:
			background-color 0.3s ease,
			color 0.3s ease;
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

	.actions-container {
		position: absolute;
		right: 1rem;
	}

	.actions {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.action-btn {
		background: none;
		border: none;
		cursor: pointer;
		padding: 0.25rem;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #041a47;
	}

	.action-btn:hover {
		color: #030f45;
	}

	.action-btn svg {
		width: 1.5rem;
		height: 1.5rem;
	}

	.theme-icon {
		width: 40px;
		height: 40px;
		transition: opacity 0.3s ease;
	}

	:global(html.dark-mode) .navbar {
		background: #041a47;
	}

	:global(html.dark-mode) .title {
		color: #d3d3d3;
	}

	:global(html.dark-mode) .action-btn {
		color: #d3d3d3;
	}

	:global(html.dark-mode) .action-btn:hover {
		color: #b0b0b0;
	}

	.admin-btn {
		margin-right: 0.05rem;
	}

	.help-btn {
		margin-left: 0.7rem;
	}
</style>
