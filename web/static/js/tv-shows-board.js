// TV Shows board interface - using shared utilities
// This file depends on utils.js being loaded first

function toggleAddForm() {
    FormUtils.toggleAddForm('tvshow');
}

function openEditModal(button) {
    ModalUtils.openEditModal(button, 'tvshow');
}

function closeEditModal() {
    ModalUtils.closeEditModal();
}

function deleteTVShow(button) {
    DeleteUtils.deleteEntity(button, 'tvshow');
}

// Initialize the TV Shows board interface
document.addEventListener('DOMContentLoaded', function() {
    Logger.info('TV Shows board interface initializing...');
    
    // Set up modal click outside to close
    ModalUtils.setupModalClickOutside('edit-modal', closeEditModal);
    
    // Make sure functions are available globally
    window.deleteTVShow = deleteTVShow;
    window.openEditModal = openEditModal;
    window.closeEditModal = closeEditModal;
    
    // Set up event delegation
    EventUtils.setupDeleteEventDelegation('tvshow');
    
    Logger.info('TV Shows board interface initialized successfully');
});

// Initialize log level
initializeLogLevel('tv_shows_board_log_level', 'INFO'); 