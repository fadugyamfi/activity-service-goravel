/**
 * Activity Service Frontend Application
 * Main entry point for loading resources and dependencies
 */

import './bootstrap';

// Log that the application has loaded
console.log('Activity Service Frontend Loaded');


/**
 * Activity Service Frontend Application
 * 
 * This module initializes the frontend for the Activity Service built with Goravel.
 * It sets up API communication, event handling, and user interface initialization.
 */

/**
 * Initialize activity management functionality
 */
function initializeActivityManagement() {
    const app = document.getElementById('app');
    if (!app) return;

    // API endpoints
    const API_BASE = '/api';
    const ACTIVITIES_ENDPOINT = `${API_BASE}/activities`;

    // DOM elements
    const activityForm = document.getElementById('activity-form');
    const activityList = document.getElementById('activity-list');
    const searchInput = document.getElementById('search-input');
    const filterStatus = document.getElementById('filter-status');
    const filterType = document.getElementById('filter-type');

    /**
     * Fetch and display activities
     */
    async function loadActivities(filters = {}) {
        try {
            let url = ACTIVITIES_ENDPOINT;
            const params = new URLSearchParams();

            if (filters.search) params.append('search', filters.search);
            if (filters.status) params.append('status', filters.status);
            if (filters.type) params.append('type', filters.type);
            if (filters.page) params.append('page', filters.page);

            if (params.toString()) {
                url += '?' + params.toString();
            }

            const response = await window.axios.get(url);
            renderActivities(response.data.data);
        } catch (error) {
            console.error('Failed to load activities:', error);
            showNotification('Failed to load activities', 'error');
        }
    }

    /**
     * Render activities list
     */
    function renderActivities(activities) {
        if (!activityList) return;

        if (!activities || activities.length === 0) {
            activityList.innerHTML = '<p class="text-gray-500">No activities found.</p>';
            return;
        }

        activityList.innerHTML = activities.map(activity => `
            <div class="activity-item border rounded p-4 mb-4 hover:shadow-lg transition">
                <div class="flex justify-between items-start">
                    <div class="flex-1">
                        <h3 class="font-bold text-lg">${escapeHtml(activity.name)}</h3>
                        <p class="text-gray-600">${escapeHtml(activity.description || '')}</p>
                        <div class="mt-2 flex gap-2">
                            <span class="badge badge-primary">${escapeHtml(activity.type)}</span>
                            <span class="badge badge-${activity.status === 'active' ? 'success' : 'secondary'}">
                                ${escapeHtml(activity.status)}
                            </span>
                        </div>
                    </div>
                    <div class="flex gap-2">
                        <button onclick="editActivity(${activity.id})" class="btn btn-sm btn-info">Edit</button>
                        <button onclick="deleteActivity(${activity.id})" class="btn btn-sm btn-danger">Delete</button>
                    </div>
                </div>
                <div class="text-xs text-gray-400 mt-2">
                    Created: ${new Date(activity.created_at).toLocaleString()}
                </div>
            </div>
        `).join('');
    }

    /**
     * Create a new activity
     */
    async function createActivity(formData) {
        try {
            const response = await window.axios.post(ACTIVITIES_ENDPOINT, formData);
            showNotification('Activity created successfully', 'success');
            activityForm.reset();
            loadActivities();
        } catch (error) {
            console.error('Failed to create activity:', error);
            showNotification('Failed to create activity', 'error');
        }
    }

    /**
     * Update an activity
     */
    async function updateActivity(id, formData) {
        try {
            const response = await window.axios.put(`${ACTIVITIES_ENDPOINT}/${id}`, formData);
            showNotification('Activity updated successfully', 'success');
            activityForm.reset();
            loadActivities();
        } catch (error) {
            console.error('Failed to update activity:', error);
            showNotification('Failed to update activity', 'error');
        }
    }

    /**
     * Delete an activity
     */
    async function deleteActivity(id) {
        if (!confirm('Are you sure you want to delete this activity?')) return;

        try {
            await window.axios.delete(`${ACTIVITIES_ENDPOINT}/${id}`);
            showNotification('Activity deleted successfully', 'success');
            loadActivities();
        } catch (error) {
            console.error('Failed to delete activity:', error);
            showNotification('Failed to delete activity', 'error');
        }
    }

    /**
     * Edit an activity
     */
    async function editActivity(id) {
        try {
            const response = await window.axios.get(`${ACTIVITIES_ENDPOINT}/${id}`);
            const activity = response.data.data;

            // Populate form (implement based on your form structure)
            if (activityForm) {
                activityForm.querySelector('[name="name"]').value = activity.name;
                activityForm.querySelector('[name="description"]').value = activity.description || '';
                activityForm.querySelector('[name="type"]').value = activity.type;
                activityForm.querySelector('[name="status"]').value = activity.status;
                activityForm.dataset.editId = id;
            }
        } catch (error) {
            console.error('Failed to load activity:', error);
            showNotification('Failed to load activity', 'error');
        }
    }

    /**
     * Show notification to user
     */
    function showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `alert alert-${type} mb-4`;
        notification.textContent = message;

        const container = document.querySelector('.notification-container') || app;
        container.insertBefore(notification, container.firstChild);

        setTimeout(() => notification.remove(), 3000);
    }

    /**
     * Escape HTML to prevent XSS
     */
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    /**
     * Event listeners
     */
    if (activityForm) {
        activityForm.addEventListener('submit', (e) => {
            e.preventDefault();
            const formData = new FormData(activityForm);
            const data = Object.fromEntries(formData);
            const editId = activityForm.dataset.editId;

            if (editId) {
                updateActivity(editId, data);
                delete activityForm.dataset.editId;
            } else {
                createActivity(data);
            }
        });
    }

    if (searchInput) {
        searchInput.addEventListener('input', (e) => {
            loadActivities({ search: e.target.value });
        });
    }

    if (filterStatus) {
        filterStatus.addEventListener('change', (e) => {
            loadActivities({ status: e.target.value });
        });
    }

    if (filterType) {
        filterType.addEventListener('change', (e) => {
            loadActivities({ type: e.target.value });
        });
    }

    /**
     * Listen for Kafka events (if using WebSocket or SSE)
     */
    window.addEventListener('activity:created', (event) => {
        console.log('Activity created via Kafka:', event.detail);
        showNotification('A new activity was created', 'info');
        loadActivities();
    });

    window.addEventListener('activity:updated', (event) => {
        console.log('Activity updated via Kafka:', event.detail);
        showNotification('An activity was updated', 'info');
        loadActivities();
    });

    window.addEventListener('activity:deleted', (event) => {
        console.log('Activity deleted via Kafka:', event.detail);
        showNotification('An activity was deleted', 'info');
        loadActivities();
    });

    // Export functions to global scope
    window.editActivity = editActivity;
    window.deleteActivity = deleteActivity;

    // Initial load
    loadActivities();
}

// Initialize when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeActivityManagement);
} else {
    initializeActivityManagement();
}

