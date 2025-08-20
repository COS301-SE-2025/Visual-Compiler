import { describe, it, expect, vi } from 'vitest';
import { redirect } from '@sveltejs/kit';
import { load } from '../../src/routes/+page.server';

// Mock the redirect function
vi.mock('@sveltejs/kit', () => ({
    redirect: vi.fn((status: number, location: string) => {
        throw new Error(`Redirecting to ${location} with status ${status}`);
    })
}));

describe('+page.server.ts', () => {
    describe('load function', () => {
        it('should redirect to landing-page with 307 status', () => {
            expect(() => load()).toThrow('Redirecting to /landing-page with status 307');
            expect(redirect).toHaveBeenCalledWith(307, '/landing-page');
        });

        it('should always throw redirect (no normal return path)', () => {
            expect(() => load()).toThrow();
        });

        it('should use temporary redirect status code', () => {
            try {
                load();
            } catch (error) {
                // Error thrown by redirect mock
            }
            
            expect(redirect).toHaveBeenCalledWith(307, '/landing-page');
            // 307 is Temporary Redirect, appropriate for this use case
        });

        it('should redirect to correct landing page path', () => {
            try {
                load();
            } catch (error) {
                // Error thrown by redirect mock
            }
            
            expect(redirect).toHaveBeenCalledWith(
                expect.any(Number),
                '/landing-page'
            );
        });
    });
});
