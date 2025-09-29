<script>
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let showScrollIndicator = true;

	// launch
	// Return type: Promise<void>
	// Parameter type(s): none
	// Navigates the user to the login page.
	async function launch() {
		await goto('/auth-page');
	}

	// Handle scroll to hide indicator when user starts scrolling
	function handleScroll() {
		if (window.scrollY > 50) {
			showScrollIndicator = false;
		}
	}

	// Smooth scroll to next section
	function scrollToNext() {
		const heroHeight = window.innerHeight;
		window.scrollTo({
			top: heroHeight,
			behavior: 'smooth'
		});
	}

	onMount(() => {
		window.addEventListener('scroll', handleScroll);
		return () => {
			window.removeEventListener('scroll', handleScroll);
		};
	});
</script>

<section class="hero_section">
	<div class="logo_container">
		<div class="logo">
			<img src="/half_stack_phoenix_grey.png" alt="Visual Compiler Logo"/>
		</div>
		<div class="project_title">
			<span>Visual Compiler</span>
		</div>
	</div>

	<div class="hero_content">
		<h1 class="main_headline">Demystifying Compilers,<br>One Block at a Time.</h1>
		<p class="sub_headline">
			An educational platform to construct, configure, visualise<br>and understand the core phases of compilation.
		</p>
		<button class="cta_button" on:click={launch}> Get Started </button>
	</div>

	<!-- Scroll Indicator -->
	{#if showScrollIndicator}
		<div class="scroll_indicator" on:click={scrollToNext} role="button" tabindex="0" on:keydown={(e) => e.key === 'Enter' && scrollToNext()}>
			<div class="scroll_text">Scroll to explore</div>
			<div class="scroll_arrow">
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<path d="m6 9 6 6 6-6"/>
				</svg>
			</div>
			<div class="scroll_dots">
				<div class="dot"></div>
				<div class="dot"></div>
				<div class="dot"></div>
			</div>
		</div>
	{/if}
</section>

<style>
	.hero_section {
		background-color: #041a47;
		color: #fafafa;
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		min-height: 100vh;
		padding: 2rem;
		gap: 3rem;
	}

	.logo_container {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}

	.logo img {
		height: 125px;
		width: auto;
	}

	.project_title span {
		font-size: 2rem;
		font-weight: 600;
	}

	.hero_content {
		max-width: 800px;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		justify-content: center;
	}

	.main_headline {
		font-size: 3rem;
		font-weight: 700;
		line-height: 1.2;
	}

	.sub_headline {
		font-size: 1.5rem;
		max-width: 700px;
		opacity: 0.9;
	}

	.cta_button {
		margin-top: 2rem;
		background-color: #fafafa;
		color: #041a47;
		border: none;
		padding: 1rem 2rem;
		font-size: 1.25rem;
		font-weight: 600;
		border-radius: 8px;
		cursor: pointer;
		transition: transform 0.2s ease;
	}

	.cta_button:hover {
		transform: scale(1.05);
	}

	/* Scroll Indicator Styles */
	.scroll_indicator {
		position: absolute;
		bottom: 2rem;
		left: 50%;
		transform: translateX(-50%);
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
		opacity: 0.8;
		transition: opacity 0.3s ease;
		animation: fadeInUp 1s ease-out 0.5s both, bounce 2s ease-in-out 2s infinite;
	}

	.scroll_indicator:hover {
		opacity: 1;
	}

	.scroll_text {
		font-size: 0.875rem;
		font-weight: 500;
		color: #fafafa;
		margin-bottom: 0.25rem;
		letter-spacing: 0.5px;
	}

	.scroll_arrow {
		color: #fafafa;
		animation: arrowBounce 1.5s ease-in-out infinite;
		margin-bottom: 0.5rem;
	}

	.scroll_dots {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		align-items: center;
	}

	.dot {
		width: 4px;
		height: 4px;
		background-color: #fafafa;
		border-radius: 50%;
		opacity: 0.6;
		animation: dotFade 2s ease-in-out infinite;
	}

	.dot:nth-child(1) {
		animation-delay: 0s;
	}

	.dot:nth-child(2) {
		animation-delay: 0.3s;
	}

	.dot:nth-child(3) {
		animation-delay: 0.6s;
	}

	/* Keyframe Animations */
	@keyframes fadeInUp {
		from {
			opacity: 0;
			transform: translateX(-50%) translateY(20px);
		}
		to {
			opacity: 0.8;
			transform: translateX(-50%) translateY(0);
		}
	}

	@keyframes bounce {
		0%, 20%, 50%, 80%, 100% {
			transform: translateX(-50%) translateY(0);
		}
		40% {
			transform: translateX(-50%) translateY(-8px);
		}
		60% {
			transform: translateX(-50%) translateY(-4px);
		}
	}

	@keyframes arrowBounce {
		0%, 100% {
			transform: translateY(0);
		}
		50% {
			transform: translateY(8px);
		}
	}

	@keyframes dotFade {
		0%, 100% {
			opacity: 0.3;
		}
		50% {
			opacity: 1;
		}
	}

	/* Responsive Design */
	@media (max-width: 768px) {
		.scroll_indicator {
			bottom: 1.5rem;
		}
		
		.scroll_text {
			font-size: 0.75rem;
		}
		
		.main_headline {
			font-size: 2.5rem;
		}
		
		.sub_headline {
			font-size: 1.1rem;
		}
	}
</style>
