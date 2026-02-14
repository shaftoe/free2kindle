import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, fetch, params }) => {
	const apiClient = locals.apiClient;

	if (!apiClient) {
		redirect(303, '/settings');
	}

	const id = params.id;

	const article = await apiClient.getArticle(id, fetch);
	return {
		article
	};
};

export const actions: Actions = {
	delete: async ({ locals, fetch, params }) => {
		const apiClient = locals.apiClient;

		if (!apiClient) {
			throw fail(401, { error: 'api key is required' });
		}

		const id = params.id;

		await apiClient.deleteArticle(id, fetch);
		throw redirect(303, '/');
	}
};
