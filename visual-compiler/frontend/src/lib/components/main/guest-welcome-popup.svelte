<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let show = false;
	
	const dispatch = createEventDispatcher();

	function handleClose() {
		show = false;
		dispatch('close');
	}

	function handleGetStarted() {
		show = false;
		dispatch('close');
	}
</script>

{#if show}
	<div class="popup-overlay" on:click={handleClose} role="dialog" aria-modal="true" aria-labelledby="guest-welcome-title">
		<div class="popup-content" on:click|stopPropagation role="document">
			<!-- Guest Icon -->
			<div class="icon-container">
				<svg class="guest-icon" viewBox="0 0 24 24" width="48" height="48">
					<circle cx="12" cy="7" r="4" stroke="currentColor" stroke-width="2" fill="none"/>
					<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" stroke="currentColor" stroke-width="2" fill="none"/>
				</svg>
			</div>

			<!-- Title -->
			<h2 id="guest-welcome-title">Welcome, Guest User!</h2>
			
			<!-- Subtitle -->
			<p class="subtitle">You're now exploring the Visual Compiler in guest mode</p>

			<!-- Features List -->
			<div class="features-section">
				<h3>What you can do:</h3>
				<ul class="features-list">
					<li>
						<svg class="check-icon" viewBox="0 0 24 24" width="16" height="16">
							<path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
						</svg>
						Create and test compiler phases
					</li>
					<li>
						<svg class="check-icon" viewBox="0 0 24 24" width="16" height="16">
							<path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
						</svg>
						Use the AI assistant for help
					</li>
					<li>
						<svg class="check-icon" viewBox="0 0 24 24" width="16" height="16">
							<path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
						</svg>
						Explore all compiler features
					</li>
				</ul>
			</div>

			<!-- Limitations Section -->
			<div class="limitations-section">
				<h3>Guest limitations:</h3>
				<ul class="limitations-list">
					<li>
						<svg class="warning-icon" viewBox="0 0 24 24" width="16" height="16">
							<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
							<line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
							<line x1="12" y1="16" x2="12.01" y2="16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
						</svg>
						Your work is temporary and won't be saved
					</li>
					<li>
						<svg class="warning-icon" viewBox="0 0 24 24" width="16" height="16">
							<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
							<line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
							<line x1="12" y1="16" x2="12.01" y2="16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
						</svg>
						Cannot save or load projects
					</li>
					<li>
						<svg class="warning-icon" viewBox="0 0 24 24" width="16" height="16">
							<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
							<line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
							<line x1="12" y1="16" x2="12.01" y2="16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
						</svg>
						Data clears when you close the browser
					</li>
				</ul>
			</div>

			<!-- Action Buttons -->
			<div class="actions">
				<button class="get-started-btn" on:click={handleGetStarted}>
					<svg viewBox="0 0 24 24" width="18" height="18">
						<path d="M12 2l3.09 6.26L22 9l-5 4.87L18.18 21 12 17.77 5.82 21 7 14.87 2 9l6.91-1.74L12 2z" fill="currentColor"/>
					</svg>
					Get Started
				</button>
			</div>

			<!-- Register Prompt -->
			<div class="register-prompt">
				<p>Want to save your work? <a href="/auth-page" class="register-link">Create an account</a></p>
			</div>
		</div>
	</div>
{/if}

<style>
	.popup-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.6);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 3000;
		animation: fadeIn 0.3s ease-out;
		padding: 1rem;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.popup-content {
		background: white;
		border-radius: 16px;
		padding: 2rem;
		max-width: 480px;
		width: 100%;
		max-height: 90vh;
		overflow-y: auto;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
		animation: slideUp 0.3s ease-out;
		position: relative;
	}

	@keyframes slideUp {
		from {
			opacity: 0;
			transform: translateY(20px) scale(0.95);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}

	.icon-container {
		text-align: center;
		margin-bottom: 1.5rem;
	}

	.guest-icon {
		stroke: #041a47;
		stroke-width: 2;
		stroke-linecap: round;
		stroke-linejoin: round;
		fill: none;
		background: linear-gradient(135deg, #E8F4FD, #E0F2FE);
		border-radius: 50%;
		padding: 12px;
		width: 72px;
		height: 72px;
	}

	h2 {
		text-align: center;
		color: #1f2937;
		font-size: 1.5rem;
		font-weight: 700;
		margin: 0 0 0.5rem 0;
		letter-spacing: -0.025em;
	}

	.subtitle {
		text-align: center;
		color: #6b7280;
		font-size: 1rem;
		margin: 0 0 2rem 0;
		line-height: 1.5;
	}

	.features-section, .limitations-section {
		margin-bottom: 1.5rem;
	}

	.features-section h3, .limitations-section h3 {
		color: #374151;
		font-size: 1rem;
		font-weight: 600;
		margin: 0 0 0.75rem 0;
	}

	.features-list, .limitations-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.features-list li, .limitations-list li {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		font-size: 0.9rem;
		line-height: 1.4;
	}

	.features-list li {
		color: #041a47;
	}

	.limitations-list li {
		color: #0369a1;
	}

	.check-icon {
		color: #041a47;
		flex-shrink: 0;
	}

	.warning-icon {
		color: #0369a1;
		flex-shrink: 0;
	}

	.actions {
		display: flex;
		gap: 0.75rem;
		margin: 2rem 0 1rem 0;
	}

	.get-started-btn {
		flex: 1;
		background: linear-gradient(135deg, #041a47, #1e40af);
		color: white;
		border: none;
		padding: 0.75rem 1.5rem;
		border-radius: 8px;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		transition: all 0.2s ease;
		box-shadow: 0 2px 4px rgba(4, 26, 71, 0.2);
	}

	.get-started-btn:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(4, 26, 71, 0.3);
	}

	.close-btn {
		background: transparent;
		color: #6b7280;
		border: 2px solid #e5e7eb;
		padding: 0.75rem 1.5rem;
		border-radius: 8px;
		font-size: 0.9rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.close-btn:hover {
		color: #374151;
		border-color: #d1d5db;
		background: #f9fafb;
	}

	.register-prompt {
		text-align: center;
		margin-top: 1rem;
		padding-top: 1rem;
		border-top: 1px solid #e5e7eb;
	}

	.register-prompt p {
		color: #6b7280;
		font-size: 0.85rem;
		margin: 0;
	}

	.register-link {
		color: #041a47;
		text-decoration: none;
		font-weight: 600;
		transition: color 0.2s ease;
	}

	.register-link:hover {
		color: #1e3a8a;
		text-decoration: underline;
	}

	/* Dark mode styles */
	:global(html.dark-mode) .popup-overlay {
		background: rgba(0, 0, 0, 0.8);
	}

	:global(html.dark-mode) .popup-content {
		background: #1f2937;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
	}

	:global(html.dark-mode) h2 {
		color: #f9fafb;
	}

	:global(html.dark-mode) .subtitle {
		color: #d1d5db;
	}

	:global(html.dark-mode) .features-section h3,
	:global(html.dark-mode) .limitations-section h3 {
		color: #e5e7eb;
	}

	:global(html.dark-mode) .features-list li {
		color: #3b82f6;
	}

	:global(html.dark-mode) .limitations-list li {
		color: #60a5fa;
	}

	:global(html.dark-mode) .check-icon {
		color: #3b82f6;
	}

	:global(html.dark-mode) .warning-icon {
		color: #60a5fa;
	}

	:global(html.dark-mode) .close-btn {
		color: #d1d5db;
		border-color: #374151;
	}

	:global(html.dark-mode) .close-btn:hover {
		color: #f3f4f6;
		border-color: #4b5563;
		background: #374151;
	}

	:global(html.dark-mode) .register-prompt {
		border-top-color: #374151;
	}

	:global(html.dark-mode) .register-prompt p {
		color: #9ca3af;
	}

	:global(html.dark-mode) .register-link {
		color: #60a5fa;
	}

	:global(html.dark-mode) .register-link:hover {
		color: #93c5fd;
	}

	/* Mobile responsiveness */
	@media (max-width: 640px) {
		.popup-content {
			padding: 1.5rem;
			margin: 1rem;
		}

		.actions {
			flex-direction: column;
		}

		.get-started-btn, .close-btn {
			width: 100%;
		}

		h2 {
			font-size: 1.25rem;
		}

		.subtitle {
			font-size: 0.9rem;
		}
	}
</style>