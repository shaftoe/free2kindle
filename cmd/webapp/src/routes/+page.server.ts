import { redirect } from '@sveltejs/kit';
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
			hasMore: response.hasMore,
			error: undefined
		};
	} catch (err) {
		return {
			articles: [],
			total: 0,
			page: 1,
			pageSize: 10,
			hasMore: false,
			error: err instanceof Error ? err.message : 'failed to load articles'
		};
	}
};
