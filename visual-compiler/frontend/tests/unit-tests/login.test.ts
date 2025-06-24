import {describe,expect,vi } from 'vitest';
import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen,waitFor } from '@testing-library/svelte';
import page_comp from '../../src/routes/auth-page/+page.svelte';
