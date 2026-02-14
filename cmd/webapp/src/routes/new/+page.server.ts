import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { requireApiKey } from '$lib/server/auth';

export const load: PageServerLoad = async ({ locals }) => {
	requireApiKey(locals);
	return {};
};

export const actions: Actions = {
	default: async ({ locals, request, fetch }) => {
		const apiClient = requireApiKey(locals);

		const data = await request.formData();
		const url = data.get('url');

		if (!url || typeof url !== 'string' || url.trim() === '') {
			return { error: 'url is required' };
		}

		await apiClient.createArticle({ url: url.trim() }, fetch);
		redirect(303, '/articles');
	}
};
