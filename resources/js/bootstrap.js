/**
 * We'll load the axios HTTP library which allows us to easily issue requests
 * to our Go Goravel back-end.
 */

import axios from 'axios';
window.axios = axios;

window.axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';

/**
 * Configure axios base URL for API requests
 */
window.axios.defaults.baseURL = '/api';

/**
 * Set up response interceptors for better error handling
 */
window.axios.interceptors.response.use(
    response => response,
    error => {
        if (error.response?.status === 401) {
            // Handle unauthorized
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

/**
 * Kafka Event Listener Configuration
 * WebSocket connection for real-time activity updates
 * Can be extended with WebSocket/SSE implementation
 */

// Example: WebSocket setup for real-time notifications
// window.addEventListener('activity:created', (event) => {
//     console.log('Activity created:', event.detail);
// });

// window.addEventListener('activity:updated', (event) => {
//     console.log('Activity updated:', event.detail);
// });

// window.addEventListener('activity:deleted', (event) => {
//     console.log('Activity deleted:', event.detail);
// });

/**
 * Echo exposes an expressive API for subscribing to channels and listening
 * for events that are broadcast by Laravel/Broadcasting services.
 * For Goravel, consider using WebSocket or Server-Sent Events (SSE).
 */

// import Echo from 'laravel-echo';
// import Pusher from 'pusher-js';
// window.Pusher = Pusher;

// window.Echo = new Echo({
//     broadcaster: 'pusher',
//     key: import.meta.env.VITE_PUSHER_APP_KEY,
//     cluster: import.meta.env.VITE_PUSHER_APP_CLUSTER ?? 'mt1',
//     wsHost: import.meta.env.VITE_PUSHER_HOST ? import.meta.env.VITE_PUSHER_HOST : `ws-${import.meta.env.VITE_PUSHER_APP_CLUSTER}.pusher.com`,
//     wsPort: import.meta.env.VITE_PUSHER_PORT ?? 80,
//     wssPort: import.meta.env.VITE_PUSHER_PORT ?? 443,
//     forceTLS: (import.meta.env.VITE_PUSHER_SCHEME ?? 'https') === 'https',
//     enabledTransports: ['ws', 'wss'],
// });
