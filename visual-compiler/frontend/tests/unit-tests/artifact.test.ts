import { render, screen } from '@testing-library/svelte';
import ArtifactViewer from '../../src/lib/components/ArtifactViewer.svelte';

describe('Test start artefact', () => {
    it('Empty', () => {
      render(ArtifactViewer,{ phase:'lexer'});
      expect(screen.getByText('Tokens will appear here after generation')).toBeInTheDocument();
    });

    it('Empty still', () => {
      render(ArtifactViewer, {phase:'parser'});
      expect(screen.queryByText('Tokens will appear here after generation')).toBeNull();
    });
});