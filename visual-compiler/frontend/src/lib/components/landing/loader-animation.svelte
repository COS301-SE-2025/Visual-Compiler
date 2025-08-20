<script>
	import { fade } from 'svelte/transition';
	import { onMount, onDestroy } from 'svelte';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	let animationStarted = false;
	let animationCompleted = false;
	let animationDuration = 3500; // Total animation duration including delays

	/**
	 * Handles the animation start event
	 */
	function handleAnimationStart() {
		animationStarted = true;
		dispatch('animationstart');
	}

	/**
	 * Handles the animation end event
	 */
	function handleAnimationEnd() {
		animationCompleted = true;
		dispatch('animationend');
	}

	/**
	 * Gets the current animation state
	 */
	function getAnimationState() {
		return {
			started: animationStarted,
			completed: animationCompleted,
			duration: animationDuration
		};
	}

	/**
	 * Sets a custom animation duration
	 * @param {number} duration - Duration in milliseconds
	 */
	function setAnimationDuration(duration) {
		if (typeof duration === 'number' && duration > 0) {
			animationDuration = duration;
			return true;
		}
		return false;
	}

	/**
	 * Resets the animation state
	 */
	function resetAnimation() {
		animationStarted = false;
		animationCompleted = false;
	}

	// Lifecycle functions
	onMount(() => {
		// Start animation when component mounts
		setTimeout(handleAnimationStart, 100);
		
		// Complete animation after duration
		setTimeout(handleAnimationEnd, animationDuration);
	});

	onDestroy(() => {
		// Clean up if needed
		resetAnimation();
	});

	// Export functions for testing
	export { getAnimationState, setAnimationDuration, resetAnimation };
</script>

<div class="loader_overlay" in:fade={{ duration: 300 }} out:fade={{ duration: 500 }}>
	<div class="animation_container">
		<div class="block block_1"></div>
		<div class="block block_2"></div>
		<div class="block block_3"></div>
	</div>
</div>

<style>
	.loader_overlay {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: #041a47;
		display: flex;
		justify-content: center;
		align-items: center;
		z-index: 9999;
	}

	.animation_container {
		display: flex;
		gap: 15px;
		animation: assemble 0.5s ease-in-out forwards;
		animation-delay: 2.5s;
	}

	.block {
		width: 45px;
		height: 45px;
		background-color: #fafafa;
		border-radius: 8px;
		opacity: 0;
		animation: slide-and-pulse 2s forwards;
	}

	@keyframes assemble {
		to {
			gap: 5px;
		}
	}

	@keyframes slide-and-pulse {
		0% {
			opacity: 0;
			transform: translateY(50px);
		}
		40% {
			opacity: 1;
			transform: translateY(0);
		}
		60% {
			transform: scale(1.15);
			box-shadow: 0 0 20px 5px rgba(250, 250, 250, 0.4);
		}
		80%,
		100% {
			transform: scale(1);
			box-shadow: none;
		}
	}

	.block_1 {
		animation-delay: 0.2s;
	}
	.block_2 {
		animation-delay: 0.4s;
	}
	.block_3 {
		animation-delay: 0.6s;
	}
</style>
