import { error as kitError, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, fetch, params }) => {
	const apiClient = locals.apiClient;

	if (!apiClient) {
		redirect(303, '/settings');
	}

	const id = params.id;

	try {
		const article = await apiClient.getArticle(id, fetch);
		return {
			article
		};
	} catch (err) {
		if (err instanceof Error) {
			throw kitError(500, `failed to load article: ${err.message}`);
		} else {
			throw kitError(500, 'failed to load article');
		}
	}
};

export const actions: Actions = {
	delete: async ({ locals, fetch, params }) => {
		const apiClient = locals.apiClient;

		if (!apiClient) {
			return fail(401, { error: 'api key is required' });
		}

		const id = params.id;

		try {
			await apiClient.deleteArticle(id, fetch);
			redirect(303, '/');
		} catch (err) {
			return fail(500, { error: err instanceof Error ? err.message : 'failed to delete article' });
		}
	}
};
