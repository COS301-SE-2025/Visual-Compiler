import { describe, expect, it, vi } from 'vitest';
import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import page_comp from '../../src/routes/auth-page/+page.svelte';
import toasts from '../../src/lib/components/toast-conatiner.svelte';

const mockSuccessResponse = {
	ok: true,
	json: async () => ({ message: 'Success' })
};
const mockFailedResponse = {
	ok: false,
	json: async () => ({ error: 'Failed' })
};

describe('page start up', () => {
	it('render login by default', () => {
		render(page_comp);
		expect(screen.getByLabelText('Username')).toBeInTheDocument();
		const u_label = screen.getByText('Username');
		expect(u_label).toHaveAttribute('for', 'loginUsername');
		expect(screen.getByLabelText('Password')).toBeInTheDocument();
		const p_label = screen.getByText('Password');
		expect(p_label).toHaveAttribute('for', 'loginPassword');
	});
});

describe('login page', () => {
	it('button disabled', () => {
		render(page_comp);
		const login_buttons = screen.getAllByRole('button', { name: 'Login', hidden: true });
		const test_button = login_buttons.find((button) => button.getAttribute('type') === 'submit');
		expect(test_button).toBeDisabled();
	});
	it('button enabled', () => {
		render(page_comp);
		const username = screen.getByLabelText('Username');
		const password = screen.getByLabelText('Password');
		fireEvent.input(username, { target: { value: 'halfstack' } });
		fireEvent.input(password, { target: { value: '12345678' } });

		const login_buttons = screen.getAllByRole('button', { name: 'Login', hidden: true });
		const test_button = login_buttons.find((button) => button.getAttribute('type') === 'submit');
		expect(test_button).toBeEnabled();
	});

	it('login function (valid)', async () => {
		Element.prototype.animate =
			Element.prototype.animate ||
			(() => ({
				finished: Promise.resolve(),
				cancel: () => {}
			}));
		global.fetch = vi.fn().mockResolvedValue(mockSuccessResponse);
		render(page_comp);
		render(toasts);
		const username = screen.getByLabelText('Username');
		const login_form = username.closest('form');
		const password = screen.getByLabelText('Password');
		await fireEvent.input(username, { target: { value: 'halfstack' } });
		await fireEvent.input(password, { target: { value: '12345678' } });
		await fireEvent.submit(login_form!);

		expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/users/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				login: 'halfstack',
				password: '12345678'
			})
		});
		await waitFor(() => {
			expect(screen.getByText('Login successful!')).toBeInTheDocument();
		});
	});

	it('login function (invalid)', async () => {
		vi.resetAllMocks();
		Element.prototype.animate =
			Element.prototype.animate ||
			(() => ({
				finished: Promise.resolve(),
				cancel: () => {}
			}));
		global.fetch = vi.fn().mockResolvedValue(mockFailedResponse);
		render(page_comp);
		const username = screen.getByLabelText('Username');
		const login_form = username.closest('form');
		const password = screen.getByLabelText('Password');
		await fireEvent.input(username, { target: { value: 'halfstack' } });
		await fireEvent.input(password, { target: { value: '12345678' } });
		await fireEvent.submit(login_form!);

		expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/users/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				login: 'halfstack',
				password: '12345678'
			})
		});

		render(toasts);
		await waitFor(() => {
			expect(screen.getByText('Login failed: Failed')).toBeInTheDocument();
		});
	});
});

describe('register page', () => {
	it('render register', async () => {
		render(page_comp);
		const reg_tab = screen.getByText('Register');
		await fireEvent.click(reg_tab);

		expect(screen.getByLabelText('Username')).toBeInTheDocument();
		const u_label = screen.getByText('Username');
		expect(u_label).toHaveAttribute('for', 'regUsername');
		expect(screen.getByLabelText('Email')).toBeInTheDocument();
		const e_label = screen.getByText('Email');
		expect(e_label).toHaveAttribute('for', 'regEmail');
		expect(screen.getByLabelText('Password')).toBeInTheDocument();
		const p_label = screen.getByText('Password');
		expect(p_label).toHaveAttribute('for', 'regPassword');
		expect(screen.getByLabelText('Confirm Password')).toBeInTheDocument();
		const c_label = screen.getByText('Confirm Password');
		expect(c_label).toHaveAttribute('for', 'regConfirmPassword');
	});

	it('button disabled', () => {
		render(page_comp);
		const reg_tab = screen.getByText('Register');
		fireEvent.click(reg_tab);

		const reg_buttons = screen.getAllByRole('button', { name: 'Register', hidden: true });
		const test_button = reg_buttons.find((button) => button.getAttribute('type') === 'submit');
		expect(test_button).toBeDisabled();
	});
	it('button enabled', () => {
		render(page_comp);
		const reg_tab = screen.getByText('Register');
		fireEvent.click(reg_tab);

		const username = screen.getByLabelText('Username');
		const email = screen.getByLabelText('Email');
		const password = screen.getByLabelText('Password');
		const c_password = screen.getByLabelText('Confirm Password');
		fireEvent.input(username, { target: { value: 'halfstack' } });
		fireEvent.input(email, { target: { value: 'halfstack@gmail.com' } });
		fireEvent.input(password, { target: { value: '12345678' } });
		fireEvent.input(c_password, { target: { value: '12345678' } });

		const reg_buttons = screen.getAllByRole('button', { name: 'Register', hidden: true });
		const test_button = reg_buttons.find((button) => button.getAttribute('type') === 'submit');
		expect(test_button).toBeEnabled();
	});

	it('register function (valid)', async () => {
		Element.prototype.animate =
			Element.prototype.animate ||
			(() => ({
				finished: Promise.resolve(),
				cancel: () => {}
			}));
		global.fetch = vi.fn().mockResolvedValue(mockSuccessResponse);
		render(page_comp);
		const reg_tab = screen.getByText('Register');
		await fireEvent.click(reg_tab);

		const username = screen.getByLabelText('Username');
		const login_form = username.closest('form');
		const email = screen.getByLabelText('Email');
		const password = screen.getByLabelText('Password');
		const c_password = screen.getByLabelText('Confirm Password');
		await fireEvent.input(username, { target: { value: 'halfstack' } });
		await fireEvent.input(email, { target: { value: 'halfstack@gmail.com' } });
		await fireEvent.input(password, { target: { value: '12345678' } });
		await fireEvent.input(c_password, { target: { value: '12345678' } });
		await fireEvent.submit(login_form!);

		expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/users/register', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				email: 'halfstack@gmail.com',
				username: 'halfstack',
				password: '12345678'
			})
		});
		render(toasts);
		await waitFor(() => {
			expect(screen.getByText('Success')).toBeInTheDocument();
		});
	});

	it('register function (invalid)', async () => {
		vi.resetAllMocks();
		Element.prototype.animate =
			Element.prototype.animate ||
			(() => ({
				finished: Promise.resolve(),
				cancel: () => {}
			}));
		global.fetch = vi.fn().mockResolvedValue(mockFailedResponse);
		render(page_comp);
		const reg_tab = screen.getByText('Register');
		await fireEvent.click(reg_tab);

		const username = screen.getByLabelText('Username');
		const login_form = username.closest('form');
		const email = screen.getByLabelText('Email');
		const password = screen.getByLabelText('Password');
		const c_password = screen.getByLabelText('Confirm Password');
		await fireEvent.input(username, { target: { value: 'halfstack' } });
		await fireEvent.input(email, { target: { value: 'halfstack@gmail.com' } });
		await fireEvent.input(password, { target: { value: '12345678' } });
		await fireEvent.input(c_password, { target: { value: '12345678' } });
		await fireEvent.submit(login_form!);

		expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/users/register', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				email: 'halfstack@gmail.com',
				username: 'halfstack',
				password: '12345678'
			})
		});

		render(toasts);
		await waitFor(() => {
			expect(screen.getByText('Registration failed: Failed')).toBeInTheDocument();
		});
	});
});
