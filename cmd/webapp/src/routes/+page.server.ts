import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, fetch, url }) => {
	const apiClient = locals.apiClient;

	if (!apiClient) {
		redirect(303, '/settings');
	}

	const pageParam = url.searchParams.get('page');
	const pageSizeParam = url.searchParams.get('page_size');

	const page = pageParam ? parseInt(pageParam, 10) : 1;
	const pageSize = pageSizeParam ? parseInt(pageSizeParam, 10) : 10;

	try {
		const response = await apiClient.getArticles(fetch, page, pageSize);
		return {
			articles: response.articles,
			total: response.total,
			page: response.page,
			page_size: response.page_size,
			has_more: response.has_more,
			error: undefined
		};
	} catch (err) {
		return {
			articles: [],
			total: 0,
			page: 1,
			page_size: 10,
			has_more: false,
			error: err instanceof Error ? err.message : 'failed to load articles'
		};
	}
};

export const actions: Actions = {};
