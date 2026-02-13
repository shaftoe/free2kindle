import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ cookies, request }) => {
		const data = await request.formData();
		const apiKey = data.get('apiKey');

		if (!apiKey || typeof apiKey !== 'string' || apiKey.trim() === '') {
			return fail(400, { error: 'api key is required' });
		}

		cookies.set('api_key', apiKey.trim(), {
			path: '/',
			httpOnly: true,
			secure: false,
			sameSite: 'lax',
			maxAge: 60 * 60 * 24 * 365
		});

		redirect(303, '/');
	}
};
