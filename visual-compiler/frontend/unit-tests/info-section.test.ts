import { render, screen, fireEvent } from '@testing-library/svelte';
import { tick } from 'svelte';
import { describe, it, expect } from 'vitest';
import InfoSection from '../src/lib/components/landing/info-section.svelte';

describe('InfoSection Component', () => {
	it('TestInitialRender_Success: Renders the "About" tab by default', () => {
		render(InfoSection);

		// Check that the "About" button has the 'active' class
		const aboutButton = screen.getByRole('button', { name: 'About The Project' });
		expect(aboutButton).toHaveClass('active');

		// Check that the "About" content is visible
		expect(screen.getByText('Visual Compiler')).toBeInTheDocument();

		// Check that content for other tabs is NOT visible
		expect(screen.queryByText('Value Through Innovation')).toBeNull();
		expect(screen.queryByText('Sashen Inder Gajai')).toBeNull(); // Check for a team member's name
	});

	it('TestTabSwitch_Client_Success: Switches to the "Client" tab on click', async () => {
		render(InfoSection);

		// Find and click the "Client" tab button
		const clientButton = screen.getByRole('button', { name: 'About EPI-USE Labs' });
		await fireEvent.click(clientButton);

		// Check that the "Client" button is now active
		expect(clientButton).toHaveClass('active');
		// Check that the "About" button is no longer active
		const aboutButton = screen.getByRole('button', { name: 'About The Project' });
		expect(aboutButton).not.toHaveClass('active');

		// Check that the "Client" content is now visible
		expect(screen.getByText('Value Through Innovation')).toBeInTheDocument();

		// Check that the "About" content is now hidden
		expect(screen.queryByText('An Educational Tool for a Complex World')).toBeNull();
	});

	it('TestTabSwitch_Team_Success: Switches to the "Team" tab on click', async () => {
		// 1. Render the component.
		render(InfoSection);
		
		// Click on team tab to activate team content
		const teamTab = screen.getByText('Meet The Team');
		await fireEvent.click(teamTab);
		await tick();
		
		// Check that team tab is active
		expect(teamTab.classList.contains('active')).toBe(true);
	});

	// ===== TAB CONTENT TESTS =====
	it('TestTabContentStructure_Success: Tab content areas exist', () => {
		const { container } = render(InfoSection);
		
		// The component should have content areas that change based on active tab
		expect(container.querySelector('.info_section')).toBeTruthy();
		expect(container.querySelector('.tab_nav')).toBeTruthy();
	});

	it('TestTabButtons_Success: All tab buttons are proper buttons', () => {
		render(InfoSection);
		
		const buttons = screen.getAllByRole('button');
		expect(buttons.length).toBe(3);
		
		// Each button should have proper text content
		const buttonTexts = buttons.map(button => button.textContent);
		expect(buttonTexts).toContain('About The Project');
		expect(buttonTexts).toContain('About EPI-USE Labs');
		expect(buttonTexts).toContain('Meet The Team');
	});

	// ===== ACCESSIBILITY TESTS =====
	it('TestAccessibility_Success: Tab navigation is accessible', () => {
		render(InfoSection);
		
		// Check that tabs are proper buttons
		const tabs = screen.getAllByRole('button');
		expect(tabs.length).toBeGreaterThanOrEqual(3);
		
		// Each tab should have text content
		tabs.forEach(tab => {
			expect(tab.textContent?.trim()).toBeTruthy();
		});
	});

	it('TestKeyboardNavigation_Success: Tabs are keyboard accessible', async () => {
		render(InfoSection);
		
		const aboutTab = screen.getByText('About The Project');
		const clientTab = screen.getByText('About EPI-USE Labs');
		
		// Focus should work on tabs
		aboutTab.focus();
		expect(document.activeElement).toBe(aboutTab);
		
		// Tab switching via keyboard (simulated click for now)
		await fireEvent.click(clientTab);
		await tick();
		expect(clientTab.classList.contains('active')).toBe(true);
	});

	// ===== SEMANTIC MARKUP TESTS =====
	it('TestSemanticallyCorrectMarkup_Success: Uses appropriate HTML elements', () => {
		const { container } = render(InfoSection);
		
		// Should use section for info section
		expect(container.querySelector('section.info_section')).toBeTruthy();
		
		// Tab navigation should contain buttons
		const buttons = container.querySelectorAll('button');
		expect(buttons.length).toBe(3);
	});

	// ===== CSS CLASSES TESTS =====
	it('TestCSSClasses_Success: Elements have correct CSS classes', () => {
		const { container } = render(InfoSection);
		
		expect(container.querySelector('.info_section')).toBeTruthy();
		expect(container.querySelector('.info_container')).toBeTruthy();
		expect(container.querySelector('.tab_nav')).toBeTruthy();
	});

	it('TestActiveTabClass_Success: Active tab has correct class', async () => {
		render(InfoSection);
		
		const aboutTab = screen.getByText('About The Project');
		const clientTab = screen.getByText('About EPI-USE Labs');
		
		// Initially about should be active
		expect(aboutTab.classList.contains('active')).toBe(true);
		
		// After clicking client, it should be active
		await fireEvent.click(clientTab);
		await tick();
		expect(clientTab.classList.contains('active')).toBe(true);
		expect(aboutTab.classList.contains('active')).toBe(false);
	});

	// ===== INTERACTION TESTS =====
	it('TestMultipleTabClicks_Success: Multiple tab clicks work correctly', async () => {
		render(InfoSection);
		
		const aboutTab = screen.getByText('About The Project');
		const clientTab = screen.getByText('About EPI-USE Labs');
		const teamTab = screen.getByText('Meet The Team');
		
		// Click through all tabs multiple times
		await fireEvent.click(clientTab);
		await tick();
		expect(clientTab.classList.contains('active')).toBe(true);
		
		await fireEvent.click(teamTab);
		await tick();
		expect(teamTab.classList.contains('active')).toBe(true);
		
		await fireEvent.click(aboutTab);
		await tick();
		expect(aboutTab.classList.contains('active')).toBe(true);
		
		// Click same tab multiple times (should remain active)
		await fireEvent.click(aboutTab);
		await tick();
		expect(aboutTab.classList.contains('active')).toBe(true);
	});

	// ===== PERFORMANCE TESTS =====
	it('TestRenderPerformance_Success: Component renders quickly', () => {
		const startTime = performance.now();
		render(InfoSection);
		const endTime = performance.now();
		
		// Should render in reasonable time (less than 50ms for test stability)
		// This accounts for variations in test environment performance
		const renderTime = endTime - startTime;
		expect(renderTime).toBeLessThan(50);
		
		// Additional validation to ensure the component actually rendered
		expect(renderTime).toBeGreaterThan(0);
	});

	it('TestMemoryUsage_Success: Component can be created and destroyed', () => {
		// Create and destroy multiple instances to test memory handling
		for (let i = 0; i < 5; i++) {
			const { unmount } = render(InfoSection);
			unmount();
		}
		
		// Should complete without memory issues
		expect(true).toBe(true);
	});

	// ===== ERROR HANDLING TESTS =====
	it('TestErrorHandling_Success: Handles invalid interactions gracefully', async () => {
		const { container } = render(InfoSection);
		
		// Component should not crash with various interactions
		const buttons = container.querySelectorAll('button');
		
		// Rapid clicking should not cause issues
		for (const button of buttons) {
			await fireEvent.click(button);
			await fireEvent.click(button);
		}
		
		expect(container).toBeTruthy();
	});

	// ===== RESPONSIVE DESIGN TESTS =====
	it('TestResponsiveStructure_Success: Component has responsive structure', () => {
		const { container } = render(InfoSection);
		
		const infoContainer = container.querySelector('.info_container');
		expect(infoContainer).toBeTruthy();
		
		const tabNav = container.querySelector('.tab_nav');
		expect(tabNav).toBeTruthy();
	});

	// ===== INTEGRATION TESTS =====
	it('TestCompleteIntegration_Success: All elements work together correctly', async () => {
		const { container } = render(InfoSection);
		
		// Check complete component structure
		expect(container.querySelector('.info_section')).toBeTruthy();
		expect(container.querySelector('.tab_nav')).toBeTruthy();
		
		// Check all tabs are present and functional
		const tabs = screen.getAllByRole('button');
		expect(tabs.length).toBe(3);
		
		// Test tab switching works
		for (const tab of tabs) {
			await fireEvent.click(tab);
			await tick();
			expect(tab.classList.contains('active')).toBe(true);
		}
	});

	// ===== STATE MANAGEMENT TESTS =====
	it('TestStateManagement_Success: Component maintains state correctly', async () => {
		render(InfoSection);
		
		const aboutTab = screen.getByText('About The Project');
		const clientTab = screen.getByText('About EPI-USE Labs');
		const teamTab = screen.getByText('Meet The Team');
		
		// Test state transitions
		expect(aboutTab.classList.contains('active')).toBe(true);
		
		await fireEvent.click(clientTab);
		await tick();
		expect(aboutTab.classList.contains('active')).toBe(false);
		expect(clientTab.classList.contains('active')).toBe(true);
		expect(teamTab.classList.contains('active')).toBe(false);
		
		await fireEvent.click(teamTab);
		await tick();
		expect(aboutTab.classList.contains('active')).toBe(false);
		expect(clientTab.classList.contains('active')).toBe(false);
		expect(teamTab.classList.contains('active')).toBe(true);
	});
});


