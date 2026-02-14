import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { requireApiKey } from '$lib/server/auth';

export const load: PageServerLoad = async ({ locals, fetch, params }) => {
	const apiClient = requireApiKey(locals);

	const id = params.id;

	const article = await apiClient.getArticle(id, fetch);
	return {
		article
	};
};

export const actions: Actions = {
	delete: async ({ locals, fetch, params }) => {
		const apiClient = requireApiKey(locals);

		const id = params.id;

		await apiClient.deleteArticle(id, fetch);
		throw redirect(303, '/');
	}
};
