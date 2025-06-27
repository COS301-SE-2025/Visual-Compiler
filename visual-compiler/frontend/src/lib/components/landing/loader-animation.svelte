<script>
	import { fade } from 'svelte/transition';
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
