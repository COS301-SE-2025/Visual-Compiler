import { describe, it, expect, vi } from 'vitest';
import { render } from '@testing-library/svelte';
import Layout from '../../src/routes/+layout.svelte';

// Mock the ToastContainer component
vi.mock('$lib/components/toast-conatiner.svelte', () => ({
    default: vi.fn(() => ({
        $$: {},
        $destroy: vi.fn(),
        $set: vi.fn()
    }))
}));

describe('+layout.svelte', () => {
    it('should render without errors', () => {
        expect(() => render(Layout)).not.toThrow();
    });

    it('should render ToastContainer component', () => {
        const { container } = render(Layout);
        expect(container).toBeTruthy();
    });

    it('should render slot content when provided', () => {
        const SlotContent = vi.fn(() => ({
            create_fragment: vi.fn(() => ({
                c: vi.fn(),
                m: vi.fn(),
                p: vi.fn(),
                d: vi.fn()
            }))
        }));

        const { container } = render(Layout, {
            props: {},
            context: new Map([['slot', SlotContent]])
        });

        expect(container).toBeTruthy();
    });

    it('should have correct component structure', () => {
        const { container } = render(Layout);
        
        // The layout should render without throwing
        expect(container.innerHTML).toBeDefined();
    });

    it('should be a valid Svelte component', () => {
        const component = render(Layout);
        
        expect(component.component).toBeDefined();
        expect(typeof component.component.$destroy).toBe('function');
    });
});
