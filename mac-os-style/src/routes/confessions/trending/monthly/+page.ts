import { confessionApi } from '$lib/api';

export async function load() {
    try {
        const confessions = await confessionApi.getTrendingMonthly();
        return {
            confessions,
            title: 'Trending This Month'
        };
    } catch (error) {
        console.error('Error loading trending monthly confessions:', error);
        return {
            confessions: [],
            title: 'Trending This Month',
            error: 'Failed to load trending confessions'
        };
    }
}
