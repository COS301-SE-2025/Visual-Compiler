import { writable } from 'svelte/store';
import type { NodeConnection as TypesNodeConnection, NodeType } from '../types';

// Define interfaces for pipeline data
export interface Position {
    x: number;
    y: number;
}

// Use the NodeConnection from types.ts for consistency
export type NodeConnection = TypesNodeConnection;

export interface PipelineNode {
    id: string;
    type: NodeType;
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
