<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fly } from 'svelte/transition';

	/**
	 * Controls the visibility of the overlay.
	 * Bind this prop in the parent component.
	 * e.g., <WelcomeOverlay bind:show={showModal} />
	 */
	export let show = false;



	
	
</script>


{#if show}

	<div
		class="backdrop"
		transition:fly={{ y: -100, duration: 400, opacity: 0.5 }}
		on:click={handleClose}
	>

		<div class="modal" on:click|stopPropagation>
			<h1 class="title">Hi user123</h1>
			
		</div>
	</div>
{/if}

<style>

	.backdrop {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: rgba(0, 0, 0, 0.5);
		-webkit-backdrop-filter: blur(4px); 
		backdrop-filter: blur(4px);
		z-index: 9999; 

		
		display: flex;
		justify-content: center;
		align-items: center;
	}


	.modal {
		width: 50vw;
		height: 50vh;
		min-width: 320px; 
		max-height: 40rem;
		background: white;
		border-radius: 16px;
		padding: 2rem 3rem;
		box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);

		
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		text-align: center;
	}

	.title {
		font-size: 2.5rem;
		font-weight: 700;
		color: #1a202c;
		margin-bottom: 1rem;
	}

	.text-content {
		font-size: 1.1rem;
		line-height: 1.6;
		color: #4a5568;
		margin-bottom: 2rem;
		max-width: 90%;
	}

	.action-button {
		padding: 0.8rem 1.6rem;
		font-size: 1rem;
		font-weight: 600;
		color: white;
		background-color: #3b82f6; /* A nice, modern blue */
		border: none;
		border-radius: 8px;
		cursor: pointer;
		transition:
			background-color 0.2s ease-in-out,
			transform 0.1s ease;
	}

	.action-button:hover {
		background-color: #2563eb;
	}

	.action-button:active {
		transform: scale(0.98);
	}

	/* Dark Mode Styles: Adapts the component to the existing dark theme in your project */
	:global(html.dark) .modal {
		background: #2d3748; /* Dark gray */
	}

	:global(html.dark) .title {
		color: #edf2f7; /* Light gray text */
	}

	:global(html.dark) .text-content {
		color: #a0aec0; /* Lighter gray for body */
	}

	:global(html.dark) .action-button {
		background-color: #4299e1; /* A brighter blue for dark mode */
	}

	:global(html.dark) .action-button:hover {
		background-color: #3182ce;
	}
</style>