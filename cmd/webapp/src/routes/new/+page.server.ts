import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ locals, request }) => {
		const apiClient = locals.apiClient;

		if (!apiClient) {
			return fail(400, { error: 'api key is required' });
		}

		const data = await request.formData();
		const url = data.get('url');

		if (!url || typeof url !== 'string' || url.trim() === '') {
			return fail(400, { error: 'url is required' });
		}

		await apiClient.createArticle({ url: url.trim() }, fetch);
		redirect(303, '/articles');
	}
};
