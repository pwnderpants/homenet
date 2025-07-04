// Movie board interface - using shared utilities
// This file depends on utils.js being loaded first

function toggleAddForm() {
    FormUtils.toggleAddForm('movie');
}

function openEditModal(button) {
    ModalUtils.openEditModal(button, 'movie');
}

function closeEditModal() {
    ModalUtils.closeEditModal();
}

function deleteMovie(button) {
    DeleteUtils.deleteEntity(button, 'movie');
}

function pickRandomMovie() {
    RandomUtils.pickRandomEntity('movie');
}

function closeRandomMovieModal() {
    RandomUtils.closeRandomEntityModal('movie');
}

// Initialize the movie board interface
document.addEventListener('DOMContentLoaded', function() {
    Logger.info('Movie board interface initializing...');
    
    // Set up modal click outside to close
    ModalUtils.setupModalClickOutside('edit-modal', closeEditModal);
    
    // Make sure functions are available globally
    window.deleteMovie = deleteMovie;
    window.openEditModal = openEditModal;
    window.closeEditModal = closeEditModal;
    window.pickRandomMovie = pickRandomMovie;
    window.closeRandomMovieModal = closeRandomMovieModal;
    
    // Set up event delegation
    EventUtils.setupDeleteEventDelegation('movie');
    EventUtils.setupRandomModalClickOutside('movie');
    
    Logger.info('Movie board interface initialized successfully');
});

// Initialize log level
initializeLogLevel('movie_board_log_level', 'INFO'); 