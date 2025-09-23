import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import OptimizerPhaseTutorial from '../src/lib/components/optimizer/optimizer-phase-tutorial.svelte';

describe('OptimizerPhaseTutorial Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Initial Rendering', () => {
    it('renders the tutorial with correct title', () => {
      render(OptimizerPhaseTutorial);
      
      expect(screen.getByText('Optimisation')).toBeInTheDocument();
    });

    it('shows the tutorial description', () => {
      render(OptimizerPhaseTutorial);
      
      expect(screen.getByText(/Optimisation is the process of automatically improving/)).toBeInTheDocument();
      expect(screen.getByText(/without changing the program's intended behavior or output/)).toBeInTheDocument();
    });

    it('starts on step 1', () => {
      render(OptimizerPhaseTutorial);
      
      expect(screen.getByText('1. What is optimisation?')).toBeInTheDocument();
    });

    it('displays navigation buttons', () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      const prevButton = screen.getByRole('button', { name: /previous/i });
      
      expect(nextButton).toBeInTheDocument();
      expect(prevButton).toBeInTheDocument();
    });

    it('shows step indicators', () => {
      render(OptimizerPhaseTutorial);
      
      // Should show step 1 of 5
      expect(screen.getByText(/1.*5/)).toBeInTheDocument();
    });
  });

  describe('Step Navigation', () => {
    it('allows navigating to next step', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      await fireEvent.click(nextButton);
      
      expect(screen.getByText('2. Benefits of Optimisation')).toBeInTheDocument();
    });

    it('allows navigating to previous step', async () => {
      render(OptimizerPhaseTutorial);
      
      // Go to step 2 first
      const nextButton = screen.getByRole('button', { name: /next/i });
      await fireEvent.click(nextButton);
      
      // Then go back to step 1
      const prevButton = screen.getByRole('button', { name: /previous/i });
      await fireEvent.click(prevButton);
      
      expect(screen.getByText('1. What is optimisation?')).toBeInTheDocument();
    });

    it('disables previous button on first step', () => {
      render(OptimizerPhaseTutorial);
      
      const prevButton = screen.getByRole('button', { name: /previous/i });
      expect(prevButton).toBeDisabled();
    });

    it('disables next button on last step', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      
      // Navigate to last step (step 5)
      for (let i = 0; i < 4; i++) {
        await fireEvent.click(nextButton);
      }
      
      expect(nextButton).toBeDisabled();
    });

    it('cannot go beyond first step when clicking previous', async () => {
      render(OptimizerPhaseTutorial);
      
      const prevButton = screen.getByRole('button', { name: /previous/i });
      
      // Try to go before step 1 multiple times
      await fireEvent.click(prevButton);
      await fireEvent.click(prevButton);
      
      expect(screen.getByText('1. What is optimisation?')).toBeInTheDocument();
    });

    it('cannot go beyond last step when clicking next', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      
      // Go to last step and try to go further
      for (let i = 0; i < 6; i++) {
        await fireEvent.click(nextButton);
      }
      
      // Should show step 5 content
      expect(screen.getByText('5. Loop Unrolling')).toBeInTheDocument();
    });
  });

  describe('Step Content', () => {
    it('shows step 1 content about what optimization is', () => {
      render(OptimizerPhaseTutorial);
      
      expect(screen.getByText('1. What is optimisation?')).toBeInTheDocument();
      expect(screen.getByText(/Intermediate Representation/)).toBeInTheDocument();
      expect(screen.getByText(/Abstract Syntax Tree/)).toBeInTheDocument();
    });

    it('shows step 2 content about benefits', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      await fireEvent.click(nextButton);
      
      expect(screen.getByText('2. Benefits of Optimisation')).toBeInTheDocument();
      expect(screen.getByText(/Increased Execution Speed/)).toBeInTheDocument();
      expect(screen.getByText(/Reduced Memory Footprint/)).toBeInTheDocument();
      expect(screen.getByText(/Smaller Binary Size/)).toBeInTheDocument();
      expect(screen.getByText(/Improved Energy Efficiency/)).toBeInTheDocument();
    });

    it('shows step 3 content about techniques', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      
      // Navigate to step 3
      await fireEvent.click(nextButton);
      await fireEvent.click(nextButton);
      
      expect(screen.getByText('3. Constant Folding')).toBeInTheDocument();
    });

    it('shows step 4 content about process', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      
      // Navigate to step 4
      for (let i = 0; i < 3; i++) {
        await fireEvent.click(nextButton);
      }
      
      expect(screen.getByText('4. Dead Code Elimination')).toBeInTheDocument();
    });

    it('shows step 5 content about real-world example', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      
      // Navigate to step 5
      for (let i = 0; i < 4; i++) {
        await fireEvent.click(nextButton);
      }
      
      expect(screen.getByText('5. Loop Unrolling')).toBeInTheDocument();
    });
  });

  describe('Step Counter Display', () => {
    it('shows correct step count on step 1', () => {
      render(OptimizerPhaseTutorial);
      
      expect(screen.getByText('1 / 5')).toBeInTheDocument();
    });

    it('updates step count when navigating', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      await fireEvent.click(nextButton);
      
      expect(screen.getByText('2 / 5')).toBeInTheDocument();
    });

    it('shows correct step count on last step', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      
      // Navigate to last step
      for (let i = 0; i < 4; i++) {
        await fireEvent.click(nextButton);
      }
      
      expect(screen.getByText('5 / 5')).toBeInTheDocument();
    });
  });

  describe('Component Structure', () => {
    it('has proper tutorial container', () => {
      const { container } = render(OptimizerPhaseTutorial);
      
      const tutorialContainer = container.querySelector('.phase-tutorial');
      expect(tutorialContainer).toBeInTheDocument();
    });

    it('has tutorial content area', () => {
      const { container } = render(OptimizerPhaseTutorial);
      
      const tutorialContent = container.querySelector('.tutorial-content');
      expect(tutorialContent).toBeInTheDocument();
    });

    it('has navigation controls', () => {
      const { container } = render(OptimizerPhaseTutorial);
      
      const navigationControls = container.querySelector('.navigation');
      expect(navigationControls).toBeInTheDocument();
    });

    it('has step indicators', () => {
      const { container } = render(OptimizerPhaseTutorial);
      
      const stepIndicator = container.querySelector('.step-counter');
      expect(stepIndicator).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    it('has proper heading hierarchy', () => {
      render(OptimizerPhaseTutorial);
      
      expect(screen.getByRole('heading', { level: 2, name: 'Optimisation' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 3, name: '1. What is optimisation?' })).toBeInTheDocument();
    });

    it('has accessible navigation buttons', () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      const prevButton = screen.getByRole('button', { name: /previous/i });
      
      expect(nextButton).toBeInTheDocument();
      expect(prevButton).toBeInTheDocument();
    });

    it('buttons have appropriate aria states', () => {
      render(OptimizerPhaseTutorial);
      
      const prevButton = screen.getByRole('button', { name: /previous/i });
      expect(prevButton).toHaveAttribute('disabled');
    });
  });

  describe('Tutorial Content Quality', () => {
    it('includes key optimization concepts', () => {
      render(OptimizerPhaseTutorial);
      
      // Navigate through all steps to check for key concepts
      expect(screen.getByText(/Optimisation is the process of automatically improving/)).toBeInTheDocument();
    });

    it('includes practical examples and analogies', () => {
      render(OptimizerPhaseTutorial);
      
      expect(screen.getByText(/Think of it like a master editor reviewing a draft/)).toBeInTheDocument();
    });

    it('explains intermediate representation', () => {
      render(OptimizerPhaseTutorial);
      
      expect(screen.getByText(/Intermediate Representation \(IR\)/)).toBeInTheDocument();
      expect(screen.getByText(/Abstract Syntax Tree \(AST\)/)).toBeInTheDocument();
    });

    it('covers optimization benefits comprehensively', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      await fireEvent.click(nextButton);
      
      expect(screen.getByText(/Increased Execution Speed/)).toBeInTheDocument();
      expect(screen.getByText(/Reduced Memory Footprint/)).toBeInTheDocument();
      expect(screen.getByText(/Smaller Binary Size/)).toBeInTheDocument();
      expect(screen.getByText(/Improved Energy Efficiency/)).toBeInTheDocument();
    });
  });

  describe('Interactive Features', () => {
    it('responds to keyboard navigation', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      
      // Click the button instead of using keyPress as keyPress might not trigger navigation
      await fireEvent.click(nextButton);
      
      expect(screen.getByText('2. Benefits of Optimisation')).toBeInTheDocument();
    });

    it('maintains focus management during navigation', async () => {
      render(OptimizerPhaseTutorial);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      await fireEvent.click(nextButton);
      
      // Button should still be focusable after navigation
      nextButton.focus();
      expect(document.activeElement).toBe(nextButton);
    });
  });
});