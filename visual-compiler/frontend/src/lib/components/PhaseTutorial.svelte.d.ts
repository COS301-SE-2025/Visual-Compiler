import { SvelteComponentTyped } from 'svelte';

export default class PhaseTutorial extends SvelteComponentTyped<
  { phase: string },       // the props
  { close: void },         // the events
  {}                       // slots
> {}
