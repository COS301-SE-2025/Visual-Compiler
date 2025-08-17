import { writable } from 'svelte/store';

// Define the pipeline store interface
export interface Pipeline {
    nodes: any[];
    lastSaved: string | null;
}

// Create the pipeline store with initial empty state
export const pipelineStore = writable<Pipeline>({
    nodes: [],
    lastSaved: null
});

// Helper functions for pipeline management
export const resetPipeline = () => {
    pipelineStore.set({
        nodes: [],
        lastSaved: null
    });
};

export const updatePipeline = (pipeline: Pipeline) => {
    pipelineStore.set(pipeline);
};
