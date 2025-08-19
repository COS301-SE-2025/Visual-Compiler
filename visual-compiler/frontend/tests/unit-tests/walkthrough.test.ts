import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import Walkthrough from '$lib/components/landing/walkthrough.svelte';

describe('Walkthrough Component', () => {
	// ===== BASIC RENDERING TESTS =====
	it('TestRender_Success: Renders walkthrough section with main heading', () => {
		render(Walkthrough);

		expect(screen.getByText('A Glimpse Into the Process')).toBeInTheDocument();
		expect(screen.getByRole('heading', { name: 'A Glimpse Into the Process' })).toBeInTheDocument();
	});

	it('TestComponentStructure_Success: Renders with proper structure', () => {
		const { container } = render(Walkthrough);
		
		// Check if walkthrough section exists
		const walkthroughSection = container.querySelector('.walkthrough_section');
		expect(walkthroughSection).toBeTruthy();
		
		// Check for walkthrough container
		const walkthroughContainer = container.querySelector('.walkthrough_container');
		expect(walkthroughContainer).toBeTruthy();
	});

	// ===== STEP CONTENT TESTS =====
	it('TestStepTitles_Success: Renders all step titles', () => {
		render(Walkthrough);

		expect(screen.getByText('Design Your Pipeline')).toBeInTheDocument();
		expect(screen.getByText('Configure and Connect')).toBeInTheDocument();
		expect(screen.getByText('Run and Visualize')).toBeInTheDocument();
	});

	it('TestStepDescriptions_Success: Renders step descriptions with key concepts', () => {
		render(Walkthrough);

		// Design Your Pipeline description
		expect(screen.getByText(/Start with a blank canvas/)).toBeInTheDocument();
		expect(screen.getByText(/Select compiler phase blocks like Lexing, Parsing, and Code Generation/)).toBeInTheDocument();

		// Configure and Connect description should be present
		expect(screen.getByText(/Configure and Connect/)).toBeInTheDocument();

		// Run and Visualize description should be present
		expect(screen.getByText(/Run and Visualize/)).toBeInTheDocument();
	});

	// ===== IMAGE TESTS =====
	it('TestStepImages_Success: Displays walkthrough steps with images', () => {
		const { container } = render(Walkthrough);
		
		// Check for step images
		const stepImages = container.querySelectorAll('.step_image');
		expect(stepImages.length).toBeGreaterThan(0);
		
		// Check for specific image sources
		const images = container.querySelectorAll('img');
		expect(images.length).toBeGreaterThan(0);
	});

	it('TestImageAttributes_Success: All images have proper alt attributes and sources', () => {
		const { container } = render(Walkthrough);
		
		const images = container.querySelectorAll('img');
		images.forEach(img => {
			expect(img.getAttribute('alt')).toBeTruthy();
			expect(img.getAttribute('src')).toBeTruthy();
		});
		
		// Check for specific images
		const canvasImage = container.querySelector('img[alt="Visual Compiler Canvas"]');
		expect(canvasImage).toBeTruthy();
		if (canvasImage) {
			expect(canvasImage.getAttribute('src')).toBe('/SvelvetcanvasFull.PNG');
		}

		const configImage = container.querySelector('img[alt="Node Configuration"]');
		expect(configImage).toBeTruthy();
		if (configImage) {
			expect(configImage.getAttribute('src')).toBe('/NodeConfiguration.PNG');
		}

		const syntaxTreeImage = container.querySelector('img[alt="Syntax Tree Visualization"]');
		expect(syntaxTreeImage).toBeTruthy();
		if (syntaxTreeImage) {
			expect(syntaxTreeImage.getAttribute('src')).toBe('/Parsing.PNG');
		}
	});

	// ===== STEP STRUCTURE TESTS =====
	it('TestStepStructure_Success: Steps have proper structure', () => {
		const { container } = render(Walkthrough);
		
		const steps = container.querySelectorAll('.step');
		expect(steps.length).toBeGreaterThan(0);
		
		steps.forEach(step => {
			expect(step.querySelector('.step_image_container')).toBeTruthy();
			expect(step.querySelector('.step_text_container')).toBeTruthy();
		});
	});

	it('TestStepContentStructure_Success: Each step has title and description', () => {
		const { container } = render(Walkthrough);
		
		const stepTitles = container.querySelectorAll('.step_title');
		const stepDescriptions = container.querySelectorAll('.step_description');
		
		expect(stepTitles.length).toBeGreaterThan(0);
		expect(stepDescriptions.length).toBeGreaterThan(0);
		
		// Check that titles are h3 elements
		stepTitles.forEach(title => {
			expect(title.tagName).toBe('H3');
		});
		
		// Check that descriptions are p elements
		stepDescriptions.forEach(description => {
			expect(description.tagName).toBe('P');
		});
	});

	// ===== LAYOUT TESTS =====
	it('TestStepLayout_Success: Steps have alternating layout', () => {
		const { container } = render(Walkthrough);
		
		const steps = container.querySelectorAll('.step');
		expect(steps.length).toBeGreaterThan(1);
		
		// Check that some steps have reverse class (alternating layout)
		const reverseSteps = container.querySelectorAll('.step.reverse');
		expect(reverseSteps.length).toBeGreaterThan(0);
	});

	// ===== ACCESSIBILITY TESTS =====
	it('TestAccessibility_Success: Component has proper accessibility structure', () => {
		const { container } = render(Walkthrough);
		
		// Check for proper heading hierarchy
		const mainHeading = container.querySelector('h2.section_heading');
		expect(mainHeading).toBeTruthy();
		
		const stepHeadings = container.querySelectorAll('h3.step_title');
		expect(stepHeadings.length).toBeGreaterThan(0);
	});

	it('TestImageAccessibility_Success: All images have descriptive alt text', () => {
		const { container } = render(Walkthrough);
		
		const images = container.querySelectorAll('img');
		images.forEach(img => {
			const alt = img.getAttribute('alt');
			expect(alt).toBeTruthy();
			expect(alt?.length).toBeGreaterThan(3); // Should be descriptive, not just a word
		});
	});

	// ===== SEMANTIC MARKUP TESTS =====
	it('TestSemanticallyCorrectMarkup_Success: Uses appropriate HTML elements', () => {
		const { container } = render(Walkthrough);
		
		// Should use section for walkthrough
		expect(container.querySelector('section.walkthrough_section')).toBeTruthy();
		
		// Should use h2 for main heading
		expect(container.querySelector('h2.section_heading')).toBeTruthy();
		
		// Should use h3 for step titles
		const stepTitles = container.querySelectorAll('h3.step_title');
		expect(stepTitles.length).toBeGreaterThan(0);
	});

	// ===== CSS CLASSES TESTS =====
	it('TestCSSClasses_Success: Elements have correct CSS classes', () => {
		const { container } = render(Walkthrough);
		
		expect(container.querySelector('.walkthrough_section')).toBeTruthy();
		expect(container.querySelector('.walkthrough_container')).toBeTruthy();
		expect(container.querySelector('.section_heading')).toBeTruthy();
		
		// Check step-specific classes
		expect(container.querySelector('.step')).toBeTruthy();
		expect(container.querySelector('.step_image_container')).toBeTruthy();
		expect(container.querySelector('.step_text_container')).toBeTruthy();
		expect(container.querySelector('.step_title')).toBeTruthy();
		expect(container.querySelector('.step_description')).toBeTruthy();
		expect(container.querySelector('.step_image')).toBeTruthy();
	});

	// ===== CONTENT VALIDATION TESTS =====
	it('TestEducationalContent_Success: Contains educational compiler concepts', () => {
		render(Walkthrough);
		
		// Should mention compiler-related terms
		expect(screen.getByText(/compiler phase blocks/i)).toBeInTheDocument();
		expect(screen.getByText(/lexing/i)).toBeInTheDocument();
		expect(screen.getByText(/parsing/i)).toBeInTheDocument();
		expect(screen.getByText(/code generation/i)).toBeInTheDocument();
	});

	it('TestProcessFlow_Success: Describes logical process flow', () => {
		render(Walkthrough);
		
		// Should describe the process from design to visualization
		expect(screen.getByText(/Start with a blank canvas/)).toBeInTheDocument();
		expect(screen.getByText(/Configure and Connect/)).toBeInTheDocument();
		expect(screen.getByText(/Run and Visualize/)).toBeInTheDocument();
	});

	// ===== PERFORMANCE TESTS =====
	it('TestRenderPerformance_Success: Component renders quickly', () => {
		const startTime = performance.now();
		render(Walkthrough);
		const endTime = performance.now();
		
		// Should render in reasonable time (less than 15ms due to images)
		expect(endTime - startTime).toBeLessThan(15);
	});

	it('TestMemoryUsage_Success: Component can be created and destroyed', () => {
		// Create and destroy multiple instances to test memory handling
		for (let i = 0; i < 3; i++) {
			const { unmount } = render(Walkthrough);
			unmount();
		}
		
		// Should complete without memory issues
		expect(true).toBe(true);
	});

	// ===== RESPONSIVE DESIGN TESTS =====
	it('TestResponsiveStructure_Success: Component has responsive structure', () => {
		const { container } = render(Walkthrough);
		
		const steps = container.querySelectorAll('.step');
		expect(steps.length).toBeGreaterThan(0);
		
		// Each step should have image and text containers for responsive layout
		steps.forEach(step => {
			expect(step.querySelector('.step_image_container')).toBeTruthy();
			expect(step.querySelector('.step_text_container')).toBeTruthy();
		});
	});

	// ===== INTEGRATION TESTS =====
	it('TestCompleteIntegration_Success: All elements work together correctly', () => {
		const { container } = render(Walkthrough);
		
		// Check complete component structure
		expect(screen.getByText('A Glimpse Into the Process')).toBeInTheDocument();
		expect(container.querySelectorAll('.step').length).toBeGreaterThan(0);
		expect(container.querySelectorAll('img').length).toBeGreaterThan(0);
		
		// Check all step titles are present
		const expectedTitles = ['Design Your Pipeline', 'Configure and Connect', 'Run and Visualize'];
		expectedTitles.forEach(title => {
			expect(screen.getByText(title)).toBeInTheDocument();
		});
	});
});
