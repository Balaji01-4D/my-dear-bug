import { confessionApi } from '$lib/api';

export async function load() {
    try {
        const confessions = await confessionApi.getTrendingWeekly();
        return {
            confessions,
            title: 'Trending This Week'
        };
    } catch (error) {
        console.error('Error loading trending weekly confessions:', error);
        return {
            confessions: [],
            title: 'Trending This Week',
            error: 'Failed to load trending confessions'
        };
    }
}
