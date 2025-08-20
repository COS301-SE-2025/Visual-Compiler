import { describe, it, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { 
    pipelineStore, 
    resetPipeline, 
    updatePipeline,
    type Pipeline,
    type PipelineNode,
    type NodeConnection,
    type Position
} from '../src/lib/stores/pipeline';

describe('Pipeline Store', () => {
    beforeEach(() => {
        // Reset pipeline store before each test
        resetPipeline();
    });

    describe('pipelineStore', () => {
        it('should have initial empty state', () => {
            const initialState = get(pipelineStore);
            expect(initialState).toEqual({
                nodes: [],
                connections: [],
                lastSaved: null
            });
        });

        it('should support subscription to changes', () => {
            const states: Pipeline[] = [];
            
            const unsubscribe = pipelineStore.subscribe((pipeline) => {
                states.push({ ...pipeline });
            });

            const testPipeline: Pipeline = {
                nodes: [
                    {
                        id: 'source-1',
                        type: 'source',
                        label: 'Source Code',
                        position: { x: 100, y: 100 }
                    }
                ],
                connections: [],
                lastSaved: '2023-01-01T10:00:00Z'
            };

            pipelineStore.set(testPipeline);

            expect(states).toHaveLength(2); // Initial state + updated state
            expect(states[1]).toEqual(testPipeline);
            
            unsubscribe();
        });

        it('should handle complex pipeline state', () => {
            const complexPipeline: Pipeline = {
                nodes: [
                    {
                        id: 'source-1',
                        type: 'source',
                        label: 'Source Code',
                        position: { x: 100, y: 100 }
                    },
                    {
                        id: 'lexer-1',
                        type: 'lexer',
                        label: 'Lexical Analyzer',
                        position: { x: 300, y: 100 }
                    },
                    {
                        id: 'parser-1',
                        type: 'parser',
                        label: 'Parser',
                        position: { x: 500, y: 100 }
                    }
                ],
                connections: [
                    {
                        id: 'conn-1',
                        sourceNodeId: 'source-1',
                        targetNodeId: 'lexer-1',
                        sourceType: 'source',
                        targetType: 'lexer',
                        sourceAnchor: 'bottom',
                        targetAnchor: 'top'
                    },
                    {
                        id: 'conn-2',
                        sourceNodeId: 'lexer-1',
                        targetNodeId: 'parser-1',
                        sourceType: 'lexer',
                        targetType: 'parser',
                        sourceAnchor: 'output',
                        targetAnchor: 'input'
                    }
                ],
                lastSaved: '2023-01-01T12:00:00Z'
            };

            pipelineStore.set(complexPipeline);
            const currentState = get(pipelineStore);

            expect(currentState).toEqual(complexPipeline);
            expect(currentState.nodes).toHaveLength(3);
            expect(currentState.connections).toHaveLength(2);
        });
    });

    describe('resetPipeline function', () => {
        it('should reset pipeline to initial empty state', () => {
            // First set some data
            const testPipeline: Pipeline = {
                nodes: [
                    {
                        id: 'source-1',
                        type: 'source',
                        label: 'Source Code',
                        position: { x: 100, y: 100 }
                    }
                ],
                connections: [
                    {
                        id: 'conn-1',
                        sourceNodeId: 'source-1',
                        targetNodeId: 'lexer-1',
                        sourceType: 'source',
                        targetType: 'lexer',
                        sourceAnchor: 'bottom',
                        targetAnchor: 'top'
                    }
                ],
                lastSaved: '2023-01-01T10:00:00Z'
            };

            pipelineStore.set(testPipeline);
            
            // Verify data is set
            let currentState = get(pipelineStore);
            expect(currentState.nodes).toHaveLength(1);
            expect(currentState.connections).toHaveLength(1);
            expect(currentState.lastSaved).toBe('2023-01-01T10:00:00Z');

            // Reset and verify
            resetPipeline();
            currentState = get(pipelineStore);
            
            expect(currentState).toEqual({
                nodes: [],
                connections: [],
                lastSaved: null
            });
        });

        it('should not affect empty pipeline', () => {
            const initialState = get(pipelineStore);
            resetPipeline();
            const afterResetState = get(pipelineStore);

            expect(initialState).toEqual(afterResetState);
            expect(afterResetState).toEqual({
                nodes: [],
                connections: [],
                lastSaved: null
            });
        });
    });

    describe('updatePipeline function', () => {
        it('should update pipeline with new data', () => {
            const newPipeline: Pipeline = {
                nodes: [
                    {
                        id: 'lexer-1',
                        type: 'lexer',
                        label: 'Lexical Analyzer',
                        position: { x: 200, y: 200 }
                    }
                ],
                connections: [],
                lastSaved: '2023-01-02T10:00:00Z'
            };

            updatePipeline(newPipeline);
            const currentState = get(pipelineStore);

            expect(currentState).toEqual(newPipeline);
        });

        it('should completely replace existing pipeline data', () => {
            // Set initial data
            const initialPipeline: Pipeline = {
                nodes: [
                    {
                        id: 'source-1',
                        type: 'source',
                        label: 'Source Code',
                        position: { x: 100, y: 100 }
                    }
                ],
                connections: [
                    {
                        id: 'conn-1',
                        sourceNodeId: 'source-1',
                        targetNodeId: 'lexer-1',
                        sourceType: 'source',
                        targetType: 'lexer',
                        sourceAnchor: 'bottom',
                        targetAnchor: 'top'
                    }
                ],
                lastSaved: '2023-01-01T10:00:00Z'
            };

            pipelineStore.set(initialPipeline);

            // Update with new data
            const newPipeline: Pipeline = {
                nodes: [
                    {
                        id: 'parser-1',
                        type: 'parser',
                        label: 'Parser',
                        position: { x: 300, y: 300 }
                    }
                ],
                connections: [],
                lastSaved: '2023-01-02T15:00:00Z'
            };

            updatePipeline(newPipeline);
            const currentState = get(pipelineStore);

            expect(currentState).toEqual(newPipeline);
            expect(currentState.nodes).toHaveLength(1);
            expect(currentState.nodes[0].id).toBe('parser-1');
            expect(currentState.connections).toHaveLength(0);
            expect(currentState.lastSaved).toBe('2023-01-02T15:00:00Z');
        });

        it('should handle empty pipeline update', () => {
            // Set some initial data
            const initialPipeline: Pipeline = {
                nodes: [
                    {
                        id: 'source-1',
                        type: 'source',
                        label: 'Source Code',
                        position: { x: 100, y: 100 }
                    }
                ],
                connections: [],
                lastSaved: '2023-01-01T10:00:00Z'
            };

            pipelineStore.set(initialPipeline);

            // Update with empty pipeline
            const emptyPipeline: Pipeline = {
                nodes: [],
                connections: [],
                lastSaved: null
            };

            updatePipeline(emptyPipeline);
            const currentState = get(pipelineStore);

            expect(currentState).toEqual(emptyPipeline);
        });
    });

    describe('Type interfaces', () => {
        it('should support Position interface', () => {
            const position: Position = { x: 100, y: 200 };
            
            expect(position).toHaveProperty('x');
            expect(position).toHaveProperty('y');
            expect(typeof position.x).toBe('number');
            expect(typeof position.y).toBe('number');
        });

        it('should support PipelineNode interface', () => {
            const node: PipelineNode = {
                id: 'test-node-1',
                type: 'analyser',
                label: 'Semantic Analyzer',
                position: { x: 400, y: 300 }
            };

            expect(node).toHaveProperty('id');
            expect(node).toHaveProperty('type');
            expect(node).toHaveProperty('label');
            expect(node).toHaveProperty('position');
            expect(typeof node.id).toBe('string');
            expect(typeof node.label).toBe('string');
            expect(node.position).toHaveProperty('x');
            expect(node.position).toHaveProperty('y');
        });

        it('should support NodeConnection interface', () => {
            const connection: NodeConnection = {
                id: 'test-conn-1',
                sourceNodeId: 'source-1',
                targetNodeId: 'lexer-1',
                sourceType: 'source',
                targetType: 'lexer',
                sourceAnchor: 'bottom',
                targetAnchor: 'top',
            };

            expect(connection).toHaveProperty('id');
            expect(connection).toHaveProperty('sourceNodeId');
            expect(connection).toHaveProperty('targetNodeId');
            expect(connection).toHaveProperty('sourceType');
            expect(connection).toHaveProperty('targetType');
            expect(connection).toHaveProperty('sourceAnchor');
            expect(connection).toHaveProperty('targetAnchor');
        });

        it('should support Pipeline interface', () => {
            const pipeline: Pipeline = {
                nodes: [
                    {
                        id: 'translator-1',
                        type: 'translator',
                        label: 'Code Generator',
                        position: { x: 600, y: 100 }
                    }
                ],
                connections: [
                    {
                        id: 'conn-trans-1',
                        sourceNodeId: 'analyser-1',
                        targetNodeId: 'translator-1',
                        sourceType: 'analyser',
                        targetType: 'translator',
                        sourceAnchor: 'bottom',
                        targetAnchor: 'top'
                    }
                ],
                lastSaved: '2023-01-03T14:30:00Z'
            };

            expect(pipeline).toHaveProperty('nodes');
            expect(pipeline).toHaveProperty('connections');
            expect(pipeline).toHaveProperty('lastSaved');
            expect(Array.isArray(pipeline.nodes)).toBe(true);
            expect(Array.isArray(pipeline.connections)).toBe(true);
        });
    });
});
