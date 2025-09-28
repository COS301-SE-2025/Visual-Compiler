import { writable } from 'svelte/store';

export interface TutorialState {
    showCanvasTutorial: boolean;
    hasSeenCanvasTutorial: boolean;
}

// Create the tutorial store
export const tutorialStore = writable<TutorialState>({
    showCanvasTutorial: false,
    hasSeenCanvasTutorial: false
});

// Function to show the canvas tutorial
export function showCanvasTutorial() {
    tutorialStore.update(state => ({
        ...state,
        showCanvasTutorial: true
    }));
}

// Function to hide the canvas tutorial
export function hideCanvasTutorial() {
    tutorialStore.update(state => ({
        ...state,
        showCanvasTutorial: false
    }));
}

// Function to mark canvas tutorial as seen
export function markCanvasTutorialAsSeen() {
    sessionStorage.setItem('hasSeenCanvasTutorial', 'true');
    tutorialStore.update(state => ({
        ...state,
        showCanvasTutorial: false,
        hasSeenCanvasTutorial: true
    }));
}

// Function to check if user has seen the tutorial
export function checkTutorialStatus() {
    const hasSeenTutorial = sessionStorage.getItem('hasSeenCanvasTutorial') === 'true';
    tutorialStore.update(state => ({
        ...state,
        hasSeenCanvasTutorial: hasSeenTutorial
    }));
    return hasSeenTutorial;
}

// Function to reset tutorial (for testing or admin purposes)
export function resetTutorialStatus() {
    sessionStorage.removeItem('hasSeenCanvasTutorial');
    tutorialStore.update(state => ({
        ...state,
        showCanvasTutorial: false,
        hasSeenCanvasTutorial: false
    }));
}

// Function to trigger tutorial for new projects
export function triggerTutorialForNewProject() {
    const hasSeenTutorial = checkTutorialStatus();
    if (!hasSeenTutorial) {
        showCanvasTutorial();
        return true;
    }
    return false;
}
