import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import Hero from '$lib/components/landing/Hero.svelte';

// Mock SvelteKit's navigation module
const mockedGoto = vi.fn();
vi.mock('$app/navigation', () => ({
	goto: (path: string) => mockedGoto(path)
}));

describe('Hero Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		mockedGoto.mockResolvedValue(undefined);
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	// ===== BASIC RENDERING TESTS =====
	it('TestInitialRender_Success: Renders the main headline and button', () => {
		render(Hero);
		expect(screen.getByText('Demystifying Compilers, One Block at a Time.')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Get Started' })).toBeInTheDocument();
	});

	it('TestHeroStructure_Success: Renders hero section with proper structure', () => {
		const { container } = render(Hero);
		
		// Check if hero section exists
		const heroSection = container.querySelector('.hero_section');
		expect(heroSection).toBeTruthy();
		
		// Check for logo container
		const logoContainer = container.querySelector('.logo_container');
		expect(logoContainer).toBeTruthy();
		
		// Check for hero content
		const heroContent = container.querySelector('.hero_content');
		expect(heroContent).toBeTruthy();
	});

	it('TestLogoElements_Success: Renders logo and project title', () => {
		const { container } = render(Hero);
		
		// Check for logo image
		const logoImage = container.querySelector('img[alt="Visual Compiler Logo"]');
		expect(logoImage).toBeTruthy();
		if (logoImage) {
			expect(logoImage.getAttribute('src')).toBe('/half_stack_phoenix_grey.png');
		}
		
		// Check for project title
		expect(screen.getByText('Visual Compiler')).toBeInTheDocument();
	});

	it('TestHeroContent_Success: Displays correct hero content text', () => {
		render(Hero);
		
		// Check for main headline
		expect(screen.getByText('Demystifying Compilers, One Block at a Time.')).toBeInTheDocument();
		
		// Check for sub headline
		expect(screen.getByText('An educational platform to construct, configure, visualise and understand the core phases of compilation.')).toBeInTheDocument();
	});

	// ===== INTERACTION TESTS =====
	it('TestLaunchButton_Success: Calls goto with the correct path on click', async () => {
		render(Hero);
		const launchButton = screen.getByRole('button', { name: 'Get Started' });
		await fireEvent.click(launchButton);

		// Assert that our mocked navigation function was called with the correct path
		expect(mockedGoto).toHaveBeenCalledWith('/auth-page');
		expect(mockedGoto).toHaveBeenCalledTimes(1);
	});

	it('TestButtonClickability_Success: CTA button is properly clickable', async () => {
		render(Hero);
		
		const ctaButton = screen.getByRole('button', { name: 'Get Started' });
		expect(ctaButton).toBeEnabled();
		
		// Multiple clicks should work
		await fireEvent.click(ctaButton);
		await fireEvent.click(ctaButton);
		
		expect(mockedGoto).toHaveBeenCalledTimes(2);
	});

	// ===== ACCESSIBILITY TESTS =====
	it('TestAccessibility_Success: Component has proper accessibility attributes', () => {
		const { container } = render(Hero);
		
		// Check for proper alt text on images
		const logoImage = container.querySelector('img');
		expect(logoImage).toBeTruthy();
		if (logoImage) {
			expect(logoImage.getAttribute('alt')).toBe('Visual Compiler Logo');
		}
		
		// Check that button is properly labeled
		const ctaButton = screen.getByRole('button', { name: 'Get Started' });
		expect(ctaButton.textContent?.trim()).toBe('Get Started');
	});

	it('TestSemanticallyCorrectMarkup_Success: Uses appropriate HTML elements', () => {
		const { container } = render(Hero);
		
		// Should use section for hero
		expect(container.querySelector('section.hero_section')).toBeTruthy();
		
		// Should use h1 for main headline
		expect(container.querySelector('h1.main_headline')).toBeTruthy();
		
		// Should use paragraph for sub headline
		expect(container.querySelector('p.sub_headline')).toBeTruthy();
	});

	// ===== VISUAL STRUCTURE TESTS =====
	it('TestCSSClasses_Success: Elements have correct CSS classes', () => {
		const { container } = render(Hero);
		
		expect(container.querySelector('.hero_section')).toBeTruthy();
		expect(container.querySelector('.logo_container')).toBeTruthy();
		expect(container.querySelector('.project_title')).toBeTruthy();
		expect(container.querySelector('.hero_content')).toBeTruthy();
		expect(container.querySelector('.main_headline')).toBeTruthy();
		expect(container.querySelector('.sub_headline')).toBeTruthy();
		expect(container.querySelector('.cta_button')).toBeTruthy();
	});

	it('TestImageAssets_Success: Logo image has correct attributes', () => {
		const { container } = render(Hero);
		
		const logoImage = container.querySelector('.logo');
		expect(logoImage).toBeTruthy();
		if (logoImage) {
			expect(logoImage.tagName).toBe('IMG');
			expect(logoImage.getAttribute('alt')).toBe('Visual Compiler Logo');
			expect(logoImage.getAttribute('src')).toBe('/half_stack_phoenix_grey.png');
			expect(logoImage.classList.contains('logo')).toBe(true);
		}
	});

	// ===== RESPONSIVE DESIGN TESTS =====
	it('TestResponsiveStructure_Success: Component has responsive structure', () => {
		const { container } = render(Hero);
		
		// Check that the hero section exists
		const heroSection = container.querySelector('.hero_section');
		expect(heroSection).toBeTruthy();
		
		// Check that logo and content containers exist
		expect(container.querySelector('.logo_container')).toBeTruthy();
		expect(container.querySelector('.hero_content')).toBeTruthy();
	});

	// ===== COMPONENT INTEGRATION TESTS =====
	it('TestComponentIntegration_Success: All elements work together correctly', async () => {
		const { container } = render(Hero);
		
		// Check complete component structure
		expect(container.querySelector('.hero_section')).toBeTruthy();
		expect(screen.getByText('Visual Compiler')).toBeInTheDocument();
		expect(screen.getByText('Demystifying Compilers, One Block at a Time.')).toBeInTheDocument();
		expect(screen.getByText('An educational platform to construct, configure, visualise and understand the core phases of compilation.')).toBeInTheDocument();
		
		// Test interaction
		const button = screen.getByRole('button', { name: 'Get Started' });
		await fireEvent.click(button);
		expect(mockedGoto).toHaveBeenCalledWith('/auth-page');
	});

	// ===== PERFORMANCE TESTS =====
	it('TestRenderPerformance_Success: Component renders quickly', () => {
		const startTime = performance.now();
		render(Hero);
		const endTime = performance.now();
		
		// Should render in reasonable time (less than 10ms for simple component)
		expect(endTime - startTime).toBeLessThan(10);
	});

	it('TestMemoryUsage_Success: Component can be created and destroyed', () => {
		// Create and destroy multiple instances to test memory handling
		for (let i = 0; i < 5; i++) {
			const { unmount } = render(Hero);
			unmount();
		}
		
		// Should complete without memory issues
		expect(true).toBe(true);
	});
});
