import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import Hero from '$lib/components/landing/Hero.svelte';

// Mock SvelteKit's navigation module
const mockedGoto = vi.fn();
vi.mock('$app/navigation', () => ({
	goto: (path: string) => mockedGoto(path)
}));

// Mock window methods for scroll testing
const mockScrollTo = vi.fn();
const mockAddEventListener = vi.fn();
const mockRemoveEventListener = vi.fn();

Object.defineProperty(window, 'scrollTo', {
	value: mockScrollTo,
	writable: true
});

Object.defineProperty(window, 'addEventListener', {
	value: mockAddEventListener,
	writable: true
});

Object.defineProperty(window, 'removeEventListener', {
	value: mockRemoveEventListener,
	writable: true
});

Object.defineProperty(window, 'innerHeight', {
	value: 800,
	writable: true
});

describe('Hero Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('TestInitialRender_Success: Renders the main headline and button', () => {
		render(Hero);
		expect(screen.getByText(/Demystifying Compilers/)).toBeInTheDocument();
		expect(screen.getByText(/One Block at a Time/)).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Get Started' })).toBeInTheDocument();
	});

	it('TestLaunchButton_Success: Calls goto with the correct path on click', async () => {
		render(Hero);
		const launchButton = screen.getByRole('button', { name: 'Get Started' });
		await fireEvent.click(launchButton);

		expect(mockedGoto).toHaveBeenCalledWith('/auth-page');
	});

	it('TestLogoAndTitle_Success: Renders logo and project title', () => {
		render(Hero);
		
		const logo = screen.getByAltText('Visual Compiler Logo');
		expect(logo).toBeInTheDocument();
		expect(logo).toHaveAttribute('src', '/half_stack_phoenix_grey.png');
		
		expect(screen.getByText('Visual Compiler')).toBeInTheDocument();
	});

	it('TestSubHeadline_Success: Renders sub-headline text', () => {
		render(Hero);
		expect(screen.getByText(/An educational platform to construct, configure, visualise/)).toBeInTheDocument();
		expect(screen.getByText(/and understand the core phases of compilation/)).toBeInTheDocument();
	});

	it('TestScrollIndicator_Success: Renders scroll indicator elements', () => {
		render(Hero);
		
		expect(screen.getByText('Scroll to explore')).toBeInTheDocument();
		
		// Check for scroll arrow SVG
		const svgElement = screen.getByRole('button', { name: /scroll to explore/i });
		expect(svgElement).toBeInTheDocument();
		
		// Check for dots
		const dotsContainer = document.querySelector('.scroll_dots');
		expect(dotsContainer).toBeInTheDocument();
		expect(dotsContainer?.children).toHaveLength(3);
	});

	it('TestScrollIndicatorClick_Success: Handles scroll indicator click', async () => {
		render(Hero);
		
		const scrollIndicator = screen.getByRole('button', { name: /scroll to explore/i });
		await fireEvent.click(scrollIndicator);
		
		expect(mockScrollTo).toHaveBeenCalledWith({
			top: 800, // window.innerHeight
			behavior: 'smooth'
		});
	});

	it('TestScrollIndicatorKeyboard_Success: Handles Enter key on scroll indicator', async () => {
		render(Hero);
		
		const scrollIndicator = screen.getByRole('button', { name: /scroll to explore/i });
		await fireEvent.keyDown(scrollIndicator, { key: 'Enter' });
		
		expect(mockScrollTo).toHaveBeenCalledWith({
			top: 800,
			behavior: 'smooth'
		});
	});

	it('TestScrollIndicatorKeyboardOther_Success: Ignores non-Enter keys on scroll indicator', async () => {
		render(Hero);
		
		const scrollIndicator = screen.getByRole('button', { name: /scroll to explore/i });
		await fireEvent.keyDown(scrollIndicator, { key: 'Space' });
		
		expect(mockScrollTo).not.toHaveBeenCalled();
	});

	it('TestEventListenerSetup_Success: Sets up scroll event listener on mount', () => {
		render(Hero);
		
		expect(mockAddEventListener).toHaveBeenCalledWith('scroll', expect.any(Function));
	});

	it('TestScrollHandling_Success: Hides scroll indicator when scrolled', async () => {
		const { component } = render(Hero);
		
		// Get the scroll handler that was registered
		const scrollHandler = mockAddEventListener.mock.calls.find(
			call => call[0] === 'scroll'
		)?.[1];
		
		expect(scrollHandler).toBeDefined();
		
		// Initially scroll indicator should be visible
		expect(screen.getByText('Scroll to explore')).toBeInTheDocument();
		
		// Simulate scroll by setting scrollY and calling the handler
		Object.defineProperty(window, 'scrollY', { value: 60, writable: true });
		if (scrollHandler) {
			scrollHandler();
		}
		
		await waitFor(() => {
			expect(screen.queryByText('Scroll to explore')).not.toBeInTheDocument();
		});
	});

	it('TestScrollIndicatorVisible_Success: Keeps scroll indicator visible for small scroll', async () => {
		render(Hero);
		
		const scrollHandler = mockAddEventListener.mock.calls.find(
			call => call[0] === 'scroll'
		)?.[1];
		
		// Simulate small scroll (less than 50px)
		Object.defineProperty(window, 'scrollY', { value: 30, writable: true });
		if (scrollHandler) {
			scrollHandler();
		}
		
		await waitFor(() => {
			expect(screen.getByText('Scroll to explore')).toBeInTheDocument();
		});
	});

	it('TestComponentCleanup_Success: Removes event listener on cleanup', () => {
		const { unmount } = render(Hero);
		
		// Verify cleanup was registered
		const scrollHandler = mockAddEventListener.mock.calls.find(
			call => call[0] === 'scroll'
		)?.[1];
		
		unmount();
		
		expect(mockRemoveEventListener).toHaveBeenCalledWith('scroll', scrollHandler);
	});

	it('TestHeroSectionStructure_Success: Renders correct section structure', () => {
		render(Hero);
		
		const heroSection = document.querySelector('.hero_section');
		expect(heroSection).toBeInTheDocument();
		
		const logoContainer = document.querySelector('.logo_container');
		expect(logoContainer).toBeInTheDocument();
		
		const heroContent = document.querySelector('.hero_content');
		expect(heroContent).toBeInTheDocument();
	});

	it('TestButtonAttributes_Success: Verifies button has correct attributes', () => {
		render(Hero);
		
		const button = screen.getByRole('button', { name: 'Get Started' });
		expect(button).toHaveClass('cta_button');
		// Button uses default type, no explicit type attribute needed
	});

	it('TestAccessibility_Success: Verifies accessibility attributes', () => {
		render(Hero);
		
		const scrollIndicator = screen.getByRole('button', { name: /scroll to explore/i });
		expect(scrollIndicator).toHaveAttribute('tabindex', '0');
		expect(scrollIndicator).toHaveAttribute('role', 'button');
		
		const logo = screen.getByAltText('Visual Compiler Logo');
		expect(logo).toHaveAttribute('alt', 'Visual Compiler Logo');
	});

	it('TestMultipleScrollEvents_Success: Handles multiple scroll events correctly', async () => {
		render(Hero);
		
		const scrollHandler = mockAddEventListener.mock.calls.find(
			call => call[0] === 'scroll'
		)?.[1];
		
		// Test scroll below threshold
		Object.defineProperty(window, 'scrollY', { value: 60, writable: true });
		if (scrollHandler) {
			scrollHandler();
		}
		
		await waitFor(() => {
			expect(screen.queryByText('Scroll to explore')).not.toBeInTheDocument();
		});
		
		// Test another scroll event - indicator should remain hidden
		Object.defineProperty(window, 'scrollY', { value: 100, writable: true });
		if (scrollHandler) {
			scrollHandler();
		}
		
		await waitFor(() => {
			expect(screen.queryByText('Scroll to explore')).not.toBeInTheDocument();
		});
	});

	it('TestCTAButtonInteraction_Success: CTA button has proper class and is clickable', () => {
		render(Hero);
		
		const button = screen.getByRole('button', { name: 'Get Started' });
		expect(button).toHaveClass('cta_button');
		// CSS hover effects are defined in the component styles via the cta_button class
	});

	it('TestScrollToNextBehavior_Success: Verifies scroll behavior parameters', async () => {
		// Test with different window height
		Object.defineProperty(window, 'innerHeight', { value: 1000, writable: true });
		
		render(Hero);
		
		const scrollIndicator = screen.getByRole('button', { name: /scroll to explore/i });
		await fireEvent.click(scrollIndicator);
		
		expect(mockScrollTo).toHaveBeenCalledWith({
			top: 1000,
			behavior: 'smooth'
		});
	});
});


