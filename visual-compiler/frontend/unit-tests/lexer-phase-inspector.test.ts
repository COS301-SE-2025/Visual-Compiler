import { render } from '@testing-library/svelte';
import PhaseInspector from '../src/lib/components/lexer/phase-inspector.svelte';
import { AddToast } from '../src/lib/stores/toast';
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../src/lib/stores/toast', () => ({
  AddToast: vi.fn()
}));

vi.mock('$app/stores', () => ({
  page: {
    subscribe: vi.fn((fn) => {
      fn({ url: { pathname: '/lexer' }, params: {}, route: { id: null } });
      return { unsubscribe: vi.fn() };
    })
  }
}));

vi.mock('$app/navigation', () => ({
  goto: vi.fn()
}));

vi.mock('$lib/stores/stores', () => ({
  loading: {
    subscribe: vi.fn((fn) => {
      fn(false);
      return { unsubscribe: vi.fn() };
    })
  },
  projectStore: {
    subscribe: vi.fn((fn) => {
      fn({
        activeProjectId: 'test-project',
        activeProjectName: 'Test Project',
        tabs: [],
        activeTabId: null
      });
      return { unsubscribe: vi.fn() };
    })
  },
  lexerOutput: {
    subscribe: vi.fn((fn) => {
      fn([]);
      return { unsubscribe: vi.fn() };
    })
  }
}));

Object.defineProperty(window, 'localStorage', {
  value: {
    getItem: vi.fn(),
    setItem: vi.fn()
  }
});

global.fetch = vi.fn();

describe('Lexer Phase Inspector', () => {
  const sourceCode = 'let x = 5;';
  let mockFetch: any;

  beforeEach(() => {
    mockFetch = global.fetch as any;
    vi.clearAllMocks();
  });

  it('renders without crashing', () => {
    render(PhaseInspector, {
      props: {
        sourceCode: sourceCode
      }
    });
    expect(true).toBe(true);
  });
});
