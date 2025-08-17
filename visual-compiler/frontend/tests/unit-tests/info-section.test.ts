import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import InfoSection from '../../src/lib/components/landing/info-section.svelte';

describe('InfoSection Component', () => {
	it('TestInitialRender_Success: Renders the "About" tab by default', () => {
		render(InfoSection);

		// Check that the "About" button has the 'active' class
		const aboutButton = screen.getByRole('button', { name: 'About The Project' });
		expect(aboutButton).toHaveClass('active');

		// Check that the "About" content is visible
		expect(screen.getByText('An Educational Tool for a Complex World')).toBeInTheDocument();

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

		// 2. Find and click the "Meet The Team" button.
		const teamButton = screen.getByRole('button', { name: 'Meet The Team' });
		await fireEvent.click(teamButton);

		// 3. Assert that the button has the 'active' class.
		expect(teamButton).toHaveClass('active');

  // 4. Assert that the team content is now visible by checking for a team member's name.
  expect(screen.getByText('Devan de Wet')).toBeInTheDocument();
  expect(screen.getAllByText('UI Engineer')).toHaveLength(2); // There are 2 people with this role
  expect(screen.getByText('Data Engineer')).toBeInTheDocument();		// 5. Assert that the content from the default "About" tab is now hidden.
		expect(screen.queryByText('An Educational Tool for a Complex World')).toBeNull();
	});
});
