import { writable } from 'svelte/store';

// Define interfaces for pipeline data
export interface Position {
    x: number;
    y: number;
}

export interface NodeConnection {
    sourceNodeId: string;
    targetNodeId: string;
}

export interface PipelineNode {
    id: string;
    type: string;
    label: string;
    position: Position;
}

export interface Pipeline {
    nodes: PipelineNode[];
    connections: NodeConnection[];
    lastSaved: string | null;
}

// Create the pipeline store with initial empty state
export const pipelineStore = writable<Pipeline>({
    nodes: [],
    connections: [],
    lastSaved: null
});

// Helper functions for pipeline management
export const resetPipeline = () => {
    pipelineStore.set({
        nodes: [],
        connections: [],
        lastSaved: null
    });
};

export const updatePipeline = (pipeline: Pipeline) => {
    pipelineStore.set(pipeline);
};
