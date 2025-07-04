// Shared utilities for the homenet application
// Structured logging system for browser console
const Logger = {
    levels: {
        DEBUG: 0,
        INFO: 1,
        WARN: 2,
        ERROR: 3
    },
    
    currentLevel: 1, // Default to INFO level
    
    setLevel(level) {
        this.currentLevel = level;
        this.info(`Log level set to ${Object.keys(this.levels)[level]}`);
    },
    
    shouldLog(level) {
        return level >= this.currentLevel;
    },
    
    formatMessage(level, message, ...args) {
        const timestamp = new Date().toISOString();
        const levelNames = ['DEBUG', 'INFO', 'WARN', 'ERROR'];
        const levelName = levelNames[level];
        const prefix = `[${timestamp}] ${levelName}:`;
        
        if (args.length > 0) {
            return [prefix, message, ...args];
        }
        return [prefix, message];
    },
    
    debug(message, ...args) {
        if (this.shouldLog(this.levels.DEBUG)) {
            console.log(...this.formatMessage(this.levels.DEBUG, message, ...args));
        }
    },
    
    info(message, ...args) {
        if (this.shouldLog(this.levels.INFO)) {
            console.log(...this.formatMessage(this.levels.INFO, message, ...args));
        }
    },
    
    warn(message, ...args) {
        if (this.shouldLog(this.levels.WARN)) {
            console.warn(...this.formatMessage(this.levels.WARN, message, ...args));
        }
    },
    
    error(message, ...args) {
        if (this.shouldLog(this.levels.ERROR)) {
            console.error(...this.formatMessage(this.levels.ERROR, message, ...args));
        }
    }
};

// Modal management utilities
const ModalUtils = {
    openEditModal(button, entityType) {
        const entityId = button.getAttribute(`data-${entityType}-id`);
        const entityTitle = button.getAttribute(`data-${entityType}-title`);
        Logger.info(`Opening edit modal for ${entityType}:`, entityTitle, 'ID:', entityId);
        
        const modal = document.getElementById('edit-modal');
        const fields = ['year', 'genre', 'streaming', 'notes', 'imdb'];
        
        // Populate common form fields
        document.getElementById(`edit-${entityType}-id`).value = entityId;
        document.getElementById('edit-title').value = entityTitle;
        
        fields.forEach(field => {
            const value = button.getAttribute(`data-${entityType}-${field}`);
            const element = document.getElementById(`edit-${field}`);
            if (element) {
                element.value = value || '';
            }
        });
        
        // Handle entity-specific fields
        if (entityType === 'movie') {
            const movieAvailableNow = button.getAttribute('data-movie-available-now');
            document.getElementById('edit-available-now').checked = movieAvailableNow === 'true';
        } else if (entityType === 'tvshow') {
            const tvshowActiveSeason = button.getAttribute('data-tvshow-active-season');
            document.getElementById('edit-active-season').checked = tvshowActiveSeason === 'true';
        }
        
        Logger.debug('Edit form populated with data');
        
        // Show modal
        modal.classList.remove('hidden');
        modal.classList.add('animate-fade-in');
    },
    
    closeEditModal() {
        Logger.debug('Closing edit modal');
        const modal = document.getElementById('edit-modal');
        modal.classList.add('hidden');
        modal.classList.remove('animate-fade-in');
    },
    
    setupModalClickOutside(modalId, closeFunction) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.addEventListener('click', function(e) {
                if (e.target === modal) {
                    closeFunction();
                }
            });
        }
    }
};

// Form toggle utilities
const FormUtils = {
    toggleAddForm(entityType) {
        Logger.info(`Toggling add ${entityType} form`);
        const form = document.getElementById(`add-${entityType}-form`);
        const toggleText = document.getElementById('toggle-text');
        const toggleButton = document.getElementById('toggle-add-form');
        const formElement = form.querySelector('form');
        
        if (form.classList.contains('hidden')) {
            // Show form
            Logger.debug(`Showing add ${entityType} form`);
            form.classList.remove('hidden');
            form.classList.add('animate-fade-in');
            toggleText.textContent = 'Hide Form';
            toggleButton.classList.remove('bg-blue-600', 'hover:bg-blue-700');
            toggleButton.classList.add('bg-gray-600', 'hover:bg-gray-700');
        } else {
            // Hide form and clear it
            Logger.debug(`Hiding add ${entityType} form`);
            form.classList.add('hidden');
            form.classList.remove('animate-fade-in');
            toggleText.textContent = `Add New ${entityType.charAt(0).toUpperCase() + entityType.slice(1)}`;
            toggleButton.classList.remove('bg-gray-600', 'hover:bg-gray-700');
            toggleButton.classList.add('bg-blue-600', 'hover:bg-blue-700');
            
            // Clear form fields
            if (formElement) {
                formElement.reset();
                Logger.debug('Form fields cleared');
            }
        }
    }
};

// Delete utilities
const DeleteUtils = {
    deleteEntity(button, entityType) {
        const entityId = button.getAttribute(`data-${entityType}-id`);
        const entityTitle = button.closest('.bg-gray-700').querySelector('h4').textContent;
        
        Logger.warn(`Delete ${entityType} requested:`, entityTitle, 'ID:', entityId);
        
        if (confirm(`Are you sure you want to delete this ${entityType}?`)) {
            Logger.info('Delete confirmed, sending delete request');
            
            fetch(`/${entityType}-board/delete/${entityId}`, {
                method: 'DELETE',
            })
            .then(response => {
                if (response.ok) {
                    Logger.info(`${entityType} deleted successfully from server`);
                    
                    // Remove the entity card from DOM
                    const entityCard = button.closest('.bg-gray-700');
                    if (entityCard) {
                        entityCard.remove();
                        Logger.debug(`${entityType} card removed from DOM`);
                    }
                    
                    // Update entity count
                    const countElement = document.getElementById(`${entityType}-count`);
                    const currentCount = parseInt(countElement.textContent.match(/\d+/)[0]);
                    const newCount = currentCount - 1;
                    
                    if (newCount === 0) {
                        Logger.info(`No ${entityType}s remaining, showing empty state`);
                        // Show "no entities" message
                        const listElement = document.getElementById(`${entityType}-list`);
                        listElement.innerHTML = `
                            <div class="text-center py-8">
                                <svg class="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m-9 0h10m-10 0a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V6a2 2 0 00-2-2"></path>
                                </svg>
                                <p class="text-gray-400">No ${entityType}s added yet. Add your first ${entityType} above!</p>
                            </div>
                        `;
                        countElement.textContent = `0 ${entityType}s in your list`;
                    } else {
                        countElement.textContent = `${newCount} ${entityType}s in your list`;
                        Logger.debug(`${entityType} count updated to:`, newCount);
                    }
                    
                } else {
                    Logger.error(`Failed to delete ${entityType} from server, status:`, response.status);
                    alert(`Failed to delete ${entityType}`);
                }
            })
            .catch(error => {
                Logger.error('Delete request failed:', error.message);
                alert(`Failed to delete ${entityType}: ` + error.message);
            });
        } else {
            Logger.info('Delete cancelled by user');
        }
    }
};

// Random entity utilities
const RandomUtils = {
    pickRandomEntity(entityType) {
        Logger.info(`Requesting random ${entityType}`);
        
        fetch(`/${entityType}-board/random`)
            .then(response => {
                Logger.debug(`Random ${entityType} response received, status:`, response.status);
                return response.text();
            })
            .then(html => {
                Logger.debug(`Random ${entityType} HTML received, length:`, html.length);
                document.getElementById(`random-${entityType}-modal-content`).innerHTML = html;
                document.getElementById(`random-${entityType}-modal`).classList.remove('hidden');
                document.getElementById(`random-${entityType}-modal`).classList.add('animate-fade-in');
                Logger.info(`Random ${entityType} modal displayed`);
            })
            .catch(error => {
                Logger.error(`Failed to fetch random ${entityType}:`, error.message);
                document.getElementById(`random-${entityType}-modal-content`).innerHTML = `<div class="text-center text-red-400">Failed to fetch a random ${entityType}.</div>`;
                document.getElementById(`random-${entityType}-modal`).classList.remove('hidden');
            });
    },
    
    closeRandomEntityModal(entityType) {
        Logger.debug(`Closing random ${entityType} modal`);
        const modal = document.getElementById(`random-${entityType}-modal`);
        modal.classList.add('hidden');
        modal.classList.remove('animate-fade-in');
    }
};

// Event delegation utilities
const EventUtils = {
    setupDeleteEventDelegation(entityType) {
        document.addEventListener('click', function(e) {
            if (e.target.closest(`.delete-${entityType}-btn`)) {
                const button = e.target.closest(`.delete-${entityType}-btn`);
                DeleteUtils.deleteEntity(button, entityType);
            }
        });
    },
    
    setupRandomModalClickOutside(entityType) {
        const modal = document.getElementById(`random-${entityType}-modal`);
        if (modal) {
            modal.addEventListener('click', function(e) {
                if (e.target === modal) {
                    RandomUtils.closeRandomEntityModal(entityType);
                }
            });
        }
    }
};

// Log level initialization
function initializeLogLevel(storageKey, defaultLevel = 'INFO') {
    const urlParams = new URLSearchParams(window.location.search);
    const logLevel = urlParams.get('log') || localStorage.getItem(storageKey) || defaultLevel;
    
    const levels = { 'DEBUG': 0, 'INFO': 1, 'WARN': 2, 'ERROR': 3 };
    if (levels[logLevel] !== undefined) {
        Logger.setLevel(levels[logLevel]);
        localStorage.setItem(storageKey, logLevel);
    }
}

// Export utilities for use in other files
window.Logger = Logger;
window.ModalUtils = ModalUtils;
window.FormUtils = FormUtils;
window.DeleteUtils = DeleteUtils;
window.RandomUtils = RandomUtils;
window.EventUtils = EventUtils;
window.initializeLogLevel = initializeLogLevel; 