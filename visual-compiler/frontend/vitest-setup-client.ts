import '@testing-library/jest-dom/vitest';
import { vi, expect } from 'vitest';
import * as matchers from '@testing-library/jest-dom/matchers';

expect.extend(matchers);

// required for svelte5 + jsdom as jsdom does not support matchMedia
Object.defineProperty(window, 'matchMedia', {
	writable: true,
	enumerable: true,
	value: vi.fn().mockImplementation((query) => ({
		matches: false,
		media: query,
		onchange: null,
		addEventListener: vi.fn(),
		removeEventListener: vi.fn(),
		dispatchEvent: vi.fn()
	}))
});

// Mock for Element.animate - JSDOM doesn't support Web Animations API
Element.prototype.animate = vi.fn().mockImplementation(() => ({
	onfinish: null,
	cancel: vi.fn(),
	finish: vi.fn(),
	play: vi.fn(),
	pause: vi.fn()
}));
