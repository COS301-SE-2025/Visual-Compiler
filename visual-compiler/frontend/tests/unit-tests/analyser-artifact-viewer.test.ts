import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import AnalyserArtifactViewer from '../../src/lib/components/analyser/analyser-artifact-viewer.svelte';

describe('AnalyserArtifactViewer Component', () => {
	const mockSymbolTable = [
		{ type: 'int', name: 'x', scope: 0 },
		{ type: 'int', name: 'y', scope: 1 },
		{ type: 'bool', name: 'isValid', scope: 0 }
	];

	it('TestRender_Success: Renders artifact viewer container', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: [],
				show_symbol_table: false,
				analyser_error: '',
				analyser_error_details: ''
			}
		});

		const container = document.querySelector('.artifact-viewer');
		expect(container).toBeInTheDocument();
	});

	it('TestEmptyState_Success: Shows empty state when no symbol table', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: [],
				show_symbol_table: false,
				analyser_error: '',
				analyser_error_details: ''
			}
		});

		expect(screen.getByText('Symbols will appear here after generation')).toBeInTheDocument();
	});

	it('TestSymbolTable_Success: Displays symbol table when available', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: mockSymbolTable,
				show_symbol_table: true,
				analyser_error: '',
				analyser_error_details: ''
			}
		});

		expect(screen.getByText('Symbols')).toBeInTheDocument();
		expect(screen.getByRole('table')).toBeInTheDocument();
	});

	it('TestSymbolTableHeaders_Success: Shows correct table headers', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: mockSymbolTable,
				show_symbol_table: true,
				analyser_error: '',
				analyser_error_details: ''
			}
		});

		expect(screen.getByText('Type')).toBeInTheDocument();
		expect(screen.getByText('Name')).toBeInTheDocument();
		expect(screen.getByText('Scope')).toBeInTheDocument();
	});

	it('TestSymbolTableContent_Success: Displays all symbols correctly', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: mockSymbolTable,
				show_symbol_table: true,
				analyser_error: '',
				analyser_error_details: ''
			}
		});

		// Check first symbol
		expect(screen.getByText('x')).toBeInTheDocument();
		expect(screen.getAllByText('int')).toHaveLength(2); // There are 2 int types
		expect(screen.getAllByText('0')).toHaveLength(2); // There are 2 scope 0 values

		// Check second symbol
		expect(screen.getByText('y')).toBeInTheDocument();
		expect(screen.getByText('1')).toBeInTheDocument();

		// Check third symbol
		expect(screen.getByText('isValid')).toBeInTheDocument();
		expect(screen.getByText('bool')).toBeInTheDocument();
	});

	it('TestEmptySymbolTable_Success: Shows no symbols message when table is empty', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: [],
				show_symbol_table: true,
				analyser_error: '',
				analyser_error_details: ''
			}
		});

		expect(screen.getByText('No symbols generated')).toBeInTheDocument();
	});

	it('TestErrorState_Success: Shows error when analyser error exists', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: [],
				show_symbol_table: true,
				analyser_error: 'Type mismatch error',
				analyser_error_details: 'Variable x is declared as int but assigned a string value'
			}
		});

		expect(screen.getByText('Semantic Error Found')).toBeInTheDocument();
		expect(screen.getByText(/The source code could not be analysed/)).toBeInTheDocument();
		expect(screen.getByText('Variable x is declared as int but assigned a string value')).toBeInTheDocument();
	});

	it('TestErrorIcon_Success: Shows error icon when error exists', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: [],
				show_symbol_table: true,
				analyser_error: 'Scope error',
				analyser_error_details: 'Variable not in scope'
			}
		});

		const errorIcon = document.querySelector('.error-icon svg');
		expect(errorIcon).toBeInTheDocument();
	});

	it('TestNonAnalyserPhase_Success: Does not render when phase is not analyser', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'lexer',
				symbol_table: mockSymbolTable,
				show_symbol_table: true,
				analyser_error: '',
				analyser_error_details: ''
			}
		});

		expect(screen.queryByText('Symbols')).not.toBeInTheDocument();
		expect(screen.queryByRole('table')).not.toBeInTheDocument();
	});

	it('TestTableStructure_Success: Table has proper structure', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: mockSymbolTable,
				show_symbol_table: true,
				analyser_error: '',
				analyser_error_details: ''
			}
		});

		const table = screen.getByRole('table');
		const thead = table.querySelector('thead');
		const tbody = table.querySelector('tbody');

		expect(thead).toBeInTheDocument();
		expect(tbody).toBeInTheDocument();

		// Check number of rows (header + data rows)
		const rows = table.querySelectorAll('tr');
		expect(rows).toHaveLength(4); // 1 header + 3 data rows
	});

	it('TestErrorAndSymbolTable_Success: Can show both error and symbols', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: mockSymbolTable,
				show_symbol_table: true,
				analyser_error: 'Warning: Type coercion',
				analyser_error_details: 'Implicit conversion from int to float'
			}
		});

		// Should show both symbols table and error
		expect(screen.getByText('Symbols')).toBeInTheDocument();
		expect(screen.getByRole('table')).toBeInTheDocument();
		expect(screen.getByText('Semantic Error Found')).toBeInTheDocument();
	});

	it('TestSymbolTableCSS_Success: Has correct CSS classes', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: mockSymbolTable,
				show_symbol_table: true,
				analyser_error: '',
				analyser_error_details: ''
			}
		});

		const table = document.querySelector('.symbol-table');
		const heading = document.querySelector('.symbol-heading');

		expect(table).toBeInTheDocument();
		expect(heading).toBeInTheDocument();
	});

	it('TestErrorStateStructure_Success: Error state has proper structure', () => {
		render(AnalyserArtifactViewer, {
			props: {
				phase: 'analyser',
				symbol_table: [],
				show_symbol_table: true,
				analyser_error: 'Test error',
				analyser_error_details: 'Test details'
			}
		});

		const errorState = document.querySelector('.error-state');
		const errorIcon = document.querySelector('.error-icon');
		const errorMessage = document.querySelector('.error-message');
		const errorDetails = document.querySelector('.error-details');

		expect(errorState).toBeInTheDocument();
		expect(errorIcon).toBeInTheDocument();
		expect(errorMessage).toBeInTheDocument();
		expect(errorDetails).toBeInTheDocument();
	});
});
