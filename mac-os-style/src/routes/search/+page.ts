import type { PageLoad } from './$types';

export const load: PageLoad = ({ url }) => {
    const searchParams = url.searchParams;
    
    return {
        initialQuery: searchParams.get('q') || '',
    };
};
