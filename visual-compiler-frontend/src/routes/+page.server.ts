import { redirect } from '@sveltejs/kit';

export function load() {
  // Permanently redirect from “/” to “/login”
  throw redirect(307, '/login');
}