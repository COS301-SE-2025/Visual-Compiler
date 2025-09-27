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
</style>
