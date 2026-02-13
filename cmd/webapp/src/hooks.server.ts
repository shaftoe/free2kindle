import type { Handle } from '@sveltejs/kit';
import { error as kitError } from '@sveltejs/kit';
import { createApiClient, ApiError } from '$lib/server/apiClient';

export const handle: Handle = async ({ event, resolve }) => {
	if (event.url.pathname !== '/settings') {
		const apiKey = event.cookies.get('api_key');

		if (!apiKey) {
			event.locals.apiClient = null;
			return resolve(event);
		}

		try {
			event.locals.apiClient = createApiClient(apiKey);
		} catch (error) {
			event.locals.apiClient = null;
			if (error instanceof ApiError) {
				throw kitError(500, `ApiError: ${error.message}`);
			}
			throw error;
		}
	}

	return resolve(event);
};
