import { error as kitError, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

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
