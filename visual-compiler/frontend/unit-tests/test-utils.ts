import { vi } from 'vitest';

// Standard mock functions that can be reused across tests
export const mockAddToast = vi.fn();
export const mockGoto = vi.fn();
export const mockInvalidateAll = vi.fn();

// Standard storage mocks
export const createStorageMock = () => {
	let store: { [key: string]: string } = {};
	return {
		getItem: vi.fn((key: string) => store[key] || null),
		setItem: vi.fn((key: string, value: string) => (store[key] = value.toString())),
		removeItem: vi.fn((key: string) => delete store[key]),
		clear: vi.fn(() => (store = {}))
	};
};

// Standard store mocks
export const createMockStore = (initialValue: any = null) => ({
	subscribe: vi.fn((fn) => {
		fn(initialValue);
		return { unsubscribe: vi.fn() };
	}),
	set: vi.fn(),
	update: vi.fn()
});

// Complete pipeline store mock
export const createPipelineStoreMock = () => ({
	pipelineStore: createMockStore({ nodes: [], connections: [], lastSaved: null }),
	phase_completion_status: createMockStore({
		source: false,
		lexer: false,
		parser: false,
		analyser: false,
		translator: false
	}),
	activePhase: createMockStore(null),
	setActivePhase: vi.fn(),
	updatePhaseStatus: vi.fn(),
	resetPipeline: vi.fn(),
	updatePipeline: vi.fn(),
	resetPhaseStatus: vi.fn()
});

// Standard project store mock
export const createProjectStoreMock = (projectName: string = 'test-project') => ({
	projectName: createMockStore(projectName)
});

// Standard toast store mock
export const createToastStoreMock = () => ({
	AddToast: mockAddToast,
	toasts: createMockStore([])
});

// Standard theme store mock
export const createThemeStoreMock = () => ({
	theme: createMockStore('light')
});

// Standard lexer store mock
export const createLexerStoreMock = () => ({
	lexerState: createMockStore({
		userInputRows: [{ type: '', regex: '', error: '' }],
		isRunning: false,
		results: null
	})
});

// SvelteKit mocks
export const createSvelteKitMocks = () => {
	// Mock SvelteKit runtime
	(globalThis as any).__SVELTEKIT_PAYLOAD__ = {
		data: {},
		errors: {}
	};

	return {
		navigation: {
			goto: mockGoto,
			invalidateAll: mockInvalidateAll
		},
		stores: {
			page: {
				subscribe: vi.fn(() => ({ unsubscribe: vi.fn() }))
			}
		}
	};
};

// Setup function for DOM mocks
export const setupDOMMocks = () => {
	const localStorageMock = createStorageMock();
	const sessionStorageMock = createStorageMock();

	Object.defineProperty(window, 'localStorage', { value: localStorageMock });
	Object.defineProperty(window, 'sessionStorage', { value: sessionStorageMock });

	// Mock fetch
	global.fetch = vi.fn();

	// Mock matchMedia
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

	// Mock Element.animate
	Element.prototype.animate = vi.fn().mockImplementation(() => ({
		onfinish: null,
		cancel: vi.fn(),
		finish: vi.fn(),
		play: vi.fn(),
		pause: vi.fn()
	}));

	return { localStorageMock, sessionStorageMock };
};

// Complete test setup function
export const setupTestEnvironment = (options: {
	projectName?: string;
	withSvelteKit?: boolean;
	withDOMMocks?: boolean;
} = {}) => {
	const {
		projectName = 'test-project',
		withSvelteKit = true,
		withDOMMocks = true
	} = options;

	let mocks: any = {};

	if (withSvelteKit) {
		mocks.svelteKit = createSvelteKitMocks();
	}

	if (withDOMMocks) {
		mocks.dom = setupDOMMocks();
	}

	// Reset all mocks
	vi.clearAllMocks();

	return {
		mocks,
		stores: {
			pipeline: createPipelineStoreMock(),
			project: createProjectStoreMock(projectName),
			toast: createToastStoreMock(),
			theme: createThemeStoreMock(),
			lexer: createLexerStoreMock()
		}
	};
};