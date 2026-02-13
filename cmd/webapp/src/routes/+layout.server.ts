import { error } from '@sveltejs/kit';
import { ApiError } from '$lib/server/apiClient';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	if (!locals.apiClient) {
		return {};
	}

	try {
		await locals.apiClient.healthCheck();
		return {};
	} catch (e) {
		if (e instanceof ApiError) {
			throw error(503, `Backend health check failed (${e.status}): ${e.message}`);
		}
		throw error(503, 'Backend health check failed');
	}
};
