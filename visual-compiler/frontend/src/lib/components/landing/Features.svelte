<script>
	// Feature data
	const FEATURES_DATA = [
		{
			icon: `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>`,
			title: 'Block-Based Construction',
			description:
				'Build compilers step-by-step using intuitive, connectable blocks that represent each phase of compilation.'
		},
		{
			icon: `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 14v8"></path><path d="M12 2v2"></path><path d="m4.93 4.93 1.41 1.41"></path><path d="M17.66 17.66l1.41 1.41"></path><path d="M2 12h2"></path><path d="M20 12h2"></path><path d="m6.34 17.66-1.41 1.41"></path><path d="m19.07 4.93-1.41 1.41"></path></svg>`,
			title: 'Interactive Learning',
			description:
				'Experiment with formal models like RegEx, Automata, and CFG with integrated tutorials to guide you.'
		},
		{
			icon: `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"></path><circle cx="12" cy="12" r="3"></circle></svg>`,
			title: 'Visualise Artefacts',
			description:
				'See the outputs of each phase like token streams, syntax trees, and symbol tables, in a clear, visual format.'
		}
	];
</script>

<section class="features_section">
	<div class="features_container">
		<h2 class="section_heading">How It Works</h2>
		<div class="features_grid">
			{#each FEATURES_DATA as feature}
				<div class="feature_card">
					<div class="feature_icon">
						{@html feature.icon}
					</div>
					<h3 class="feature_title">{feature.title}</h3>
					<p class="feature_description">{feature.description}</p>
				</div>
			{/each}
		</div>
	</div>
</section>

<style>
	.features_section {
		background-color: #fafafa;
		color: #041a47;
		padding: 6rem 2rem;
		text-align: center;
	}

	.features_container {
		max-width: 1200px;
		margin: 0 auto;
	}

	.section_heading {
		font-size: 2.5rem;
		font-weight: 700;
		margin-bottom: 4rem;
	}

	.features_grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: 2rem;
	}

	.feature_card {
		background-color: #ffffff;
		padding: 2.5rem;
		border-radius: 12px;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}

	.feature_icon {
		color: #041a47;
		margin-bottom: 0.5rem;
	}

	.feature_title {
		font-size: 1.5rem;
		font-weight: 600;
	}

	.feature_description {
		font-size: 1rem;
		line-height: 1.6;
		opacity: 0.8;
	}
</style>
