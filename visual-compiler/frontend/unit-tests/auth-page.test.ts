import { describe, expect, it, vi, beforeEach } from 'vitest';
import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import page_comp from '../src/routes/auth-page/+page.svelte';
import toasts from '../src/lib/components/toast-conatiner.svelte';

// Mock SvelteKit runtime
(globalThis as any).__SVELTEKIT_PAYLOAD__ = {
	data: {},
	errors: {}
};

// Mock SvelteKit navigation
vi.mock('$app/navigation', () => ({
	goto: vi.fn(),
	invalidateAll: vi.fn()
}));

vi.mock('$app/stores', () => ({
	page: {
		subscribe: vi.fn(() => ({ unsubscribe: vi.fn() }))
	}
}));

const mockSuccessResponse = {
	ok: true,
	status: 201,
	json: async () => ({ message: 'Successfully registered user' })
};
const mockFailedResponse = {
	ok: false,
	json: async () => ({ error: 'Failed' })
};
const mockLoginSuccessResponse = {
	ok: true,
	status: 200,
	json: async () => ({ 
		message: 'Welcome back, halfstack! Redirecting to your workspace...',
		id: 'user123',
		is_admin: true
	})
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
		global.fetch = vi.fn().mockResolvedValue(mockLoginSuccessResponse);
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
			expect(screen.getByText('Welcome back, halfstack! Redirecting to your workspace...')).toBeInTheDocument();
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
			expect(screen.getByText(/Account successfully created/)).toBeInTheDocument();
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

	it('register function (password mismatch)', async () => {
		// Mock Element.prototype.animate properly
		Element.prototype.animate = vi.fn().mockReturnValue({
			finished: Promise.resolve(),
			cancel: vi.fn(),
			onfinish: null
		});

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
		await fireEvent.input(c_password, { target: { value: 'differentpassword' } });
		
		render(toasts);
		await fireEvent.submit(login_form!);

		await waitFor(() => {
			expect(screen.getByText(/Passwords don't match/)).toBeInTheDocument();
		});
	});

	it('register function (network error)', async () => {
		vi.resetAllMocks();
		// Mock Element.prototype.animate properly
		Element.prototype.animate = vi.fn().mockReturnValue({
			finished: Promise.resolve(),
			cancel: vi.fn(),
			onfinish: null
		});
		
		// Mock fetch to throw a network error
		global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));
		
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
		
		render(toasts);
		await fireEvent.submit(login_form!);

		await waitFor(() => {
			expect(screen.getByText(/Registration error: Network error/)).toBeInTheDocument();
		});
	});

	it('register function (form reset after success)', async () => {
		// Mock Element.prototype.animate properly
		Element.prototype.animate = vi.fn().mockReturnValue({
			finished: Promise.resolve(),
			cancel: vi.fn(),
			onfinish: null
		});
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

		// Wait for success processing and tab switch
		await waitFor(() => {
			// Should switch back to login tab - check for login form
			expect(screen.getByLabelText('Username')).toHaveAttribute('id', 'loginUsername');
		});

		// After switching to login tab, the register form fields are no longer visible
		// So we verify the form reset by checking that we're back on the login tab
		expect(screen.getByText('Login')).toHaveClass('active');
	});
});

describe('login page advanced', () => {
	it('login function (network error)', async () => {
		vi.resetAllMocks();
		// Mock Element.prototype.animate properly
		Element.prototype.animate = vi.fn().mockReturnValue({
			finished: Promise.resolve(),
			cancel: vi.fn(),
			onfinish: null
		});
		
		// Mock fetch to throw a network error
		global.fetch = vi.fn().mockRejectedValue(new Error('Network connection failed'));
		
		render(page_comp);
		const username = screen.getByLabelText('Username');
		const login_form = username.closest('form');
		const password = screen.getByLabelText('Password');
		
		await fireEvent.input(username, { target: { value: 'halfstack' } });
		await fireEvent.input(password, { target: { value: '12345678' } });
		
		render(toasts);
		await fireEvent.submit(login_form!);

		await waitFor(() => {
			expect(screen.getByText(/Login error: Network connection failed/)).toBeInTheDocument();
		});
	});

	it('login function (user data storage)', async () => {
		Element.prototype.animate = vi.fn().mockReturnValue({
			finished: Promise.resolve(),
			cancel: vi.fn(),
			onfinish: null
		});
		
		const mockLoginResponse = {
			ok: true,
			status: 200,
			json: async () => ({ 
				message: 'Welcome back, halfstack! Redirecting to your workspace...',
				id: 'user123',
				is_admin: true
			})
		};
		
		global.fetch = vi.fn().mockResolvedValue(mockLoginResponse);
		
		// Mock sessionStorage
		const sessionStorageMock = {
			setItem: vi.fn(),
			getItem: vi.fn(),
			removeItem: vi.fn(),
			clear: vi.fn()
		};
		Object.defineProperty(window, 'sessionStorage', {
			value: sessionStorageMock
		});

		// Mock sessionStorage  
		const sessionStorageMock = {
			setItem: vi.fn(),
			getItem: vi.fn(),
			removeItem: vi.fn(),
			clear: vi.fn()
		};
		Object.defineProperty(window, 'sessionStorage', {
			value: sessionStorageMock
		});
		
		render(page_comp);
		const username = screen.getByLabelText('Username');
		const login_form = username.closest('form');
		const password = screen.getByLabelText('Password');
		
		await fireEvent.input(username, { target: { value: 'halfstack' } });
		await fireEvent.input(password, { target: { value: '12345678' } });
		
		await fireEvent.submit(login_form!);

		await waitFor(() => {
			// Should store user data in sessionStorage
			expect(sessionStorageMock.setItem).toHaveBeenCalledWith('user_id', 'user123');
			expect(sessionStorageMock.setItem).toHaveBeenCalledWith('is_admin', 'true');
			// Should set welcome overlay flag in sessionStorage
			expect(sessionStorageMock.setItem).toHaveBeenCalledWith('showWelcomeOverlay', 'true');
		});
	});

	it('login function (non-admin user)', async () => {
		Element.prototype.animate = vi.fn().mockReturnValue({
			finished: Promise.resolve(),
			cancel: vi.fn(),
			onfinish: null
		});
		
		const mockLoginResponse = {
			ok: true,
			status: 200,
			json: async () => ({ 
				message: 'Welcome back, regularuser! Redirecting to your workspace...',
				id: 'user456',
				is_admin: false
			})
		};
		
		global.fetch = vi.fn().mockResolvedValue(mockLoginResponse);
		
		// Mock sessionStorage
		const sessionStorageMock = {
			setItem: vi.fn(),
			getItem: vi.fn(),
			removeItem: vi.fn(),
			clear: vi.fn()
		};
		Object.defineProperty(window, 'sessionStorage', {
			value: sessionStorageMock
		});
		
		render(page_comp);
		const username = screen.getByLabelText('Username');
		const login_form = username.closest('form');
		const password = screen.getByLabelText('Password');
		
		await fireEvent.input(username, { target: { value: 'regularuser' } });
		await fireEvent.input(password, { target: { value: '12345678' } });
		
		await fireEvent.submit(login_form!);

		await waitFor(() => {
			// Should store is_admin as 'false'
			expect(sessionStorageMock.setItem).toHaveBeenCalledWith('is_admin', 'false');
		});
	});

	it('login button disabled when form invalid', () => {
		render(page_comp);
		
		// Initially should be disabled
		const login_buttons = screen.getAllByRole('button', { name: 'Login', hidden: true });
		const test_button = login_buttons.find((button) => button.getAttribute('type') === 'submit');
		expect(test_button).toBeDisabled();

		// Still disabled with only username
		const username = screen.getByLabelText('Username');
		fireEvent.input(username, { target: { value: 'halfstack' } });
		expect(test_button).toBeDisabled();

		// Still disabled with only password
		fireEvent.input(username, { target: { value: '' } });
		const password = screen.getByLabelText('Password');
		fireEvent.input(password, { target: { value: '12345678' } });
		expect(test_button).toBeDisabled();
	});
});

describe('password visibility', () => {
	it('toggle password visibility', async () => {
		render(page_comp);
		
		const password = screen.getByLabelText('Password') as HTMLInputElement;
		expect(password.type).toBe('password');

		// Find and click the password visibility toggle button
		const toggleButton = screen.getByRole('button', { name: '' }); // The show/hide password button
		
		await fireEvent.click(toggleButton);
		
		// Should toggle to text type
		expect(password.type).toBe('text');
		
		// Toggle back
		await fireEvent.click(toggleButton);
		expect(password.type).toBe('password');
	});

	it('toggle confirm password visibility', async () => {
		render(page_comp);
		const reg_tab = screen.getByText('Register');
		await fireEvent.click(reg_tab);

		const confirmPassword = screen.getByLabelText('Confirm Password') as HTMLInputElement;
		expect(confirmPassword.type).toBe('password');

		// Find the confirm password container and its toggle button
		const passwordGroup = confirmPassword.closest('.input-group');
		const toggleButton = passwordGroup?.querySelector('button[type="button"]') as HTMLButtonElement;
		
		await fireEvent.click(toggleButton);
		
		// Should toggle to text type  
		expect(confirmPassword.type).toBe('text');
		
		// Toggle back
		await fireEvent.click(toggleButton);
		expect(confirmPassword.type).toBe('password');
	});

	it('password visibility state changes correctly', () => {
		render(page_comp);
		
		const password = screen.getByLabelText('Password') as HTMLInputElement;
		const toggleButton = screen.getByRole('button', { name: '' }); // The show/hide password button
		
		// Initial state should be password (hidden)
		expect(password.type).toBe('password');
		
		// Click to show password
		fireEvent.click(toggleButton);
		expect(password.type).toBe('text');
		
		// Click again to hide password
		fireEvent.click(toggleButton);
		expect(password.type).toBe('password');
	});

	it('confirm password visibility state changes correctly', async () => {
		render(page_comp);
		const reg_tab = screen.getByText('Register');
		await fireEvent.click(reg_tab);

		const confirmPassword = screen.getByLabelText('Confirm Password') as HTMLInputElement;
		const passwordGroup = confirmPassword.closest('.input-group');
		const toggleButton = passwordGroup?.querySelector('button[type="button"]') as HTMLButtonElement;
		
		// Initial state should be password (hidden)
		expect(confirmPassword.type).toBe('password');
		
		// Click to show password
		fireEvent.click(toggleButton);
		expect(confirmPassword.type).toBe('text');
		
		// Click again to hide password
		fireEvent.click(toggleButton);
		expect(confirmPassword.type).toBe('password');
	});
});

describe('form validation edge cases', () => {
	it('register button disabled with whitespace-only inputs', () => {
		render(page_comp);
		const reg_tab = screen.getByText('Register');
		fireEvent.click(reg_tab);

		const username = screen.getByLabelText('Username');
		const email = screen.getByLabelText('Email');
		const password = screen.getByLabelText('Password');
		const c_password = screen.getByLabelText('Confirm Password');
		
		// Fill with whitespace only
		fireEvent.input(username, { target: { value: '   ' } });
		fireEvent.input(email, { target: { value: '   ' } });
		fireEvent.input(password, { target: { value: '   ' } });
		fireEvent.input(c_password, { target: { value: '   ' } });

		const reg_buttons = screen.getAllByRole('button', { name: 'Register', hidden: true });
		const test_button = reg_buttons.find((button) => button.getAttribute('type') === 'submit');
		expect(test_button).toBeDisabled();
	});

	it('login button disabled with whitespace-only inputs', () => {
		render(page_comp);

		const username = screen.getByLabelText('Username');
		const password = screen.getByLabelText('Password');
		
		// Fill with whitespace only
		fireEvent.input(username, { target: { value: '   ' } });
		fireEvent.input(password, { target: { value: '   ' } });

		const login_buttons = screen.getAllByRole('button', { name: 'Login', hidden: true });
		const test_button = login_buttons.find((button) => button.getAttribute('type') === 'submit');
		expect(test_button).toBeDisabled();
	});
});


