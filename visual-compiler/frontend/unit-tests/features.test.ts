import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import Features from '../src/lib/components/landing/Features.svelte';

describe('Features Component', () => {
	it('TestRender_Success: Renders features section with heading', () => {
		render(Features);

		expect(screen.getByText('How It Works')).toBeInTheDocument();
		expect(screen.getByRole('heading', { name: 'How It Works' })).toBeInTheDocument();
	});

	it('TestFeatureCards_Success: Renders all feature cards with correct content', () => {
		render(Features);

		// Test Block-Based Construction feature
		expect(screen.getByText('Block-Based Construction')).toBeInTheDocument();
		expect(screen.getByText(/Build compilers step-by-step using intuitive/)).toBeInTheDocument();

		// Test Visualize Artifacts feature
		expect(screen.getByText('Visualize Artifacts')).toBeInTheDocument();
		expect(screen.getByText(/See the outputs of each compilation phase/)).toBeInTheDocument();

		// Test Interactive Learning feature
		expect(screen.getByText('Interactive Learning')).toBeInTheDocument();
		expect(screen.getByText(/Experiment with formal concepts like RegEx/)).toBeInTheDocument();
	});

	it('TestFeatureCount_Success: Renders exactly 3 feature cards', () => {
		render(Features);

		const featureCards = screen.getAllByRole('heading', { level: 3 });
		expect(featureCards).toHaveLength(3);
	});

	it('TestFeatureTitles_Success: All feature titles are present as h3 headings', () => {
		render(Features);

		expect(screen.getByRole('heading', { name: 'Block-Based Construction' })).toBeInTheDocument();
		expect(screen.getByRole('heading', { name: 'Visualize Artifacts' })).toBeInTheDocument();
		expect(screen.getByRole('heading', { name: 'Interactive Learning' })).toBeInTheDocument();
	});

	it('TestFeatureDescriptions_Success: All feature descriptions are present', () => {
		render(Features);

		expect(screen.getByText(/Build compilers step-by-step using intuitive, connectable blocks/)).toBeInTheDocument();
		expect(screen.getByText(/See the outputs of each compilation phase, such as token streams/)).toBeInTheDocument();
		expect(screen.getByText(/Experiment with formal concepts like RegEx, CFGs, and Automata/)).toBeInTheDocument();
	});

	it('TestSectionStructure_Success: Has proper semantic structure', () => {
		const { container } = render(Features);

		const section = container.querySelector('section.features_section');
		expect(section).toBeInTheDocument();

		const featuresContainer = container.querySelector('.features_container');
		expect(featuresContainer).toBeInTheDocument();

		const featuresGrid = container.querySelector('.features_grid');
		expect(featuresGrid).toBeInTheDocument();
	});

	it('TestFeatureCards_Success: Each feature card has proper structure', () => {
		const { container } = render(Features);

		const featureCards = container.querySelectorAll('.feature_card');
		expect(featureCards).toHaveLength(3);

		featureCards.forEach(card => {
			expect(card.querySelector('.feature_icon')).toBeInTheDocument();
			expect(card.querySelector('.feature_title')).toBeInTheDocument();
			expect(card.querySelector('.feature_description')).toBeInTheDocument();
		});
	});

	it('TestSVGIcons_Success: SVG icons are rendered in feature cards', () => {
		const { container } = render(Features);

		const svgElements = container.querySelectorAll('svg');
		expect(svgElements).toHaveLength(3);

		// Check that each SVG has proper attributes
		svgElements.forEach(svg => {
			expect(svg).toHaveAttribute('width', '32');
			expect(svg).toHaveAttribute('height', '32');
			expect(svg).toHaveAttribute('viewBox', '0 0 24 24');
		});
	});

	it('TestAccessibility_Success: Has proper heading hierarchy', () => {
		render(Features);

		// Main heading should be h2
		expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('How It Works');
		
		// Feature titles should be h3
		const h3Headings = screen.getAllByRole('heading', { level: 3 });
		expect(h3Headings).toHaveLength(3);
	});
});


