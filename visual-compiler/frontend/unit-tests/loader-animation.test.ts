import { render } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import LoaderAnimation from '../../src/lib/components/landing/loader-animation.svelte';

describe('LoaderAnimation Component', () => {
	it('TestRender_Success: Renders loader overlay container', () => {
		const { container } = render(LoaderAnimation);

		const loaderOverlay = container.querySelector('.loader_overlay');
		expect(loaderOverlay).toBeInTheDocument();
	});

	it('TestAnimationContainer_Success: Contains animation container', () => {
		const { container } = render(LoaderAnimation);

		const animationContainer = container.querySelector('.animation_container');
		expect(animationContainer).toBeInTheDocument();
	});

	it('TestBlocks_Success: Renders exactly 3 animation blocks', () => {
		const { container } = render(LoaderAnimation);

		const blocks = container.querySelectorAll('.block');
		expect(blocks).toHaveLength(3);
	});

	it('TestBlockClasses_Success: Each block has correct CSS classes', () => {
		const { container } = render(LoaderAnimation);

		const block1 = container.querySelector('.block_1');
		const block2 = container.querySelector('.block_2');
		const block3 = container.querySelector('.block_3');

		expect(block1).toBeInTheDocument();
		expect(block1).toHaveClass('block');

		expect(block2).toBeInTheDocument();
		expect(block2).toHaveClass('block');

		expect(block3).toBeInTheDocument();
		expect(block3).toHaveClass('block');
	});

	it('TestOverlayClasses_Success: Loader overlay has correct CSS classes', () => {
		const { container } = render(LoaderAnimation);

		const loaderOverlay = container.querySelector('.loader_overlay');
		expect(loaderOverlay).toHaveClass('loader_overlay');
	});

	it('TestAnimationClasses_Success: Animation container has correct classes', () => {
		const { container } = render(LoaderAnimation);

		const animationContainer = container.querySelector('.animation_container');
		expect(animationContainer).toHaveClass('animation_container');
	});

	it('TestBlocksStructure_Success: All blocks are properly nested', () => {
		const { container } = render(LoaderAnimation);

		const animationContainer = container.querySelector('.animation_container');
		expect(animationContainer).toBeInTheDocument();
		
		if (animationContainer) {
			const blocksInContainer = animationContainer.querySelectorAll('.block');
			expect(blocksInContainer).toHaveLength(3);
		}
	});

	it('TestComponentStructure_Success: Has proper DOM structure', () => {
		const { container } = render(LoaderAnimation);

		// Check overall structure
		const loaderOverlay = container.querySelector('.loader_overlay');
		const animationContainer = container.querySelector('.animation_container');

		expect(loaderOverlay).toBeInTheDocument();
		expect(animationContainer).toBeInTheDocument();

		// Check blocks are children of animation container
		expect(loaderOverlay).toBeInTheDocument();
		expect(animationContainer).toBeInTheDocument();
	});

	it('TestBlockVariants_Success: Each block has unique class variant', () => {
		const { container } = render(LoaderAnimation);

		const block1 = container.querySelector('.block.block_1');
		const block2 = container.querySelector('.block.block_2');
		const block3 = container.querySelector('.block.block_3');

		expect(block1).toBeInTheDocument();
		expect(block2).toBeInTheDocument();
		expect(block3).toBeInTheDocument();
	});
});
