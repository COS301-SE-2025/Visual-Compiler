<script>
	const CANVAS_IMAGE_URL = '/SvelvetcanvasFull.PNG';
	const CONFIGURATION_IMAGE_URL = '/NodeConfiguration.PNG';
	const SYNTAX_TREE_IMAGE_URL = '/Parsing.PNG';
</script>

<section class="walkthrough_section">
	<div class="walkthrough_container">
		<h2 class="section_heading">A Glimpse Into the Process</h2>
		<div class="step">
			<div class="step_image_container">
				<img src={CANVAS_IMAGE_URL} alt="Visual Compiler Canvas" class="step_image" />
			</div>
			<div class="step_text_container">
				<h3 class="step_title">Design Your Pipeline</h3>
				<p class="step_description">
					Start with a blank canvas. Select compiler phase blocks like Lexing, Parsing, and Code
					Generation from the toolbox and drag them into your workspace to begin mapping out your
					compiler's structure.
				</p>
			</div>
		</div>

		<div class="step reverse">
			<div class="step_image_container">
				<img src={CONFIGURATION_IMAGE_URL} alt="Node Configuration" class="step_image" />
			</div>
			<div class="step_text_container">
				<h3 class="step_title">Configure and Connect</h3>
				<p class="step_description">
					Click on each block to open its configuration panel. Define the rules for each phase using
					formal tools like Regular Expressions or Context-Free Grammars. Then, connect the blocks
					to create a data flow.
				</p>
			</div>
		</div>

		<div class="step">
			<div class="step_image_container">
				<img src={SYNTAX_TREE_IMAGE_URL} alt="Syntax Tree Visualization" class="step_image" />
			</div>
			<div class="step_text_container">
				<h3 class="step_title">Run and Visualize</h3>
				<p class="step_description">
					Execute your pipeline and watch as the application generates and visualizes the artifacts
					from each step. Instantly see the tokens, syntax trees, and final output to deepen your
					understanding.
				</p>
			</div>
		</div>
	</div>
</section>

<style>
	.walkthrough_section {
		background-color: #ffffff;
		color: #041a47;
		padding: 6rem 2rem;
	}

	.walkthrough_container {
		max-width: 1200px;
		margin: 0 auto;
		display: flex;
		flex-direction: column;
		gap: 5rem;
	}

	.section_heading {
		font-size: 2.5rem;
		font-weight: 700;
		text-align: center;
		margin-bottom: 2rem;
	}

	.step {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 4rem;
		align-items: center;
	}

	.step.reverse .step_image_container {
		order: 2;
	}

	.step.reverse .step_text_container {
		order: 1;
	}

	.step_image {
		width: 100%;
		height: auto;
		border-radius: 12px;
		box-shadow: 0 10px 20px rgba(4, 26, 71, 0.1);
	}

	.step_title {
		font-size: 2rem;
		font-weight: 600;
		margin-bottom: 1rem;
	}

	.step_description {
		font-size: 1.1rem;
		line-height: 1.7;
		opacity: 0.8;
	}

	@media (max-width: 768px) {
		.step,
		.step.reverse {
			grid-template-columns: 1fr;
			gap: 2rem;
		}

		.step.reverse .step_image_container {
			order: 1;
		}

		.step.reverse .step_text_container {
			order: 2;
		}

		.step_text_container {
			text-align: center;
		}
	}
</style>
