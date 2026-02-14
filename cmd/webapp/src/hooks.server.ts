import type { Handle } from '@sveltejs/kit';
import { error as kitError } from '@sveltejs/kit';
import { createApiClient, ApiError } from '$lib/server/apiClient';

const HEALTH_CACHE_TTL_SUCCESS = 5 * 60 * 1000;
const HEALTH_CACHE_TTL_FAILURE = 30 * 1000;

interface HealthCache {
	status: 'healthy' | 'unhealthy';
	timestamp: number;
	error?: string;
}

let healthCache: HealthCache | null = null;
let healthCheckPromise: Promise<void> | null = null;

async function getHealthStatus(fetch: typeof globalThis.fetch): Promise<void> {
	const now = Date.now();

	if (healthCache) {
		const ttl =
			healthCache.status === 'healthy' ? HEALTH_CACHE_TTL_SUCCESS : HEALTH_CACHE_TTL_FAILURE;
		if (now - healthCache.timestamp < ttl) {
			if (healthCache.status === 'unhealthy') {
				throw new ApiError(503, healthCache.error || 'Backend health check failed');
			}
			return;
		}
	}

	if (!healthCheckPromise) {
		healthCheckPromise = (async () => {
			try {
				const client = createApiClient('placeholder');
				await client.healthCheck(fetch);
				healthCache = { status: 'healthy', timestamp: Date.now() };
			} catch (error) {
				const errorMsg = error instanceof ApiError ? error.message : 'Backend health check failed';
				healthCache = { status: 'unhealthy', timestamp: Date.now(), error: errorMsg };
				throw error;
			} finally {
				healthCheckPromise = null;
			}
		})();
	}

	await healthCheckPromise;
}

export const handle: Handle = async ({ event, resolve }) => {
	if (event.url.pathname !== '/settings') {
		const apiKey = event.cookies.get('api_key');

		if (!apiKey) {
			event.locals.apiClient = null;
			return resolve(event);
		}

		try {
			event.locals.apiClient = createApiClient(apiKey);
			await getHealthStatus(event.fetch);
		} catch (error) {
			event.locals.apiClient = null;
			if (error instanceof ApiError) {
				throw kitError(error.status, `ApiError: ${error.message}`);
			}
			throw error;
		}
	}

	return resolve(event);
};
