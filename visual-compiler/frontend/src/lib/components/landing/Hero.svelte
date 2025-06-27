<script>
	import { goto } from '$app/navigation';

	// launch
	// Return type: Promise<void>
	// Parameter type(s): none
	// Navigates the user to the login page.
	async function launch() {
		await goto('/auth-page');
	}
</script>

<section class="hero_section">
	<div class="logo_container">
		<img src="/half_stack_phoenix_grey.png" alt="Visual Compiler Logo" class="logo" />
		<span class="project_title">Visual Compiler</span>
	</div>

	<div class="hero_content">
		<h1 class="main_headline">Demystifying Compilers, One Block at a Time.</h1>
		<p class="sub_headline">
			An educational platform to construct, visualize, and understand the core phases of compilation
			using an intuitive, block-based interface.
		</p>
		<button class="cta_button" on:click={launch}> Launch Application </button>
	</div>
</section>

<style>
	.hero_section {
		background-color: #041a47;
		color: #fafafa;
		display: flex;
		align-items: center;
		justify-content: center;
		text-align: center;
		height: 100vh;
		padding: 0 2rem;
		position: relative;
	}

	.logo_container {
		position: absolute;
		top: 2rem;
		left: 50%;
		transform: translateX(-50%);
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.logo {
		height: 50px;
		width: auto;
	}

	.project_title {
		font-size: 1.75rem; /* Slightly adjusted size */
		font-weight: 600;
	}

	.hero_content {
		max-width: 800px;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1.5rem;
	}

	.main_headline {
		font-size: 3.5rem;
		font-weight: 700;
		line-height: 1.2;
	}

	.sub_headline {
		font-size: 1.25rem;
		max-width: 600px;
		opacity: 0.9;
	}

	.cta_button {
		background-color: #fafafa;
		color: #041a47;
		border: none;
		padding: 1rem 2rem;
		font-size: 1rem;
		font-weight: 600;
		border-radius: 8px;
		cursor: pointer;
		transition: transform 0.2s ease;
	}

	.cta_button:hover {
		transform: scale(1.05);
	}
</style>
