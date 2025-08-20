import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import Hero from '$lib/components/landing/Hero.svelte';

// Mock SvelteKit's navigation module
const mockedGoto = vi.fn();
vi.mock('$app/navigation', () => ({
	goto: (path: string) => mockedGoto(path)
}));

describe('hero Component', () => {
	it('TestInitialRender_Success: Renders the main headline and button', () => {
		render(Hero);
		expect(screen.getByText('Demystifying Compilers, One Block at a Time.')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Get Started' })).toBeInTheDocument();
	});

	it('TestLaunchButton_Success: Calls goto with the correct path on click', async () => {
		render(Hero);
		const launchButton = screen.getByRole('button', { name: 'Get Started' });
		await fireEvent.click(launchButton);

		// Assert that our mocked navigation function was called with the correct path
		expect(mockedGoto).toHaveBeenCalledWith('/auth-page');
	});
});


