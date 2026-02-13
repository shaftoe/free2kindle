import type { Handle } from '@sveltejs/kit';
import { createApiClient } from '$lib/server/apiClient';

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
			console.error('error creating api client:', error);
			event.locals.apiClient = null;
		}
	}

	return resolve(event);
};
