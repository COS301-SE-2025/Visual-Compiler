<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fly, fade } from 'svelte/transition';
	import { quintOut } from 'svelte/easing';

	export let show: boolean = false;

	const dispatch = createEventDispatcher();

	interface TutorialStep {
		id: number;
		title: string;
		description: string;
		icon: string;
		details: string[];
	}

	const tutorialSteps: TutorialStep[] = [
		{
			id: 1,
			title: "Welcome to Visual Compiler!",
			description: "Learn how to build your first compiler pipeline using our visual canvas.",
			icon: "",
			details: [
				"This tutorial will guide you through the canvas interface",
				"You'll learn to create nodes and connect them",
				"Each step builds upon the previous one",
				"Take your time to understand each concept"
			]
		},
		{
			id: 2,
			title: "Understanding the Toolbox",
			description: "The toolbox contains all the compiler phase blocks you can add to your canvas.",
			icon: "",
			details: [
				"Source Code: Input your program here",
				"Lexer: Breaks code into tokens",
				"Parser: Creates syntax trees from tokens",
				"Analyser: Performs semantic analysis",
				"Translator: Generates target code"
			]
		},
		{
			id: 3,
			title: "Adding Nodes to Canvas",
			description: "Click on any block in the toolbox to add it to your canvas workspace.",
			icon: "",
			details: [
				"Each phase can only be added once",
				"Nodes appear in the center of the canvas",
				"You can drag nodes to reposition them",
				"Click on a node to configure its settings"
			]
		},
		{
			id: 4,
			title: "Connecting Your Pipeline",
			description: "Create data flow by connecting nodes from their output to input anchors.",
			icon: "",
			details: [
				"Drag from a node's output anchor to another's input",
				"Follow the sequential order: Source → Lexer → Parser → Analyser → Translator",
				"Connections show how data flows through your compiler",
				"Invalid connections will be prevented automatically"
			]
		},
		{
			id: 5,
			title: "Configuring Phases",
			description: "Double-click any node to open its configuration panel and set up the phase.",
			icon: "",
			details: [
				"Source: Enter your code to compile",
				"Lexer: Define token rules and patterns",
				"Parser: Create grammar rules",
				"Analyser: Set up symbol tables",
				"Translator: Configure output settings"
			]
		},
		{
			id: 6,
			title: "Ready to Build!",
			description: "You're now ready to create your first compiler pipeline. Start by adding a Source Code node!",
			icon: "",
			details: [
				"Begin with the Source Code block from the toolbox",
				"Add phases sequentially and connect them",
				"Configure each phase as you build",
				"Use the help menu for additional guidance"
			]
		}
	];

	let currentStep = 0;
	let isVisible = false;

	$: if (show) {
		currentStep = 0;
		isVisible = true;
	}

	function nextStep() {
		if (currentStep < tutorialSteps.length - 1) {
			currentStep++;
		} else {
			closeTutorial();
		}
	}

	function previousStep() {
		if (currentStep > 0) {
			currentStep--;
		}
	}

	function skipTutorial() {
		// Mark as completed and close
		localStorage.setItem('hasSeenCanvasTutorial', 'true');
		closeTutorial();
	}

	function closeTutorial() {
		isVisible = false;
		// Mark as completed
		localStorage.setItem('hasSeenCanvasTutorial', 'true');
		setTimeout(() => {
			dispatch('close');
		}, 300);
	}

	$: currentTutorialStep = tutorialSteps[currentStep];
	$: isLastStep = currentStep === tutorialSteps.length - 1;
	$: isFirstStep = currentStep === 0;
</script>

{#if show && isVisible}
	<div 
		class="tutorial-backdrop" 
		transition:fade={{ duration: 300 }}
		on:click={skipTutorial}
		role="dialog"
		aria-modal="true"
		aria-labelledby="tutorial-title"
		aria-describedby="tutorial-description"
	>
		<div 
			class="tutorial-modal" 
			transition:fly={{ y: 50, duration: 400, easing: quintOut }}
			on:click|stopPropagation
		>
			<!-- Header -->
			<div class="tutorial-header">
				<div class="step-indicator">
					<span class="step-text">TUTORIAL ({currentStep + 1}/{tutorialSteps.length})</span>
				</div>
				<button 
					class="close-button" 
					on:click={skipTutorial}
					aria-label="Close tutorial"
					title="Close tutorial"
				>
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<line x1="18" y1="6" x2="6" y2="18"></line>
						<line x1="6" y1="6" x2="18" y2="18"></line>
					</svg>
				</button>
			</div>

			<!-- Content -->
			<div class="tutorial-content">
				{#if currentTutorialStep.icon}
					<div class="icon-container">
						<span class="tutorial-icon" role="img" aria-label="Step icon">{currentTutorialStep.icon}</span>
					</div>
				{/if}
				
				<h2 class="tutorial-title" id="tutorial-title">
					{currentTutorialStep.title}
				</h2>
				
				<p class="tutorial-description" id="tutorial-description">
					{currentTutorialStep.description}
				</p>

				<!-- Feature Details -->
				<div class="feature-details">
					{#each currentTutorialStep.details as detail, index}
						<div 
							class="detail-item"
							in:fly={{ x: -20, duration: 300, delay: 0 }}
						>
							<div class="detail-bullet"></div>
							<span class="detail-text">{detail}</span>
						</div>
					{/each}
				</div>
			</div>

			<!-- Progress Bar -->
			<div class="progress-container">
				<div class="progress-bar">
					<div 
						class="progress-fill" 
						style="width: {((currentStep + 1) / tutorialSteps.length) * 100}%"
					></div>
				</div>
				<span class="progress-text">{currentStep + 1} of {tutorialSteps.length}</span>
			</div>

			<!-- Navigation -->
			<div class="tutorial-navigation">
				<button 
					class="nav-button secondary" 
					on:click={previousStep}
					disabled={isFirstStep}
					aria-label="Previous step"
				>
					Previous
				</button>

				<button 
					class="skip-button" 
					on:click={skipTutorial}
					aria-label="Skip tutorial"
				>
					Skip the tutorial
				</button>

				<button 
					class="nav-button primary" 
					on:click={nextStep}
					aria-label={isLastStep ? "Finish tutorial" : "Next step"}
				>
					{isLastStep ? "Get Started!" : "Next"}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.tutorial-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		background: rgba(0, 0, 0, 0.75);
		backdrop-filter: blur(8px);
		-webkit-backdrop-filter: blur(8px);
		z-index: 10000;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
	}

	.tutorial-modal {
		background: linear-gradient(145deg, #ffffff, #f8fafc);
		border-radius: 1.5rem;
		box-shadow: 
			0 25px 50px -12px rgba(0, 0, 0, 0.25),
			0 0 0 1px rgba(255, 255, 255, 0.1);
		width: 100%;
		max-width: 600px;
		height: 700px;
		overflow: hidden;
		position: relative;
		display: flex;
		flex-direction: column;
	}

	.tutorial-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1.5rem 2rem 0;
		flex-shrink: 0;
	}

	.step-indicator {
		background: linear-gradient(135deg, #041a47, #0c2d5e);
		color: white;
		padding: 0.5rem 1rem;
		border-radius: 2rem;
		font-size: 0.875rem;
		font-weight: 600;
		letter-spacing: 0.05em;
	}

	.step-text {
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
	}

	.close-button {
		background: none;
		border: none;
		color: #64748b;
		cursor: pointer;
		padding: 0.5rem;
		border-radius: 0.5rem;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.close-button:hover {
		background: rgba(148, 163, 184, 0.1);
		color: #334155;
		transform: scale(1.05);
	}

	.tutorial-content {
		padding: 2rem;
		text-align: center;
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}

	.icon-container {
		margin-bottom: 1.5rem;
	}

	.tutorial-icon {
		font-size: 4rem;
		display: block;
		filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.1));
	}

	.tutorial-title {
		font-size: 2rem;
		font-weight: 700;
		color: #1e293b;
		margin: 0 0 1rem 0;
		line-height: 1.2;
	}

	.tutorial-description {
		font-size: 1.125rem;
		color: #64748b;
		margin: 0 0 2rem 0;
		line-height: 1.6;
		max-width: 500px;
		margin-left: auto;
		margin-right: auto;
	}

	.feature-details {
		text-align: left;
		max-width: 480px;
		margin: 0 auto;
	}

	.detail-item {
		display: flex;
		align-items: flex-start;
		margin-bottom: 1rem;
		padding: 0.75rem;
		background: rgba(4, 26, 71, 0.05);
		border-radius: 0.75rem;
		border-left: 3px solid #041a47;
	}

	.detail-bullet {
		width: 8px;
		height: 8px;
		background: linear-gradient(135deg, #041a47, #0c2d5e);
		border-radius: 50%;
		margin-right: 0.75rem;
		margin-top: 0.375rem;
		flex-shrink: 0;
	}

	.detail-text {
		color: #374151;
		font-size: 0.95rem;
		line-height: 1.5;
	}

	.progress-container {
		padding: 0 2rem 1rem;
		display: flex;
		align-items: center;
		gap: 1rem;
		flex-shrink: 0;
	}

	.progress-bar {
		flex: 1;
		height: 6px;
		background: #e2e8f0;
		border-radius: 3px;
		overflow: hidden;
	}

	.progress-fill {
		height: 100%;
		background: linear-gradient(90deg, #041a47, #0c2d5e);
		border-radius: 3px;
		transition: width 0.3s ease;
	}

	.progress-text {
		font-size: 0.875rem;
		color: #64748b;
		font-weight: 500;
		min-width: fit-content;
	}

	.tutorial-navigation {
		padding: 1.5rem 2rem 2rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: rgba(248, 250, 252, 0.8);
		border-top: 1px solid rgba(226, 232, 240, 0.8);
		flex-shrink: 0;
	}

	.nav-button {
		padding: 0.75rem 1.5rem;
		border-radius: 0.75rem;
		font-weight: 600;
		font-size: 0.95rem;
		cursor: pointer;
		transition: all 0.2s ease;
		border: none;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.nav-button.primary {
		background: linear-gradient(135deg, #041a47, #0c2d5e);
		color: white;
		box-shadow: 0 4px 14px 0 rgba(4, 26, 71, 0.39);
	}

	.nav-button.primary:hover:not(:disabled) {
		transform: translateY(-2px);
		box-shadow: 0 8px 25px 0 rgba(4, 26, 71, 0.5);
	}

	.nav-button.secondary {
		background: white;
		color: #64748b;
		border: 1px solid #e2e8f0;
		box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
	}

	.nav-button.secondary:hover:not(:disabled) {
		background: #f8fafc;
		border-color: #cbd5e1;
		transform: translateY(-1px);
	}

	.nav-button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none !important;
	}

	.skip-button {
		background: none;
		border: none;
		color: #94a3b8;
		font-size: 0.875rem;
		cursor: pointer;
		padding: 0.5rem 1rem;
		border-radius: 0.5rem;
		transition: all 0.2s ease;
		text-decoration: underline;
		text-underline-offset: 2px;
	}

	.skip-button:hover {
		color: #64748b;
		background: rgba(148, 163, 184, 0.1);
		text-decoration: none;
	}

	/* Dark mode support */
	:global(html.dark-mode) .tutorial-modal {
		background: linear-gradient(145deg, #1a2a4a, #2d3748);
		color: #f0f0f0;
	}

	:global(html.dark-mode) .tutorial-title {
		color: #ebeef1;
	}

	:global(html.dark-mode) .tutorial-description {
		color: #a0aec0;
	}

	:global(html.dark-mode) .detail-item {
		background: rgba(0, 26, 110, 0.15);
		border-left-color: #001A6E;
	}

	:global(html.dark-mode) .detail-text {
		color: #d1d5db;
	}

	:global(html.dark-mode) .progress-bar {
		background: #4a5568;
	}

	:global(html.dark-mode) .tutorial-navigation {
		background: rgba(26, 42, 74, 0.9);
		border-top-color: rgba(74, 85, 104, 0.8);
	}

	:global(html.dark-mode) .nav-button.secondary {
		background: #4a5568;
		color: #e2e8f0;
		border-color: #718096;
	}

	:global(html.dark-mode) .nav-button.secondary:hover:not(:disabled) {
		background: #2d3748;
		border-color: #4b5563;
	}

	:global(html.dark-mode) .nav-button.primary {
		background: linear-gradient(135deg, #001A6E, #002a8e);
		color: #ffffff;
		border-color: #374151;
	}

	:global(html.dark-mode) .nav-button.primary:hover:not(:disabled) {
		background: linear-gradient(135deg, #002a8e, #003bb3);
	}

	:global(html.dark-mode) .close-button {
		color: #a0aec0;
	}

	:global(html.dark-mode) .close-button:hover {
		background: rgba(74, 85, 104, 0.2);
		color: #d1d5db;
	}

	:global(html.dark-mode) .step-indicator {
		background: linear-gradient(135deg, #001A6E, #002a8e);
	}

	:global(html.dark-mode) .progress-text {
		color: #a0aec0;
	}

	/* Responsive design */
	@media (max-width: 640px) {
		.tutorial-modal {
			margin: 1rem;
			border-radius: 1rem;
			height: 700px;
		}

		.tutorial-header,
		.tutorial-content,
		.tutorial-navigation {
			padding-left: 1.5rem;
			padding-right: 1.5rem;
		}

		.tutorial-title {
			font-size: 1.5rem;
		}

		.tutorial-description {
			font-size: 1rem;
		}

		.tutorial-navigation {
			flex-direction: column;
			gap: 1rem;
		}

		.nav-button {
			width: 100%;
			justify-content: center;
		}

		.skip-button {
			order: -1;
		}
	}
</style>
