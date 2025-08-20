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
    lastSaved: null,
});

// Add phase completion status interface
export interface PhaseCompletionStatus {
    source: boolean;
    lexer: boolean;
    parser: boolean;
    analyser: boolean;
    translator: boolean;
}

// Create the phase completion status store
export const phase_completion_status = writable<PhaseCompletionStatus>({
    source: false,
    lexer: false,
    parser: false,
    analyser: false,
    translator: false
});

// Add a function to update phase status
export function updatePhaseStatus(phases: Partial<PhaseCompletionStatus>) {
    phase_completion_status.update(status => ({
        ...status,
        ...phases
    }));
}

// Helper functions for pipeline management
export const resetPipeline = () => {
    pipelineStore.set({
        nodes: [],
        connections: [],
        lastSaved: null,
    });
    resetPhaseStatus();
};

export const updatePipeline = (pipeline: Pipeline) => {
    pipelineStore.set(pipeline);
};

// Add function to reset phase completion status
export const resetPhaseStatus = () => {
    phase_completion_status.set({
        source: false,
        lexer: false,
        parser: false,
        analyser: false,
        translator: false
    });
};
