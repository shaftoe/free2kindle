import { error as kitError, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, fetch }) => {
	const apiClient = locals.apiClient;

	if (!apiClient) {
		redirect(303, '/settings');
	}

	try {
		const response = await apiClient.getArticles(fetch);
		return {
			articles: response.articles,
			total: response.total,
			page: response.page,
			pageSize: response.pageSize,
			hasMore: response.hasMore
		};
	} catch (err) {
		if (err instanceof Error) {
			throw kitError(500, `failed to load articles: ${err.message}`);
		} else {
			throw kitError(500, 'failed to load articles');
		}
	}
};
