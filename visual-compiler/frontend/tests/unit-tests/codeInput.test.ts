import { render, fireEvent, waitFor } from '@testing-library/svelte';
import test_comp from '../../src/lib/components/CodeInput.svelte';
import {vi} from 'vitest';

vi.mock('$lib/stores/toast', () => {
  return {
    addToast: vi.fn(),
    toasts: [],
  };
});
import {addToast} from '$lib/stores/toast';

describe("Testing code input", () => {
    it("Error in not .txt file", async() => {
        const {getByLabelText} = render(test_comp);
        const fileInput = getByLabelText(/upload file/i) as HTMLInputElement;

        const file = new File(['halfstack'], 'test.pdf', {type:'application/pdf'});
        await fireEvent.change(fileInput, {target: {files:[file]}});
        expect(addToast).toHaveBeenCalledWith('Only .txt files are allowed. Please upload a valid plain text file.','error');
    });

    it("Success for .txt file", async() => {
        const comp = render(test_comp);
        const fileInput = comp.getByLabelText(/upload file/i) as HTMLInputElement;

        const file = new File(['halfsatck'], 'test.txt', {type:'text/plain'});
        await fireEvent.change(fileInput, {target: {files:[file]}});
        await waitFor(() => {expect(addToast).toHaveBeenCalledWith('File uploaded successfully!','success');});
    })

    it("Empty code", async() => {
        const comp = render(test_comp);
        const c_button = comp.getByText(/confirm code/i);
        expect(c_button).toBeDisabled();
    })
})
