import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/svelte';
import Features from '$lib/components/landing/Features.svelte';

describe('Features Component', () => {
	beforeEach(() => {
		// Reset any state before each test
		document.body.innerHTML = '';
	});

	it('renders the features section', () => {
		const { container } = render(Features);
		
		const section = container.querySelector('section.features_section');
		expect(section).toBeInTheDocument();
	});

	it('displays the correct heading', () => {
		render(Features);
		
		const heading = screen.getByText('How It Works');
		expect(heading).toBeInTheDocument();
		expect(heading.tagName).toBe('H2');
	});

	it('renders all feature cards', () => {
		const { container } = render(Features);
		
		// Check for feature cards by class
		const featureCards = container.querySelectorAll('.feature_card');
		expect(featureCards).toHaveLength(3);
	});

	it('displays feature icons', () => {
		render(Features);
		
		// Check for SVG icons
		const svgIcons = document.querySelectorAll('svg');
		expect(svgIcons.length).toBeGreaterThan(0);
	});

	it('has proper structure for feature content', () => {
		render(Features);
		
		// Check for feature titles and descriptions
		const featureTitles = screen.getAllByRole('heading', { level: 3 });
		expect(featureTitles.length).toBe(3);
		
		const featureDescriptions = document.querySelectorAll('p.feature_description');
		expect(featureDescriptions.length).toBe(3);
	});

	it('applies proper CSS classes', () => {
		const { container } = render(Features);
		
		const section = container.querySelector('section.features_section');
		expect(section).toHaveClass('features_section');
	});

	it('is responsive', () => {
		const { container } = render(Features);
		
		// Check for responsive grid classes
		const gridContainer = container.querySelector('.features_grid');
		expect(gridContainer).toBeInTheDocument();
	});

	it('has accessible headings hierarchy', () => {
		render(Features);
		
		// Main heading should be h2
		const mainHeading = screen.getByRole('heading', { level: 2 });
		expect(mainHeading).toBeInTheDocument();
		
		// Feature headings should be h3
		const featureHeadings = screen.getAllByRole('heading', { level: 3 });
		expect(featureHeadings.length).toBe(3);
	});

	it('renders without errors', () => {
		expect(() => render(Features)).not.toThrow();
	});

	it('has proper semantic structure', () => {
		const { container } = render(Features);
		
		// Should have proper section element
		const section = container.querySelector('section');
		expect(section).toBeInTheDocument();
		
		// Should have container for features
		const featuresContainer = container.querySelector('.features_container');
		expect(featuresContainer).toBeInTheDocument();
		
		// Should have grid for layout
		const featuresGrid = container.querySelector('.features_grid');
		expect(featuresGrid).toBeInTheDocument();
	});

	it('displays all feature titles correctly', () => {
		render(Features);
		
		expect(screen.getByText('Block-Based Construction')).toBeInTheDocument();
		expect(screen.getByText('Visualize Artifacts')).toBeInTheDocument();
		expect(screen.getByText('Interactive Learning')).toBeInTheDocument();
	});

	it('has proper icon structure', () => {
		const { container } = render(Features);
		
		const featureIcons = container.querySelectorAll('.feature_icon');
		expect(featureIcons).toHaveLength(3);
		
		featureIcons.forEach(icon => {
			const svg = icon.querySelector('svg');
			expect(svg).toBeInTheDocument();
		});
	});
});
