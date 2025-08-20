import { render } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import LandingPage from '../src/routes/landing-page/+page.svelte';

// Mock all the landing page components with simple function mocks
vi.mock('$lib/components/landing/loader-animation.svelte', () => ({
	default: () => null
}));

vi.mock('$lib/components/landing/hero.svelte', () => ({
	default: () => null
}));

vi.mock('$lib/components/landing/features.svelte', () => ({
	default: () => null
}));

vi.mock('$lib/components/landing/walkthrough.svelte', () => ({
	default: () => null
}));

vi.mock('$lib/components/landing/info-section.svelte', () => ({
	default: () => null
}));

describe('LandingPage Component', () => {
	it('TestBasicRender_Success: Component renders without crashing', () => {
		const { container } = render(LandingPage);
		expect(container).toBeTruthy();
	});

	it('TestComponentStructure_Success: Component has basic structure', () => {
		const { container } = render(LandingPage);
		// Component should render some content
		expect(container.innerHTML).toBeTruthy();
	});
});


