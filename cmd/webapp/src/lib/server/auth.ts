import { redirect } from '@sveltejs/kit';
import type { ApiClient } from './apiClient';

export function requireApiKey(locals: { apiClient: ApiClient | null }): ApiClient {
	if (!locals.apiClient) {
		redirect(303, '/settings');
	}
	return locals.apiClient;
}
