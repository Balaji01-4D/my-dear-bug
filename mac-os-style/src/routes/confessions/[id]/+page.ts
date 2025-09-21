import api from '$lib/api.js';
import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
    try {
        const confessionId = parseInt(params.id);
        
        if (isNaN(confessionId) || confessionId <= 0) {
            throw error(400, 'Invalid confession ID');
        }
        
        const confession = await api.confessions.getById(confessionId);
        
        return {
            confession
        };
    } catch (err: any) {
        console.error('Failed to load confession:', err);
        
        if (err.message?.includes('404')) {
            throw error(404, 'Confession not found');
        }
        
        throw error(500, 'Failed to load confession');
    }
};
